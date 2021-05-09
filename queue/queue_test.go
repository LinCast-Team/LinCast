package queue

import (
	"lincast/database"
	"lincast/podcasts"
	"os"
	"runtime"
	"sync"
	"testing"

	"github.com/joomcode/errorx"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type QueueTestSuite struct {
	length      int
	dbPath      string
	dbFilename  string
	sampleFeeds []string

	suite.Suite
}

func (s *QueueTestSuite) SetupTest() {
	s.length = runtime.NumCPU()
	s.dbPath = "./test_queue"
	s.dbFilename = "queue_test.sqlite"
	s.sampleFeeds = []string{
		"https://changelog.com/gotime/feed",
		"https://feeds.emilcar.fm/daily",
		"https://www.ivoox.com/podcast-despeja-x-by-xataka_fg_f1579492_filtro_1.xml",
		"https://www.ivoox.com/podcast-tortulia-podcast-episodios_fg_f1157653_filtro_1.xml",
	}

	err := os.Mkdir(s.dbPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (s *QueueTestSuite) BeforeTest(_, _ string) {}

func (s *QueueTestSuite) TestNewUpdateQueue() {
	assert := assert2.New(s.T())

	db, err := database.New(s.dbPath, "new_update_queue_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	q, err := NewUpdateQueue(db, s.length)

	assert.NoError(err, "the queue should be created correctly")
	assert.NotNil(q, "the returned Queue instance shouldn't be nil")

	q, err = NewUpdateQueue(db, -1)

	if assert.Error(err, "if the argument that corresponds to the length of the queue is a negative"+
		" number or 0, an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.IllegalArgument), "an error of type "+
			"errorx.IllegalArgument should be returned")
	}
	assert.Nil(q, "if the argument is incorrect, the instance of Queue returned should be nil")
}

func (s *QueueTestSuite) TestWorker() {
	assert := assert2.New(s.T())

	db, err := database.New(s.dbPath, "worker_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	q := UpdateQueue{
		dbInstance: db,
		q:          make(chan Job),
	}

	var wg sync.WaitGroup

	for i, feed := range s.sampleFeeds {
		wg.Add(i)

		go func(i int, feed string, wg *sync.WaitGroup) {
			defer wg.Done()

			// Try to process a podcast
			p, err := podcasts.GetPodcast(feed)
			if err != nil {
				panic(errorx.Decorate(err, "error when trying to get podcast"))
			}

			err = db.InsertPodcast(p)
			if err != nil {
				panic(errorx.Decorate(err, "error when trying to store podcast"))
			}

			p, err = db.GetPodcastByID(i + 1)
			if err != nil {
				panic(errorx.Decorate(err, "error when trying to get the podcast with ID %d"+
					" from the database", i+1))
			}

			job := NewJob(p)

			go q.worker(i)

			q.Send(job)

			// Wait until the job has been reported as done.
			<-job.Done

			// Check if the episodes of the processed podcast has been saved in the database.
			eps, err := db.GetEpisodesByPodcast(i + 1)
			if err != nil {
				panic(errorx.Decorate(err, "error when trying to get episodes of podcast with ID 0"))
			}

			assert.True(len(*eps) > 0, "the episodes of the processed podcast should be stored in the database")
		}(i, feed, &wg)
	}
}

func (s *QueueTestSuite) TestNewJob() {
	assert := assert2.New(s.T())

	db, err := database.New(s.dbPath, "new_job_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	p, err := podcasts.GetPodcast(s.sampleFeeds[0])
	if err != nil {
		panic(errorx.Decorate(err, "error when trying to get the podcast"))
	}

	job := NewJob(p)

	if assert.NotNil(job, "the returned Job should not be nil") {
		assert.NotNil(job.Done, "Job.Done should not be nil")
	}
}

func (s *QueueTestSuite) AfterTest(_, _ string) {}

func (s *QueueTestSuite) TearDownTest() {
	err := os.RemoveAll(s.dbPath)
	if err != nil {
		panic(err)
	}
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}

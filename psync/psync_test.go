package psync

import (
	"fmt"
	"os"
	"testing"
	"time"

	"lincast/database"
	lTesting "lincast/utils/testing"

	"github.com/joomcode/errorx"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SynchronizerTestSuite struct {
	dbPath       string
	dbFilename   string
	dbInstance1  *database.Database
	dbInstance2  *database.Database
	dbInstanceQ1 *database.Database
	dbInstanceQ2 *database.Database
	dbInstanceQ3 *database.Database
	dbInstanceQ4 *database.Database
	dbInstanceQ5 *database.Database

	suite.Suite
}

func (s *SynchronizerTestSuite) SetupTest() {
	s.dbPath = "./test_synchronizer"
	s.dbFilename = "test.sqlite"

	err := os.Mkdir(s.dbPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	db, err := database.New(s.dbPath, s.dbFilename)
	if err != nil {
		panic(err)
	}

	s.dbInstance1 = db

	db2, err := database.New(s.dbPath, "1_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	s.dbInstance2 = db2

	db3, err := database.New(s.dbPath, "q1_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	s.dbInstanceQ1 = db3

	db4, err := database.New(s.dbPath, "q2_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	s.dbInstanceQ2 = db4

	db5, err := database.New(s.dbPath, "q3_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	s.dbInstanceQ3 = db5

	db6, err := database.New(s.dbPath, "q4_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	s.dbInstanceQ4 = db6

	db7, err := database.New(s.dbPath, "q5_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	s.dbInstanceQ5 = db7
}

func (s *SynchronizerTestSuite) BeforeTest(_, _ string) {}

func (s *SynchronizerTestSuite) TestNew() {
	assert := assert2.New(s.T())

	pSync, err := New(s.dbInstance1)

	assert.NoError(err, "a new instance of the Synchronizer should be returned without errors")
	assert.NotNil(pSync, "the instance returned shouldn't be nil")

	pSync, err = New(nil)

	if assert.Error(err, "if the used instance of the database is nil, an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.IllegalState), "the returned error should be of type"+
			" errorx.IllegalState")
	}
	assert.Nil(pSync, "if the instance of the database used is nil, the returned instance of the"+
		" Synchronizer should be nil too")
}

func (s *SynchronizerTestSuite) TestUpdateProgress() {
	assert := assert2.New(s.T())

	pSync, err := New(s.dbInstance2)
	if err != nil {
		panic(err)
	}

	newProgress := time.Second * 100
	episodeGUID := "abc"
	podcastID := 2

	err = pSync.UpdateProgress(newProgress, episodeGUID, podcastID)

	assert.NoError(err, "the progress should be updated without errors")

	p := pSync.GetProgress()
	if err != nil {
		panic(err)
	}

	assert.Equal(newProgress, p.Progress, "the progress stored in the database should be the same as"+
		" the one introduced")
	assert.Equal(episodeGUID, p.EpisodeGUID, "the episodeGUID stored in the database should be the same as"+
		" the one introduced")
	assert.Equal(podcastID, p.PodcastID, "the podcastID stored in the database should be the same as"+
		" the one introduced")

	err = pSync.UpdateProgress(newProgress, "", podcastID)

	if assert.Error(err, "if the episodeGUID used is an empty string, an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.IllegalArgument), "the returned error should be of"+
			" type errorx.IllegalArgument")
	}
}

func (s *SynchronizerTestSuite) TestGetProgress() {
	assert := assert2.New(s.T())

	pSync, err := New(s.dbInstance2)
	if err != nil {
		panic(err)
	}

	progress := time.Second * 134
	episodeGUID := "1234"
	podcastID := 8

	err = pSync.UpdateProgress(progress, episodeGUID, podcastID)
	if err != nil {
		panic(err)
	}

	p := pSync.GetProgress()

	assert.Equal(progress, p.Progress, "the progress should be obtained correctly")
	assert.Equal(episodeGUID, p.EpisodeGUID, "the episodeGUID should be obtained correctly")
	assert.Equal(podcastID, p.PodcastID, "the podcastID should be obtained correctly")
}

func (s *SynchronizerTestSuite) TestGetQueue() {
	assert := assert2.New(s.T())

	eps, err := s.insertRandomQueueEps(s.dbInstanceQ1)
	if err != nil {
		panic(err)
	}

	pSync, err := New(s.dbInstanceQ1)
	if err != nil {
		panic(err)
	}

	q := pSync.GetQueue()

	if assert.NotNil(q, "the queue returned should not be nil") {
		assert.Equal(*eps, q.Content, "the returned queue should be the same as the inserted one")
	}
}

func (s *SynchronizerTestSuite) TestSetQueue() {
	assert := assert2.New(s.T())

	eps1, err := s.insertRandomQueueEps(s.dbInstanceQ2)
	if err != nil {
		panic(err)
	}

	pSync, err := New(s.dbInstanceQ2)
	if err != nil {
		panic(err)
	}

	// This will generate the episodes (without IDs) to be inserted in the database.
	eps2 := s.generateDummyQueueEp(10, false)

	// Method to test.
	err = pSync.SetQueue(eps2)

	// Get the queue to check if the content has been set correctly.
	newQ := pSync.GetQueue()

	// Remove the IDs of the obtained queue to avoid a false positive.
	for i := range newQ.Content {
		newQ.Content[i].ID = 0
	}

	assert.NoError(err, "the queue should be overwritten without errors")
	assert.Equal(*eps2, newQ.Content, "the content of the queue should be set correctly")
	assert.NotEqual(*eps1, newQ.Content, "the previous content of the queue should be overwritten by the new one")
}

func (s *SynchronizerTestSuite) TestCleanQueue() {
	assert := assert2.New(s.T())

	_, err := s.insertRandomQueueEps(s.dbInstanceQ3)
	if err != nil {
		panic(err)
	}

	pSync, err := New(s.dbInstanceQ3)
	if err != nil {
		panic(err)
	}

	err = pSync.CleanQueue()

	assert.NoError(err, "the queue should be cleaned without errors")

	sqlDB := s.dbInstanceQ3.GetInstance()

	rows, err := sqlDB.Query("SELECT * FROM player_queue;")
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = rows.Close()
	}()

	assert.False(rows.Next(), "there should be no rows because the table is supposed to be empty")
}

func (s *SynchronizerTestSuite) TestAddToQueue() {
	assert := assert2.New(s.T())

	_, err := s.insertRandomQueueEps(s.dbInstanceQ4)
	if err != nil {
		panic(err)
	}

	pSync, err := New(s.dbInstanceQ4)
	if err != nil {
		panic(err)
	}

	anotherEp1 := QueueEpisode{
		ID:        0,
		PodcastID: int(lTesting.RandomInt63(100)),
		EpisodeID: lTesting.RandomString(5),
		// The position should be managed by the method, so this value should not have any effect.
		Position: 0,
	}

	anotherEp2 := QueueEpisode{
		ID:        0,
		PodcastID: int(lTesting.RandomInt63(100)),
		EpisodeID: lTesting.RandomString(5),
		// The position should be managed by the method, so this value should not have any effect.
		Position: 0,
	}

	id1, err := pSync.AddToQueue(anotherEp1, false)

	assert.NoError(err, "the episode should be added to the queue without errors")
	assert.True(id1 > 0, "the ID given by the database to the new episode should be returned")

	anotherEp1DB, err := s.getQueueEpByEpisodeID(s.dbInstanceQ4, anotherEp1.EpisodeID)
	if err != nil {
		panic(err)
	}

	assert.Equal(anotherEp1DB.ID, id1, "the ID returned by the method 'AddToQueue' should be correct")
	assert.Equal(11, anotherEp1DB.Position, "the position of the inserted episode should equal to the"+
		" length of the queue")
	assert.Equal(anotherEp1.PodcastID, anotherEp1DB.PodcastID, "the 'PodcastID' of the episode returned"+
		" should be the same as the inserted one")
	assert.Equal(anotherEp1.EpisodeID, anotherEp1DB.EpisodeID, "the 'EpisodeID' of the episode returned"+
		" should be the same as the inserted one")

	id2, err := pSync.AddToQueue(anotherEp2, true)

	assert.NoError(err, "the episode should be added to the queue without errors")
	assert.True(id2 > 0, "the ID given by the database to the new episode should be returned")

	anotherEp2DB, err := s.getQueueEpByEpisodeID(s.dbInstanceQ4, anotherEp2.EpisodeID)
	if err != nil {
		panic(err)
	}

	assert.Equal(anotherEp2DB.ID, id2, "the ID returned by the method 'AddToQueue' should be correct")
	// Positions start at 1
	assert.Equal(1, anotherEp2DB.Position, "the position of the inserted episode should be the"+
		" position of the last element in the queue + 1. Note: 'last element in the queue' refers to the episode that"+
		" is literally at the end of the queue, no the episode stored in the last row of the table 'player_queue'.")
	assert.Equal(anotherEp2.PodcastID, anotherEp2DB.PodcastID, "the 'PodcastID' of the episode returned"+
		" should be the same as the inserted one")
	assert.Equal(anotherEp2.EpisodeID, anotherEp2DB.EpisodeID, "the 'EpisodeID' of the episode returned"+
		" should be the same as the inserted one")
}

func (s *SynchronizerTestSuite) TestRemoveFromQueue() {
	assert := assert2.New(s.T())

	eps, err := s.insertRandomQueueEps(s.dbInstanceQ5)
	if err != nil {
		panic(err)
	}

	pSync, err := New(s.dbInstanceQ5)
	if err != nil {
		panic(err)
	}

	err = pSync.RemoveFromQueue((*eps)[1].ID)

	assert.NoError(err, "the episode should be removed without errors")

	_, err = s.getQueueEpByEpisodeID(s.dbInstanceQ5, (*eps)[1].EpisodeID)

	assert.True(errorx.IsOfType(err, errorx.IllegalArgument), "the row should be removed from the table"+
		" 'player_queue'")

	err = pSync.RemoveFromQueue(9999)

	if assert.Error(err, "if a non-existent ID is used, an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.IllegalArgument), "the returned error should be of type errorx.IllegalArgument")
	}
}

func (s *SynchronizerTestSuite) AfterTest(_, _ string) {}

func (s *SynchronizerTestSuite) TearDownTest() {
	_ = s.dbInstance1.Close()
	_ = s.dbInstance2.Close()
	_ = s.dbInstanceQ1.Close()
	_ = s.dbInstanceQ2.Close()
	_ = s.dbInstanceQ3.Close()
	_ = s.dbInstanceQ4.Close()
	_ = s.dbInstanceQ5.Close()

	err := os.RemoveAll(s.dbPath)
	if err != nil {
		panic(err)
	}
}

func TestSynchronizerTestSuite(t *testing.T) {
	suite.Run(t, new(SynchronizerTestSuite))
}

func (s *SynchronizerTestSuite) insertRandomQueueEps(db *database.Database) (*[]QueueEpisode, error) {
	eps := s.generateDummyQueueEp(10, true)

	sqlDB := db.GetInstance()

	for _, e := range *eps {
		r, err := sqlDB.Exec("INSERT INTO player_queue (podcast_id, episode_id, position) VALUES (?, ?, ?);",
			e.PodcastID, e.EpisodeID, e.Position)
		if err != nil {
			panic(err)
		}

		rowsAffected, err := r.RowsAffected()
		if err != nil {
			return nil, err
		}

		if rowsAffected != 1 {
			return nil, fmt.Errorf("unexpected number of rows affected: %d", rowsAffected)
		}
	}

	return eps, nil
}

func (s *SynchronizerTestSuite) generateDummyQueueEp(n int, withID bool) *[]QueueEpisode {
	var eps []QueueEpisode

	for i := 0; i < n; i++ {
		var ID int

		if withID {
			ID = i + 1
		}

		ep := QueueEpisode{
			ID:        ID,
			PodcastID: int(lTesting.RandomInt63(20)),
			EpisodeID: lTesting.RandomString(10),
			Position:  i + 1,
		}

		eps = append(eps, ep)
	}

	return &eps
}

func (s *SynchronizerTestSuite) getQueueEpByEpisodeID(db *database.Database, id string) (QueueEpisode, error) {
	sqlDB := db.GetInstance()

	rows, err := sqlDB.Query("SELECT * FROM player_queue WHERE episode_id = ?;", id)
	if err != nil {
		return QueueEpisode{}, err
	}

	defer func() {
		_ = rows.Close()
	}()

	if !rows.Next() {
		return QueueEpisode{}, errorx.IllegalArgument.New("the episode ID '%s' does not exist in the database",
			id)
	}

	var ep QueueEpisode
	err = rows.Scan(&ep.ID, &ep.PodcastID, &ep.EpisodeID, &ep.Position)
	if err != nil {
		return QueueEpisode{}, err
	}

	return ep, nil
}

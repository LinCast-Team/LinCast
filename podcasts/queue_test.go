package podcasts

import (
	"os"
	"runtime"
	"testing"

	"github.com/joomcode/errorx"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type QueueTestSuite struct {
	length     int
	dbPath     string
	dbFilename string

	suite.Suite
}

func (s *QueueTestSuite) SetupTest() {
	s.length = runtime.NumCPU()
	s.dbPath = "./test_queue"
	s.dbFilename = "queue_test.sqlite"

	err := os.Mkdir(s.dbPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (s *QueueTestSuite) BeforeTest(_, _ string) {}

func (s *QueueTestSuite) TestNewUpdateQueue() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "new_update_queue_"+s.dbFilename)
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

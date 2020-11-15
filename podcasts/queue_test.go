package podcasts

import (
	"runtime"
	"testing"

	"github.com/joomcode/errorx"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type QueueTestSuite struct {
	length int

	suite.Suite
}

func (s *QueueTestSuite) SetupTest() {
	s.length = runtime.NumCPU()
}

func (s *QueueTestSuite) BeforeTest(_, _ string) {}

func (s *QueueTestSuite) TestNewQueue() {
	assert := assert2.New(s.T())

	q, err := NewQueue(s.length)

	assert.NoError(err, "the queue should be created correctly")
	assert.NotNil(q, "the returned Queue instance shouldn't be nil")

	q, err = NewQueue(-1)

	if assert.Error(err, "if the argument that corresponds to the length of the queue is a negative"+
		" number or 0, an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.IllegalArgument), "an error of type "+
			"errorx.IllegalArgument should be returned")
	}
	assert.Nil(q, "if the argument is incorrect, the instance of Queue returned should be nil")
}

func (s *QueueTestSuite) AfterTest(_, _ string) {}

func (s *QueueTestSuite) TearDownTest() {}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}

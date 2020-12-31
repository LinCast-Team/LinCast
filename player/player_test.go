package player

import (
	"github.com/stretchr/testify/suite"
	"lincast/database"
	"testing"
)

type SynchronizerTestSuite struct {
	dbPath     string
	dbFilename string
	dbInstance *database.Database

	suite.Suite
}

func (s *SynchronizerTestSuite) SetupTest() {
	s.dbPath = "./test_synchronizer"
}

func (s *SynchronizerTestSuite) BeforeTest(_, _ string) {

}

func (s *SynchronizerTestSuite) TestNew() {

}

func (s *SynchronizerTestSuite) TestUpdateProgress() {

}

func (s *SynchronizerTestSuite) TestGetProgress() {

}

func (s *SynchronizerTestSuite) AfterTest(_, _ string) {

}

func (s *SynchronizerTestSuite) TearDownTest() {

}

func TestSynchronizerTestSuite(t *testing.T) {
	suite.Run(t, new(SynchronizerTestSuite))
}

	package psync

import (
	"os"
	"testing"
	"time"

	"lincast/database"

	"github.com/joomcode/errorx"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SynchronizerTestSuite struct {
	dbPath      string
	dbFilename  string
	dbInstance1 *database.Database
	dbInstance2 *database.Database

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

	db2, err := database.New(s.dbPath, "1"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	s.dbInstance2 = db2
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

func (s *SynchronizerTestSuite) AfterTest(_, _ string) {}

func (s *SynchronizerTestSuite) TearDownTest() {
	_ = s.dbInstance1.Close()
	_ = s.dbInstance2.Close()

	err := os.RemoveAll(s.dbPath)
	if err != nil {
		panic(err)
	}
}

func TestSynchronizerTestSuite(t *testing.T) {
	suite.Run(t, new(SynchronizerTestSuite))
}

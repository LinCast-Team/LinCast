package podcasts

import (
	"github.com/joomcode/errorx"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PodcastsTestSuite struct {
	feedURL string

	suite.Suite
}

func (s *PodcastsTestSuite) SetupTest() {
	s.feedURL = "https://feeds.emilcar.fm/daily"
}

func (s *PodcastsTestSuite) BeforeTest(_, _ string) {}

func (s *PodcastsTestSuite) TestQueue() {

}

func (s *PodcastsTestSuite) TestNewPodcast() {
	assert := assert2.New(s.T())

	p, err := GetPodcast(s.feedURL)

	assert.NoError(err, "the podcast should be created without errors")
	assert.NotNil(p, "the struct returned should contain the info of the podcast")

	wrongURL := "something-wrong.com"
	p, err = GetPodcast(wrongURL)

	if assert.Error(err, "if the passed url is incorrect an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.IllegalFormat), "the error should be of type IllegalFormat")
	}
	assert.Nil(p, "the returned struct should be nil")

	wrongURL = "http://localhost:8080"
	p, err = GetPodcast(wrongURL)

	if assert.Error(err, "if there is a problem with the request an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.DataUnavailable), "the returned error should be of type "+
			"DataUnavailable")
	}
	assert.Nil(p, "the returned struct should be nil")
}

func (s *PodcastsTestSuite) AfterTest(_, _ string) {}

func (s *PodcastsTestSuite) TearDownTest() {}

func TestPodcastsTestSuite(t *testing.T) {
	suite.Run(t, new(PodcastsTestSuite))
}

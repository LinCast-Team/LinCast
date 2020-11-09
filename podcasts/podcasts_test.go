package podcasts

import (
	"github.com/joomcode/errorx"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"reflect"
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

func (s *PodcastsTestSuite) TestGetPodcast() {
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
		assert.True(errorx.IsOfType(err, errorx.ExternalError), "the returned error should be of type "+
			"ExternalError")
	}
	assert.Nil(p, "the returned struct should be nil")
}

func (s *PodcastsTestSuite) TestGetEpisodes() {
	assert := assert2.New(s.T())

	p, err := GetPodcast(s.feedURL)
	if err != nil {
		assert.FailNow(err.Error(), "a podcast to test method `GetEpisodes` is needed")
	}
	if reflect.ValueOf(p).IsNil() {
		assert.FailNow("a pointer to a Podcast instance is needed")
	}

	eps, err := p.GetEpisodes()

	assert.NoError(err, "episodes should be obtained without errors")
	assert.True(len(*eps) != 0, "a slice with episodes should be returned")
}

func (s *PodcastsTestSuite) AfterTest(_, _ string) {}

func (s *PodcastsTestSuite) TearDownTest() {}

func TestPodcastsTestSuite(t *testing.T) {
	suite.Run(t, new(PodcastsTestSuite))
}

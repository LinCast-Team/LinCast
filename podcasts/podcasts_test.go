package podcasts

import (
	"github.com/joomcode/errorx"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type PodcastsTestSuite struct {
	sampleFeeds []string

	suite.Suite
}

func (s *PodcastsTestSuite) SetupTest() {
	s.sampleFeeds = []string{
		"https://changelog.com/gotime/feed",
		"https://feeds.emilcar.fm/daily",
		"https://www.ivoox.com/podcast-despeja-x-by-xataka_fg_f1579492_filtro_1.xml",
		"https://www.ivoox.com/podcast-tortulia-podcast-episodios_fg_f1157653_filtro_1.xml",
	}
}

func (s *PodcastsTestSuite) BeforeTest(_, _ string) {}

func (s *PodcastsTestSuite) TestGetPodcast() {
	assert := assert2.New(s.T())

	for _, feed := range s.sampleFeeds {
		p, err := GetPodcast(feed)

		assert.NoErrorf(err, "the podcast should be created without errors (feed %s)", feed)
		assert.NotNil(p, "the struct returned should contain the info of the podcast and"+
			" not be nil (feed %s)", feed)
	}

	wrongURL := "something-wrong.com"
	p, err := GetPodcast(wrongURL)

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

	for _, feed := range s.sampleFeeds {
		p, err := GetPodcast(feed)
		if err != nil {
			panic(errorx.Decorate(err, "podcast can't be obtained (feed %s)", feed))
		}
		if reflect.ValueOf(p).IsNil() {
			panic(errorx.Decorate(err, "the returned Podcast struct is nil (feed %s)", feed))
		}

		eps, err := p.GetEpisodes()

		assert.NoError(err, "episodes should be obtained without errors")
		assert.True(len(*eps) != 0, "a slice with episodes should be returned")
	}
}

func (s *PodcastsTestSuite) AfterTest(_, _ string) {}

func (s *PodcastsTestSuite) TearDownTest() {}

func TestPodcastsTestSuite(t *testing.T) {
	suite.Run(t, new(PodcastsTestSuite))
}

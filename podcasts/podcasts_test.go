package podcasts

import (
	"reflect"
	"testing"

	"github.com/joomcode/errorx"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
		p, originalFeed, err := GetPodcastData(feed)

		assert.NoErrorf(err, "the podcast should be created without errors (feed %s)", feed)
		assert.NotNil(p, "the struct returned should contain the info of the podcast and"+
			" not be nil (feed %s)", feed)

		assert.NotNil(originalFeed, "the original feed (parsed on a structure gofeed.Feed) should be returned")
	}

	wrongURL := "something-wrong.com"
	p, _, err := GetPodcastData(wrongURL)

	if assert.Error(err, "if the passed url is incorrect an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.ExternalError), "the error should be of type "+
			"errorx.ExternalError")
	}
	assert.Nil(p, "the returned struct should be nil")

	wrongURL = "http://localhost:8080"
	p, _, err = GetPodcastData(wrongURL)

	if assert.Error(err, "if there is a problem with the request an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.ExternalError), "the returned error should be of type "+
			"ExternalError")
	}
	assert.Nil(p, "the returned struct should be nil")
}

func (s *PodcastsTestSuite) TestGetEpisodes() {
	assert := assert2.New(s.T())

	for _, feed := range s.sampleFeeds {
		_, originalFeed, err := GetPodcastData(feed)
		if err != nil {
			panic(errorx.Decorate(err, "the feed '%s' can't be obtained", feed))
		}

		if reflect.ValueOf(originalFeed).IsNil() {
			panic(errorx.Decorate(err, "the returned gofeed.Feed struct is nil (feed %s), which is not expected", feed))
		}

		eps, err := GetEpisodes(originalFeed)

		assert.NoError(err, "episodes should be obtained without errors")
		assert.True(len(*eps) != 0, "a slice with episodes should be returned")
	}
}

func (s *PodcastsTestSuite) AfterTest(_, _ string) {}

func (s *PodcastsTestSuite) TearDownTest() {}

func TestPodcastsTestSuite(t *testing.T) {
	suite.Run(t, new(PodcastsTestSuite))
}

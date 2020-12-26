package backend

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"lincast/podcasts"

	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HandlersTestSuite struct {
	podcastsDBPath     string
	podcastsDBFilename string
	sampleFeeds        []string

	suite.Suite
}

func (s *HandlersTestSuite) SetupTest() {
	s.podcastsDBPath = "./backend_test"
	s.podcastsDBFilename = "handlers_podcastsDB_test.sqlite"

	err := os.Mkdir(s.podcastsDBPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	db, err := podcasts.NewDB(s.podcastsDBPath, s.podcastsDBFilename)
	if err != nil {
		panic(err)
	}

	_podcastsDB = db

	// Prepare the database for tests.
	s.sampleFeeds = []string{
		"https://changelog.com/gotime/feed",
		"https://feeds.emilcar.fm/daily",
		"https://www.ivoox.com/podcast-despeja-x-by-xataka_fg_f1579492_filtro_1.xml",
	}

	for _, feed := range s.sampleFeeds {
		p, err := podcasts.GetPodcast(feed)
		if err != nil {
			panic(err)
		}

		err = _podcastsDB.InsertPodcast(p)
		if err != nil {
			panic(err)
		}
	}
}

func (s *HandlersTestSuite) BeforeTest(_, _ string) {}

func (s *HandlersTestSuite) TestSubscribeToPodcastHandler() {
	type subscription struct {
		URL string `json:"url"`
	}

	assert := assert2.New(s.T())
	res1 := httptest.NewRecorder()
	res2 := httptest.NewRecorder()

	sub := subscription{
		URL: "https://www.ivoox.com/podcast-tortulia-podcast-episodios_fg_f1157653_filtro_1.xml",
	}

	c, err := json.Marshal(sub)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/v0/feeds/subscribe", bytes.NewReader(c))
	newRouter(false, false).ServeHTTP(res1, req)

	assert.Equal(http.StatusOK, res1.Code, "if the body of the request has no issues, it should be"+
		" responded with the HTTP status code 200 (OK)")

	sub.URL = "ivoox.asdfsadf/podcast-tortulia-podcast-episodios_fg_f1157653_filtro_1.xml"

	c, err = json.Marshal(sub)
	if err != nil {
		panic(err)
	}

	req = httptest.NewRequest("POST", "/api/v0/feeds/subscribe", bytes.NewReader(c))
	newRouter(false, false).ServeHTTP(res2, req)

	assert.Equal(http.StatusBadRequest, res2.Code, "if the body of the request contains a not valid URL,"+
		" it should be responded with the HTTP status code 400 (Bad Request)")
}

func (s *HandlersTestSuite) TestUnSubscribeToPodcastHandler() {

}

func (s *HandlersTestSuite) TestGetPodcastHandler() {

}

func (s *HandlersTestSuite) TestGetUserPodcastsHandler() {

}

func (s *HandlersTestSuite) AfterTest(_, _ string) {}

func (s *HandlersTestSuite) TearDownTest() {
	err := os.RemoveAll(s.podcastsDBPath)
	if err != nil {
		panic(err)
	}
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

package backend

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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
		"https://changelog.com/gotime/feed", // Will be removed by TestUnSubscribeToPodcastHandler
		"https://feeds.emilcar.fm/daily",
		"https://www.ivoox.com/podcast-despeja-x-by-xataka_fg_f1579492_filtro_1.xml",
		"https://anchor.fm/s/14d75a0/podcast/rss",
		"http://feeds.feedburner.com/LaVinetaEnDiscoInferno",
		"https://anchor.fm/s/564ad20/podcast/rss/",
		"http://feeds.feedburner.com/StacktraceBy9to5mac",
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
	type reqBody struct {
		URL string `json:"url"`
	}

	assert := assert2.New(s.T())
	res := httptest.NewRecorder()

	r := reqBody{
		URL: "https://www.ivoox.com/podcast-tortulia-podcast-episodios_fg_f1157653_filtro_1.xml",
	}

	c, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/v0/feeds/subscribe", bytes.NewReader(c))
	newRouter(false, false).ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code, "if the body of the request has no issues, it should be"+
		" responded with the HTTP status code 200 (OK)")

	r.URL = "ivoox.asdfsadf/podcast-tortulia-podcast-episodios_fg_f1157653_filtro_1.xml"

	c, err = json.Marshal(r)
	if err != nil {
		panic(err)
	}

	req = httptest.NewRequest("POST", "/api/v0/feeds/subscribe", bytes.NewReader(c))
	res = httptest.NewRecorder()
	newRouter(false, false).ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code, "if the body of the request contains a not valid URL,"+
		" it should be responded with the HTTP status code 400 (Bad Request)")
}

func (s *HandlersTestSuite) TestUnSubscribeToPodcastHandler() {
	assert := assert2.New(s.T())
	res := httptest.NewRecorder()
	id := 1

	req := httptest.NewRequest("PUT", "/api/v0/feeds/unsubscribe?id="+strconv.Itoa(id), nil)
	newRouter(false, false).ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code, "the request should be processed correctly, returning"+
		" a 200 HTTP status code")

	res = httptest.NewRecorder()
	req = httptest.NewRequest("PUT", "/api/v0/feeds/unsubscribe?id="+strconv.Itoa(100), nil)
	newRouter(false, false).ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code, "the usage of a non existent ID should return"+
		" a 400 HTTP status code")

	res = httptest.NewRecorder()
	req = httptest.NewRequest("PUT", "/api/v0/feeds/unsubscribe", nil)
	newRouter(false, false).ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code, "if the request does not contain an ID param,"+
		" it should return a 400 HTTP status code (Bad Request)")
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

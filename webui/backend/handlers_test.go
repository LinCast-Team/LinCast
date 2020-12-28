package backend

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
		"https://changelog.com/gotime/feed", // Will be unsubscribed by TestUnSubscribeToPodcastHandler
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

	req := httptest.NewRequest("GET", "/api/v0/podcasts/subscribe", bytes.NewReader(c))
	newRouter(false, false).ServeHTTP(res, req)

	assert.Equal(http.StatusNotFound, res.Code, "the usage of a incorrect method should return"+
		" a 404 HTTP status code")

	req = httptest.NewRequest("POST", "/api/v0/podcasts/subscribe", bytes.NewReader(c))
	res = httptest.NewRecorder()
	newRouter(false, false).ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code, "if the body of the request has no issues, it should be"+
		" responded with the HTTP status code 200 (OK)")

	r.URL = "ivoox.asdfsadf/podcast-tortulia-podcast-episodios_fg_f1157653_filtro_1.xml"

	c, err = json.Marshal(r)
	if err != nil {
		panic(err)
	}

	req = httptest.NewRequest("POST", "/api/v0/podcasts/subscribe", bytes.NewReader(c))
	res = httptest.NewRecorder()
	newRouter(false, false).ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code, "if the body of the request contains a not valid URL,"+
		" it should be responded with the HTTP status code 400 (Bad Request)")
}

func (s *HandlersTestSuite) TestUnSubscribeToPodcastHandler() {
	assert := assert2.New(s.T())
	res := httptest.NewRecorder()
	id := 1

	req := httptest.NewRequest("GET", "/api/v0/podcasts/unsubscribe?id="+strconv.Itoa(id), nil)
	newRouter(false, false).ServeHTTP(res, req)

	assert.Equal(http.StatusNotFound, res.Code, "the usage of a incorrect method should return"+
		" a 404 HTTP status code")

	res = httptest.NewRecorder()
	req = httptest.NewRequest("PUT", "/api/v0/podcasts/unsubscribe?id="+strconv.Itoa(id), nil)
	newRouter(false, false).ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code, "the request should be processed correctly, returning"+
		" a 200 HTTP status code")

	res = httptest.NewRecorder()
	req = httptest.NewRequest("PUT", "/api/v0/podcasts/unsubscribe?id="+strconv.Itoa(100), nil)
	newRouter(false, false).ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code, "the usage of a non existent ID should return"+
		" a 400 HTTP status code")

	res = httptest.NewRecorder()
	req = httptest.NewRequest("PUT", "/api/v0/podcasts/unsubscribe", nil)
	newRouter(false, false).ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code, "if the request does not contain an ID param,"+
		" it should return a 400 HTTP status code (Bad Request)")
}

func (s *HandlersTestSuite) TestGetPodcastHandler() {

}

func (s *HandlersTestSuite) TestGetUserPodcastsHandler() {
	assert := assert2.New(s.T())
	res := httptest.NewRecorder()

	req := httptest.NewRequest("POST", "/api/v0/podcasts/user?subscribed=true&unsubscribed=true", nil)
	newRouter(false, false).ServeHTTP(res, req)

	assert.Equal(http.StatusNotFound, res.Code, "the usage of a incorrect method should return"+
		" a 404 HTTP status code")

	req = httptest.NewRequest("GET", "/api/v0/podcasts/user?subscribed=true&unsubscribed=true", nil)
	res = httptest.NewRecorder()
	newRouter(false, false).ServeHTTP(res, req)

	if assert.Equal(http.StatusOK, res.Code, "the request should be processed correctly, returning"+
		" a 200 HTTP status code") {
		assert.Equal("application/json", res.Header().Get("Content-Type"), "the response"+
			" should contain the appropriate 'Content-Type' headers'")

		p := map[string][]podcasts.Podcast{
			"subscribed":   {},
			"unsubscribed": {},
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(body, &p)
		if err != nil {
			panic(err)
		}

		// The length can variate by one if the test of subscription is executed before this (remember that they're run
		// in parallel).
		if !assert.True(len(p["subscribed"])+len(p["unsubscribed"]) == len(s.sampleFeeds) ||
			len(p["subscribed"])+len(p["unsubscribed"]) == len(s.sampleFeeds)+1,
			"a slice of podcasts should be returned as the response's body") {
			s.T().Logf("Len of received subscribed podcasts: %d, Len of received unsubscribed podcasts: %d,"+
				" Expacted total len: %d (+/- 1)\n", len(p["subscribed"]), len(p["unsubscribed"]), len(s.sampleFeeds))
		}
	}

	req = httptest.NewRequest("GET", "/api/v0/podcasts/user", nil)
	res = httptest.NewRecorder()
	newRouter(false, false).ServeHTTP(res, req)

	if assert.Equal(http.StatusOK, res.Code, "if the request does not contain parameters it should be"+
		" responded with all the podcasts and a HTTP status code 200 (OK)") {
		assert.Equal("application/json", res.Header().Get("Content-Type"), "the response"+
			" should contain the appropriate 'Content-Type' headers'")

		p := map[string][]podcasts.Podcast{
			"subscribed":   {},
			"unsubscribed": {},
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(body, &p)
		if err != nil {
			panic(err)
		}

		// The length can variate by one if the test of subscription is executed before this (remember that they're run
		// in parallel).
		if !assert.True(len(p["subscribed"])+len(p["unsubscribed"]) == len(s.sampleFeeds) ||
			len(p["subscribed"])+len(p["unsubscribed"]) == len(s.sampleFeeds)+1,
			"a slice of podcasts should be returned as the response's body") {
			s.T().Logf("Len of received subscribed podcasts: %d, Len of received unsubscribed podcasts: %d,"+
				" Expacted total len: %d (+/- 1)\n", len(p["subscribed"]), len(p["unsubscribed"]), len(s.sampleFeeds))
		}
	}

	req = httptest.NewRequest("GET", "/api/v0/podcasts/user?subscribed=true&unsubscribed=hello", nil)
	res = httptest.NewRecorder()
	newRouter(false, false).ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code, "if the content of a parameter is incorrect, the"+
		" response's HTTP code should be 400 (Bad Request)")
	assert.Equal("", res.Header().Get("Content-Type"), "the response should"+
		" not contain the 'Content-Type' headers'")
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

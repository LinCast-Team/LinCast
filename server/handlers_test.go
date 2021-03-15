package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"lincast/database"
	"lincast/podcasts"
	"lincast/psync"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HandlersTestSuite struct {
	podcastsDBPath     string
	podcastsDBFilename string
	sampleFeeds        []string

	podcastsMMutex sync.Mutex // Use this only for tests related with podcasts management.
	queueMutex     sync.Mutex // Use this only for tests related with the queue of the player.
	suite.Suite
}

func (s *HandlersTestSuite) SetupTest() {
	log.SetOutput(io.Discard)

	s.podcastsDBPath = "./backend_test"
	s.podcastsDBFilename = "handlers_podcastsDB_test.sqlite"

	err := os.Mkdir(s.podcastsDBPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	db, err := database.New(s.podcastsDBPath, s.podcastsDBFilename)
	if err != nil {
		panic(err)
	}

	_podcastsDB = db

	// Prepare the database for tests.
	s.sampleFeeds = []string{
		"https://changelog.com/gotime/feed", // Will be unsubscribed by TestUnSubscribeToPodcastHandler
		"https://feeds.emilcar.fm/daily",    // Will be used by TestGetPodcastHandler & TestGetEpisodeHandler
		"https://www.ivoox.com/podcast-despeja-x-by-xataka_fg_f1579492_filtro_1.xml",
		"https://anchor.fm/s/14d75a0/podcast/rss",
		"http://feeds.feedburner.com/LaVinetaEnDiscoInferno",
		"https://anchor.fm/s/564ad20/podcast/rss/",
		"http://feeds.feedburner.com/StacktraceBy9to5mac",
	}

	for i, feed := range s.sampleFeeds {
		p, err := podcasts.GetPodcast(feed)
		if err != nil {
			panic(err)
		}

		err = _podcastsDB.InsertPodcast(p)
		if err != nil {
			panic(err)
		}

		if i != 1 {
			continue
		}

		p.ID = i

		eps, err := p.GetEpisodes()
		if err != nil {
			panic(err)
		}

		for _, ep := range *eps {
			err = _podcastsDB.InsertEpisode(&ep)
			if err != nil {
				panic(err)
			}
		}
	}

	ps, err := psync.New(_podcastsDB)
	if err != nil {
		assert2.FailNowf(s.T(), "the player synchronizer can't be instantiated", "error when trying to"+
			" instantiate the player synchronizer: %s", errorx.EnsureStackTrace(err))
	}

	_playerSync = ps
}

func (s *HandlersTestSuite) BeforeTest(_, _ string) {}

func (s *HandlersTestSuite) TestSubscribeToPodcastHandler() {
	type reqBody struct {
		URL string `json:"url"`
	}

	assert := assert2.New(s.T())

	r := reqBody{
		URL: "https://www.ivoox.com/podcast-tortulia-podcast-episodios_fg_f1157653_filtro_1.xml",
	}

	c, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	res := s.newRequest(http.MethodGet, "/api/v0/podcasts/subscribe", bytes.NewReader(c))

	assert.Equal(http.StatusNotFound, res.Code, "the usage of an incorrect method should return"+
		" a 404 HTTP status code")

	s.podcastsMMutex.Lock()

	res = s.newRequest(http.MethodPost, "/api/v0/podcasts/subscribe", bytes.NewReader(c))

	s.podcastsMMutex.Unlock()

	assert.Equal(http.StatusCreated, res.Code, "if the body of the request has no issues, it should be"+
		" responded with the HTTP status code 201 (Created)")

	res = s.newRequest(http.MethodPost, "/api/v0/podcasts/subscribe", bytes.NewReader(c))

	assert.Equal(http.StatusConflict, res.Code, "if the submitted feed already exists on the database"+
		" the request should be responded with the HTTP status code 409 (Conflict)")

	r.URL = "ivoox.something-wrong/podcast-tortulia-podcast-episodios_fg_f1157653_filtro_1.xml"

	c, err = json.Marshal(r)
	if err != nil {
		panic(err)
	}

	res = s.newRequest(http.MethodPost, "/api/v0/podcasts/subscribe", bytes.NewReader(c))

	assert.Equal(http.StatusBadRequest, res.Code, "if the body of the request contains a not valid URL,"+
		" it should be responded with the HTTP status code 400 (Bad Request)")
}

func (s *HandlersTestSuite) TestUnSubscribeToPodcastHandler() {
	assert := assert2.New(s.T())
	id := 1

	res := s.newRequest(http.MethodGet, "/api/v0/podcasts/unsubscribe?id="+strconv.Itoa(id), nil)

	assert.Equal(http.StatusNotFound, res.Code, "the usage of an incorrect method should return"+
		" a 404 HTTP status code")

	s.podcastsMMutex.Lock()

	res = s.newRequest(http.MethodPut, "/api/v0/podcasts/unsubscribe?id="+strconv.Itoa(id), nil)

	s.podcastsMMutex.Unlock()

	assert.Equal(http.StatusNoContent, res.Code, "the request should be processed correctly, returning"+
		" a 204 HTTP status code")

	res = s.newRequest(http.MethodPut, "/api/v0/podcasts/unsubscribe?id="+strconv.Itoa(100), nil)

	assert.Equal(http.StatusBadRequest, res.Code, "the usage of a non existent ID should return"+
		" a 400 HTTP status code")

	res = s.newRequest(http.MethodPut, "/api/v0/podcasts/unsubscribe", nil)

	assert.Equal(http.StatusBadRequest, res.Code, "if the request does not contain an ID param,"+
		" it should return a 400 HTTP status code (Bad Request)")
}

func (s *HandlersTestSuite) TestGetPodcastHandler() {
	assert := assert2.New(s.T())
	id := 2

	res := s.newRequest(http.MethodPost, "/api/v0/podcasts/"+strconv.Itoa(id)+"/details", nil)

	assert.Equal(http.StatusNotFound, res.Code, "the usage of an incorrect method should return"+
		" a 404 HTTP status code")
	assert.Equal("", res.Header().Get("Content-Type"), "the response should not contain"+
		" the 'Content-Type' headers'")

	res = s.newRequest(http.MethodGet, "/api/v0/podcasts/"+strconv.Itoa(id)+"/details", nil)

	assert.Equal(http.StatusOK, res.Code, "the request should be processed correctly, returning"+
		" a 200 HTTP status code")
	assert.Equal("application/json", res.Header().Get("Content-Type"), "the response"+
		" should contain the appropriate 'Content-Type' headers'")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var receivedPodcast podcasts.Podcast
	err = json.Unmarshal(body, &receivedPodcast)
	if err != nil {
		panic(err)
	}

	p, err := _podcastsDB.GetPodcastByID(id)
	if err != nil {
		panic(err)
	}

	receivedPodcast.Updated = time.Time{}
	receivedPodcast.LastCheck = time.Time{}
	receivedPodcast.Added = time.Time{}
	p.Updated = time.Time{}
	p.LastCheck = time.Time{}
	p.Added = time.Time{}

	assert.Equal(*p, receivedPodcast, "the response should contain the same podcast as the one stored"+
		" in the database")

	res = s.newRequest(http.MethodGet, "/api/v0/podcasts/"+strconv.Itoa(100)+"/details", nil)

	assert.Equal(http.StatusNotFound, res.Code, "if the used ID does not exist, the response should be"+
		" with the HTTP status code 404 (Not Found)")
	assert.Equal("text/plain; charset=utf-8", res.Header().Get("Content-Type"), "the response should not contain"+
		" the 'Content-Type' headers'")
}

func (s *HandlersTestSuite) TestGetUserPodcastsHandler() {
	assert := assert2.New(s.T())

	res := s.newRequest(http.MethodPost, "/api/v0/podcasts/user?subscribed=true&unsubscribed=true", nil)

	assert.Equal(http.StatusNotFound, res.Code, "the usage of a incorrect method should return"+
		" a 404 HTTP status code")
	assert.Equal("", res.Header().Get("Content-Type"), "the response should not contain"+
		" the 'Content-Type' headers'")

	res = s.newRequest(http.MethodGet, "/api/v0/podcasts/user?subscribed=true&unsubscribed=true", nil)

	if assert.Equal(http.StatusOK, res.Code, "the request should be processed correctly, returning"+
		" a 200 HTTP status code") {
		assert.Equal("application/json", res.Header().Get("Content-Type"), "the response"+
			" should contain the appropriate 'Content-Type' headers'")

		var p map[string][]podcasts.Podcast

		body, err := io.ReadAll(res.Body)
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

	res = s.newRequest(http.MethodGet, "/api/v0/podcasts/user", nil)

	if assert.Equal(http.StatusOK, res.Code, "if the request does not contain parameters it should be"+
		" responded with all the podcasts and a HTTP status code 200 (OK)") {
		assert.Equal("application/json", res.Header().Get("Content-Type"), "the response"+
			" should contain the appropriate 'Content-Type' headers'")

		var p map[string][]podcasts.Podcast

		body, err := io.ReadAll(res.Body)
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

	res = s.newRequest(http.MethodGet, "/api/v0/podcasts/user?subscribed=true&unsubscribed=hello", nil)

	assert.Equal(http.StatusBadRequest, res.Code, "if the content of a parameter is incorrect, the"+
		" response's HTTP code should be 400 (Bad Request)")
	assert.Equal("text/plain; charset=utf-8", res.Header().Get("Content-Type"), "the response should"+
		" not contain the 'Content-Type' headers'")

	s.podcastsMMutex.Lock()

	res = s.newRequest(http.MethodGet, "/api/v0/podcasts/user?subscribed=false&unsubscribed=true", nil)

	if assert.Equal(http.StatusOK, res.Code, "the request should be responded with a 200 HTTP status"+
		" code (OK)") {
		assert.Equal("application/json", res.Header().Get("Content-Type"), "the response"+
			" should contain the appropriate 'Content-Type' headers'")

		dbUnsubscribedPodcasts, err := _podcastsDB.GetPodcastsBySubscribedStatus(false)
		if err != nil {
			panic(errorx.EnsureStackTrace(err))
		}
		s.podcastsMMutex.Unlock()

		var p map[string][]podcasts.Podcast

		body, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(body, &p)
		if err != nil {
			panic(err)
		}

		if assert.Len(p["unsubscribed"], len(*dbUnsubscribedPodcasts), "the length of returned"+
			" podcasts should be the same as the length of podcasts stored in the database (with the required"+
			" subscription status)") {
			// The absence of metadata on fields of type time.Time will cause a false positive on the test, so we
			// should overwrite them. Here, we're assuming that the time stored on those fields is correct
			// something that can be incorrect. However, the correct storage and return of data must be checked
			// on the package that manages it (lincast/podcasts).
			for i := range *dbUnsubscribedPodcasts {
				(*dbUnsubscribedPodcasts)[i].Added = time.Time{}
				(*dbUnsubscribedPodcasts)[i].LastCheck = time.Time{}
				(*dbUnsubscribedPodcasts)[i].Updated = time.Time{}

				p["unsubscribed"][i].Added = time.Time{}
				p["unsubscribed"][i].LastCheck = time.Time{}
				p["unsubscribed"][i].Updated = time.Time{}
			}

			assert.Equal(*dbUnsubscribedPodcasts, p["unsubscribed"], "the returned unsubscribed podcasts"+
				" should be the same as the ones stored in the database")
		}
	} else {
		s.podcastsMMutex.Unlock()
	}
}

func (s *HandlersTestSuite) TestGetEpisodesHandler() {
	assert := assert2.New(s.T())
	id := 2

	res := s.newRequest(http.MethodPost, "/api/v0/podcasts/"+strconv.Itoa(id)+"/episodes", nil)

	assert.Equal(http.StatusNotFound, res.Code, "the usage of an incorrect method should return"+
		" a 404 HTTP status code")
	assert.Equal("", res.Header().Get("Content-Type"), "the response should not contain"+
		" the 'Content-Type' headers'")

	res = s.newRequest(http.MethodGet, "/api/v0/podcasts/"+strconv.Itoa(id)+"/episodes", nil)

	assert.Equal(http.StatusOK, res.Code, "the request should be processed correctly, returning"+
		" a 200 HTTP status code")
	assert.Equal("application/json", res.Header().Get("Content-Type"), "the response"+
		" should contain the appropriate 'Content-Type' headers'")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var receivedEps podcasts.Episodes
	err = json.Unmarshal(body, &receivedEps)
	if err != nil {
		panic(err)
	}

	for i := range receivedEps {
		receivedEps[i].Updated = time.Time{}
		receivedEps[i].Published = time.Time{}
	}

	eps, err := _podcastsDB.GetEpisodesByPodcast(id)
	if err != nil {
		panic(err)
	}

	for i := range *eps {
		(*eps)[i].Published = time.Time{}
		(*eps)[i].Updated = time.Time{}
	}

	assert.Equal(*eps, receivedEps, "the response should contain the same episodes as the ones that"+
		" are stored on the database")

	res = s.newRequest(http.MethodGet, "/api/v0/podcasts/"+strconv.Itoa(100)+"/episodes", nil)

	assert.Equal(http.StatusNotFound, res.Code, "if the used ID does not exist, the response should be"+
		" with the HTTP status code 404 (Not Found)")
	assert.Equal("text/plain; charset=utf-8", res.Header().Get("Content-Type"), "the response should not contain"+
		" the 'Content-Type' headers'")
}

func (s *HandlersTestSuite) TestPlayerProgressHandler() {
	assert := assert2.New(s.T())

	res := s.newRequest(http.MethodDelete, "/api/v0/player/progress", nil)

	assert.Equal(http.StatusNotFound, res.Code, "the usage of an incorrect method should return"+
		" a 404 HTTP status code")
	assert.Equal("", res.Header().Get("Content-Type"), "the response should not contain"+
		" the 'Content-Type' headers'")

	res = s.newRequest(http.MethodGet, "/api/v0/player/progress", nil)

	assert.Equal(http.StatusOK, res.Code, "the request should be processed correctly, returning"+
		" a 200 HTTP status code")
	assert.Equal("application/json", res.Header().Get("Content-Type"), "the response"+
		" should contain the appropriate 'Content-Type' headers'")

	var progress psync.CurrentProgress

	assert.NoError(s.parseProgressReq(res.Body, &progress), "the body of the request should have the"+
		" required format")

	p := psync.CurrentProgress{
		Progress:    time.Duration(16898),
		EpisodeGUID: "1234",
		PodcastID:   2,
	}

	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	res = s.newRequest(http.MethodPut, "/api/v0/player/progress", bytes.NewReader(b))

	assert.Equal(http.StatusCreated, res.Code, "the request should be processed correctly, returning"+
		" a 201 HTTP status code")
	assert.Equal(p, _playerSync.GetProgress(), "the progress should be correctly updated")

	res = s.newRequest(http.MethodGet, "/api/v0/player/progress", nil)

	assert.Equal(http.StatusOK, res.Code, "the request should be processed correctly, returning"+
		" a 200 HTTP status code")
	assert.Equal("application/json", res.Header().Get("Content-Type"), "the response"+
		" should contain the appropriate 'Content-Type' headers'")

	var p1 psync.CurrentProgress

	assert.NoError(s.parseProgressReq(res.Body, &p1), "the body of the request should have the"+
		" required format")

	assert.Equal(p, p1, "the returned progress should be the same as the one that we have")
}

func (s *HandlersTestSuite) TestQueueHandler() {
	assert := assert2.New(s.T())

	res := s.newRequest(http.MethodPost, "/api/v0/player/queue", nil)

	assert.Equal(http.StatusNotFound, res.Code, "the usage of an incorrect method should return"+
		" a 404 HTTP status code")
	assert.Equal("", res.Header().Get("Content-Type"), "the response should not contain"+
		" the 'Content-Type' headers'")

	/* ---------------------------- Test PUT requests --------------------------- */
	q := []psync.QueueEpisode{
		{
			ID:        0,
			PodcastID: 10,
			EpisodeID: "10d981dd-20dd-4bc8-9999-e12586a561d8",
			Position:  0,
		},
		{
			ID:        0,
			PodcastID: 8,
			EpisodeID: "5f069078-eb21-480d-9e8e-be92742d90f5",
			Position:  2, // the position is repeated on purpose
		},
		{
			ID:        0,
			PodcastID: 1,
			EpisodeID: "1c547723-43a5-46fe-91ae-4d410b0ebc79",
			Position:  2,
		},
		{
			ID:        0,
			PodcastID: 24,
			EpisodeID: "e0bab9d0-bad6-4b4f-afa4-4af40939b5c5",
			Position:  3,
		},
	}

	c, err := json.Marshal(q)
	if err != nil {
		panic(err)
	}

	s.queueMutex.Lock()

	res = s.newRequest(http.MethodPut, "/api/v0/player/queue", bytes.NewReader(c))

	assert.Equal(http.StatusBadRequest, res.Code, "the request should be rejected because one position"+
		" is repeated, returning a 400 HTTP status code")
	assert.Equal("text/plain; charset=utf-8", res.Header().Get("Content-Type"), "the response"+
		" should not contain 'Content-Type' headers'")

	// Set the correct position
	q[1].Position = 1

	c, err = json.Marshal(q)
	if err != nil {
		panic(err)
	}

	res = s.newRequest(http.MethodPut, "/api/v0/player/queue", bytes.NewReader(c))

	assert.Equal(http.StatusCreated, res.Code, "the request should be processed correctly, returning"+
		" a 201 HTTP status code")
	assert.Equal("", res.Header().Get("Content-Type"), "if the response is successful"+
		" should not contain 'Content-Type' headers'")
	assert.Equal("/api/v0/player/queue", res.Header().Get("Location"), "the response should return"+
		" a Location header with the respective value indicating the location of the recently added resource")

	localQ := _playerSync.GetQueue()

	// Overwrite the IDs to avoid a false positive.
	for i := range localQ.Content {
		localQ.Content[i].ID = 0
	}

	assert.Equal(q, localQ.Content, "the queue should be inserted correctly")
	/* -------------------------------------------------------------------------- */

	/* -------------------------- Test GET requests -------------------------- */
	res = s.newRequest(http.MethodGet, "/api/v0/player/queue", nil)

	assert.Equal(http.StatusOK, res.Code, "the request should be processed correctly, returning"+
		" a 200 HTTP status code")
	assert.Equal("application/json", res.Header().Get("Content-Type"), "if the response is successful"+
		" should contain 'Content-Type' headers'")

	returnedQBytes, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var returnedQ []psync.QueueEpisode
	err = json.Unmarshal(returnedQBytes, &returnedQ)
	if err != nil {
		panic(err)
	}

	// Remove the ID to avoid inconsistencies (because we won't be sure of which test is executed first)
	for i := range returnedQ {
		returnedQ[i].ID = 0
	}

	assert.Equal(q, returnedQ, "the returned queue should be the same as the original")
	/* -------------------------------------------------------------------------- */

	/* -------------------------- Test DELETE requests -------------------------- */
	res = s.newRequest(http.MethodDelete, "/api/v0/player/queue", nil)

	assert.Equal(http.StatusNoContent, res.Code, "the request should be processed correctly, returning"+
		" a 204 HTTP status code")
	assert.Equal("", res.Header().Get("Content-Type"), "if the response is successful"+
		" should not contain 'Content-Type' headers'")

	localQ = _playerSync.GetQueue()

	assert.Equal([]psync.QueueEpisode{}, localQ.Content, "the queue should be removed correctly")
	/* -------------------------------------------------------------------------- */

	s.queueMutex.Unlock()
}

func (s *HandlersTestSuite) TestAddToQueueHandler() {
	assert := assert2.New(s.T())

	res := s.newRequest(http.MethodGet, "/api/v0/player/queue/add", nil)

	assert.Equal(http.StatusNotFound, res.Code, "the usage of an incorrect method should return"+
		" a 404 HTTP status code")
	assert.Equal("", res.Header().Get("Content-Type"), "the response should not contain"+
		" the 'Content-Type' header")

	ep := psync.QueueEpisode{
		ID:        0,
		PodcastID: 88,
		EpisodeID: "1ab36235-ddc2-4ee3-b2a7-382b379f0cba",
		Position:  0,
	}

	epB, err := json.Marshal(ep)
	if err != nil {
		panic(err)
	}

	s.queueMutex.Lock()
	defer s.queueMutex.Unlock()

	res = s.newRequest(http.MethodPost, "/api/v0/player/queue/add", bytes.NewReader(epB))

	assert.Equal(http.StatusBadRequest, res.Code, "if the request does not have the 'append' variable, it should"+
		" return a 400 HTTP status code")
	assert.Equal("text/plain; charset=utf-8", res.Header().Get("Content-Type"), "the response should not contain"+
		" the 'Content-Type' header")

	res = s.newRequest(http.MethodPost, "/api/v0/player/queue/add?append=1", bytes.NewReader(epB))

	assert.Equal(http.StatusCreated, res.Code, "the request should be processed correctly, returning"+
		" a 201 HTTP status code")
	assert.Equal("application/json", res.Header().Get("Content-Type"), "if the response is successful"+
		" should contain the appropriate 'Content-Type' headers'")
	assert.Equal("/api/v0/player/queue", res.Header().Get("Location"), "the response should return"+
		" a Location header with the respective value indicating the location of the recently added resource")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	assert.NotEmpty(body, "the body must not be empty and should contain the ID of the recently added episode")

	storedQ := _playerSync.GetQueue().Content

	var qEp psync.QueueEpisode
	for _, e := range storedQ {
		if ep.EpisodeID == e.EpisodeID {
			qEp = e
		}
	}

	// Remove the ID and Position to avoid a false positive when comparing the structures.
	qEp.ID = 0
	qEp.Position = 0

	assert.Equal(ep, qEp, "the episode should be correctly added to the queue stored in the database")

	/**
	 * We are not checking if the episode has been added at the end or the beginning of the queue
	 * because other tests are taking care of that task.
	 */
}

func (s *HandlersTestSuite) TestDelFromQueueHandler() {

}

func (s *HandlersTestSuite) parseProgressReq(body *bytes.Buffer, progressVar *psync.CurrentProgress) error {
	b, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}

	return json.Unmarshal(b, progressVar)
}

func (s *HandlersTestSuite) newRequest(method, target string, body io.Reader) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, body)
	newRouter(false, false).ServeHTTP(res, req)

	return res
}

func (s *HandlersTestSuite) AfterTest(_, _ string) {}

func (s *HandlersTestSuite) TearDownTest() {
	_ = _podcastsDB.Close()
	err := os.RemoveAll(s.podcastsDBPath)
	if err != nil {
		panic(err)
	}
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

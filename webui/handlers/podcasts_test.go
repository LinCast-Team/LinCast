package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"lincast/database"
	"lincast/models"
	"lincast/podcasts"

	testUtils "lincast/utils/testing"

	"github.com/mmcdole/gofeed"
	assert2 "github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSubscribeToPodcastHandler(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db)

	method := "POST"
	body := struct {
		Url string `json:"url"`
	}{
		Url: "https://gotime.fm/rss",
	}

	r := testUtils.NewRequest(mng.SubscribeToPodcastHandler, method, "", testUtils.NewBody(t, body))

	assert.Equal(http.StatusCreated, r.StatusCode, "The status code returned on the subscription of a new podcast should be 201 Created")
	assert.Equal("", r.Header.Get("Content-Type"), "Since the response should not return a body, the 'Content-Type' headers should not be there")

	r = testUtils.NewRequest(mng.SubscribeToPodcastHandler, method, "", testUtils.NewBody(t, body))

	assert.Equal(http.StatusNoContent, r.StatusCode, "The status code returned on the subscription of a podcast that is already on the database should be 204 No Content")
	assert.Equal("", r.Header.Get("Content-Type"), "Since the response should not return a body, the 'Content-Type' headers should not be there")

	body.Url = "abc123"
	r = testUtils.NewRequest(mng.SubscribeToPodcastHandler, method, "", testUtils.NewBody(t, body))

	assert.Equal(http.StatusBadRequest, r.StatusCode, "The status code returned if the feed provided is invalid should be 400 Bad Request")
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"), "Since the response should contain an error msg in plain text, the "+
		"'Content-Type' headers should be 'text/plain; charset=utf-8'")
}

func TestUnsubscribeToPodcastHandler(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db)

	addPodcastToDB("https://gotime.fm/rss", true, db, t) // ID: 1
	id := 1

	method := "PUT"

	// Usage of an ID that is supposed to exist
	r := testUtils.NewRequest(mng.UnsubscribeToPodcastHandler, method, "?id="+fmt.Sprint(id), testUtils.NewBody(t, nil))

	assert.Equal(http.StatusNoContent, r.StatusCode, "The status code returned on the unsubscription of a podcast should be 204 No Content")
	assert.Equal("", r.Header.Get("Content-Type"), "Since the response should not return a body, the 'Content-Type' headers should not be there")

	var pFromDB models.Podcast
	res := db.First(&pFromDB, id)
	if res.Error != nil {
		assert.FailNow(err.Error())
	}

	assert.False(pFromDB.Subscribed, "The subscription of the podcast should be altered correctly")

	// Usage of an ID that does not exist
	r = testUtils.NewRequest(mng.UnsubscribeToPodcastHandler, method, "?id="+fmt.Sprint(10), testUtils.NewBody(t, nil))

	assert.Equal(http.StatusBadRequest, r.StatusCode, "The status code returned when trying to unsubscribe of a podcast with an ID that does not exist should be 400 Bad Request")
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"), "Since the response should contain an error msg in plain text, the 'Content-Type'"+
		" headers should be 'text/plain; charset=utf-8'")
}

func TestGetUserPodcastsHandler(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db)

	feeds := map[string]bool{
		"https://gotime.fm/rss":                     true,
		"https://rustacean-station.org/podcast.rss": false,
		"https://realpython.com/podcasts/rpp/feed":  true,
	}

	for k, v := range feeds {
		addPodcastToDB(k, v, db, t)
	}

	method := "GET"

	r := testUtils.NewRequest(mng.GetUserPodcastsHandler, method, "", testUtils.NewBody(t, nil))

	assert.Equal(http.StatusOK, r.StatusCode, "The status code returned when returning the user's subscriptions should be 200 OK")
	assert.Equal("application/json", r.Header.Get("Content-Type"), "Since the response should have a body with the requested data (json), the 'Content-Type' headers should be 'application/json'")

	var userPodcasts []models.Podcast
	err = json.NewDecoder(r.Body).Decode(&userPodcasts)
	if err != nil {
		panic(err)
	}

	assert.Len(userPodcasts, 2, "There are only 2 podcasts with an active subscription")

	for _, p := range userPodcasts {
		if assert.NotNil(p, "There shouldn't be nil subscriptions") {
			assert.True(p.Subscribed, "Only subscribed podcasts should be returned")
		}
	}
}

func TestGetPodcastHandler(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db)

	url := "https://gotime.fm/rss"
	method := "GET"

	parsedFeed, _, err := podcasts.GetPodcastData(url)
	if err != nil {
		panic(err)
	}

	addOfflinePodcastToDB(parsedFeed, db, t)

	vars := map[string]string{
		"id": fmt.Sprint(parsedFeed.ID),
	}
	r := testUtils.NewRequestWithVars(mng.GetPodcastHandler, method, "/api/v0/podcasts/{id:[0-9]+}/details", vars, testUtils.NewBody(t, nil))

	assert.Equal(http.StatusOK, r.StatusCode, "The response should have the status code 200")
	assert.Equal("application/json", r.Header.Get("Content-Type"), "Since the response should have a body with the requested data (json), the 'Content-Type' headers should be 'application/json'")

	var receivedData models.Podcast
	err = json.NewDecoder(r.Body).Decode(&receivedData)
	if err != nil {
		panic(err)
	}

	// Check data of type time.Time independently, since it will throw a false positive (metadata diff).
	if assert.True(parsedFeed.Added.Equal(receivedData.Added)) {
		receivedData.Added = parsedFeed.Added
	}
	if assert.True(parsedFeed.LastCheck.Equal(receivedData.LastCheck)) {
		receivedData.LastCheck = parsedFeed.LastCheck
	}
	if assert.True(parsedFeed.Model.CreatedAt.Equal(receivedData.Model.CreatedAt)) {
		receivedData.Model.CreatedAt = parsedFeed.Model.CreatedAt
	}
	if assert.True(parsedFeed.Model.UpdatedAt.Equal(receivedData.Model.UpdatedAt)) {
		receivedData.Model.UpdatedAt = parsedFeed.Model.UpdatedAt
	}

	assert.Equal(*parsedFeed, receivedData, "The received data about the podcast should be the same as the stored one")

	// Usage of an ID that is not an integer
	vars = map[string]string{
		"id": "abc",
	}
	r = testUtils.NewRequestWithVars(mng.GetPodcastHandler, method, "/api/v0/podcasts/{id:[0-9]+}/details", vars, testUtils.NewBody(t, nil))

	assert.Equal(http.StatusBadRequest, r.StatusCode)
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"))

	// Request with a non-existent ID
	id := 10

	vars = map[string]string{
		"id": fmt.Sprint(id),
	}
	r = testUtils.NewRequestWithVars(mng.GetPodcastHandler, method, "/api/v0/podcasts/{id:[0-9]+}/details", vars, testUtils.NewBody(t, nil))

	assert.Equal(http.StatusNotFound, r.StatusCode)
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"))

	// Request without ID
	r = testUtils.NewRequest(mng.GetPodcastHandler, method, "", testUtils.NewBody(t, nil))

	assert.Equal(http.StatusBadRequest, r.StatusCode)
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"))
}

func TestGetEpisodesHandler(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db)

	url := "https://gotime.fm/rss"
	method := "GET"

	parsedFeed, originalFeed, err := podcasts.GetPodcastData(url)
	if err != nil {
		panic(err)
	}

	addOfflinePodcastToDB(parsedFeed, db, t)
	addEpisodesToDB(originalFeed, parsedFeed.ID, db, t)

	var epsFromDB []models.Episode
	res := db.Where("parent_podcast_id", parsedFeed.ID).Find(&epsFromDB)
	if res.Error != nil {
		assert.FailNow("Cannot get the episodes stored on the database: %s", res.Error.Error())
	}

	vars := map[string]string{
		"id": fmt.Sprint(parsedFeed.ID),
	}
	r := testUtils.NewRequestWithVars(mng.GetEpisodesHandler, method, "", vars, testUtils.NewBody(t, nil))

	assert.Equal(http.StatusOK, r.StatusCode, "The response should have the status code 200")
	assert.Equal("application/json", r.Header.Get("Content-Type"), "Since the response should have a body with the requested data (json), the 'Content-Type' headers should be 'application/json'")

	var receivedData []models.Episode
	err = json.NewDecoder(r.Body).Decode(&receivedData)
	if err != nil {
		assert.FailNow(err.Error())
	}

	compareEpisodes(&epsFromDB, &receivedData, t)

	// Usage of an ID that is not an integer
	vars = map[string]string{
		"id": "abc",
	}
	r = testUtils.NewRequestWithVars(mng.GetEpisodesHandler, method, "", vars, testUtils.NewBody(t, nil))

	assert.Equal(http.StatusBadRequest, r.StatusCode)
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"))

	// Request with a non-existent ID
	id := 9999999

	vars = map[string]string{
		"id": fmt.Sprint(id),
	}
	r = testUtils.NewRequestWithVars(mng.GetEpisodesHandler, method, "", vars, testUtils.NewBody(t, nil))

	assert.Equal(http.StatusNotFound, r.StatusCode)
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"))

	// Request without ID
	r = testUtils.NewRequest(mng.GetEpisodesHandler, method, "", testUtils.NewBody(t, nil))

	assert.Equal(http.StatusBadRequest, r.StatusCode)
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"))
}

func TestEpisodeProgressHandler_GET(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db)

	url := "https://feeds.megaphone.fm/darknetdiaries"
	method := "GET"

	parsedFeed, originalFeed, err := podcasts.GetPodcastData(url)
	if err != nil {
		panic(err)
	}

	addOfflinePodcastToDB(parsedFeed, db, t)
	addEpisodesToDB(originalFeed, parsedFeed.ID, db, t)

	podcastID := 1
	episodeID := 5
	expectedProgress := time.Minute * 46

	res := db.Model(&models.Episode{}).Where("id = ?", episodeID).UpdateColumn("current_progress", expectedProgress)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	vars := map[string]string{
		"pID":  fmt.Sprint(podcastID),
		"epID": fmt.Sprint(episodeID),
	}

	r := testUtils.NewRequestWithVars(mng.EpisodeProgressHandler, method, "", vars, testUtils.NewBody(t, nil))

	var response struct {
		Progress time.Duration `json:"progress"`
	}

	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		assert.FailNow(err.Error())
	}

	assert.Equal(http.StatusOK, r.StatusCode)
	assert.Equal("application/json", r.Header.Get("Content-Type"))
	assert.Equal(expectedProgress, response.Progress)

	vars = map[string]string{
		"pID":  fmt.Sprint(999999),
		"epID": fmt.Sprint(999999999),
	}

	r = testUtils.NewRequestWithVars(mng.EpisodeProgressHandler, method, "", vars, testUtils.NewBody(t, nil))

	assert.Equal(http.StatusBadRequest, r.StatusCode)
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"))
}

func addPodcastToDB(feedURL string, subscribed bool, db *gorm.DB, t *testing.T) {
	p, _, err := podcasts.GetPodcastData(feedURL)
	if err != nil {
		assert2.FailNow(t, err.Error())
	}

	p.Subscribed = subscribed

	res := db.Save(p)
	if res.Error != nil {
		assert2.FailNow(t, res.Error.Error())
	}
}

func addOfflinePodcastToDB(p *models.Podcast, db *gorm.DB, t *testing.T) {
	res := db.Save(p)
	if res.Error != nil {
		assert2.FailNow(t, res.Error.Error())
	}
}

func addEpisodesToDB(originalFeed *gofeed.Feed, parentPodcastID uint, db *gorm.DB, t *testing.T) {
	eps, err := podcasts.GetEpisodes(originalFeed)
	if err != nil {
		assert2.FailNow(t, "Error trying to get episodes of the given feed: %s", err.Error())
	}

	for _, e := range *eps {
		e.ParentPodcastID = parentPodcastID

		res := db.Create(&e)
		if res.Error != nil {
			assert2.FailNow(t, "Error when trying to store an episode: %s", res.Error.Error())
		}
	}
}

func compareEpisodes(expected *[]models.Episode, current *[]models.Episode, t *testing.T) {
	assert := assert2.New(t)

	for i := range *expected {
		// Check data of type time.Time independently, since it will throw a false positive (metadata diff).
		if assert.True((*expected)[i].Updated.Equal((*current)[i].Updated), `The field "Updated" of the current episode %d does not match with the original`, i) {
			(*current)[i].Updated = (*expected)[i].Updated
		}

		if assert.True((*expected)[i].Published.Equal((*current)[i].Published), `The field "Published" of the current episode %d does not match with the original`, i) {
			(*current)[i].Published = (*expected)[i].Published
		}

		if assert.True((*expected)[i].Model.UpdatedAt.Equal((*current)[i].Model.UpdatedAt), `The field "Model.UpdatedAt" of the current episode %d does not match with the original`, i) {
			(*current)[i].Model.UpdatedAt = (*expected)[i].Model.UpdatedAt
		}

		if assert.True((*expected)[i].Model.CreatedAt.Equal((*current)[i].Model.CreatedAt), `The field "Model.CreatedAt" of the current episode %d does not match with the original`, i) {
			(*current)[i].Model.CreatedAt = (*expected)[i].Model.CreatedAt
		}
	}

	assert.Equal(*expected, *current, "There's a mismatch between the expected episodes and the received ones")
}

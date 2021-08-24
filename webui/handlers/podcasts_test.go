package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"lincast/database"
	"lincast/models"
	"lincast/podcasts"

	testUtils "lincast/utils/testing"

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
	id := 1

	parsedFeed, _, err := podcasts.GetPodcastData(url)
	if err != nil {
		panic(err)
	}

	addOfflinePodcastToDB(parsedFeed, db, t) // ID: 1

	vars := map[string]string{
		"id": fmt.Sprint(id),
	}
	r := testUtils.NewRequestWithVars(mng.GetPodcastHandler, method, "/api/v0/podcasts/{id:[0-9]+}/details", vars, testUtils.NewBody(t, nil))

	assert.Equal(http.StatusOK, r.StatusCode, "The response should have the status code 200")
	assert.Equal("application/json", r.Header.Get("Content-Type"), "Since the response should have a body with the requested data (json), the 'Content-Type' headers should be 'application/json'")

	var receivedData models.Podcast
	err = json.NewDecoder(r.Body).Decode(&receivedData)
	if err != nil {
		panic(err)
	}

	// Check time data independently, since it will throw a false positive (metadata diff).
	if assert.True(parsedFeed.Added.Equal(receivedData.Added)) {
		receivedData.Added = parsedFeed.Added
	}
	if assert.True(parsedFeed.LastCheck.Equal(receivedData.LastCheck)) {
		receivedData.LastCheck = parsedFeed.LastCheck
	}

	assert.Equal(*parsedFeed, receivedData, "The received data about the podcast should be the same as the stored one")

	// Request with a non-existent ID
	id = 10

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

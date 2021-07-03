package handlers

import (
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
	assert.Equal(-1, r.ContentLength, "The response should not have a content") // -1 means that the content length is unknown

	r = testUtils.NewRequest(mng.SubscribeToPodcastHandler, method, "", testUtils.NewBody(t, body))

	assert.Equal(http.StatusNoContent, r.StatusCode, "The status code returned on the subscription of a podcast that is already on the database should be 204 No Content")
	assert.Equal("", r.Header.Get("Content-Type"), "Since the response should not return a body, the 'Content-Type' headers should not be there")
	assert.Equal(-1, r.ContentLength, "The response should not have a content") // -1 means that the content length is unknown

	body.Url = "abc123"
	r = testUtils.NewRequest(mng.SubscribeToPodcastHandler, method, "", testUtils.NewBody(t, body))

	assert.Equal(http.StatusBadRequest, r.StatusCode, "The status code returned if the feed provided is invalid should be 400 Bad Request")
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"), "Since the response should contain an error msg in plain text, the "+
		"'Content-Type' headers should be 'text/plain; charset=utf-8'")
	assert.NotEqual(-1, r.ContentLength, "The response should have a content, since the body must contain the description of the error")
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
	assert.Equal(-1, r.ContentLength, "The response should not have a content") // -1 means that the content length is unknown

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
	assert.NotEqual(-1, r.ContentLength, "The response should have a content, since the body must contain the description of the error")
}

func addPodcastToDB(feedURL string, subscribed bool, db *gorm.DB, t *testing.T) {
	p, _, err := podcasts.GetPodcastData(feedURL)
	if err != nil {
		assert2.FailNow(t, err.Error())
	}

	p.Subscribed = subscribed

	res := db.Save(p)
	if res.Error != nil {
		assert2.FailNow(t, err.Error())
	}
}

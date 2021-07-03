package handlers

import (
	"lincast/database"
	"net/http"
	"testing"

	testUtils "lincast/utils/testing"

	assert2 "github.com/stretchr/testify/assert"
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
	assert.Equal(0, r.ContentLength, "The ContentLength of the response should be 0, since it should not have a body")

	r = testUtils.NewRequest(mng.SubscribeToPodcastHandler, method, "", testUtils.NewBody(t, body))

	assert.Equal(http.StatusNoContent, r.StatusCode, "The status code returned on the subscription of a podcast that is already on the database should be 204 No Content")
	assert.Equal("", r.Header.Get("Content-Type"), "Since the response should not return a body, the 'Content-Type' headers should not be there")
	assert.Equal(0, r.ContentLength, "The ContentLength of the response should be 0, since it should not have a body")

	body.Url = "abc123"
	r = testUtils.NewRequest(mng.SubscribeToPodcastHandler, method, "", testUtils.NewBody(t, body))

	assert.Equal(http.StatusBadRequest, r.StatusCode, "The status code returned if the feed provided is invalid should be 400 Bad Request")
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"), "Since the response should contain an error msg in plain text, the "+
		"'Content-Type' headers should be 'text/plain; charset=utf-8'")
	assert.NotEqual(0, r.ContentLength, "The ContentLength of the response should not be 0, since the body must contain the description of the error")
}

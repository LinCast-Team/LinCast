package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"lincast/database"
	"lincast/models"
	testUtils "lincast/utils/testing"

	assert2 "github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPlayerPlaybackInfoHandler_GET(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db)
	method := "GET"

	// If nothing is being played, an error should be returned
	r := testUtils.NewRequest(mng.PlayerPlaybackInfoHandler, method, "", testUtils.NewBody(t, nil))

	assert.Equal(http.StatusNotFound, r.StatusCode, "The status code returned when nothing is being played should be 404 Not Found")
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"), "Since the response should contain the description of the error, the expected 'Content-Type' headers are 'text/plain; charset=utf-8'")

	// Check if the playback info is correctly returned.
	expectedProgress := models.PlaybackInfo{
		PodcastID: 11,
		EpisodeID: 123,
	}

	res := db.Save(&expectedProgress)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	r = testUtils.NewRequest(mng.PlayerPlaybackInfoHandler, method, "", testUtils.NewBody(t, nil))

	var receivedProgress models.PlaybackInfo
	err = json.NewDecoder(r.Body).Decode(&receivedProgress)
	if err != nil {
		assert.FailNow(err.Error())
	}

	// Remove the fields that are not included on the response
	expectedProgress.Model = gorm.Model{}

	assert.Equal(http.StatusOK, r.StatusCode, "Since the playback info should be returned without problems, the status code must be 200 OK")
	assert.Equal("application/json", r.Header.Get("Content-Type"), "Since the response should contain the playback info, the 'Content-Type' haders should have the value of 'application/json'")
	assert.Equal(expectedProgress, receivedProgress, "The returned progress should be the same as the stored one")
}

func TestPlayerPlaybackInfoHandler_PUT(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db)
	method := "PUT"

	expectedProgress := models.PlaybackInfo{
		PodcastID: 10,
		EpisodeID: 123,
	}

	// Send the request to update the progress of the player
	r := testUtils.NewRequest(mng.PlayerPlaybackInfoHandler, method, "", testUtils.NewBody(t, &expectedProgress))

	// Get the progress of the player stored in the database to check if it was updated correctly
	var progressInDB models.PlaybackInfo
	res := db.First(&progressInDB)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	// Remove the fields that are not used
	progressInDB.Model = gorm.Model{}

	assert.Equal(http.StatusCreated, r.StatusCode, "Since the progress should be updated without problems, the expected status code in the response is 201 Created")
	assert.Equal("", r.Header.Get("Content-Type"), "Since the response should not have a body, the 'Content-Type' headers should be empty")
	assert.Equal(expectedProgress, progressInDB, "The progress of the player should be updated correctly")

	wrongBody := "{'something': 'else', '1': '2', 'foo': 'bar'}"
	r = testUtils.NewRequest(mng.PlayerPlaybackInfoHandler, method, "", testUtils.NewBody(t, &wrongBody))

	assert.Equal(http.StatusBadRequest, r.StatusCode, "Since the request has an unexpected content, it should be rejected by using the HTTP status code 400 Bad Request")
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"), "Since the response should contain the description of the error, the expected 'Content-Type' headers are 'text/plain; charset=utf-8'")
}

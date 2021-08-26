package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"lincast/database"
	"lincast/models"
	testUtils "lincast/utils/testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestPlayerProgressHandler_GET(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db)
	method := "GET"

	// There should be allways an empty row in the database
	var firstResult models.CurrentProgress
	var expectedProgress models.CurrentProgress
	res := db.First(&firstResult)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	// The unique field that should variate is the ID, so we avoid that
	// expected difference by setting the same on both
	expectedProgress.ID = firstResult.ID

	// Remove the filds of type time.Time, since we don't need them (and it would cause a false positive)
	expectedProgress.Model.CreatedAt = time.Time{}
	expectedProgress.Model.UpdatedAt = time.Time{}
	firstResult.Model.CreatedAt = time.Time{}
	firstResult.Model.UpdatedAt = time.Time{}

	assert.Equal(expectedProgress, firstResult, "The db should store an empty progress")

	// Now, we should check if the stored progress is correctly returned from the handler.
	expectedProgress = models.CurrentProgress{
		PodcastID:   1,
		EpisodeGUID: "guid-123",
		Progress:    time.Duration(time.Minute * 87),
	}
	expectedProgress.ID = 1

	res = db.Model(&models.CurrentProgress{}).Where("id = ?", expectedProgress.ID).Updates(&expectedProgress)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	r := testUtils.NewRequest(mng.PlayerProgressHandler, method, "", testUtils.NewBody(t, nil))

	var receivedProgress models.CurrentProgress
	err = json.NewDecoder(r.Body).Decode(&receivedProgress)
	if err != nil {
		assert.FailNow(err.Error())
	}

	// Remove the filds of type time.Time, since we don't need them (and it would cause a false positive)
	receivedProgress.Model.CreatedAt = time.Time{}
	receivedProgress.Model.UpdatedAt = time.Time{}

	assert.Equal(http.StatusOK, r.StatusCode, "Since the current progress should be returned without problems, the status code must be 200 OK")
	assert.Equal("application/json", r.Header.Get("Content-Type"), "Since the response should contain the current progress of the player, the 'Content-Type' haders should have the value of 'application/json'")
	assert.Equal(expectedProgress, receivedProgress, "The returned progress should be the same as the stored one")
}

func TestPlayerProgressHandler_PUT(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db)
	method := "PUT"

	expectedProgress := models.CurrentProgress{
		PodcastID:   10,
		EpisodeGUID: "some_guid-123",
		Progress:    time.Duration(time.Minute * 30),
	}
	expectedProgress.ID = 1

	// Send the request to update the progress of the player
	r := testUtils.NewRequest(mng.PlayerProgressHandler, method, "", testUtils.NewBody(t, &expectedProgress))

	// Get the progress of the player stored in the database to check if it was updated correctly
	var progressInDB models.CurrentProgress
	res := db.First(&progressInDB)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	// Remove the filds of type time.Time, since we don't need them (and it would cause a false positive)
	progressInDB.Model.CreatedAt = time.Time{}
	progressInDB.Model.UpdatedAt = time.Time{}

	assert.Equal(http.StatusCreated, r.StatusCode, "Since the progress should be updated without problems, the expected status code in the response is 201 Created")
	assert.Equal("", r.Header.Get("Content-Type"), "Since the response should not have a body, the 'Content-Type' headers should be empty")
	assert.Equal(expectedProgress, progressInDB, "The progress of the player should be updated correctly")

	wrongBody := "{'something': 'else', '1': '2', 'foo': 'bar'}"
	r = testUtils.NewRequest(mng.PlayerProgressHandler, method, "", testUtils.NewBody(t, &wrongBody))

	assert.Equal(http.StatusBadRequest, r.StatusCode, "Since the request has an unexpected content, it should be rejected by using the HTTP status code 400 Bad Request")
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"), "Since the response should contain the description of the error, the expected 'Content-Type' headers are 'text/plain; charset=utf-8'")
}

package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"lincast/database"
	"lincast/models"
	testUtils "lincast/utils/testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestQueueHandler_GET(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db)
	method := "GET"

	expectedQueue := []models.QueueEpisode{
		{
			Position:  1,
			PodcastID: 10,
			EpisodeID: "guid1",
		},
		{
			Position:  2,
			PodcastID: 1,
			EpisodeID: "guid2",
		},
		{
			Position:  3,
			PodcastID: 18,
			EpisodeID: "guid3",
		},
	}

	res := db.Save(&expectedQueue)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	r := testUtils.NewRequest(mng.QueueHandler, method, "", testUtils.NewBody(t, nil))

	var receivedQueue []models.QueueEpisode
	err = json.NewDecoder(r.Body).Decode(&receivedQueue)
	if err != nil {
		assert.FailNow(err.Error())
	}

	// Check the fields of type time.Time independently to avoid a false positive
	for i := range receivedQueue {
		if assert.True(receivedQueue[i].Model.CreatedAt.Equal(expectedQueue[i].Model.CreatedAt)) {
			receivedQueue[i].Model.CreatedAt = expectedQueue[i].Model.CreatedAt
		}

		if assert.True(receivedQueue[i].Model.UpdatedAt.Equal(expectedQueue[i].Model.UpdatedAt)) {
			receivedQueue[i].Model.UpdatedAt = expectedQueue[i].Model.UpdatedAt
		}
	}

	assert.Equal(http.StatusOK, r.StatusCode)
	// This happens locally (check if the CI haves the same issue): For some reason the contents of the "Content-Type" headers are not being recorded by the http.ResponseRecorder.
	assert.Equal("application/json", r.Header.Get("Content-Type"))
	assert.Equal(expectedQueue, receivedQueue)
}

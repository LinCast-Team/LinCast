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
	"gorm.io/gorm"
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
	assert.Equal("application/json", r.Header.Get("Content-Type"))
	assert.Equal(expectedQueue, receivedQueue)
}

func TestQueueHandler_PUT(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db)
	method := "PUT"

	expectedQueue := []models.QueueEpisode{
		{
			Position:  1,
			PodcastID: 10,
			EpisodeID: "guid1",
			Model: gorm.Model{
				ID: 1,
			},
		},
		{
			Position:  2,
			PodcastID: 1,
			EpisodeID: "guid2",
			Model: gorm.Model{
				ID: 2,
			},
		},
		{
			Position:  3,
			PodcastID: 18,
			EpisodeID: "guid3",
			Model: gorm.Model{
				ID: 3,
			},
		},
	}

	r := testUtils.NewRequest(mng.QueueHandler, method, "", testUtils.NewBody(t, &expectedQueue))

	var queueOnDB []models.QueueEpisode
	res := db.Find(&queueOnDB)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	for i := range queueOnDB {
		queueOnDB[i].Model.CreatedAt = time.Time{}
		queueOnDB[i].Model.UpdatedAt = time.Time{}
	}

	assert.Equal(http.StatusCreated, r.StatusCode)
	assert.Equal("", r.Header.Get("Content-Type"))
	assert.Equal(expectedQueue, queueOnDB)

	// The usage of a repeated position should cause the rejection of the request
	wrongQueue := expectedQueue
	wrongQueue[1].Position = 1
	r = testUtils.NewRequest(mng.QueueHandler, method, "", testUtils.NewBody(t, &wrongQueue))

	assert.Equal(http.StatusBadRequest, r.StatusCode)
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"))

	wrongQueue2 := "[{'foo': 'bar', '1': 2}]"
	r = testUtils.NewRequest(mng.QueueHandler, method, "", testUtils.NewBody(t, &wrongQueue2))

	assert.Equal(http.StatusBadRequest, r.StatusCode)
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"))
}

func TestQueueHandler_DELETE(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db)
	method := "DELETE"

	queueToStore := []models.QueueEpisode{
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

	res := db.Save(&queueToStore)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	r := testUtils.NewRequest(mng.QueueHandler, method, "", testUtils.NewBody(t, nil))

	var queueFromDB []models.QueueEpisode
	res = db.Find(&queueFromDB)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	assert.Equal(http.StatusNoContent, r.StatusCode)
	assert.Equal("", r.Header.Get("Content-Type"), "Since the response should not have a body, the 'Content-Type' headers must be empty")
	assert.Len(queueFromDB, 0, "The queue should be empty")
}

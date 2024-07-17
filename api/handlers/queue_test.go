package handlers

import (
	"encoding/json"
	"fmt"
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
	mng := NewManager(db, make(chan *models.Podcast))
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
	mng := NewManager(db, make(chan *models.Podcast))
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
	mng := NewManager(db, make(chan *models.Podcast))
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

func TestAddToQueueHandler(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db, make(chan *models.Podcast))
	method := "POST"

	baseQueue := []models.QueueEpisode{
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

	res := db.Save(&baseQueue)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	extraEp := models.QueueEpisode{
		Position:  baseQueue[len(baseQueue)-1].Position + 1,
		PodcastID: 99,
		EpisodeID: "guid99",
		Model: gorm.Model{
			ID: baseQueue[len(baseQueue)-1].ID + 1,
		},
	}

	/* Test 1 - Try to add the episode at the end of the queue */

	// Request to add the episode at the end of the queue (append)
	r := testUtils.NewRequest(mng.AddToQueueHandler, method, "?append=1", testUtils.NewBody(t, &extraEp))

	var extraEpFromDB models.QueueEpisode
	res = db.Last(&extraEpFromDB)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	// Remove fields of type time.Time to avoid a false positive (due to metadata diff)
	extraEpFromDB.Model.CreatedAt = time.Time{}
	extraEpFromDB.Model.UpdatedAt = time.Time{}

	assert.Equal(http.StatusCreated, r.StatusCode)
	assert.Equal("application/json", r.Header.Get("Content-Type"))
	assert.Equal("/api/v0/player/queue", r.Header.Get("Location"))
	assert.Equal(extraEp, extraEpFromDB)

	/* Test 2 - Try to add the episode at the beginning of the queue */
	extraEp2 := models.QueueEpisode{
		Position:  1,
		PodcastID: 100,
		EpisodeID: "guid100",
		Model: gorm.Model{
			ID: baseQueue[len(baseQueue)-1].ID + 2,
		},
	}

	// Request to add the episode at the beginning of the queue
	r = testUtils.NewRequest(mng.AddToQueueHandler, method, "?append=0", testUtils.NewBody(t, &extraEp2))

	var extraEp2FromDB models.QueueEpisode
	res = db.Last(&extraEp2FromDB)
	if res.Error != nil {
		assert.FailNow(err.Error())
	}

	// Remove fields of type time.Time to avoid a false positive (due to metadata diff)
	extraEp2FromDB.Model.CreatedAt = time.Time{}
	extraEp2FromDB.Model.UpdatedAt = time.Time{}

	assert.Equal(http.StatusCreated, r.StatusCode)
	assert.Equal("application/json", r.Header.Get("Content-Type"))
	assert.Equal("/api/v0/player/queue", r.Header.Get("Location"))
	assert.Equal(extraEp2, extraEp2FromDB)

	var allEps []models.QueueEpisode
	res = db.Model(&models.QueueEpisode{}).Order("position asc").Find(&allEps)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	for i, e := range allEps {
		assert.Equal(i+1, e.Position, "One episode has the incorrect position (guid '%s')", e.EpisodeID)
	}
}

func TestDelFromQueueHandler(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}
	mng := NewManager(db, make(chan *models.Podcast))
	method := http.MethodDelete

	baseQueue := []models.QueueEpisode{
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

	res := db.Save(&baseQueue)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	idToRemove := baseQueue[1].ID

	expectedQueue := []models.QueueEpisode{
		baseQueue[0],
		baseQueue[2],
	}

	r := testUtils.NewRequest(mng.DelFromQueueHandler, method, "?id="+fmt.Sprint(idToRemove), testUtils.NewBody(t, nil))

	var queueFromDB []models.QueueEpisode
	res = db.Find(&queueFromDB)
	if res.Error != nil {
		assert.FailNow(res.Error.Error())
	}

	// Remove fields of type time.Time to avoid inconsistencies between them
	for i := range queueFromDB {
		queueFromDB[i].Model.CreatedAt = time.Time{}
		queueFromDB[i].Model.UpdatedAt = time.Time{}
		expectedQueue[i].Model.CreatedAt = time.Time{}
		expectedQueue[i].Model.UpdatedAt = time.Time{}
	}

	assert.Equal(http.StatusNoContent, r.StatusCode)
	assert.Equal("", r.Header.Get("Content-Type"))
	assert.Equal(expectedQueue, queueFromDB)

	// Try to remove an episode with a non-existent ID
	r = testUtils.NewRequest(mng.DelFromQueueHandler, method, "?id="+fmt.Sprint(99), testUtils.NewBody(t, nil))

	assert.Equal(http.StatusNotFound, r.StatusCode)
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"))

	// Try to use an ID that is not an integer
	r = testUtils.NewRequest(mng.DelFromQueueHandler, method, "?id=abc", testUtils.NewBody(t, nil))

	assert.Equal(http.StatusBadRequest, r.StatusCode)
	assert.Equal("text/plain; charset=utf-8", r.Header.Get("Content-Type"))
}

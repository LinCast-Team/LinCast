package podcasts

import (
	"github.com/joomcode/errorx"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type DBTestSuite struct {
	dbPath     string
	dbFilename string

	suite.Suite
}

func (s *DBTestSuite) SetupTest() {
	s.dbPath = "./test"
	s.dbFilename = "test.sqlite"

	err := os.Mkdir(s.dbPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (s *DBTestSuite) BeforeTest(_, _ string) {}

func (s *DBTestSuite) TestNewDB() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, s.dbFilename)

	assert.NoError(err, "a new instance of the database should be returned without errors")
	if !assert.NotNil(db.DB, "the returned instance of *sql.DB should not be nil") {
		assert.FailNow("nil instance of db", "test can't continue if the supposed "+
			"valid instance of the database is nil")
	}

	err = db.Close()

	assert.NoError(err, "the database should be closed without errors")

	db, err = NewDB(s.dbPath, "")

	if assert.Error(err, "if the argument 'filename' is empty, an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.IllegalArgument), "the type of the error returned "+
			"should be errorx.IllegalArgument")
	}
	assert.Nil(db, "if there is an error, the instance of the database returned should be nil")
}

func (s *DBTestSuite) TestInsertPodcast() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "insert_podcast_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	p := Podcast{
		Subscribed:  false,
		AuthorName:  "Martin Diaz",
		AuthorEmail: "something@myemail.com",
		Title:       "Random Podcast",
		Description: "Just something random.",
		Categories:  []string{"Tech", "Another Thing", "Sky is Blue"},
		ImageURL:    "https://dogsimages.dog/dog.png",
		ImageTitle:  "A beauty dog",
		Link:        "https://random-podcast1234.org",
		FeedLink:    "https://random-podcast1234.org/feed",
		FeedType:    "rss",
		FeedVersion: "2.0",
		Language:    "en",
		Updated:     time.Now(),
		LastCheck:   time.Now(),
		Added:       time.Now(),
	}

	err = db.InsertPodcast(&p)

	assert.NoError(err, "the podcast should be added without problems")

	_ = db.Close()
}

func (s *DBTestSuite) TestDeletePodcast() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "delete_podcast_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	p := Podcast{
		ID:          1,
		Subscribed:  false,
		AuthorName:  "Martin Diaz",
		AuthorEmail: "something@myemail.com",
		Title:       "Random Podcast",
		Description: "Just something random.",
		Categories:  []string{"Tech", "Another Thing", "Sky is Blue"},
		ImageURL:    "https://dogsimages.dog/dog.png",
		ImageTitle:  "A beauty dog",
		Link:        "https://random-podcast1234.org",
		FeedLink:    "https://random-podcast1234.org/feed",
		FeedType:    "rss",
		FeedVersion: "2.0",
		Language:    "en",
		Updated:     time.Now(),
		LastCheck:   time.Now(),
		Added:       time.Now(),
	}

	err = db.InsertPodcast(&p)
	if err != nil {
		panic(err)
	}

	err = db.DeletePodcast(p.ID)

	assert.NoError(err, "the podcast should be deleted without errors")

	// Try to use a non-existent ID
	err = db.DeletePodcast(858)

	if assert.Error(err, "try to delete a podcast by an ID that doesn't exist should return an error") {
		assert.True(errorx.IsOfType(err, errorx.IllegalArgument), "the returned error should be of "+
			"type errorx.IllegalArgument")
	}

	_ = db.Close()
}

func (s *DBTestSuite) TestGetPodcastByID() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "get_podcast_by_ID_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	p := Podcast{
		ID:          1,
		Subscribed:  false,
		AuthorName:  "Martin Diaz",
		AuthorEmail: "something@myemail.com",
		Title:       "Random Podcast",
		Description: "Just something random.",
		Categories:  []string{"Tech", "Another Thing", "Sky is Blue"},
		ImageURL:    "https://dogsimages.dog/dog.png",
		ImageTitle:  "A beauty dog",
		Link:        "https://random-podcast1234.org",
		FeedLink:    "https://random-podcast1234.org/feed",
		FeedType:    "rss",
		FeedVersion: "2.0",
		Language:    "en",
		Updated:     time.Now(),
		LastCheck:   time.Now(),
		Added:       time.Now(),
	}

	err = db.InsertPodcast(&p)
	if err != nil {
		panic(err)
	}

	retrievedPodcast, err := db.GetPodcastByID(p.ID)

	assert.NoError(err, "the podcast should be obtained without errors")
	if assert.NotNil(retrievedPodcast, "a valid podcast should be retrieved") {
		// We should avoid compare the structures in a general way because the comparison of fields of type
		// time.Time will always cause a failed test. This is due to absence of part of the metadata in the
		// retrieved fields of that type.
		assert.True(p.Updated.Equal(retrievedPodcast.Updated), "the returned value of the field "+
			"Updated should be the same as the original")
		assert.True(p.LastCheck.Equal(retrievedPodcast.LastCheck), "the returned value of the field"+
			" LastCheck should be the same as the original")
		assert.True(p.Added.Equal(retrievedPodcast.Added), "the returned value of the field Added "+
			"should be the same as the original")

		// Clean fields of type time.Time
		p.Updated = time.Time{}
		p.LastCheck = time.Time{}
		p.Added = time.Time{}
		retrievedPodcast.Updated = time.Time{}
		retrievedPodcast.LastCheck = time.Time{}
		retrievedPodcast.Added = time.Time{}

		assert.Equal(p, *retrievedPodcast, "the podcast retrieved should be the same as the inserted one")
	}

	// Try to use a non-existent ID
	retrievedPodcast, err = db.GetPodcastByID(858)

	if assert.Error(err, "try to get a podcast by an ID that doesn't exist should return an error") {
		assert.True(errorx.IsOfType(err, errorx.IllegalArgument), "the returned error should be of "+
			"type errorx.IllegalArgument")
	}
	assert.Nil(retrievedPodcast, "if there is an error, a nil Podcast struct should be returned")

	_ = db.Close()
}

func (s *DBTestSuite) AfterTest(_, _ string) {}

func (s *DBTestSuite) TearDownTest() {
	err := os.RemoveAll(s.dbPath)
	if err != nil {
		panic(err)
	}
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}

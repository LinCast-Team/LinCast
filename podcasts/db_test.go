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
	if !assert.NotNil(db, "the returned instance of *sql.DB should not be nil") {
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

	defer func() {
		_ = db.Close()
	}()

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
}

func (s *DBTestSuite) TestGetAllPodcasts() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "get_all_podcasts_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	podcasts := []Podcast{
		{
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
		},
		{
			ID:          2,
			Subscribed:  false,
			AuthorName:  "Another Author",
			AuthorEmail: "something2@myemail.com",
			Title:       "Random Podcast v2",
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
		},
		{
			ID:          3,
			Subscribed:  false,
			AuthorName:  "Someone Else",
			AuthorEmail: "something3@myemail.com",
			Title:       "Random Podcast v3",
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
		},
	}

	for _, p := range podcasts {
		err = db.InsertPodcast(&p)
		if err != nil {
			panic(err)
		}
	}

	retrievedPodcasts, err := db.GetAllPodcasts()

	assert.NoError(err, "the podcast should be obtained without errors")
	assert.Len(*retrievedPodcasts, len(podcasts), "the number of podcasts returned does not match")
}

func (s *DBTestSuite) TestGetPodcastsBySubscribedStatus() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "get_podcasts_by_subscribed_status"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	podcastsSubscribed := []Podcast{
		{
			ID:          1,
			Subscribed:  true,
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
		},
		{
			ID:          2,
			Subscribed:  true,
			AuthorName:  "Another Author",
			AuthorEmail: "something2@myemail.com",
			Title:       "Random Podcast v2",
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
		},
		{
			ID:          3,
			Subscribed:  true,
			AuthorName:  "Someone Else",
			AuthorEmail: "something3@myemail.com",
			Title:       "Random Podcast v3",
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
		},
	}

	podcastsNotSubscribed := []Podcast{
		{
			ID:          4,
			Subscribed:  false,
			AuthorName:  "Martin Diaz 2",
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
		},
		{
			ID:          5,
			Subscribed:  false,
			AuthorName:  "Another Author 2",
			AuthorEmail: "something2@myemail.com",
			Title:       "Random Podcast v2",
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
		},
		{
			ID:          6,
			Subscribed:  false,
			AuthorName:  "Someone Else 2",
			AuthorEmail: "something3@myemail.com",
			Title:       "Random Podcast v3",
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
		},
	}

	for _, p := range podcastsSubscribed {
		err = db.InsertPodcast(&p)
		if err != nil {
			panic(err)
		}
	}

	for _, p := range podcastsNotSubscribed {
		err = db.InsertPodcast(&p)
		if err != nil {
			panic(err)
		}
	}

	retrievedSubscribedPodcasts, err := db.GetPodcastsBySubscribedStatus(true)

	if !assert.NoError(err, "podcasts should be obtained without errors") {
		assert.FailNow("error not admitted", "tests can't continue with an error")
	}

	for _, p := range *retrievedSubscribedPodcasts {
		var found bool

		for _, p1 := range podcastsSubscribed {
			if p.AuthorName == p1.AuthorName {
				found = true
				// We should avoid compare the structures directly because the comparison of fields of type
				// time.Time will always cause a failed test (in our case). This is due to absence of part of
				// the metadata in the retrieved fields of that type.
				assert.True(p1.Updated.Equal(p.Updated), "the returned value of the field "+
					"Updated should be the same as the original")
				assert.True(p1.LastCheck.Equal(p.LastCheck), "the returned value of the field"+
					" LastCheck should be the same as the original")
				assert.True(p1.Added.Equal(p.Added), "the returned value of the field Added "+
					"should be the same as the original")

				// Clean fields of type time.Time
				p1.Updated = time.Time{}
				p1.LastCheck = time.Time{}
				p1.Added = time.Time{}
				p.Updated = time.Time{}
				p.LastCheck = time.Time{}
				p.Added = time.Time{}

				if assert.Equal(p1, p, "the returned podcast should be the same as the inserted one") {
					break
				}
			}
		}

		if !found {
			assert.Failf("podcast not returned from the database", "the podcast with ID '%d' wasn't"+
				" returned from the database", p.ID)
		}
	}

	retrievedNotSubscribedPodcasts, err := db.GetPodcastsBySubscribedStatus(false)

	if !assert.NoError(err, "podcasts should be obtained without errors") {
		assert.FailNow("error not admitted", "tests can't continue with an error")
	}

	for _, p := range *retrievedNotSubscribedPodcasts {
		var found bool
		for _, p1 := range podcastsNotSubscribed {
			if p.AuthorName == p1.AuthorName {
				found = true
				// We should avoid compare the structures directly because the comparison of fields of type
				// time.Time will always cause a failed test (in our case). This is due to absence of part of the
				// metadata in the retrieved fields of that type.
				assert.True(p1.Updated.Equal(p.Updated), "the returned value of the field "+
					"Updated should be the same as the original")
				assert.True(p1.LastCheck.Equal(p.LastCheck), "the returned value of the field"+
					" LastCheck should be the same as the original")
				assert.True(p1.Added.Equal(p.Added), "the returned value of the field Added "+
					"should be the same as the original")

				// Clean fields of type time.Time
				p1.Updated = time.Time{}
				p1.LastCheck = time.Time{}
				p1.Added = time.Time{}
				p.Updated = time.Time{}
				p.LastCheck = time.Time{}
				p.Added = time.Time{}

				if assert.Equal(p1, p, "the returned podcast should be the same as the inserted one") {
					break
				}
			}
		}

		if !found {
			assert.Failf("podcast not returned from the database", "the podcast with ID '%d' wasn't"+
				" returned from the database", p.ID)
		}
	}
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

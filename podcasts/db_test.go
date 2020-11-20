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
	podcasts   []Podcast
	episodes   Episodes

	suite.Suite
}

func (s *DBTestSuite) SetupTest() {
	s.dbPath = "./test"
	s.dbFilename = "test.sqlite"

	err := os.Mkdir(s.dbPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	s.podcasts = []Podcast{
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
			AuthorName:  "Another Person",
			AuthorEmail: "something@myemail2.com",
			Title:       "Random Podcast 2",
			Description: "Just something random.",
			Categories:  []string{"Tech", "Another Thing", "Cats"},
			ImageURL:    "https://catsimages.dog/cat.png",
			ImageTitle:  "A beauty cat",
			Link:        "https://another-random-podcast1234.org",
			FeedLink:    "https://another-random-podcast1234.org/feed",
			FeedType:    "rss",
			FeedVersion: "2.0",
			Language:    "en",
			Updated:     time.Now(),
			LastCheck:   time.Now(),
			Added:       time.Now(),
		},
	}

	s.episodes = Episodes{
		{
			ID:              1,
			ParentPodcastID: s.podcasts[0].ID,
			Title:           "My First Episode",
			Description:     "This is te description of my awesome episode.",
			Link:            "https://random-podcast1234.org/episodes/1",
			AuthorName:      "Martin Diaz",
			GUID:            "some-random-guid-1234",
			ImageURL:        "https://random-podcast1234.org/assets/episodes/1.png",
			ImageTitle:      "Me wearing a tutu",
			Categories:      []string{"Tech", "Tech News"},
			EnclosureURL:    "https://random-podcast1234.org/assets/episodes/1.mp3",
			EnclosureLength: "14098",
			EnclosureType:   "audio/mp3",
			Season:          "1",
			Published:       time.Time{},
			Played:          false,
			CurrentProgress: "",
		},
		{
			ID:              2,
			ParentPodcastID: s.podcasts[0].ID,
			Title:           "My Second Episode",
			Description:     "This is te description of my second awesome episode.",
			Link:            "https://random-podcast1234.org/episodes/2",
			AuthorName:      "Martin Diaz",
			GUID:            "some-random-guid-12345",
			ImageURL:        "https://random-podcast1234.org/assets/episodes/2.png",
			ImageTitle:      "Me wearing a tutu",
			Categories:      []string{"Tech", "Tech News"},
			EnclosureURL:    "https://random-podcast1234.org/assets/episodes/2.mp3",
			EnclosureLength: "14100",
			EnclosureType:   "audio/mp3",
			Season:          "1",
			Published:       time.Time{},
			Played:          false,
			CurrentProgress: "",
		},
		{
			ID:              3,
			ParentPodcastID: s.podcasts[1].ID,
			Title:           "My Second Episode",
			Description:     "This is te description of my second awesome episode.",
			Link:            "https://another-random-podcast1234.org/episodes/2",
			AuthorName:      "Someone Else",
			GUID:            "some-random-guid-123456",
			ImageURL:        "https://another-random-podcast1234.org/assets/episodes/2.png",
			ImageTitle:      "Me wearing a tutu",
			Categories:      []string{"Tech", "Tech News"},
			EnclosureURL:    "https://another-random-podcast1234.org/assets/episodes/2.mp3",
			EnclosureLength: "14100",
			EnclosureType:   "audio/mp3",
			Season:          "1",
			Published:       time.Time{},
			Played:          false,
			CurrentProgress: "",
		},
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

	err = db.InsertPodcast(&p)

	if assert.Error(err, "if the podcast already exists (has the same FeedLink), an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.RejectedOperation), "the returned error should be of "+
			"type errorx.RejectedOperation")
	}
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

	for _, p := range s.podcasts {
		err = db.InsertPodcast(&p)
		if err != nil {
			panic(err)
		}
	}

	retrievedPodcasts, err := db.GetAllPodcasts()

	assert.NoError(err, "the podcast should be obtained without errors")
	assert.Len(*retrievedPodcasts, len(s.podcasts), "the number of podcasts returned does not match")
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
			Link:        "https://random-podcast12345.org",
			FeedLink:    "https://random-podcast12345.org/feed",
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
			Link:        "https://random-podcast123456.org",
			FeedLink:    "https://random-podcast123456.org/feed",
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
			Link:        "https://random-podcast1234567.org",
			FeedLink:    "https://random-podcast1234567.org/feed",
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
			Link:        "https://random-podcast12345678.org",
			FeedLink:    "https://random-podcast12345678.org/feed",
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
			Link:        "https://random-podcast123456789.org",
			FeedLink:    "https://random-podcast123456789.org/feed",
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

func (s *DBTestSuite) TestPodcastExists() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "podcast_exists_"+s.dbFilename)
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
	if err != nil {
		panic(err)
	}

	exists, err := db.PodcastExists(p.FeedLink)

	assert.NoError(err, "there shouldn't be errors")
	assert.True(exists, "the podcast exists on the database, so the method should reflect that and "+
		"return true")

	exists, err = db.PodcastExists("https://another-feed.org/feed")

	assert.NoError(err, "there shouldn't be errors")
	assert.False(exists, "the podcast doesn't exist on the database, so the method should return false")
}

func (s *DBTestSuite) TestPodcastExistsByID() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "podcast_exists_by_id_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	err = db.InsertPodcast(&s.podcasts[0])
	if err != nil {
		panic(err)
	}

	exists, err := db.PodcastExistsByID(s.podcasts[0].ID)

	assert.NoError(err, "there shouldn't be errors")
	assert.True(exists, "the podcast exists on the database, so the method should reflect that and "+
		"return true")

	exists, err = db.PodcastExistsByID(10)

	assert.NoError(err, "there shouldn't be errors")
	assert.False(exists, "the podcast doesn't exist on the database, so the method should return false")
}

func (s *DBTestSuite) TestSetPodcastSubscription() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "set_podcast_subscription_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	err = db.InsertPodcast(&s.podcasts[0])
	if err != nil {
		panic(err)
	}

	err = db.SetPodcastSubscription(s.podcasts[0].ID, true)

	assert.NoError(err, "the subscription of a podcast should be changed without errors")

	p, err := db.GetPodcastsBySubscribedStatus(true)
	if err != nil {
		panic(err)
	}

	if assert.Len(*p, 1, "the subscribed status should be changed in the database") {
		assert.True((*p)[0].Subscribed, "the subscription should be reflected on the database")
	}

	err = db.SetPodcastSubscription(200, true)

	if assert.Error(err, "if a non recognized ID is used, an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.IllegalArgument), "the returned error should be of"+
			" type errorx.IllegalArgument")
	}
}

func (s *DBTestSuite) TestUpdatePodcastLastCheck() {
	assert := assert2.New(s.T())
	ltTime := time.Now()

	db, err := NewDB(s.dbPath, "update_podcast_last_check_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	err = db.InsertPodcast(&s.podcasts[0])
	if err != nil {
		panic(err)
	}

	err = db.UpdatePodcastLastCheck(s.podcasts[0].ID, ltTime)

	assert.NoError(err, "the field LastCheck should be updated without issues on the database")

	err = db.UpdatePodcastLastCheck(200, ltTime)

	if assert.Error(err, "if an unrecognized ID is used, an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.IllegalArgument), "the returned error should be of"+
			" type errorx.IllegalArgument")
	}

	err = db.UpdatePodcastLastCheck(s.podcasts[0].ID, time.Time{})

	if assert.Error(err, "if a time with value zero is passed as argument, an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.IllegalState), "the returned error should be of"+
			" type errorx.IllegalState")
	}
}

func (s *DBTestSuite) TestEpisodeExists() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "episode_exists_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	err = db.InsertPodcast(&s.podcasts[0])
	if err != nil {
		panic(err)
	}

	err = db.InsertEpisode(&s.episodes[0])
	if err != nil {
		panic(err)
	}

	exists, err := db.EpisodeExists(s.episodes[0].GUID)

	assert.NoError(err, "there shouldn't be an error")
	assert.True(exists, "the method should return true because the episode is already stored in the database")

	exists, err = db.EpisodeExists(s.episodes[0].GUID + "1")

	assert.NoError(err, "there shouldn't be an error")
	assert.False(exists, "the method should return false because the episode is not stored in the database")
}

func (s *DBTestSuite) TestInsertEpisode() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "insert_episode_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	err = db.InsertPodcast(&s.podcasts[0])
	if err != nil {
		panic(err)
	}

	err = db.InsertEpisode(&s.episodes[0])

	assert.NoError(err, "the episode should be inserted into the database without errors")

	err = db.InsertEpisode(&s.episodes[0])

	if assert.Error(err, "the episode shouldn't be inserted into the database because it already exists") {
		assert.True(errorx.IsOfType(err, errorx.RejectedOperation), "the error should be of type"+
			" errorx.RejectedOperation")
	}

	// Use the ID of a non-existent podcast
	ep := s.episodes[0]
	ep.ParentPodcastID = 10
	err = db.InsertEpisode(&ep)

	if assert.Error(err, "the episode shouldn't be inserted into the database because the parent podcast"+
		" doesn't exist") {
		assert.True(errorx.IsOfType(err, errorx.RejectedOperation), "the error should be of type"+
			" errorx.RejectedOperation")
	}
}

func (s *DBTestSuite) TestGetEpisodesByPodcast() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "get_episodes_by_podcast_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	for _, p := range s.podcasts {
		err = db.InsertPodcast(&p)
		if err != nil {
			panic(err)
		}
	}

	for _, e := range s.episodes {
		err = db.InsertEpisode(&e)
		if err != nil {
			panic(err)
		}
	}

	retrievedEps, err := db.GetEpisodesByPodcast(s.podcasts[0].ID)

	assert.NoError(err, "the requested episodes of the podcast with ID '%d' should be returned"+
		" without errors", s.podcasts[0].ID)
	assert.Equal(s.episodes[:2], *retrievedEps, "retrieved episodes should be the same as the inserted ones")

	retrievedEps, err = db.GetEpisodesByPodcast(10)

	if assert.Error(err, "if the podcast ID used as target doesn't exist an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.IllegalArgument), "the returned error must be of type"+
			" errorx.IllegalArgument")
	}
	assert.Nil(retrievedEps, "if there is an error no episodes should be returned")
}

func (s *DBTestSuite) TestSetEpisodePlayed() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "set_episode_played_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	err = db.InsertPodcast(&s.podcasts[0])
	if err != nil {
		panic(err)
	}

	err = db.InsertEpisode(&s.episodes[0])
	if err != nil {
		panic(err)
	}

	err = db.SetEpisodePlayed(s.episodes[0].GUID, true)

	assert.NoError(err, "the episode should be set as played without errors")

	eps, err := db.GetEpisodesByPodcast(s.podcasts[0].ID)
	if err != nil {
		panic(err)
	}

	e := (*eps)[0]
	assert.True(e.Played, "if the episode was set as played, the changes should be "+
		"reflected on the database")

	err = db.SetEpisodePlayed(s.episodes[0].GUID+"1234", true)

	if assert.Error(err, "if the passed GUID does not exist, an error should be returned") {
		assert.True(errorx.IsOfType(err, errorx.IllegalArgument), "the type of the error should "+
			"be errorx.IllegalArgument")
	}
}

func (s *DBTestSuite) TestUpdateEpisodeProgress() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "update_episode_progress_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	err = db.InsertPodcast(&s.podcasts[0])
	if err != nil {
		panic(err)
	}

	err = db.InsertEpisode(&s.episodes[0])
	if err != nil {
		panic(err)
	}

	//              hh:mm:ss
	newProgress := "00:10:50"
	err = db.UpdateEpisodeProgress(newProgress, s.episodes[0].GUID)

	assert.NoError(err, "the episode of the podcast should be updated without problems")

	eps, err := db.GetEpisodesByPodcast(s.podcasts[0].ID)
	if err != nil {
		panic(err)
	}

	assert.Equal(newProgress, (*eps)[0].CurrentProgress)

	cases := []string{
		"00:15",
		"16",
		"00:21:50:09",
		"00:15:90",
		"00:80:10",
		"-80:15:10",
		"This is not supposed to be here",
	}

	for _, c := range cases {
		err = db.UpdateEpisodeProgress(c, s.episodes[0].GUID)

		if assert.Errorf(err, "the usage of a wrong format ('%s') should return an error", c) {
			assert.True(errorx.IsOfType(err, errorx.IllegalFormat), "the returned error should be of"+
				" type errorx.IllegalFormat")
		}
	}
}

func (s *DBTestSuite) TestGetEpisodeUpdated() {
	assert := assert2.New(s.T())

	db, err := NewDB(s.dbPath, "get_episode_updated_"+s.dbFilename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	err = db.InsertPodcast(&s.podcasts[0])
	if err != nil {
		panic(err)
	}

	ep := s.episodes[0]
	ep.Updated = time.Now()

	err = db.InsertEpisode(&ep)
	if err != nil {
		panic(err)
	}

	ut, err := db.GetEpisodeUpdated(ep.GUID)

	assert.NoError(err, "the updated time of the episode should be returned without errors")
	assert.True(ep.Updated.Equal(ut), "the updated time returned should be the same as the"+
		" introduced on the episode")

	ut, err = db.GetEpisodeUpdated("some-random-GUID")

	if assert.Error(err, "the usage of a non existent GUID should return an error") {
		assert.True(errorx.IsOfType(err, errorx.IllegalArgument), "the returned error should be"+
			" of type errorx.IllegalArgument")
	}
	assert.True(ut.IsZero(), "if there is an error a nil time should be returned")
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

package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"lincast/podcasts"

	"github.com/joomcode/errorx"
	_ "github.com/mattn/go-sqlite3" // SQLite3 package
	log "github.com/sirupsen/logrus"
)

var _progressRgx = regexp.MustCompile("^[0-9][0-9]:[0-5][0-9]:[0-5][0-9]$")

type Database struct {
	Path     string
	instance *sql.DB
}

func (db *Database) InsertPodcast(p *podcasts.Podcast) error {
	exists, err := db.PodcastExists(p.FeedLink)
	if err != nil {
		return err
	}

	if exists {
		return errorx.RejectedOperation.New("the podcast with feed '%s' already exists", p.FeedLink)
	}

	query := `
INSERT INTO podcasts (
	subscribed,
	author_name,
    author_email,
	title,
	description,
	categories,
	image_url,
	image_title,
	link,
	feed_link,
	feed_type,
	feed_version,
	lang,
	updated,
	last_check,
	added
) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
`
	stmt, err := db.instance.Prepare(query)
	if err != nil {
		return err
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error(errorx.Decorate(err, "error when trying to execute statement's close method"))
		}
	}()

	categories := strings.Join(p.Categories, ",")

	r, err := stmt.Exec(
		p.Subscribed,
		p.AuthorName,
		p.AuthorEmail,
		p.Title,
		p.Description,
		categories,
		p.ImageURL,
		p.ImageTitle,
		p.Link,
		p.FeedLink,
		p.FeedType,
		p.FeedVersion,
		p.Language,
		p.Updated,
		p.LastCheck,
		p.Added,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return errorx.InternalError.WrapWithNoMessage(err)
	}

	if rowsAffected == 0 {
		return errorx.InternalError.New("no rows has been affected")
	}

	return nil
}

func (db *Database) DeletePodcast(id int) error {
	query := `
DELETE FROM podcasts
WHERE id = ?;
`

	result, err := db.instance.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errorx.IllegalArgument.New("id '%d' does not exist", id)
	}

	return nil
}

func (db *Database) GetPodcastByID(id int) (*podcasts.Podcast, error) {
	query := `
SELECT * FROM podcasts
WHERE id = ?;
`

	row, err := db.instance.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = row.Close()
		if err != nil {
			log.Error(errorx.Decorate(err, "error when trying to close rows"))
		}
	}()

	if !row.Next() {
		return nil, errorx.IllegalArgument.New("can't find the podcast with ID '%d'", id)
	}

	var categories string
	var p podcasts.Podcast

	err = row.Scan(
		&p.ID,
		&p.Subscribed,
		&p.AuthorName,
		&p.AuthorEmail,
		&p.Title,
		&p.Description,
		&categories,
		&p.ImageURL,
		&p.ImageTitle,
		&p.Link,
		&p.FeedLink,
		&p.FeedType,
		&p.FeedVersion,
		&p.Language,
		&p.Updated,
		&p.LastCheck,
		&p.Added,
	)
	if err != nil {
		return nil, err
	}

	p.Categories = strings.Split(categories, ",")

	return &p, nil
}

func (db *Database) GetAllPodcasts() (*[]podcasts.Podcast, error) {
	query := "SELECT * FROM podcasts;"

	rows, err := db.instance.Query(query)
	if err != nil {
		return nil, err
	}

	return db.scanRowsToPodcasts(rows)
}

func (db *Database) GetPodcastsBySubscribedStatus(subscribed bool) (*[]podcasts.Podcast, error) {
	query := `
SELECT * FROM podcasts
WHERE subscribed = ?;
`

	rows, err := db.instance.Query(query, subscribed)
	if err != nil {
		return nil, err
	}

	return db.scanRowsToPodcasts(rows)
}

func (db *Database) PodcastExists(feedURL string) (bool, error) {
	query := `
SELECT feed_link FROM podcasts
WHERE feed_link = ?;
`

	row, err := db.instance.Query(query, feedURL)
	if err != nil {
		return false, err
	}

	defer func() {
		err = row.Close()
		if err != nil {
			log.Error(errorx.Decorate(err, "error when trying to close rows"))
		}
	}()

	return row.Next(), nil
}

func (db *Database) PodcastExistsByID(id int) (bool, error) {
	query := `
SELECT id FROM podcasts
WHERE id = ?;
`

	row, err := db.instance.Query(query, id)
	if err != nil {
		return false, err
	}

	defer func() {
		err = row.Close()
		if err != nil {
			log.Error(errorx.Decorate(err, "error when trying to close rows"))
		}
	}()

	return row.Next(), nil
}

func (db *Database) SetPodcastSubscription(id int, subscribed bool) error {
	query := `
UPDATE podcasts
SET subscribed = ?
WHERE id = ?;
`
	result, err := db.instance.Exec(query, subscribed, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errorx.IllegalArgument.New("podcast with ID '%d' does not exist", id)
	}

	return nil
}

func (db *Database) UpdatePodcastLastCheck(id int, t time.Time) error {
	if t.IsZero() {
		return errorx.IllegalState.New("time (t) can't be zero")
	}

	query := `
UPDATE podcasts
SET last_check = ?
WHERE id = ?;
`
	result, err := db.instance.Exec(query, t, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errorx.IllegalArgument.New("podcast with ID '%d' does not exist", id)
	}

	return nil
}

func (db *Database) EpisodeExists(guid string) (bool, error) {
	query := `
SELECT guid FROM episodes
WHERE guid = ?;
`

	row, err := db.instance.Query(query, guid)
	if err != nil {
		return false, err
	}

	defer func() {
		err = row.Close()
		if err != nil {
			log.Error(errorx.Decorate(err, "error when trying to close rows"))
		}
	}()

	return row.Next(), nil
}

func (db *Database) InsertEpisode(e *podcasts.Episode) error {
	exists, err := db.EpisodeExists(e.GUID)
	if err != nil {
		return err
	}

	if exists {
		return errorx.RejectedOperation.New("the episode with GUID '%s' already exists", e.GUID)
	}

	ppExists, err := db.PodcastExistsByID(e.ParentPodcastID)
	if err != nil {
		return err
	}

	if !ppExists {
		return errorx.RejectedOperation.New("the parent podcast with ID '%d' doesn't exist", e.ParentPodcastID)
	}

	query := `
INSERT INTO episodes (
	parent_podcast_id,
    title,
    description,
	link,
	author_name,
	guid,
	image_url,
	image_title,
	categories,
	enclosure_url,
	enclosure_length,
	enclosure_type,
	season,
	published,
    updated,
	played,
	current_progress
) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
`

	stmt, err := db.instance.Prepare(query)
	if err != nil {
		return err
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error(errorx.Decorate(err, "error when trying to execute statement's close method"))
		}
	}()

	categories := strings.Join(e.Categories, ",")

	r, err := stmt.Exec(
		e.ParentPodcastID,
		e.Title,
		e.Description,
		e.Link,
		e.AuthorName,
		e.GUID,
		e.ImageURL,
		e.ImageTitle,
		categories,
		e.EnclosureURL,
		e.EnclosureLength,
		e.EnclosureType,
		e.Season,
		e.Published,
		e.Updated,
		e.Played,
		e.CurrentProgress,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return errorx.InternalError.WrapWithNoMessage(err)
	}

	if rowsAffected == 0 {
		return errorx.InternalError.New("no rows has been affected")
	}

	return nil
}

func (db *Database) GetEpisodesByPodcast(id int) (*podcasts.Episodes, error) {
	exists, err := db.PodcastExistsByID(id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errorx.IllegalArgument.New("the podcast with ID '%d' doesn't exist", id)
	}

	query := `
SELECT * FROM episodes
WHERE parent_podcast_id = ?;
`
	rows, err := db.instance.Query(query, id)
	if err != nil {
		return nil, err
	}

	return db.scanRowsToEpisodes(rows)
}

func (db *Database) SetEpisodePlayed(guid string, played bool) error {
	query := `
UPDATE episodes
SET played = ?
WHERE guid = ?;
`
	result, err := db.instance.Exec(query, played, guid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errorx.IllegalArgument.New("episode with GUID '%s' does not exist", guid)
	}

	return nil
}

func (db *Database) UpdateEpisodeProgress(newProgress, guid string) error {
	if !_progressRgx.MatchString(newProgress) {
		return errorx.IllegalFormat.New("illegal format on argument newProgress ('%s')", newProgress)
	}

	query := `
UPDATE episodes
SET current_progress = ?
WHERE guid = ?;
`
	result, err := db.instance.Exec(query, newProgress, guid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errorx.IllegalArgument.New("episode with guid '%s' does not exist", guid)
	}

	return nil
}

func (db *Database) GetEpisodeUpdated(guid string) (time.Time, error) {
	query := `
SELECT updated FROM episodes
WHERE guid = ?;
`

	row, err := db.instance.Query(query, guid)
	if err != nil {
		return time.Time{}, err
	}

	defer func() {
		err = row.Close()
		if err != nil {
			log.Error(errorx.Decorate(err, "error when trying to close rows"))
		}
	}()

	if !row.Next() {
		return time.Time{}, errorx.IllegalArgument.New("episode with GUID '%s' does not exist", guid)
	}

	var t time.Time

	err = row.Scan(&t)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

// Close closes the database executing sql.DB.Close().
func (db *Database) Close() error {
	return db.instance.Close()
}

func New(path, filename string) (*Database, error) {
	if filename == "" {
		return nil, errorx.IllegalArgument.New("filename argument can't be an empty string")
	}

	var db Database

	sqlDB, err := initDB(path, filename)
	if err != nil {
		return nil, errorx.Decorate(err, "error when trying to initialize the database '%s'",
			filepath.Join(path, filename))
	}

	db.instance = sqlDB

	err = createTables(db.instance)
	if err != nil {
		return nil, errorx.Decorate(err, "error when trying to create tables")
	}

	db.Path = filepath.Join(path, filename)

	return &db, nil
}

func (db *Database) scanRowsToPodcasts(rows *sql.Rows) (*[]podcasts.Podcast, error) {
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error(errorx.Decorate(err, "error when trying to close rows"))
		}
	}()

	var ps []podcasts.Podcast

	for rows.Next() {
		var p podcasts.Podcast
		var categories string

		err := rows.Scan(
			&p.ID,
			&p.Subscribed,
			&p.AuthorName,
			&p.AuthorEmail,
			&p.Title,
			&p.Description,
			&categories,
			&p.ImageURL,
			&p.ImageTitle,
			&p.Link,
			&p.FeedLink,
			&p.FeedType,
			&p.FeedVersion,
			&p.Language,
			&p.Updated,
			&p.LastCheck,
			&p.Added,
		)
		if err != nil {
			return nil, err
		}

		p.Categories = strings.Split(categories, ",")

		ps = append(ps, p)
	}

	return &ps, nil
}

func (db *Database) scanRowsToEpisodes(rows *sql.Rows) (*podcasts.Episodes, error) {
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error(errorx.Decorate(err, "error when trying to close rows"))
		}
	}()

	var episodes podcasts.Episodes

	for rows.Next() {
		var e podcasts.Episode
		var categories string

		err := rows.Scan(
			&e.ID,
			&e.ParentPodcastID,
			&e.Title,
			&e.Description,
			&e.Link,
			&e.AuthorName,
			&e.GUID,
			&e.ImageURL,
			&e.ImageTitle,
			&categories,
			&e.EnclosureURL,
			&e.EnclosureLength,
			&e.EnclosureType,
			&e.Season,
			&e.Published,
			&e.Updated,
			&e.Played,
			&e.CurrentProgress,
		)
		if err != nil {
			return nil, err
		}

		e.Categories = strings.Split(categories, ",")

		episodes = append(episodes, e)
	}

	return &episodes, nil
}

func initDB(path, filename string) (*sql.DB, error) {
	// Check if the directory is accessible
	dir := filepath.Clean(path)
	_, err := os.Stat(dir)
	if err != nil {
		return nil, errorx.IllegalState.New("the directory '%s' is not accessible", path)
	}

	db, err := sql.Open("sqlite3", filepath.Join(dir, filename))
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errorx.IllegalState.New("instance of sql.DB is nil")
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	query := `
CREATE TABLE IF NOT EXISTS podcasts (
	id 				INTEGER PRIMARY KEY AUTOINCREMENT,
	subscribed 		BOOLEAN NOT NULL DEFAULT false,
	author_name 	TEXT NOT NULL,
	author_email 	TEXT,
	title 			TEXT NOT NULL,
	description 	TEXT,
	categories 		TEXT,
	image_url 		TEXT NOT NULL,
	image_title 	TEXT,
	link 			TEXT,
	feed_link 		TEXT NOT NULL,
	feed_type 		TEXT,
	feed_version 	TEXT,
	lang 			TEXT,
	updated 		DATETIME, 
	last_check 		DATETIME,
	added 			DATETIME
);

CREATE TABLE IF NOT EXISTS episodes (
	id 					INTEGER PRIMARY KEY AUTOINCREMENT,
	parent_podcast_id 	INTEGER NOT NULL,
	title 				TEXT NOT NULL,
	description 		TEXT,
	link 				TEXT NOT NULL,
	author_name 		TEXT NOT NULL,
	guid 				TEXT NOT NULL,
	image_url 			TEXT NOT NULL,
	image_title 		TEXT,
	categories 			TEXT,
	enclosure_url 		TEXT,
	enclosure_length 	TEXT,
	enclosure_type 		TEXT,
	season 				TEXT,
	published 			DATETIME NOT NULL,
	updated 			DATETIME NOT NULL,
	played 				BOOLEAN NOT NULL DEFAULT false,
	current_progress 	TEXT NOT NULL DEFAULT '00:00:00'
);

CREATE TABLE IF NOT EXISTS player_progress (
   id 			 INTEGER PRIMARY KEY CHECK (id = 0),
   progress 	 INTEGER,
   episode_guid  TEXT,
   podcast_id 	 INTEGER,
   user 		 TEXT
);
`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

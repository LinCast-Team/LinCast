package podcasts

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/joomcode/errorx"
	_ "github.com/mattn/go-sqlite3" // SQLite3 package
	log "github.com/sirupsen/logrus"
)

type Database struct {
	Path string
	*sql.DB
}

func (db *Database) InsertPodcast(p *Podcast) error {
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
) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`
	stmt, err := db.Prepare(query)
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

	result, err := db.Exec(query, id)
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

func (db *Database) GetPodcastByID(id int) (*Podcast, error) {
	query := `
SELECT * FROM podcasts
WHERE id = ?;
`

	row, err := db.Query(query, id)
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
	var p Podcast

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

func NewDB(path, filename string) (*Database, error) {
	if filename == "" {
		return nil, errorx.IllegalArgument.New("filename argument can't be an empty string")
	}

	var db Database

	sqlDB, err := initDB(path, filename)
	if err != nil {
		return nil, errorx.Decorate(err, "error when trying to initialize the database '%s'",
			filepath.Join(path, filename))
	}

	db.DB = sqlDB

	err = createTables(db.DB)
	if err != nil {
		return nil, errorx.Decorate(err, "error when trying to create tables")
	}

	db.Path = filepath.Join(path, filename)

	return &db, nil
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
	played 				BOOLEAN NOT NULL DEFAULT false,
	current_progress 	TEXT NOT NULL DEFAULT '00:00:00'
);
`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

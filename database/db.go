package database

import (
	"os"
	"path/filepath"
	"time"

	"lincast/models"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(path, filename string) (*gorm.DB, error) {
	if filename == "" {
		return nil, errorx.IllegalArgument.New("filename argument can't be an empty string")
	}

	// Check if the directory is accessible
	dir := filepath.Clean(path)
	_, err := os.Stat(dir)
	if err != nil {
		return nil, errorx.IllegalState.New("the directory '%s' is not accessible", path)
	}

	dbpath := filepath.Join(dir, filename)

	l := logger.New(
		log.StandardLogger(),
		logger.Config{
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Error,
			SlowThreshold:             time.Second,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{Logger: l})
	if err != nil {
		return nil, errorx.Decorate(err, "error when trying to open the database '%s'",
			filepath.Join(path, filename))
	}

	migrate(db)

	return db, nil
}

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.Podcast{})
	if err != nil {
		log.WithError(errorx.EnsureStackTrace(err)).Panic("error on execution pf the automatic migration of the table " +
			"that contains the podcasts")
	}

	err = db.AutoMigrate(&models.Episode{})
	if err != nil {
		log.WithError(errorx.EnsureStackTrace(err)).Panic("error when executing the automatic migration of the table " +
			"that contains the episodes")
	}

	err = db.AutoMigrate(&models.PlaybackInfo{})
	if err != nil {
		log.WithError(errorx.EnsureStackTrace(err)).Panic("error when executing the automatic migration of the table " +
			"that contains the current progress of the player")
	}

	err = db.AutoMigrate(&models.QueueEpisode{})
	if err != nil {
		log.WithError(errorx.EnsureStackTrace(err)).Panic("error when executing the automatic migration of the table " +
			"that contains the queue of the player")
	}
}

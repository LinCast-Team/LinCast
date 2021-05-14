package database

import (
	"os"
	"path/filepath"
	
	"lincast/models"

	"github.com/joomcode/errorx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		return nil,  errorx.Decorate(err, "error when trying to open the database '%s'",
		filepath.Join(path, filename))
	}

	migrate(db)

	return db, nil
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Podcast{})
	db.AutoMigrate(&models.Episode{})
	db.AutoMigrate(&models.CurrentProgress{})
	db.AutoMigrate(&models.QueueEpisode{})

	// We need to make sure that there will be allways one row on the table that stores
	// the progress of the player, otherwise, we can have issues trying to update the
	// progress on a non-existent row.
	db.FirstOrCreate(&models.CurrentProgress{})
}

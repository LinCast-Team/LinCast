package database

import (
	"lincast/internal/domain/entities"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.User{},
		&entities.Podcast{},
		&entities.Episode{},
		&entities.PlaybackInfo{},
		&entities.QueueEpisode{},
		&entities.EpisodeProgress{},
	)
}

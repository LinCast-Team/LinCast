package repositories

import (
	"lincast/models"

	"gorm.io/gorm"
)

type PodcastRepository interface {
	IRepository[uint, models.Podcast]
}

type podcastRepository struct {
	repository[uint, models.Podcast]
}

func NewPodcastRepository(db *gorm.DB) PodcastRepository {
	return &podcastRepository{
		repository[uint, models.Podcast]{
			db,
		},
	}
}

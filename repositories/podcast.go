package repositories

import (
	"lincast/models"

	"gorm.io/gorm"
)

type PodcastRepository interface {
	GetById(id uint) (*models.Podcast, error)
	Create(podcast models.Podcast) error
	Update(podcast models.Podcast) error
	Delete(id uint) error
}

type podcastRepository struct {
	db *gorm.DB
}

func NewPodcastRepository(db *gorm.DB) PodcastRepository {
	return &podcastRepository{
		db,
	}
}

func (dao *podcastRepository) GetById(id uint) (*models.Podcast, error) {
	var p models.Podcast

	if err := dao.db.First(&p, id).Error; err != nil {
		return nil, err
	}

	return &p, nil
}

func (dao *podcastRepository) Create(podcast models.Podcast) error {
	if err := dao.db.Create(podcast).Error; err != nil {
		return err
	}

	return nil
}

func (dao *podcastRepository) Update(podcast models.Podcast) error {
	if err := dao.db.Save(podcast).Error; err != nil {
		return nil
	}

	return nil
}

func (dao *podcastRepository) Delete(id uint) error {
	if err := dao.db.Delete(&models.Podcast{}, id).Error; err != nil {
		return err
	}

	return nil
}

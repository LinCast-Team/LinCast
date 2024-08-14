package daos

import (
	"lincast/models"

	"gorm.io/gorm"
)

type PodcastDao interface {
	Create(podcast models.Podcast) error
	GetById(id uint) (*models.Podcast, error)
	Update(podcast models.Podcast) error
	Delete(id uint) error
}

type podcastDao struct {
	db *gorm.DB
}

func NewPodcastsDAO(db *gorm.DB) PodcastDao {
	return &podcastDao{
		db,
	}
}

func (dao *podcastDao) GetById(id uint) (*models.Podcast, error) {
	var p models.Podcast

	if err := dao.db.First(&p, id).Error; err != nil {
		return nil, err
	}

	return &p, nil
}

func (dao *podcastDao) Create(podcast models.Podcast) error {
	if err := dao.db.Create(podcast).Error; err != nil {
		return err
	}

	return nil
}

func (dao *podcastDao) Update(podcast models.Podcast) error {
	if err := dao.db.Save(podcast).Error; err != nil {
		return nil
	}

	return nil
}

func (dao *podcastDao) Delete(id uint) error {
	if err := dao.db.Delete(&models.Podcast{}, id).Error; err != nil {
		return err
	}

	return nil
}

package repositories

import (
	"lincast/models"

	"gorm.io/gorm"
)

type QueueRepository interface {
	GetByUser(userId uint) (*[]models.QueueEpisode, error)
	Add(queueEpisode models.QueueEpisode) error
	RemoveEpisode(userID uint, queueEpisodeID uint) error
	RemoveAll(userID uint) error
}

type queueRepository struct {
	db *gorm.DB
}

func NewQueueRepository(db *gorm.DB) QueueRepository {
	return &queueRepository{
		db,
	}
}

func (qr *queueRepository) GetByUser(userId uint) (*[]models.QueueEpisode, error) {
	var q []models.QueueEpisode

	err := qr.db.Find(&q, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (qr *queueRepository) Add(queueEpisode models.QueueEpisode) error {
	err := qr.db.Save(queueEpisode).Error
	if err != nil {
		return err
	}

	return nil
}

func (qr *queueRepository) RemoveEpisode(userID uint, queueEpisodeID uint) error {
	err := qr.db.Delete(&models.QueueEpisode{
		Model: gorm.Model{
			ID: queueEpisodeID,
		},
		UserID: userID,
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func (qr *queueRepository) RemoveAll(userID uint) error {
	err := qr.db.Delete(&models.QueueEpisode{}, "user_id = ?", userID).Error
	if err != nil {
		return err
	}

	return nil
}

package repositories

import (
	"lincast/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QueueRepository interface {
	GetByUser(userId uuid.UUID) (*[]entities.QueueEpisode, error)
	Add(queueEpisode entities.QueueEpisode) error
	RemoveEpisode(userID uuid.UUID, queueEpisodeID uint) error
	RemoveAll(userID uuid.UUID) error
}

type queueRepository struct {
	db *gorm.DB
}

func NewQueueRepository(db *gorm.DB) QueueRepository {
	return &queueRepository{
		db,
	}
}

func (qr *queueRepository) GetByUser(userId uuid.UUID) (*[]entities.QueueEpisode, error) {
	var q []entities.QueueEpisode

	err := qr.db.Find(&q, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (qr *queueRepository) Add(queueEpisode entities.QueueEpisode) error {
	err := qr.db.Save(queueEpisode).Error
	if err != nil {
		return err
	}

	return nil
}

func (qr *queueRepository) RemoveEpisode(userID uuid.UUID, queueEpisodeID uint) error {
	err := qr.db.Delete(&entities.QueueEpisode{
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

func (qr *queueRepository) RemoveAll(userID uuid.UUID) error {
	err := qr.db.Delete(&entities.QueueEpisode{}, "user_id = ?", userID).Error
	if err != nil {
		return err
	}

	return nil
}

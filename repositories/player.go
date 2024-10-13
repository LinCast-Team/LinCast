package repositories

import (
	"lincast/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlayerRepository interface {
	GetByUserId(userId uuid.UUID) (*models.PlaybackInfo, error)
}

type playerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepository{
		db,
	}
}

func (pr *playerRepository) GetByUserId(userID uuid.UUID) (*models.PlaybackInfo, error) {
	p := models.User{
		ID: userID,
	}

	if err := pr.db.Preload("Player").First(&p).Error; err != nil {
		return nil, err
	}

	return &p.Player, nil
}

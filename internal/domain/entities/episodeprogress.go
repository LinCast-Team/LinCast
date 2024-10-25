package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EpisodeProgress struct {
	EpisodeID uint      `json:"episodeID"`
	Episode   Episode   `json:"episode" gorm:"foreignKey:EpisodeID"`
	UserID    uuid.UUID `json:"userID"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
	Progress  uint      `json:"progress"`

	gorm.Model
}

func NewEpisodeProgress(episodeID uint, userID uuid.UUID, progress uint) *EpisodeProgress {
	return &EpisodeProgress{
		EpisodeID: episodeID,
		UserID:    userID,
		Progress:  progress,
	}
}

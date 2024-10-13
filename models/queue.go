package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// QueueEpisode represents an episode of the queue.
type QueueEpisode struct {
	EpisodeID uint      `json:"episodeID"`
	Episode   Episode   `json:"episode" gorm:"foreignKey:EpisodeID"`
	Position  uint      `json:"position"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
	UserID    uuid.UUID `json:"userID"`

	gorm.Model
}

package models

import "gorm.io/gorm"

// QueueEpisode represents an episode of the queue.
type QueueEpisode struct {
	EpisodeID string `json:"episodeID"`
	Position  int    `json:"position"`

	User   User `json:"-" gorm:"foreignKey:UserID"`
	UserID uint `json:"userID"`

	gorm.Model
}

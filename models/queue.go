package models

import "gorm.io/gorm"

// QueueEpisode represents an episode of the queue.
type QueueEpisode struct {
	EpisodeID string `json:"episodeID"`
	Position  int    `json:"position"`
	UserID    uint   `json:"userID"`

	gorm.Model
}

package models

import "gorm.io/gorm"

// QueueEpisode represents an episode of the queue.
type QueueEpisode struct {
	PodcastID int    `json:"podcastID"`
	EpisodeID string `json:"episodeID"`
	Position  int    `json:"position"`

	gorm.Model
}

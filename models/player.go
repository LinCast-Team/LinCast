package models

import (
	"gorm.io/gorm"
)

// PlaybackInfo is the structure used to store and parse the information related with the episode that is being
// played by the player.
type PlaybackInfo struct {
	EpisodeID int `json:"episodeID"`
	UserID    int `json:"userID"`

	gorm.Model `json:"-"`
}

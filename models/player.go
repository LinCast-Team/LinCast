package models

import (
	"time"

	"gorm.io/gorm"
)

// CurrentProgress is the structure used to store and parse the information related with the episode that is being
// currently playing on the player.
type CurrentProgress struct {
	Progress    time.Duration `json:"progress"`
	EpisodeGUID string        `json:"episodeID"`
	PodcastID   int           `json:"podcastID"`

	gorm.Model
}

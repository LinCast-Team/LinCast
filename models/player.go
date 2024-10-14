package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PlaybackInfo is the structure used to store and parse the information related with the episode that is being
// played by the player.
type PlaybackInfo struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primarykey"`
	EpisodeID uint      `json:"episodeID"`
	Episode   Episode   `json:"-" gorm:"foreignKey:EpisodeID"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt time.Time `json:"deletedAt" gorm:"autoDeleteTime"`
}

func (pi *PlaybackInfo) BeforeCreate(tx *gorm.DB) (err error) {
	pi.ID = uuid.New()
	return
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID         `json:"id" gorm:"type:char(36);primaryKey"`
	Username        string            `json:"username" gorm:"unique"`
	PasswordHash    string            `json:"-"`
	PasswordSalt    string            `json:"-"`
	Email           string            `json:"email" gorm:"unique"`
	Name            string            `json:"name"`
	Player          PlaybackInfo      `json:"player"`
	Queue           []QueueEpisode    `json:"queue"`
	EpisodeProgress []EpisodeProgress `json:"episodeProgress"`
	SubscribedTo    []*Podcast        `json:"subscribedTo" gorm:"many2many:subscriptions;"`
	CreatedAt       time.Time         `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time         `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt       time.Time         `json:"deletedAt" gorm:"autoDeleteTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

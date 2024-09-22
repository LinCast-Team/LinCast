package models

import "gorm.io/gorm"

type User struct {
	Username        string            `json:"username" gorm:"unique"`
	PasswordHash    string            `json:"-"`
	PasswordSalt    string            `json:"-"`
	Email           string            `json:"email" gorm:"unique"`
	Name            string            `json:"name"`
	Player          PlaybackInfo      `json:"player"`
	Queue           []QueueEpisode    `json:"queue"`
	EpisodeProgress []EpisodeProgress `json:"episodeProgress"`

	gorm.Model
}

package models

import "gorm.io/gorm"

type EpisodeProgress struct {
	EpisodeID int  `json:"episodeID"`
	UserID    uint `json:"userID"`
	Progress  uint `json:"progress"`

	gorm.Model
}

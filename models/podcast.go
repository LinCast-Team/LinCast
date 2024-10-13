package models

import (
	"time"

	"gorm.io/gorm"
)

// Podcast is the structure that represents a podcast.
type Podcast struct {
	AuthorName      string    `json:"authorName"`
	AuthorEmail     string    `json:"authorEmail"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Categories      string    `json:"categories"`
	ImageURL        string    `json:"imageURL"`
	ImageTitle      string    `json:"imageTitle"`
	Link            string    `json:"link"`
	FeedLink        string    `json:"feedLink" gorm:"unique"`
	FeedType        string    `json:"feedType"`
	FeedVersion     string    `json:"feedVersion"`
	Language        string    `json:"language"`
	Updated         time.Time `json:"updated"` // Mirror of gofeed.Feed.UpdatedParsed
	LastCheck       time.Time `json:"lastCheck"`
	Added           time.Time `json:"added"`
	Episodes        []Episode `json:"episodes"`
	AddedBy         User      `json:"-" gorm:"foreignKey:AddedByID"`
	AddedByID       uint      `json:"addedByID"`
	LastUpdatedBy   User      `json:"-" gorm:"foreignKey:LastUpdatedByID"`
	LastUpdatedByID uint      `json:"lastUpdatedByID"`
	Subscriptions   []*User   `json:"-" gorm:"many2many:subscriptions;"`

	gorm.Model
}

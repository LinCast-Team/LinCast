package models

import (
	"time"

	"gorm.io/gorm"
)

// Podcast is the structure that represents a podcast.
type Podcast struct {
	AuthorName    string    `json:"authorName"`
	AuthorEmail   string    `json:"authorEmail"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Categories    string    `json:"categories"`
	ImageURL      string    `json:"imageURL"`
	ImageTitle    string    `json:"imageTitle"`
	Link          string    `json:"link"`
	FeedLink      string    `json:"feedLink" gorm:"unique"`
	FeedType      string    `json:"feedType"`
	FeedVersion   string    `json:"feedVersion"`
	Language      string    `json:"language"`
	Updated       time.Time `json:"updated"` // Mirror of gofeed.Feed.UpdatedParsed
	LastCheck     time.Time `json:"lastCheck"`
	Added         time.Time `json:"added"`
	Episodes      []Episode `json:"episodes"`

	AddedBy       User      `json:"-" gorm:"foreignKey:AddedByID"`
	AddedByID     uint      `json:"addedByID"`
	
	LastUpdatedBy User      `json:"-" gorm:"foreignKey:LastUpdatedByID"`
	LastUpdatedByID uint     `json:"lastUpdatedByID"`

	Subscriptions []*User    `json:"-" gorm:"many2many:subscriptions;"`

	gorm.Model
}

// Episode is the structure that represent an episode of a podcast.
type Episode struct {
	PodcastID       uint          `json:"podcastID"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	Link            string        `json:"link"`
	AuthorName      string        `json:"authorName"`
	GUID            string        `json:"guid"` // Unique identifier for an item
	ImageURL        string        `json:"imageURL"`
	ImageTitle      string        `json:"imageTitle"`
	Categories      string        `json:"categories"`
	EnclosureURL    string        `json:"enclosureURL"`
	EnclosureLength string        `json:"enclosureLength"`
	EnclosureType   string        `json:"enclosureType"`
	Season          string        `json:"season"`    // Comes from gofeed.Item.ITunesExt.Season - can be empty
	Published       time.Time     `json:"published"` // Mirror of gofeed.Item.PublishedParsed
	Updated         time.Time     `json:"updated"`   // Mirror of gofeed.Item.UpdatedParsed
	Played          bool          `json:"played"`
	CurrentProgress time.Duration `json:"currentProgress"`

	QueuesAddedTo   []QueueEpisode    `json:"queuesAddedTo"`
	BeingPlayedOn   []PlaybackInfo    `json:"beingPlayedOn"`
	EpisodeProgress []EpisodeProgress `json:"episodeProgress"`

	gorm.Model
}

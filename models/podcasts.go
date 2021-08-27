package models

import (
	"time"

	"gorm.io/gorm"
)

// Podcast is the structure that represents a podcast.
type Podcast struct {
	Subscribed  bool
	AuthorName  string
	AuthorEmail string
	Title       string
	Description string
	Categories  string
	ImageURL    string
	ImageTitle  string
	Link        string
	FeedLink    string
	FeedType    string
	FeedVersion string
	Language    string
	Updated     time.Time // Use the field gofeed.Feed.UpdatedParsed
	LastCheck   time.Time
	Added       time.Time

	gorm.Model
}

// Episode is the structure that represent an episode of a podcast.
type Episode struct {
	ParentPodcastID uint
	Title           string
	Description     string
	Link            string
	AuthorName      string
	GUID            string // Unique identifier for an item
	ImageURL        string
	ImageTitle      string
	Categories      string
	EnclosureURL    string
	EnclosureLength string
	EnclosureType   string
	Season          string    // Comes from gofeed.Item.ITunesExt.Season - can be empty
	Published       time.Time // Use the field gofeed.Item.PublishedParsed
	Updated         time.Time // Use the field gofeed.Item.UpdatedParsed
	Played          bool
	CurrentProgress time.Duration

	gorm.Model
}

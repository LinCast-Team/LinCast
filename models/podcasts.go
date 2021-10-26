package models

import (
	"time"

	"gorm.io/gorm"
)

// Podcast is the structure that represents a podcast.
type Podcast struct {
	Subscribed  bool      `json:"subscribed"`
	AuthorName  string    `json:"authorName"`
	AuthorEmail string    `json:"authorEmail"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Categories  string    `json:"categories"`
	ImageURL    string    `json:"imageURL"`
	ImageTitle  string    `json:"imageTitle"`
	Link        string    `json:"link"`
	FeedLink    string    `json:"feedLink"`
	FeedType    string    `json:"feedType"`
	FeedVersion string    `json:"feedVersion"`
	Language    string    `json:"language"`
	Updated     time.Time `json:"updated"` // Use the field gofeed.Feed.UpdatedParsed
	LastCheck   time.Time `json:"lastCheck"`
	Added       time.Time `json:"added"`

	gorm.Model
}

// Episode is the structure that represent an episode of a podcast.
type Episode struct {
	ParentPodcastID uint          `json:"parentPodcastID"`
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
	Published       time.Time     `json:"published"` // Use the field gofeed.Item.PublishedParsed
	Updated         time.Time     `json:"updated"`   // Use the field gofeed.Item.UpdatedParsed
	Played          bool          `json:"played"`
	CurrentProgress time.Duration `json:"currentProgress"`

	gorm.Model
}

package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
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
	AddedByID     uuid.UUID `json:"addedByID" `
	Subscriptions []*User   `json:"-" gorm:"many2many:subscriptions;"`

	gorm.Model
}

func NewPodcast(authorName, authorEmail, title, description, categories, imageURL, imageTitle, link, feedLink, feedType, feedVersion, language string, addedByID uuid.UUID, updated time.Time) Podcast {
	return Podcast{
		AuthorName:  authorName,
		AuthorEmail: authorEmail,
		Title:       title,
		Description: description,
		Categories:  categories,
		ImageURL:    imageURL,
		ImageTitle:  imageTitle,
		Link:        link,
		FeedLink:    feedLink,
		FeedType:    feedType,
		FeedVersion: feedVersion,
		Language:    language,
		AddedByID:   addedByID,
		Updated:     updated,
	}
}

func NewPodcastFromFeed(feed *gofeed.Feed, addedByID uuid.UUID, lastCheck, added time.Time) *Podcast {
	p := NewPodcast(
		feed.Author.Name,
		feed.Author.Email,
		feed.Title,
		feed.Description,
		strings.Join(feed.Categories, ","),
		feed.Image.URL,
		feed.Image.Title,
		feed.Link,
		feed.FeedLink,
		feed.FeedType,
		feed.FeedVersion,
		feed.Language,
		addedByID,
		*feed.UpdatedParsed,
	)

	p.LastCheck = lastCheck
	p.Added = added

	return &p
}

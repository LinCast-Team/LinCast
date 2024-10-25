package entities

import (
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"gorm.io/gorm"
)

// Episode is the structure that represent an episode of a podcast.
type Episode struct {
	PodcastID       uint      `json:"podcastID"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Link            string    `json:"link"`
	AuthorName      string    `json:"authorName"`
	GUID            string    `json:"guid"` // Unique identifier for an item
	ImageURL        string    `json:"imageURL"`
	ImageTitle      string    `json:"imageTitle"`
	Categories      string    `json:"categories"`
	EnclosureURL    string    `json:"enclosureURL"`
	EnclosureLength string    `json:"enclosureLength"`
	EnclosureType   string    `json:"enclosureType"`
	Season          string    `json:"season"`    // Comes from gofeed.Item.ITunesExt.Season - can be empty
	Published       time.Time `json:"published"` // Mirror of gofeed.Item.PublishedParsed
	Updated         time.Time `json:"updated"`   // Mirror of gofeed.Item.UpdatedParsed

	QueuesAddedTo   []QueueEpisode    `json:"queuesAddedTo"`
	BeingPlayedOn   []PlaybackInfo    `json:"beingPlayedOn"`
	EpisodeProgress []EpisodeProgress `json:"episodeProgress"`

	gorm.Model
}

func NewEpisode(podcastID uint, title, description, link, authorName, guid, imageURL, imageTitle, categories, enclosureURL, enclosureLength, enclosureType, season string, published, updated time.Time) *Episode {
	return &Episode{
		PodcastID:       podcastID,
		Title:           title,
		Description:     description,
		Link:            link,
		AuthorName:      authorName,
		GUID:            guid,
		ImageURL:        imageURL,
		ImageTitle:      imageTitle,
		Categories:      categories,
		EnclosureURL:    enclosureURL,
		EnclosureLength: enclosureLength,
		EnclosureType:   enclosureType,
		Season:          season,
		Published:       published,
		Updated:         updated,
	}
}

func NewEpisodeFromItem(podcastID uint, item *gofeed.Item) *Episode {
	return NewEpisode(
		podcastID,
		item.Title,
		item.Description,
		item.Link,
		item.Author.Name,
		item.GUID,
		item.Image.URL,
		item.Image.Title,
		strings.Join(item.Categories, ","),
		item.Enclosures[0].URL,
		item.Enclosures[0].Length,
		item.Enclosures[0].Type,
		item.ITunesExt.Season,
		*item.PublishedParsed,
		*item.UpdatedParsed,
	)
}

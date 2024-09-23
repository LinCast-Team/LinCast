package podcasts

import (
	"strings"
	"time"

	"lincast/models"

	"github.com/joomcode/errorx"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
	log "github.com/sirupsen/logrus"
)

// GetPodcastData returns the data from the feed's URL, doing the parsing of the feed itself (into a struct of type *gofeed.Feed) and the podcast.
// Possible errors:
// 	- errorx.ExternalError: if the request to `feedURL` or the parsing of the response fails.
func GetPodcastData(feedURL string) (parsedPodcast *models.Podcast, originalFeed *gofeed.Feed, err error) {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(feedURL)
	if err != nil {
		return nil, nil, errorx.ExternalError.New("the feed can't be obtained/parsed")
	}

	now := time.Now()

	if feed.UpdatedParsed == nil {
		feed.UpdatedParsed = new(time.Time)
	}

	p := &models.Podcast{
		// Subscribed:  false,
		AuthorName:  feed.Author.Name,
		AuthorEmail: feed.Author.Email,
		Title:       feed.Title,
		Description: feed.Description,
		Categories:  strings.Join(feed.Categories, ","),
		ImageURL:    feed.Image.URL,
		ImageTitle:  feed.Image.Title,
		Link:        feed.Link,
		FeedLink:    feed.FeedLink,
		FeedType:    feed.FeedType,
		FeedVersion: feed.FeedVersion,
		Language:    feed.Language,
		Updated:     *feed.UpdatedParsed,
		LastCheck:   now,
		Added:       now,
	}

	return p, feed, nil
}

// GetEpisodes returns the episodes (struct Episodes) of the given Podcast.
// Possible errors:
// 	- errorx.ExternalError: if the request to `p.FeedLink` or the parsing of the response fails.
func GetEpisodes(feed *gofeed.Feed) (*[]models.Episode, error) {
	var episodes []models.Episode

	for _, item := range feed.Items {
		if len(item.Enclosures) == 0 {
			log.WithFields(log.Fields{
				"podcastFeed": feed.FeedLink,
				"episodeGUID": item.GUID,
				"error": errorx.DataUnavailable.New("the episode (GUID '%s') doesn't have"+
					" enclosures", item.GUID),
			}).Error("Episode with no enclosures")

			continue
		}

		if item.UpdatedParsed == nil {
			item.UpdatedParsed = new(time.Time)
		}

		if item.PublishedParsed == nil {
			item.PublishedParsed = new(time.Time)
		}

		if item.Author == nil {
			item.Author = new(gofeed.Person)
		}

		if item.Image == nil {
			item.Image = new(gofeed.Image)
		}

		if item.ITunesExt == nil {
			item.ITunesExt = new(ext.ITunesItemExtension)
		}

		e := models.Episode{
			Title:           item.Title,
			Description:     item.Description,
			Link:            item.Link,
			Published:       *item.PublishedParsed,
			Updated:         *item.UpdatedParsed,
			AuthorName:      item.Author.Name,
			GUID:            item.GUID,
			ImageURL:        item.Image.URL,
			ImageTitle:      item.Image.Title,
			Categories:      strings.Join(item.Categories, ","),
			EnclosureURL:    item.Enclosures[0].URL,
			EnclosureLength: item.Enclosures[0].Length,
			EnclosureType:   item.Enclosures[0].Type,
			Season:          item.ITunesExt.Season,
			Played:          false,
			CurrentProgress: 0,
		}

		episodes = append(episodes, e)
	}

	// Not needed (I think)
	//sort.SliceStable(episodes, func(i, j int) bool {
	//	return episodes[i].Published.Before(*episodes[j].Published)
	//})

	return &episodes, nil
}

/*
func isValidURL(url1 string) (bool, url.URL) {
	if _, err := url.ParseRequestURI(url1); err != nil {
		return false, url.URL{}
	}

	u, err := url.Parse(url1)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false, url.URL{}
	}

	return true, *u
}
*/

package podcasts

import (
	"net/url"
	"time"

	"github.com/joomcode/errorx"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
)

// Podcast is the structure that represents a podcast.
type Podcast struct {
	ID          int
	Subscribed  bool
	AuthorName  string
	AuthorEmail string
	Title       string
	Description string
	Categories  []string
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
}

// Episode is the structure that represent an episode of a podcast.
type Episode struct {
	ID              int
	ParentPodcastID int
	Title           string
	Description     string
	Link            string
	AuthorName      string
	GUID            string // Unique identifier for an item
	ImageURL        string
	ImageTitle      string
	Categories      []string
	EnclosureURL    string
	EnclosureLength string
	EnclosureType   string
	Season          string    // Comes from gofeed.Item.ITunesExt.Season - can be empty
	Published       time.Time // Use the field gofeed.Item.PublishedParsed
	Updated         time.Time // Use the field gofeed.Item.UpdatedParsed
	Played          bool
	CurrentProgress string
}

// Episodes is a slice of structures of type Episode.
type Episodes []Episode

// GetPodcast returns information of a podcast parsed into a struct of type *gofeed.Feed. The information will be
// obtained from the feed's URL.
// Possible errors:
// 	- errorx.IllegalFormat: if the format of the URL is incorrect.
// 	- errorx.ExternalError: if the request to `feedURL` or the parsing of the response fails.
func GetPodcast(feedURL string) (*Podcast, error) {
	valid, parsedURL := isValidURL(feedURL)
	if !valid {
		return nil, errorx.IllegalFormat.New("the url '%s' is not correctly formatted", feedURL)
	}

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(parsedURL.String())
	if err != nil {
		return nil, errorx.ExternalError.New("the feed can't be obtained/parsed")
	}

	now := time.Now()

	if feed.UpdatedParsed == nil {
		feed.UpdatedParsed = new(time.Time)
	}

	p := &Podcast{
		ID:          0,
		Subscribed:  false,
		AuthorName:  feed.Author.Name,
		AuthorEmail: feed.Author.Email,
		Title:       feed.Title,
		Description: feed.Description,
		Categories:  feed.Categories,
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

	return p, nil
}

// GetEpisodes returns the episodes (struct Episodes) of the given Podcast.
// Possible errors:
// 	- errorx.ExternalError: if the request to `p.FeedLink` or the parsing of the response fails.
func (p *Podcast) GetEpisodes() (*Episodes, error) {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(p.FeedLink)
	if err != nil {
		return nil, errorx.ExternalError.New("the feed can't be obtained/parsed")
	}

	var episodes Episodes
	for _, item := range feed.Items {
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

		if len(item.Enclosures) == 0 {
			item.Enclosures = []*gofeed.Enclosure{
				new(gofeed.Enclosure),
			}
		}

		if item.ITunesExt == nil {
			item.ITunesExt = new(ext.ITunesItemExtension)
		}

		e := Episode{
			ParentPodcastID: p.ID,
			Title:           item.Title,
			Description:     item.Description,
			Link:            item.Link,
			Published:       *item.PublishedParsed,
			Updated:         *item.UpdatedParsed,
			AuthorName:      item.Author.Name,
			GUID:            item.GUID,
			ImageURL:        item.Image.URL,
			ImageTitle:      item.Image.Title,
			Categories:      item.Categories,
			EnclosureURL:    item.Enclosures[0].URL,
			EnclosureLength: item.Enclosures[0].Length,
			EnclosureType:   item.Enclosures[0].Type,
			Season:          item.ITunesExt.Season,
			Played:          false,
			CurrentProgress: "00:00:00",
		}

		episodes = append(episodes, e)
	}

	// Not needed (I think)
	//sort.SliceStable(episodes, func(i, j int) bool {
	//	return episodes[i].Published.Before(*episodes[j].Published)
	//})

	return &episodes, nil
}

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

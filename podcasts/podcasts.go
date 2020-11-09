package podcasts

import (
	"net/url"

	"github.com/joomcode/errorx"
	"github.com/mmcdole/gofeed"
)

// GetPodcast returns information of a podcast parsed into a struct of type *gofeed.Feed. The information will be
// obtained from the feed's URL. Possible errors:
// 	- errorx.IllegalFormat: if the format of the URL is incorrect.
// 	- errorx.ExternalError: if the request to `feedURL` or the parsing of the response fails.
func GetPodcast(feedURL string) (*gofeed.Feed, error) {
	valid, parsedURL := isValidURL(feedURL)
	if !valid {
		return nil, errorx.IllegalFormat.New("the url '%s' is not correctly formatted", feedURL)
	}

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(parsedURL.String())
	if err != nil {
		return nil, errorx.ExternalError.New("the feed can't be obtained/parsed")
	}

	return feed, nil
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

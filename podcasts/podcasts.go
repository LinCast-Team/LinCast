package podcasts

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"

	"github.com/joomcode/errorx"
	"github.com/mmcdole/gofeed"
)

// GetPodcast returns information of a podcast parsed into a struct of type *gofeed.Feed. The information will be
// obtained from the feed's URL. Possible errors:
// 	- errorx.IllegalFormat: if the format of the URL or the feed are incorrect.
// 	- errorx.DataUnavailable: if the request to `feedURL` fails.
func GetPodcast(feedURL string) (*gofeed.Feed, error) {
	valid, parsedURL := isValidURL(feedURL)
	if !valid {
		return nil, errorx.IllegalFormat.New("the url '%s' is not correctly formatted", feedURL)
	}

	res, err := http.Get(parsedURL.String())
	if err != nil {
		return nil, errorx.DataUnavailable.Wrap(err, "the request has failed")
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Error(errorx.Decorate(err, "error when trying to close response's body"))
		}
	}()

	parser := gofeed.NewParser()
	feed, err := parser.Parse(res.Body)
	if err != nil {
		return nil, errorx.IllegalFormat.Wrap(err, "the feed can't be parsed")
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

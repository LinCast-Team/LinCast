package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"lincast/podcasts"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
)

var _podcastsDB *podcasts.Database

func subscribeToPodcastHandler(w http.ResponseWriter, r *http.Request) {
	u := struct {
		URL string `json:"url"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      err.Error(),
		}).Error("Error when trying to decode the body of the request")

		return
	}

	p, err := podcasts.GetPodcast(u.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr":  r.RemoteAddr,
			"requestURI":  r.RequestURI,
			"method":      r.Method,
			"request.url": u.URL,
			"error":       errorx.EnsureStackTrace(err),
		}).Error("Error when trying to decode the body of the request")

		return
	}

	p.Subscribed = true

	err = _podcastsDB.InsertPodcast(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr":    r.RemoteAddr,
			"requestURI":    r.RequestURI,
			"method":        r.Method,
			"error":         errorx.EnsureStackTrace(err),
			"request.url":   u.URL,
			"podcastStruct": fmt.Sprintf("%+v", *p),
		}).Error("Error when trying to decode the body of the request")

		return
	}

	log.WithFields(log.Fields{
		"remoteAddr":  r.RemoteAddr,
		"requestURI":  r.RequestURI,
		"method":      r.Method,
		"request.url": u.URL,
	}).Debug("Request to add a new subscription (and a podcast to the database) processed correctly")

	w.WriteHeader(http.StatusOK)
}

func unsubscribeToPodcastHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.IllegalFormat.New("param 'id' is missing"),
		}).Error("Request rejected due to absence of parameter 'id'")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	idStr := keys[0]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"usedID":     idStr,
			"error":      errorx.EnsureStackTrace(err),
		}).Error("Cannot parse the ID of the podcast to unsubscribe")
	}

	err = _podcastsDB.SetPodcastSubscription(int(id), false)
	if err != nil {
		if errorx.IsOfType(err, errorx.IllegalArgument) {
			w.WriteHeader(http.StatusBadRequest)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"requestURI": r.RequestURI,
				"method":     r.Method,
				"error":      errorx.Decorate(err, "the podcast with the given ID does not exist"),
				"usedID":     id,
			}).Error("Error when trying to change the subscription status of the podcast")
		} else {
			w.WriteHeader(http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"requestURI": r.RequestURI,
				"method":     r.Method,
				"error":      errorx.EnsureStackTrace(err),
				"podcastID":  id,
			}).Error("Error when trying to change the subscription status of the podcast")
		}

		return
	}

	log.WithFields(log.Fields{
		"remoteAddr": r.RemoteAddr,
		"requestURI": r.RequestURI,
		"method":     r.Method,
		"podcastID":  id,
	}).Debug("Request to unsubscribe from a podcast processed correctly")

	w.WriteHeader(http.StatusOK)
}

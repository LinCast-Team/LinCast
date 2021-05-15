package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"lincast/podcasts"
	"lincast/models"
	"lincast/utils/safe"

	"github.com/gorilla/mux"
	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _db *gorm.DB

var SubscribeToPodcastHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	u := struct {
		URL string `json:"url"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      err.Error(),
		}).Error("Error when trying to decode the body of the request")

		return
	}

	p, _, err := podcasts.GetPodcastData(u.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

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

	// TODO avoid possible duplicate of podcasts.
	_db.Create(&p)

	w.WriteHeader(http.StatusCreated)
})

var UnsubscribeToPodcastHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		err := errorx.IllegalFormat.New("param 'id' is missing")

		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      err.Error(),
		}).Error("Request rejected due to absence of parameter 'id'")

		return
	}

	idStr := keys[0]

	id := safe.SafeParseInt(idStr)
	if id == safe.DefaultAllocate {
		err := errorx.IllegalArgument.New("the value '%s' is over the limit of int values", idStr)

		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      err.Error(),
			"givenID":    idStr,
		}).Error("Cannot parse the ID of the podcast to unsubscribe")
	}

	if res := _db.Model(&models.Podcast{}).Where("id = ?", id).Update("subscribed", false); res.Error != nil {
		http.Error(w, "the podcast with the given ID does not exist", http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.Decorate(res.Error, "the podcast with the given ID does not exist"),
			"usedID":     id,
		}).Error("Error when trying to change the subscription status of the podcast")

		return
	}

	w.WriteHeader(http.StatusNoContent)
})

var GetUserPodcastsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var subscribed bool
	var unsubscribed bool

	keys, ok := r.URL.Query()["subscribed"]
	if !ok || len(keys[0]) < 1 {
		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
		}).Warn("Parameter 'subscribed' is not present in the request")
	} else {
		s, err := strconv.ParseBool(keys[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			log.WithFields(log.Fields{
				"remoteAddr":      r.RemoteAddr,
				"requestURI":      r.RequestURI,
				"method":          r.Method,
				"error":           errorx.Decorate(err, "'subscribed' cannot be parsed"),
				"subscribedValue": keys[0],
			}).Error("Error when trying to parse the value of 'subscribed' param")

			return
		}

		subscribed = s
	}

	keys, ok = r.URL.Query()["unsubscribed"]
	if !ok || len(keys[0]) < 1 {
		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
		}).Warn("Parameter 'unsubscribed' is not present in the request")
	} else {
		uns, err := strconv.ParseBool(keys[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			log.WithFields(log.Fields{
				"remoteAddr":        r.RemoteAddr,
				"requestURI":        r.RequestURI,
				"method":            r.Method,
				"error":             errorx.Decorate(err, "'unsubscribed' cannot be parsed"),
				"unsubscribedValue": keys[0],
			}).Error("Error when trying to parse the value of 'unsubscribed' param")

			return
		}

		unsubscribed = uns
	}

	var subscribedPodcasts []models.Podcast
	if subscribed || (!subscribed && !unsubscribed) {
		if res := _db.Where("subscribed = ?", true).Find(&subscribedPodcasts); res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"requestURI": r.RequestURI,
				"method":     r.Method,
				"error":      errorx.EnsureStackTrace(res.Error),
			}).Error("Error when trying to get subscribed podcasts from db")

			return
		}
	}

	var unsubscribedPodcasts []models.Podcast
	if unsubscribed || (!subscribed && !unsubscribed) {
		if res := _db.Where("subscribed = ?", false).Find(&unsubscribedPodcasts); res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"requestURI": r.RequestURI,
				"method":     r.Method,
				"error":      errorx.EnsureStackTrace(res.Error),
			}).Error("Error when trying to get unsubscribed podcasts from db")

			return
		}
	}

	p := map[string][]models.Podcast{
		"subscribed":   subscribedPodcasts,
		"unsubscribed": unsubscribedPodcasts,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.EnsureStackTrace(err),
		}).Error("Error when trying to encode the response to the request")

		return
	}
})

var GetPodcastHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id := safe.SafeParseInt(idStr)
	if id == safe.DefaultAllocate {
		err := errorx.IllegalArgument.New("the value '%s' is over the limit of int values", idStr)

		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      err.Error(),
			"givenID":    idStr,
		}).Error("The given ID cannot be parsed")

		return
	}

	var p models.Podcast

	if res := _db.Where("id = ?", id).First(&p); res.Error != nil {
		http.Error(w, "the podcast with the given ID does not exist", http.StatusNotFound)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.Decorate(res.Error, "the podcast with the given ID does not exist"),
			"givenID":    id,
		}).Error("Error when trying to get the requested podcast")

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.EnsureStackTrace(err),
		}).Error("Error when trying to encode the response to the request")

		return
	}
})

var GetEpisodesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id := safe.SafeParseInt(idStr)
	if id == safe.DefaultAllocate {
		err := errorx.IllegalArgument.New("the value '%s' is over the limit of int values", idStr)

		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      err.Error(),
			"givenID":    idStr,
		}).Error("The given ID cannot be parsed")

		return
	}

	var eps []models.Episode

	if res := _db.Where("parent_podcast_id", id).Find(&eps); res.Error != nil {
		http.Error(w, "the podcast with the given ID does not exist", http.StatusNotFound)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.Decorate(res.Error, "the podcast with the given ID does not exist"),
			"givenID":    id,
		}).Error("Error when trying to get the requested episodes")

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(eps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.EnsureStackTrace(err),
		}).Error("Error when trying to encode the response to the request")

		return
	}
})

package handlers

import (
	"encoding/json"
	"net/http"

	"lincast/models"
	"lincast/podcasts"
	"lincast/utils/safe"

	"github.com/gorilla/mux"
	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
)

func (m *Manager) SubscribeToPodcastHandler(w http.ResponseWriter, r *http.Request) {
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

	// Resolve the given URL first, and then decide if we will save the data or not.
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

	// Check if the feed's URL is already on the database by counting their appearances.
	var alreadyOnDB int64
	res := m.db.Model(&models.Podcast{}).Where("feed_link = ?", p.FeedLink).Count(&alreadyOnDB)
	if res.Error != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr":  r.RemoteAddr,
			"requestURI":  r.RequestURI,
			"method":      r.Method,
			"request.url": u.URL,
			"error":       errorx.EnsureStackTrace(err),
		}).Error("Error when checking if the feed url is already in the database")

		return
	}

	// If alreadyOnDB equals to 0, then the feed is not in db.
	if alreadyOnDB == 0 {
		p.Subscribed = true

		res = m.db.Create(&p)
		if res.Error != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr":  r.RemoteAddr,
				"requestURI":  r.RequestURI,
				"method":      r.Method,
				"request.url": u.URL,
				"error":       errorx.EnsureStackTrace(err),
			}).Error("Error when trying to store the new subscribed podcast")

			return
		}

		w.WriteHeader(http.StatusCreated)
	} else { // If alreadyOnDB is different of 0 (should be 1), the feed is already on db and we should just update the 'subscribed' column.
		res = m.db.Model(&models.Podcast{}).Where("feed_link = ?", p.FeedLink).Update("subscribed", true)
		if res.Error != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr":  r.RemoteAddr,
				"requestURI":  r.RequestURI,
				"method":      r.Method,
				"request.url": u.URL,
				"error":       errorx.EnsureStackTrace(err),
			}).Error("Error when trying to update the subscription of a podcast")

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (m *Manager) UnsubscribeToPodcastHandler(w http.ResponseWriter, r *http.Request) {
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

	res := m.db.Model(&models.Podcast{}).Where("id = ?", id).Update("subscribed", false)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.EnsureStackTrace(res.Error),
			"usedID":     id,
		}).Error("Unexpected error when trying to change the subscription status of the podcast")

		return
	}

	if res.RowsAffected == 0 {
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
}

func (m *Manager) GetUserPodcastsHandler(w http.ResponseWriter, r *http.Request) {
	var p []models.Podcast

	if res := m.db.Where("subscribed = ?", true).Find(&p); res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.EnsureStackTrace(res.Error),
		}).Error("Error when trying to get subscribed podcasts from db")

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(&p)
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
}

func (m *Manager) GetPodcastHandler(w http.ResponseWriter, r *http.Request) {
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

	if res := m.db.Where("id = ?", id).First(&p); res.Error != nil {
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
}

func (m *Manager) GetEpisodesHandler(w http.ResponseWriter, r *http.Request) {
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

	if res := m.db.Where("parent_podcast_id", id).Find(&eps); res.Error != nil {
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
}

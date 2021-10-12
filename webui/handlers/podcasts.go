package handlers

import (
	"encoding/json"
	"net/http"
	"time"

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

		var id uint
		res = m.db.Model(&models.Podcast{}).Where("feed_link = ?", p.FeedLink).Select("id").Find(&id)
		if res.Error != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr":  r.RemoteAddr,
				"requestURI":  r.RequestURI,
				"method":      r.Method,
				"request.url": u.URL,
				"error":       errorx.EnsureStackTrace(err),
			}).Error("Error when trying to get the ID of the updated podcast")

			return
		}

		p.ID = id

		w.WriteHeader(http.StatusNoContent)
	}

	select {
	case m.updateChannel <- p:
	default:
		{ // Avoid blocking if, for some reason, the channel is busy
			log.WithFields(log.Fields{
				"remoteAddr":  r.RemoteAddr,
				"requestURI":  r.RequestURI,
				"method":      r.Method,
				"request.url": u.URL,
				"podcastID":   p.ID,
				"podcastFeed": p.FeedLink,
			}).Warning("The channel used to update the recently subscribed podcasts is busy; skipping feed...")
		}
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
		err := errorx.IllegalArgument.New("the value '%s' is over the limit of int values or can't be parsed", idStr)

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
		err := errorx.IllegalArgument.New("the value '%s' is over the limit of int values or can't be parsed", idStr)

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

func (m *Manager) GetEpisodesHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id := safe.SafeParseInt(idStr)
	if id == safe.DefaultAllocate {
		err := errorx.IllegalArgument.New("the value '%s' is over the limit of int values or can't be parsed", idStr)

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

	res := m.db.Where("parent_podcast_id", id).Find(&eps)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.Decorate(res.Error, "unexpected error when trying to fetch episodes from db"),
			"givenID":    id,
		}).Error("Error when trying to get the requested episodes")

		return
	}

	if len(eps) == 0 {
		errmsg := "the podcast with the given ID does not exist or has no episodes fetched"

		http.Error(w, errmsg, http.StatusNotFound)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.Decorate(res.Error, errmsg),
			"givenID":    id,
		}).Error("Error when trying to get the requested episodes")

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(&eps)
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

func (m *Manager) EpisodeProgressHandler(w http.ResponseWriter, r *http.Request) {
	podcastIDStr := mux.Vars(r)["pID"]
	epIDStr := mux.Vars(r)["epID"]

	podcastID := safe.SafeParseInt(podcastIDStr)
	if podcastID == safe.DefaultAllocate {
		err := errorx.IllegalArgument.New("the value '%s' is over the limit of int values or can't be parsed", podcastIDStr)

		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      err.Error(),
			"givenID":    podcastIDStr,
		}).Error("The given podcastID cannot be parsed")

		return
	}

	episodeID := safe.SafeParseInt(epIDStr)
	if episodeID == safe.DefaultAllocate {
		err := errorx.IllegalArgument.New("the value '%s' is over the limit of int values or can't be parsed", epIDStr)

		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      err.Error(),
			"givenID":    epIDStr,
		}).Error("The given episodeID cannot be parsed")

		return
	}

	switch r.Method {
	case http.MethodGet:
		{
			// Returns the progress of the episode
			var ep models.Episode

			res := m.db.Model(&models.Episode{}).Where("id = ? AND parent_podcast_id = ?", episodeID, podcastID).Select("current_progress").Find(&ep)
			if res.Error != nil {
				http.Error(w, res.Error.Error(), http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"requestURI": r.RequestURI,
					"method":     r.Method,
					"error":      res.Error.Error(),
					"podcastID":  podcastID,
					"episodeID":  episodeID,
				}).Error("Error when trying to get the progress of the requested episode")

				return
			}

			if res.RowsAffected == 0 {
				e := "the requested episode does not exist"

				http.Error(w, e, http.StatusBadRequest)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"requestURI": r.RequestURI,
					"method":     r.Method,
					"podcastID":  podcastID,
					"episodeID":  episodeID,
				}).Error(e)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			err := json.NewEncoder(w).Encode(map[string]time.Duration{"progress": ep.CurrentProgress})
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

	case http.MethodPut:
		{
			// Update the progress of the episode
			var requestBody struct {
				Progress time.Duration `json:"progress"`
			}

			err := json.NewDecoder(r.Body).Decode(&requestBody)
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

			res := m.db.Model(&models.Episode{}).Where("id = ? AND parent_podcast_id = ?", episodeID, podcastID).UpdateColumn("current_progress", requestBody.Progress)
			if res.Error != nil {
				http.Error(w, res.Error.Error(), http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"requestURI": r.RequestURI,
					"method":     r.Method,
					"error":      res.Error.Error(),
					"podcastID":  podcastID,
					"episodeID":  episodeID,
				}).Error("Error when trying to update the progress of an episode")

				return
			}

			if res.RowsAffected == 0 {
				e := "the episode to update does not exist"

				http.Error(w, e, http.StatusBadRequest)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"requestURI": r.RequestURI,
					"method":     r.Method,
					"podcastID":  podcastID,
					"episodeID":  episodeID,
				}).Error(e)

				return
			}

			w.WriteHeader(http.StatusCreated)
		}
	}
}

func (m *Manager) LatestEpisodesHandler(w http.ResponseWriter, r *http.Request) {
	const dateLayout = "2006-01-02"
	var from time.Time
	var to time.Time

	// Query parameter "from" processing
	keys, ok := r.URL.Query()["from"]
	if !ok || len(keys[0]) < 1 {
		err := "query parameter 'from' is missing"

		http.Error(w, err, http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      err,
		}).Error("Request rejected due to absence of the query parameter 'from'")

		return
	}

	from, err := time.Parse(dateLayout, keys[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr":   r.RemoteAddr,
			"requestURI":   r.RequestURI,
			"method":       r.Method,
			"error":        err.Error(),
			"queryContent": keys[0],
		}).Error("The query parameter 'from' can't be parsed")

		return
	}

	// Query parameter "to" processing
	keys, ok = r.URL.Query()["to"]
	if !ok || len(keys[0]) < 1 {
		err := "query parameter 'to' is missing"

		http.Error(w, err, http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      err,
		}).Error("Request rejected due to absence of the query parameter 'to'")

		return
	}

	to, err = time.Parse(dateLayout, keys[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr":   r.RemoteAddr,
			"requestURI":   r.RequestURI,
			"method":       r.Method,
			"error":        err.Error(),
			"queryContent": keys[0],
		}).Error("The query parameter 'to' can't be parsed")

		return
	}

	var eps []models.Episode

	res := m.db.Where("published BETWEEN ? AND ?", from, to).Order("published DESC").Find(&eps)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.Decorate(res.Error, "unexpected error when trying to fetch episodes from db"),
			"fromDate":   from.String(),
			"toDate":     to.String(),
		}).Error("Error when trying to get the latest episodes")

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(&eps)
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

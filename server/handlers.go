package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"lincast/database"
	"lincast/podcasts"
	"lincast/psync"

	"github.com/gorilla/mux"
	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
)

var _podcastsDB *database.Database
var _pSynchronizer *psync.Synchronizer

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
		if errorx.IsOfType(err, errorx.RejectedOperation) {
			w.WriteHeader(http.StatusConflict)

			log.WithFields(log.Fields{
				"remoteAddr":  r.RemoteAddr,
				"requestURI":  r.RequestURI,
				"method":      r.Method,
				"request.url": u.URL,
			}).Warn("The user has tried to subscribe to an already subscribed podcast")

			return
		}

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

	w.WriteHeader(http.StatusCreated)
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

	w.WriteHeader(http.StatusNoContent)
}

func getUserPodcastsHandler(w http.ResponseWriter, r *http.Request) {
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
			w.WriteHeader(http.StatusBadRequest)

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
			w.WriteHeader(http.StatusBadRequest)

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

	var subscribedPodcasts []podcasts.Podcast
	if subscribed || (!subscribed && !unsubscribed) {
		sp, err := _podcastsDB.GetPodcastsBySubscribedStatus(true)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"requestURI": r.RequestURI,
				"method":     r.Method,
				"error":      errorx.EnsureStackTrace(err),
			}).Error("Error when trying to get subscribed podcasts from db")

			return
		}

		subscribedPodcasts = *sp
	}

	var unsubscribedPodcasts []podcasts.Podcast
	if unsubscribed || (!subscribed && !unsubscribed) {
		up, err := _podcastsDB.GetPodcastsBySubscribedStatus(false)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"requestURI": r.RequestURI,
				"method":     r.Method,
				"error":      errorx.EnsureStackTrace(err),
			}).Error("Error when trying to get unsubscribed podcasts from db")

			return
		}

		unsubscribedPodcasts = *up
	}

	p := map[string][]podcasts.Podcast{
		"subscribed":   subscribedPodcasts,
		"unsubscribed": unsubscribedPodcasts,
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.EnsureStackTrace(err),
		}).Error("Error when trying to encode the response to the request")

		return
	}

	w.WriteHeader(http.StatusOK)
}

func getPodcastHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.EnsureStackTrace(err),
			"givenID":    idStr,
		}).Error("The given ID cannot be parsed")

		return
	}

	p, err := _podcastsDB.GetPodcastByID(int(id))
	if err != nil {
		if errorx.IsOfType(err, errorx.IllegalArgument) {
			w.WriteHeader(http.StatusNotFound)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"requestURI": r.RequestURI,
				"method":     r.Method,
				"error":      errorx.Decorate(err, "the podcast with the given ID does not exist"),
				"givenID":    id,
			}).Error("Error when trying to get the requested podcast")
		} else {
			w.WriteHeader(http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"requestURI": r.RequestURI,
				"method":     r.Method,
				"error":      errorx.EnsureStackTrace(err),
				"givenID":    id,
			}).Error("Unexpected error when trying to get the requested podcast")
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.EnsureStackTrace(err),
		}).Error("Error when trying to encode the response to the request")

		return
	}

	w.WriteHeader(http.StatusOK)
}

func getEpisodesHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.EnsureStackTrace(err),
			"givenID":    idStr,
		}).Error("The given ID cannot be parsed")

		return
	}

	eps, err := _podcastsDB.GetEpisodesByPodcast(int(id))
	if err != nil {
		if errorx.IsOfType(err, errorx.IllegalArgument) {
			w.WriteHeader(http.StatusNotFound)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"requestURI": r.RequestURI,
				"method":     r.Method,
				"error":      errorx.Decorate(err, "the podcast with the given ID does not exist"),
				"givenID":    id,
			}).Error("Error when trying to get the requested episodes")
		} else {
			w.WriteHeader(http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"requestURI": r.RequestURI,
				"method":     r.Method,
				"error":      errorx.EnsureStackTrace(err),
				"givenID":    id,
			}).Error("Unexpected error when trying to get the requested episodes")
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(eps)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.EnsureStackTrace(err),
		}).Error("Error when trying to encode the response to the request")

		return
	}

	w.WriteHeader(http.StatusOK)
}

func playerProgressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")

		p := _pSynchronizer.GetProgress()

		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"requestURI": r.RequestURI,
				"method":     r.Method,
				"error":      errorx.EnsureStackTrace(err),
			}).Error("Error when trying to decode the request's body")
		}

		return
	}

	var p psync.CurrentProgress
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.EnsureStackTrace(err),
		}).Error("Error when trying to decode the request's body")

		return
	}

	err = _pSynchronizer.UpdateProgress(p.Progress, p.EpisodeGUID, p.PodcastID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.EnsureStackTrace(err),
		}).Error("Error when trying to update the progress of the player")

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func queueHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		{
			// REVIEW maybe can be better the usage of io.LimitReader().
			var q []psync.QueueEpisode

			err := json.NewDecoder(r.Body).Decode(&q)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"requestURI": r.RequestURI,
					"method":     r.Method,
					"error":      errorx.EnsureStackTrace(err),
				}).Error("Error when trying to decode the request's body")

				return
			}

			// Check if the position is repeated before set the new queue.
			var positions []int
			for _, ep := range q {
				for _, p := range positions {
					if ep.Position == p {
						w.WriteHeader(http.StatusBadRequest)

						log.WithFields(log.Fields{
							"remoteAddr": r.RemoteAddr,
							"requestURI": r.RequestURI,
							"method":     r.Method,
						}).Error("The user tried to use a repeated position. Request rejected")

						return
					}
				}

				positions = append(positions, ep.Position)
			}

			err = _pSynchronizer.SetQueue(&q)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"requestURI": r.RequestURI,
					"method":     r.Method,
					"error":      errorx.EnsureStackTrace(err),
				}).Error("Error when trying to set the queue")

				return
			}

			w.WriteHeader(http.StatusCreated)
		}

	case http.MethodDelete:
		{
			err := _pSynchronizer.CleanQueue()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"requestURI": r.RequestURI,
					"method":     r.Method,
					"error":      errorx.EnsureStackTrace(err),
				}).Error("Error when trying to clean the queue")

				return
			}

			w.WriteHeader(http.StatusNoContent)
		}

	default:
		{
			q := _pSynchronizer.GetQueue()

			w.Header().Set("Content-Type", "application/json")

			err := json.NewEncoder(w).Encode(q.Content)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"requestURI": r.RequestURI,
					"method":     r.Method,
					"error":      errorx.EnsureStackTrace(err),
				}).Error("Error when trying to encode the response")

				return
			}

			w.WriteHeader(http.StatusOK)
		}
	}
}

func addToQueueHandler(w http.ResponseWriter, r *http.Request) {
}

func delFromQueueHandler(w http.ResponseWriter, r *http.Request) {
}

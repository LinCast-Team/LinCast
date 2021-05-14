package handlers

import (
	"encoding/json"
	"net/http"

	"lincast/models"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
)

var PlayerProgressHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var p models.CurrentProgress

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")

		if res := _db.First(&p); res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"requestURI": r.RequestURI,
				"method":     r.Method,
				"error":      errorx.EnsureStackTrace(res.Error),
			}).Error("Error when trying to get the current progress of the player")

			return
		}

		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"requestURI": r.RequestURI,
				"method":     r.Method,
				"error":      errorx.EnsureStackTrace(err),
			}).Error("Error when trying to decode the request's body")
		}

		return
	}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.EnsureStackTrace(err),
		}).Error("Error when trying to decode the request's body")

		return
	}

	// Update the first (and only) row of the table that stores the current progress
	// of the player.
	if res := _db.Model(&models.CurrentProgress{}).Where("id = ?", 1).Updates(&p); res.Error != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"requestURI": r.RequestURI,
			"method":     r.Method,
			"error":      errorx.EnsureStackTrace(err),
		}).Error("Error when trying to update the progress of the player")

		return
	}

	w.WriteHeader(http.StatusCreated)
})

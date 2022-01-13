package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"lincast/models"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (m *Manager) PlayerPlaybackInfoHandler(w http.ResponseWriter, r *http.Request) {
	var p models.PlaybackInfo

	switch r.Method {
	case http.MethodGet:
		{
			w.Header().Set("Content-Type", "application/json")

			if res := m.db.First(&p); res.Error != nil {
				// If the error is of type gorm.ErrRecordNotFound, it means that there is no episode
				// being played, so we should let it know it to the user that there is no content
				if errors.Is(res.Error, gorm.ErrRecordNotFound) {
					http.Error(w, "there is no episode being played", http.StatusNotFound)

					log.WithFields(log.Fields{
						"remoteAddr": r.RemoteAddr,
					}).Warning("The client requested the player's playback info, but no episode is being played")

					return
				}

				http.Error(w, res.Error.Error(), http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"error":      errorx.EnsureStackTrace(res.Error),
				}).Error("Error when trying to get the player's playback info")

				return
			}

			w.WriteHeader(http.StatusOK)

			err := json.NewEncoder(w).Encode(&p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"error":      errorx.EnsureStackTrace(err),
				}).Error("Error when trying to decode the request's body")
			}

			return
		}
	case http.MethodPut:
		{
			err := json.NewDecoder(r.Body).Decode(&p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"error":      errorx.EnsureStackTrace(err),
				}).Error("Error when trying to decode the request's body")

				return
			}

			// Try to update the first (and only) row of the table that stores the playback info
			// of the player.
			res := m.db.Model(&models.PlaybackInfo{}).Where("id = ?", 1).Updates(&p)
			if res.Error != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"error":      errorx.EnsureStackTrace(err),
				}).Error("Error when trying to update the player's playback info")

				return
			}

			// If no rows have been afected (which means that there are no records in the table), we should
			// create the first record.
			if res.RowsAffected == 0 {
				res = m.db.Create(&p)
				if res.Error != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)

					log.WithFields(log.Fields{
						"remoteAddr": r.RemoteAddr,
						"error":      errorx.EnsureStackTrace(err),
					}).Error("Error when trying to add the first entry of the table that stores the player's playback info")

					return
				}
			}

			w.WriteHeader(http.StatusCreated)
		}
	}
}

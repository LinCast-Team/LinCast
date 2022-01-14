package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"lincast/models"
	"lincast/utils/safe"

	"github.com/gorilla/mux"
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

func (m *Manager) SetEpisodeStatusHandler(w http.ResponseWriter, r *http.Request)  {
	podcastIDStr := mux.Vars(r)["pID"]
	epIDStr := mux.Vars(r)["epID"]

	podcastID := safe.SafeParseInt(podcastIDStr)
	if podcastID == safe.DefaultAllocate {
		err := errorx.IllegalArgument.New("value is over the limit of int values or can't be parsed")

		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"error":      err.Error(),
		}).Error("The given podcastID cannot be parsed")

		return
	}

	episodeID := safe.SafeParseInt(epIDStr)
	if episodeID == safe.DefaultAllocate {
		err := errorx.IllegalArgument.New("value is over the limit of int values or can't be parsed")

		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"error":      err.Error(),
		}).Error("The given episodeID cannot be parsed")

		return
	}

	reqBody := struct {
		Played bool `json:"played"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"error":      err.Error(),
		}).Error("Error when trying to decode the body of the request")

		return
	}

	res := m.db.Model(&models.Episode{}).Where("id = ? AND parent_podcast_id = ?", episodeID, podcastID).UpdateColumn("played", reqBody.Played)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"error":      res.Error.Error(),
			"podcastID":  podcastID,
			"episodeID":  episodeID,
		}).Error("Error when trying to update the progress of an episode")

		return
	}

	if res.RowsAffected == 0 {
		err := "episode not found"
		
		http.Error(w, err, http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"error":      err,
			"podcastID":  podcastID,
			"episodeID":  episodeID,
		}).Error("Error when trying to update the progress of an episode")

		return
	}

	w.Header().Set("Location", fmt.Sprintf("/api/v0/podcasts/%d/episodes/%d", podcastID, episodeID))
	w.WriteHeader(http.StatusCreated)
}

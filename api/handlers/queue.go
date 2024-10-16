package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"lincast/models"
	"lincast/utils/safe"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (m *Manager) QueueHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		{
			// REVIEW maybe can be better the usage of io.LimitReader().
			var q []models.QueueEpisode

			err := json.NewDecoder(r.Body).Decode(&q)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"error":      errorx.EnsureStackTrace(err),
				}).Error("Error when trying to decode the request's body")

				return
			}

			// Check if the position is repeated before set the new queue.
			var positions []uint
			for _, ep := range q {
				for _, p := range positions {
					if ep.Position == p {
						http.Error(w, "two or more episodes have the same position", http.StatusBadRequest)

						log.WithFields(log.Fields{
							"remoteAddr": r.RemoteAddr,
						}).Error("The user tried to use a repeated position. Request rejected")

						return
					}
				}

				positions = append(positions, ep.Position)
			}

			// First we delete all the rows of the table.
			if res := m.db.Unscoped().Where("1 = 1").Delete(&models.QueueEpisode{}); res.Error != nil {
				http.Error(w, res.Error.Error(), http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"error":      errorx.EnsureStackTrace(err),
				}).Error("Error when trying to clean the queue (before set the new content)")

				return
			}

			// And later we introduce the new elements of the queue.
			if res := m.db.Create(&q); res.Error != nil {
				http.Error(w, res.Error.Error(), http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"error":      errorx.EnsureStackTrace(err),
				}).Error("Error when trying to set the new queue")

				return
			}

			w.Header().Set("Location", "/api/v0/player/queue")
			w.WriteHeader(http.StatusCreated)
		}

	case http.MethodDelete:
		{
			// Delete all the rows of the table.
			if res := m.db.Unscoped().Where("1 = 1").Delete(&models.QueueEpisode{}); res.Error != nil {
				http.Error(w, res.Error.Error(), http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"error":      errorx.EnsureStackTrace(res.Error),
				}).Error("Error when trying to clean the queue")

				return
			}

			w.WriteHeader(http.StatusNoContent)
		}

	default:
		{
			var q []models.QueueEpisode

			res := m.db.Find(&q)
			if res.Error != nil {
				http.Error(w, res.Error.Error(), http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"error":      errorx.EnsureStackTrace(res.Error),
				}).Error("Error when trying to fetch the queue from the database")

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			err := json.NewEncoder(w).Encode(&q)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"error":      errorx.EnsureStackTrace(err),
				}).Error("Error when trying to encode the response")

				return
			}
		}
	}
}

func (m *Manager) AddToQueueHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["append"]
	if !ok || len(keys[0]) < 1 {
		err := errorx.IllegalFormat.New("param 'append' is missing")

		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"error":      err.Error(),
		}).Error("Request rejected due to absence of parameter 'append'")

		return
	}

	appendStr := keys[0]

	append, err := strconv.ParseBool(appendStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"error":      errorx.EnsureStackTrace(err),
		}).Error("The variable 'append' is not present in the request or the value cannot be parsed")

		return
	}

	var ep models.QueueEpisode

	err = json.NewDecoder(r.Body).Decode(&ep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"error":      errorx.EnsureStackTrace(err),
		}).Error("Error when trying to decode the request's body")

		return
	}

	if append {
		var refPosition uint
		// To append the new episode to the queue we need to know which is the bigger position
		// stored, so the position of the new episode will be that + 1.
		if res := m.db.Model(&models.QueueEpisode{}).Select("position").Limit(1).Order("position desc").First(&refPosition); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				// If the error is because there are no episodes stored (which means that the queue is
				// empty), we know that the new episode should be stored with the position 1.
				ep.Position = 1

				if res := m.db.Create(&ep); res.Error != nil {
					http.Error(w, res.Error.Error(), http.StatusInternalServerError)

					log.WithFields(log.Fields{
						"remoteAddr": r.RemoteAddr,
						"error":      errorx.EnsureStackTrace(res.Error),
					}).Error("Error when trying to add an episode to the queue")

					return
				}

				// Since we need to avoid the rest of the body of the underlying if (the one that checks
				// the variable `append`), we'll use a `goto` statement.
				goto response
			} else {
				// If there is an unexpected error then we need to abort the operation, log and notify the user.
				http.Error(w, res.Error.Error(), http.StatusInternalServerError)

				log.WithFields(log.Fields{
					"remoteAddr": r.RemoteAddr,
					"error":      errorx.EnsureStackTrace(res.Error),
				}).Error("Error when trying to find the last position in the queue")

				return
			}
		}

		// Last position + 1
		ep.Position = refPosition + 1

		if res := m.db.Create(&ep); res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"error":      errorx.EnsureStackTrace(res.Error),
			}).Error("Error when trying to add an episode to the queue")

			return
		}
	} else {
		// To add the episode with the position 1 we'll need to update the position of the rest
		// of episodes, adding 1 to each one.
		if res := m.db.Model(&models.QueueEpisode{}).Where("1 = 1").Update("position", gorm.Expr("position + 1")); res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"error":      errorx.EnsureStackTrace(res.Error),
			}).Error("Error when trying to update the position of the episodes that are already in the queue")

			return
		}

		// Now that all the positions have been updated, we need to insert the new episode with the position 1.
		ep.Position = 1
		if res := m.db.Create(&ep); res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusInternalServerError)

			log.WithFields(log.Fields{
				"remoteAddr": r.RemoteAddr,
				"error":      errorx.EnsureStackTrace(res.Error),
			}).Error("Error when trying to add an episode to the queue")

			return
		}
	}

response:

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/api/v0/player/queue")
	w.WriteHeader(http.StatusCreated)

	response := map[string]uint{
		"episodeID": ep.ID,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"error":      errorx.EnsureStackTrace(err),
		}).Error("Error when trying to encode the response")

		return
	}
}

func (m *Manager) DelFromQueueHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		err := errorx.IllegalFormat.New("param 'id' is missing")

		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"error":      err.Error(),
		}).Error("Request rejected due to absence of parameter 'id'")

		return
	}

	idStr := keys[0]

	id := safe.SafeParseInt(idStr)
	if id == safe.DefaultAllocate {
		err := errorx.IllegalArgument.New("the value '%s' is over the limit of int values or can't be parsed", safe.Sanitize(idStr))

		http.Error(w, err.Error(), http.StatusBadRequest)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"error":      err.Error(),
		}).Error("The variable 'id' is not present in the request or the value cannot be parsed")

		return
	}

	res := m.db.Unscoped().Delete(&models.QueueEpisode{}, id)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"error":      errorx.Decorate(res.Error, res.Error.Error()),
			"usedID":     id,
		}).Error("Error when trying to remove an episode from the queue")

		return
	}

	if res.RowsAffected == 0 {
		errmsg := "the episode of the queue with the given ID does not exist"

		http.Error(w, errmsg, http.StatusNotFound)

		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"error":      errorx.Decorate(res.Error, errmsg),
			"usedID":     id,
		}).Warning("Usage of the wrong ID when trying to remove an episode from the queue")

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

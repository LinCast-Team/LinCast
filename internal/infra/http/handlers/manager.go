package handlers

import (
	"lincast/internal/domain/repositories"
	"lincast/models"

	"github.com/go-chi/jwtauth/v5"
)

type Manager struct {
	updateChannel     chan *models.Podcast
	tokenAuth         *jwtauth.JWTAuth
	userRepository    *repositories.UserRepository
	podcastRepository *repositories.PodcastRepository
	playerRepository  *repositories.PlayerRepository
	queueRepository   *repositories.QueueRepository
}

// NewManager returns a new Manager. The `Manager` is who provides the access to the handlers. The unique function of
// this is to provide the access to the database in an ordered way to all the handlers, without the usage of global
// variables.
func NewManager(
	manualUpdate chan *models.Podcast,
	tokenAuth *jwtauth.JWTAuth,
	userRepository *repositories.UserRepository,
	podcastRepository *repositories.PodcastRepository,
	playerRepository *repositories.PlayerRepository,
	queueRepository *repositories.QueueRepository,
) *Manager {
	m := Manager{
		updateChannel:     manualUpdate,
		tokenAuth:         tokenAuth,
		userRepository:    userRepository,
		podcastRepository: podcastRepository,
		playerRepository:  playerRepository,
		queueRepository:   queueRepository,
	}

	return &m
}

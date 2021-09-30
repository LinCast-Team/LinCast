package handlers

import (
	"lincast/models"

	"gorm.io/gorm"
)

type Manager struct {
	updateChannel chan *models.Podcast
	db *gorm.DB
}

// NewManager returns a new Manager. The `Manager` is who provides the access to the handlers. The unique function of
// this is to provide the access to the database in an ordered way to all the handlers, without the usage of global
// variables.
func NewManager(db *gorm.DB, manualUpdate chan *models.Podcast) *Manager {
	m := Manager{
		updateChannel: manualUpdate,
		db: db,
	}

	return &m
}

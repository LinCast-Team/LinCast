package handlers

import "gorm.io/gorm"

type Manager struct {
	db *gorm.DB
}

// NewManager returns a new Manager. The `Manager` is who provides the access to the handlers. The unique function of
// this is to provide the access to the database in an ordered way to all the handlers, without the usage of global
// variables.
func NewManager(db *gorm.DB) *Manager {
	m := Manager{
		db: db,
	}

	return &m
}

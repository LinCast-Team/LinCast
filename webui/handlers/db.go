package handlers

import "gorm.io/gorm"

// SetGlobalDB stores the instance of the database in a variable used by the handlers
// to process the data.
func SetGlobalDB(db *gorm.DB) {
	_db = db
}

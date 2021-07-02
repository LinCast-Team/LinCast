package database

import (
	"errors"
	"lincast/models"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNew(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()

	db, err := New(tempDir, "test.db")

	assert.NoError(err, "The database should be initialized without errors")
	assert.NotNil(db, "A valid instance of the database should be returned")

	var p models.CurrentProgress
	tx := db.First(&p)

	if !assert.NoError(tx.Error, "There shouldn't be errors when trying to get the first record of the table that contains the current progress of the player") {
		assert.False(errors.Is(tx.Error, gorm.ErrRecordNotFound), "The table that contains the current progress of the player should contain one row by default")
	}
	assert.NotNil(p, "The table that contains the current progress of the player should contain one row by default")
}

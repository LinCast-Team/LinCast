package handlers

import (
	"testing"

	"lincast/database"

	assert2 "github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()
	db, err := database.New(tempDir, "test.db")
	if err != nil {
		assert.FailNow(err.Error())
	}

	mng := NewManager(db)

	assert.NotNil(mng, "A valid instance of Manager should be returned")
}

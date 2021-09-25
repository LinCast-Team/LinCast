package database

import (
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert2.New(t)
	tempDir := t.TempDir()

	db, err := New(tempDir, "test.db")

	assert.NoError(err, "The database should be initialized without errors")
	assert.NotNil(db, "A valid instance of the database should be returned")
}

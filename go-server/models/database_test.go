package models

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewModels(t *testing.T) {
	// Open database mock
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	testModels := NewModels(db)

	assert.NotNil(t, testModels.DB)
}

package handlers

import (
	"net/http"
	"testing"

	"github.com/acornak/healthcare-poc/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewHandler(t *testing.T) {
	logger := zap.NewExample()
	handler := NewHandler(logger, models.Models{}, http.Get)

	assert.NotNil(t, handler)
}

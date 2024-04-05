package handlers

import (
	"github.com/acornak/healthcare-poc/models"
	"go.uber.org/zap"
)

type Handler struct {
	Logger *zap.Logger
	Models models.Models
}

type ErrorResponse struct {
	Error string `json:"error"`
}

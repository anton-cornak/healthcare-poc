package handlers

import (
	"github.com/acornak/poc-gpt/models"
	"go.uber.org/zap"
)

type Handler struct {
	Logger *zap.Logger
	Models models.Models
}

type SuccessResponse struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

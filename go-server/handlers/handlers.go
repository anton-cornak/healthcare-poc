package handlers

import (
	"net/http"

	"github.com/acornak/healthcare-poc/models"
	"go.uber.org/zap"
)

type Handler struct {
	Logger *zap.Logger
	Models models.Models
	Get    func(url string) (resp *http.Response, err error)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

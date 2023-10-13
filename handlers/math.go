package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AddPayload struct {
	Numbers []float64 `json:"numbers"`
}

type Handler struct {
	Logger *zap.Logger
}

func (h *Handler) Add(c *gin.Context) {
	var payload AddPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	result := 0.0
	for _, num := range payload.Numbers {
		result += num
	}

	h.Logger.Info("Endpoint /compute called",
		zap.Any("payload", payload),
		zap.Float64("result", result),
	)

	c.JSON(http.StatusOK, gin.H{"result": result})
}

type SubstractPayload struct {
	Number   float64   `json:"number"`
	Subtract []float64 `json:"subtract"`
}

func (h *Handler) Subtract(c *gin.Context) {
	var payload SubstractPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	result := payload.Number
	for _, num := range payload.Subtract {
		result -= num
	}

	h.Logger.Info("Endpoint /compute called",
		zap.Any("payload", payload),
		zap.Float64("result", result),
	)

	c.JSON(http.StatusOK, gin.H{"result": result})
}

type ComputePayload struct {
	Add      []float64 `json:"add"`
	Subtract []float64 `json:"subtract"`
}

func (h *Handler) Compute(c *gin.Context) {
	var payload ComputePayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	result := float64(0)
	for _, num := range payload.Add {
		result += num
	}

	for _, num := range payload.Subtract {
		result -= num
	}

	h.Logger.Info("Endpoint /compute called",
		zap.Any("payload", payload),
		zap.Float64("result", result),
	)

	c.JSON(http.StatusOK, gin.H{"result": result})
}

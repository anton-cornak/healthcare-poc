package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SuccessResponse struct {
	Result float64 `json:"result"`
}

type AddPayload struct {
	Numbers []float64 `json:"numbers"`
}

// @Summary		Add numbers
// @Description	Add all numbers provided in the payload
// @ID			add-operation
// @Accept		json
// @Produce		json
// @Param		payload	body		AddPayload	true	"Numbers to add"
// @Success		200		{object}	SuccessResponse
// @Failure		400		{object}	ErrorResponse
// @Router		/math/add [post]
func (h *Handler) Add(c *gin.Context) {
	var payload AddPayload
	var errResp ErrorResponse
	var successResp SuccessResponse

	if err := c.ShouldBindJSON(&payload); err != nil {
		errResp.Error = "Invalid JSON payload"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	if len(payload.Numbers) == 0 {
		errResp.Error = "Invalid payload: missing numbers field"
		c.JSON(http.StatusBadRequest, errResp)
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

	successResp.Result = result
	c.JSON(http.StatusOK, successResp)
}

type SubtractPayload struct {
	Number   float64   `json:"number"`
	Subtract []float64 `json:"subtract"`
}

// Subtract performs subtraction operation on a given number and a list of numbers.
//
//	@Summary		Subtract numbers
//	@Description	Subtract all numbers in the 'subtract' list from the 'number'.
//	@Tags			Math Operations
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		SubtractPayload	true	"Numbers to substract from the 'number'"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/math/subtract [post]
func (h *Handler) Subtract(c *gin.Context) {
	var payload SubtractPayload
	var errResp ErrorResponse
	var successResp SuccessResponse

	if err := c.ShouldBindJSON(&payload); err != nil {
		errResp.Error = "Invalid JSON payload"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	result := payload.Number
	for _, num := range payload.Subtract {
		result -= num
	}

	if len(payload.Subtract) == 0 {
		errResp.Error = "Invalid payload: missing subtract field"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	h.Logger.Info("Endpoint /compute called",
		zap.Any("payload", payload),
		zap.Float64("result", result),
	)

	successResp.Result = result
	c.JSON(http.StatusOK, successResp)
}

type ComputePayload struct {
	Add      []float64 `json:"add"`
	Subtract []float64 `json:"subtract"`
}

// Compute performs addition and subtraction on lists of numbers.
//
//	@Summary		Compute result
//	@Description	Adds all numbers in the 'add' list and subtracts all numbers in the 'subtract' list.
//	@Tags			Math Operations
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		ComputePayload	true	"Compute payload"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/math/compute [post]
func (h *Handler) Compute(c *gin.Context) {
	var payload ComputePayload
	var errResp ErrorResponse
	var successResp SuccessResponse

	if err := c.ShouldBindJSON(&payload); err != nil {
		errResp.Error = "Invalid JSON payload"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	if len(payload.Add) == 0 && len(payload.Subtract) == 0 {
		errResp.Error = "Invalid payload: missing add and subtract fields"
		c.JSON(http.StatusBadRequest, errResp)
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

	successResp.Result = result
	c.JSON(http.StatusOK, successResp)
}

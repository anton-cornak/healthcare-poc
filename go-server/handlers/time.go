package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GetCurrentTimePayload struct {
	Timezone string `json:"timezone"`
}

type GetCurrentTimeResponse struct {
	Time string `json:"time"`
}

// @Summary		Get current time
// @Description	Get the current time in a specific timezone
// @ID			current-time
// @Accept		json
// @Produce		json
// @Success		200		{object}	GetCurrentTimeResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/time/current [post]
func (h *Handler) GetCurrentTime(c *gin.Context) {
	var payload GetCurrentTimePayload
	var errResp ErrorResponse

	if err := c.ShouldBindJSON(&payload); err != nil {
		errResp.Error = "Invalid JSON payload"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	if payload.Timezone == "" {
		errResp.Error = "Invalid payload: missing timezone field"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	loc, err := time.LoadLocation(payload.Timezone)
	if err != nil {
		errResp.Error = "Invalid timezone"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	currentTime := time.Now().In(loc).Format("2006-01-02 15:04:05")

	c.JSON(http.StatusOK, GetCurrentTimeResponse{Time: currentTime})
}

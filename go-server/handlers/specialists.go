package handlers

import (
	"net/http"

	"github.com/acornak/healthcare-poc/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FindSpecialistPayload struct {
	SpecialtyId  int    `json:"specialty_id"`
	Radius       int    `json:"radius"`
	UserLocation string `json:"user_location"`
}

type FindSpecialistResponse struct {
	Specialists []*types.Specialist `json:"specialists"`
}

// @Summary		Find specialist
// @Description	Find a specialist based on the user's location, specialty, and radius
// @ID			find-specialist
// @Accept		json
// @Produce		json
// @Param		payload	body		FindSpecialistPayload	true	"Specialty, radius, and user location"
// @Success		200		{object}	SuccessResponse
// @Failure		400		{object}	ErrorResponse
// @Router		/find-specialist [post]
func (h *Handler) FindSpecialist(c *gin.Context) {
	var payload FindSpecialistPayload
	var errResp ErrorResponse

	if err := c.ShouldBindJSON(&payload); err != nil {
		errResp.Error = "Invalid JSON payload"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	h.Logger.Info("Finding specialist", zap.Any("payload", payload))

	specialists, err := h.Models.DB.GetSpecialistBySpecialtyAndLocation(payload.SpecialtyId, payload.Radius, payload.UserLocation)
	if err != nil {
		errResp.Error = "Specialists not found"
		c.JSON(http.StatusNotFound, errResp)
		return
	}

	c.JSON(http.StatusOK, FindSpecialistResponse{Specialists: specialists})

}

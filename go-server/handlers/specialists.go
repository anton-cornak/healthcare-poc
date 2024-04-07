package handlers

import (
	"net/http"
	"strings"

	"github.com/acornak/healthcare-poc/types"
	"github.com/gin-gonic/gin"
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
// @Failure		500		{object}	ErrorResponse
// @Router		/specialist/find [post]
func (h *Handler) FindSpecialist(c *gin.Context) {
	var payload FindSpecialistPayload
	var errResp ErrorResponse

	if err := c.ShouldBindJSON(&payload); err != nil {
		errResp.Error = "Invalid JSON payload"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	missingParams := []string{}
	if payload.SpecialtyId == 0 {
		missingParams = append(missingParams, "specialty_id")
	}
	if payload.Radius == 0 {
		missingParams = append(missingParams, "radius")
	}
	if payload.UserLocation == "" {
		missingParams = append(missingParams, "user_location")
	}

	if len(missingParams) > 0 {
		errResp.Error = "Invalid payload: missing " + strings.Join(missingParams, ", ")
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	specialists, err := h.Models.DB.GetSpecialistBySpecialtyAndLocation(payload.SpecialtyId, payload.Radius, payload.UserLocation)
	if err != nil {
		errResp.Error = err.Error()
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	c.JSON(http.StatusOK, FindSpecialistResponse{Specialists: specialists})

}

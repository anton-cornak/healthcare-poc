package handlers

import (
	"net/http"

	"github.com/acornak/healthcare-poc/types"
	"github.com/gin-gonic/gin"
)

type GetSpecialtiesResponse struct {
	Specialties []*types.Specialty `json:"specialties"`
}

// @Summary		Get specialties
// @Description	Get all specialties
// @ID			specialties
// @Accept		json
// @Produce		json
// @Success		200		{object}	GetSpecialtiesResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/specialty/all [post]
func (h *Handler) GetSpecialties(c *gin.Context) {
	var errResp ErrorResponse

	specialties, err := h.Models.DB.AllSpecialties()
	if err != nil {
		errResp.Error = err.Error()
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	c.JSON(http.StatusOK, GetSpecialtiesResponse{Specialties: specialties})
}

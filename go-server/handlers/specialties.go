package handlers

import (
	"net/http"

	"github.com/acornak/healthcare-poc/types"
	"github.com/gin-gonic/gin"
)

type GetSpecialtiesResponse struct {
	Specialties []*types.Specialty `json:"specialties"`
}

func (h *Handler) GetSpecialties(c *gin.Context) {
	var errResp ErrorResponse

	specialties, err := h.Models.DB.AllSpecialties()
	if err != nil {
		errResp.Error = "Specialty not found"
		c.JSON(http.StatusNotFound, errResp)
		return
	}

	c.JSON(http.StatusOK, GetSpecialtiesResponse{Specialties: specialties})
}

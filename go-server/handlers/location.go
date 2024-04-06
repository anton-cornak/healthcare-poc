package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type locationResult struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type GetWKTLocationPayload struct {
	UserLocation string `json:"user_location"`
}

type GetWKTLocationResponse struct {
	WKTLocation string `json:"wkt_location"`
}

// @Summary		Get WKT location
// @Description	Get the WKT location based on the user's location
// @ID			location
// @Accept		json
// @Produce		json
// @Param		payload	body		GetWKTLocationPayload	true	"User location"
// @Success		200		{object}	SuccessResponse
// @Failure		400		{object}	ErrorResponse
// @Router		/location [post]
func (h *Handler) GetWKTLocation(c *gin.Context) {
	var payload GetWKTLocationPayload
	var errResp ErrorResponse

	if err := c.ShouldBindJSON(&payload); err != nil {
		errResp.Error = "Invalid JSON payload"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	if payload.UserLocation == "" {
		errResp.Error = "Invalid payload: missing user_location field"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	h.Logger.Info("Generating WKT location", zap.Any("payload", payload))

	apiKey := os.Getenv("GEOCODE_API_KEY")
	if apiKey == "" {
		h.Logger.Error("API key is required")
		errResp.Error = "Something went wrong, please try again later"
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	geocodeURL := fmt.Sprintf("https://geocode.maps.co/search?q=%s&api_key=%s", payload.UserLocation, apiKey)
	geocodeResp, err := h.Get(geocodeURL)
	if err != nil {
		h.Logger.Error("Failed to get geocode location", zap.Error(err))
		errResp.Error = "Failed to get geocode location"
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	if geocodeResp.StatusCode != http.StatusOK {
		h.Logger.Error("Failed to get geocode location", zap.Any("status_code", geocodeResp.StatusCode))
		errResp.Error = "Failed to get geocode location"
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}
	defer geocodeResp.Body.Close()

	var geocodeData []locationResult
	if err := json.NewDecoder(geocodeResp.Body).Decode(&geocodeData); err != nil {
		h.Logger.Error("Failed to decode geocode response", zap.Error(err))
		errResp.Error = "Failed to get geocode location"
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	if len(geocodeData) == 0 {
		h.Logger.Error("No geocode data found for the location", zap.Any("location", payload.UserLocation))
		errResp.Error = "No geocode data found for the location provided"
		c.JSON(http.StatusOK, errResp)
		return
	}

	c.JSON(http.StatusOK, GetWKTLocationResponse{WKTLocation: fmt.Sprintf("POINT(%s %s)", geocodeData[0].Lon, geocodeData[0].Lat)})
}

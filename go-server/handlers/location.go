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
// @ID			get-wkt-location
// @Accept		json
// @Produce		json
// @Param		payload	body		GetWKTLocationPayload	true	"User location"
// @Success		200		{object}	SuccessResponse
// @Failure		400		{object}	ErrorResponse
// @Router		/get-wkt-location [post]
func (h *Handler) GetWKTLocation(c *gin.Context) {
	var payload GetWKTLocationPayload
	var errResp ErrorResponse

	if err := c.ShouldBindJSON(&payload); err != nil {
		errResp.Error = "Invalid JSON payload"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	h.Logger.Info("Generating WKT location", zap.Any("payload", payload))

	apiKey := os.Getenv("GEOCODE_API_KEY")
	if apiKey == "" {
		h.Logger.Error("API key is required")
		errResp.Error = "API key is required"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	geocodeURL := fmt.Sprintf("https://geocode.maps.co/search?q=%s&api_key=%s", payload.UserLocation, apiKey)
	geocodeResp, err := http.Get(geocodeURL)
	if err != nil {
		h.Logger.Error("Failed to get geocode response", zap.Error(err))
		errResp.Error = "Failed to get geocode response"
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}
	defer geocodeResp.Body.Close()

	if geocodeResp.StatusCode != http.StatusOK {
		h.Logger.Error("Failed to get geocode response", zap.Any("status_code", geocodeResp.StatusCode))
		errResp.Error = "Failed to get geocode response"
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	var geocodeData []locationResult
	if err := json.NewDecoder(geocodeResp.Body).Decode(&geocodeData); err != nil {
		h.Logger.Error("Failed to decode geocode response", zap.Error(err))
		errResp.Error = "Failed to decode geocode response"
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	if len(geocodeData) == 0 {
		h.Logger.Error("No geocode data found")
		errResp.Error = "No geocode data found"
		c.JSON(http.StatusNotFound, errResp)
		return
	}

	firstResult := geocodeData[0]
	wktLocation := fmt.Sprintf("POINT(%s %s)", firstResult.Lon, firstResult.Lat)

	c.JSON(http.StatusOK, GetWKTLocationResponse{WKTLocation: wktLocation})
}

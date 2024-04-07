package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type wktResult struct {
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
// @Router		/location/wkt [post]
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

	apiKey := os.Getenv("GEOCODE_API_KEY")
	if apiKey == "" {
		h.Logger.Error("API key is required")
		errResp.Error = "Something went wrong, please try again later"
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	geocodeUrl := os.Getenv("GEOCODE_URL")
	if geocodeUrl == "" {
		h.Logger.Error("Geocode URL is required")
		errResp.Error = "Something went wrong, please try again later"
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	geocodeResp, err := h.Get(fmt.Sprintf("%s/search?q=%s&api_key=%s", geocodeUrl, payload.UserLocation, apiKey))
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

	var geocodeData []wktResult
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

type GetAddressFromWKTPayload struct {
	WKTLocation string `json:"wkt_location"`
}

type GetAddressFromWKTResponse struct {
	Address string `json:"display_name"`
}

// @Summary		Get address from WKT
// @Description	Get the address based on the WKT location
// @ID			address
// @Accept		json
// @Produce		json
// @Param		payload	body		GetAddressFromWKTPayload	true	"WKT location"
// @Success		200		{object}	GetAddressFromWKTResponse
// @Failure		400		{object}	ErrorResponse
// @Router		/location/address [post]
func (h *Handler) GetAddressFromWKT(c *gin.Context) {
	var payload GetAddressFromWKTPayload
	var errResp ErrorResponse

	if err := c.ShouldBindJSON(&payload); err != nil {
		errResp.Error = "Invalid JSON payload"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	if payload.WKTLocation == "" {
		errResp.Error = "Invalid payload: missing wkt_location field"
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	apiKey := os.Getenv("GEOCODE_API_KEY")
	if apiKey == "" {
		h.Logger.Error("API key is required")
		errResp.Error = "Something went wrong, please try again later"
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	geocodeUrl := os.Getenv("GEOCODE_URL")
	if geocodeUrl == "" {
		h.Logger.Error("Geocode URL is required")
		errResp.Error = "Something went wrong, please try again later"
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	lon, lat, err := getLatAndLonFromWKT(payload.WKTLocation)
	if err != nil {
		h.Logger.Error("Failed to get lat lon from WKT", zap.Error(err))
		errResp.Error = "Param 'wkt_location' is not a valid WKT location"
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	geocodeResp, err := h.Get(fmt.Sprintf("%s/reverse?lat=%s&lon=%s&api_key=%s", geocodeUrl, lat, lon, apiKey))
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

	var geocodeData GetAddressFromWKTResponse
	if err := json.NewDecoder(geocodeResp.Body).Decode(&geocodeData); err != nil {
		h.Logger.Error("Failed to decode geocode response", zap.Error(err))
		errResp.Error = "Failed to get geocode location"
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	if geocodeData.Address == "" {
		h.Logger.Error("No geocode data found for the location", zap.Any("location", payload.WKTLocation))
		errResp.Error = "No geocode data found for the location provided"
		c.JSON(http.StatusOK, errResp)
		return
	}

	c.JSON(http.StatusOK, GetAddressFromWKTResponse{Address: geocodeData.Address})
}

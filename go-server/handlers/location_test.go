package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetWKTLocationHandler_InvalidJson(t *testing.T) {
	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	// Define an invalid payload
	payload := "{invalid_json}"

	// Create a test request with the JSON payload
	req, _ := http.NewRequest("POST", "/location/wkt", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/wkt", handler.GetWKTLocation)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Invalid JSON payload", response.Error)
}

func TestGetWKTLocationHandler_IncompletePayload(t *testing.T) {
	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := GetWKTLocationPayload{}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/wkt", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/wkt", handler.GetWKTLocation)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Invalid payload: missing user_location field", response.Error)
}

func TestGetWKTLocationHandler_MissingGeocodeApiKey(t *testing.T) {
	os.Unsetenv("GEOCODE_API_KEY")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := GetWKTLocationPayload{UserLocation: "New York, NY"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/wkt", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/wkt", handler.GetWKTLocation)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Something went wrong, please try again later", response.Error)
}

func TestGetWKTLocationHandler_MissingGeocodeUrl(t *testing.T) {
	os.Setenv("GEOCODE_API_KEY", "some-key")
	os.Unsetenv("GEOCODE_URL")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := GetWKTLocationPayload{UserLocation: "New York, NY"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/wkt", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/wkt", handler.GetWKTLocation)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Something went wrong, please try again later", response.Error)
}

func TestGetWKTLocationHandler_GetError(t *testing.T) {
	os.Setenv("GEOCODE_API_KEY", "some-key")
	os.Setenv("GEOCODE_URL", "http://geocode.url")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
		Get:    func(url string) (*http.Response, error) { return nil, errors.New("http get error") },
	}

	payload := GetWKTLocationPayload{UserLocation: "New York, NY"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/wkt", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/wkt", handler.GetWKTLocation)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Failed to get geocode location", response.Error)
}

func TestGetWKTLocationHandler_GetInternalServerError(t *testing.T) {
	os.Setenv("GEOCODE_API_KEY", "some-key")
	os.Setenv("GEOCODE_URL", "http://geocode.url")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusInternalServerError}, nil
		},
	}

	payload := GetWKTLocationPayload{UserLocation: "New York, NY"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/wkt", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/wkt", handler.GetWKTLocation)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Failed to get geocode location", response.Error)
}

func TestGetWKTLocationHandler_ErrorDecodingBody(t *testing.T) {
	os.Setenv("GEOCODE_API_KEY", "some-key")
	os.Setenv("GEOCODE_URL", "http://geocode.url")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("{invalid_json}")),
			}, nil
		},
	}

	payload := GetWKTLocationPayload{UserLocation: "New York, NY"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/wkt", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/wkt", handler.GetWKTLocation)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Failed to get geocode location", response.Error)
}

func TestGetWKTLocationHandler_EmptyGeocodeResponse(t *testing.T) {
	os.Setenv("GEOCODE_API_KEY", "some-key")
	os.Setenv("GEOCODE_URL", "http://geocode.url")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
			}, nil
		},
	}

	payload := GetWKTLocationPayload{UserLocation: "New York, NY"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/wkt", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/wkt", handler.GetWKTLocation)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "No geocode data found for the location provided", response.Error)
}

func TestGetWKTLocationHandler_Success(t *testing.T) {
	os.Setenv("GEOCODE_API_KEY", "some-key")
	os.Setenv("GEOCODE_URL", "http://geocode.url")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	geocodeResp := []wktResult{
		{
			Lat: "12.34567",
			Lon: "-12.34567",
		},
		{
			Lat: "98.76543",
			Lon: "-98.76543",
		},
	}

	geocodeDataJSON, err := json.Marshal(geocodeResp)
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(geocodeDataJSON)),
			}, nil
		},
	}

	payload := GetWKTLocationPayload{UserLocation: "New York, NY"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/wkt", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/wkt", handler.GetWKTLocation)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal the response JSON
	var response GetWKTLocationResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the success message
	assert.Equal(t, "POINT(-12.34567 12.34567)", response.WKTLocation)
}

func TestGetAddressFromWKTHandler_InvalidJson(t *testing.T) {
	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	// Define an invalid payload
	payload := "{invalid_json}"

	// Create a test request with the JSON payload
	req, _ := http.NewRequest("POST", "/location/address", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/address", handler.GetAddressFromWKT)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Invalid JSON payload", response.Error)
}

func TestGetAddressFromWKTHandler_IncompletePayload(t *testing.T) {
	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := GetAddressFromWKTPayload{}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/address", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/address", handler.GetAddressFromWKT)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Invalid payload: missing wkt_location field", response.Error)
}

func TestGetAddressFromWKTHandler_MissingGeocodeApiKey(t *testing.T) {
	os.Unsetenv("GEOCODE_API_KEY")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := GetAddressFromWKTPayload{WKTLocation: "POINT(-12.34567 12.34567)"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/address", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/address", handler.GetAddressFromWKT)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Something went wrong, please try again later", response.Error)
}

func TestGetAddressFromWKTHandler_MissingGeocodeUrl(t *testing.T) {
	os.Setenv("GEOCODE_API_KEY", "some-key")
	os.Unsetenv("GEOCODE_URL")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := GetAddressFromWKTPayload{WKTLocation: "POINT(-12.34567 12.34567)"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/address", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/address", handler.GetAddressFromWKT)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Something went wrong, please try again later", response.Error)
}

func TestGetAddressFromWKTHandler_InvalidWKTLocation(t *testing.T) {
	os.Setenv("GEOCODE_API_KEY", "some-key")
	os.Setenv("GEOCODE_URL", "http://geocode.url")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := GetAddressFromWKTPayload{WKTLocation: "POINT(-12.34567)"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/address", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/address", handler.GetAddressFromWKT)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Param 'wkt_location' is not a valid WKT location", response.Error)
}

func TestGetAddressFromWKTHandler_GetError(t *testing.T) {
	os.Setenv("GEOCODE_API_KEY", "some-key")
	os.Setenv("GEOCODE_URL", "http://geocode.url")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
		Get:    func(url string) (*http.Response, error) { return nil, errors.New("http get error") },
	}

	payload := GetAddressFromWKTPayload{WKTLocation: "POINT(-12.34567 12.34567)"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/address", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/address", handler.GetAddressFromWKT)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Failed to get geocode location", response.Error)
}

func TestGetAddressFromWKTHandler_GetInternalServerError(t *testing.T) {
	os.Setenv("GEOCODE_API_KEY", "some-key")
	os.Setenv("GEOCODE_URL", "http://geocode.url")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusInternalServerError}, nil
		},
	}

	payload := GetAddressFromWKTPayload{WKTLocation: "POINT(-12.34567 12.34567)"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/address", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/address", handler.GetAddressFromWKT)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Failed to get geocode location", response.Error)
}

func TestGetAddressFromWKTHandler_ErrorDecodingBody(t *testing.T) {
	os.Setenv("GEOCODE_API_KEY", "some-key")
	os.Setenv("GEOCODE_URL", "http://geocode.url")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("{invalid_json}")),
			}, nil
		},
	}

	payload := GetAddressFromWKTPayload{WKTLocation: "POINT(-12.34567 12.34567)"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/address", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/address", handler.GetAddressFromWKT)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "Failed to get geocode location", response.Error)
}

func TestGetAddressFromWKTHandler_EmptyGeocodeResponse(t *testing.T) {
	os.Setenv("GEOCODE_API_KEY", "some-key")
	os.Setenv("GEOCODE_URL", "http://geocode.url")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	geocodeResp := GetAddressFromWKTResponse{
		Address: "",
	}

	geocodeDataJSON, err := json.Marshal(geocodeResp)
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(geocodeDataJSON)),
			}, nil
		},
	}

	payload := GetAddressFromWKTPayload{WKTLocation: "POINT(-12.34567 12.34567)"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/address", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/address", handler.GetAddressFromWKT)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal the response JSON
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the error message
	assert.Equal(t, "No geocode data found for the location provided", response.Error)
}

func TestGetAddressFromWKTHandler_Success(t *testing.T) {
	os.Setenv("GEOCODE_API_KEY", "some-key")
	os.Setenv("GEOCODE_URL", "http://geocode.url")

	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	geocodeResp := GetAddressFromWKTResponse{
		Address: "New York, NY",
	}

	geocodeDataJSON, err := json.Marshal(geocodeResp)
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(geocodeDataJSON)),
			}, nil
		},
	}

	payload := GetAddressFromWKTPayload{WKTLocation: "POINT(-12.34567 12.34567)"}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/location/address", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/location/address", handler.GetAddressFromWKT)
	r.ServeHTTP(w, req)

	// Assert the HTTP response code
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal the response JSON
	var response GetAddressFromWKTResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Assert the success message
	assert.Equal(t, "New York, NY", response.Address)
}

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetCurrentTimeHandler_InvalidJson(t *testing.T) {

	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := "{invalid_json}"

	req, _ := http.NewRequest("POST", "/time/current", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.POST("/time/current", handler.GetCurrentTime)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.Equal(t, "Invalid JSON payload", response.Error)
}

func TestGetCurrentTime_IncompletePayload(t *testing.T) {

	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := GetCurrentTimePayload{}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/time/current", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.POST("/time/current", handler.GetCurrentTime)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.Equal(t, "Invalid payload: missing timezone field", response.Error)
}

func TestGetCurrentTime_InvalidTimezone(t *testing.T) {

	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := GetCurrentTimePayload{
		Timezone: "Invalid/Timezone",
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/time/current", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.POST("/time/current", handler.GetCurrentTime)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.Equal(t, "Invalid timezone", response.Error)
}

func TestGetCurrentTime_Success(t *testing.T) {

	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := GetCurrentTimePayload{
		Timezone: "Europe/Bratislava",
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/time/current", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.POST("/time/current", handler.GetCurrentTime)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response GetCurrentTimeResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.NotNil(t, response.Time)

	expectedFormat := "2006-01-02 15:04:05"
	_, parseErr := time.Parse(expectedFormat, response.Time)
	assert.Nil(t, parseErr, "The time should be in the expected format")
}

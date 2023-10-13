package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAddHandlerValidInput(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal("Logger initialization failed")
	}

	tests := []struct {
		payload  AddPayload
		expected float64
		code     int
	}{
		{payload: AddPayload{Numbers: []float64{1.0, 2.0, 3.0}}, expected: 6.0, code: http.StatusOK},
		{payload: AddPayload{Numbers: []float64{0.0, 0.0, 0.0}}, expected: 0.0, code: http.StatusOK},
		{payload: AddPayload{Numbers: []float64{1.0, -2.0, 3.0}}, expected: 2.0, code: http.StatusOK},
	}

	for _, tt := range tests {
		r := gin.New() // Create a new Gin router for each test case
		handler := &Handler{
			Logger: logger,
		}

		payloadJSON, err := json.Marshal(tt.payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(payloadJSON))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.POST("/add", handler.Add)
		r.ServeHTTP(w, req)

		if tt.code != w.Code {
			t.Errorf("Expected status code %d, but got %d", tt.code, w.Code)
		}

		var response SuccessResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Error unmarshaling response: %v", err)
		}

		if tt.expected != response.Result {
			t.Errorf("Expected result %f, but got %f", tt.expected, response.Result)
		}
	}
}

func TestAddHandlerInvalidInput(t *testing.T) {
	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal("Logger initialization failed")
	}

	handler := &Handler{
		Logger: logger,
	}

	// Define an invalid payload (missing 'numbers' field)
	payload := map[string]interface{}{}

	// Marshal the payload to JSON
	payloadJSON, _ := json.Marshal(payload)

	// Create a test request with the JSON payload
	req, _ := http.NewRequest("POST", "/add", bytes.NewBuffer(payloadJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/add", handler.Add)
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
	assert.Equal(t, "Missing 'numbers' field in the payload", response.Error)
}

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAddHandler_InvalidJson(t *testing.T) {

	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := "{invalid_json}"

	req, _ := http.NewRequest("POST", "/add", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.POST("/add", handler.Add)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.Equal(t, "Invalid JSON payload", response.Error)
}

func TestAddHandler_MissingInput(t *testing.T) {

	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := map[string]interface{}{}

	payloadJSON, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/add", bytes.NewBuffer(payloadJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.POST("/add", handler.Add)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.Equal(t, "Invalid payload: missing numbers field", response.Error)
}

func TestAddHandler_ValidInput(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
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
		r := gin.New()
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

func TestSubtractHandler_InvalidJson(t *testing.T) {

	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := "{invalid_json}"

	req, _ := http.NewRequest("POST", "/subtract", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.POST("/subtract", handler.Subtract)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.Equal(t, "Invalid JSON payload", response.Error)
}

func TestSubtractHandler_MissingInput(t *testing.T) {

	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := map[string]interface{}{}

	payloadJSON, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/subtract", bytes.NewBuffer(payloadJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.POST("/subtract", handler.Subtract)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.Equal(t, "Invalid payload: missing subtract field", response.Error)
}

func TestSubtractHandler_ValidInput(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		payload  SubtractPayload
		expected float64
		code     int
	}{
		{payload: SubtractPayload{Number: 10.0, Subtract: []float64{5.0}}, expected: 5.0, code: http.StatusOK},
		{payload: SubtractPayload{Number: 10.0, Subtract: []float64{5.0, 5.0}}, expected: 0.0, code: http.StatusOK},
		{payload: SubtractPayload{Number: 10.0, Subtract: []float64{5.0, 5.0, 5.0}}, expected: -5.0, code: http.StatusOK},
	}

	for _, tt := range tests {
		r := gin.New()
		handler := &Handler{
			Logger: logger,
		}

		payloadJSON, err := json.Marshal(tt.payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/subtract", bytes.NewBuffer(payloadJSON))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.POST("/subtract", handler.Subtract)
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

func TestComputeHandler_InvalidJson(t *testing.T) {

	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := "{invalid_json}"

	req, _ := http.NewRequest("POST", "/compute", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.POST("/compute", handler.Compute)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.Equal(t, "Invalid JSON payload", response.Error)
}

func TestComputeHandler_MissingInput(t *testing.T) {

	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := map[string]interface{}{}

	payloadJSON, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/compute", bytes.NewBuffer(payloadJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.POST("/compute", handler.Compute)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.Equal(t, "Invalid payload: missing add and subtract fields", response.Error)
}

func TestComputeHandler_ValidInput(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		payload  ComputePayload
		expected float64
		code     int
	}{
		{payload: ComputePayload{Add: []float64{5.0}, Subtract: []float64{5.0}}, expected: 0.0, code: http.StatusOK},
		{payload: ComputePayload{Add: []float64{5.0, 10.0}, Subtract: []float64{5.0, 5.0}}, expected: 5.0, code: http.StatusOK},
		{payload: ComputePayload{Add: []float64{5.0, 10.0, 10.0}, Subtract: []float64{10.0, 10.0, 10.0}}, expected: -5.0, code: http.StatusOK},
	}

	for _, tt := range tests {
		r := gin.New()
		handler := &Handler{
			Logger: logger,
		}

		payloadJSON, err := json.Marshal(tt.payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/compute", bytes.NewBuffer(payloadJSON))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.POST("/compute", handler.Compute)
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

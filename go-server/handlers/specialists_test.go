package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/acornak/healthcare-poc/models"
	"github.com/acornak/healthcare-poc/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestFindSpecialistHandler_InvalidJson(t *testing.T) {
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
	req, _ := http.NewRequest("POST", "/specialist", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/specialist", handler.FindSpecialist)
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

func TestFindSpecialistHandler_IncompletePayload(t *testing.T) {
	// Create a new Gin router and handler
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	// Define an incomplete payload
	payload := FindSpecialistPayload{}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/specialist", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a test recorder
	w := httptest.NewRecorder()

	// Handle the request
	r.POST("/specialist", handler.FindSpecialist)
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
	assert.Equal(t, "Invalid payload: missing specialty_id, radius, user_location", response.Error)
}

func TestFindSpecialistHandler_SqlError(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM specialist").WithArgs(1, "POINT(-71.060316 48.432044)", 10).WillReturnError(errors.New("mocked error"))

	modelsDB := models.NewModels(db)
	payload := FindSpecialistPayload{SpecialtyId: 1, Radius: 10, UserLocation: "POINT(-71.060316 48.432044)"}

	r := gin.New() // Create a new Gin router for each test case
	handler := &Handler{
		Logger: logger,
		Models: modelsDB,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/specialist", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.POST("/specialist", handler.FindSpecialist)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.Equal(t, "mocked error", response.Error)
}

func TestFindSpecialistHandler_EmptyResponse(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email"})

	mock.ExpectQuery("SELECT (.+) FROM specialist").WithArgs(1, "POINT(-71.060316 48.432044)", 10).WillReturnRows(rows)

	modelsDB := models.NewModels(db)
	payload := FindSpecialistPayload{SpecialtyId: 1, Radius: 10, UserLocation: "POINT(-71.060316 48.432044)"}

	r := gin.New() // Create a new Gin router for each test case
	handler := &Handler{
		Logger: logger,
		Models: modelsDB,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/specialist", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.POST("/specialist", handler.FindSpecialist)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response FindSpecialistResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.Nil(t, response.Specialists)
}

func TestFindSpecialistHandler_Success(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	specialist := types.Specialist{
		ID:          1,
		Name:        "Dr. John Doe",
		SpecialtyID: 1,
		Location:    "POINT(-71.060316 48.432044)",
		Address:     "123 Main St, Boston, MA 02110",
		Url:         "https://example.com",
		Telephone:   "123-456-7890",
		Email:       "john@doe.com",
	}

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email"}).
		AddRow(specialist.ID, specialist.Name, specialist.SpecialtyID, specialist.Location, specialist.Address, specialist.Url, specialist.Telephone, specialist.Email)

	mock.ExpectQuery("SELECT (.+) FROM specialist").WithArgs(1, "POINT(-71.060316 48.432044)", 10).WillReturnRows(rows)

	modelsDB := models.NewModels(db)
	payload := FindSpecialistPayload{SpecialtyId: 1, Radius: 10, UserLocation: "POINT(-71.060316 48.432044)"}

	r := gin.New() // Create a new Gin router for each test case
	handler := &Handler{
		Logger: logger,
		Models: modelsDB,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/specialist", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.POST("/specialist", handler.FindSpecialist)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response FindSpecialistResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	fields := []struct {
		expected, got interface{}
		name          string
	}{
		{specialist.ID, response.Specialists[0].ID, "ID"},
		{specialist.Name, response.Specialists[0].Name, "Name"},
		{specialist.SpecialtyID, response.Specialists[0].SpecialtyID, "SpecialtyID"},
		{specialist.Location, response.Specialists[0].Location, "Location"},
		{specialist.Address, response.Specialists[0].Address, "Address"},
		{specialist.Url, response.Specialists[0].Url, "Url"},
		{specialist.Telephone, response.Specialists[0].Telephone, "Telephone"},
		{specialist.Email, response.Specialists[0].Email, "Email"},
	}

	for _, field := range fields {
		assert.Equal(t, field.expected, field.got)
	}
}

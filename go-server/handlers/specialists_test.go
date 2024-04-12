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
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

	payload := "{invalid_json}"

	req, _ := http.NewRequest("POST", "/specialist", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.POST("/specialist", handler.FindSpecialist)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.Equal(t, "Invalid JSON payload", response.Error)
}

func TestFindSpecialistHandler_IncompletePayload(t *testing.T) {
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Logger: logger,
	}

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

	w := httptest.NewRecorder()

	r.POST("/specialist", handler.FindSpecialist)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

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

	r := gin.New()
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

	r := gin.New()
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
		Name:        "John Doe",
		SpecialtyID: 1,
		Location:    "New York",
		Address:     "123 Main St",
		Url:         "https://example.com",
		Telephone:   "123-456-7890",
		Email:       "me@example.com",
		Monday:      "7:00 - 12:00, 13:00 - 15:00",
		Tuesday:     "7:00 - 12:00, 13:00 - 15:00",
		Wednesday:   "7:00 - 12:00, 13:00 - 15:00",
		Thursday:    "7:00 - 12:00, 13:00 - 15:00",
		Friday:      "7:00 - 12:00, 13:00 - 15:00",
		Saturday:    "",
		Sunday:      "",
	}

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}).
		AddRow(specialist.ID, specialist.Name, specialist.SpecialtyID, specialist.Location, specialist.Address, specialist.Url, specialist.Telephone, specialist.Email, specialist.Monday, specialist.Tuesday, specialist.Wednesday, specialist.Thursday, specialist.Friday, specialist.Saturday, specialist.Sunday)

	mock.ExpectQuery("SELECT (.+) FROM specialist").WithArgs(1, "POINT(-71.060316 48.432044)", 10).WillReturnRows(rows)

	modelsDB := models.NewModels(db)
	payload := FindSpecialistPayload{SpecialtyId: 1, Radius: 10, UserLocation: "POINT(-71.060316 48.432044)"}

	r := gin.New()
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
		{specialist.Monday, response.Specialists[0].Monday, "Monday"},
		{specialist.Tuesday, response.Specialists[0].Tuesday, "Tuesday"},
		{specialist.Wednesday, response.Specialists[0].Wednesday, "Wednesday"},
		{specialist.Thursday, response.Specialists[0].Thursday, "Thursday"},
		{specialist.Friday, response.Specialists[0].Friday, "Friday"},
		{specialist.Saturday, response.Specialists[0].Saturday, "Saturday"},
		{specialist.Sunday, response.Specialists[0].Sunday, "Sunday"},
	}

	for _, field := range fields {
		assert.Equal(t, field.expected, field.got)
	}
}

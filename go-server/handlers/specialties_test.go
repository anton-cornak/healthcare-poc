package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/acornak/healthcare-poc/models"
	"github.com/acornak/healthcare-poc/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetSpecialtiesHandler_SqlError(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM specialty").WithoutArgs().WillReturnError(errors.New("mocked error"))

	modelsDB := models.NewModels(db)

	r := gin.New() // Create a new Gin router for each test case
	handler := &Handler{
		Logger: logger,
		Models: modelsDB,
	}

	req, err := http.NewRequest("POST", "/specialties", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.POST("/specialties", handler.GetSpecialties)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.Equal(t, "mocked error", response.Error)
}

func TestGetSpecialtiesHandler_EmptyResponse(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description"})

	mock.ExpectQuery("SELECT (.+) FROM specialty").WithoutArgs().WillReturnRows(rows)

	modelsDB := models.NewModels(db)

	r := gin.New() // Create a new Gin router for each test case
	handler := &Handler{
		Logger: logger,
		Models: modelsDB,
	}

	req, err := http.NewRequest("POST", "/specialties", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.POST("/specialties", handler.GetSpecialties)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response GetSpecialtiesResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	assert.Nil(t, response.Specialties)
}

func TestGetSpecialtiesHandler_Success(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	specialty := types.Specialty{
		ID:          1,
		Name:        "Dr. John Doe",
		Description: "Specialty description",
	}

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(specialty.ID, specialty.Name, specialty.Description)

	mock.ExpectQuery("SELECT (.+) FROM specialty").WithoutArgs().WillReturnRows(rows)

	modelsDB := models.NewModels(db)

	r := gin.New() // Create a new Gin router for each test case
	handler := &Handler{
		Logger: logger,
		Models: modelsDB,
	}
	req, err := http.NewRequest("POST", "/specialties", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.POST("/specialties", handler.GetSpecialties)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response GetSpecialtiesResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	fields := []struct {
		expected, got interface{}
		name          string
	}{
		{specialty.ID, response.Specialties[0].ID, "ID"},
		{specialty.Name, response.Specialties[0].Name, "Name"},
		{specialty.Description, response.Specialties[0].Description, "SpecialtyID"},
	}

	for _, field := range fields {
		assert.Equal(t, field.expected, field.got)
	}

}

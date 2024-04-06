package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/acornak/healthcare-poc/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	// Set up test environment variables
	os.Setenv("PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASS", "password")
	os.Setenv("DB_NAME", "database")
	os.Setenv("SSL_MODE", "disable")

	// Run the tests
	code := m.Run()

	// Clean up environment variables after the tests
	os.Unsetenv("PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASS")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("SSL_MODE")

	// Exit with the appropriate exit code
	os.Exit(code)
}

func TestMainRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := zap.NewExample()

	// Create a sample database instance
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	// Create the server
	s := newServer(logger, &handlers.Handler{Logger: logger})

	// Assert that the returned server instance is not nil
	if s == nil || s.Router == nil {
		t.Error("server or router instance is nil")
	}

	tests := []struct {
		method string
		path   string
		code   int
	}{
		{"GET", "/api/v1/math/add", http.StatusNotFound},

		{"POST", "/api/v1/math/add", http.StatusBadRequest},
		{"POST", "/api/v1/math/subtract", http.StatusBadRequest},
		{"POST", "/api/v1/math/compute", http.StatusBadRequest},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		s.Router.ServeHTTP(rr, req)

		assert.Equal(t, tt.code, rr.Code, "Status code should match")
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Create a config instance to hold the loaded values
	cfg := config{}

	// Call the function to load config from environment variables
	err := loadConfigFromEnv(&cfg)
	if err != nil {
		t.Errorf("Unexpected error loading config: %v", err)
	}

	// Validate the loaded config
	err = validateConfig(&cfg)
	if err != nil {
		t.Errorf("Loaded config is invalid: %v", err)
	}

	// Perform assertions on the loaded config
	if cfg.port != "8080" {
		t.Errorf("Expected port to be '8080', got '%s'", cfg.port)
	}
	if cfg.dbConn.host != "localhost" {
		t.Errorf("Expected DB_HOST to be 'localhost', got '%s'", cfg.dbConn.host)
	}
	if cfg.dbConn.port != "5432" {
		t.Errorf("Expected DB_PORT to be '5432', got '%s'", cfg.dbConn.port)
	}
	if cfg.dbConn.user != "user" {
		t.Errorf("Expected DB_USER to be 'user', got '%s'", cfg.dbConn.user)
	}
	if cfg.dbConn.password != "password" {
		t.Errorf("Expected DB_PASS to be 'password', got '%s'", cfg.dbConn.password)
	}
	if cfg.dbConn.dbname != "database" {
		t.Errorf("Expected DB_NAME to be 'database', got '%s'", cfg.dbConn.dbname)
	}
	if cfg.dbConn.sslmode != "disable" {
		t.Errorf("Expected SSL_MODE to be 'disable', got '%s'", cfg.dbConn.sslmode)
	}
}

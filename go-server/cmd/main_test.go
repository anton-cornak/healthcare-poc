package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMainRoutes(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal("Logger initialization failed")
	}
	server := NewServer(logger)
	if server == nil || server.Router == nil {
		t.Fatal("Server or server.Router is nil")
	}

	gin.SetMode(gin.TestMode)

	tests := []struct {
		method string
		path   string
		code   int
	}{
		{"GET", "/add", http.StatusNotFound},

		{"POST", "/add", http.StatusBadRequest},
		{"POST", "/subtract", http.StatusBadRequest},
		{"POST", "/compute", http.StatusBadRequest},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		server.Router.ServeHTTP(rr, req)

		assert.Equal(t, tt.code, rr.Code, "Status code should match")
	}
}

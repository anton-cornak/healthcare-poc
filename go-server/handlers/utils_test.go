package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatAndLonFromWKT(t *testing.T) {
	tests := []struct {
		name        string
		wkt         string
		wantLat     string
		wantLon     string
		expectError bool
	}{
		{
			name:        "Valid POINT positive",
			wkt:         "POINT(12.34536 34.56789)",
			wantLat:     "12.34536",
			wantLon:     "34.56789",
			expectError: false,
		},
		{
			name:        "Valid POINT negative",
			wkt:         "POINT(-12.34536 -34.56789)",
			wantLat:     "-12.34536",
			wantLon:     "-34.56789",
			expectError: false,
		},
		{
			name:        "Invalid format",
			wkt:         "INVALID(-12.34536 -34.56789)",
			wantLat:     "",
			wantLon:     "",
			expectError: true,
		},
		{
			name:        "Missing coordinates",
			wkt:         "POINT()",
			wantLat:     "",
			wantLon:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLat, gotLon, err := getLatAndLonFromWKT(tt.wkt)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLat, gotLat, "Latitude does not match")
				assert.Equal(t, tt.wantLon, gotLon, "Longitude does not match")
			}
		})
	}
}

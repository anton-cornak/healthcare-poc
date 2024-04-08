package scrapers

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetCurrentTimeHandler_MissingUrl(t *testing.T) {
	os.Unsetenv("SCRAPER_TIME_URL")

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	scraper := &Scraper{
		Logger: logger,
	}

	err = scraper.GetSpecialists()
	assert.Equal(t, "SCRAPER_SPECIALISTS_URL not set", err.Error())
}

func TestGetCurrentTimeHandler_GetError(t *testing.T) {
	os.Setenv("SCRAPER_SPECIALISTS_URL", "http://example.com")

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	scraper := &Scraper{
		Logger: logger,
		Get:    func(url string) (*http.Response, error) { return nil, errors.New("http get error") },
	}

	err = scraper.GetSpecialists()
	assert.Equal(t, "http get error", err.Error())
}

func TestGetCurrentTimeHandler_GetInternalServerError(t *testing.T) {
	os.Setenv("SCRAPER_SPECIALISTS_URL", "http://example.com")

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	scraper := &Scraper{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusInternalServerError}, nil
		},
	}

	err = scraper.GetSpecialists()
	assert.Equal(t, "error getting specialists", err.Error())
}

func TestGetCurrentTimeHandler_ErrorDecodingBody(t *testing.T) {
	os.Setenv("SCRAPER_SPECIALISTS_URL", "http://example.com")

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	scraper := &Scraper{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("{invalid_json}")),
			}, nil
		},
	}

	err = scraper.GetSpecialists()
	assert.NotNil(t, err)
}

func TestGetCurrentTimeHandler_Success(t *testing.T) {
	os.Setenv("SCRAPER_SPECIALISTS_URL", "http://example.com")

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	scraper := &Scraper{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"sometype":"sometype"}`)),
			}, nil
		},
	}

	err = scraper.GetSpecialists()
	assert.Nil(t, err)
}

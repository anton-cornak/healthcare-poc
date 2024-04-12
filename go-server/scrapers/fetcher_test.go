package scrapers

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/acornak/healthcare-poc/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetSpecialists_MissingUrl(t *testing.T) {
	os.Unsetenv("SCRAPER_SPECIALISTS_URL")

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	scraper := &Scraper{
		Logger: logger,
	}

	_, err = scraper.GetSpecialists()
	assert.Equal(t, "SCRAPER_SPECIALISTS_URL not set", err.Error())
}

func TestGetSpecialists_GetError(t *testing.T) {
	os.Setenv("SCRAPER_SPECIALISTS_URL", "http://example.com")

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	scraper := &Scraper{
		Logger: logger,
		Get:    func(url string) (*http.Response, error) { return nil, errors.New("http get error") },
	}

	_, err = scraper.GetSpecialists()
	assert.Equal(t, "http get error", err.Error())
}

func TestGetSpecialists_GetInternalServerError(t *testing.T) {
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

	_, err = scraper.GetSpecialists()
	assert.Equal(t, "error getting specialists", err.Error())
}

func TestGetSpecialists_ErrorDecodingBody(t *testing.T) {
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

	_, err = scraper.GetSpecialists()
	assert.NotNil(t, err)
}

func TestGetSpecialists_Success(t *testing.T) {
	os.Setenv("SCRAPER_SPECIALISTS_URL", "http://example.com")

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	response := `{"features":[{"properties":{"id":1, "druh_zariadenia": "ortoped"}}]}`

	mock.ExpectExec("INSERT INTO specialty").WithArgs("ortoped").WillReturnResult(sqlmock.NewResult(1, 1))

	scraper := &Scraper{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(response)),
			}, nil
		},
		Models: models.NewModels(db),
	}

	resp, err := scraper.GetSpecialists()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp.Features))
	assert.Equal(t, 1, resp.Features[0].Properties.ID)
}
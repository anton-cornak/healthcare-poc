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

func TestScraperHandler_GetError(t *testing.T) {
	os.Setenv("SCRAPER_SPECIALISTS_URL", "http://example.com")

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	scraper := &Scraper{
		Logger: logger,
		Get:    func(url string) (*http.Response, error) { return nil, errors.New("http get error") },
	}

	err = scraper.ScrapeHandler()
	assert.Equal(t, "http get error", err.Error())
}

// TODO rewrite this to be per function
func TestScraperHandler_Success(t *testing.T) {
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
	rows := sqlmock.NewRows([]string{"id", "name", "description"})

	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE name=\$1`).WithArgs("ortoped").WillReturnRows(rows)
	mock.ExpectExec("INSERT INTO specialty").WithArgs("ortoped", "").WillReturnResult(sqlmock.NewResult(1, 1))

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

	err = scraper.ScrapeHandler()
	assert.Nil(t, err)
}

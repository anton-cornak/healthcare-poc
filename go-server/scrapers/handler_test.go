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
	"github.com/acornak/healthcare-poc/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInsertSpecialties_GetSpecialtyError(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM specialty WHERE name").WithArgs("ortoped").WillReturnError(errors.New("mocked error"))

	scraper := &Scraper{
		Logger: logger,
		Models: models.NewModels(db),
	}

	err = scraper.insertSpecialties([]struct {
		Properties types.GeoportalSpecialist
	}{{Properties: types.GeoportalSpecialist{Specialization: "ortoped"}}})

	assert.Equal(t, "mocked error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertSpecialties_SpecialtyFound(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(1, "ortoped", "test")

	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE name=\$1`).WithArgs("ortoped").WillReturnRows(rows)

	scraper := &Scraper{
		Logger: logger,
		Models: models.NewModels(db),
	}

	err = scraper.insertSpecialties([]struct {
		Properties types.GeoportalSpecialist
	}{{Properties: types.GeoportalSpecialist{Specialization: "ortoped"}}})

	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertSpecialties_InsertError(t *testing.T) {
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

	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE name=\$1`).WithArgs("ortoped").WillReturnRows(rows)
	mock.ExpectExec("INSERT INTO specialty").WithArgs("ortoped", "").WillReturnError(errors.New("mocked error"))

	scraper := &Scraper{
		Logger: logger,
		Models: models.NewModels(db),
	}

	err = scraper.insertSpecialties([]struct {
		Properties types.GeoportalSpecialist
	}{{Properties: types.GeoportalSpecialist{Specialization: "ortoped"}}})

	assert.Equal(t, "mocked error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertSpecialties_Success(t *testing.T) {
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

	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE name=\$1`).WithArgs("ortoped").WillReturnRows(rows)
	mock.ExpectExec("INSERT INTO specialty").WithArgs("ortoped", "").WillReturnResult(sqlmock.NewResult(1, 1))

	scraper := &Scraper{
		Logger: logger,
		Models: models.NewModels(db),
	}

	err = scraper.insertSpecialties([]struct {
		Properties types.GeoportalSpecialist
	}{{Properties: types.GeoportalSpecialist{Specialization: "ortoped"}}})

	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScraperHandler_GetSpecialistsError(t *testing.T) {
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

func TestScraperHandler_InsertSpecialtiesError(t *testing.T) {
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

	mock.ExpectQuery("SELECT (.+) FROM specialty WHERE name").WithArgs("ortoped").WillReturnError(errors.New("mocked error"))

	resp := `{"features":[{"properties":{"id":1, "druh_zariadenia": "ortoped"}}]}`

	scraper := &Scraper{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(resp)),
			}, nil
		},
		Models: models.NewModels(db),
	}

	err = scraper.ScrapeHandler()

	assert.Equal(t, "mocked error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScraperHandler_GetSpecialistByNameError(t *testing.T) {
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

	rows := sqlmock.NewRows([]string{"id", "name", "description"})

	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE name=\$1`).WithArgs("ortoped").WillReturnRows(rows)
	mock.ExpectExec("INSERT INTO specialty").WithArgs("ortoped", "").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE name=`).WithArgs("John Doe, Md.").WillReturnError(errors.New("mocked error"))

	resp := `{"features":[{"properties":{"id":1, "nazov_zariadenia": "John Doe, Md.", "druh_zariadenia": "ortoped"}}]}`

	scraper := &Scraper{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(resp)),
			}, nil
		},
		Models: models.NewModels(db),
	}

	err = scraper.ScrapeHandler()

	assert.Equal(t, "mocked error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScraperHandler_SpecialistFound(t *testing.T) {
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

	rowsSpecialty := sqlmock.NewRows([]string{"id", "name", "description"})
	rowsSpecialist := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}).
		AddRow(1, "John Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "me@example.com", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "", "")

	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE name=\$1`).WithArgs("ortoped").WillReturnRows(rowsSpecialty)
	mock.ExpectExec("INSERT INTO specialty").WithArgs("ortoped", "").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE name=`).WithArgs("John Doe, Md.").WillReturnRows(rowsSpecialist)

	resp := `{"features":[{"properties":{"id":1, "nazov_zariadenia": "John Doe, Md.", "druh_zariadenia": "ortoped"}}]}`

	scraper := &Scraper{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(resp)),
			}, nil
		},
		Models: models.NewModels(db),
	}

	err = scraper.ScrapeHandler()

	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScraperHandler_GetSpecialtyByNameError(t *testing.T) {
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

	rowsSpecialty := sqlmock.NewRows([]string{"id", "name", "description"})
	rowsSpecialist := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"})

	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE name=\$1`).WithArgs("ortoped").WillReturnRows(rowsSpecialty)
	mock.ExpectExec("INSERT INTO specialty").WithArgs("ortoped", "").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE name=`).WithArgs("John Doe, Md.").WillReturnRows(rowsSpecialist)
	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE name=\$1`).WithArgs("ortoped").WillReturnError(errors.New("mocked error"))

	resp := `{"features":[{"properties":{"id":1, "nazov_zariadenia": "John Doe, Md.", "druh_zariadenia": "ortoped"}}]}`

	scraper := &Scraper{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(resp)),
			}, nil
		},
		Models: models.NewModels(db),
	}

	err = scraper.ScrapeHandler()

	assert.Equal(t, "mocked error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScraperHandler_SpecialtyNotFound(t *testing.T) {
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

	rowsSpecialty := sqlmock.NewRows([]string{"id", "name", "description"})
	rowsSpecialist := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"})

	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE name=\$1`).WithArgs("ortoped").WillReturnRows(rowsSpecialty)
	mock.ExpectExec("INSERT INTO specialty").WithArgs("ortoped", "").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE name=`).WithArgs("John Doe, Md.").WillReturnRows(rowsSpecialist)
	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE name=\$1`).WithArgs("ortoped").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}))

	resp := `{"features":[{"properties":{"id":1, "nazov_zariadenia": "John Doe, Md.", "druh_zariadenia": "ortoped"}}]}`

	scraper := &Scraper{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(resp)),
			}, nil
		},
		Models: models.NewModels(db),
	}

	err = scraper.ScrapeHandler()

	assert.Equal(t, "specialty not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScraperHandler_InsertSpecialistError(t *testing.T) {
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

	rowsSpecialty := sqlmock.NewRows([]string{"id", "name", "description"})
	rowsSpecialist := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"})
	rowsSpecialtyByName := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(1, "ortoped", "")

	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE name=\$1`).WithArgs("ortoped").WillReturnRows(rowsSpecialty)
	mock.ExpectExec("INSERT INTO specialty").WithArgs("ortoped", "").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE name=`).WithArgs("John Doe, Md.").WillReturnRows(rowsSpecialist)
	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE name=\$1`).WithArgs("ortoped").WillReturnRows(rowsSpecialtyByName)

	mock.ExpectExec(`INSERT INTO specialist`).
		WithArgs("John Doe, Md.", 1, "POINT(0 0)", " ,  , Slovenská republika", "", ", ", "", "", "", "", "", "", "", "").
		WillReturnError(errors.New("mocked error"))

	resp := `{"features":[{"properties":{"id":1, "nazov_zariadenia": "John Doe, Md.", "druh_zariadenia": "ortoped"}}]}`

	scraper := &Scraper{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(resp)),
			}, nil
		},
		Models: models.NewModels(db),
	}

	err = scraper.ScrapeHandler()

	assert.Equal(t, "mocked error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

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

	rowsSpecialty := sqlmock.NewRows([]string{"id", "name", "description"})
	rowsSpecialist := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"})
	rowsSpecialtyByName := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(1, "ortoped", "")

	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE name=\$1`).WithArgs("ortoped").WillReturnRows(rowsSpecialty)
	mock.ExpectExec("INSERT INTO specialty").WithArgs("ortoped", "").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE name=`).WithArgs("John Doe, Md.").WillReturnRows(rowsSpecialist)
	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE name=\$1`).WithArgs("ortoped").WillReturnRows(rowsSpecialtyByName)
	mock.ExpectExec(`INSERT INTO specialist`).
		WithArgs("John Doe, Md.", 1, "POINT(0 0)", " ,  , Slovenská republika", "", ", ", "", "", "", "", "", "", "", "").
		WillReturnResult(sqlmock.NewResult(1, 1))

	resp := `{"features":[{"properties":{"id":1, "nazov_zariadenia": "John Doe, Md.", "druh_zariadenia": "ortoped"}}]}`

	scraper := &Scraper{
		Logger: logger,
		Get: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(resp)),
			}, nil
		},
		Models: models.NewModels(db),
	}

	err = scraper.ScrapeHandler()

	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

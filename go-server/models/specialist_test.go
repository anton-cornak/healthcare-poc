package models

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/acornak/healthcare-poc/types"
	"github.com/stretchr/testify/assert"
)

func TestGetAllSpecialists_SqlError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM specialist").WithoutArgs().WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetAllSpecialists()

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllSpecialists_RowsScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}).
		AddRow(1, "John Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "me@example.com", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "", "")
	rows.RowError(0, errors.New("rows scan error"))

	mock.ExpectQuery("SELECT (.+) FROM specialist").WithoutArgs().WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetAllSpecialists()

	assert.Error(t, err)
	assert.EqualError(t, err, "rows scan error")
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllSpecialists_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}).
		AddRow(1, "John Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "me@example.com", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "", "")

	mock.ExpectQuery("SELECT (.+) FROM specialist").WithoutArgs().WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetAllSpecialists()

	expected := []*types.Specialist{
		{
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
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistBySpecialty_SqlError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE specialty_id=`).WithArgs(1).WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialistBySpecialty(1)

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistBySpecialty_RowsScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}).
		AddRow(1, "John Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "me@example.com", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "", "").
		AddRow(2, "Jane Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "me@example.com", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "", "")

	rows.RowError(0, errors.New("rows scan error"))

	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE specialty_id=`).WithArgs(1).WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialistBySpecialty(1)

	assert.Error(t, err)
	assert.EqualError(t, err, "rows scan error")
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistBySpecialty_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}).
		AddRow(1, "John Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "me@example.com", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "", "").
		AddRow(2, "Jane Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "jane@example.com", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "", "")

	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE specialty_id=`).WithArgs(1).WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialistBySpecialty(1)

	expected := []*types.Specialist{
		{
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
		},
		{
			ID:          2,
			Name:        "Jane Doe",
			SpecialtyID: 1,
			Location:    "New York",
			Address:     "123 Main St",
			Url:         "https://example.com",
			Telephone:   "123-456-7890",
			Email:       "jane@example.com",
			Monday:      "7:00 - 12:00, 13:00 - 15:00",
			Tuesday:     "7:00 - 12:00, 13:00 - 15:00",
			Wednesday:   "7:00 - 12:00, 13:00 - 15:00",
			Thursday:    "7:00 - 12:00, 13:00 - 15:00",
			Friday:      "7:00 - 12:00, 13:00 - 15:00",
			Saturday:    "",
			Sunday:      "",
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistByID_SqlError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE id=`).WithArgs(1).WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialistByID(1)

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistByID_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"})

	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE id=`).WithArgs(1).WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialistByID(1)

	assert.NoError(t, err)
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}).
		AddRow(1, "John Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "me@example.com", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "", "")

	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE id=`).WithArgs(1).WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialistByID(1)

	expected := &types.Specialist{
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

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistByName_SqlError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE name=`).WithArgs("test").WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialistByName("test")

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistByName_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"})

	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE name=`).WithArgs("test").WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialistByName("test")

	assert.NoError(t, err)
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistByName_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}).
		AddRow(1, "John Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "me@example.com", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "", "")

	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE name=`).WithArgs("test").WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialistByName("test")

	expected := &types.Specialist{
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

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertSpecialist_SqlError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	s := types.Specialist{
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

	mock.ExpectExec(`INSERT INTO specialist`).
		WithArgs(s.Name, s.SpecialtyID, s.Location, s.Address, s.Url, s.Telephone, s.Email, s.Monday, s.Tuesday, s.Wednesday, s.Thursday, s.Friday, s.Saturday, s.Sunday).
		WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	err = modelsDB.DB.InsertSpecialist(s)

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertSpecialist_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	s := types.Specialist{
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

	mock.ExpectExec(`INSERT INTO specialist`).
		WithArgs(s.Name, s.SpecialtyID, s.Location, s.Address, s.Url, s.Telephone, s.Email, s.Monday, s.Tuesday, s.Wednesday, s.Thursday, s.Friday, s.Saturday, s.Sunday).
		WillReturnResult(sqlmock.NewResult(1, 1))

	modelsDB := NewModels(db)
	err = modelsDB.DB.InsertSpecialist(s)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSpecialist_SqlError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectExec(`DELETE FROM specialist WHERE id=\$1`).WithArgs(1).WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	err = modelsDB.DB.DeleteSpecialist(1)

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSpecialist_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectExec(`DELETE FROM specialist WHERE id=\$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	modelsDB := NewModels(db)
	err = modelsDB.DB.DeleteSpecialist(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateSpecialist_SqlError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	s := types.Specialist{
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

	mock.ExpectExec(`UPDATE specialist SET name=\$1, specialty_id=\$2, location=\$3, address=\$4, url=\$5, telephone=\$6, email=\$7, monday=\$8, tuesday=\$9, wednesday=\$10, thursday=\$11, friday=\$12, saturday=\$13, sunday=\$14 WHERE id=\$8`).
		WithArgs(s.Name, s.SpecialtyID, s.Location, s.Address, s.Url, s.Telephone, s.Email, s.ID, s.Monday, s.Tuesday, s.Wednesday, s.Thursday, s.Friday, s.Saturday, s.Sunday).
		WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	err = modelsDB.DB.UpdateSpecialist(s)

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateSpecialist_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	s := types.Specialist{
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

	mock.ExpectExec(`UPDATE specialist SET name=\$1, specialty_id=\$2, location=\$3, address=\$4, url=\$5, telephone=\$6, email=\$7, monday=\$8, tuesday=\$9, wednesday=\$10, thursday=\$11, friday=\$12, saturday=\$13, sunday=\$14 WHERE id=\$8`).
		WithArgs(s.Name, s.SpecialtyID, s.Location, s.Address, s.Url, s.Telephone, s.Email, s.ID, s.Monday, s.Tuesday, s.Wednesday, s.Thursday, s.Friday, s.Saturday, s.Sunday).
		WillReturnResult(sqlmock.NewResult(1, 1))

	modelsDB := NewModels(db)
	err = modelsDB.DB.UpdateSpecialist(s)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistBySpecialtyAndLocation_SqlError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM specialist").WithArgs(1, "123 Main St", 10000).WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialistBySpecialtyAndLocation(1, 10000, "123 Main St")

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistBySpecialtyAndLocation_RowsScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email"}).
		AddRow(1, "John Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "me@example.com")
	rows.RowError(0, errors.New("rows scan error"))

	mock.ExpectQuery("SELECT (.+) FROM specialist").WithArgs(1, "123 Main St", 10000).WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialistBySpecialtyAndLocation(1, 10000, "123 Main St")

	assert.Error(t, err)
	assert.EqualError(t, err, "rows scan error")
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistBySpecialtyAndLocation_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}).
		AddRow(1, "John Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "me@example.com", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "7:00 - 12:00, 13:00 - 15:00", "", "")

	mock.ExpectQuery("SELECT (.+) FROM specialist").WithArgs(1, "123 Main St", 10000).WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialistBySpecialtyAndLocation(1, 10000, "123 Main St")

	expected := []*types.Specialist{
		{
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
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

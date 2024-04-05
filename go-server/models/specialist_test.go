package models

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/acornak/healthcare-poc/types"
	"github.com/stretchr/testify/assert"
)

func TestAllSpecialists_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email"}).
		AddRow(1, "John Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "me@example.com")

	mock.ExpectQuery("SELECT (.+) FROM specialist").WithoutArgs().WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.AllSpecialists()

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
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAllSpecialists_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM specialist").WithoutArgs().WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.AllSpecialists()

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistBySpecialty_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email"}).
		AddRow(1, "John Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "me@example.com").
		AddRow(2, "Jane Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "jane@example.com")

	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE specialty_id=`).WithArgs(1).WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialistBySpecialty(1)

	expected := []*types.Specialist{
		{
			ID:   1,
			Name: "John Doe",

			SpecialtyID: 1,
			Location:    "New York",
			Address:     "123 Main St",
			Url:         "https://example.com",
			Telephone:   "123-456-7890",
			Email:       "me@example.com",
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
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistBySpecialty_Error(t *testing.T) {
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

func TestGetSpecialistByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email"}).
		AddRow(1, "John Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "me@example.com")

	mock.ExpectQuery(`SELECT (.+) FROM specialist WHERE id=`).WithArgs(1).WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialistByID(1)

	expected := &types.Specialist{
		ID:   1,
		Name: "John Doe",

		SpecialtyID: 1,
		Location:    "New York",
		Address:     "123 Main St",
		Url:         "https://example.com",
		Telephone:   "123-456-7890",
		Email:       "me@example.com",
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistByID_Error(t *testing.T) {
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
	}

	mock.ExpectExec(`INSERT INTO specialist`).WithArgs(s.Name, s.SpecialtyID, s.Location, s.Address, s.Url, s.Telephone, s.Email).WillReturnResult(sqlmock.NewResult(1, 1))

	modelsDB := NewModels(db)
	err = modelsDB.DB.InsertSpecialist(s)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertSpecialist_Error(t *testing.T) {
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
	}

	mock.ExpectExec(`INSERT INTO specialist`).WithArgs(s.Name, s.SpecialtyID, s.Location, s.Address, s.Url, s.Telephone, s.Email).WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	err = modelsDB.DB.InsertSpecialist(s)

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

func TestDeleteSpecialist_Error(t *testing.T) {
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
	}

	mock.ExpectExec(`UPDATE specialist SET name=\$1, specialty_id=\$2, location=\$3, address=\$4, url=\$5, telephone=\$6, email=\$7 WHERE id=\$8`).WithArgs(s.Name, s.SpecialtyID, s.Location, s.Address, s.Url, s.Telephone, s.Email, s.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	modelsDB := NewModels(db)
	err = modelsDB.DB.UpdateSpecialist(s)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateSpecialist_Error(t *testing.T) {
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
	}

	mock.ExpectExec(`UPDATE specialist SET name=\$1, specialty_id=\$2, location=\$3, address=\$4, url=\$5, telephone=\$6, email=\$7 WHERE id=\$8`).WithArgs(s.Name, s.SpecialtyID, s.Location, s.Address, s.Url, s.Telephone, s.Email, s.ID).WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	err = modelsDB.DB.UpdateSpecialist(s)

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistBySpecialtyAndLocation_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "specialty_id", "location", "address", "url", "telephone", "email"}).
		AddRow(1, "John Doe", 1, "New York", "123 Main St", "https://example.com", "123-456-7890", "me@example.com")

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
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialistBySpecialtyAndLocation_Error(t *testing.T) {
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

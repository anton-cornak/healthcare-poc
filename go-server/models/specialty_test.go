package models

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAllSpecialties_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(1, "test", "test")

	mock.ExpectQuery("SELECT (.+) FROM specialty").WithoutArgs().WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.AllSpecialties()

	expected := []*Specialty{
		{
			ID:          1,
			Name:        "test",
			Description: "test",
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAllSpecialties_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT (.+) FROM specialty`).WithoutArgs().WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.AllSpecialties()

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialtyByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(1, "test", "test")

	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE id=\$1`).WithArgs(1).WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialtyByID(1)

	expected := &Specialty{
		ID:          1,
		Name:        "test",
		Description: "test",
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSpecialtyByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT (.+) FROM specialty WHERE id=\$1`).WithArgs(1).WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetSpecialtyByID(1)

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertSpecialty_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	s := Specialty{
		Name:        "test",
		Description: "test",
	}

	mock.ExpectExec("INSERT INTO specialty").WithArgs(s.Name, s.Description).WillReturnResult(sqlmock.NewResult(1, 1))

	modelsDB := NewModels(db)
	err = modelsDB.DB.InsertSpecialty(s)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertSpecialty_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	s := Specialty{
		Name:        "test",
		Description: "test",
	}

	mock.ExpectExec("INSERT INTO specialty").WithArgs(s.Name, s.Description).WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	err = modelsDB.DB.InsertSpecialty(s)

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSpecialty_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM specialties").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	modelsDB := NewModels(db)
	err = modelsDB.DB.DeleteSpecialty(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSpecialty_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM specialties").WithArgs(1).WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	err = modelsDB.DB.DeleteSpecialty(1)

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateSpecialty_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	s := Specialty{
		ID:          1,
		Name:        "test",
		Description: "test",
	}

	mock.ExpectExec(`UPDATE specialty SET name=\$1, description=\$2 WHERE id=\$3`).WithArgs(s.Name, s.Description, s.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	modelsDB := NewModels(db)
	err = modelsDB.DB.UpdateSpecialty(s)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateSpecialty_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	s := Specialty{
		ID:          1,
		Name:        "test",
		Description: "test",
	}

	mock.ExpectExec(`UPDATE specialty SET name=\$1, description=\$2 WHERE id=\$3`).WithArgs(s.Name, s.Description, s.ID).WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	err = modelsDB.DB.UpdateSpecialty(s)

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

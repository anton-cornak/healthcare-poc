package models

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/acornak/healthcare-poc/types"
	"github.com/stretchr/testify/assert"
)

func TestAllReviews_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "specialist_id", "url", "rating", "comment"}).AddRow(1, 1, "test", 4.5, "test")

	mock.ExpectQuery("SELECT (.+) FROM review").WithoutArgs().WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.AllReviews()

	expected := []*types.Review{
		{
			ID:           1,
			SpecialistId: 1,
			Url:          "test",
			Rating:       4.5,
			Comment:      "test",
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAllReviews_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT (.+) FROM review`).WithoutArgs().WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.AllReviews()

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetReviewBySpecialistId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "specialist_id", "url", "rating", "comment"}).AddRow(1, 1, "test", 4.5, "test")

	mock.ExpectQuery(`SELECT (.+) FROM review WHERE specialist_id=\$1`).WithArgs(1).WillReturnRows(rows)

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetReviewBySpecialistId(1)

	expected := &types.Review{
		ID:           1,
		SpecialistId: 1,
		Url:          "test",
		Rating:       4.5,
		Comment:      "test",
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetReviewBySpecialistId_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT (.+) FROM review WHERE specialist_id=\$1`).WithArgs(1).WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	res, err := modelsDB.DB.GetReviewBySpecialistId(1)

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.Nil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertReview_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	r := types.Review{
		SpecialistId: 1,
		Url:          "test",
		Rating:       4.5,
		Comment:      "test",
	}

	mock.ExpectExec("INSERT INTO review").WithArgs(r.SpecialistId, r.Url, r.Rating, r.Comment).WillReturnResult(sqlmock.NewResult(1, 1))

	modelsDB := NewModels(db)
	err = modelsDB.DB.InsertReview(r)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertReview_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	r := types.Review{
		SpecialistId: 1,
		Url:          "test",
		Rating:       4.5,
		Comment:      "test",
	}

	mock.ExpectExec("INSERT INTO review").WithArgs(r.SpecialistId, r.Url, r.Rating, r.Comment).WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	err = modelsDB.DB.InsertReview(r)

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteReview_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM review").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	modelsDB := NewModels(db)
	err = modelsDB.DB.DeleteReview(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteReview_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM review").WithArgs(1).WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	err = modelsDB.DB.DeleteReview(1)

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateReview_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	r := types.Review{
		SpecialistId: 1,
		Url:          "test",
		Rating:       4.5,
		Comment:      "test",
	}

	mock.ExpectExec(`UPDATE review SET specialist_id=\$1, url=\$2, rating=\$3, comment=\$4 WHERE id=\$5`).WithArgs(r.SpecialistId, r.Url, r.Rating, r.Comment, r.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	modelsDB := NewModels(db)
	err = modelsDB.DB.UpdateReview(r)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateReview_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	r := types.Review{
		SpecialistId: 1,
		Url:          "test",
		Rating:       4.5,
		Comment:      "test",
	}

	mock.ExpectExec(`UPDATE review SET specialist_id=\$1, url=\$2, rating=\$3, comment=\$4 WHERE id=\$5`).WithArgs(r.SpecialistId, r.Url, r.Rating, r.Comment, r.ID).WillReturnError(errors.New("mocked error"))

	modelsDB := NewModels(db)
	err = modelsDB.DB.UpdateReview(r)

	assert.Error(t, err)
	assert.EqualError(t, err, "mocked error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

package models

import "github.com/acornak/healthcare-poc/types"

/*
AllReviews returns all reviews from the database
The function returns a slice of pointers to Review structs
The function returns an error if there was an issue with the database
*/
func (m *DBModel) AllReviews() ([]*types.Review, error) {
	stmt := `
	SELECT id, specialist_id, url, rating, comment
	FROM review
	`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*types.Review

	for rows.Next() {
		var s types.Review
		rows.Scan(&s.ID, &s.SpecialistId, &s.Url, &s.Rating, &s.Comment)
		reviews = append(reviews, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

/*
GetReview returns a review from the database with a specific id
The id is the id of the review
The function returns a pointer to a Review struct
The function returns an error if there was an issue with the database
*/
func (m *DBModel) GetReviewBySpecialistId(id int) (*types.Review, error) {
	stmt := `
	SELECT id, specialist_id, url, rating, comment
	FROM review
	WHERE specialist_id=$1
	`

	row := m.DB.QueryRow(stmt, id)

	var r types.Review
	err := row.Scan(&r.ID, &r.SpecialistId, &r.Url, &r.Rating, &r.Comment)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

/*
GetReviewByID returns a review from the database with a specific id
The id is the id of the review
The function returns a pointer to a Review struct
The function returns an error if there was an issue with the database
*/
func (m *DBModel) InsertReview(r types.Review) error {
	stmt := `
	INSERT INTO review (specialist_id, url, rating, comment)
	VALUES ($1, $2, $3, $4)
	`

	_, err := m.DB.Exec(stmt, r.SpecialistId, r.Url, r.Rating, r.Comment)
	if err != nil {
		return err
	}

	return nil
}

/*
DeleteReview deletes a review from the database with a specific id
The id is the id of the review
The function returns an error if there was an issue with the database
*/
func (m *DBModel) DeleteReview(id int) error {
	stmt := `
	DELETE FROM review
	WHERE id = $1
	`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

/*
UpdateReview updates a review in the database
The r parameter is a Review struct
The function returns an error if there was an issue with the database
*/
func (m *DBModel) UpdateReview(r types.Review) error {
	stmt := `
	UPDATE review
	SET specialist_id=$1, url=$2, rating=$3, comment=$4
	WHERE id=$5
	`

	_, err := m.DB.Exec(stmt, r.SpecialistId, r.Url, r.Rating, r.Comment, r.ID)
	if err != nil {
		return err
	}

	return nil
}

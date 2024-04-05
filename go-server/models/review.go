package models

type Review struct {
	ID           int     `json:"id"`
	SpecialistId int     `json:"specialist_id"`
	Url          string  `json:"url"`
	Rating       float64 `json:"rating"`
	Comment      string  `json:"comment,omitempty"`
}

func (m *DBModel) AllReviews() ([]*Review, error) {
	stmt := "SELECT id, specialist_id, url, rating, comment FROM review"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*Review

	for rows.Next() {
		var s Review
		err = rows.Scan(&s.ID, &s.SpecialistId, &s.Url, &s.Rating, &s.Comment)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (m *DBModel) GetReviewBySpecialistId(id int) (*Review, error) {
	stmt := `SELECT id, specialist_id, url, rating, comment FROM review WHERE specialist_id=$1`

	row := m.DB.QueryRow(stmt, id)

	var r Review
	err := row.Scan(&r.ID, &r.SpecialistId, &r.Url, &r.Rating, &r.Comment)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (m *DBModel) InsertReview(r Review) error {
	stmt := `INSERT INTO review (specialist_id, url, rating, comment) VALUES ($1, $2, $3, $4)`

	_, err := m.DB.Exec(stmt, r.SpecialistId, r.Url, r.Rating, r.Comment)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) DeleteReview(id int) error {
	stmt := `DELETE FROM review WHERE id = $1`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) UpdateReview(r Review) error {
	stmt := `UPDATE review SET specialist_id=$1, url=$2, rating=$3, comment=$4 WHERE id=$5`

	_, err := m.DB.Exec(stmt, r.SpecialistId, r.Url, r.Rating, r.Comment, r.ID)
	if err != nil {
		return err
	}

	return nil
}

package models

type Specialty struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (m *DBModel) AllSpecialties() ([]*Specialty, error) {
	stmt := "SELECT id, name, description FROM specialty"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specialties []*Specialty

	for rows.Next() {
		var s Specialty
		err = rows.Scan(&s.ID, &s.Name, &s.Description)
		if err != nil {
			return nil, err
		}
		specialties = append(specialties, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return specialties, nil
}

func (m *DBModel) GetSpecialtyByID(id int) (*Specialty, error) {
	stmt := `SELECT id, name, description FROM specialty WHERE id=$1`

	row := m.DB.QueryRow(stmt, id)

	var s Specialty
	err := row.Scan(&s.ID, &s.Name, &s.Description)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (m *DBModel) InsertSpecialty(s Specialty) error {
	stmt := `INSERT INTO specialty (name, description) VALUES ($1, $2)`

	_, err := m.DB.Exec(stmt, s.Name, s.Description)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) DeleteSpecialty(id int) error {
	stmt := `DELETE FROM specialties WHERE id = $1`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) UpdateSpecialty(s Specialty) error {
	stmt := `UPDATE specialty SET name=$1, description=$2 WHERE id=$3`

	_, err := m.DB.Exec(stmt, s.Name, s.Description, s.ID)
	if err != nil {
		return err
	}

	return nil
}

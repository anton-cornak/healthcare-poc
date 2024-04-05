package models

type Specialist struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	SpecialtyID int    `json:"specialty_id"`
	Location    string `json:"location,omitempty"`
	Address     string `json:"address,omitempty"`
	Url         string `json:"url,omitempty"`
	Telephone   string `json:"telephone,omitempty"`
	Email       string `json:"email,omitempty"`
}

func (m *DBModel) AllSpecialists() ([]*Specialist, error) {
	stmt := "SELECT id, name, specialty_id, location, address, url, telephone, email FROM specialist"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specialists []*Specialist

	for rows.Next() {
		var s Specialist
		err = rows.Scan(&s.ID, &s.Name, &s.SpecialtyID, &s.Location, &s.Address, &s.Url, &s.Telephone, &s.Email)
		if err != nil {
			return nil, err
		}
		specialists = append(specialists, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return specialists, nil
}

func (m *DBModel) GetSpecialistBySpecialty(specialty_id int) ([]*Specialist, error) {
	stmt := `SELECT id, name, specialty_id, location, address, url, telephone, email FROM specialist WHERE specialty_id=$1`

	rows, err := m.DB.Query(stmt, specialty_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specialists []*Specialist

	for rows.Next() {
		var s Specialist
		err = rows.Scan(&s.ID, &s.Name, &s.SpecialtyID, &s.Location, &s.Address, &s.Url, &s.Telephone, &s.Email)
		if err != nil {
			return nil, err
		}
		specialists = append(specialists, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return specialists, nil
}

func (m *DBModel) GetSpecialistByID(id int) (*Specialist, error) {
	stmt := `SELECT id, name, specialty_id, location, address, url, telephone, email FROM specialist WHERE id=$1`

	row := m.DB.QueryRow(stmt, id)

	var s Specialist
	err := row.Scan(&s.ID, &s.Name, &s.SpecialtyID, &s.Location, &s.Address, &s.Url, &s.Telephone, &s.Email)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (m *DBModel) InsertSpecialist(s Specialist) error {
	stmt := `INSERT INTO specialist (name, specialty_id, location, address, url, telephone, email) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.Exec(stmt, s.Name, s.SpecialtyID, s.Location, s.Address, s.Url, s.Telephone, s.Email)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) DeleteSpecialist(id int) error {
	stmt := `DELETE FROM specialist WHERE id=$1`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) UpdateSpecialist(s Specialist) error {
	stmt := `UPDATE specialist SET name=$1, specialty_id=$2, location=$3, address=$4, url=$5, telephone=$6, email=$7 WHERE id=$8`

	_, err := m.DB.Exec(stmt, s.Name, s.SpecialtyID, s.Location, s.Address, s.Url, s.Telephone, s.Email, s.ID)
	if err != nil {
		return err
	}

	return nil
}

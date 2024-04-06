package models

import "github.com/acornak/healthcare-poc/types"

/*
AllSpecialties returns all specialties from the database
The function returns a slice of pointers to Specialty structs
The function returns an error if there was an issue with the database
*/
func (m *DBModel) AllSpecialties() ([]*types.Specialty, error) {
	stmt := `
	SELECT id, name, description
	FROM specialty
	`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specialties []*types.Specialty

	for rows.Next() {
		var s types.Specialty
		rows.Scan(&s.ID, &s.Name, &s.Description)
		specialties = append(specialties, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return specialties, nil
}

/*
GetSpecialtyByID returns a specialty from the database with a specific id
The id is the id of the specialty
The function returns a pointer to a Specialty struct
The function returns an error if there was an issue with the database
*/
func (m *DBModel) GetSpecialtyByID(id int) (*types.Specialty, error) {
	stmt := `
	SELECT id, name, description
	FROM specialty
	WHERE id=$1
	`

	row := m.DB.QueryRow(stmt, id)

	var s types.Specialty
	err := row.Scan(&s.ID, &s.Name, &s.Description)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

/*
GetSpecialtyByName returns a specialty from the database with a specific name
The name is the name of the specialty
The function returns an error if there was an issue with the database
*/
func (m *DBModel) InsertSpecialty(s types.Specialty) error {
	stmt := `
	INSERT INTO specialty (name, description)
	VALUES ($1, $2)
	`

	_, err := m.DB.Exec(stmt, s.Name, s.Description)
	if err != nil {
		return err
	}

	return nil
}

/*
DeleteSpecialty deletes a specialty from the database with a specific id
The id is the id of the specialty
The function returns an error if there was an issue with the database
*/
func (m *DBModel) DeleteSpecialty(id int) error {
	stmt := `
	DELETE FROM specialties
	WHERE id = $1
	`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

/*
UpdateSpecialty updates a specialty in the database
The s parameter is a Specialty struct
The function returns an error if there was an issue with the database
*/
func (m *DBModel) UpdateSpecialty(s types.Specialty) error {
	stmt := `
	UPDATE specialty
	SET name=$1, description=$2
	WHERE id=$3
	`

	_, err := m.DB.Exec(stmt, s.Name, s.Description, s.ID)
	if err != nil {
		return err
	}

	return nil
}

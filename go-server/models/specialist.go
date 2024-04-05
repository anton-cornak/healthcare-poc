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

// AllSpecialists returns all specialists from the database
func (m *DBModel) AllSpecialists() ([]*Specialist, error) {
	stmt := `
	SELECT id, name, specialty_id, location, address, url, telephone, email
	FROM specialist
	`

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

// GetSpecialistBySpecialty returns all specialists from the database with a specific specialty
// The specialtyID is the id of the specialty
// The function returns a slice of pointers to Specialist structs
// The function returns an error if there was an issue with the database
func (m *DBModel) GetSpecialistBySpecialty(specialtyID int) ([]*Specialist, error) {
	stmt := `
	SELECT id, name, specialty_id, location, address, url, telephone, email
	FROM specialist
	WHERE specialty_id=$1
	`

	rows, err := m.DB.Query(stmt, specialtyID)
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

/*
GetSpecialistByID returns a specialist from the database with a specific id
The id is the id of the specialist
The function returns a pointer to a Specialist struct
The function returns an error if there was an issue with the database
*/
func (m *DBModel) GetSpecialistByID(id int) (*Specialist, error) {
	stmt := `
	SELECT id, name, specialty_id, location, address, url, telephone, email
	FROM specialist
	WHERE id=$1
	`

	row := m.DB.QueryRow(stmt, id)

	var s Specialist
	err := row.Scan(&s.ID, &s.Name, &s.SpecialtyID, &s.Location, &s.Address, &s.Url, &s.Telephone, &s.Email)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

/*
GetSpecialistBySpecialtyAndLocation returns all specialists from the database with a specific specialty and within a certain radius of a location
The specialtyID is the id of the specialty
The radius is the radius in meters
The userLocation is the location in WKT format
The function returns a slice of pointers to Specialist structs
The function returns an error if there was an issue with the database
*/
func (m *DBModel) GetSpecialistBySpecialtyAndLocation(specialtyID, radius int, userLocation string) ([]*Specialist, error) {
	stmt := `
    SELECT id, name, specialty_id, location, address, url, telephone, email 
    FROM specialist 
    WHERE specialty_id=$1 AND ST_DWithin(location, ST_GeogFromText($2), $3)
    `

	rows, err := m.DB.Query(stmt, specialtyID, userLocation, radius)
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

/*
InsertSpecialist inserts a new specialist into the database
The s parameter is a Specialist struct
The function returns an error if there was an issue with the database
*/
func (m *DBModel) InsertSpecialist(s Specialist) error {
	stmt := `
	INSERT INTO specialist (name, specialty_id, location, address, url, telephone, email)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

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

/*
UpdateSpecialist updates a specialist in the database
The s parameter is a Specialist struct
The function returns an error if there was an issue with the database
*/
func (m *DBModel) UpdateSpecialist(s Specialist) error {
	stmt := `
	UPDATE specialist
	SET name=$1, specialty_id=$2, location=$3, address=$4, url=$5, telephone=$6, email=$7
	WHERE id=$8
	`

	_, err := m.DB.Exec(stmt, s.Name, s.SpecialtyID, s.Location, s.Address, s.Url, s.Telephone, s.Email, s.ID)
	if err != nil {
		return err
	}

	return nil
}

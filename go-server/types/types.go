package types

/*
Review represents a review of a specialist
The struct contains the following fields:
- ID: the id of the review
- SpecialistId: the id of the specialist
- Url: the url of the review
- Rating: the rating of the review
- Comment: the comment of the review
*/
type Review struct {
	ID           int     `json:"id"`
	SpecialistId int     `json:"specialist_id"`
	Url          string  `json:"url"`
	Rating       float64 `json:"rating"`
	Comment      string  `json:"comment,omitempty"`
}

/*
Specialist represents a specialist
The struct contains the following fields:
- ID: the id of the specialist
- Name: the name of the specialist
- SpecialtyID: the id of the specialty
- Location: the location of the specialist in the WKT format
- Address: the address of the specialist
- Url: the url of the specialist
- Telephone: the telephone of the specialist
- Email: the email of the specialist
*/
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

/*
Specialty represents a specialty
The struct contains the following fields:
- ID: the id of the specialty
- Name: the name of the specialty
- Description: the description of the specialty
*/
type Specialty struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

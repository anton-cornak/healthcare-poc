package types

import (
	"strconv"
	"strings"
)

type GeoportalSpecialist struct {
	ID             int       `json:"id"`
	Identifier     string    `json:"identifikator"`
	KPZS           string    `json:"kpzs"`
	Specialization string    `json:"druh_zariadenia"`
	Name           string    `json:"nazov_zariadenia"`
	Latitude       float64   `json:"poloha_lat"`
	Longitude      float64   `json:"poloha_lon"`
	Address        string    `json:"addressline"`
	Municipality   string    `json:"municipality"`
	BuildingNumber string    `json:"buildingnumber"`
	County         string    `json:"county"`
	StreetName     string    `json:"streetname"`
	PostalCode     string    `json:"postalcode"`
	Email          string    `json:"email"`
	Cellphone      string    `json:"mobil"`
	Phone          string    `json:"telefon"`
	Staff          string    `json:"odborni_zastupcovia"`
	MondayHours    string    `json:"pondelok"`
	TuesdayHours   string    `json:"utorok"`
	WednesdayHours string    `json:"streda"`
	ThursdayHours  string    `json:"stvrtok"`
	FridayHours    string    `json:"piatok"`
	SaturdayHours  string    `json:"sobota"`
	SundayHours    string    `json:"nedela"`
	AbsenceFrom    string    `json:"nepritomnost_od"`
	AbsenceTo      string    `json:"nepritomnost_do"`
	Info           string    `json:"info"`
	Union          string    `json:"union"`
	Vszp           string    `json:"vszp"`
	Dovera         string    `json:"dovera"`
	Bbox           []float64 `json:"bbox"`
}

func (g *GeoportalSpecialist) getWKTLocation() string {
	lat := strconv.FormatFloat(g.Latitude, 'f', -1, 64)
	lon := strconv.FormatFloat(g.Longitude, 'f', -1, 64)
	return "POINT(" + lon + " " + lat + ")"
}

func (g *GeoportalSpecialist) getAddress() string {
	if g.Address == "" {
		return g.StreetName + " " + g.BuildingNumber + ", " + g.PostalCode + " " + g.Municipality + ", Slovenská republika"
	}
	return g.Address
}

func (g *GeoportalSpecialist) getSpecialistNames() string {
	toBeRemoved := []string{
		" ako lekár",
		" ako sestra",
		" ako iný zdravotnícky pracovník - psychológ",
		" ako zubný lekár",
		" ako iný zdravotnícky pracovník - logopéd",
		" ako iný zdravotnícky pracovník - liečebný pedagóg",
		" ako dentálna hygienička",
		" ako zdravotnícky laborant",
		" ako pôrodná asistentka",
		" ako rádiologický technik",
	}

	for _, remove := range toBeRemoved {
		g.Staff = strings.Replace(g.Staff, remove, "", -1)
	}

	return g.Staff
}

func (g *GeoportalSpecialist) getSpecialistPhones() string {
	return g.Phone + ", " + g.Cellphone
}

func (g *GeoportalSpecialist) getUnion() bool {
	if g.Union == "áno" {
		return true
	} else {
		return false
	}
}

func (g *GeoportalSpecialist) getVszp() bool {
	if g.Vszp == "áno" {
		return true
	} else {
		return false
	}
}

func (g *GeoportalSpecialist) getDovera() bool {
	if g.Dovera == "áno" {
		return true
	} else {
		return false
	}
}

func (g *GeoportalSpecialist) CastToDbType(specialtyID int) Specialist {
	return Specialist{
		Name:        g.Name,
		SpecialtyID: specialtyID,
		Location:    g.getWKTLocation(),
		Address:     g.getAddress(),
		Telephone:   g.getSpecialistPhones(),
		Email:       g.Email,
		Monday:      g.MondayHours,
		Tuesday:     g.TuesdayHours,
		Wednesday:   g.WednesdayHours,
		Thursday:    g.ThursdayHours,
		Friday:      g.FridayHours,
		Saturday:    g.SaturdayHours,
		Sunday:      g.SundayHours,
	}
}

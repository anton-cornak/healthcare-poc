package scrapers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestCase() geoportalSpecialist {
	testCase := geoportalSpecialist{
		ID:             1,
		Identifier:     "68-44869223-A0002",
		KPZS:           "P27489001201",
		Specialization: "ambulancia vnútorného lekárstva",
		Name:           "Ambulancia vnútorného lekárstva, MUDr. Milena Zidanova, Michalovce, (Milzid s.r.o.)",
		Latitude:       48.789456,
		Longitude:      21.123456,
		Address:        "Ulica 9, 07101 Michalovce, Slovenská republika",
		Municipality:   "Michalovce",
		BuildingNumber: "9",
		County:         "Michalovce",
		StreetName:     "Ulica",
		PostalCode:     "07101",
		Email:          "me@example.com",
		Cellphone:      "+421 987 654 321",
		Phone:          "+421 123 456 789",
		Staff:          "MUDr. Milena Zidanova ako lekár, Zita Triuma ako sestra, Jozef Kralik ako iný zdravotnícky pracovník - psychológ, MUDr. Jana Kralikova ako zubný lekár",
		MondayHours:    "7:00 - 13:00, 13:30 - 15:00",
		TuesdayHours:   "7:00 - 13:00, 13:30 - 15:00",
		WednesdayHours: "7:00 - 13:00, 13:30 - 15:00",
		ThursdayHours:  "7:00 - 13:00, 13:30 - 15:00",
		FridayHours:    "7:00 - 13:00, 13:30 - 15:00",
		SaturdayHours:  "",
		SundayHours:    "",
		AbsenceFrom:    "",
		AbsenceTo:      "",
		Info:           "",
		Union:          "áno",
		Vszp:           "áno",
		Dovera:         "nie",
		Bbox:           []float64{1.0, 2.0, 3.0, 4.0},
	}

	testCaseCopy := testCase
	return testCaseCopy
}

func TestGetWKTLocation(t *testing.T) {
	testCase := setupTestCase()

	expected := "POINT(48.789456 21.123456)"
	actual := testCase.getWKTLocation()
	assert.Equal(t, expected, actual)
}

func TestGetAddress(t *testing.T) {
	testCase := setupTestCase()

	expected := "Ulica 9, 07101 Michalovce, Slovenská republika"
	actual := testCase.getAddress()
	assert.Equal(t, expected, actual)

	testCase.Address = ""
	expected = "Ulica 9, 07101 Michalovce, Slovenská republika"
	actual = testCase.getAddress()
	assert.Equal(t, expected, actual)

	testCase.Address = expected
}

func TestGetSpecialistNames(t *testing.T) {
	testCase := setupTestCase()

	expected := "MUDr. Milena Zidanova, Zita Triuma, Jozef Kralik, MUDr. Jana Kralikova"
	actual := testCase.getSpecialistNames()
	assert.Equal(t, expected, actual)
}

func TestGetSpecialistPhones(t *testing.T) {
	testCase := setupTestCase()

	expected := "+421 123 456 789, +421 987 654 321"
	actual := testCase.getSpecialistPhones()
	assert.Equal(t, expected, actual)
}

func TestGetUnion(t *testing.T) {
	testCase := setupTestCase()

	expected := true
	actual := testCase.getUnion()
	assert.Equal(t, expected, actual)

	testCase.Union = "nie"
	expected = false
	actual = testCase.getUnion()
	assert.Equal(t, expected, actual)
}

func TestGetVszp(t *testing.T) {
	testCase := setupTestCase()

	expected := true
	actual := testCase.getVszp()
	assert.Equal(t, expected, actual)

	testCase.Vszp = "nie"
	expected = false
	actual = testCase.getVszp()
	assert.Equal(t, expected, actual)
}

func TestGetDovera(t *testing.T) {
	testCase := setupTestCase()

	expected := false
	actual := testCase.getDovera()
	assert.Equal(t, expected, actual)

	testCase.Dovera = "áno"
	expected = true
	actual = testCase.getDovera()
	assert.Equal(t, expected, actual)
}

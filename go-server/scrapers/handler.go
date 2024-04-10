package scrapers

import (
	"github.com/acornak/healthcare-poc/types"
	"go.uber.org/zap"
)

func (s *Scraper) insertSpecialties(specialists []struct {
	Properties geoportalSpecialist
}) error {
	specialtiesMap := make(map[string]bool)

	for _, specialist := range specialists {
		specialtiesMap[specialist.Properties.Specialization] = true
	}

	for specialty := range specialtiesMap {
		typedSpecialty := types.Specialty{Name: specialty}

		found, err := s.Models.DB.GetSpecialtyByName(typedSpecialty.Name)
		if err != nil {
			return err
		}

		if found == nil {
			err = s.Models.DB.InsertSpecialty(typedSpecialty)
			if err != nil {
				return err
			}

			s.Logger.Info("Inserted specialty", zap.String("name", specialty))
		}
	}

	return nil
}

func (s *Scraper) ScrapeHandler() error {
	// Scrape specialists and save to DB
	specialists, err := s.GetSpecialists()
	if err != nil {
		return err
	}

	// get all existing specialties
	err = s.insertSpecialties(specialists.Features)
	if err != nil {
		return err
	}

	return nil
}

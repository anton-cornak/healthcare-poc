package scrapers

import (
	"errors"

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
		castedSpecialty := types.Specialty{Name: specialty}

		found, err := s.Models.DB.GetSpecialtyByName(castedSpecialty.Name)
		if err != nil {
			return err
		}

		if found == nil {
			err = s.Models.DB.InsertSpecialty(castedSpecialty)
			if err != nil {
				return err
			}

			s.Logger.Info("specialty inserted", zap.String("name", castedSpecialty.Name))
		}

	}

	return nil
}

func (s *Scraper) ScrapeHandler() error {
	// Scrape data from geoportal API
	specialists, err := s.GetSpecialists()
	if err != nil {
		return err
	}

	// get all existing specialties
	err = s.insertSpecialties(specialists.Features)
	if err != nil {
		return err
	}

	for _, specialist := range specialists.Features {
		// check if specialist already exists
		found, err := s.Models.DB.GetSpecialistByName(specialist.Properties.Name)
		if err != nil {
			return err
		}

		if found != nil {
			continue
		}

		// get specialty by name
		specialty, err := s.Models.DB.GetSpecialtyByName(specialist.Properties.Specialization)
		if err != nil {
			return err
		}

		if specialty == nil {
			return errors.New("specialty not found")
		}

		// insert specialist
		castedSpecialist := specialist.Properties.castToDbType(specialty.ID)
		err = s.Models.DB.InsertSpecialist(castedSpecialist)
		if err != nil {
			return err
		}

	}

	return nil
}

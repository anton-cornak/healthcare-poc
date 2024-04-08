package scrapers

import (
	"encoding/json"
	"errors"
	"os"

	"go.uber.org/zap"
)

type GetSpecialistsResponse struct {
	Sometype string `json:"sometype"`
}

func (s *Scraper) GetSpecialists() error {
	url := os.Getenv("SCRAPER_SPECIALISTS_URL")
	if url == "" {
		s.Logger.Error("SCRAPER_SPECIALISTS_URL not set")
		return errors.New("SCRAPER_SPECIALISTS_URL not set")
	}

	resp, err := s.Get(url)
	if err != nil {
		s.Logger.Error("Error getting specialists", zap.Error(err))
		return err
	}

	if resp.StatusCode != 200 {
		s.Logger.Error("Error getting specialists", zap.Int("status_code", resp.StatusCode))
		return errors.New("error getting specialists")
	}

	defer resp.Body.Close()

	var specialistsData GetSpecialistsResponse
	if err := json.NewDecoder(resp.Body).Decode(&specialistsData); err != nil {
		s.Logger.Error("Error decoding body", zap.Error(err))
		return err
	}

	s.Logger.Info("Specialists scraped")

	return nil
}

package scrapers

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/acornak/healthcare-poc/types"
	"go.uber.org/zap"
)

type GetSpecialistsResponse struct {
	Features []struct {
		Properties types.GeoportalSpecialist
	}
}

func (s *Scraper) GetSpecialists() (GetSpecialistsResponse, error) {
	url := os.Getenv("SCRAPER_SPECIALISTS_URL")
	if url == "" {
		s.Logger.Error("SCRAPER_SPECIALISTS_URL not set")
		return GetSpecialistsResponse{}, errors.New("SCRAPER_SPECIALISTS_URL not set")
	}

	resp, err := s.Get(url)
	if err != nil {
		s.Logger.Error("Error getting specialists", zap.Error(err))
		return GetSpecialistsResponse{}, err
	}

	if resp.StatusCode != 200 {
		s.Logger.Error("Error getting specialists", zap.Int("status_code", resp.StatusCode))
		return GetSpecialistsResponse{}, errors.New("error getting specialists")
	}

	defer resp.Body.Close()

	var specialistsData GetSpecialistsResponse
	if err := json.NewDecoder(resp.Body).Decode(&specialistsData); err != nil {
		s.Logger.Error("Error decoding body", zap.Error(err))
		return GetSpecialistsResponse{}, err
	}

	return specialistsData, nil
}

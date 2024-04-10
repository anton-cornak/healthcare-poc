package scrapers

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/acornak/healthcare-poc/models"
	"go.uber.org/zap"
)

type Scraper struct {
	Logger *zap.Logger
	Models models.Models
	Get    func(url string) (resp *http.Response, err error)
}

func NewScraper(logger *zap.Logger, models models.Models) *Scraper {
	return &Scraper{
		Logger: logger,
		Models: models,
		Get: func(url string) (*http.Response, error) {
			customTransport := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			httpClient := &http.Client{
				Transport: customTransport,
				Timeout:   30 * time.Second,
			}

			resp, err := httpClient.Get(url)
			if err != nil {
				return nil, err
			}
			return resp, nil
		},
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

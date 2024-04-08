package scrapers

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Scraper struct {
	Logger *zap.Logger
	Get    func(url string) (resp *http.Response, err error)
	Post   func(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

func NewScraper(logger *zap.Logger, Get func(url string) (resp *http.Response, err error)) *Scraper {
	return &Scraper{
		Logger: logger,
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

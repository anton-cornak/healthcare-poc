package scrapers

import (
	"testing"

	"github.com/acornak/healthcare-poc/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewScraper_NonExistingUrl(t *testing.T) {
	logger := zap.NewExample()
	scraper := NewScraper(logger, models.Models{})

	assert.NotNil(t, scraper)

	_, err := scraper.Get("http://notexistingurl.com")
	assert.NotNil(t, err)
}

func TestNewScraper_Success(t *testing.T) {
	logger := zap.NewExample()
	scraper := NewScraper(logger, models.Models{})

	assert.NotNil(t, scraper)

	_, err := scraper.Get("http://google.com")
	assert.Nil(t, err)
}

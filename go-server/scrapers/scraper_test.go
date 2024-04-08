package scrapers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewScraper(t *testing.T) {
	logger := zap.NewExample()
	scraper := NewScraper(logger, nil)

	assert.NotNil(t, scraper)
}

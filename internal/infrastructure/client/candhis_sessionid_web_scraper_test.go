package client_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tul1/candhis_api/internal/infrastructure/client"
	scrapermock "github.com/tul1/candhis_api/internal/infrastructure/client/scraper_mock"
	"go.uber.org/mock/gomock"
)

func TestGetCandhisSessionID_Failure(t *testing.T) {
	mockScraper := scrapermock.NewMockScraper(gomock.NewController(t))
	scraper := client.NewCandhisSessionIDWebScraper(mockScraper, "https://example.com")

	mockScraper.EXPECT().
		Run(gomock.Any(), "https://example.com", gomock.Any()).
		Return(errors.New("run failed"))

	_, err := scraper.GetCandhisSessionID(context.Background())
	assert.EqualError(t, err, "failed while running chromedp tasks to retrieve session id: run failed")
}

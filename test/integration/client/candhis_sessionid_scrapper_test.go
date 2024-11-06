package client_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/application/repository"
	"github.com/tul1/candhis_api/internal/infrastructure/client"
	"github.com/tul1/candhis_api/internal/pkg/chrome"
)

func TestGetCandhisSessionID_Success(t *testing.T) {
	ctx := context.Background()
	sessionScraper := setupCandhisSessionIDWebScraperClient(t)

	sessionID, err := sessionScraper.GetCandhisSessionID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, sessionID)

	assert.NotEmpty(t, sessionID.ID())
}

func setupCandhisSessionIDWebScraperClient(t *testing.T) repository.CandhisSessionIDWebScraper {
	t.Helper()

	chromeURL := os.Getenv("CHROME_URL")
	require.NotEmpty(t, chromeURL, "CHROME_URL must be set")

	targetWeb := os.Getenv("TARGET_WEB")
	require.NotEmpty(t, targetWeb, "TARGET_WEB must be set")

	chromeScraper, err := chrome.NewChromedpScraper(&http.Client{}, chromeURL)
	require.NoError(t, err)

	return client.NewCandhisSessionIDWebScraper(chromeScraper, targetWeb)
}

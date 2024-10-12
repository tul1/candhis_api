package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/application/repository"
	"github.com/tul1/candhis_api/internal/infrastructure/client"
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

	chromeID := getChromeID(t, chromeURL)
	require.NotEmpty(t, chromeID, "Chrome ID should not be empty")

	targetWeb := os.Getenv("TARGET_WEB")
	require.NotEmpty(t, targetWeb, "TARGET_WEB must be set")

	return client.NewCandhisSessionIDWebScraper(chromeURL, chromeID, targetWeb)
}

func getChromeID(t *testing.T, chromeURL string) string {
	t.Helper()

	resp, err := http.Get(fmt.Sprintf("http://%s/json/version", chromeURL))
	require.NoError(t, err, "failed to get Chrome version info")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read response body")

	var versionInfo struct {
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}
	err = json.Unmarshal(body, &versionInfo)
	require.NoError(t, err, "failed to unmarshal version info")

	parts := strings.Split(versionInfo.WebSocketDebuggerURL, "/")
	require.GreaterOrEqual(t, len(parts), 3, "invalid WebSocket URL format")

	return parts[len(parts)-1]
}

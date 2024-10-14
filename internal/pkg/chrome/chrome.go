package chrome

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const webSocketDebuggerURLFieldSplitedParts = 6

func GetChromeID(client *http.Client, chromeURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/json/version", chromeURL), http.NoBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request, url: %s, error: %w", chromeURL, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform request, url: %s, error: %w", chromeURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var versionInfo struct {
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}
	err = json.Unmarshal(body, &versionInfo)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal version info: %w", err)
	}

	parts := strings.Split(versionInfo.WebSocketDebuggerURL, "/")
	if len(parts) != webSocketDebuggerURLFieldSplitedParts {
		return "", fmt.Errorf("invalid WebSocket URL format: %s", versionInfo.WebSocketDebuggerURL)
	}

	return parts[len(parts)-1], nil
}

package chrome

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type chromedpScraper struct {
	chromodpWS string
}

func NewChromedpScraper(client *http.Client, chromeURL string) (*chromedpScraper, error) {
	if !strings.HasPrefix(chromeURL, "http://") && !strings.HasPrefix(chromeURL, "https://") {
		chromeURL = "http://" + chromeURL
	}

	chromeID, err := getChromeID(client, chromeURL)
	if err != nil {
		return nil, err
	}

	chromeURL = strings.TrimPrefix(chromeURL, "http://")
	chromeURL = strings.TrimPrefix(chromeURL, "https://")

	return &chromedpScraper{
		chromodpWS: fmt.Sprintf("ws://%s/devtools/browser/%s", chromeURL, chromeID),
	}, nil
}

const webSocketDebuggerURLFieldSplitedParts = 6

func getChromeID(client *http.Client, chromeURL string) (string, error) {
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

func (cs *chromedpScraper) Run(ctx context.Context, targetWeb string, actionFunc func(context.Context) error) error {
	ctx, cancel := chromedp.NewRemoteAllocator(ctx, cs.chromodpWS)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	return chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(targetWeb),
		chromedp.ActionFunc(actionFunc),
	)
}

package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/tul1/candhis_api/internal/application/model"
)

const (
	scrapingbeeURL = "https://app.scrapingbee.com/api/v1/"
	// In order to create and activate the sessionID cookie we need to click in any buttom of the web page.
	// The ID "#idBtnAr" is the ID of the buttom "Archives" and the ID "#idBtnTR" is the buttom "temps reel".
	scrapingbeeJSScenario = `{"instructions":[{"click":"#idBtnAr"},{"wait":1000},{"click":"#idBtnTR"},{"wait":1000}]}`
)

type scrapingBeeClient struct {
	client            *http.Client
	scrapingbeeAPIKey string
}

func NewScrapingBeeClient(client *http.Client, scrapingbeeAPIKey string) *scrapingBeeClient {
	return &scrapingBeeClient{client, scrapingbeeAPIKey}
}

func (s *scrapingBeeClient) GetCandhisSessionID() (model.CandhisSessionID, error) {
	params := url.Values{}
	params.Add("api_key", s.scrapingbeeAPIKey)
	params.Add("url", candhisURL)
	params.Add("js_scenario", scrapingbeeJSScenario)

	reqURL := fmt.Sprintf("%s?%s", scrapingbeeURL, params.Encode())
	req, err := http.NewRequest(http.MethodGet, reqURL, http.NoBody)
	if err != nil {
		return model.CandhisSessionID{}, fmt.Errorf("failed to create request, url: %s, error: %w", reqURL, err)
	}

	err = req.ParseForm()
	if err != nil {
		return model.CandhisSessionID{}, fmt.Errorf("failed to parse form, url: %s, error: %w", reqURL, err)
	}

	resp, err := s.client.Do(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	if err != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return model.CandhisSessionID{},
				fmt.Errorf("failed to read error response body, status: %d, error: %w, url: %s", resp.StatusCode, err, reqURL)
		}
		return model.CandhisSessionID{},
			fmt.Errorf("error response from server, status: %d, response: %s, url: %s", resp.StatusCode, string(body), reqURL)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return model.CandhisSessionID{},
				fmt.Errorf("failed to read error response body, status: %d, error: %w, url: %s", resp.StatusCode, err, reqURL)
		}
		return model.CandhisSessionID{},
			fmt.Errorf("error response from server, status: %d, response: %s, url: %s", resp.StatusCode, string(body), reqURL)
	}

	cookies := resp.Header["Set-Cookie"]
	for _, cookie := range cookies {
		if strings.Contains(cookie, "PHPSESSID") {
			phpSessionID := strings.Split(strings.TrimPrefix(cookie, "PHPSESSID="), ";")[0]
			return model.NewCandhisSessionID(phpSessionID, nil)
		}
	}

	return model.CandhisSessionID{}, fmt.Errorf("failed to retrieve cookie PHPSESSID, url: %s", reqURL)
}

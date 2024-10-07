package client_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tul1/candhis_api/internal/application/repository"
	"github.com/tul1/candhis_api/internal/infrastructure/client"
)

func TestGetCandhisSessionID_Success(t *testing.T) {
	mockHandler := func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Set-Cookie": {"PHPSESSID=valid-session-id; Path=/; HttpOnly"},
			},
		}
	}

	scrapingBeeClientClient := setupMockScrapingBeeClient(t, mockHandler)

	sessionID, err := scrapingBeeClientClient.GetCandhisSessionID()
	assert.NoError(t, err)
	assert.Equal(t, "valid-session-id", sessionID.ID())
	assert.NotNil(t, sessionID.CreatedAt())
}

func TestGetCandhisSessionID_NoCookie(t *testing.T) {
	mockHandler := func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
		}
	}

	scrapingBeeClientClient := setupMockScrapingBeeClient(t, mockHandler)

	_, err := scrapingBeeClientClient.GetCandhisSessionID()
	assert.ErrorContains(t, err, "failed to retrieve cookie PHPSESSID, url:")
}

func TestGetCandhisSessionID_RequestError(t *testing.T) {
	mockHandler := func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader(`{"message": "Internal Server Error"}`)),
		}
	}

	scrapingBeeClient := setupMockScrapingBeeClient(t, mockHandler)

	_, err := scrapingBeeClient.GetCandhisSessionID()
	assert.ErrorContains(t, err,
		`error response from server, status: 500, response: {"message": "Internal Server Error"}, url: `)
}

func setupMockScrapingBeeClient(t *testing.T, mockHandler func(req *http.Request) *http.Response) repository.ScrapingBeeClient {
	t.Helper()

	mockClient := &http.Client{
		Transport: &mockRoundTripper{mockHandler: mockHandler},
	}

	return client.NewScrapingBeeClient(mockClient, "test-api-key")
}

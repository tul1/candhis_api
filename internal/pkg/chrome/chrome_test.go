package chrome_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/pkg/chrome"
)

func MockHTTPResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestGetChromeID_Success(t *testing.T) {
	mockHandler := func(req *http.Request) *http.Response {
		return MockHTTPResponse(200, `{"webSocketDebuggerUrl": "ws://localhost:9222/devtools/browser/abc123"}`)
	}
	client := setupMockHTTPClient(mockHandler)

	chromeID, err := chrome.GetChromeID(client, "http://fake.url")
	require.NoError(t, err)
	assert.Equal(t, "abc123", chromeID, "chrome ID did not match expected value")
}

func TestGetChromeID_Failures(t *testing.T) {
	testCases := map[string]struct {
		mockResponse string
		expectedErr  string
	}{
		"invalid WebSocket URL format": {
			mockResponse: `{"webSocketDebuggerUrl": "invalid-websocket-url"}`,
			expectedErr:  "invalid WebSocket URL format: invalid-websocket-url",
		},
		"empty response": {
			mockResponse: ``,
			expectedErr:  "failed to unmarshal version info: unexpected end of JSON input",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockHandler := func(req *http.Request) *http.Response {
				return MockHTTPResponse(200, tc.mockResponse)
			}
			client := setupMockHTTPClient(mockHandler)

			_, err := chrome.GetChromeID(client, "http://fake.url")
			require.EqualError(t, err, tc.expectedErr)
		})
	}
}

type mockRoundTripper struct {
	mockHandler func(req *http.Request) *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.mockHandler(req), nil
}

func setupMockHTTPClient(mockHandler func(req *http.Request) *http.Response) *http.Client {
	return &http.Client{Transport: &mockRoundTripper{mockHandler: mockHandler}}
}
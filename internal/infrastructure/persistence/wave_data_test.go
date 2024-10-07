package persistence_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/assert"
	"github.com/tul1/candhis_api/internal/domain/model/modeltest"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence"
)

func TestAdd_Success(t *testing.T) {
	waveDataStore := setupMockWaveData(func(req *http.Request) (*http.Response, error) {
		return MockResponse(201, `{"result": "created"}`), nil
	})

	waveData := modeltest.MustCreateWaveData(t, "17/09/2024", "09:00", "0.6", "1.1", "4.7", "8", "32", "15")

	err := waveDataStore.Add(context.Background(), waveData, "test-index")
	assert.NoError(t, err)
}

func TestAdd_Error(t *testing.T) {
	waveDataStore := setupMockWaveData(func(req *http.Request) (*http.Response, error) {
		return MockResponse(500, `{"error": "internal server error"}`), nil
	})

	waveData := modeltest.MustCreateWaveData(t, "17/09/2024", "09:00", "0.6", "1.1", "4.7", "8", "32", "15")

	err := waveDataStore.Add(context.Background(), waveData, "test-index")
	assert.EqualError(t, err, `error indexing document: 500 Internal Server Error, body: {"error": "internal server error"}`)
}

func TestAdd_EmptyIndexName(t *testing.T) {
	waveDataStore := setupMockWaveData(func(req *http.Request) (*http.Response, error) {
		return MockResponse(201, `{"result": "created"}`), nil
	})

	waveData := modeltest.MustCreateWaveData(t, "17/09/2024", "09:00", "0.6", "1.1", "4.7", "8", "32", "15")

	err := waveDataStore.Add(context.Background(), waveData, "")
	assert.EqualError(t, err, "indexName cannot be empty")
}

type MockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func MockResponse(statusCode int, body string) *http.Response {
	resp := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
	// Add the header that the Elasticsearch client expects
	resp.Header.Set("X-Elastic-Product", "Elasticsearch")
	return resp
}

func setupMockWaveData(mockHandler func(req *http.Request) (*http.Response, error)) *persistence.WaveData {
	mockTransport := &MockTransport{
		RoundTripFunc: mockHandler,
	}

	mockClient, _ := elasticsearch.NewClient(elasticsearch.Config{
		Transport: mockTransport,
	})

	return persistence.NewWaveData(mockClient)
}

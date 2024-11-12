package persistence_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

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

func TestList_Success(t *testing.T) {
	waveDataStore := setupMockWaveData(func(req *http.Request) (*http.Response, error) {
		responseBody := `{
            "hits": {
                "hits": [
					{"_source": {
						"timestamp": "2024-09-17T09:00:00Z", 
						"h1_3": 0.6, 
						"hmax": 1.1, 
						"th1_3": 4.7, 
						"peak_direction": 8, 
						"peak_directional_spread": 32, 
						"temperature": 15
						}
					},
                    {"_source": {
						"timestamp": "2024-09-18T10:00:00Z", 
						"h1_3": 0.8, 
						"hmax": 1.3, 
						"th1_3": 5.0, 
						"peak_direction": 10, 
						"peak_directional_spread": 35, 
						"temperature": 16
						}
					}
                ]
            }
        }`
		return MockResponse(200, responseBody), nil
	})

	startDate := time.Date(2024, 9, 17, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 9, 19, 0, 0, 0, 0, time.UTC)

	waveDataList, err := waveDataStore.List(context.Background(), "test-index", startDate, endDate)
	assert.NoError(t, err)
	assert.Len(t, waveDataList, 2)

	expectedTimestamp1, err := time.Parse(time.RFC3339, "2024-09-17T09:00:00Z")
	assert.NoError(t, err)
	assert.Equal(t, expectedTimestamp1, waveDataList[0].Timestamp())
	assert.Equal(t, 0.6, waveDataList[0].AverageTopThirdWaveHeight())
	assert.Equal(t, 1.1, waveDataList[0].MaxHeight())
	assert.Equal(t, 4.7, waveDataList[0].AverageTopThirdWavePeriod())
	assert.Equal(t, 8, waveDataList[0].PeakDirection())
	assert.Equal(t, 32, waveDataList[0].PeakDirectionalSpread())
	assert.Equal(t, 15.0, waveDataList[0].Temperature())

	expectedTimestamp2, _ := time.Parse(time.RFC3339, "2024-09-18T10:00:00Z")
	assert.Equal(t, expectedTimestamp2, waveDataList[1].Timestamp())
	assert.Equal(t, 0.8, waveDataList[1].AverageTopThirdWaveHeight())
	assert.Equal(t, 1.3, waveDataList[1].MaxHeight())
	assert.Equal(t, 5.0, waveDataList[1].AverageTopThirdWavePeriod())
	assert.Equal(t, 10, waveDataList[1].PeakDirection())
	assert.Equal(t, 35, waveDataList[1].PeakDirectionalSpread())
	assert.Equal(t, 16.0, waveDataList[1].Temperature())
}

func TestList_EmptyIndexName(t *testing.T) {
	waveDataStore := setupMockWaveData(func(req *http.Request) (*http.Response, error) {
		return MockResponse(200, `{"hits": {"hits": []}}`), nil
	})

	startDate := time.Date(2024, 9, 17, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 9, 19, 0, 0, 0, 0, time.UTC)

	_, err := waveDataStore.List(context.Background(), "", startDate, endDate)
	assert.EqualError(t, err, "indexName cannot be empty")
}

func TestList_SearchError(t *testing.T) {
	waveDataStore := setupMockWaveData(func(req *http.Request) (*http.Response, error) {
		return MockResponse(500, `{"error": "internal server error"}`), nil
	})

	startDate := time.Date(2024, 9, 17, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 9, 19, 0, 0, 0, 0, time.UTC)

	_, err := waveDataStore.List(context.Background(), "test-index", startDate, endDate)
	assert.EqualError(t, err, "error executing search request: 500 Internal Server Error, body: {\"error\": \"internal server error\"}")
}

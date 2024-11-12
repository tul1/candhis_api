package persistence

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/tul1/candhis_api/internal/domain/model"
)

type WaveData struct {
	client *elasticsearch.Client
}

func NewWaveData(client *elasticsearch.Client) *WaveData {
	return &WaveData{
		client: client,
	}
}

func (w *WaveData) Add(ctx context.Context, waveData model.WaveData, indexName string) error {
	if indexName == "" {
		return errors.New("indexName cannot be empty")
	}

	dataJSON, err := json.Marshal(waveData)
	if err != nil {
		return fmt.Errorf("failed to marshal wave data to JSON: %v", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: waveData.Timestamp().Format("2006-01-02T15:04:05Z07:00"),
		Body:       bytes.NewReader(dataJSON),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, w.client)
	if err != nil {
		return fmt.Errorf("error indexing document: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("error indexing document: %s, body: %s", res.Status(), string(body))
	}

	return nil
}

func (w *WaveData) List(ctx context.Context, indexName string, startDate, endDate time.Time) ([]model.WaveData, error) {
	if indexName == "" {
		return nil, fmt.Errorf("indexName cannot be empty")
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"timestamp": map[string]interface{}{
					"gte": startDate.Format("2006-01-02T15:04:05Z07:00"),
					"lte": endDate.Format("2006-01-02T15:04:05Z07:00"),
				},
			},
		},
	}

	dataJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query to JSON: %v", err)
	}

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  bytes.NewReader(dataJSON),
	}

	res, err := req.Do(ctx, w.client)
	if err != nil {
		return nil, fmt.Errorf("error executing search request: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("error executing search request: %s, body: %s", res.Status(), string(body))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source model.WaveData `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	waveDataList := make([]model.WaveData, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		waveDataList[i] = hit.Source
	}

	return waveDataList, nil
}

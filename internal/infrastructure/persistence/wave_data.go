package persistence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

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
		return fmt.Errorf("indexName cannot be empty")
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

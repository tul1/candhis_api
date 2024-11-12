package persistencetest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/domain/model"
)

type waveDataPersistor struct {
	t        *testing.T
	esClient *elasticsearch.Client
}

func NewWaveDataPersistor(t *testing.T, esClient *elasticsearch.Client) *waveDataPersistor {
	t.Helper()

	return &waveDataPersistor{
		t:        t,
		esClient: esClient,
	}
}

func (p *waveDataPersistor) Add(ctx context.Context, waveData model.WaveData, index string) {
	p.t.Helper()

	data, err := json.Marshal(&waveData)
	require.NoError(p.t, err, "failed to marshal WaveData")

	req := esapi.IndexRequest{
		Index:   index,
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}
	res, err := req.Do(ctx, p.esClient)
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()
	require.NoError(p.t, err, "failed to clear wave_data index in Elasticsearch: %v", err)
	require.False(p.t, res.IsError())
}

func (p *waveDataPersistor) Get(ctx context.Context, waveDataID, index string) model.WaveData {
	p.t.Helper()

	req := esapi.GetRequest{
		Index:      index,
		DocumentID: waveDataID,
	}
	res, err := req.Do(ctx, p.esClient)
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()
	require.NoError(p.t, err, "failed to get wave_data from Elasticsearch")
	require.False(p.t, res.IsError(), "Elasticsearch returned an error during get request")

	body, err := io.ReadAll(res.Body)
	require.NoError(p.t, err, "failed to read response body from Elasticsearch")

	var esResponse struct {
		Source model.WaveData `json:"_source"`
	}

	err = json.Unmarshal(body, &esResponse)
	require.NoError(p.t, err, "failed to unmarshal wave_data from Elasticsearch")

	return esResponse.Source
}

func (p *waveDataPersistor) List(ctx context.Context, index string) []model.WaveData {
	p.t.Helper()

	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(`{"query": {"match_all": {}}}`),
	}

	res, err := req.Do(ctx, p.esClient)
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()
	require.NoError(p.t, err, "failed to search wave_data in Elasticsearch: %v", err)
	require.False(p.t, res.IsError(), "Elasticsearch returned an error during search")

	body, err := io.ReadAll(res.Body)
	require.NoError(p.t, err, "failed to read response body from Elasticsearch")

	var esResponse struct {
		Hits struct {
			Hits []struct {
				Source model.WaveData `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	err = json.Unmarshal(body, &esResponse)
	require.NoError(p.t, err, "failed to unmarshal search response from Elasticsearch")

	waveDataList := make([]model.WaveData, 0)
	for _, hit := range esResponse.Hits.Hits {
		waveDataList = append(waveDataList, hit.Source)
	}

	return waveDataList
}

func (p *waveDataPersistor) Clear(ctx context.Context) {
	p.t.Helper()

	req := esapi.DeleteByQueryRequest{
		Index: []string{"wave_data_test"},
		Body:  strings.NewReader(`{"query": {"match_all": {}}}`),
	}
	res, err := req.Do(ctx, p.esClient)
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()
	require.NoError(p.t, err, "failed to clear wave_data index in Elasticsearch: %v", err)
	require.False(p.t, res.IsError())
}

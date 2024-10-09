package persistencetest

import (
	"bytes"
	"context"
	"encoding/json"
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

func (p *waveDataPersistor) Add(ctx context.Context, waveData *model.WaveData, index string) {
	p.t.Helper()

	data, err := json.Marshal(waveData)
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
	require.True(p.t, res.IsError())
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
	require.NoError(p.t, err, "failed to clear wave_data index in Elasticsearch: %v", err)
	require.True(p.t, res.IsError())

	var esResponse struct {
		Source model.WaveData `json:"_source"`
	}
	err = json.NewDecoder(res.Body).Decode(&esResponse)
	require.NoError(p.t, err, "failed to decode wave_data from Elasticsearch: %v", err)

	return esResponse.Source
}

func (p *waveDataPersistor) ListWaveData(ctx context.Context, index string) []model.WaveData {
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

	var esResponse struct {
		Hits struct {
			Hits []struct {
				Source model.WaveData `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	err = json.NewDecoder(res.Body).Decode(&esResponse)
	require.NoError(p.t, err, "failed to decode search response from Elasticsearch: %v", err)

	waveDataList := make([]model.WaveData, 0)
	for _, hit := range esResponse.Hits.Hits {
		waveDataList = append(waveDataList, hit.Source)
	}

	return waveDataList
}

func (p *waveDataPersistor) Clear(ctx context.Context) {
	p.t.Helper()

	req := esapi.DeleteByQueryRequest{
		Index: []string{"_all"},
		Body:  strings.NewReader(`{"query": {"match_all": {}}}`),
	}
	res, err := req.Do(ctx, p.esClient)
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()
	require.NoError(p.t, err, "failed to clear wave_data index in Elasticsearch: %v", err)
	require.True(p.t, res.IsError())
}

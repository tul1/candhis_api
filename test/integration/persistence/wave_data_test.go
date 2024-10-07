package persistence_test

import (
	"context"
	"os"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/application/repository"
	"github.com/tul1/candhis_api/internal/domain/model/modeltest"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence/persistencetest"
)

func TestWaveData_Add_Success(t *testing.T) {
	ctx := context.Background()
	persistor, waveDataStore := setupWaveDataTest(t)

	waveData1 := modeltest.MustCreateWaveData(t, "17/09/2024", "09:00", "0.6", "1.1", "4.7", "8", "32", "14")
	waveData2 := modeltest.MustCreateWaveData(t, "18/09/2024", "10:00", "0.8", "1.3", "5.0", "10", "35", "14")

	err := waveDataStore.Add(ctx, waveData1, "wave_data_test")
	require.NoError(t, err)
	err = waveDataStore.Add(ctx, waveData2, "wave_data_test")
	require.NoError(t, err)

	retrievedWaveDataList := persistor.WaveData().List(context.Background(), "wave_data_test")

	require.Len(t, retrievedWaveDataList, 2)
	assert.Equal(t, retrievedWaveDataList[0], waveData1)
	assert.Equal(t, retrievedWaveDataList[1], waveData2)
}

func setupWaveDataTest(t *testing.T) (*persistencetest.ESPersistor, repository.WaveData) {
	t.Helper()

	esURL := os.Getenv("ELASTICSEARCH_URL")
	require.NotEmpty(t, esURL, "ELASTICSEARCH_URL must be set")

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{esURL},
	})
	require.NoError(t, err, "failed to create Elasticsearch client")

	waveData := persistence.NewWaveData(es)
	persistor := persistencetest.NewESPersistor(t, es)

	t.Cleanup(func() {
		persistor.Clear(context.Background())
	})

	return persistor, waveData
}

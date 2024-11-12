package persistence_test

import (
	"context"
	"os"
	"testing"
	"time"

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

func TestWaveData_List_Success(t *testing.T) {
	ctx := context.Background()
	persistor, waveDataStore := setupWaveDataTest(t)

	waveData1 := modeltest.MustCreateWaveData(t, "17/09/2024", "09:00", "0.6", "1.1", "4.7", "8", "32", "14")
	waveData2 := modeltest.MustCreateWaveData(t, "18/09/2024", "10:00", "0.8", "1.3", "5.0", "10", "35", "14")
	waveData3 := modeltest.MustCreateWaveData(t, "19/09/2024", "11:00", "1.0", "1.5", "5.5", "12", "40", "15")

	persistor.WaveData().Add(ctx, waveData1, "wave_data_test")
	persistor.WaveData().Add(ctx, waveData2, "wave_data_test")
	persistor.WaveData().Add(ctx, waveData3, "wave_data_test")

	startDate := waveData1.Timestamp()
	endDate := waveData3.Timestamp()
	retrievedWaveDataList, err := waveDataStore.List(ctx, "wave_data_test", startDate, endDate)
	require.NoError(t, err)

	require.Len(t, retrievedWaveDataList, 3)
	assert.Equal(t, waveData1, retrievedWaveDataList[0])
	assert.Equal(t, waveData2, retrievedWaveDataList[1])
	assert.Equal(t, waveData3, retrievedWaveDataList[2])
}

func TestWaveData_List_NoResults(t *testing.T) {
	ctx := context.Background()
	_, waveDataStore := setupWaveDataTest(t)

	startDate := time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 9, 16, 0, 0, 0, 0, time.UTC)

	retrievedWaveDataList, err := waveDataStore.List(ctx, "wave_data_test", startDate, endDate)
	require.NoError(t, err)
	require.Len(t, retrievedWaveDataList, 0)
}

func TestWaveData_List_PartialResults(t *testing.T) {
	ctx := context.Background()
	persistor, waveDataStore := setupWaveDataTest(t)

	waveData1 := modeltest.MustCreateWaveData(t, "17/09/2024", "09:00", "0.6", "1.1", "4.7", "8", "32", "14")
	waveData2 := modeltest.MustCreateWaveData(t, "18/09/2024", "10:00", "0.8", "1.3", "5.0", "10", "35", "14")
	waveData3 := modeltest.MustCreateWaveData(t, "19/09/2024", "11:00", "1.0", "1.5", "5.5", "12", "40", "15")

	persistor.WaveData().Add(ctx, waveData1, "wave_data_test")
	persistor.WaveData().Add(ctx, waveData2, "wave_data_test")
	persistor.WaveData().Add(ctx, waveData3, "wave_data_test")

	startDate := waveData1.Timestamp()
	endDate := waveData2.Timestamp()
	retrievedWaveDataList, err := waveDataStore.List(ctx, "wave_data_test", startDate, endDate)
	require.NoError(t, err)

	require.Len(t, retrievedWaveDataList, 2)
	assert.Equal(t, waveData1, retrievedWaveDataList[0])
	assert.Equal(t, waveData2, retrievedWaveDataList[1])
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

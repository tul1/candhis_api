package modeltest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/domain/model"
)

func MustCreateWaveData(t *testing.T, dateStr, timeStr, h13, hmax, th13, direction, spread, temp string) model.WaveData {
	t.Helper()

	waveData, err := model.NewWaveData(dateStr, timeStr, h13, hmax, th13, direction, spread, temp)
	require.NoError(t, err, "failed to create WaveData")

	return waveData
}

package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tul1/candhis_api/internal/domain/model"
)

func TestNewWaveDataSuccess(t *testing.T) {
	waveData, err := model.NewWaveData("07/10/2024", "14:00", 2.5, 4.0, 10.5, 90, 30, 20.0)
	require.NoError(t, err)

	expectedTimestamp, err := time.Parse(time.RFC3339, "2024-10-07T14:00:00Z")
	require.NoError(t, err)

	assert.Equal(t, expectedTimestamp, waveData.Timestamp())
	assert.Equal(t, 2.5, waveData.AverageTopThirdWaveHeight())
	assert.Equal(t, 4.0, waveData.MaxHeight())
	assert.Equal(t, 10.5, waveData.AverageTopThirdWavePeriod())
	assert.Equal(t, 90, waveData.PeakDirection())
	assert.Equal(t, 30, waveData.PeakDirectionalSpread())
	assert.Equal(t, 20.0, waveData.Temperature())
}

func TestNewWaveDataFailure(t *testing.T) {
	testCases := map[string]struct {
		date                      string
		timeStr                   string
		averageTopThirdWaveHeight float64
		maxHeight                 float64
		averageTopThirdWavePeriod float64
		peakDirection             int
		peakDirectionalSpread     int
		temperature               float64
		errMsg                    string
	}{
		"invalid date format YYYY/MM/DD": {
			date:                      "2024/10/07", // Invalid format
			timeStr:                   "14:00",
			averageTopThirdWaveHeight: 2.5,
			maxHeight:                 4.0,
			averageTopThirdWavePeriod: 10.5,
			peakDirection:             90,
			peakDirectionalSpread:     30,
			temperature:               20.0,
			errMsg:                    "invalid date or time format, expected DD/MM/YYYY and HH:MM",
		},
		"invalid date format DD-MM-YYYY": {
			date:                      "07-10-2024", // Invalid format
			timeStr:                   "14:00",
			averageTopThirdWaveHeight: 2.5,
			maxHeight:                 4.0,
			averageTopThirdWavePeriod: 10.5,
			peakDirection:             90,
			peakDirectionalSpread:     30,
			temperature:               20.0,
			errMsg:                    "invalid date or time format, expected DD/MM/YYYY and HH:MM",
		},
		"invalid time format": {
			date:                      "07/10/2024",
			timeStr:                   "14:00:00", // Invalid time format
			averageTopThirdWaveHeight: 2.5,
			maxHeight:                 4.0,
			averageTopThirdWavePeriod: 10.5,
			peakDirection:             90,
			peakDirectionalSpread:     30,
			temperature:               20.0,
			errMsg:                    "invalid date or time format, expected DD/MM/YYYY and HH:MM",
		},
		"negative wave height": {
			date:                      "07/10/2024",
			timeStr:                   "14:00",
			averageTopThirdWaveHeight: -1.0, // Invalid: negative value
			maxHeight:                 4.0,
			averageTopThirdWavePeriod: 10.5,
			peakDirection:             90,
			peakDirectionalSpread:     30,
			temperature:               20.0,
			errMsg:                    "invalid input: negative values for heights, periods, or temperature below absolute zero",
		},
		"negative max wave height": {
			date:                      "07/10/2024",
			timeStr:                   "14:00",
			averageTopThirdWaveHeight: 2.5,
			maxHeight:                 -4.0, // Invalid: negative value
			averageTopThirdWavePeriod: 10.5,
			peakDirection:             90,
			peakDirectionalSpread:     30,
			temperature:               20.0,
			errMsg:                    "invalid input: negative values for heights, periods, or temperature below absolute zero",
		},
		"negative wave period": {
			date:                      "07/10/2024",
			timeStr:                   "14:00",
			averageTopThirdWaveHeight: 2.5,
			maxHeight:                 4.0,
			averageTopThirdWavePeriod: -10.5, // Invalid: negative value
			peakDirection:             90,
			peakDirectionalSpread:     30,
			temperature:               20.0,
			errMsg:                    "invalid input: negative values for heights, periods, or temperature below absolute zero",
		},
		"temperature below absolute zero": {
			date:                      "07/10/2024",
			timeStr:                   "14:00",
			averageTopThirdWaveHeight: 2.5,
			maxHeight:                 4.0,
			averageTopThirdWavePeriod: 10.5,
			peakDirection:             90,
			peakDirectionalSpread:     30,
			temperature:               -300.0, // Invalid: temperature below absolute zero
			errMsg:                    "invalid input: negative values for heights, periods, or temperature below absolute zero",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			waveData, err := model.NewWaveData(
				tc.date, tc.timeStr, tc.averageTopThirdWaveHeight, tc.maxHeight,
				tc.averageTopThirdWavePeriod, tc.peakDirection, tc.peakDirectionalSpread, tc.temperature,
			)
			assert.EqualError(t, err, tc.errMsg)
			assert.Equal(t, model.WaveData{}, waveData)
		})
	}
}

func TestWaveDataMarshalJSON(t *testing.T) {
	waveData, err := model.NewWaveData(
		"07/10/2024", "14:00", 2.5, 4.0, 10.5, 90, 30, 20.0,
	)
	require.NoError(t, err)

	jsonData, err := json.Marshal(waveData)
	require.NoError(t, err)

	expectedJSON := `{
		"timestamp": "2024-10-07T14:00:00Z",
		"h1_3": 2.5,
		"hmax": 4.0,
		"th1_3": 10.5,
		"peak_direction": 90,
		"peak_directional_spread": 30,
		"temperature": 20.0
	}`
	assert.JSONEq(t, expectedJSON, string(jsonData))
}

func TestWaveDataUnmarshalJSON(t *testing.T) {
	jsonData := `{
		"timestamp": "2024-10-07T14:00:00Z",
		"h1_3": 2.5,
		"hmax": 4.0,
		"th1_3": 10.5,
		"peak_direction": 90,
		"peak_directional_spread": 30,
		"temperature": 20.0
	}`

	var waveData model.WaveData
	err := json.Unmarshal([]byte(jsonData), &waveData)
	require.NoError(t, err)

	expectedTimestamp, err := time.Parse(time.RFC3339, "2024-10-07T14:00:00Z")
	require.NoError(t, err)

	assert.Equal(t, expectedTimestamp, waveData.Timestamp())
	assert.Equal(t, 2.5, waveData.AverageTopThirdWaveHeight())
	assert.Equal(t, 4.0, waveData.MaxHeight())
	assert.Equal(t, 10.5, waveData.AverageTopThirdWavePeriod())
	assert.Equal(t, 90, waveData.PeakDirection())
	assert.Equal(t, 30, waveData.PeakDirectionalSpread())
	assert.Equal(t, 20.0, waveData.Temperature())
}

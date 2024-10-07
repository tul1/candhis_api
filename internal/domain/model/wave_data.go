package model

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type WaveData struct {
	// timestamp of the observation.
	timestamp time.Time
	// Significant wave height, the average value of the highest one-third of wave heights observed over a 30-minute period.
	averageTopThirdWaveHeight float64
	// Height of the largest wave observed over a 30-minute period.
	maxHeight float64
	// Significant period, defined by the average value of the periods of the highest one-third of the largest waves observed
	// over a 30-minute period.
	averageTopThirdWavePeriod float64
	// Average direction of wave origin at the peak of the energy spectrum. The angle is measured positively, in a clockwise
	// direction, between geographic north and the direction of the wave origin.
	peakDirection int
	// Directional width, characterizing the directional spread of energy around the average direction at the peak (angular
	// distribution function of energy associated with the peak frequency of the energy spectrum).
	peakDirectionalSpread int
	// Water temperature in degrees Celsius at the time of the observation.
	temperature float64
}

func NewWaveData(
	dateStr,
	timeStr,
	averageTopThirdWaveHeightStr,
	maxHeightStr,
	averageTopThirdWavePeriodStr,
	peakDirectionStr,
	peakDirectionalSpreadStr,
	temperatureStr string,
) (WaveData, error) {
	datetimeStr := dateStr + " " + timeStr
	timestamp, err := time.Parse("02/01/2006 15:04", datetimeStr)
	if err != nil {
		return WaveData{}, errors.New("invalid date or time format, expected DD/MM/YYYY and HH:MM")
	}

	averageTopThirdWaveHeight, err := strconv.ParseFloat(averageTopThirdWaveHeightStr, 64)
	if err != nil {
		return WaveData{}, errors.New("invalid value for averageTopThirdWaveHeight")
	}

	maxHeight, err := strconv.ParseFloat(maxHeightStr, 64)
	if err != nil {
		return WaveData{}, errors.New("invalid value for maxHeight")
	}

	averageTopThirdWavePeriod, err := strconv.ParseFloat(averageTopThirdWavePeriodStr, 64)
	if err != nil {
		return WaveData{}, errors.New("invalid value for averageTopThirdWavePeriod")
	}

	peakDirection, err := strconv.Atoi(peakDirectionStr)
	if err != nil {
		return WaveData{}, errors.New("invalid value for peakDirection")
	}

	peakDirectionalSpread, err := strconv.Atoi(peakDirectionalSpreadStr)
	if err != nil {
		return WaveData{}, errors.New("invalid value for peakDirectionalSpread")
	}

	temperature, err := strconv.ParseFloat(temperatureStr, 64)
	if err != nil {
		return WaveData{}, errors.New("invalid value for temperature")
	}

	if averageTopThirdWaveHeight < 0 || maxHeight < 0 || averageTopThirdWavePeriod < 0 || temperature < -273.15 {
		return WaveData{}, errors.New("invalid input: negative values for heights, periods, or temperature below absolute zero")
	}

	return WaveData{
		timestamp,
		averageTopThirdWaveHeight,
		maxHeight,
		averageTopThirdWavePeriod,
		peakDirection,
		peakDirectionalSpread,
		temperature,
	}, nil
}

func (w WaveData) Timestamp() time.Time {
	return w.timestamp
}

func (w WaveData) AverageTopThirdWaveHeight() float64 {
	return w.averageTopThirdWaveHeight
}

func (w WaveData) MaxHeight() float64 {
	return w.maxHeight
}

func (w WaveData) AverageTopThirdWavePeriod() float64 {
	return w.averageTopThirdWavePeriod
}

func (w WaveData) PeakDirection() int {
	return w.peakDirection
}

func (w WaveData) PeakDirectionalSpread() int {
	return w.peakDirectionalSpread
}

func (w WaveData) Temperature() float64 {
	return w.temperature
}

type waveDataJSON struct {
	Timestamp                 string  `json:"timestamp"`
	AverageTopThirdWaveHeight float64 `json:"h1_3"`
	MaxHeight                 float64 `json:"hmax"`
	AverageTopThirdWavePeriod float64 `json:"th1_3"`
	PeakDirection             int     `json:"peak_direction"`
	PeakDirectionalSpread     int     `json:"peak_directional_spread"`
	Temperature               float64 `json:"temperature"`
}

func (w WaveData) MarshalJSON() ([]byte, error) {
	data := waveDataJSON{
		Timestamp:                 w.timestamp.Format(time.RFC3339),
		AverageTopThirdWaveHeight: w.averageTopThirdWaveHeight,
		MaxHeight:                 w.maxHeight,
		AverageTopThirdWavePeriod: w.averageTopThirdWavePeriod,
		PeakDirection:             w.peakDirection,
		PeakDirectionalSpread:     w.peakDirectionalSpread,
		Temperature:               w.temperature,
	}
	return json.Marshal(data)
}

func (w *WaveData) UnmarshalJSON(data []byte) error {
	var aux waveDataJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	timestamp, err := time.Parse(time.RFC3339, aux.Timestamp)
	if err != nil {
		return err
	}

	*w = WaveData{
		timestamp:                 timestamp,
		averageTopThirdWaveHeight: aux.AverageTopThirdWaveHeight,
		maxHeight:                 aux.MaxHeight,
		averageTopThirdWavePeriod: aux.AverageTopThirdWavePeriod,
		peakDirection:             aux.PeakDirection,
		peakDirectionalSpread:     aux.PeakDirectionalSpread,
		temperature:               aux.Temperature,
	}

	return nil
}

package candhisapi

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tul1/candhis_api/openapi"
)

const elasticSearchIndexLesPierresNoires = "les-pierres-noires"

func (s candhisAPI) Waves(c *gin.Context) {
	waves, err := s.waveData.List(
		c.Request.Context(),
		elasticSearchIndexLesPierresNoires,
		time.Now().Add(-24*time.Hour),
		time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.Error{Message: err.Error()})
	}

	wavesResp := make(openapi.Waves, len(waves))
	for i, wave := range waves {
		wavesResp[i] = struct {
			AverageTopThirdWaveHeight float32 `json:"averageTopThirdWaveHeight"`
			AverageTopThirdWavePeriod float32 `json:"averageTopThirdWavePeriod"`
			MaxHeight                 float32 `json:"maxHeight"`
			PeakDirection             float32 `json:"peakDirection"`
			PeakDirectionalSpread     int     `json:"peakDirectionalSpread"`
			Temperature               float32 `json:"temperature"`
			Timestamp                 string  `json:"timestamp"`
		}{
			AverageTopThirdWaveHeight: float32(wave.AverageTopThirdWaveHeight()),
			AverageTopThirdWavePeriod: float32(wave.AverageTopThirdWavePeriod()),
			MaxHeight:                 float32(wave.MaxHeight()),
			PeakDirection:             float32(wave.PeakDirection()),
			PeakDirectionalSpread:     wave.PeakDirectionalSpread(),
			Temperature:               float32(wave.Temperature()),
			Timestamp:                 wave.Timestamp().Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, wavesResp)
}

package repository

import (
	"context"
	"time"

	"github.com/tul1/candhis_api/internal/domain/model"
)

//go:generate mockgen -package persistencemock -destination=./persistence_mock/wave_data.go -source=wave_data.go WaveData
type WaveData interface {
	Add(ctx context.Context, waveData model.WaveData, indexName string) error
	List(ctx context.Context, indexName string, startDate, endDate time.Time) ([]model.WaveData, error)
}

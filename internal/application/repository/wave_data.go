package repository

import (
	"context"

	"github.com/tul1/candhis_api/internal/domain/model"
)

//go:generate mockgen -package persistencemock -destination=./persistence_mock/wave_data.go -source=wave_data.go WaveData
type WaveData interface {
	Add(ctx context.Context, waveData model.WaveData, indexName string) error
}

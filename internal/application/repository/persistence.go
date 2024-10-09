package repository

import (
	"context"

	appmodel "github.com/tul1/candhis_api/internal/application/model"
	domainmodel "github.com/tul1/candhis_api/internal/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -package persistencemock -destination=./persistence_mock/sessionid.go -source=sessionid.go SessionID
type SessionID interface {
	Get(ctx context.Context) (*appmodel.CandhisSessionID, error)
	Update(ctx context.Context, sessionID *appmodel.CandhisSessionID) error
}

type WaveData interface {
	AddWaveData(ctx context.Context, waveData domainmodel.WaveData, indexName string) error
}

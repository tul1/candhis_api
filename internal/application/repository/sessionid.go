package repository

import (
	"context"

	appmodel "github.com/tul1/candhis_api/internal/application/model"
)

//go:generate mockgen -package persistencemock -destination=./persistence_mock/sessionid.go -source=sessionid.go SessionID
type SessionID interface {
	Get(ctx context.Context) (*appmodel.CandhisSessionID, error)
	Update(ctx context.Context, sessionID appmodel.CandhisSessionID) error
}

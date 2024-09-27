package repository

import (
	"context"

	"github.com/tul1/candhis_api/internal/application/model"
)

//go:generate mockgen -package repositorymocks -destination=./repositorymocks/sessionid.go . SessionID
type SessionID interface {
	Get(ctx context.Context) (*model.CandhisSessionID, error)
	Update(ctx context.Context, sessionID *model.CandhisSessionID) error
}

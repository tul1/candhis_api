package repository

import (
	"context"

	"github.com/tul1/candhis_api/internal/application/model"
)

//go:generate mockgen -package clientmock -destination=./client_mock/candhis_sessionid_scraper.go -source=candhis_sessionid_scraper.go CandhisSessionIDWebScraper
type CandhisSessionIDWebScraper interface {
	GetCandhisSessionID(ctx context.Context) (model.CandhisSessionID, error)
}

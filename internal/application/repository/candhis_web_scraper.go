package repository

import (
	appmodel "github.com/tul1/candhis_api/internal/application/model"
	"github.com/tul1/candhis_api/internal/domain/model"
)

//go:generate mockgen -package clientmock -destination=./client_mock/candhis_web_scraper.go -source=candhis_web_scraper.go CandhisCampaignsWebScraper
type CandhisCampaignsWebScraper interface {
	GatherWavesDataFromWebTable(candhisSessionID appmodel.CandhisSessionID, candhisURL string) ([]model.WaveData, error)
}

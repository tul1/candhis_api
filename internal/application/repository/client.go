package repository

import (
	appmodel "github.com/tul1/candhis_api/internal/application/model"
	"github.com/tul1/candhis_api/internal/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -package clientmock -destination=./client_mock/scraping_bee_client.go -source=client.go ScrapingBeeClient
type ScrapingBeeClient interface {
	GetCandhisSessionID() (appmodel.CandhisSessionID, error)
}

//go:generate go run go.uber.org/mock/mockgen -package clientmock -destination=./client_mock/candhis_web_scrapper.go -source=client.go CandhisWebScrapper
type CandhisWebScrapper interface {
	GatherWaveDataFromWebTable(phpsessid, candhisURL string) ([]model.WaveData, error)
}

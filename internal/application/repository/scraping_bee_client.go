package repository

import (
	"github.com/tul1/candhis_api/internal/application/model"
)

//go:generate mockgen -package clientmock -destination=./client_mock/scraping_bee_client.go -source=scraping_bee_client.go ScrapingBeeClient
type ScrapingBeeClient interface {
	GetCandhisSessionID() (model.CandhisSessionID, error)
}

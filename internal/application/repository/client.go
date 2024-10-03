package repository

import "github.com/tul1/candhis_api/internal/application/model"

//go:generate go run go.uber.org/mock/mockgen -package clientmock -destination=./client_mock/scraping_bee_client.go -source=client.go ScrapingBeeClient
type ScrapingBeeClient interface {
	GetCandhisSessionID() (model.CandhisSessionID, error)
}

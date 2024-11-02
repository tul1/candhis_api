package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/tul1/candhis_api/internal/application/service"
	"github.com/tul1/candhis_api/internal/infrastructure/client"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence"
	"github.com/tul1/candhis_api/internal/pkg/configuration"
	"github.com/tul1/candhis_api/internal/pkg/db"
	"github.com/tul1/candhis_api/internal/pkg/logger"
)

func main() {
	log := logger.NewWithDefaultLogger()
	ctx := context.Background()

	// Parse the config file path from the command line arguments
	configFile := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	// Load configuration
	config, err := configuration.Load[Config](*configFile)
	if err != nil {
		log.Errorf("Configuration error: %v", err)
		return
	}

	// Create cConnect to the PostgreSQL database
	dbConn, err := db.NewDBConnection(
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName, db.DefaultDBConnector, log)
	if err != nil {
		log.Errorf("Database connection error: %v", err)
		return
	}
	defer dbConn.CloseWithLog()

	// Create candhis scraper service
	httpClient := http.Client{}
	defer httpClient.CloseIdleConnections()

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{config.ElasticsearchURL}})
	if err != nil {
		log.Errorf("Failed to create Elasticsearch client: %v", err)
		return
	}

	candhisCampaignsScraper := service.NewCandhisCampaignsScraper(
		persistence.NewSessionID(dbConn.DB),
		persistence.NewWaveData(esClient),
		client.NewCandhisCampaignsWebScraper(&httpClient),
	)

	// Scraping and store campaigns from Candhis web
	log.Info("Start scraping Candhis web to fetch and store wave data from campaigns")
	err = candhisCampaignsScraper.FetchAndStoreWaveData(ctx)
	if err != nil {
		log.Errorf("Failed scraping Candhis web to fetch and store wave data from campaigns: %v", err)
		return
	}
	log.Info("Finished scraping Candhis web to fetch and store wave data from campaigns Successfully")
}

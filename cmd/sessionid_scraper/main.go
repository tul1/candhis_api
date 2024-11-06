package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/tul1/candhis_api/internal/application/service"
	"github.com/tul1/candhis_api/internal/infrastructure/client"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence"
	"github.com/tul1/candhis_api/internal/pkg/chrome"
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

	// Create connect to the PostgreSQL database
	dbConn, err := db.NewDBConnection(
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName, db.DefaultDBConnector, log)
	if err != nil {
		log.Errorf("Database connection error: %v", err)
		return
	}
	defer dbConn.CloseWithLog()

	// Get Get Chrome ID from headless-chrome service
	httpClient := http.Client{}
	defer httpClient.CloseIdleConnections()

	chromeScraper, err := chrome.NewChromedpScraper(&httpClient, config.ChromeURL)
	if err != nil {
		log.Errorf("Chrome scraper initialization error: %v", err)
		return
	}

	// Create candhisScraper service
	candhisScraper := service.NewCandhisSessionIDScraper(
		persistence.NewSessionID(dbConn.DB),
		client.NewCandhisSessionIDWebScraper(chromeScraper, config.TargetWeb),
	)

	// Retrieve and Store CandhisSessionID
	log.Info("Start scraping Candhis web to fetch and store session id")
	if err = candhisScraper.FetchAndStoreSessionID(ctx); err != nil {
		log.Errorf("Failed scraping Candhis web to fetch and store session id: %v", err)
		return
	}
	log.Info("Finished scraping Candhis web to fetch and store session id successfully")
}

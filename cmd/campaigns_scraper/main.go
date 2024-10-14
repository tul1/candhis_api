package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/tul1/candhis_api/internal/application/service"
	"github.com/tul1/candhis_api/internal/infrastructure/client"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence"
	"github.com/tul1/candhis_api/internal/pkg/db"
	"github.com/tul1/candhis_api/internal/pkg/loadconfig"
	"github.com/tul1/candhis_api/internal/pkg/logger"
)

func loadConfig(configFile string) (*Config, error) {
	var config *Config
	if err := loadconfig.LoadConfig(configFile, config); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return config, nil
}

func createDBConnection(config Config) (*sql.DB, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	dbConn, err := db.NewDatabaseConnection(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	return dbConn, nil
}

func main() {
	log := logger.NewWithDefaultLogger()
	ctx := context.Background()

	// Parse the config file path from the command line arguments
	configFile := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Errorf("Configuration error: %v", err)
		return
	}

	// Create cConnect to the PostgreSQL database
	dbConn, err := createDBConnection(*config)
	defer func() {
		if err = dbConn.Close(); err != nil {
			log.Errorf("Failed closing database: %v", err)
			return
		}
	}()
	if err != nil {
		log.Errorf("Database connection error: %v", err)
		return
	}

	// Create candhis scraper service
	httpClient := http.Client{}
	defer httpClient.CloseIdleConnections()

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{config.ElasticsearchURL}})
	if err != nil {
		log.Errorf("Failed to create Elasticsearch client: %v", err)
		return
	}

	candhisCampaignsScraper := service.NewCandhisCampaignsScraper(
		persistence.NewSessionID(dbConn),
		persistence.NewWaveData(esClient),
		client.NewCandhisCampaignsWebScraper(&httpClient),
	)

	// Scraping and store campaigns from Candhis web
	log.Info("Start scraping Candhis web campaigns")
	err = candhisCampaignsScraper.ScrapingCandhisCampaigns(ctx)
	if err != nil {
		log.Errorf("Failed while scraping Candhis web campaigns: %v", err)
		return
	}
	log.Info("Finished scraping Candhis web campaigns Successfully")
}

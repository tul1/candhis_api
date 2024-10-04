package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"

	"github.com/tul1/candhis_api/internal/application/service"
	"github.com/tul1/candhis_api/internal/application/service/client"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence"
	"github.com/tul1/candhis_api/internal/pkg/db"
	"github.com/tul1/candhis_api/internal/pkg/loadconfig"
	"github.com/tul1/candhis_api/internal/pkg/logger"
)

func loadConfig(configFile string) (*Config, error) {
	var config *Config
	err := loadconfig.LoadConfig(configFile, config)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return config, nil
}

func createDBConnection(config Config) (*sql.DB, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	db, err := db.NewDatabaseConnection(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	return db, nil
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
	db, err := createDBConnection(*config)
	defer db.Close()
	if err != nil {
		log.Errorf("Database connection error: %v", err)
		return
	}

	// Create candhisScraper
	httpClient := http.Client{}
	defer httpClient.CloseIdleConnections()

	candhisScraper := service.NewCandhisScraper(
		persistence.NewSessionIDRepository(db),
		client.NewScrapingBeeClient(&httpClient, config.ScrapingbeeAPIKey),
	)

	// Retrieve and Store CandhisSessionID
	log.Info("Start to retrieve from Candhis web and store candhisSessionID to DB")
	err = candhisScraper.RetrieveAndStoreCandhisSessionID(ctx)
	if err != nil {
		log.Errorf("Failed to retrieve candhis session ID: %v", err)
		return
	}
	log.Info("Finished retrieving from Candhis web and storing candhisSessionID to DB: Successfully")
}

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/tul1/candhis_api/internal/application/service/client"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence"
	"github.com/tul1/candhis_api/internal/pkg/db"
	"github.com/tul1/candhis_api/internal/pkg/loadconfig"
)

const (
	candhisURL                         = "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ=="
	elasticSearchIndexLesPierresNoires = "les-pierres-noires"
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

	dbConn, err := db.NewDatabaseConnection(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	return dbConn, nil
}

func main() {
	ctx := context.Background()

	// Parse the config file path from the command line arguments
	configFile := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configFile)
	if err != nil {
		// log.Errorf("Configuration error: %v", err)
		return
	}

	// Create cConnect to the PostgreSQL database
	dbConn, err := createDBConnection(*config)
	defer func() {
		err = dbConn.Close()
		if err != nil {
			// log.Errorf("Failed closing database: %v", err)
			return
		}
	}()
	if err != nil {
		// log.Errorf("Database connection error: %v", err)
		return
	}
	// Retrieve the latest session ID from the database
	sessionIDRepo := persistence.NewSessionIDStore(dbConn)
	sessionID, err := sessionIDRepo.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to retrieve session ID: %v", err)
	}

	// Create an HTTP client and use the session ID to scrap data
	candhisWebScrapperClient := client.NewCandhisWebScrapper(&http.Client{})
	waveDataList, err := candhisWebScrapperClient.GatherWaveDataFromWebTable(sessionID.ID(), candhisURL)
	if err != nil {
		log.Fatalf("Failed to get table data: %v", err)
	}

	// Initialize the Elasticsearch client with custom URL and port
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.ElasticsearchURL},
	})
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	// Initialize the WaveDataStore with the Elasticsearch client
	waveDataStore := persistence.NewWaveData(esClient)

	// Iterate over waveDataList and index each wave data element individually
	for _, waveData := range waveDataList {
		err := waveDataStore.AddWaveData(ctx, waveData, elasticSearchIndexLesPierresNoires)
		if err != nil {
			log.Printf("Failed to push wave data to Elasticsearch: %v", err)
		}
	}
}

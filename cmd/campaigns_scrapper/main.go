package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/tul1/candhis_api/internal/domain/model"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence"
	"github.com/tul1/candhis_api/internal/pkg/db"
	"github.com/tul1/candhis_api/internal/pkg/loadconfig"
)

const (
	candhisURL                         = "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ=="
	elasticSearchIndexLesPierresNoires = "les-pierres-noires"
)

func parseRow(cells *goquery.Selection) (model.WaveData, error) {
	if cells.Length() != 8 {
		return model.WaveData{}, fmt.Errorf("expected 8 cells, but got %d", cells.Length())
	}

	values := make([]string, 8)
	cells.Each(func(cellIndex int, cell *goquery.Selection) {
		values[cellIndex] = strings.TrimSpace(cell.Text())
	})

	return model.NewWaveData(values[0], values[1], values[2], values[3],
		values[4], values[5], values[6], values[7])
}

func getTableData(client *http.Client, phpsessid string) ([]model.WaveData, error) {
	var waveDataList []model.WaveData

	reqURL := candhisURL
	req, err := http.NewRequest(http.MethodGet, candhisURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request, url: %s, error: %w", reqURL, err)
	}

	req.Header.Set("Accept", "text/html")
	req.Header.Set("Cookie", fmt.Sprintf("acceptCookies=true; %s", phpsessid))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request, url: %s, error: %w", reqURL, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed parse HTML: %w", err)
	}

	doc.Find("table.table-striped.table-bordered.table-sm").Each(func(index int, table *goquery.Selection) {
		table.Find("tr").Each(func(rowIndex int, row *goquery.Selection) {
			waveData, err := parseRow(row.Find("td"))
			if err != nil {
				log.Printf("Skipping row due to error: %v", err)
				return
			}

			waveDataList = append(waveDataList, waveData)
		})
	})

	return waveDataList, nil
}

func pushWaveDataToES(waveDataList []model.WaveData) error {
	// Initialize Elasticsearch client
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return fmt.Errorf("error creating Elasticsearch client: %v", err)
	}

	// Loop over the waveDataList and insert each document into Elasticsearch
	for _, waveData := range waveDataList {
		dataJSON, err := json.Marshal(waveData)
		if err != nil {
			log.Printf("Failed to marshal wave data to JSON: %v", err)
			continue
		}

		// Use the RFC3339 formatted timestamp as the DocumentID
		documentID := waveData.Timestamp().Format(time.RFC3339)

		// Prepare the request to index the data
		req := esapi.IndexRequest{
			Index:      elasticSearchIndexLesPierresNoires,
			DocumentID: documentID,
			Body:       bytes.NewReader(dataJSON),
			Refresh:    "true",
		}

		// Perform the request
		res, err := req.Do(context.Background(), es)
		if err != nil {
			log.Printf("Error indexing document: %v", err)
			continue
		}
		defer res.Body.Close()

		// Check if the request was successful
		if res.IsError() {
			body, _ := io.ReadAll(res.Body)
			log.Printf("Error indexing document: %s\nResponse body: %s", res.Status(), string(body))
		} else {
			log.Printf("Successfully indexed document: %s", documentID)
		}
	}

	return nil
}

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
	sessionIDRepo := persistence.NewSessionIDRepository(dbConn)
	sessionID, err := sessionIDRepo.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to retrieve session ID: %v", err)
	}

	// Create an HTTP client and use the session ID to scrape data
	client := &http.Client{}
	waveDataList, err := getTableData(client, sessionID.ID())
	if err != nil {
		log.Fatalf("Failed to get table data: %v", err)
	}

	// Push the data to Elasticsearch
	err = pushWaveDataToES(waveDataList)
	if err != nil {
		log.Fatalf("Failed to push data to Elasticsearch: %v", err)
	}
}

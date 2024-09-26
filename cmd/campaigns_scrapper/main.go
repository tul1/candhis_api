package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/jackc/pgx/v4"
	"github.com/tul1/candhis_api/internal/pkg/loadconfig"
)

const (
	candhisURL                         = "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ=="
	elasticSearchIndexLesPierresNoires = "les-pierres-noires"
)

type WaveData struct {
	Date        string  `json:"date"`
	Time        string  `json:"time"`
	H1_3        float64 `json:"h1_3"`
	Hmax        float64 `json:"hmax"`
	Th1_3       float64 `json:"th1_3"`
	DirPeak     int     `json:"dir_peak"`
	EtalPic     int     `json:"etal_pic"`
	Temperature float64 `json:"temperature"`
}

func validateWaveData(waveData WaveData) error {
	if waveData.Date == "" || waveData.Time == "" || waveData.H1_3 == 0 || waveData.Hmax == 0 || waveData.Th1_3 == 0 || waveData.DirPeak == 0 || waveData.EtalPic == 0 || waveData.Temperature == 0 {
		return errors.New("one or more required fields are missing or have invalid values")
	}
	return nil
}

func parseRow(cells *goquery.Selection) (WaveData, error) {
	var waveData WaveData
	cellCount := cells.Length()

	if cellCount != 8 {
		return waveData, fmt.Errorf("expected 8 cells, but got %d", cellCount)
	}

	fieldMap := map[int]func(string) error{
		0: func(val string) error { waveData.Date = val; return nil },
		1: func(val string) error { waveData.Time = val; return nil },
		2: func(val string) error {
			v, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			waveData.H1_3 = v
			return nil
		},
		3: func(val string) error {
			v, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			waveData.Hmax = v
			return nil
		},
		4: func(val string) error {
			v, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			waveData.Th1_3 = v
			return nil
		},
		5: func(val string) error {
			v, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			waveData.DirPeak = v
			return nil
		},
		6: func(val string) error {
			v, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			waveData.EtalPic = v
			return nil
		},
		7: func(val string) error {
			v, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			waveData.Temperature = v
			return nil
		},
	}

	cells.Each(func(cellIndex int, cell *goquery.Selection) {
		cellText := strings.TrimSpace(cell.Text())

		if parseFunc, ok := fieldMap[cellIndex]; ok {
			if err := parseFunc(cellText); err != nil {
				log.Printf("Failed to parse value for cell %d: %v", cellIndex, err)
			}
		}
	})

	// Validate the parsed data
	if err := validateWaveData(waveData); err != nil {
		return waveData, err
	}

	return waveData, nil
}

func getTableData(client *http.Client, phpsessid string) ([]WaveData, error) {
	var waveDataList []WaveData

	reqURL := candhisURL
	req, err := http.NewRequest(http.MethodGet, candhisURL, nil)
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

	// Parse the table data
	doc.Find("table.table-striped.table-bordered.table-sm").Each(func(index int, table *goquery.Selection) {
		table.Find("tr").Each(func(rowIndex int, row *goquery.Selection) {
			// Parse each row of cells into WaveData
			waveData, err := parseRow(row.Find("td"))
			if err != nil {
				log.Printf("Skipping row due to error: %v", err)
				return
			}

			// Add the valid wave data to the list
			waveDataList = append(waveDataList, waveData)
		})
	})

	return waveDataList, nil
}

func getSessionIDFromDB(ctx context.Context, conn *pgx.Conn) (string, error) {
	var sessionID string
	query := `SELECT id FROM candhis_session`

	err := conn.QueryRow(ctx, query).Scan(&sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve session ID from DB: %w", err)
	}

	return sessionID, nil
}

func sanitizeDocumentID(date, time string) string {
	// Replace slashes with hyphens in the date and concatenate with time
	return strings.ReplaceAll(date, "/", "-") + "-" + time
}

func pushWaveDataToES(waveDataList []WaveData) error {
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

		// Sanitize the DocumentID
		documentID := sanitizeDocumentID(waveData.Date, waveData.Time)

		// Prepare the request to index the data
		req := esapi.IndexRequest{
			Index:      elasticSearchIndexLesPierresNoires,
			DocumentID: documentID,
			Body:       strings.NewReader(string(dataJSON)),
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

func loadConfig(configFile string) Config {
	var config Config
	err := loadconfig.LoadConfig(configFile, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return config
}

func main() {
	ctx := context.Background()

	// Parse the config file path from the command line arguments
	configFile := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	// Load configuration
	config := loadConfig(*configFile)

	// Connect to the PostgreSQL database
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName)

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close(ctx)

	// Retrieve the latest session ID from the database
	sessionID, err := getSessionIDFromDB(ctx, conn)
	if err != nil {
		log.Fatalf("Failed to retrieve session ID: %v", err)
	}

	// Create an HTTP client and use the session ID to scrape data
	client := &http.Client{}
	waveDataList, err := getTableData(client, sessionID)
	if err != nil {
		log.Fatalf("Failed to get table data: %v", err)
	}

	// Push the data to Elasticsearch
	err = pushWaveDataToES(waveDataList)
	if err != nil {
		log.Fatalf("Failed to push data to Elasticsearch: %v", err)
	}
}

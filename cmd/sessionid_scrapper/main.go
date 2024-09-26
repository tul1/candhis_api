package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/tul1/candhis_api/internal/pkg/loadconfig"
)

const (
	scrapingbeeURL = "https://app.scrapingbee.com/api/v1/"

	// In order to create and activate the sessionID cookie we need to click in any buttom of the web page.
	// The ID "#idBtnAr" is the ID of the buttom "Archives" and the ID "#idBtnTR" is the buttom "temps reel".
	scrapingbeeJSScenario = `{"instructions":[{"click":"#idBtnAr"},{"wait":1000},{"click":"#idBtnTR"},{"wait":1000}]}`

	candhisURL = "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ=="
)

func getSessionIDWithScrapingbee(client *http.Client, scrapingbeeAPIKey string) (string, error) {
	params := url.Values{}
	params.Add("api_key", scrapingbeeAPIKey)
	params.Add("url", candhisURL)
	params.Add("js_scenario", scrapingbeeJSScenario)

	reqURL := fmt.Sprintf("%s?%s", scrapingbeeURL, params.Encode())
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request, url: %s, error: %w", reqURL, err)
	}

	err = req.ParseForm()
	if err != nil {
		return "", fmt.Errorf("failed to parse form, url: %s, error: %w", reqURL, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform request, url: %s, error: %w", reqURL, err)
	}

	cookies := resp.Header["Set-Cookie"]
	for _, cookie := range cookies {
		if strings.Contains(cookie, "PHPSESSID") {
			return strings.Split(cookie, ";")[0], nil
		}
	}

	return "", fmt.Errorf("failed to retrieve cookie PHPSESSID, url: %s", reqURL)
}

func updateSessionID(ctx context.Context, db *pgx.Conn, sessionID string) error {
	query := `UPDATE candhis_session SET id = $1, created_at = $2`
	_, err := db.Exec(ctx, query, sessionID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update session ID: %w", err)
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

	client := &http.Client{}
	phpsessid, err := getSessionIDWithScrapingbee(client, config.ScrapingbeeAPIKey)
	if err != nil {
		log.Fatalf("Failed to get sessionID: %v", err)
	}

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

	err = updateSessionID(ctx, conn, phpsessid)
	if err != nil {
		log.Fatalf("Failed to update session ID in database: %v", err)
	}

	log.Print("Session ID inserted into PostgreSQL")
}

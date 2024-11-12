package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/elastic/go-elasticsearch/v8"
	candhisapi "github.com/tul1/candhis_api/internal/application/candhis_api"
	"github.com/tul1/candhis_api/internal/infrastructure/persistence"
	"github.com/tul1/candhis_api/internal/pkg/configuration"
	"github.com/tul1/candhis_api/internal/pkg/logger"
	"github.com/tul1/candhis_api/internal/pkg/server"
)

func main() {
	log := logger.NewWithDefaultLogger()

	// Parse the config file path from the command line arguments
	configFile := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	// Load configuration
	config, err := configuration.Load[Config](*configFile)
	if err != nil {
		log.Errorf("Configuration error: %v", err)
		return
	}

	// Create Gin server
	s, err := server.NewGinServer(log, config.PublicURL, config.ServerPort)
	if err != nil {
		log.Errorf("Failed to create Gin server: %v", err)
		return
	}

	// Register candhis API handlers
	httpClient := http.Client{}
	defer httpClient.CloseIdleConnections()

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{config.ElasticsearchURL}})
	if err != nil {
		log.Errorf("Failed to create Elasticsearch client: %v", err)
		return
	}

	_ = candhisapi.NewCandhisAPI(
		s.GetRouter(),
		persistence.NewWaveData(esClient),
	)

	// Start server
	errCh := make(chan error)
	go func() {
		err := s.Start()
		if err != nil {
			errCh <- err
			return
		}
	}()

	// Manage app interruption to close server
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	select {
	case signalErr := <-signalCh:
		log.Infof("System interruption signal received: %s\n", signalErr.String())
	case err := <-errCh:
		log.Errorf("Error while running the application: %s\n", err)
	}

	// Stop server
	if err = s.Close(); err != nil {
		log.Errorf("Error while closing the application: %s\n", err)
		os.Exit(1)
	}
}

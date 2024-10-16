package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	candhisapi "github.com/tul1/candhis_api/internal/application/candhis_api"
	"github.com/tul1/candhis_api/internal/pkg/loadconfig"
	"github.com/tul1/candhis_api/internal/pkg/logger"
	"github.com/tul1/candhis_api/internal/pkg/server"
)

func loadConfig(configFile string) (*Config, error) {
	var config *Config
	if err := loadconfig.LoadConfig(configFile, config); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return config, nil
}

func main() {
	log := logger.NewWithDefaultLogger()

	// Parse the config file path from the command line arguments
	configFile := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Errorf("Configuration error: %v", err)
		return
	}

	// Create Gin server
	s, err := server.NewGinServer(config.PublicURL, config.ServerPort)
	if err != nil {
		log.Errorf("Failed to create Gin server: %v", err)
		return
	}

	// Register candhis API handlers
	_ = candhisapi.NewCandhisAPI(s.GetRouter())

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
		log.Infof("system interruption signal received: %s\n", signalErr.String())
	case err := <-errCh:
		log.Errorf("error while running the application: %s\n", err)
	}

	// Stop server
	if err = s.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

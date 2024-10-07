package service

import (
	"context"
	"fmt"

	"github.com/tul1/candhis_api/internal/application/repository"
)

type CandhisScraper interface {
	RetrieveAndStoreCandhisSessionID(ctx context.Context) error
	ScrapingCandhisCampaigns(ctx context.Context) error
}

type candhisScraper struct {
	sessionID               repository.SessionID
	waveData                repository.WaveData
	scrapingBeeClient       repository.ScrapingBeeClient
	candhisWebScraperClient repository.CandhisWebScraper
}

func NewCandhisScraper(
	sessionIDRepo repository.SessionID,
	waveDataRepo repository.WaveData,
	scrapingBeeClient repository.ScrapingBeeClient,
	candhisWebScraperClient repository.CandhisWebScraper,
) *candhisScraper {
	return &candhisScraper{
		sessionIDRepo,
		waveDataRepo,
		scrapingBeeClient,
		candhisWebScraperClient,
	}
}

func (s *candhisScraper) RetrieveAndStoreCandhisSessionID(ctx context.Context) error {
	candhisSessionID, err := s.scrapingBeeClient.GetCandhisSessionID()
	if err != nil {
		return fmt.Errorf("failed to get session ID from candhis web: %w", err)
	}

	err = s.sessionID.Update(ctx, candhisSessionID)
	if err != nil {
		return fmt.Errorf("failed to update session ID in database: %w", err)
	}

	return nil
}

const (
	candhisURL                         = "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ=="
	elasticSearchIndexLesPierresNoires = "les-pierres-noires"
)

func (s *candhisScraper) ScrapingCandhisCampaigns(ctx context.Context) error {
	candhisSessionID, err := s.sessionID.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get session ID from db: %w", err)
	}

	waveDataList, err := s.candhisWebScraperClient.GatherWavesDataFromWebTable(
		*candhisSessionID, candhisURL)
	if err != nil {
		return fmt.Errorf("failed to gather waves data from candhis web: %w", err)
	}

	for _, waveData := range waveDataList {
		err := s.waveData.Add(ctx, waveData, elasticSearchIndexLesPierresNoires)
		if err != nil {
			return fmt.Errorf("failed to push wave data to Elasticsearch: %w", err)
		}
	}

	return nil
}

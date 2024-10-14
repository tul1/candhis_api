package service

import (
	"context"
	"fmt"

	"github.com/tul1/candhis_api/internal/application/repository"
)

type CandhisCampaignsScraper interface {
	ScrapingCandhisCampaigns(ctx context.Context) error
}

type candhisCampaignsScraper struct {
	sessionID                        repository.SessionID
	waveData                         repository.WaveData
	candhisCampaignsWebScraperClient repository.CandhisCampaignsWebScraper
}

func NewCandhisCampaignsScraper(
	sessionIDRepo repository.SessionID,
	waveDataRepo repository.WaveData,
	candhisCampaignsWebScraperClient repository.CandhisCampaignsWebScraper,
) *candhisCampaignsScraper {
	return &candhisCampaignsScraper{
		sessionIDRepo,
		waveDataRepo,
		candhisCampaignsWebScraperClient,
	}
}

const (
	candhisURL                         = "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ=="
	elasticSearchIndexLesPierresNoires = "les-pierres-noires"
)

func (s *candhisCampaignsScraper) ScrapingCandhisCampaigns(ctx context.Context) error {
	candhisSessionID, err := s.sessionID.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get session ID from db: %w", err)
	}

	waveDataList, err := s.candhisCampaignsWebScraperClient.GatherWavesDataFromWebTable(
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

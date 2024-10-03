package service

import (
	"context"
	"fmt"

	"github.com/tul1/candhis_api/internal/application/repository"
)

type CandhisScraper interface {
	RetrieveAndStoreCandhisSessionID(ctx context.Context) error
}

type candhisScraper struct {
	sessionIDRepo     repository.SessionID
	scrapingBeeClient repository.ScrapingBeeClient
}

func NewCandhisScraper(sessionIDRepo repository.SessionID, scrapingBeeClient repository.ScrapingBeeClient) *candhisScraper {
	return &candhisScraper{sessionIDRepo, scrapingBeeClient}
}

func (s *candhisScraper) RetrieveAndStoreCandhisSessionID(ctx context.Context) error {
	candhisSessionID, err := s.scrapingBeeClient.GetCandhisSessionID()
	if err != nil {
		return fmt.Errorf("failed to get session ID from candhis web: %w", err)
	}

	err = s.sessionIDRepo.Update(ctx, &candhisSessionID)
	if err != nil {
		return fmt.Errorf("failed to update session ID in database: %w", err)
	}

	return nil
}

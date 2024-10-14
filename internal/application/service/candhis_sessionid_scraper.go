package service

import (
	"context"
	"fmt"

	"github.com/tul1/candhis_api/internal/application/repository"
)

type CandhisSessionIDScraper interface {
	FetchAndStoreSessionID(ctx context.Context) error
}

type candhisSessionIDScraper struct {
	sessionID                        repository.SessionID
	candhisSessionIDWebScraperClient repository.CandhisSessionIDWebScraper
}

func NewCandhisSessionIDScraper(
	sessionID repository.SessionID,
	candhisSessionIDWebScraperClient repository.CandhisSessionIDWebScraper,
) *candhisSessionIDScraper {
	return &candhisSessionIDScraper{sessionID, candhisSessionIDWebScraperClient}
}

func (s *candhisSessionIDScraper) FetchAndStoreSessionID(ctx context.Context) error {
	candhisSessionID, err := s.candhisSessionIDWebScraperClient.GetCandhisSessionID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get session ID from candhis web: %w", err)
	}

	err = s.sessionID.Update(ctx, candhisSessionID)
	if err != nil {
		return fmt.Errorf("failed to update session ID in database: %w", err)
	}

	return nil
}

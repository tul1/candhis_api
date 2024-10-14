package service_test

import (
	"context"
	"errors"
	"testing"

	clientmock "github.com/tul1/candhis_api/internal/application/repository/client_mock"
	persistencemock "github.com/tul1/candhis_api/internal/application/repository/persistence_mock"
	"github.com/tul1/candhis_api/internal/application/service"

	"github.com/stretchr/testify/assert"
	appmodel "github.com/tul1/candhis_api/internal/application/model"
	appmodeltest "github.com/tul1/candhis_api/internal/application/model/modeltest"

	"go.uber.org/mock/gomock"
)

func TestCandhisSessionIDScraper_FetchAndStoreSessionID_Success(t *testing.T) {
	mocks, candhisScraper := setupCandhisSessionIDScraperAndMocks(t)

	sessionID := appmodeltest.MustCreateCandhisSessionID(t, "valid-session-id")

	mocks.candhisSessionIDWebScraperClient.EXPECT().GetCandhisSessionID(gomock.Any()).Return(sessionID, nil)
	mocks.sessionID.EXPECT().Update(gomock.Any(), sessionID).Return(nil)

	err := candhisScraper.FetchAndStoreSessionID(context.Background())
	assert.NoError(t, err)
}

func TestCandhisSessionIDScraper_FetchAndStoreSessionID_ScrapingBeeClientFailure(t *testing.T) {
	mocks, candhisScraper := setupCandhisSessionIDScraperAndMocks(t)

	mocks.candhisSessionIDWebScraperClient.EXPECT().GetCandhisSessionID(gomock.Any()).Return(
		appmodel.CandhisSessionID{}, errors.New("error scraping bee"))

	err := candhisScraper.FetchAndStoreSessionID(context.Background())
	assert.EqualError(t, err, "failed to get session ID from candhis web: error scraping bee")
}

func TestCandhisSessionIDScraper_FetchAndStoreSessionID_UpdateSessionIDFailure(t *testing.T) {
	mocks, candhisScraper := setupCandhisSessionIDScraperAndMocks(t)

	sessionID := appmodeltest.MustCreateCandhisSessionID(t, "valid-session-id")

	mocks.candhisSessionIDWebScraperClient.EXPECT().GetCandhisSessionID(gomock.Any()).Return(sessionID, nil)
	mocks.sessionID.EXPECT().Update(gomock.Any(), sessionID).Return(errors.New("error db"))

	err := candhisScraper.FetchAndStoreSessionID(context.Background())
	assert.EqualError(t, err, "failed to update session ID in database: error db")
}

type sessionIDTestingMocks struct {
	sessionID                        *persistencemock.MockSessionID
	candhisSessionIDWebScraperClient *clientmock.MockCandhisSessionIDWebScraper
}

func setupCandhisSessionIDScraperAndMocks(t *testing.T) (sessionIDTestingMocks, service.CandhisSessionIDScraper) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockSessionIDRepo := persistencemock.NewMockSessionID(ctrl)
	mockCandhisSessionIDWebScraperClient := clientmock.NewMockCandhisSessionIDWebScraper(ctrl)

	return sessionIDTestingMocks{
		sessionID:                        mockSessionIDRepo,
		candhisSessionIDWebScraperClient: mockCandhisSessionIDWebScraperClient,
	}, service.NewCandhisSessionIDScraper(mockSessionIDRepo, mockCandhisSessionIDWebScraperClient)
}

package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tul1/candhis_api/internal/application/model"
	"github.com/tul1/candhis_api/internal/application/model/modeltest"
	clientmock "github.com/tul1/candhis_api/internal/application/repository/client_mock"
	persistencemock "github.com/tul1/candhis_api/internal/application/repository/persistence_mock"
	"github.com/tul1/candhis_api/internal/application/service"
	"go.uber.org/mock/gomock"
)

func TestRetrieveAndStoreCandhisSessionID_Success(t *testing.T) {
	mocks, candhisScraper := setupCandhisScrapperAndMocks(t)

	sessionID := modeltest.MustCreateCandhisSessionID(t, "valid-session-id")

	mocks.mockScrapingBeeClient.EXPECT().GetCandhisSessionID().Return(sessionID, nil)
	mocks.sessionID.EXPECT().Update(gomock.Any(), &sessionID).Return(nil)

	err := candhisScraper.RetrieveAndStoreCandhisSessionID(context.Background())
	assert.NoError(t, err)
}

func TestRetrieveAndStoreCandhisSessionID_ScrapingBeeClientFailure(t *testing.T) {
	mocks, candhisScraper := setupCandhisScrapperAndMocks(t)

	mocks.mockScrapingBeeClient.EXPECT().GetCandhisSessionID().Return(
		model.CandhisSessionID{}, errors.New("error scraping bee"))

	err := candhisScraper.RetrieveAndStoreCandhisSessionID(context.Background())
	assert.EqualError(t, err, "failed to get session ID from candhis web: error scraping bee")
}

func TestRetrieveAndStoreCandhisSessionID_UpdateSessionIDFailure(t *testing.T) {
	mocks, candhisScraper := setupCandhisScrapperAndMocks(t)

	sessionID := modeltest.MustCreateCandhisSessionID(t, "valid-session-id")

	mocks.mockScrapingBeeClient.EXPECT().GetCandhisSessionID().Return(sessionID, nil)
	mocks.sessionID.EXPECT().Update(gomock.Any(), &sessionID).Return(errors.New("error db"))

	err := candhisScraper.RetrieveAndStoreCandhisSessionID(context.Background())
	assert.EqualError(t, err, "failed to update session ID in database: error db")
}

type testingMocks struct {
	sessionID             *persistencemock.MockSessionID
	mockScrapingBeeClient *clientmock.MockScrapingBeeClient
}

func setupCandhisScrapperAndMocks(t *testing.T) (testingMocks, service.CandhisScraper) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockSessionIDRepo := persistencemock.NewMockSessionID(ctrl)
	mockScrapingBeeClient := clientmock.NewMockScrapingBeeClient(ctrl)

	return testingMocks{mockSessionIDRepo, mockScrapingBeeClient},
		service.NewCandhisScraper(mockSessionIDRepo, mockScrapingBeeClient)

}

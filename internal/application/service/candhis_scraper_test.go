package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	appmodel "github.com/tul1/candhis_api/internal/application/model"
	appmodeltest "github.com/tul1/candhis_api/internal/application/model/modeltest"
	"github.com/tul1/candhis_api/internal/domain/model"
	"github.com/tul1/candhis_api/internal/domain/model/modeltest"

	clientmock "github.com/tul1/candhis_api/internal/application/repository/client_mock"
	persistencemock "github.com/tul1/candhis_api/internal/application/repository/persistence_mock"
	"github.com/tul1/candhis_api/internal/application/service"
	"go.uber.org/mock/gomock"
)

func TestRetrieveAndStoreCandhisSessionID_Success(t *testing.T) {
	mocks, candhisScraper := setupCandhisScraperAndMocks(t)

	sessionID := appmodeltest.MustCreateCandhisSessionID(t, "valid-session-id")

	mocks.mockCandhisSessionIDWebScraperClient.EXPECT().GetCandhisSessionID(gomock.Any()).Return(sessionID, nil)
	mocks.sessionID.EXPECT().Update(gomock.Any(), sessionID).Return(nil)

	err := candhisScraper.RetrieveAndStoreCandhisSessionID(context.Background())
	assert.NoError(t, err)
}

func TestRetrieveAndStoreCandhisSessionID_ScrapingBeeClientFailure(t *testing.T) {
	mocks, candhisScraper := setupCandhisScraperAndMocks(t)

	mocks.mockCandhisSessionIDWebScraperClient.EXPECT().GetCandhisSessionID(gomock.Any()).Return(
		appmodel.CandhisSessionID{}, errors.New("error scraping bee"))

	err := candhisScraper.RetrieveAndStoreCandhisSessionID(context.Background())
	assert.EqualError(t, err, "failed to get session ID from candhis web: error scraping bee")
}

func TestRetrieveAndStoreCandhisSessionID_UpdateSessionIDFailure(t *testing.T) {
	mocks, candhisScraper := setupCandhisScraperAndMocks(t)

	sessionID := appmodeltest.MustCreateCandhisSessionID(t, "valid-session-id")

	mocks.mockCandhisSessionIDWebScraperClient.EXPECT().GetCandhisSessionID(gomock.Any()).Return(sessionID, nil)
	mocks.sessionID.EXPECT().Update(gomock.Any(), sessionID).Return(errors.New("error db"))

	err := candhisScraper.RetrieveAndStoreCandhisSessionID(context.Background())
	assert.EqualError(t, err, "failed to update session ID in database: error db")
}

func TestScrapingCandhisCampaigns_Success(t *testing.T) {
	mocks, candhisScraper := setupCandhisScraperAndMocks(t)

	sessionID := appmodeltest.MustCreateCandhisSessionID(t, "valid-session-id")
	wavesData := []model.WaveData{
		modeltest.MustCreateWaveData(t, "17/09/2024", "09:00", "0.6", "1.1", "4.7", "8", "32", "15"),
		modeltest.MustCreateWaveData(t, "17/09/2024", "08:30", "0.5", "0.9", "4.8", "4", "47", "15"),
	}

	mocks.sessionID.EXPECT().Get(gomock.Any()).Return(&sessionID, nil)
	mocks.mockCandhisWebScraperClient.EXPECT().
		GatherWavesDataFromWebTable(sessionID, "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ==").
		Return(wavesData, nil)
	mocks.waveData.EXPECT().Add(gomock.Any(), wavesData[0], "les-pierres-noires").Return(nil)
	mocks.waveData.EXPECT().Add(gomock.Any(), wavesData[1], "les-pierres-noires").Return(nil)

	err := candhisScraper.ScrapingCandhisCampaigns(context.Background())
	assert.NoError(t, err)
}

func TestScrapingCandhisCampaigns_SessionIDFailure(t *testing.T) {
	mocks, candhisScraper := setupCandhisScraperAndMocks(t)

	mocks.sessionID.EXPECT().Get(gomock.Any()).Return(nil, errors.New("error db"))

	err := candhisScraper.ScrapingCandhisCampaigns(context.Background())
	assert.EqualError(t, err, "failed to get session ID from db: error db")
}

func TestScrapingCandhisCampaigns_GatherWavesDataFailure(t *testing.T) {
	mocks, candhisScraper := setupCandhisScraperAndMocks(t)

	sessionID := appmodeltest.MustCreateCandhisSessionID(t, "valid-session-id")

	mocks.sessionID.EXPECT().Get(gomock.Any()).Return(&sessionID, nil)
	mocks.mockCandhisWebScraperClient.EXPECT().
		GatherWavesDataFromWebTable(sessionID, "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ==").
		Return(nil, errors.New("error web"))

	err := candhisScraper.ScrapingCandhisCampaigns(context.Background())
	assert.EqualError(t, err, "failed to gather waves data from candhis web: error web")
}

func TestScrapingCandhisCampaigns_AddWaveDataFailure(t *testing.T) {
	mocks, candhisScraper := setupCandhisScraperAndMocks(t)

	sessionID := appmodeltest.MustCreateCandhisSessionID(t, "valid-session-id")
	wavesData := []model.WaveData{
		modeltest.MustCreateWaveData(t, "17/09/2024", "09:00", "0.6", "1.1", "4.7", "8", "32", "15"),
		modeltest.MustCreateWaveData(t, "17/09/2024", "08:30", "0.5", "0.9", "4.8", "4", "47", "15"),
	}

	mocks.sessionID.EXPECT().Get(gomock.Any()).Return(&sessionID, nil)
	mocks.mockCandhisWebScraperClient.EXPECT().
		GatherWavesDataFromWebTable(sessionID, "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ==").
		Return(wavesData, nil)
	mocks.waveData.EXPECT().Add(gomock.Any(), wavesData[0], "les-pierres-noires").Return(errors.New("error elasticsearch"))

	err := candhisScraper.ScrapingCandhisCampaigns(context.Background())
	assert.EqualError(t, err, "failed to push wave data to Elasticsearch: error elasticsearch")
}

type testingMocks struct {
	sessionID                            *persistencemock.MockSessionID
	waveData                             *persistencemock.MockWaveData
	mockCandhisSessionIDWebScraperClient *clientmock.MockCandhisSessionIDWebScraper
	mockCandhisWebScraperClient          *clientmock.MockCandhisWebScraper
}

func setupCandhisScraperAndMocks(t *testing.T) (testingMocks, service.CandhisScraper) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockSessionIDRepo := persistencemock.NewMockSessionID(ctrl)
	mockWaveDataRepo := persistencemock.NewMockWaveData(ctrl)
	mockCandhisSessionIDWebScraperClient := clientmock.NewMockCandhisSessionIDWebScraper(ctrl)
	mockCandhisWebScraperClient := clientmock.NewMockCandhisWebScraper(ctrl)

	return testingMocks{
		sessionID:                            mockSessionIDRepo,
		waveData:                             mockWaveDataRepo,
		mockCandhisSessionIDWebScraperClient: mockCandhisSessionIDWebScraperClient,
		mockCandhisWebScraperClient:          mockCandhisWebScraperClient,
	}, service.NewCandhisScraper(mockSessionIDRepo, mockWaveDataRepo, mockCandhisSessionIDWebScraperClient, mockCandhisWebScraperClient)
}

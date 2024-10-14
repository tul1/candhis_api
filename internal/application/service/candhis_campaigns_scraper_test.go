package service_test

import (
	"context"
	"errors"
	"testing"

	clientmock "github.com/tul1/candhis_api/internal/application/repository/client_mock"
	persistencemock "github.com/tul1/candhis_api/internal/application/repository/persistence_mock"

	"github.com/stretchr/testify/assert"
	appmodeltest "github.com/tul1/candhis_api/internal/application/model/modeltest"
	"github.com/tul1/candhis_api/internal/application/service"
	"github.com/tul1/candhis_api/internal/domain/model"
	"github.com/tul1/candhis_api/internal/domain/model/modeltest"
	"go.uber.org/mock/gomock"
)

func TestCandhisCampaignsScraper_FetchAndStoreWaveData_Success(t *testing.T) {
	mocks, candhisScraper := setupCandhisCampaignsScraperAndMocks(t)

	sessionID := appmodeltest.MustCreateCandhisSessionID(t, "valid-session-id")
	wavesData := []model.WaveData{
		modeltest.MustCreateWaveData(t, "17/09/2024", "09:00", "0.6", "1.1", "4.7", "8", "32", "15"),
		modeltest.MustCreateWaveData(t, "17/09/2024", "08:30", "0.5", "0.9", "4.8", "4", "47", "15"),
	}

	mocks.sessionID.EXPECT().Get(gomock.Any()).Return(&sessionID, nil)
	mocks.candhisCampaignsWebScraper.EXPECT().
		GatherWavesDataFromWebTable(sessionID, "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ==").
		Return(wavesData, nil)
	mocks.waveData.EXPECT().Add(gomock.Any(), wavesData[0], "les-pierres-noires").Return(nil)
	mocks.waveData.EXPECT().Add(gomock.Any(), wavesData[1], "les-pierres-noires").Return(nil)

	err := candhisScraper.FetchAndStoreWaveData(context.Background())
	assert.NoError(t, err)
}

func TestCandhisCampaignsScraper_FetchAndStoreWaveData_SessionIDFailure(t *testing.T) {
	mocks, candhisScraper := setupCandhisCampaignsScraperAndMocks(t)

	mocks.sessionID.EXPECT().Get(gomock.Any()).Return(nil, errors.New("error db"))

	err := candhisScraper.FetchAndStoreWaveData(context.Background())
	assert.EqualError(t, err, "failed to get session ID from db: error db")
}

func TestCandhisCampaignsScraper_FetchAndStoreWaveData_GatherWavesDataFailure(t *testing.T) {
	mocks, candhisScraper := setupCandhisCampaignsScraperAndMocks(t)

	sessionID := appmodeltest.MustCreateCandhisSessionID(t, "valid-session-id")

	mocks.sessionID.EXPECT().Get(gomock.Any()).Return(&sessionID, nil)
	mocks.candhisCampaignsWebScraper.EXPECT().
		GatherWavesDataFromWebTable(sessionID, "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ==").
		Return(nil, errors.New("error web"))

	err := candhisScraper.FetchAndStoreWaveData(context.Background())
	assert.EqualError(t, err, "failed to gather waves data from candhis web: error web")
}

func TestCandhisCampaignsScraper_FetchAndStoreWaveData_AddWaveDataFailure(t *testing.T) {
	mocks, candhisScraper := setupCandhisCampaignsScraperAndMocks(t)

	sessionID := appmodeltest.MustCreateCandhisSessionID(t, "valid-session-id")
	wavesData := []model.WaveData{
		modeltest.MustCreateWaveData(t, "17/09/2024", "09:00", "0.6", "1.1", "4.7", "8", "32", "15"),
		modeltest.MustCreateWaveData(t, "17/09/2024", "08:30", "0.5", "0.9", "4.8", "4", "47", "15"),
	}

	mocks.sessionID.EXPECT().Get(gomock.Any()).Return(&sessionID, nil)
	mocks.candhisCampaignsWebScraper.EXPECT().
		GatherWavesDataFromWebTable(sessionID, "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ==").
		Return(wavesData, nil)
	mocks.waveData.EXPECT().Add(gomock.Any(), wavesData[0], "les-pierres-noires").Return(errors.New("error elasticsearch"))

	err := candhisScraper.FetchAndStoreWaveData(context.Background())
	assert.EqualError(t, err, "failed to push wave data to Elasticsearch: error elasticsearch")
}

type campaignsTestingMocks struct {
	sessionID                  *persistencemock.MockSessionID
	waveData                   *persistencemock.MockWaveData
	candhisCampaignsWebScraper *clientmock.MockCandhisCampaignsWebScraper
}

func setupCandhisCampaignsScraperAndMocks(t *testing.T) (campaignsTestingMocks, service.CandhisCampaignsScraper) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockSessionIDRepo := persistencemock.NewMockSessionID(ctrl)
	mockWaveDataRepo := persistencemock.NewMockWaveData(ctrl)
	mockCandhisCampaignsWebScraperClient := clientmock.NewMockCandhisCampaignsWebScraper(ctrl)

	return campaignsTestingMocks{
		sessionID:                  mockSessionIDRepo,
		waveData:                   mockWaveDataRepo,
		candhisCampaignsWebScraper: mockCandhisCampaignsWebScraperClient,
	}, service.NewCandhisCampaignsScraper(mockSessionIDRepo, mockWaveDataRepo, mockCandhisCampaignsWebScraperClient)
}

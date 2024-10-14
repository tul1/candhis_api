package client_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	appmodeltest "github.com/tul1/candhis_api/internal/application/model/modeltest"
	repo "github.com/tul1/candhis_api/internal/application/repository"
	"github.com/tul1/candhis_api/internal/domain/model"
	"github.com/tul1/candhis_api/internal/domain/model/modeltest"
	"github.com/tul1/candhis_api/internal/infrastructure/client"
)

const mockHTMLResponse = `
<!DOCTYPE html>
<html>
<body>
	<table class="table table-striped table-bordered table-sm">
	<thead>
		<tr class="table-warning text-center">
		<th class="clALGTab"><strong class="clALGTab">Date</strong></th>
		<th class="clALGTab"><strong class="clALGTab">Heure (TU)</strong></th>
		<th class="clALGTab"><strong class="clALGTab">H1/3 (m)</strong></th>
		<th class="clALGTab"><strong class="clALGTab">Hmax (m)</strong></th>
		<th class="clALGTab"><strong class="clALGTab">Th1/3 (s)</strong></th>
		<th class="clALGTab"><strong class="clALGTab">Dir. au pic (°)</strong></th>
		<th class="clALGTab"><strong class="clALGTab">Etal. au pic (°)</strong></th>
		<th class="clALGTab"><strong class="clALGTab">Temp. mer (°C)</strong></th>
		</tr>
	</thead>
	<tbody>
		<tr>
		<td class="text-center clALGTab"><span class="clALGTab">17/09/2024</span></td>
		<td class="text-center clALGTab"><span class="clALGTab">09:00</span></td>
		<td class="text-center clALGTab"><span class="clALGTab">0.6</span></td>
		<td class="text-center clALGTab"><span class="clALGTab">1.1</span></td>
		<td class="text-center clALGTab"><span class="clALGTab">4.7</span></td>
		<td class="text-center clALGTab"><span class="clALGTab">8</span></td>
		<td class="text-center clALGTab"><span class="clALGTab">32</span></td>
		<td class="text-center clALGTab"><span class="clALGTab">15</span></td>
		</tr>
		<tr>
		<td class="text-center clALGTab"><span class="clALGTab">17/09/2024</span></td>
		<td class="text-center clALGTab"><span class="clALGTab">08:30</span></td>
		<td class="text-center clALGTab"><span class="clALGTab">0.5</span></td>
		<td class="text-center clALGTab"><span class="clALGTab">0.9</span></td>
		<td class="text-center clALGTab"><span class="clALGTab">4.8</span></td>
		<td class="text-center clALGTab"><span class="clALGTab">4</span></td>
		<td class="text-center clALGTab"><span class="clALGTab">47</span></td>
		<td class="text-center clALGTab"><span class="clALGTab">15</span></td>
		</tr> 
	</tbody>
	</table>
</body>
</html>
`

func MockHTTPResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestGatherWavesDataFromWebTable_Success(t *testing.T) {
	mockHandler := func(req *http.Request) *http.Response {
		return MockHTTPResponse(200, mockHTMLResponse)
	}
	scraper := setupMockCandhisCampaignsWebScraper(t, mockHandler)

	waveData, err := scraper.GatherWavesDataFromWebTable(
		appmodeltest.MustCreateCandhisSessionID(t, "valid-session-id"), "http://fake.url")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(waveData))

	expected := []model.WaveData{
		modeltest.MustCreateWaveData(t, "17/09/2024", "09:00", "0.6", "1.1", "4.7", "8", "32", "15"),
		modeltest.MustCreateWaveData(t, "17/09/2024", "08:30", "0.5", "0.9", "4.8", "4", "47", "15"),
	}

	assert.Equal(t, expected, waveData, "Expected correct parsed wave data")
}

func TestGatherWavesDataFromWebTable_EmptyResponse(t *testing.T) {
	mockHandler := func(req *http.Request) *http.Response {
		return MockHTTPResponse(200, "")
	}
	scraper := setupMockCandhisCampaignsWebScraper(t, mockHandler)

	waveData, err := scraper.GatherWavesDataFromWebTable(
		appmodeltest.MustCreateCandhisSessionID(t, "valid-session-id"), "http://fake.url")
	assert.NoError(t, err)
	assert.Empty(t, waveData)
}

type mockRoundTripper struct {
	mockHandler func(req *http.Request) *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.mockHandler(req), nil
}

func setupMockCandhisCampaignsWebScraper(t *testing.T, mockHandler func(req *http.Request) *http.Response) repo.CandhisCampaignsWebScraper {
	t.Helper()

	mockClient := &http.Client{Transport: &mockRoundTripper{mockHandler: mockHandler}}

	return client.NewCandhisCampaignsWebScraper(mockClient)
}

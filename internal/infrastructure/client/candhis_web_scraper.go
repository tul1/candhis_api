package client

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	appmodel "github.com/tul1/candhis_api/internal/application/model"
	"github.com/tul1/candhis_api/internal/domain/model"
)

const expectedCellsNum = 8

type candhisWebScraper struct {
	client *http.Client
}

func NewCandhisWebScraper(client *http.Client) *candhisWebScraper {
	return &candhisWebScraper{client}
}

func (c *candhisWebScraper) GatherWavesDataFromWebTable(
	candhisSessionID appmodel.CandhisSessionID,
	candhisURL string,
) ([]model.WaveData, error) {
	req, err := http.NewRequest(http.MethodGet, candhisURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request, url: %s, error: %w", candhisURL, err)
	}

	req.Header.Set("Accept", "text/html")
	req.Header.Set("Cookie", fmt.Sprintf("acceptCookies=true; %s", candhisSessionID.ID()))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request, url: %s, error: %w", candhisURL, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed parse HTML: %w", err)
	}

	var waveDataList []model.WaveData
	doc.Find("table.table-striped.table-bordered.table-sm").Each(func(index int, table *goquery.Selection) {
		table.Find("tr").Each(func(rowIndex int, row *goquery.Selection) {
			waveData, err := c.parseRowOfWebTable(row.Find("td"))
			if err != nil {
				log.Printf("Skipping row due to error: %v", err)
				return
			}

			waveDataList = append(waveDataList, waveData)
		})
	})

	return waveDataList, nil
}

func (c *candhisWebScraper) parseRowOfWebTable(cells *goquery.Selection) (model.WaveData, error) {
	if cells.Length() != expectedCellsNum {
		return model.WaveData{}, fmt.Errorf("expected %d cells, but got %d", expectedCellsNum, cells.Length())
	}

	values := make([]string, expectedCellsNum)
	cells.Each(func(cellIndex int, cell *goquery.Selection) {
		values[cellIndex] = strings.TrimSpace(cell.Text())
	})

	return model.NewWaveData(values[0], values[1], values[2], values[3],
		values[4], values[5], values[6], values[7])
}

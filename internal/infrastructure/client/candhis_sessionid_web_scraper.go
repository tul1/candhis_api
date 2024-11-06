package client

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/network"
	"github.com/tul1/candhis_api/internal/application/model"
)

//go:generate mockgen -package scrapermock -destination=./scraper_mock/scraper_mock.go -source=candhis_sessionid_web_scraper.go ScraperMock
type Scraper interface {
	Run(ctx context.Context, targetWeb string, actionFunc func(context.Context) error) error
}

type candhisSessionIDWebScraper struct {
	chromeScraper Scraper
	targetWeb     string
}

func NewCandhisSessionIDWebScraper(chromeScraper Scraper, targetWeb string) *candhisSessionIDWebScraper {
	return &candhisSessionIDWebScraper{chromeScraper, targetWeb}
}

func (c *candhisSessionIDWebScraper) GetCandhisSessionID(ctx context.Context) (model.CandhisSessionID, error) {
	var cookies []*network.Cookie
	getCookies := func(ctx context.Context) error {
		var err error
		cookies, err = network.GetCookies().Do(ctx)
		return err
	}

	err := c.chromeScraper.Run(ctx, c.targetWeb, getCookies)
	if err != nil {
		return model.CandhisSessionID{}, fmt.Errorf("failed while running chromedp tasks to retrieve session id: %w", err)
	}

	for _, cookie := range cookies {
		if cookie.Name == "PHPSESSID" {
			return model.NewCandhisSessionID(cookie.Value, nil)
		}
	}

	return model.CandhisSessionID{}, fmt.Errorf("failed to retrieve session id: %w", err)
}

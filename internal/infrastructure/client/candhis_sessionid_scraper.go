package client

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/tul1/candhis_api/internal/application/model"
)

type candhisSessionIDWebScraper struct {
	chromeURL string
	chromeID  string
	targetWeb string
}

func NewCandhisSessionIDWebScraper(chromeURL, chromeID, targetWeb string) *candhisSessionIDWebScraper {
	return &candhisSessionIDWebScraper{chromeURL, chromeID, targetWeb}
}

func (c *candhisSessionIDWebScraper) GetCandhisSessionID(ctx context.Context) (model.CandhisSessionID, error) {
	chromodpWS := fmt.Sprintf("ws://%s/devtools/browser/%s", c.chromeURL, c.chromeID)
	ctx, cancel := chromedp.NewRemoteAllocator(ctx, chromodpWS)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var cookies []*network.Cookie
	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(c.targetWeb),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			cookies, err = network.GetCookies().Do(ctx)
			return err
		}),
	)
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

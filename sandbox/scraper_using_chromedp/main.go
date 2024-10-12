package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func main() {
	ctx, cancel := chromedp.NewRemoteAllocator(context.Background(), "ws://0.0.0.0:9222/devtools/browser/681683b0-eb0d-428c-a05b-6eedfc0f4674")
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var cookies []*network.Cookie
	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate("https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ=="),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			cookies, err = network.GetCookies().Do(ctx)
			return err
		}),
	)
	if err != nil {
		log.Fatal("Failed to retrieve cookies: ", err)
	}

	for _, cookie := range cookies {
		if cookie.Name == "PHPSESSID" {
			fmt.Println("PHPSESSID:", cookie.Value)
		}
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func getPHPSESSID(url string) (string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),              // Run in headless mode
		chromedp.Flag("disable-gpu", true),           // Disable GPU
		chromedp.Flag("no-sandbox", true),            // Disable sandboxing
		chromedp.Flag("disable-dev-shm-usage", true), // Overcome limited resource problems in Docker
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var sessionCookie string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetCookies().Do(ctx)
			if err != nil {
				return err
			}
			for _, v := range cookies {
				if v.Name == "PHPSESSID" {
					sessionCookie = v.Name + "=" + v.Value
				}
			}
			return nil
		}),
		chromedp.Reload(),
	)

	if err != nil {
		return "", fmt.Errorf("error gathering cookies: %v", err)
	}

	return sessionCookie, nil
}

// getTableData uses the PHPSESSID and other cookies to make a request and extract the table data
func getTableData(url, cookieHeader string) error {
	client := &http.Client{}

	// Create a new request with the dynamically retrieved PHPSESSID and cookies
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Add the cookies to the request header
	req.Header.Set("Cookie", fmt.Sprintf("acceptCookies=true; %s", cookieHeader))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:130.0) Gecko/20100101 Firefox/130.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform data request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the HTML response using goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Find the table and extract the data
	doc.Find("table.table-striped.table-bordered.table-sm").Each(func(index int, table *goquery.Selection) {
		table.Find("tr").Each(func(rowIndex int, row *goquery.Selection) {
			row.Find("td").Each(func(cellIndex int, cell *goquery.Selection) {
				fmt.Printf("Row %d, Cell %d: %s\n", rowIndex, cellIndex, strings.TrimSpace(cell.Text()))
			})
		})
	})

	return nil
}

func main() {
	url := "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ=="

	// Step 1: Use Chromedp to dynamically get PHPSESSID and other cookies
	cookieHeader, err := getPHPSESSID(url)
	if err != nil {
		log.Fatalf("Error gathering cookies: %v", err)
	}
	fmt.Println("Cookies:", cookieHeader)

	// Step 2: Use the cookies to gather table data
	err = getTableData(url, cookieHeader)
	if err != nil {
		log.Fatalf("Error gathering table data: %v", err)
	}
}

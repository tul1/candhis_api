package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// getPHPSESSID makes the initial request and retrieves the PHPSESSID from the Set-Cookie header
func getPHPSESSID(url string) (string, error) {
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set User-Agent to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:130.0) Gecko/20100101 Firefox/130.0")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	// Extract the Set-Cookie header
	cookies := resp.Header["Set-Cookie"]
	var phpsessid string
	for _, cookie := range cookies {
		if strings.Contains(cookie, "PHPSESSID") {
			// Extract the PHPSESSID part
			phpsessid = strings.Split(cookie, ";")[0]
			break
		}
	}

	if phpsessid == "" {
		return "", fmt.Errorf("PHPSESSID not found in the response")
	}

	return phpsessid, nil
}

// getTableData uses the PHPSESSID and additional headers to request the page and scrape the table data
func getTableData(url, phpsessid string) error {
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Add the necessary headers
	req.Header.Set("Cookie", fmt.Sprintf("lang=fr; _pk_id.42.9752=d3fff9e6205683e7.1719322715.; acceptCookies=true; %s", phpsessid))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:130.0) Gecko/20100101 Firefox/130.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br") // Ensure we accept compressed responses
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMTEwMQ==")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform data request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response is compressed and decompress it if needed
	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %v", err)
		}
		defer reader.Close()
	} else {
		reader = resp.Body
	}

	// DEBUG: Print the response body to verify the HTML content
	body, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	fmt.Println("HTML Response:\n", string(body)) // Print the full HTML response for debugging

	// Parse the HTML response using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
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

	// Step 1: Get the PHPSESSID dynamically
	phpsessid, err := getPHPSESSID(url)
	if err != nil {
		log.Fatalf("Error gathering PHPSESSID: %v", err)
	}
	fmt.Println("PHPSESSID:", phpsessid)

	// Step 2: Use the PHPSESSID to gather table data
	err = getTableData(url, phpsessid)
	if err != nil {
		log.Fatalf("Error gathering table data: %v", err)
	}
}

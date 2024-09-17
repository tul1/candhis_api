package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ=="

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Accept", "text/html")
	req.Header.Set("Cookie", "acceptCookies=true; PHPSESSID=eqki3teu3froaqdo7p0ld0iadh") //TODO: I cannot gather this sessionID automatically

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	doc.Find("table.table-striped.table-bordered.table-sm").Each(func(index int, table *goquery.Selection) {
		table.Find("tr").Each(func(rowIndex int, row *goquery.Selection) {
			row.Find("td").Each(func(cellIndex int, cell *goquery.Selection) {
				fmt.Printf("Row %d, Cell %d: %s\n", rowIndex, cellIndex, strings.TrimSpace(cell.Text()))
			})
		})
	})
}

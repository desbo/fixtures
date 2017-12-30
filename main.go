package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/desbo/fixtures/scraper"

	"github.com/namsral/microdata"
)

func isSportsEvent(item *microdata.Item) bool {
	for _, t := range item.Types {
		if t == "http://schema.org/SportsEvent" {
			return true
		}
	}

	return false
}

func main() {
	file, err := os.Open("./example.html")

	if err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse(scraper.BaseURL)

	if err != nil {
		log.Fatal(err)
	}

	data, err := microdata.ParseHTML(file, "text/html", u)

	for _, item := range data.Items[:2] {
		if isSportsEvent(item) {
			fixture, _ := scraper.NewFixture(item)
			fmt.Println(fixture)
			fmt.Println("!!!")
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}

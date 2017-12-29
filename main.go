package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/go-openapi/strfmt"

	"github.com/desbo/fixtures/models"
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

	for _, item := range data.Items {
		if isSportsEvent(item) {
			fixture := &models.Fixture{
				Name:   item.Properties["description"][0].(string),
				Status: scraper.ParseStatus(item),
			}

			time, err := scraper.ParseTime(item)

			fmt.Println(time, err)

			if err != nil {
				fixture.Time = strfmt.DateTime(time)
			}

			// fmt.Println(fixture)
			// fmt.Println(item.Properties)
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}

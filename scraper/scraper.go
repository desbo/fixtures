package scraper

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/desbo/fixtures/restapi/operations/fixtures"
	"github.com/namsral/microdata"
)

const BaseURL = "https://www.tabletennis365.com/CentralLondon/Fixtures/Winter_2017-18/All_Divisions?vm=1"

var centralLondonDivisionIDs = map[int64]int{
	1: 5596,
	2: 5597,
	3: 5598,
	4: 5599,
	5: 5600,
	6: 5601,
}

func boolToParam(b bool) string {
	if b {
		return "True"
	}

	return "False"
}

func createURL(params *fixtures.ListFixturesParams) (*url.URL, error) {
	u, err := url.Parse(BaseURL)

	if err != nil {
		return nil, err
	}

	setQuery := func(key string, value interface{}) {
		q := u.Query()
		q.Set(key, fmt.Sprintf("%v", value))
		u.RawQuery = q.Encode()
	}

	if params.ClDivision != nil {
		if d, ok := centralLondonDivisionIDs[*params.ClDivision]; ok {
			setQuery("d", d)
		}
	} else if params.DivisionID != nil {
		setQuery("d", *params.DivisionID)
	}

	if params.ClubID != nil {
		setQuery("cl", *params.ClubID)
	}

	if params.ShowCompleted != nil {
		setQuery("hc", boolToParam(!*params.ShowCompleted))
	}

	return u, nil
}

func isSportsEvent(item *microdata.Item) bool {
	for _, t := range item.Types {
		if t == "http://schema.org/SportsEvent" {
			return true
		}
	}

	return false
}

func ParseStatus(item *microdata.Item) string {
	s := item.Properties["eventStatus"][0].(string)
	return strings.ToLower(strings.Replace(s, "http://schema.org/Event", "", 1))
}

func ParseTime(item *microdata.Item) (time.Time, error) {
	layout := "Monday 2 January 2006 @ 15:04"
	cleaned := regexp.MustCompile(`(\d{1,2})(st|nd|rd|th)`).ReplaceAllString(item.Properties["startDate"][0].(string), "$1")
	return time.Parse(layout, cleaned)
}

// func AsFixture(item *microdata.Item) (error, *models.Fixture) {
// 	fixture := &models.Fixture{}

// 	fixture.Name = item.Properties["description"][0].(string)
// }

func Scrape(params *fixtures.ListFixturesParams) error {
	u, err := createURL(params)

	if err != nil {
		return err
	}

	_, err = microdata.ParseURL(u.String())

	if err != nil {
		return err
	}

	return nil
}

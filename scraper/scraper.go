package scraper

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/desbo/fixtures/models"
	"github.com/desbo/fixtures/restapi/operations/fixtures"
	"github.com/go-openapi/strfmt"
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

func parseStatus(item *microdata.Item) string {
	s, ok := item.Properties["eventStatus"]

	if !ok || len(s) == 0 {
		// default status for events missing a status
		return "scheduled"
	}

	return strings.ToLower(strings.Replace(s[0].(string), "http://schema.org/Event", "", 1))
}

func parseTime(item *microdata.Item) (time.Time, error) {
	layout := "Monday 2 January 2006 @ 15:04"
	fixed := strings.Replace(item.Properties["startDate"][0].(string), "\u00a0", " ", -1)
	cleaned := regexp.MustCompile(`(\d{1,2})(st|nd|rd|th)`).ReplaceAllString(fixed, "$1")
	return time.Parse(layout, cleaned)
}

// returns home, away, error
func parseTeams(item *microdata.Item, eventName string) (*models.Team, *models.Team, error) {
	names := regexp.MustCompile(`(.*) v's (.*) at .*`)
	matches := names.FindStringSubmatch(eventName)

	if len(matches) < 3 {
		return nil, nil, fmt.Errorf("Unable to find 2 team names in %s", eventName)
	}

	home := &models.Team{
		Name: matches[1],
	}

	away := &models.Team{
		Name: matches[2],
	}

	if perf := item.Properties["performer"]; perf != nil && len(perf) > 0 {
		playerR := regexp.MustCompile(`([\w\s]+) \((\d+)\)`)
		team := home

		// performer contains team and player names, with the following structure:
		// [HomeTeam, HomePlayer, HomePlayer, HomePlayer, AwayTeam, AwayPlayer, AwayPlayer, AwayPlayer]
		// we start writing players to Home (skipping HomeTeam) and switch to away after we hit AwayTeam
		for _, p := range perf[1:] {
			if p == away.Name && team == home {
				team = away
			} else {
				matches := playerR.FindStringSubmatch(p.(string))
				player := &models.Player{Name: matches[1]}

				score, err := strconv.ParseInt(matches[2], 10, 64)

				if err == nil {
					player.Score = score
				}

				team.Players = append(team.Players, player)
			}
		}
	}

	return home, away, nil
}

func parseVenue(item *microdata.Item) string {
	return item.Properties["location"][0].(string)
}

func NewFixture(item *microdata.Item) (*models.Fixture, error) {
	fixture := &models.Fixture{
		Name:   item.Properties["description"][0].(string),
		Status: parseStatus(item),
		Venue:  item.Properties["location"][0].(string),
	}

	time, err := parseTime(item)

	if err != nil {
		return nil, err
	}

	fixture.Time = strfmt.DateTime(time)

	home, away, err := parseTeams(item, fixture.Name)

	if err != nil {
		return nil, err
	}

	fixture.Home = home
	fixture.Away = away

	return fixture, nil
}

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

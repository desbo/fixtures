package scraper

import (
	"context"
	"errors"
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

	"google.golang.org/appengine/urlfetch"
)

const BaseURL = "https://www.tabletennis365.com"

func boolToParam(b bool) string {
	if b {
		return "True"
	}

	return "False"
}

func createURL(params fixtures.ListFixturesParams) (*url.URL, error) {
	if params.League == "" || params.Season == "" {
		return nil, errors.New("please provide both a League and Season")
	}

	u, err := url.Parse(fmt.Sprintf("%s/%s/Fixtures/%s/All_Divisions", BaseURL, params.League, params.Season))

	if err != nil {
		return nil, err
	}

	setQuery := func(key string, value interface{}) {
		q := u.Query()
		q.Set(key, fmt.Sprintf("%v", value))
		u.RawQuery = q.Encode()
	}

	setQuery("vm", 1)

	if params.DivisionID != nil {
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

func safeGet(item *microdata.Item, key string) (interface{}, error) {
	values, ok := item.Properties[key]

	if !ok || len(values) == 0 {
		return nil, fmt.Errorf("property %v not found (%v)", key, item.Properties)
	}

	return values[0], nil
}

func getOrElse(item *microdata.Item, key string, orElse interface{}) interface{} {
	v, err := safeGet(item, key)

	if err != nil {
		return orElse
	}

	return v
}

func parseStatus(item *microdata.Item) string {
	s, err := safeGet(item, "eventStatus")

	if err != nil {
		// default status for events missing a status
		return "scheduled"
	}

	return strings.ToLower(strings.Replace(s.(string), "http://schema.org/Event", "", 1))
}

func parseTime(item *microdata.Item) (*time.Time, error) {
	t, err := safeGet(item, "startDate")

	if err != nil {
		return nil, err
	}

	layout := "Monday 2 January 2006 @ 15:04"
	fixed := strings.Replace(t.(string), "\u00a0", " ", -1)
	cleaned := regexp.MustCompile(`(\d{1,2})(st|nd|rd|th)`).ReplaceAllString(fixed, "$1")

	tp, err := time.Parse(layout, cleaned)

	if err != nil {
		return nil, err
	}

	return &tp, nil
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

func NewFixture(item *microdata.Item) (*models.Fixture, error) {
	fixture := &models.Fixture{
		Name:   getOrElse(item, "description", "unknown").(string),
		Status: parseStatus(item),
		Venue:  getOrElse(item, "location", "unknown").(string),
	}

	time, err := parseTime(item)

	if err != nil {
		return nil, err
	}

	fixture.Time = strfmt.DateTime(*time)
	home, away, err := parseTeams(item, fixture.Name)

	if err != nil {
		return nil, err
	}

	fixture.Home = home
	fixture.Away = away

	return fixture, nil
}

func CacheKey(params fixtures.ListFixturesParams) string {
	return fmt.Sprintf(
		"%s/%s:%d:%d:%b",
		params.League,
		params.Season,
		params.ClubID,
		params.DivisionID,
		params.ShowCompleted,
	)
}

func Scrape(ctx context.Context, params fixtures.ListFixturesParams) ([]*models.Fixture, error) {
	u, err := createURL(params)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	client := urlfetch.Client(ctx)
	resp, err := client.Get(u.String())

	if err != nil {
		return nil, err
	}

	data, err := microdata.ParseHTML(resp.Body, "text/html", u)

	if err != nil {
		return nil, err
	}

	var fixtures []*models.Fixture

	for _, item := range data.Items {
		fixture, err := NewFixture(item)
		if err == nil {
			fixtures = append(fixtures, fixture)
		}
	}

	return fixtures, nil
}

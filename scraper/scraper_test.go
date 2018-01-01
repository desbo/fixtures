package scraper

import "testing"
import "github.com/desbo/fixtures/restapi/operations/fixtures"

func TestCreateURL(t *testing.T) {
	londonDivision := int64(1)
	d := int64(5600)
	club := int64(5123)
	tru := bool(true)
	fals := bool(false)

	params1 := fixtures.ListFixturesParams{
		League:        "CentralLondon",
		Season:        "Winter_2017-18",
		ClDivision:    &londonDivision,
		ClubID:        &club,
		ShowCompleted: &tru,
	}

	params2 :=
		fixtures.ListFixturesParams{
			League:        "CentralLondon",
			Season:        "Winter_2017-18",
			DivisionID:    &d,
			ShowCompleted: &fals,
		}

	url1, err := createURL(params1)

	if err != nil {
		t.Fatal(err)
	}

	url2, err := createURL(params2)

	if err != nil {
		t.Fatal(err)
	}

	if url1.String() != "https://www.tabletennis365.com/CentralLondon/Fixtures/Winter_2017-18/All_Divisions?cl=5123&d=5596&hc=False&vm=1" {
		t.Fatalf("expected https://www.tabletennis365.com/CentralLondon/Fixtures/Winter_2017-18/All_Divisions?cl=5123&d=5596&hc=False&vm=1, got %s", url1)
	}

	if url2.String() != "https://www.tabletennis365.com/CentralLondon/Fixtures/Winter_2017-18/All_Divisions?d=5600&hc=True&vm=1" {
		t.Fatalf("expected https://www.tabletennis365.com/CentralLondon/Fixtures/Winter_2017-18/All_Divisions?d=5600&hc=True&vm=1, got %s", url1)
	}
}

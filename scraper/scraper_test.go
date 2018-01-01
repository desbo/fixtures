package scraper

import "testing"
import "github.com/desbo/fixtures/restapi/operations/fixtures"

func TestCreateURL(t *testing.T) {
	d := int64(5600)
	club := int64(5123)
	tru := bool(true)
	fals := bool(false)

	params1 := fixtures.ListFixturesParams{
		League:        "CentralLondon",
		Season:        "Winter_2017-18",
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

func TestCacheKey(t *testing.T) {
	club := int64(10)
	d := int64(20)
	tru := bool(true)

	params1 := fixtures.ListFixturesParams{
		League:        "CentralLondon",
		Season:        "Winter_2017-18",
		DivisionID:    &d,
		ClubID:        &club,
		ShowCompleted: &tru,
	}

	params2 := fixtures.ListFixturesParams{
		League: "a",
		Season: "b",
	}

	key := CacheKey(params1)
	expected := "CentralLondon/Winter_2017-18:10:20:1"

	key2 := CacheKey(params2)
	expected2 := "a/b:0:0:0"

	if key != expected {
		t.Fatalf("expected %s, got %s", expected, key)
	}

	if key2 != expected2 {
		t.Fatalf("expected %s, got %s", expected2, key2)
	}
}

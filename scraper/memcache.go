package scraper

import (
	"context"
	"encoding/json"
	"time"

	"github.com/desbo/fixtures/models"
	"github.com/desbo/fixtures/restapi/operations/fixtures"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

const Lifetime = time.Duration(24) * time.Hour

// GetOrUpdateCache gets an item from the cache or calls the get function and stores it
func GetOrUpdateCachedFixtures(ctx context.Context, params fixtures.ListFixturesParams, fixtures *[]*models.Fixture) error {
	key := CacheKey(params)
	s, err := memcache.Get(ctx, key)

	if err == nil {
		log.Infof(ctx, "cache HIT for %s", key)
		return json.Unmarshal(s.Value, &fixtures)
	}

	log.Infof(ctx, "cache MISS for %s", key)

	*fixtures, err = Scrape(ctx, params)
	b, err := json.Marshal(fixtures)

	if err == nil {
		memcache.Add(ctx, &memcache.Item{
			Key:        key,
			Value:      b,
			Expiration: Lifetime,
		})
	}

	return nil
}

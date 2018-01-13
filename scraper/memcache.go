package scraper

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/desbo/fixtures/models"
	"github.com/desbo/fixtures/restapi/operations/fixtures"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

// GetOrUpdateCache gets an item from the cache or calls the get function and stores it
func GetOrUpdateCachedFixtures(ctx context.Context, params fixtures.ListFixturesParams, fixtures *[]*models.Fixture) error {
	key := CacheKey(params)
	s, err := memcache.Get(ctx, key)
	expiry := int64(24) // default cache expiry (hours)

	if v := os.Getenv("CACHE_TIME"); v != "" {
		v, err := strconv.ParseInt(v, 10, 0)
		if err == nil {
			expiry = v
		}
	}

	if err == nil {
		log.Infof(ctx, "cache HIT for %s", key)
		return json.Unmarshal(s.Value, &fixtures)
	}

	log.Infof(ctx, "cache MISS for %s (error: %s)", key, err)

	*fixtures, err = Scrape(ctx, params)
	b, err := json.Marshal(fixtures)

	if err == nil {
		log.Infof(ctx, "adding cache entry for %s", key)
		memcache.Add(ctx, &memcache.Item{
			Key:        key,
			Value:      b,
			Expiration: time.Duration(expiry) * time.Hour,
		})
	}

	return nil
}

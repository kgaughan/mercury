package manifest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kgaughan/mercury/internal/utils"
	"github.com/kgaughan/mercury/internal/version"
	"github.com/mmcdole/gofeed"
)

// TODO if a feed is fetched, it shouldn't need to be loaded.
type cacheItem struct {
	Name         string
	UUID         string    // Used to identify the cached feed
	LastModified string    // Used for conditional GET
	ETag         string    // Also used for conditional GET
	Expires      time.Time // Date after which we should ignore the cache
}

func (ci *cacheItem) Fetch(feedURL, cacheDir string, timeout time.Duration) error {
	req, err := http.NewRequest("GET", feedURL, nil)
	if err != nil {
		return fmt.Errorf("cannot construct request: %w", err)
	}

	cacheFile := filepath.Join(cacheDir, ci.UUID+".json")
	// Blank headers if the cached feed doesn't exist
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		ci.LastModified = ""
		ci.ETag = ""
		ci.Expires = time.Time{}
	}

	// Avoid fetching stuff in the cache.
	if ci.Expires.After(time.Now()) {
		log.Printf("%s: cache not expired", feedURL)
		return nil
	}

	req = req.WithContext(context.Background())
	req.Header.Set("User-Agent", version.UserAgent())
	if ci.LastModified != "" {
		req.Header.Set("If-Modified-Since", ci.LastModified)
	}
	if ci.ETag != "" {
		req.Header.Set("If-None-Match", ci.ETag)
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot fetch feed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotModified:
		log.Printf("%s: conditional GET", feedURL)
		return nil

	case http.StatusNotFound:
	case http.StatusUnauthorized:
	case http.StatusForbidden:
		log.Printf("%s: %s, so skipping", feedURL, resp.Status)
		return nil

	case http.StatusOK:
		// Save for next time
		ci.ETag = resp.Header.Get("ETag")
		ci.LastModified = resp.Header.Get("Last-Modified")
		if err := utils.ParseCacheControlExpiration(resp.Header.Get("Cache-Control"), &ci.Expires); err != nil {
			log.Printf("Issue with %s (%v): ignoring Cache-Control", feedURL, err)
		}

		parser := gofeed.NewParser()
		if feed, err := parser.Parse(resp.Body); err != nil {
			return fmt.Errorf("can't parse %s: %w", feedURL, err)
		} else if file, err := json.Marshal(feed); err != nil {
			return fmt.Errorf("can't marshal %s: %w", feedURL, err)
		} else {
			// Save to the cache
			if err := os.WriteFile(cacheFile, file, 0o600); err != nil {
				return fmt.Errorf("can't write to cache: %w", err)
			}
		}

	default:
		// Not sure yet: do later...
		log.Fatal(resp)
	}

	return nil
}

func (ci *cacheItem) Load(cacheDir string) (*gofeed.Feed, error) {
	cacheFile := filepath.Join(cacheDir, ci.UUID+".json")
	file, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, nil
	}
	feed := &gofeed.Feed{}
	if err := json.Unmarshal(file, feed); err != nil {
		return nil, fmt.Errorf("cannot read cached feed: %w", err)
	}
	return feed, nil
}

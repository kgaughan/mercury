package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pquerna/cachecontrol/cacheobject"
)

// TODO if a feed is fetched, it shouldn't need to be loaded
type cacheItem struct {
	Name         string
	UUID         string    // Used to identify the cached feed
	LastModified string    // Used for conditional GET
	ETag         string    // Also used for conditional GET
	Expires      time.Time // Date after which we should ignore the cache
}

func (ci *cacheItem) Fetch(feedURL string, cacheDir string, timeout time.Duration) error {
	req, err := http.NewRequest("GET", feedURL, nil)
	if err != nil {
		return err
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
	req.Header.Set("User-Agent", fmt.Sprintf("planet-mercury/%v (%v)", Version, repo))
	if ci.LastModified != "" {
		req.Header.Set("If-Modified-Since", ci.LastModified)
	}
	if ci.ETag != "" {
		req.Header.Set("If-None-Match", ci.ETag)
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 304 {
		log.Printf("%s: conditional GET", feedURL)
		return nil
	}

	if resp.StatusCode == 200 {
		// Save for next time
		ci.ETag = resp.Header.Get("ETag")
		ci.LastModified = resp.Header.Get("Last-Modified")

		if resDir, err := cacheobject.ParseResponseCacheControl(resp.Header.Get("Cache-Control")); err != nil {
			log.Printf("Issue with %s (%v): ignoring Cache-Control", feedURL, err)
		} else if resDir.MaxAge > 0 {
			ci.Expires = time.Now().UTC().Add(time.Second * time.Duration(resDir.MaxAge))
		}

		parser := gofeed.NewParser()
		if feed, err := parser.Parse(resp.Body); err != nil {
			return err
		} else if file, err := json.Marshal(feed); err != nil {
			return err
		} else {
			// Save to the cache
			return ioutil.WriteFile(cacheFile, file, 0600)
		}
	}

	// Not sure yet: do later...
	log.Fatal(resp)
	return nil
}

func (ci *cacheItem) Load(cacheDir string) (*gofeed.Feed, error) {
	cacheFile := filepath.Join(cacheDir, ci.UUID+".json")
	file, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}
	feed := &gofeed.Feed{}
	if err := json.Unmarshal(file, feed); err != nil {
		return nil, err
	}
	return feed, nil
}

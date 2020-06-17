package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
)

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

var configPath = flag.String("config", "./mercury.toml", "Path to configuration")

type feed struct {
	Name string
	Feed string
}

type Config struct {
	Name         string
	URL          string `toml:url`
	Owner        string
	Email        string
	Cache        string
	Timeout      duration
	Theme        string
	Output       string
	Feed         []feed
	ItemsPerPage uint `toml:items`
	MaxPages     uint `toml:max_pages`
}

func (c *Config) Load(path string) error {
	if _, err := toml.DecodeFile(path, c); err != nil {
		return err
	}

	if configDir, err := filepath.Abs(filepath.Dir(path)); err != nil {
		return err
	} else {
		c.Cache = filepath.Join(configDir, c.Cache)
		c.Theme = filepath.Join(configDir, c.Theme)
		c.Output = filepath.Join(configDir, c.Output)
	}

	return nil
}

type cacheItem struct {
	UUID         string // Used to identify the cached feed
	LastModified string // Used for conditional GET
	ETag         string // Also used for conditional GET
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
	}

	req = req.WithContext(context.Background())
	req.Header.Set("User-Agent", "planet-mercury/1.0")
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
		return nil
	}

	if resp.StatusCode == 200 {
		// Save for next time
		ci.ETag = resp.Header.Get("ETag")
		ci.LastModified = resp.Header.Get("Last-Modified")

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
	if file, err := ioutil.ReadFile(cacheFile); err != nil {
		return nil, err
	} else {
		feed := &gofeed.Feed{}
		if err := json.Unmarshal(file, feed); err != nil {
			return nil, err
		} else {
			return feed, nil
		}
	}
}

type manifest map[string]*cacheItem

func (m *manifest) Load(path string) error {
	if file, err := ioutil.ReadFile(path); err == nil {
		if err := json.Unmarshal(file, m); err != nil {
			return err
		}
	}
	return nil
}

func (m *manifest) Save(path string) error {
	if file, err := json.Marshal(m); err == nil {
		return ioutil.WriteFile(path, file, 0600)
	} else {
		return err
	}
}

func main() {
	var config Config
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if err := config.Load(*configPath); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(config.Theme); os.IsNotExist(err) {
		log.Fatalf("Theme directory '%v' not found", config.Theme)
	}

	if fileInfo, err := os.Stat(config.Cache); os.IsNotExist(err) {
		if err := os.MkdirAll(config.Cache, 0700); err != nil {
			log.Fatal(err)
		}
	} else if !fileInfo.IsDir() {
		log.Fatalf("%s must be a directory\n", config.Cache)
	}

	manifestPath := path.Join(config.Cache, "manifest.json")
	cachedManifest := make(manifest)
	if err := cachedManifest.Load(manifestPath); err != nil {
		log.Fatal(err)
	}

	// Populate the manifest with the contents of the config file
	manifest := make(manifest)
	for _, feed := range config.Feed {
		if item, ok := cachedManifest[feed.Feed]; ok {
			// Copy over the extant cache entry
			manifest[feed.Feed] = item
		} else {
			// New feed: create a new record
			manifest[feed.Feed] = &cacheItem{
				UUID: uuid.New().String(),
			}
		}
	}

	for feedURL, cacheItem := range manifest {
		if err := cacheItem.Fetch(feedURL, config.Cache, config.Timeout.Duration); err != nil {
			log.Fatal(err)
		}
		if _, err := cacheItem.Load(config.Cache); err != nil {
			log.Fatal(err)
		}
	}

	if err := manifest.Save(manifestPath); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/google/uuid"
)

var configPath = flag.String("config", "./mercury.toml", "Path to configuration")

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

	for feedURL, item := range manifest {
		if err := item.Fetch(feedURL, config.Cache, config.Timeout.Duration); err != nil {
			log.Fatal(err)
		}
		if _, err := item.Load(config.Cache); err != nil {
			log.Fatal(err)
		}
	}

	if err := manifest.Save(manifestPath); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"container/heap"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
)

var configPath = flag.String("config", "./mercury.toml", "Path to configuration")

func ensureDir(path string) {
	if fileInfo, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatal(err)
		}
	} else if !fileInfo.IsDir() {
		log.Fatalf("%s must be a directory\n", path)
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
	tmpl, err := template.ParseFiles(path.Join(config.Theme, "index.html"))
	if err != nil {
		log.Fatal(err)
	}

	ensureDir(config.Cache)
	ensureDir(config.Output)

	manifestPath := path.Join(config.Cache, "manifest.json")
	cachedManifest := make(manifest)
	if err := cachedManifest.Load(manifestPath); err != nil {
		log.Fatal(err)
	}

	// Populate the manifest with the contents of the config file
	manifest := make(manifest)
	manifest.Populate(cachedManifest, config.Feed)
	if err := manifest.Prime(config.Cache, config.Timeout.Duration); err != nil {
		log.Fatal(err)
	}
	if err := manifest.Save(manifestPath); err != nil {
		log.Fatal(err)
	}

	// Load everything from the cache
	var fq feedQueue
	for _, item := range manifest {
		if feed, err := item.Load(config.Cache); err != nil {
			log.Fatal(err)
		} else {
			fq.AppendFeed(feed)
		}
	}

	heap.Init(&fq)
PageLoop:
	for iPage := 0; iPage < config.MaxPages; iPage++ {
		for iEntry := 0; iEntry < config.ItemsPerPage; iEntry++ {
			item := heap.Pop(&fq).(*feedAndEntry)
			if item == nil {
				break PageLoop
			}
			fmt.Println(item.Entry)
		}
	}
}

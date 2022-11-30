package main

import (
	"container/heap"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"time"

	"github.com/kgaughan/mercury/internal"
	"github.com/kgaughan/mercury/internal/flags"
	"github.com/kgaughan/mercury/internal/templates"
	"github.com/kgaughan/mercury/internal/utils"
	"github.com/mmcdole/gofeed"
)

func main() {
	var config internal.Config
	flag.Parse()

	if *flags.PrintVersion {
		fmt.Println(internal.Version)
		return
	}

	if err := config.Load(*flags.ConfigPath); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(config.Theme); os.IsNotExist(err) {
		log.Fatalf("Theme directory '%v' not found", config.Theme)
	}

	tmpl, err := templates.Configure(config.Theme)
	if err != nil {
		log.Fatal(err)
	}

	utils.EnsureDir(config.Cache)
	utils.EnsureDir(config.Output)

	manifestPath := path.Join(config.Cache, "manifest.json")
	cachedManifest := make(internal.Manifest)
	if err := cachedManifest.Load(manifestPath); err != nil {
		log.Fatal(err)
	}

	// Populate the manifest with the contents of the config file
	manifest := make(internal.Manifest)
	manifest.Populate(cachedManifest, config.Feed)
	if !*flags.NoFetch {
		manifest.Prime(config.Cache, config.Timeout.Duration)
	}
	if err := manifest.Save(manifestPath); err != nil {
		log.Fatal(err)
	}

	// Load everything from the cache
	var fq internal.FeedQueue
	var feeds []*gofeed.Feed
	for _, item := range manifest {
		if feed, err := item.Load(config.Cache); err != nil {
			log.Fatal(err)
		} else if feed == nil {
			log.Printf("Could not load cache for %q; skipping", item.Name)
		} else {
			fq.AppendFeed(feed)
			feeds = append(feeds, feed)
		}
	}

	now := time.Now()

	heap.Init(&fq)
	for iPage := 0; iPage < config.MaxPages; iPage++ {
		var pageName string
		if iPage == 0 {
			pageName = "index.html"
		} else {
			pageName = fmt.Sprintf("index%d.html", iPage)
		}

		lastPage := false
		var items []*internal.FeedEntry
		for iEntry := 0; iEntry < config.ItemsPerPage; iEntry++ {
			item := fq.Top()
			if item == nil {
				lastPage = true
				break
			}
			items = append(items, item.(*internal.FeedEntry))
			heap.Fix(&fq, 0)
		}

		f, err := os.Create(path.Join(config.Output, pageName))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		vars := struct {
			Generator string
			Name      string
			URL       template.URL
			Owner     string
			Email     string
			PageNo    int
			Items     []*internal.FeedEntry
			Generated time.Time
			Feeds     []*gofeed.Feed
		}{
			Generator: internal.Generator(),
			Name:      config.Name,
			URL:       template.URL(config.URL),
			Owner:     config.Owner,
			Email:     config.Email,
			PageNo:    iPage + 1,
			Items:     items,
			Generated: now,
			Feeds:     feeds,
		}
		if err := tmpl.ExecuteTemplate(f, "index.html", vars); err != nil {
			log.Fatal(err)
		}
		if lastPage {
			break
		}
	}

	// Generate OPML
	opml := internal.NewOPML(len(feeds))
	for url, item := range manifest {
		opml.Append(item.Name, url)
	}
	if err := opml.MarshalToFile(path.Join(config.Output, "opml.xml")); err != nil {
		log.Fatal(err)
	}
}

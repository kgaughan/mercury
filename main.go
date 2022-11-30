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

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
)

const repo = "https://github.com/kgaughan/mercury/"

// Version contains the version (set during build)
var Version string

var printVersion = flag.Bool("version", false, "Print version and exit")
var configPath = flag.String("config", "./mercury.toml", "Path to configuration")
var noFetch = flag.Bool("no-fetch", false, "Don't fetch, just use what's in the cache")

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
		out := flag.CommandLine.Output()
		name := path.Base(os.Args[0])
		fmt.Fprintf(out, "%s - Generates an aggregated site from a set of feeds.\n\n", name)
		fmt.Fprintln(out, "Usage:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *printVersion {
		fmt.Println(Version)
		return
	}

	if err := config.Load(*configPath); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(config.Theme); os.IsNotExist(err) {
		log.Fatalf("Theme directory '%v' not found", config.Theme)
	}

	// This is just a starting point so there's a reasonable policy
	p := bluemonday.UGCPolicy()

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"isodatefmt": func(t time.Time) string {
			return t.Format(time.RFC3339)
		},
		"datefmt": func(fmt string, t time.Time) string {
			return t.Format(fmt)
		},
		"safe": func(text string) template.HTML {
			return template.HTML(text)
		},
		"sanitize": func(text template.HTML) template.HTML {
			return template.HTML(p.Sanitize(string(text)))
		},
		"excerpt": func(max int, text template.HTML) template.HTML {
			return template.HTML(excerpt(string(text), max))
		},
	}).ParseFiles(path.Join(config.Theme, "index.html"))
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
	if !*noFetch {
		manifest.Prime(config.Cache, config.Timeout.Duration)
	}
	if err := manifest.Save(manifestPath); err != nil {
		log.Fatal(err)
	}

	// Load everything from the cache
	var fq feedQueue
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
		var items []*FeedEntry
		for iEntry := 0; iEntry < config.ItemsPerPage; iEntry++ {
			item := fq.Top()
			if item == nil {
				lastPage = true
				break
			}
			items = append(items, item.(*FeedEntry))
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
			Items     []*FeedEntry
			Generated time.Time
			Feeds     []*gofeed.Feed
		}{
			Generator: fmt.Sprintf("Planet Mercury %v (%v)", Version, repo),
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
	opml := NewOPML(len(feeds))
	for url, item := range manifest {
		opml.Append(item.Name, url)
	}
	if err := opml.MarshalToFile(path.Join(config.Output, "opml.xml")); err != nil {
		log.Fatal(err)
	}
}

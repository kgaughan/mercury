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
)

const REPO = "https://github.com/kgaughan/mercury/"

var Version string

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

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"isodatefmt": func(t time.Time) string {
			return t.Format(time.RFC3339)
		},
		"datefmt": func(fmt string, t time.Time) string {
			return t.Format(fmt)
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
		var items []*feedEntry
		for iEntry := 0; iEntry < config.ItemsPerPage; iEntry++ {
			item := fq.Top()
			if item == nil {
				lastPage = true
				break
			}
			items = append(items, item.(*feedEntry))
			heap.Fix(&fq, 0)
		}

		f, err := os.OpenFile(path.Join(config.Output, pageName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
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
			Items     []*feedEntry
			Generated time.Time
		}{
			Generator: fmt.Sprintf("Planet Mercury %v (%v)", Version, REPO),
			Name:      config.Name,
			URL:       template.URL(config.URL),
			Owner:     config.Owner,
			Email:     config.Email,
			PageNo:    iPage + 1,
			Items:     items,
			Generated: now,
		}
		if err := tmpl.ExecuteTemplate(f, "index.html", vars); err != nil {
			log.Fatal(err)
		}
		if lastPage {
			break
		}
	}
}

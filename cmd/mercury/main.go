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
	"github.com/kgaughan/mercury/internal/feed"
	"github.com/kgaughan/mercury/internal/flags"
	"github.com/kgaughan/mercury/internal/manifest"
	"github.com/kgaughan/mercury/internal/opml"
	"github.com/kgaughan/mercury/internal/templates"
	"github.com/kgaughan/mercury/internal/theme"
	"github.com/kgaughan/mercury/internal/utils"
	"github.com/kgaughan/mercury/internal/version"
	"github.com/mmcdole/gofeed"
)

func main() {
	var config internal.Config
	flag.Parse()

	if *flags.PrintVersion {
		fmt.Println(version.Version)
		return
	}

	if err := config.Load(*flags.ConfigPath); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(config.Theme); os.IsNotExist(err) {
		log.Fatalf("Theme directory '%v' not found", config.Theme)
	}

	var themeConfig theme.Config
	if err := themeConfig.Load(config.Theme); err != nil {
		log.Fatal(err)
	}

	tmpl, err := templates.Configure(config.Theme)
	if err != nil {
		log.Fatal(err)
	}

	utils.EnsureDir(config.Cache)

	if !*flags.NoBuild {
		utils.EnsureDir(config.Output)
		if err := themeConfig.CopyTo(config.Output); err != nil {
			log.Fatal(err)
		}
	}

	manifestPath := path.Join(config.Cache, "manifest.json")
	manifest, err := manifest.LoadManifest(manifestPath)
	if err != nil {
		log.Fatal(err)
	}

	// Populate the manifest with the contents of the config file
	manifest.Populate(config.Feeds)
	if !*flags.NoFetch {
		manifest.Prime(config.Cache, config.Timeout.Duration)
	}
	if err := manifest.Save(manifestPath); err != nil {
		log.Fatal(err)
	}

	fq, feeds, err := populate(manifest, config.Cache)
	if err != nil {
		log.Fatal(err)
	}

	if !*flags.NoBuild {
		if err := writePages(fq, feeds, config, tmpl); err != nil {
			log.Fatal(err)
		}

		if err := writeOPML(manifest, path.Join(config.Output, "opml.xml")); err != nil {
			log.Fatal(err)
		}
	}
}

func populate(manifest *manifest.Manifest, cache string) (*feed.Queue, []*gofeed.Feed, error) {
	fq := &feed.Queue{}
	var feeds []*gofeed.Feed
	for _, item := range *manifest {
		if feed, err := item.Load(cache); err != nil {
			return nil, nil, fmt.Errorf("could not load feed from cache: %w", err)
		} else if feed == nil {
			log.Printf("Could not load cache for %q; skipping", item.Name)
		} else {
			fq.Append(feed)
			feeds = append(feeds, feed)
		}
	}
	return fq, feeds, nil
}

func writePages(fq *feed.Queue, feeds []*gofeed.Feed, config internal.Config, tmpl *template.Template) error {
	now := time.Now()

	heap.Init(fq)
	for iPage := 0; iPage < config.MaxPages; iPage++ {
		var pageName string
		if iPage == 0 {
			pageName = "index.html"
		} else {
			pageName = fmt.Sprintf("index%d.html", iPage)
		}

		lastPage := false
		var items []*feed.Entry
		for iEntry := 0; iEntry < config.ItemsPerPage; iEntry++ {
			item := fq.Top()
			if item == nil {
				lastPage = true
				break
			}
			items = append(items, item.(*feed.Entry))
			heap.Fix(fq, 0)
		}

		f, err := os.Create(path.Join(config.Output, pageName))
		if err != nil {
			return fmt.Errorf("could not create page: %w", err)
		}
		defer f.Close()

		vars := struct {
			Generator string
			Name      string
			URL       template.URL
			Owner     string
			Email     string
			PageNo    int
			Items     []*feed.Entry
			Generated time.Time
			Feeds     []*gofeed.Feed
		}{
			Generator: version.Generator(),
			Name:      config.Name,
			URL:       template.URL(config.URL), //nolint:gosec
			Owner:     config.Owner,
			Email:     config.Email,
			PageNo:    iPage + 1,
			Items:     items,
			Generated: now,
			Feeds:     feeds,
		}
		if err := tmpl.ExecuteTemplate(f, "index.html", vars); err != nil {
			return fmt.Errorf("could not render template: %w", err)
		}
		if lastPage {
			break
		}
	}
	return nil
}

func writeOPML(manifest *manifest.Manifest, path string) error {
	opml := opml.New(manifest.Len())
	for url, item := range *manifest {
		opml.Append(item.Name, url)
	}
	if err := utils.MarshalToFile(path, opml); err != nil {
		return fmt.Errorf("can't write %s: %w", path, err)
	}
	return nil
}

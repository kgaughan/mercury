package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"time"

	"github.com/kgaughan/mercury/internal"
	"github.com/kgaughan/mercury/internal/atom"
	"github.com/kgaughan/mercury/internal/feed"
	"github.com/kgaughan/mercury/internal/flags"
	"github.com/kgaughan/mercury/internal/manifest"
	"github.com/kgaughan/mercury/internal/templates"
	"github.com/kgaughan/mercury/internal/theme"
	"github.com/kgaughan/mercury/internal/utils"
	"github.com/kgaughan/mercury/internal/version"
	"github.com/mmcdole/gofeed"
	flag "github.com/spf13/pflag"
)

var errNoEntries = errors.New("no entries to write to feed")

func main() {
	var config internal.Config
	flag.Parse()

	if *flags.PrintVersion {
		fmt.Println(version.Version)
		return
	}
	if *flags.ShowHelp {
		flag.Usage()
		os.Exit(0)
	}

	if err := config.Load(*flags.ConfigPath); err != nil {
		log.Fatal(err)
	}

	var themeConfig theme.Config
	if err := themeConfig.Load(config.Theme); err != nil {
		log.Fatal(err)
	}

	tmpl, err := templates.Configure(config.Theme)
	if err != nil {
		log.Fatal(err)
	}

	utils.EnsureCache(config.Cache)

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
		manifest.Prime(config.Cache, config.Timeout.Duration, config.Parallelism, config.JobQueueDepth)
	}
	if err := manifest.Save(manifestPath); err != nil {
		log.Fatal(err)
	}

	fq, feeds, err := populate(manifest, config.Cache)
	if err != nil {
		log.Fatal(err)
	}

	if !*flags.NoBuild {
		entries := fq.Shuffle(config.ItemsPerPage * config.MaxPages)
		if err := writePages(entries, feeds, config, tmpl); err != nil {
			log.Fatal(err)
		}

		if err := writeFeed(entries, config); err != nil {
			log.Fatal(err)
		}

		if err := manifest.AsOPML().Save(path.Join(config.Output, "opml.xml")); err != nil {
			log.Fatal(err)
		}
	}
}

// populate loads cached feeds from the manifest.
func populate(manifest *manifest.Manifest, cache string) (*feed.Queue, []*gofeed.Feed, error) {
	fq := &feed.Queue{}
	var feeds []*gofeed.Feed
	for _, item := range *manifest {
		if feed, err := item.Load(cache); err != nil {
			return nil, nil, fmt.Errorf("cannot load feed from cache: %w", err)
		} else if feed == nil {
			log.Printf("Could not load cache for %q; skipping", item.Name)
		} else {
			fq.Append(feed)
			feeds = append(feeds, feed)
		}
	}
	return fq, feeds, nil
}

// writePages generates paginated HTML pages from the feed entries.
func writePages(entries []*feed.Entry, feeds []*gofeed.Feed, config internal.Config, tmpl *template.Template) error {
	now := time.Now()

	nPages := len(entries) / config.ItemsPerPage
	for iPage := range nPages {
		offset := iPage * config.ItemsPerPage
		var pageName string
		if iPage == 0 {
			pageName = "index.html"
		} else {
			pageName = fmt.Sprintf("index%d.html", iPage)
		}

		f, err := os.Create(path.Join(config.Output, pageName))
		if err != nil {
			return fmt.Errorf("cannot create page: %w", err)
		}
		defer f.Close()

		var prevPage, nextPage string
		if iPage == 1 {
			prevPage = "index.html"
		} else if iPage > 1 {
			prevPage = fmt.Sprintf("index%d.html", iPage-1)
		}
		if iPage < nPages-1 {
			nextPage = fmt.Sprintf("index%d.html", iPage+1)
		}

		end := min(offset+config.ItemsPerPage, len(entries))

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
			PrevPage  string
			NextPage  string
		}{
			Generator: version.Generator(),
			Name:      config.Name,
			URL:       template.URL(config.URL), //nolint:gosec
			Owner:     config.Owner,
			Email:     config.Email,
			PageNo:    iPage + 1,
			Items:     entries[offset:end],
			Generated: now,
			Feeds:     feeds,
			PrevPage:  prevPage,
			NextPage:  nextPage,
		}
		if err := tmpl.ExecuteTemplate(f, "index.html", vars); err != nil {
			return fmt.Errorf("cannot render template: %w", err)
		}
	}
	return nil
}

// writeFeed generates an Atom feed from the entries.
func writeFeed(entries []*feed.Entry, config internal.Config) error {
	if len(entries) == 0 {
		return errNoEntries
	}
	feed := atom.Feed{
		Title:   config.Name,
		ID:      config.FeedID,
		Updated: atom.Time(*entries[0].Updated()),
		Links: []atom.Link{{
			Rel:  "self",
			Href: config.URL + "feed.atom",
		}},
	}

	for _, entry := range entries {
		summary := &atom.Text{
			Type: "html",
			Body: string(entry.Summary()),
		}
		if summary.Body == "" {
			summary = nil
		}
		atomEntry := &atom.Entry{
			Title: entry.Title(),
			ID:    entry.ID(),
			Links: []atom.Link{{
				Rel:  "alternate",
				Href: string(entry.Link()),
			}},
			Published: atom.Time(*entry.Published()),
			Updated:   atom.Time(*entry.Updated()),
			Author: &atom.Person{
				Name: entry.Author(),
			},
			Summary: summary,
			Content: &atom.Text{
				Type: "html",
				Body: string(entry.Content()),
			},
		}
		feed.Entries = append(feed.Entries, atomEntry)
	}

	if err := feed.Save(path.Join(config.Output, "atom.xml")); err != nil {
		return fmt.Errorf("cannot create Atom feed: %w", err)
	}
	return nil
}

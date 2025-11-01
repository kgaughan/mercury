package manifest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kgaughan/mercury/internal/opml"
)

// Manifest maps feed URLs to their cache metadata.
type Manifest map[string]*cacheEntry

// fetchJob represents a job to fetch a feed and update its cache item.
type fetchJob struct {
	URL  string
	Item *cacheEntry
}

// LoadManifest loads the manifest from a file.
func LoadManifest(path string) (*Manifest, error) {
	manifest := &Manifest{}
	file, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// No manifest yet: return empty one
			return manifest, nil
		}
		return nil, fmt.Errorf("cannot read manifest: %w", err)
	}
	if err := json.Unmarshal(file, manifest); err != nil {
		return nil, fmt.Errorf("cannot unmarshal manifest: %w", err)
	}
	return manifest, nil
}

// Populate adds feeds to the manifest.
func (m *Manifest) Populate(feeds []Feed) {
	// The value is a dummy: we're just using the map as a set
	liveFeedUrls := make(map[string]struct{})
	for _, feed := range feeds {
		liveFeedUrls[feed.Feed] = struct{}{}
		if _, ok := (*m)[feed.Feed]; !ok {
			// New feed: create a new record
			(*m)[feed.Feed] = &cacheEntry{
				Name: feed.Name,
				UUID: uuid.New().String(),
			}
		}
	}
	// Remove any feeds no longer in the config
	for url := range *m {
		if _, ok := liveFeedUrls[url]; !ok {
			log.Printf("Removing feed %q from manifest", url)
			delete(*m, url)
		}
	}
}

// Len returns the number of feeds in the manifest.
func (m *Manifest) Len() int {
	return len(*m)
}

// Save writes the manifest to a file.
func (m *Manifest) Save(path string) error {
	file, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("can't marshal manifest: %w", err)
	}
	if err = os.WriteFile(path, file, 0o600); err != nil {
		return fmt.Errorf("can't save manifest: %w", err)
	}
	return nil
}

// AsOPML converts the manifest to an OPML document.
func (m *Manifest) AsOPML() *opml.OPML {
	opml := opml.New(m.Len())
	for url, item := range *m {
		opml.Append(item.Name, url)
	}
	return opml
}

// Prime fetches and caches all feeds in the manifest concurrently.
func (m *Manifest) Prime(cache string, timeout time.Duration, parallelism, jobQueueDepth int) {
	var wg sync.WaitGroup
	jobs := make(chan *fetchJob, jobQueueDepth)

	log.Printf("Priming manifest with %d feeds using %d workers, with a queue depth of %d", len(*m), parallelism, jobQueueDepth)
	for range parallelism {
		wg.Add(1)
		go func() {
			ctx := context.Background()
			defer wg.Done()
			for job := range jobs {
				if job == nil || job.Item == nil {
					continue
				}
				log.Print("Fetching ", job.URL)
				if err := job.Item.Fetch(ctx, job.URL, cache, timeout); err != nil {
					log.Print(err)
				}
			}
		}()
	}

	for feedURL, item := range *m {
		if item != nil {
			jobs <- &fetchJob{URL: feedURL, Item: item}
		}
	}

	close(jobs)
	wg.Wait()
}

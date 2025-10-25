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

type Manifest map[string]*cacheItem

type fetchJob struct {
	URL  string
	Item *cacheItem
}

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

func (m *Manifest) Populate(feeds []Feed) {
	for _, feed := range feeds {
		if _, ok := (*m)[feed.Feed]; !ok {
			// New feed: create a new record
			(*m)[feed.Feed] = &cacheItem{
				Name: feed.Name,
				UUID: uuid.New().String(),
			}
		}
	}
}

func (m *Manifest) Len() int {
	return len(*m)
}

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

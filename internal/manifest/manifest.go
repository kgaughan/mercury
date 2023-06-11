package manifest

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Manifest map[string]*cacheItem

type fetchJob struct {
	URL  string
	Item *cacheItem
}

func LoadManifest(path string) (*Manifest, error) {
	manifest := &Manifest{}
	if file, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(file, manifest); err != nil {
			return nil, fmt.Errorf("could not load manifest: %w", err)
		}
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

func (m *Manifest) Prime(cache string, timeout time.Duration) {
	var wg sync.WaitGroup

	// The channel depth is kind of arbitrary.
	jobs := make(chan *fetchJob, 2*runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				if err := job.Item.Fetch(job.URL, cache, timeout); err != nil {
					log.Print(err)
				}
			}
		}()
	}

	for feedURL, item := range *m {
		jobs <- &fetchJob{
			URL:  feedURL,
			Item: item,
		}
	}
	close(jobs)
	wg.Wait()
}

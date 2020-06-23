package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
)

type manifest map[string]*cacheItem

func (m *manifest) Load(path string) error {
	if file, err := ioutil.ReadFile(path); err == nil {
		if err := json.Unmarshal(file, m); err != nil {
			return err
		}
	}
	return nil
}

func (m *manifest) Save(path string) error {
	if file, err := json.Marshal(m); err == nil {
		return ioutil.WriteFile(path, file, 0600)
	} else {
		return err
	}
}

func (m manifest) Populate(cache manifest, feeds []feed) {
	for _, feed := range feeds {
		if item, ok := cache[feed.Feed]; ok {
			// Copy over the extant cache entry
			m[feed.Feed] = item
		} else {
			// New feed: create a new record
			m[feed.Feed] = &cacheItem{
				UUID: uuid.New().String(),
			}
		}
	}
}

type fetchJob struct {
	URL  string
	Item *cacheItem
}

func (m manifest) Prime(cache string, timeout time.Duration) {
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

	for feedURL, item := range m {
		jobs <- &fetchJob{
			URL:  feedURL,
			Item: item,
		}
	}
	close(jobs)
	wg.Wait()
}

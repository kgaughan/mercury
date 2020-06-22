package main

import (
	"encoding/json"
	"io/ioutil"
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

func (m manifest) Prime(cache string, timeout time.Duration) error {
	for feedURL, item := range m {
		// TODO these can be fetched in parallel
		if err := item.Fetch(feedURL, cache, timeout); err != nil {
			return err
		}
	}
	return nil
}

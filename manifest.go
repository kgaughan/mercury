package main

import (
	"encoding/json"
	"io/ioutil"
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

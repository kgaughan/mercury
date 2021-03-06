package main

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type feed struct {
	Name string
	Feed string
}

// Config describes our configuration
type Config struct {
	Name         string
	URL          string `toml:"url"`
	Owner        string
	Email        string
	Cache        string
	Timeout      duration
	Theme        string
	Output       string
	Feed         []feed
	ItemsPerPage int `toml:"items"`
	MaxPages     int `toml:"max_pages"`
}

// Load loads our configuration file
func (c *Config) Load(path string) error {
	if _, err := toml.DecodeFile(path, c); err != nil {
		return err
	}
	configDir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return err
	}
	c.Cache = filepath.Join(configDir, c.Cache)
	c.Theme = filepath.Join(configDir, c.Theme)
	c.Output = filepath.Join(configDir, c.Output)

	return nil
}

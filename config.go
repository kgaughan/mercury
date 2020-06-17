package main

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type feed struct {
	Name string
	Feed string
}

type Config struct {
	Name         string
	URL          string `toml:url`
	Owner        string
	Email        string
	Cache        string
	Timeout      duration
	Theme        string
	Output       string
	Feed         []feed
	ItemsPerPage uint `toml:items`
	MaxPages     uint `toml:max_pages`
}

func (c *Config) Load(path string) error {
	if _, err := toml.DecodeFile(path, c); err != nil {
		return err
	}

	if configDir, err := filepath.Abs(filepath.Dir(path)); err != nil {
		return err
	} else {
		c.Cache = filepath.Join(configDir, c.Cache)
		c.Theme = filepath.Join(configDir, c.Theme)
		c.Output = filepath.Join(configDir, c.Output)
	}

	return nil
}

package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/kgaughan/mercury/internal/manifest"
	dflt "github.com/kgaughan/mercury/internal/theme/default"
	"github.com/kgaughan/mercury/internal/utils"
)

// Config describes our configuration.
type Config struct {
	Name         string          `toml:"name"`
	URL          string          `toml:"url"`
	Owner        string          `toml:"owner"`
	Email        string          `toml:"email"`
	FeedID       string          `toml:"feed_id"`
	Cache        string          `toml:"cache"`
	Timeout      utils.Duration  `toml:"timeout"`
	theme        string          `toml:"theme"`
	Output       string          `toml:"output"`
	Feeds        []manifest.Feed `toml:"feed"`
	ItemsPerPage int             `toml:"items"`
	MaxPages     int             `toml:"max_pages"`
}

func (c Config) GetThemeFS() fs.FS {
	if c.theme == "" {
		return dflt.Theme
	}
	return os.DirFS(c.theme)
}

// Load loads our configuration file.
func (c *Config) Load(path string) error {
	c.Name = "Planet"
	c.Cache = "./cache"
	c.theme = ""
	c.Output = "./output"
	c.ItemsPerPage = 10
	c.MaxPages = 5

	if _, err := toml.DecodeFile(path, c); err != nil {
		return fmt.Errorf("cannot load configuration: %w", err)
	}

	configDir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("cannot normalize configuration path: %w", err)
	}

	c.Cache = filepath.Join(configDir, c.Cache)
	if c.theme != "" {
		c.theme = filepath.Join(configDir, c.theme)
		if _, err := os.Stat(c.theme); os.IsNotExist(err) {
			return fmt.Errorf("theme %q not found: %w", c.theme, err)
		}
	}
	c.Output = filepath.Join(configDir, c.Output)
	return nil
}

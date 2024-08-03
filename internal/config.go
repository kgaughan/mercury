package internal

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/kgaughan/mercury/internal/manifest"
	dflt "github.com/kgaughan/mercury/internal/theme/default"
	"github.com/kgaughan/mercury/internal/utils"
)

const cpuLimit = 32 // Cap on the number of CPUs/cores to use

// Config describes our configuration.
type Config struct {
	Name          string          `toml:"name"`
	URL           string          `toml:"url"`
	Owner         string          `toml:"owner"`
	Email         string          `toml:"email"`
	FeedID        string          `toml:"feed_id"`
	Cache         string          `toml:"cache"`
	Timeout       utils.Duration  `toml:"timeout"`
	ThemePath     string          `toml:"theme"`
	Theme         fs.FS           `toml:"-"`
	Output        string          `toml:"output"`
	Feeds         []manifest.Feed `toml:"feed"`
	ItemsPerPage  int             `toml:"items"`
	MaxPages      int             `toml:"max_pages"`
	JobQueueDepth int             `toml:"job_queue_depth"`
	Parallelism   int             `toml:"parallelism"`
}

// Load loads our configuration file.
func (c *Config) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open configuration %q: %w", path, err)
	}
	defer file.Close()
	return c.LoadFromReader(file, path)
}

// LoadFromReader loads our configuration file from an arbitrary source.
func (c *Config) LoadFromReader(r io.Reader, path string) error {
	c.Name = "Planet"
	c.Cache = "./cache"
	c.ThemePath = ""
	c.Output = "./output"
	c.ItemsPerPage = 10
	c.MaxPages = 5
	c.JobQueueDepth = 0
	c.Parallelism = runtime.NumCPU()

	if _, err := toml.NewDecoder(r).Decode(c); err != nil {
		return fmt.Errorf("cannot load configuration: %w", err)
	}

	configDir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("cannot normalize configuration path: %w", err)
	}

	// Enforce some sensible lower bounds on feed fetching parallelism
	c.Parallelism = min(max(1, c.Parallelism), cpuLimit)
	c.JobQueueDepth = max(2*c.Parallelism, c.JobQueueDepth)

	// Normalize all the other paths
	if !filepath.IsAbs(c.Cache) {
		c.Cache = filepath.Join(configDir, c.Cache)
	}
	if c.ThemePath == "" {
		c.Theme = dflt.Theme
	} else {
		var themePath string
		if filepath.IsAbs(c.ThemePath) {
			themePath = c.ThemePath
		} else {
			themePath = filepath.Join(configDir, c.ThemePath)
		}
		if _, err := os.Stat(themePath); os.IsNotExist(err) {
			return fmt.Errorf("theme %q not found: %w", themePath, err)
		}
		c.Theme = os.DirFS(themePath)
	}
	if !filepath.IsAbs(c.Output) {
		c.Output = filepath.Join(configDir, c.Output)
	}

	// Compile all the filters
	for _, feed := range c.Feeds {
		for _, filter := range feed.Filters {
			if err := filter.Compile(); err != nil {
				return fmt.Errorf("issue with filter for %q: %w", feed.Name, err)
			}
		}
	}

	return nil
}

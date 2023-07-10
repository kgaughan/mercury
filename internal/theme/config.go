package theme

import (
	"fmt"
	"io/fs"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/kgaughan/mercury/internal/utils"
)

// Config describes the configuration of a theme configuration.
type Config struct {
	root fs.FS
	Name string     `toml:"name"`
	BOM  []BOMEntry `toml:"bom"`
}

type BOMEntry struct {
	Path string `toml:"path"`
}

// Load loads our configuration file.
func (c *Config) Load(themeFS fs.FS) error {
	c.root = themeFS
	if _, err := toml.DecodeFS(themeFS, "theme.toml", c); err != nil {
		return fmt.Errorf("cannot load theme configuration: %w", err)
	}
	return nil
}

func (c *Config) CopyTo(destDir string) error {
	for _, entry := range c.BOM {
		src, err := c.root.Open(entry.Path)
		if err != nil {
			return fmt.Errorf("cannot read %q from theme: %w", entry.Path, err)
		}
		defer src.Close()
		if err = utils.Copy(src, path.Join(destDir, entry.Path)); err != nil {
			return fmt.Errorf("Failed to copy %q into %q: %w", entry.Path, destDir, err)
		}
	}
	return nil
}

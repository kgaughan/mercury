package theme

import (
	"fmt"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/kgaughan/mercury/internal/utils"
)

// Config describes the configuration of a theme configuration.
type Config struct {
	root string
	Name string     `toml:"name"`
	BOM  []BOMEntry `toml:"bom"`
}

type BOMEntry struct {
	Path string `toml:"path"`
}

// Load loads our configuration file.
func (c *Config) Load(themeDir string) error {
	c.root = themeDir
	if _, err := toml.DecodeFile(path.Join(themeDir, "theme.toml"), c); err != nil {
		return fmt.Errorf("could not load theme configuration: %w", err)
	}
	return nil
}

func (c *Config) CopyTo(destDir string) error {
	for _, entry := range c.BOM {
		if err := utils.Copy(path.Join(c.root, entry.Path), path.Join(destDir, entry.Path)); err != nil {
			return fmt.Errorf("Failed to copy %q into %q: %w", entry.Path, destDir)
		}
	}
	return nil
}

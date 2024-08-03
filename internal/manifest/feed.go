package manifest

import "github.com/kgaughan/mercury/internal/filters"

// Feed represents a feed entry in the manifest.
type Feed struct {
	Name    string           `toml:"name"`
	Feed    string           `toml:"feed"`
	Filters []filters.Filter `toml:"filters"`
}

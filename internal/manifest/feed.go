package manifest

import "github.com/kgaughan/mercury/internal/filters"

// Feed represents a feed entry in the manifest.
type Feed struct {
	Name    string
	Feed    string
	Filters []filters.Filter
}

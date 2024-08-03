package manifest

import "github.com/kgaughan/mercury/internal/filters"

type Feed struct {
	Name    string
	Feed    string
	Filters []filters.Filter
}

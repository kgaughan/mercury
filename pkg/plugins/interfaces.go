// Package interfaces provides interfaces and utilities for plugin
// implementors to implement and used.
package plugins

import (
	"html/template"

	"github.com/mmcdole/gofeed"
)

// Templating defines the methods plugins wishing to extend the template
// function map should implement.
type Templating interface {
	Funcs() template.FuncMap
}

// Filter defined the methods plugins wishing to perform filtering on feeds
// should implement.
type Filter interface {
	// Not super sure about this. It can modify the item in place if need
	// be, I think. The return values are a boolean indicating whether the
	// item should be skipped or not (assuming there's no error), and an
	// error.
	FilterItem(item *gofeed.Item) (bool, error)
}

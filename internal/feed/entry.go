package feed

import (
	"html/template"
	"time"

	"github.com/mmcdole/gofeed"
)

// Entry describes an entry in a feed in a form suitable for templating.
type Entry struct {
	feed  *gofeed.Feed
	entry *gofeed.Item
}

// NewEntry creates a new feed entry object.
func NewEntry(feed *gofeed.Feed, entry *gofeed.Item) *Entry {
	return &Entry{feed: feed, entry: entry}
}

// FeedName returns the name of the feed.
func (e Entry) FeedName() string {
	return e.feed.Title
}

// SiteLink returns the site URL in a form usable in a template.
func (e Entry) SiteLink() template.URL {
	return template.URL(e.feed.Link) //nolint:gosec
}

// FeedLink returns the feed URL in a form usable in a template.
func (e Entry) FeedLink() template.URL {
	return template.URL(e.feed.FeedLink) //nolint:gosec
}

// FeedPublished returns a best guess at the correct feed publication date.
func (e Entry) FeedPublished() *time.Time {
	if e.feed.PublishedParsed != nil {
		return e.feed.PublishedParsed
	}
	// Fallback
	return e.feed.UpdatedParsed
}

// FeedUpdated returns when the feed was updated.
func (e Entry) FeedUpdated() *time.Time {
	return e.feed.UpdatedParsed
}

// FeedDescription returns the description of the feed, if any.
func (e Entry) FeedDescription() string {
	return e.feed.Description
}

// Title returns the entry title.
func (e Entry) Title() string {
	return e.entry.Title
}

// Summary returns the entry summary.
func (e Entry) Summary() template.HTML {
	return template.HTML(e.entry.Description) //nolint:gosec
}

// Content returns the full content of the entry.
func (e Entry) Content() template.HTML {
	return template.HTML(e.entry.Content) //nolint:gosec
}

// Link returns a link to the entry.
func (e Entry) Link() template.URL {
	return template.URL(e.entry.Link) //nolint:gosec
}

// Author returns the entry's author, and if none, returns the feed's author.
func (e Entry) Author() string {
	if e.entry.Author != nil {
		return e.entry.Author.Name
	}
	if e.feed.Author != nil {
		return e.feed.Author.Name
	}
	return ""
}

// Published returns a best guess at the entry's publication date.
func (e Entry) Published() *time.Time {
	if e.entry.PublishedParsed != nil {
		return e.entry.PublishedParsed
	}
	// Fallback
	return e.entry.UpdatedParsed
}

// Updated returns the update date of the entry.
func (e Entry) Updated() *time.Time {
	if e.entry.UpdatedParsed != nil {
		return e.entry.UpdatedParsed
	}
	// Fallback
	return e.entry.PublishedParsed
}

// Categories returns the entry's categories.
func (e Entry) Categories() []string {
	return e.entry.Categories
}

func (e Entry) ID() string {
	if e.entry.GUID != "" {
		return e.entry.GUID
	}
	return e.entry.Link
}

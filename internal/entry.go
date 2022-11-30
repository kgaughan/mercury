package internal

import (
	"html/template"
	"time"

	"github.com/mmcdole/gofeed"
)

// FeedEntry describes an entry in a feed in a form suitable for templating
type FeedEntry struct {
	feed  *gofeed.Feed
	entry *gofeed.Item
}

// NewEntry creates a new feed entry object
func NewEntry(feed *gofeed.Feed, entry *gofeed.Item) *FeedEntry {
	return &FeedEntry{feed: feed, entry: entry}
}

// FeedName returns the name of the feed
func (e FeedEntry) FeedName() string {
	return e.feed.Title
}

// SiteLink returns the site URL in a form usable in a template
func (e FeedEntry) SiteLink() template.URL {
	return template.URL(e.feed.Link)
}

// FeedLink returns the feed URL in a form usable in a template
func (e FeedEntry) FeedLink() template.URL {
	return template.URL(e.feed.FeedLink)
}

// FeedPublished returns a best guess at the correct feed publication date
func (e FeedEntry) FeedPublished() *time.Time {
	if e.feed.PublishedParsed != nil {
		return e.feed.PublishedParsed
	}
	// Fallback
	return e.feed.UpdatedParsed
}

// FeedUpdated returns when the feed was updated
func (e FeedEntry) FeedUpdated() *time.Time {
	return e.feed.UpdatedParsed
}

// FeedDescription returns the description of the feed, if any
func (e FeedEntry) FeedDescription() string {
	return e.feed.Description
}

// Title returns the entry title
func (e FeedEntry) Title() string {
	return e.entry.Title
}

// Summary returns the entry summary
func (e FeedEntry) Summary() template.HTML {
	return template.HTML(e.entry.Description)
}

// Content returns the full content of the entry
func (e FeedEntry) Content() template.HTML {
	return template.HTML(e.entry.Content)
}

// Link returns a link to the entry
func (e FeedEntry) Link() template.URL {
	return template.URL(e.entry.Link)
}

// Author returns the entry's author, and if none, returns the feed's author
func (e FeedEntry) Author() string {
	if e.entry.Author != nil {
		return e.entry.Author.Name
	}
	if e.feed.Author != nil {
		return e.feed.Author.Name
	}
	return ""
}

// Published returns a best guess at the entry's publication date
func (e FeedEntry) Published() *time.Time {
	if e.entry.PublishedParsed != nil {
		return e.entry.PublishedParsed
	}
	// Fallback
	return e.entry.UpdatedParsed
}

// Updated returns the update date of the entry
func (e FeedEntry) Updated() *time.Time {
	return e.entry.UpdatedParsed
}

// Categories returns the entry's categories
func (e FeedEntry) Categories() []string {
	return e.entry.Categories
}

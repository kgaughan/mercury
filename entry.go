package main

import (
	"html/template"
	"time"

	"github.com/mmcdole/gofeed"
)

type feedEntry struct {
	feed  *gofeed.Feed
	entry *gofeed.Item
}

func NewEntry(feed *gofeed.Feed, entry *gofeed.Item) *feedEntry {
	return &feedEntry{feed: feed, entry: entry}
}

func (e feedEntry) FeedName() string {
	return e.feed.Title
}

func (e feedEntry) SiteLink() template.URL {
	return template.URL(e.feed.Link)
}

func (e feedEntry) FeedLink() template.URL {
	return template.URL(e.feed.FeedLink)
}

func (e feedEntry) FeedPublished() *time.Time {
	if e.feed.PublishedParsed != nil {
		return e.feed.PublishedParsed
	}
	// Fallback
	return e.feed.UpdatedParsed
}

func (e feedEntry) FeedUpdated() *time.Time {
	return e.feed.UpdatedParsed
}

func (e feedEntry) FeedDescription() string {
	return e.feed.Description
}

func (e feedEntry) Title() string {
	return e.entry.Title
}

func (e feedEntry) Summary() template.HTML {
	return template.HTML(e.entry.Description)
}

func (e feedEntry) Content() template.HTML {
	return template.HTML(e.entry.Content)
}

func (e feedEntry) Link() template.URL {
	return template.URL(e.entry.Link)
}

func (e feedEntry) Author() string {
	if e.entry.Author != nil {
		return e.entry.Author.Name
	}
	if e.feed.Author != nil {
		return e.feed.Author.Name
	}
	return ""
}

func (e feedEntry) Published() *time.Time {
	if e.entry.PublishedParsed != nil {
		return e.entry.PublishedParsed
	}
	// Fallback
	return e.entry.UpdatedParsed
}

func (e feedEntry) Updated() *time.Time {
	return e.entry.UpdatedParsed
}

func (e feedEntry) Categories() []string {
	return e.entry.Categories
}

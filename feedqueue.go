package main

import (
	"time"

	"github.com/mmcdole/gofeed"
)

type feedQueue struct {
	feeds   []*gofeed.Feed
	indices []int
}

type feedAndEntry struct {
	Feed  *gofeed.Feed
	Entry *gofeed.Item
}

func (fq feedQueue) Len() int {
	return len(fq.feeds)
}

func (fq feedQueue) remaining(i int) int {
	return len(fq.feeds[i].Items) - fq.indices[i]
}

func (fq feedQueue) getPublished(i int) *time.Time {
	return fq.feeds[i].Items[fq.indices[i]].PublishedParsed
}

func (fq feedQueue) Less(i, j int) bool {
	// Exhausted feeds are considered 'greater'
	if fq.remaining(j) == 0 {
		return false
	}
	if fq.remaining(i) == 0 {
		return true
	}
	// Otherwise, it's considered 'less' if it's more recent
	return fq.getPublished(i).After(*fq.getPublished(j))
}

func (fq feedQueue) Swap(i, j int) {
	fq.feeds[i], fq.feeds[j] = fq.feeds[j], fq.feeds[i]
	fq.indices[i], fq.indices[j] = fq.indices[j], fq.indices[i]
}

// Does nothing: we only ever pop items off
func (fq *feedQueue) Push(x interface{}) {}

func (fq *feedQueue) Pop() interface{} {
	i := fq.indices[0]
	// If there's nothing more to process, we return nil
	if i == len(fq.feeds[0].Items) {
		return nil
	}
	fq.indices[0]++
	return &feedAndEntry{
		Feed:  fq.feeds[0],
		Entry: fq.feeds[0].Items[i],
	}
}

func (fq *feedQueue) AppendFeed(feed *gofeed.Feed) {
	(*fq).feeds = append((*fq).feeds, feed)
	(*fq).indices = append((*fq).indices, 0)
}

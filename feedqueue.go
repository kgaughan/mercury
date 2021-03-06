package main

import (
	"time"

	"github.com/mmcdole/gofeed"
)

type feedQueue struct {
	feeds   []*gofeed.Feed
	indices []int
}

func (fq feedQueue) Len() int {
	return len(fq.feeds)
}

func (fq feedQueue) remaining(i int) int {
	return len(fq.feeds[i].Items) - fq.indices[i]
}

func (fq feedQueue) getPublished(i int) *time.Time {
	entry := fq.feeds[i].Items[fq.indices[i]]
	// This is ridiculous.
	if entry.PublishedParsed != nil {
		return entry.PublishedParsed
	}
	return entry.UpdatedParsed
}

func (fq feedQueue) Less(i, j int) bool {
	iRemaining := fq.remaining(i)
	jRemaining := fq.remaining(j)
	if iRemaining > 0 && jRemaining > 0 {
		return fq.getPublished(i).After(*fq.getPublished(j))
	}
	return iRemaining > jRemaining
}

func (fq feedQueue) Swap(i, j int) {
	fq.feeds[i], fq.feeds[j] = fq.feeds[j], fq.feeds[i]
	fq.indices[i], fq.indices[j] = fq.indices[j], fq.indices[i]
}

// Does nothing
func (fq *feedQueue) Push(x interface{}) {}

// Does nothing
func (fq *feedQueue) Pop() interface{} {
	return nil
}

func (fq *feedQueue) Top() interface{} {
	i := fq.indices[0]
	// If there's nothing more to process, we return nil
	if i == len(fq.feeds[0].Items) {
		return nil
	}
	fq.indices[0]++
	return NewEntry(fq.feeds[0], fq.feeds[0].Items[i])
}

func (fq *feedQueue) AppendFeed(feed *gofeed.Feed) {
	(*fq).feeds = append((*fq).feeds, feed)
	(*fq).indices = append((*fq).indices, 0)
}

package internal

import (
	"time"

	"github.com/mmcdole/gofeed"
)

type FeedQueue struct {
	feeds   []*gofeed.Feed
	indices []int
}

func (fq FeedQueue) Len() int {
	return len(fq.feeds)
}

func (fq FeedQueue) remaining(i int) int {
	return len(fq.feeds[i].Items) - fq.indices[i]
}

func (fq FeedQueue) getPublished(i int) *time.Time {
	entry := fq.feeds[i].Items[fq.indices[i]]
	// This is ridiculous.
	if entry.PublishedParsed != nil {
		return entry.PublishedParsed
	}
	return entry.UpdatedParsed
}

func (fq FeedQueue) Less(i, j int) bool {
	iRemaining := fq.remaining(i)
	jRemaining := fq.remaining(j)
	if iRemaining > 0 && jRemaining > 0 {
		return fq.getPublished(i).After(*fq.getPublished(j))
	}
	return iRemaining > jRemaining
}

func (fq FeedQueue) Swap(i, j int) {
	fq.feeds[i], fq.feeds[j] = fq.feeds[j], fq.feeds[i]
	fq.indices[i], fq.indices[j] = fq.indices[j], fq.indices[i]
}

// Does nothing
func (fq *FeedQueue) Push(x interface{}) {}

// Does nothing
func (fq *FeedQueue) Pop() interface{} {
	return nil
}

func (fq *FeedQueue) Top() interface{} {
	i := fq.indices[0]
	// If there's nothing more to process, we return nil
	if i == len(fq.feeds[0].Items) {
		return nil
	}
	fq.indices[0]++
	return NewEntry(fq.feeds[0], fq.feeds[0].Items[i])
}

func (fq *FeedQueue) AppendFeed(feed *gofeed.Feed) {
	(*fq).feeds = append((*fq).feeds, feed)
	(*fq).indices = append((*fq).indices, 0)
}

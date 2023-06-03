package feed

import (
	"time"

	"github.com/mmcdole/gofeed"
)

type Queue struct {
	feeds   []*gofeed.Feed
	indices []int
}

func (fq Queue) Len() int {
	return len(fq.feeds)
}

func (fq Queue) remaining(i int) int {
	return len(fq.feeds[i].Items) - fq.indices[i]
}

func (fq Queue) getPublished(i int) *time.Time {
	entry := fq.feeds[i].Items[fq.indices[i]]
	// This is ridiculous.
	if entry.PublishedParsed != nil {
		return entry.PublishedParsed
	}
	return entry.UpdatedParsed
}

func (fq Queue) Less(i, j int) bool {
	iRemaining := fq.remaining(i)
	jRemaining := fq.remaining(j)
	if iRemaining > 0 && jRemaining > 0 {
		this := fq.getPublished(i)
		that := fq.getPublished(j)
		// The dates can potentially be non-existent, so we treat nil as
		// infinitely old.
		if this == nil {
			return false
		}
		if that == nil {
			return true
		}
		return this.After(*that)
	}
	return iRemaining > jRemaining
}

func (fq Queue) Swap(i, j int) {
	fq.feeds[i], fq.feeds[j] = fq.feeds[j], fq.feeds[i]
	fq.indices[i], fq.indices[j] = fq.indices[j], fq.indices[i]
}

// Does nothing
func (fq *Queue) Push(_ interface{}) {}

// Does nothing
func (fq *Queue) Pop() interface{} {
	return nil
}

func (fq *Queue) Top() interface{} {
	i := fq.indices[0]
	// If there's nothing more to process, we return nil
	if i == len(fq.feeds[0].Items) {
		return nil
	}
	fq.indices[0]++
	return NewEntry(fq.feeds[0], fq.feeds[0].Items[i])
}

func (fq *Queue) Append(feed *gofeed.Feed) {
	fq.feeds = append(fq.feeds, feed)
	fq.indices = append(fq.indices, 0)
}

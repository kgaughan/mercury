package feed

import (
	"container/heap"
	"time"

	"github.com/mmcdole/gofeed"
)

// queueItem pairs a feed with the index of the entry to be processed next.
type queueItem struct {
	feed *gofeed.Feed
	idx  int
}

// hasMore returns true when this queue item has unprocessed entries.
func (qi queueItem) hasMore() bool {
	if qi.feed == nil || qi.feed.Items == nil {
		return false
	}
	return qi.idx < len(qi.feed.Items)
}

// remaining returns the number of unprocessed entries in this queue item.
func (qi queueItem) remaining() int {
	if !qi.hasMore() {
		return 0
	}
	return max(0, len(qi.feed.Items)-qi.idx)
}

// getCurrentEntry returns the current entry in this queue item, or nil if none.
func (qi queueItem) getCurrentEntry() *gofeed.Item {
	if !qi.hasMore() {
		return nil
	}
	return qi.feed.Items[qi.idx]
}

// Queue is a priority queue for merging multiple feeds.
type Queue struct {
	items []queueItem
}

func (fq *Queue) Len() int {
	return len(fq.items)
}

// isSafeIndex returns true when i is a valid index into fq.items.
func (fq *Queue) isSafeIndex(i int) bool {
	return i >= 0 && i < len(fq.items)
}

// remaining returns the number of unprocessed entries in item i.
func (fq *Queue) remaining(i int) int {
	if !fq.isSafeIndex(i) {
		return 0
	}
	return fq.items[i].remaining()
}

// getPublished returns the published/updated time of the current entry in item i,
// or nil when it cannot be determined safely.
func (fq *Queue) getPublished(i int) *time.Time {
	if !fq.isSafeIndex(i) {
		return nil
	}
	entry := fq.items[i].getCurrentEntry()
	if entry == nil {
		return nil
	}
	// This is ridiculous.
	if entry.PublishedParsed != nil {
		return entry.PublishedParsed
	}
	return entry.UpdatedParsed
}

// Less orders items by the timestamp of their current entry (newest first).
// When timestamps are not available, compare by remaining count.
func (fq *Queue) Less(i, j int) bool {
	iRemaining := fq.remaining(i)
	jRemaining := fq.remaining(j)
	if iRemaining > 0 && jRemaining > 0 {
		this := fq.getPublished(i)
		that := fq.getPublished(j)
		// treat nil as infinitely old
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

func (fq *Queue) Swap(i, j int) {
	fq.items[i], fq.items[j] = fq.items[j], fq.items[i]
}

// Push adds an element to the heap (required by heap.Interface).
func (fq *Queue) Push(x any) {
	fq.items = append(fq.items, x.(queueItem))
}

// Pop removes and returns the last element (required by heap.Interface).
func (fq *Queue) Pop() any {
	n := len(fq.items)
	if n == 0 {
		return nil
	}
	v := fq.items[n-1]
	fq.items = fq.items[:n-1]
	return v
}

// allEmpty returns true when no item has remaining entries.
func (fq *Queue) allEmpty() bool {
	for i := range fq.items {
		if fq.remaining(i) > 0 {
			return false
		}
	}
	return true
}

// Top returns the next Entry to emit (advancing that feed's index).
// It will re-heapify and skip exhausted items until a valid entry is found
// or there are no entries left. Returns nil when finished.
func (fq *Queue) Top() any {
	if len(fq.items) == 0 {
		return nil
	}

	for {
		if fq.remaining(0) > 0 {
			it := &fq.items[0]
			if it.hasMore() {
				entry := it.getCurrentEntry()
				it.idx++
				return NewEntry(it.feed, entry)
			}
			// treat as exhausted and fallthrough to re-heapify
		}

		if fq.allEmpty() {
			return nil
		}

		heap.Fix(fq, 0)
	}
}

// Append adds a new feed to the queue.
func (fq *Queue) Append(feed *gofeed.Feed) {
	fq.items = append(fq.items, queueItem{feed: feed, idx: 0})
}

func (fq *Queue) Shuffle(nEntries int) []*Entry {
	if nEntries <= 0 {
		return nil
	}
	var items []*Entry
	heap.Init(fq)
	for range nEntries {
		item := fq.Top()
		if item == nil {
			break
		}
		items = append(items, item.(*Entry))
		heap.Fix(fq, 0)
	}
	return items
}

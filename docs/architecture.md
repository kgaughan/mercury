# Architecture

## Splicing

To efficiently splice together multiple feeds by date, Mercury arranges the feeds into a [heap][], and for each feed maintains the index of the current topmost item in the feed. As it pulls out entries, it performs a _heapify_ operation using the publication date of the topmost item in each feed as the key. This guarantees that the first feed is always the one with the most recent entry while doing the least amount of processing. It then takes the topmost item in that feed, increments its counter, and starts the cycle again. This effectively does a partial _heapsort_ of the total collection of feed items.

[heap]: https://en.wikipedia.org/wiki/Heap_(data_structure)

One weakness with the current implementaton is that it doesn't yet deal with feeds that don't sort their items in reverse chronological order. Two options would be to sort feeds as they're pulled down or to heapify the entries similarly to what's done with the feeds themselves. This latter approach is probably the best one to take if there are many feeds with many items. For feeds that are already sorted, the heap property should already hold, so the impact for them would be negligible.

## Cache directory

The cache directory stores the cache manifest in `manifest.json`, and cached copies of feeds serialised as JSON files.

The manifest file consists of a map of feed URLs to metadata, including the UUID used for the cache filename, the last modified date, and the entity tag, both of which are used to make [conditional GET][] requests against the target feed to minimise Mercury's bandwidth impact.

[conditional GET]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Conditional_requests

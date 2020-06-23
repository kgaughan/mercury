# TODO

- [ ] If a cache miss causes an item to be fetched, keep the feed in memory rather than loading the copy saved to disc.
- [ ] Ensure fetched feeds are sorted from newest to oldest item. Could this be done lazily with container/heap? This would be more complicated than using the sort package, but would have the advantage of bubbling up only the required items in the feed.
- [ ] Is there a good way to expose further feed/entry information in the template?
- [ ] Themes should include a BOM indicating which files should be copied across to the output directory.
- [ ] Give a way to specify particular categories that should be included from the source feed. This would be useful for blogs that lack per-category feeds.

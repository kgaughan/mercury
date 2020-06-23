# Planet Mercury

<img src="assets/mercury.png" align="right">_Mercury_ is intended a replacement for Sam Ruby's [Planet Venus](https://github.com/rubys/venus/).

A _planet_ is a kind of feed aggregator. It takes a list of newsfeeds (Atom, RSS, &c.), splices them together, and spits a set of HTML pages and/or a feed.

## Quickstart

By default, _mercury_ will look for look for a file called _mercury.toml_ in the current directory. This feed is in [TOML][] format, but the key thing you need to know is that keys and values are separated with an `=`, string values must be quoted, and `[[feed]]` introduces new feed configuration.

If you want to use an explicitly named configuration file, you can pass this with the `-config` flag.

[TOML]: https://en.wikipedia.org/wiki/TOML

### Configuration

The top-level configuration fields are:

| Name | Type | Description |
| ---- | ---- | ----------- |
| name | string | The name of your planet |
| url | string | The base URL of your planet |
| owner | string | Your name |
| email | string | Your email |
| cache | string | The path, relative to _mercury.toml_ of the feed cache |
| timeout | duration | How long to wait when fetching a feed |
| theme | string | The path, relative to _mercury.toml_ of the theme to use |
| output | string | The path, relative to _mercury.toml_ to which _mercury_ should write the files it generates |
| items | number | The number of items to include per page |
| max_pages | number | The maximum number of pages to generate |

A _duration_ is a sequence of numbers followed by a unit, with 's' being 'second', 'm' being 'minute', and 'h' being 'hour'. Thus '5m30' would mean five minutes and thirty seconds.

Each feed is introduced with `[[feed]]`, and can contain the following fields:

| Name | Type | Description |
| ---- | ---- | ----------- |
| name | string | The name of the feed |
| feed | string | The URL of the feed. Note that this must be the URL of the _feed_ itself and no attempt is made to do feed discovery if all that's provided is the site's homepage |

### Themes

A theme is a directory that contains the template files and assets needed to generate the site. Currently, there is only one template `index.html`. _Mercury_ uses Go's [html/template][] library, which is built upon [text/template][]. You should read the documentation for the latter to get a feel for the templating language and read the former for any HTML-specific behavioural differences differences.

[html/template]: https://golang.org/pkg/html/template/
[text/template]: https://golang.org/pkg/text/template/

The template supplied in [theme/index.html](theme/index.html) should be a good jumping-off point. The top-level fields available are:

| Name | Description |
| ---- | ----------- |
| .Name | The _name_ field from _config.toml_ |
| .URL | The _url_ field from _config.toml_ |
| .Owner | The _owner_ field from _config.toml_ |
| .Email | The _email_ field from _config.toml_ |
| .PageNo | The current page number |
| .Items | A collection of feed items to be rendered |

Each feed item has the following fields:

| Name | Description | Extras |
| ---- | ----------- | ------ |
| .FeedName | The name of the feed this item is from | |
| .SiteLink | The site homepage for the feed | |
| .FeedLink | The URL of the feed | |
| .FeedPublished | The publication date/time of the feed | formattable |
| .FeedUpdated | The date/time of when the feed was updated | formattable |
| .Title | The title of the item | |
| .Summary | A summary of the feed item, if available | |
| .Content | The entire content of the feed item, if available | |
| .Link | A lnk back to the original post | |
| .Author | The author name, if available | |
| .Published | The publication date/time of the entry | formattable |
| .Updated | The date/time of when entry was updated | formattable |
| .Categories | A collection of categories associated with the entry | rangeable |

You can use [text/template][]'s `{{range}}` action to iterate over the values in `.Items`. For instance, to print the title of each entry and its link, you'd do:

```
<ul>
{{range .Items}}
<li><a href="{{.Link}}">{{.Title}}</a></li>
{{end}}
</ul>
```

_Mercury_ provides a number of utility filter functions you can use for formatting dates. `isodatefmt` formats the date in [ISO 8601][]/[RFC 3339][] format, which is machine readable and useful for the `<time>` tag's `datetime` attribute. `datefmt` takes a single argument, a format specification for Go's [time/Time.Format][] function and renders it accordingly.

[ISO 8601]: https://en.wikipedia.org/wiki/ISO_8601
[RFC 3339]: https://www.ietf.org/rfc/rfc3339.txt
[time/Time.Format]: https://golang.org/pkg/time/#Time.Format

Here's an example of both being used:

```
<time datetime="{{.Published | isodatefmt}}">{{.Published | datefmt "January 2, 2006 at 15:04:05 MST"}}</time>
```

## Architecture

### Splicing

To efficiently splice together multiple feeds by date, Mercury arranges the feeds into a [heap][], and for each feed maintains the index of the current topmost item in the feed. As it pulls out entries, it performs a _heapify_ operation using the publication date of the topmost item in each feed as the key. This guarantees that the first feed is always the one with the most recent entry while doing the least amount of processing. It then takes the topmost item in that feed, increments its counter, and starts the cycle again. This effectively does a partial _heapsort_ of the total collection of feed items.

[heap]: https://en.wikipedia.org/wiki/Heap_(data_structure)

One weakness with the current implementaton is that it doesn't yet deal with feeds that don't sort their items in reverse chronological order. Two options would be to sort feeds as they're pulled down or to heapify the entries similarly to what's done with the feeds themselves. This latter approach is probably the best one to take if there are many feeds with many items. For feeds that are already sorted, the heap property should already hold, so the impact for them would be negligible.

### Cache directory

The cache directory stores the cache manifest in `manifest.json`, and cached copies of feeds serialised as JSON files.

The manifest file consists of a map of feed URLs to metadata, including the UUID used for the cache filename, the last modified date, and the entity tag, both of which are used to make [conditional GET][] requests against the target feed to minimise Mercury's bandwidth impact.

[conditional GET]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Conditional_requests

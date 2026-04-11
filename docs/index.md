---
title: Planet Mercury
author: Keith Gaughan
date: 2026-04-11
lang: en
abstract: |
  _Mercury_ is intended as a replacement for Sam Ruby's [Planet Venus](https://github.com/rubys/venus/).

  A _planet_ is a kind of feed aggregator. It takes a list of newsfeeds (Atom, RSS, &c.), splices them together, and spits a set of HTML pages and/or a feed.
---

# Quickstart

If you have the Go toolchain already configured, you can `go install` the binary:

```console
go install github.com/kgaughan/mercury/cmd/@latest
```

By default, _mercury_ will look for look for a file called _mercury.toml_ in the current directory. This feed is in [TOML][] format, but the key thing you need to know is that keys and values are separated with an `=`, string values must be quoted, and `[[feed]]` introduces new feed configuration.

If you want to use an explicitly named configuration file, you can pass this with the `--config` flag.

Here is an example file:

```toml
name = "My Planet!"
url = "https://example.com/"
owner = "Jane Doe"
email = "jane.doe@example.com"
cache = "./cache"
timeout = "20s"
theme = "./theme"
output = "./output"
items = 10
max_pages = 2

[[feed]]
name = "Keith Gaughan"
feed = "https://keith.gaughan.ie/feeds/all.xml"

[[feed]]
name = "Inklings"
feed = "https://talideon.com/inklings/feed"
```

See [the configuration details](configuration.html) for more details on the meaning of each field.

Then run:

```console
mercury
```

This will fetch all the feeds to the cache directory and write the site to the output directory.

[TOML]: https://en.wikipedia.org/wiki/TOML

## Command line

The `--help` flag will show you the help information:

```console
$ ./mercury --help
mercury - Generates an aggregated site from a set of feeds.

Flags:
  -C, --clean-cache     clean any obsolete entries from the cache
  -c, --config string   path to configuration (default "./mercury.toml")
  -h, --help            show help
  -B, --no-build        don't build anything
  -F, --no-fetch        don't fetch, just use what's in the cache
  -V, --version         print version and exit
```

Usually, the default behaviour is what you want: mercury will try to intelligently fetch any feeds and regenerate the site. Use `--no-build` if you just want to prime the cache but don't want to generate the site. Use `--no-fetch` if you want to regenerate the site without fetching any feeds. This can be useful if you're testing out a new theme.

If you want to free up some disc space, supply `--clean-cache`, which will remove any cached feeds that no longer appear in the configuration. If you want to do this without also updating any feeds or rebuilding the site, you should supply the `--no-fetch` and `--no-build` flags at the same time.

# Configuration

The top-level configuration fields are:

| Name | Type | Description | Default |
| ---- | ---- | ----------- | ------- |
| name | string | The name of your planet | "Planet" |
| url | string | The base URL of your planet | "" |
| owner | string | Your name | "" |
| email | string | Your email | "" |
| feed_id | string | Unique ID to use for the Atom feed | "" |
| cache | string | The path, relative to _mercury.toml_ of the feed cache | "./cache" |
| generate_feed | boolean | Should a feed be generated? | true |
| timeout | duration | How long to wait when fetching a feed | - |
| theme | string | The path, relative to _mercury.toml_ of the theme to use | _use default theme_ |
| output | string | The path, relative to _mercury.toml_ to which _mercury_ should write the files it generates | "./output" |
| items | number | The number of items to include per page | 10 |
| max_pages | number | The maximum number of pages to generate | 5 |

Note that the `theme`, `output`, and `cache` paths are assumed to be relative to the directory in which the configuration file is found, not the current working directory. You can specify absolute paths in these fields, however.

A _duration_ is a sequence of numbers followed by a unit, with 's' being 'second', 'm' being 'minute', and 'h' being 'hour'. Thus '5m30s' would mean five minutes and thirty seconds.

The feed ID is a URI identifying the feed. I would recommend using a [tag URI](https://en.wikipedia.org/wiki/Tag_URI_scheme), or a [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier) [URN](https://en.wikipedia.org/wiki/Uniform_Resource_Name). In the latter case, use a UUID generator such as `uuidgen` to generate a UUID, prefix it with `urn:uuid:`, and use the result as the value of `feed_id`.

Each feed is introduced with `[[feed]]`, and can contain the following fields:

| Name | Type | Description |
| ---- | ---- | ----------- |
| name | string | The name of the feed |
| feed | string | The URL of the feed. Note that this must be the URL of the _feed_ itself and no attempt is made to do feed discovery if all that's provided is the site's homepage |

## Filters

Filters are defined by adding configuration sections named `[[feed.filter]]` subsequent to the corresponding `[[feed]]` entry. Filters are defined using [Expr](https://expr-lang.org/docs/language-definition), and your filter is expected to take a feed entry and return true if the entry should be kept, or false if not.

| Name | Description |
| ---- | ----------- |
| when | An expression to determine whether the entry should be kept or skipped. This should evaluate to a boolean. Defaults to `true` |
<!--
| transform | A transformation to apply to each entry in the feed. This is only executed if `when` evaluates to `true`. |
-->

The entry is available in your filter's environment in the variable `entry`.  See the [gofeed documentation on the Item type](https://pkg.go.dev/github.com/mmcdole/gofeed#Item) for details on the fields you can expect.

Here's an example filter that only allows through an entry if its `Title` field contains the letter 'e':

```toml
[[feed.filter]]
when = "entry.Title contains 'e'"
```

## Converting an OPML file into Mercury configuration

The `tools/opml2config.py` script can be used to take multiple [OPML](https://en.wikipedia.org/wiki/OPML) files.
Taking `blogs.opml` as an example:
```xml
<?xml version="1.0" encoding="utf-8"?>
<opml>
  <head>
    <title>My Subscriptions</title>
  </head>
  <body>
    <outline text="Feeds">
      <outline type="rss" text="Keith Gaughan" xmlUrl="https://keith.gaughan.ie/feeds/all.xml"/>
      <outline type="rss" text="Inklings" xmlUrl="https://talideon.com/inklings/feed"/>
    </outline>
  </body>
</opml>
```
You can process it into Mercury configuration like so:
```console
$ tools/opml2config.py blogs.opml
[[feed]]
name = "Keith Gaughan"
feed = "https://keith.gaughan.ie/feeds/all.xml"

[[feed]]
name = "Inklings"
feed = "https://talideon.com/inklings/feed"
```

# Deploying as a container

There's a container image published at `ghcr.io/kgaughan/mercury`.

In the following, we'll assume that you're mounting your configuration at `/config` and a data volume at `/data`.

Here's an example configuration file you can use to try things out.
Save this as `mercury.toml`:

```toml
name = "My Planet!"
url = "http://localhost/"
feed_id = "urn:uuid:032a6e90-899c-4d27-aa94-b99e2c1c343f"
owner = "Jane Doe"
email = "jane@example.com"
cache = "/data/cache"
timeout = "20s"
output = "/data/output"
items = 10
max_pages = 2

[[feed]]
name = "Keith Gaughan"
feed = "https://keith.gaughan.ie/feeds/all.xml"

[[feed]]
name = "Inklings"
feed = "https://talideon.com/inklings/feed"
```

Here's a quick demonstration of how to use the configuration file and mount volumes within the container.
Note the use of `-u "$(id -u):$(id -g)`: the image is based off of a Distroless image that defaults to the `nonroot` user, so this is necessary to run the _mercury_ binary as your user, otherwise it'll have issues accessing `/data` within the container.
You will need to provide the path to your configuration within the container with the `--config` flag.

```console
$ mkdir -p volumes/data volumes/config
$ cp mercury.toml volumes/config
$ docker run --rm --user "$(id -u):$(id -g)" \
    --volume ./volumes/data:/data --volume ./volumes/config:/config \
    ghcr.io/kgaughan/mercury:latest --config /config/mercury.toml
Unable to find image 'ghcr.io/kgaughan/mercury:latest' locally
latest: Pulling from kgaughan/mercury
259db2ee6b87: Pull complete
2e4cf50eeb92: Pull complete
56ce5a7a0a8c: Pull complete
e1089d61b200: Pull complete
0f8b424aa0b9: Pull complete
d557676654e5: Pull complete
d82bc7a76a83: Pull complete
d858cbc252ad: Pull complete
1069fc2daed1: Pull complete
b40161cd83fc: Pull complete
3f4e2c586348: Pull complete
eb8f5749650b: Pull complete
6a6214ee1035: Pull complete
Digest: sha256:1668181ece1cf6c5db042eff4a59bf741c65cdac823629408a044e0252d148e8
Status: Downloaded newer image for ghcr.io/kgaughan/mercury:latest
2025/10/27 23:24:16 Priming manifest with 2 feeds using 8 workers, with a queue depth of 16
2025/10/27 23:24:16 Fetching https://keith.gaughan.ie/feeds/all.xml
2025/10/27 23:24:16 https://keith.gaughan.ie/feeds/all.xml: cache not expired
2025/10/27 23:24:16 Fetching https://talideon.com/inklings/feed
2025/10/27 23:24:16 https://talideon.com/inklings/feed: cache not expired
2025/10/27 23:24:16 Finding most recent 20 entries across 2 feeds
2025/10/27 23:24:16 Writing Atom feed
2025/10/27 23:24:16 Writing OPML file
```

If you now list the contents of `volumes/data/output`, you'll see the newly-generated site.

# Themes

A theme is a directory that contains the template files and assets needed to generate the site. Currently, there is only one template `index.html`. _Mercury_ uses Go's [html/template][] library, which is built upon [text/template][]. You should read the documentation for the latter to get a feel for the templating language and read the former for any HTML-specific behavioural differences.

[html/template]: https://golang.org/pkg/html/template/
[text/template]: https://golang.org/pkg/text/template/

The template supplied in [the default built-in theme](https://github.com/kgaughan/mercury/blob/master/internal/theme/default/index.html) should be a good jumping-off point. The top-level fields available are:

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
| .Link | A link back to the original post | |
| .Author | The author name, if available | |
| .Published | The publication date/time of the entry | formattable |
| .Updated | The date/time when the entry was updated | formattable |
| .Categories | A collection of categories associated with the entry | rangeable |

You can use [text/template][]'s `{{range}}` action to iterate over the values in `.Items`. For instance, to print the title of each entry and its link, you'd do:

```html
<ul>
{{range .Items}}
<li><a href="{{.Link}}">{{.Title}}</a></li>
{{end}}
</ul>
```

_Mercury_ provides a number of utility filter functions you can use for formatting dates. `isodate` formats the date in [ISO 8601][]/[RFC 3339][] format, which is machine readable and useful for the `<time>` tag's `datetime` attribute. It also includes the [Sprig][] template function library.

[ISO 8601]: https://en.wikipedia.org/wiki/ISO_8601
[RFC 3339]: https://www.ietf.org/rfc/rfc3339.txt
[Sprig]: https://masterminds.github.io/sprig/

Here's an example of both being used:

```html
<time datetime="{{.Published | isodate}}">{{.Published | date "January 2, 2006 at 15:04:05 MST"}}</time>
```

## Bill of Materials

A theme directory must also contain a `theme.toml` file. This contains metadata about the theme (currently just its name, which is given in the `name` field) and optionally a [bill of materials] listing the files to be copied across, each entry in which is introduced with `[[bom]]`. Here's an example file:

```toml
name = "Community"

[[bom]]
path = "static/style.css"

[[bom]]
path = "static/images/banner.png"
```

This lists two files in its BOM, which are copied across when the output directory is populated.

[Bill of materials]: https://en.wikipedia.org/wiki/Bill_of_materials

# Architecture

## Splicing

To efficiently splice together multiple feeds by date, Mercury arranges the feeds into a [heap][], and for each feed maintains the index of the current topmost item in the feed. As it pulls out entries, it performs a _heapify_ operation using the publication date of the topmost item in each feed as the key. This guarantees that the first feed is always the one with the most recent entry while doing the least amount of processing. It then takes the topmost item in that feed, increments its counter, and starts the cycle again. This effectively does a partial _heapsort_ of the total collection of feed items.

[heap]: https://en.wikipedia.org/wiki/Heap_(data_structure)

One weakness with the current implementation is that it doesn't yet deal with feeds that don't sort their items in reverse chronological order. Two options would be to sort feeds as they're pulled down or to heapify the entries similarly to what's done with the feeds themselves. This latter approach is probably the best one to take if there are many feeds with many items. For feeds that are already sorted, the heap property should already hold, so the impact for them would be negligible.

## Cache directory

The cache directory stores the cache manifest in `manifest.json`, and cached copies of feeds serialised as JSON files.

The manifest file consists of a map of feed URLs to metadata, including the UUID used for the cache filename, the last modified date, and the entity tag, both of which are used to make [conditional GET][] requests against the target feed to minimise Mercury's bandwidth impact.

[conditional GET]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Conditional_requests

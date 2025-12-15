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

Filters are defined by adding sections to the configuration with sections named `[[feed.filter]]` subsequent to the corresponding `[[feed]]` entry. Filters are defined using [Expr](https://expr-lang.org/), and your filter is expected to take the entries of feed and return a filtered list of those entries.

| Name | Description |
| ---- | ----------- |
| when | An expression to determine whether the entry should be kept or skipped. This should evaluate to a boolean. Defaults to `true` |
<!--
| transform | A transformation to apply to each entry in the feed. This is only executed if `when` evaluates to `true`. |
-->

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

![](mercury.png){: width="318" height"319" align="right" style="border:0;background:transparent" }

# Planet Mercury

_Mercury_ is intended a replacement for Sam Ruby's [Planet Venus](https://github.com/rubys/venus/).

A _planet_ is a kind of feed aggregator. It takes a list of newsfeeds (Atom, RSS, &c.), splices them together, and spits a set of HTML pages and/or a feed.

## Quickstart

By default, _mercury_ will look for look for a file called _mercury.toml_ in the current directory. This feed is in [TOML][] format, but the key thing you need to know is that keys and values are separated with an `=`, string values must be quoted, and `[[feed]]` introduces new feed configuration.

If you want to use an explicitly named configuration file, you can pass this with the `-config` flag.

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

See [the configuration details](configuration.md) for more details on the meaning of each field.

[TOML]: https://en.wikipedia.org/wiki/TOML

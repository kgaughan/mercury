# Themes

A theme is a directory that contains the template files and assets needed to generate the site. Currently, there is only one template `index.html`. _Mercury_ uses Go's [html/template][] library, which is built upon [text/template][]. You should read the documentation for the latter to get a feel for the templating language and read the former for any HTML-specific behavioural differences differences.

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
| .Link | A lnk back to the original post | |
| .Author | The author name, if available | |
| .Published | The publication date/time of the entry | formattable |
| .Updated | The date/time of when entry was updated | formattable |
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

# Configuration

The top-level configuration fields are:

| Name | Type | Description | Default |
| ---- | ---- | ----------- | ------- |
| name | string | The name of your planet | "Planet" |
| url | string | The base URL of your planet | "" |
| owner | string | Your name | "" |
| email | string | Your email | "" |
| cache | string | The path, relative to _mercury.toml_ of the feed cache | "./cache" |
| timeout | duration | How long to wait when fetching a feed | - |
| theme | string | The path, relative to _mercury.toml_ of the theme to use | "./theme" |
| output | string | The path, relative to _mercury.toml_ to which _mercury_ should write the files it generates | "./output" |
| items | number | The number of items to include per page | 10
| max_pages | number | The maximum number of pages to generate | 5 |

A _duration_ is a sequence of numbers followed by a unit, with 's' being 'second', 'm' being 'minute', and 'h' being 'hour'. Thus '5m30' would mean five minutes and thirty seconds.

Each feed is introduced with `[[feed]]`, and can contain the following fields:

| Name | Type | Description |
| ---- | ---- | ----------- |
| name | string | The name of the feed |
| feed | string | The URL of the feed. Note that this must be the URL of the _feed_ itself and no attempt is made to do feed discovery if all that's provided is the site's homepage |

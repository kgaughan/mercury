# Deploying as a container

There's a container image published at `ghcr.io/kgaughan/mercury`.
It expects a volume mounted at `/data` and one at `/config` for your data and configuration respectively.

Here's an example configuration file you can use to try things out.
Save this as `mercury.toml`:

```toml
name = "My Planet!"
url = "http://localhost/"
feed_id = "uri:urn:032a6e90-899c-4d27-aa94-b99e2c1c343f"
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
Note the use of `-u "$(id -u):$(id -G | cut -f1 -d' ')`: the image is based off of a Distroless image that defaults to the root user, so this is necessary to run the _mercury_ binary as your user, otherwise it'll have issues accessing `/data` within the container.

```console
$ mkdir -p volumes/data volumes/config
$ cp mercury.toml volumes/config
$ docker run --rm --user "$(id -u):$(id -g)" \
    --volume ./volumes/data:/data --volume ./volumes/config:/config \
    ghcr.io/kgaughan/mercury:latest
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

package theme

import "embed"

//go:embed *.toml *.html robots.txt static
var Theme embed.FS

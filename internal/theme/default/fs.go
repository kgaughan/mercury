package theme

import "embed"

//go:embed *.toml *.html static
var Theme embed.FS

package flags

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var (
	PrintVersion = flag.Bool("version", false, "Print version and exit")
	ConfigPath   = flag.String("config", "./mercury.toml", "Path to configuration")
	NoFetch      = flag.Bool("no-fetch", false, "Don't fetch, just use what's in the cache")
)

func init() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		name := path.Base(os.Args[0])
		fmt.Fprintf(out, "%s - Generates an aggregated site from a set of feeds.\n\n", name)
		fmt.Fprintf(out, "Usage:\n\n")
		flag.PrintDefaults()
	}
}

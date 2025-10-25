package flags

import (
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

var (
	PrintVersion = flag.BoolP("version", "V", false, "print version and exit")
	ConfigPath   = flag.StringP("config", "c", "./mercury.toml", "path to configuration")
	NoFetch      = flag.BoolP("no-fetch", "F", false, "don't fetch, just use what's in the cache")
	NoBuild      = flag.BoolP("no-build", "B", false, "don't build anything")
	ShowHelp     = flag.BoolP("help", "h", false, "show help")
)

func init() {
	flag.Usage = func() {
		name := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "%v - Generates an aggregated site from a set of feeds.\n\n", name)
		fmt.Fprintf(os.Stderr, "Usage:\n  %v [options]\n\n", name)
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
}

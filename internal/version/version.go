package version

import "fmt"

// nolint: gochecknoglobals
var Version string

const repo = "https://github.com/kgaughan/mercury/"

// UserAgent returns the User-Agent string for HTTP requests.
func UserAgent() string {
	return fmt.Sprintf("planet-mercury/%v (%v)", Version, repo)
}

// Generator returns the generator string for feed generation.
func Generator() string {
	return fmt.Sprintf("Planet Mercury %v (%v)", Version, repo)
}

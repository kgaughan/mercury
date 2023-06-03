package version

import "fmt"

// nolint: gochecknoglobals
var Version string

const repo = "https://github.com/kgaughan/mercury/"

func UserAgent() string {
	return fmt.Sprintf("planet-mercury/%v (%v)", Version, repo)
}

func Generator() string {
	return fmt.Sprintf("Planet Mercury %v (%v)", Version, repo)
}

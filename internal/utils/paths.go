package utils

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
)

// cacheTagMarker is the standard marker for cache directories.
// See: https://bford.info/cachedir/
const cacheTagMarker = "Signature: 8a477f597d28d172789f06886806bc55\n"

// EnsureDir ensures that the specified path exists and is a directory.
func EnsureDir(path string) {
	if fileInfo, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0o755); err != nil {
			log.Fatal(err)
		}
	} else if !fileInfo.IsDir() {
		log.Fatalf("%q must be a directory\n", path)
	}
}

// writeCacheTag writes the cache tag file at the specified path.
func writeCacheTag(path string) error {
	return os.WriteFile(path, []byte(cacheTagMarker), 0o600) //nolint:wrapcheck
}

// EnsureCache ensures that the specified path is a valid cache directory.
func EnsureCache(path string) {
	EnsureDir(path)
	cacheTag := filepath.Join(path, "CACHEDIR.TAG")
	if _, err := os.Stat(cacheTag); os.IsNotExist(err) {
		if err = writeCacheTag(cacheTag); err != nil {
			log.Fatal(err)
		}
	} else {
		contents, err := os.ReadFile(cacheTag)
		if err != nil {
			log.Fatal(err)
		}
		if !bytes.HasPrefix(contents, []byte(cacheTagMarker)) {
			if err = writeCacheTag(cacheTag); err != nil {
				log.Fatal(err)
			}
		}
	}
}

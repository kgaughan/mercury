package utils

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
)

const cacheTagMarker = "Signature: 8a477f597d28d172789f06886806bc55\n"

func EnsureDir(path string) {
	if fileInfo, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0o755); err != nil {
			log.Fatal(err)
		}
	} else if !fileInfo.IsDir() {
		log.Fatalf("%s must be a directory\n", path)
	}
}

func writeCacheTag(path string) error {
	return os.WriteFile(path, []byte(cacheTagMarker), 0o600)
}

func EnsureCache(path string) {
	EnsureDir(path)
	cacheTag := filepath.Join(path, "CACHE.TAG")
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

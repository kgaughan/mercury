package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// ErrNotADir is returned when a path exists but is not a directory.
var ErrNotADir = fmt.Errorf("is not a directory")

// cacheTagMarker is the standard marker for cache directories.
// See: https://bford.info/cachedir/
const cacheTagMarker = "Signature: 8a477f597d28d172789f06886806bc55\n"

// EnsureDir ensures that the specified path exists and is a directory.
func EnsureDir(path string) error {
	if fileInfo, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0o755); err != nil {
			return fmt.Errorf("failed to create %q: %w", path, err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to stat %q: %w", path, err)
	} else if !fileInfo.IsDir() {
		return fmt.Errorf("%q: %w", path, ErrNotADir)
	}
	return nil
}

// writeCacheTag writes the cache tag file at the specified path.
func writeCacheTag(path string) error {
	if err := os.WriteFile(path, []byte(cacheTagMarker), 0o600); err != nil {
		return fmt.Errorf("failed to write cache tag %q: %w", path, err)
	}
	return nil
}

// EnsureCache ensures that the specified path is a valid cache directory.
func EnsureCache(path string) error {
	if err := EnsureDir(path); err != nil {
		return err //nolint:wrapcheck
	}
	cacheTag := filepath.Join(path, "CACHEDIR.TAG")
	if _, err := os.Stat(cacheTag); os.IsNotExist(err) {
		// No cache tag; create one
		if err = writeCacheTag(cacheTag); err != nil {
			return err //nolint:wrapcheck
		}
	} else if err != nil {
		return fmt.Errorf("failed to stat %q: %w", cacheTag, err)
	} else {
		// Cache tag exists; verify contents and write one out if invalid
		contents, err := os.ReadFile(cacheTag)
		if err != nil {
			return fmt.Errorf("failed to read cache tag %q: %w", cacheTag, err)
		}
		if !bytes.HasPrefix(contents, []byte(cacheTagMarker)) {
			if err = writeCacheTag(cacheTag); err != nil {
				return err //nolint:wrapcheck
			}
		}
	}
	return nil
}

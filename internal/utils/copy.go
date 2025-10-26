package utils

import (
	"fmt"
	"io"
	"os"
	"path"
)

// Copy copies data from src to a file at destPath, creating any necessary
// directories along the way.
func Copy(src io.Reader, destPath string) error {
	dirPath := path.Dir(destPath)
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return fmt.Errorf("cannot make %q: %w", dirPath, err)
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("cannot create %q: %w", destPath, err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return err //nolint:wrapcheck
	}

	if err := dest.Sync(); err != nil {
		return err //nolint:wrapcheck
	}
	return nil
}

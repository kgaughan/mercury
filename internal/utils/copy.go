package utils

import (
	"fmt"
	"io"
	"os"
	"path"
)

func Copy(srcPath, destPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("cannot open %q: %w", srcPath, err)
	}
	defer src.Close()

	dirPath := path.Dir(destPath)
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return fmt.Errorf("cannot make %q: %w", dirPath, err)
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("cannot create %q: %w", destPath, err)
	}
	defer src.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return err //nolint:wrapcheck
	}

	if err := dest.Sync(); err != nil {
		return err //nolint:wrapcheck
	}
	return nil
}

package utils

import (
	"log"
	"os"
)

func EnsureDir(path string) {
	if fileInfo, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0o755); err != nil {
			log.Fatal(err)
		}
	} else if !fileInfo.IsDir() {
		log.Fatalf("%s must be a directory\n", path)
	}
}

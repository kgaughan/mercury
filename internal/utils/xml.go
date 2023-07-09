package utils

import (
	"encoding/xml"
	"fmt"
	"os"
)

// MarshalToFile serialises an XML document to a file.
func MarshalToFile(filename string, o interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err //nolint:wrapcheck
	}
	defer f.Close()

	if _, err := f.WriteString(xml.Header); err != nil {
		return fmt.Errorf("cannot write XML header: %w", err)
	}
	encoder := xml.NewEncoder(f)
	encoder.Indent("", "\t")
	return encoder.Encode(o) //nolint:wrapcheck
}

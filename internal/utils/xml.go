package utils

import (
	"encoding/xml"
	"os"
)

// MarshalToFile serialises an XML document to a file
func MarshalToFile(filename string, o interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write([]byte(xml.Header)); err != nil {
		return err
	}
	encoder := xml.NewEncoder(f)
	encoder.Indent("", "\t")
	return encoder.Encode(o)
}
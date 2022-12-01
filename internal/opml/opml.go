package opml

import (
	"encoding/xml"
	"io"
	"os"
)

// OPML represents an OPML document
type OPML struct {
	XMLName  xml.Name   `xml:"opml"`
	Version  string     `xml:"version,attr"`
	Outlines []*Outline `xml:"body>outline"`
}

// Outline represents an outline within an OPML document
type Outline struct {
	Text   string `xml:"text,attr"`
	Type   string `xml:"type,attr"`
	XMLURL string `xml:"xmlUrl,attr"`
}

// NewOPML creates a new, empty OPML document
func New(size int) *OPML {
	return &OPML{
		Version:  "2.0",
		Outlines: make([]*Outline, 0, size),
	}
}

// Append appends a feed entry to the OPML document
func (o *OPML) Append(text, xmlURL string) {
	o.Outlines = append(o.Outlines, &Outline{
		Text:   text,
		Type:   "rss",
		XMLURL: xmlURL,
	})
}

// Marshal serialises the OPML document to w
func (o *OPML) Marshal(w io.Writer) error {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "\t")
	return encoder.Encode(o)
}

// MarshalToFile serialises the OPML document to a file
func (o *OPML) MarshalToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := o.Marshal(f); err != nil {
		return err
	}
	return nil
}

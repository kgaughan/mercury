package atom

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

// Feed represents an Atom feed.
type Feed struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Title   string   `xml:"title"`
	ID      string   `xml:"id"`
	Links   []Link   `xml:"link"`
	Updated TimeStr  `xml:"updated"`
	Author  *Person  `xml:"author,omitempty"`
	Entries []*Entry `xml:"entry"`
}

// Entry represents a single entry in an Atom feed.
type Entry struct {
	Title     string  `xml:"title"`
	ID        string  `xml:"id"`
	Links     []Link  `xml:"link"`
	Published TimeStr `xml:"published,omitempty"`
	Updated   TimeStr `xml:"updated,omitempty"`
	Author    *Person `xml:"author,omitempty"`
	Summary   *Text   `xml:"summary,omitempty"`
	Content   *Text   `xml:"content,omitempty"`
}

// Link represents a link in an Atom feed or entry.
type Link struct {
	Rel      string `xml:"rel,attr,omitempty"`
	Href     string `xml:"href,attr"`
	Type     string `xml:"type,attr,omitempty"`
	HrefLang string `xml:"hreflang,attr,omitempty"`
	Title    string `xml:"title,attr,omitempty"`
	Length   uint   `xml:"length,attr,omitempty"`
}

// Person represents an author or contributor in an Atom feed or entry.
type Person struct {
	Name     string `xml:"name"`
	URI      string `xml:"uri,omitempty"`
	Email    string `xml:"email,omitempty"`
	InnerXML string `xml:",innerxml"`
}

// Text represents text content in an Atom feed or entry.
type Text struct {
	Type string `xml:"type,attr,omitempty"` // defaults to "text"
	Body string `xml:",chardata"`
}

// TimeStr is a string representation of time in RFC3339 format.
type TimeStr string

// Time converts a time.Time to a TimeStr in RFC3339 format.
func Time(t time.Time) TimeStr {
	return TimeStr(t.Format(time.RFC3339))
}

// Marshal marshals the Atom feed to XML.
func (feed *Feed) Marshal() ([]byte, error) {
	out, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Atom feed: %w", err)
	}
	return append([]byte(xml.Header), out...), nil
}

// Save writes the Atom feed to a file.
func (feed *Feed) Save(filename string) error {
	data, err := feed.Marshal()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0o600) //nolint:wrapcheck
}

// Unmarshal parses an Atom feed from XML data.
func Unmarshal(data []byte) (*Feed, error) {
	feed := &Feed{}
	if err := xml.Unmarshal(data, feed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Atom feed: %w", err)
	}
	return feed, nil
}

// Load reads and parses an Atom feed from a file.
func Load(filename string) (*Feed, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read Atom feed file: %w", err)
	}
	return Unmarshal(data)
}

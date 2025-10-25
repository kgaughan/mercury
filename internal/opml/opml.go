package opml

import (
	"encoding/xml"
	"fmt"
	"os"
)

// OPML represents an OPML document.
type OPML struct {
	XMLName  xml.Name   `xml:"opml"`
	Version  string     `xml:"version,attr"`
	Head     *Head      `xml:"head,omitempty"`
	Outlines []*Outline `xml:"body>outline"`
}

// Head represents the head section of an OPML document.
type Head struct {
	Title        string `xml:"title,omitempty"`
	OwnerName    string `xml:"ownerName,omitempty"`
	OwnerEmail   string `xml:"ownerEmail,omitempty"`
	DateCreated  string `xml:"dateCreated,omitempty"`
	DateModified string `xml:"dateModified,omitempty"`
}

// Outline represents an outline within an OPML document.
type Outline struct {
	Text   string `xml:"text,attr,omitempty"`
	Type   string `xml:"type,attr,omitempty"`
	XMLURL string `xml:"xmlUrl,attr,omitempty"`
}

// NewOPML creates a new, empty OPML document.
func New(size int) *OPML {
	return &OPML{
		Version:  "2.0",
		Outlines: make([]*Outline, 0, size),
	}
}

// SetTitle sets the title of the OPML document.
func (o *OPML) SetTitle(title string) {
	if o.Head == nil {
		o.Head = &Head{}
	}
	o.Head.Title = title
}

// SetOwner sets the owner name and email of the OPML document.
func (o *OPML) SetOwner(name, email string) {
	if o.Head == nil {
		o.Head = &Head{}
	}
	o.Head.OwnerName = name
	o.Head.OwnerEmail = email
}

// Append appends a feed entry to the OPML document.
func (o *OPML) Append(text, xmlURL string) {
	if xmlURL == "" {
		return
	}
	o.Outlines = append(o.Outlines, &Outline{
		Text:   text,
		Type:   "rss",
		XMLURL: xmlURL,
	})
}

// Marshal serializes the OPML document to XML.
func (o *OPML) Marshal() ([]byte, error) {
	out, err := xml.MarshalIndent(o, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OPML document: %w", err)
	}
	return append([]byte(xml.Header), out...), nil
}

// Save writes the OPML document to a file.
func (o *OPML) Save(filename string) error {
	data, err := o.Marshal()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0o600) //nolint:wrapcheck
}

// Unmarshal parses an OPML document from XML data.
func Unmarshal(data []byte) (*OPML, error) {
	o := &OPML{}
	if err := xml.Unmarshal(data, &o); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OPML document: %w", err)
	}
	return o, nil
}

// Load reads and parses an OPML document from a file.
func Load(filename string) (*OPML, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read OPML file: %w", err)
	}
	return Unmarshal(data)
}

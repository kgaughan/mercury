package main

import (
	"encoding/xml"
	"io"
	"os"
)

type Opml struct {
	XMLName  xml.Name   `xml:"opml"`
	Version  string     `xml:"version,attr"`
	Outlines []*Outline `xml:"body>outline"`
}

type Outline struct {
	Text   string `xml:"text,attr"`
	Type   string `xml:"type,attr"`
	XmlUrl string `xml:"xmlUrl,attr"`
}

func NewOpml(size int) *Opml {
	return &Opml{
		Version:  "2.0",
		Outlines: make([]*Outline, 0, size),
	}
}

func (o *Opml) Append(text, xmlUrl string) {
	o.Outlines = append(o.Outlines, &Outline{
		Text:   text,
		Type:   "rss",
		XmlUrl: xmlUrl,
	})
}

func (o *Opml) Marshal(w io.Writer) error {
	w.Write([]byte(xml.Header))
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "\t")
	return encoder.Encode(o)
}

func (o *Opml) MarshalToFile(filename string) error {
	if f, err := os.Create(filename); err != nil {
		return err
	} else {
		defer f.Close()
		if err := o.Marshal(f); err != nil {
			return err
		}
	}
	return nil
}

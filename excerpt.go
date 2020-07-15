package main

import (
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

func excerpt(text string, max int) string {
	var b strings.Builder

	t := html.NewTokenizer(strings.NewReader(text))

	remaining := max
	tagStack := make([]string, 0, 16) // This should be plenty to avoid reallocation

Loop:
	for {
		tt := t.Next()
		token := t.Token()

		switch tt {
		case html.ErrorToken:
			if t.Err() != io.EOF {
				log.Println(t.Err())
			}
			break Loop

		case html.TextToken:
			rs := []rune(token.Data)
			nr := len(rs)
			var toAppend string
			if nr <= remaining {
				toAppend = token.Data
				remaining -= nr
			} else {
				toAppend = string(rs[:remaining]) + "\u2026" // ellipsis
				remaining = 0
			}
			b.WriteString(html.EscapeString(toAppend))
			if remaining == 0 {
				break Loop
			}

		case html.StartTagToken:
			b.WriteString(token.String())
			tagStack = append(tagStack, token.Data)

		case html.EndTagToken:
			b.WriteString(token.String())
			tagStack = tagStack[:len(tagStack)-1]
		}
	}

	for i := len(tagStack) - 1; i >= 0; i-- {
		b.WriteString("</")
		b.WriteString(tagStack[i])
		b.WriteByte('>')
	}

	return b.String()
}

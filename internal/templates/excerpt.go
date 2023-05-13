package templates

import (
	"errors"
	"io"
	"log"
	"strings"
	"unicode"

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
			if !errors.Is(t.Err(), io.EOF) {
				log.Println(t.Err())
			}
			break Loop

		case html.TextToken:
			toAppend := token.Data
			lastSpace := -1
			for n, r := range toAppend {
				if unicode.IsSpace(r) {
					lastSpace = n
				}
				remaining--
				if remaining == 0 {
					// Get the biggest slice we can if no space was found up
					// to the truncation point.
					if lastSpace == -1 {
						lastSpace = n
					}
					toAppend = toAppend[:lastSpace]
					break
				}
			}

			b.WriteString(html.EscapeString(toAppend))
			if remaining == 0 {
				b.WriteRune('\u2026') // ellipsis
				break Loop
			}

		case html.StartTagToken:
			b.WriteString(token.String())
			tagStack = append(tagStack, token.Data)

		case html.EndTagToken:
			b.WriteString(token.String())
			tagStack = tagStack[:len(tagStack)-1]

		case html.CommentToken:
		case html.SelfClosingTagToken:
		case html.DoctypeToken:
			// Ignore
		}
	}

	for i := len(tagStack) - 1; i >= 0; i-- {
		b.WriteString("</")
		b.WriteString(tagStack[i])
		b.WriteByte('>')
	}

	return b.String()
}

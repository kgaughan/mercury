package templates

import (
	"errors"
	"io"
	"log"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

const ellipsis = '\u2026'

func excerpt(text string, maxlen int) string {
	var b strings.Builder

	t := html.NewTokenizer(strings.NewReader(text))

	remaining := maxlen
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
			toAppend, remaining := truncateText(token.Data, remaining)
			b.WriteString(html.EscapeString(toAppend))
			if remaining == 0 {
				b.WriteRune(ellipsis)
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

func truncateText(text string, remaining int) (string, int) {
	lastSpace := -1
	for n, r := range text {
		if unicode.IsSpace(r) {
			lastSpace = n
		}
		remaining--
		if remaining == 0 {
			// Get the biggest slice we can if no space was found up
			// to the truncation point.
			if lastSpace == -1 {
				lastSpace = n + 1
			}
			text = text[:lastSpace]
			break
		}
	}

	return text, remaining
}

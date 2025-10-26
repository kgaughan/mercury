package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncateText(t *testing.T) {
	result, n := truncateText("this has multiple words", 10)
	assert.Equal(t, "this has", result)
	assert.Equal(t, 0, n)

	result, n = truncateText("this has multiple words", 2)
	assert.Equal(t, "th", result)
	assert.Equal(t, 0, n)

	result, n = truncateText("this is short", 20)
	assert.Equal(t, "this is short", result)
	assert.Equal(t, 7, n)
}

func TestExcerpt(t *testing.T) {
	htmlInput := `<p>This is a <strong>test</strong> of the excerpt function. It should <em>correctly</em> handle <a href="#">HTML</a> tags and truncate the text appropriately.</p>`

	result := excerpt(htmlInput, 5)
	expected := "<p>This\u2026</p>"
	assert.Equal(t, expected, result)

	result = excerpt(htmlInput, 20)
	expected = "<p>This is a <strong>test</strong> of the excerpt\u2026</p>"
	assert.Equal(t, expected, result)

	result = excerpt(htmlInput, 100)
	assert.Equal(t, htmlInput, result)
}

package templates

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTruncateText(t *testing.T) {
	result, n := truncateText("this has multiple words", 10)
	if diff := cmp.Diff("this has", result); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}
	if n != 0 {
		t.Errorf("expected no remaining characters; got %v", n)
	}

	result, n = truncateText("this has multiple words", 2)
	if diff := cmp.Diff("th", result); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}
	if n != 0 {
		t.Errorf("expected no remaining characters; got %v", n)
	}

	result, n = truncateText("this is short", 20)
	if diff := cmp.Diff("this is short", result); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}
	if n != 7 {
		t.Errorf("expected 7 remaining characters; got %v", n)
	}
}

func TestExcerpt(t *testing.T) {
	htmlInput := `<p>This is a <strong>test</strong> of the excerpt function. It should <em>correctly</em> handle <a href="#">HTML</a> tags and truncate the text appropriately.</p>`

	result := excerpt(htmlInput, 5)
	expected := "<p>This\u2026</p>"
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}

	result = excerpt(htmlInput, 20)
	expected = "<p>This is a <strong>test</strong> of the excerpt\u2026</p>"
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}

	result = excerpt(htmlInput, 100)
	if diff := cmp.Diff(htmlInput, result); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}
}

package templates

import "testing"

func TestTruncateText(t *testing.T) {
	if result, n := truncateText("this has multiple words", 10); result != "this has" || n != 0 {
		t.Errorf("got %s with %d remaining", result, n)
	}
	if result, n := truncateText("this has multiple words", 2); result != "th" || n != 0 {
		t.Errorf("got %s with %d remaining", result, n)
	}
	if result, n := truncateText("this is short", 20); result != "this is short" || n != 7 {
		t.Errorf("got %s with %d remaining", result, n)
	}
}

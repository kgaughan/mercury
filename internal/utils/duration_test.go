package utils

import (
	"strings"
	"testing"
)

func TestUnmarshalDurationHappyPath(t *testing.T) {
	d := &Duration{}
	if err := d.UnmarshalText([]byte("2h")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if d.String() != "2h0m0s" {
		t.Errorf("unexpected duration value: %q", d.String())
	}
}

func TestUnmarshalDurationSadPath(t *testing.T) {
	d := &Duration{}
	if err := d.UnmarshalText([]byte("a")); !strings.Contains(err.Error(), `invalid duration "a"`) {
		t.Errorf("unexpected error: %v", err)
	}
}

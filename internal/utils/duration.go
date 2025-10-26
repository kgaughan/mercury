package utils

import (
	"fmt"
	"time"
)

// Duration is a wrapper around time.Duration that supports unmarshaling from
// text.
type Duration struct {
	time.Duration
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for our
// Duration type.
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	if d.Duration, err = time.ParseDuration(string(text)); err != nil {
		return fmt.Errorf("cannot unmartial duration: %w", err)
	}
	return nil
}

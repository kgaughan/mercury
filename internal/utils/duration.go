package utils

import (
	"fmt"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	if d.Duration, err = time.ParseDuration(string(text)); err != nil {
		return fmt.Errorf("cannot unmartial duration: %w", err)
	}
	return nil
}

package utils

import (
	"testing"
	"time"
)

func TestParseCacheControlExpiration(t *testing.T) {
	tests := []struct {
		name          string
		cc            string
		expectExpires bool
		minDuration   time.Duration
		maxDuration   time.Duration
	}{
		{
			name:          "no cache control",
			cc:            "",
			expectExpires: false,
		},
		{
			name:          "no-cache directive",
			cc:            "no-cache",
			expectExpires: false,
		},
		{
			name:          "max-age directive",
			cc:            "max-age=60",
			expectExpires: true,
			minDuration:   59 * time.Second,
			maxDuration:   61 * time.Second,
		},
		{
			name:          "no-cache and max-age directives",
			cc:            "no-cache, max-age=120",
			expectExpires: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expires time.Time
			err := ParseCacheControlExpiration(tt.cc, &expires)
			if err != nil {
				t.Fatalf("ParseCacheControlExpiration returned error: %v", err)
			}

			if tt.expectExpires {
				if expires.IsZero() {
					t.Errorf("expected expires to be set, but it is zero")
				} else {
					duration := time.Until(expires)
					if duration < tt.minDuration || duration > tt.maxDuration {
						t.Errorf("expires duration %v not within expected range [%v, %v]", duration, tt.minDuration, tt.maxDuration)
					}
				}
			} else {
				if !expires.IsZero() {
					t.Errorf("expected expires to be zero, but got %v", expires)
				}
			}
		})
	}
}

package internal

import (
	"runtime"
	"strings"
	"testing"

	dflt "github.com/kgaughan/mercury/internal/theme/default"
	"github.com/stretchr/testify/assert"
)

func TestLoadDefaults(t *testing.T) {
	cfg := &Config{}
	reader := strings.NewReader(``)
	if err := cfg.LoadFromReader(reader, "/cfg/mercury.toml"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, "Planet", cfg.Name)
	assert.Equal(t, "/cfg/cache", cfg.Cache)
	assert.Equal(t, "/cfg/output", cfg.Output)
	assert.Equal(t, 10, cfg.ItemsPerPage)
	assert.Equal(t, 5, cfg.MaxPages)
	assert.Equal(t, cfg.Parallelism*2, cfg.JobQueueDepth)
	assert.Equal(t, runtime.NumCPU(), cfg.Parallelism)
	assert.Len(t, cfg.Feeds, 0)
	assert.Equal(t, dflt.Theme, cfg.Theme)
}

func TestLoadAbsolutePaths(t *testing.T) {
	src := `
	cache = "/data/cache"
	output = "/data/output"
	`
	cfg := &Config{}
	reader := strings.NewReader(src)
	if err := cfg.LoadFromReader(reader, "/cfg/mercury.toml"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, "/data/cache", cfg.Cache)
	assert.Equal(t, "/data/output", cfg.Output)
}

func TestClampCPU(t *testing.T) {
	src := `
	parallelism = 64
	`
	cfg := &Config{}
	reader := strings.NewReader(src)
	if err := cfg.LoadFromReader(reader, "/cfg/mercury.toml"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, cpuLimit, cfg.Parallelism)
}

func TestLoadFeeds(t *testing.T) {
	src := `
	[[feed]]
	feed = "https://example.com/feed1.xml"
	[[feed]]
	feed = "https://example.com/feed2.xml"
	`
	cfg := &Config{}
	reader := strings.NewReader(src)
	if err := cfg.LoadFromReader(reader, "/cfg/mercury.toml"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Len(t, cfg.Feeds, 2)
	assert.Equal(t, "https://example.com/feed1.xml", cfg.Feeds[0].Feed)
	assert.Equal(t, "https://example.com/feed2.xml", cfg.Feeds[1].Feed)
}

func TestExternalTheme(t *testing.T) {
	// I need a directory that's guaranteed to exist for this test
	src := `
	theme = "/tmp"
	`
	cfg := &Config{}
	reader := strings.NewReader(src)
	if err := cfg.LoadFromReader(reader, "/cfg/mercury.toml"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.NotEqual(t, dflt.Theme, cfg.Theme)
	assert.Equal(t, "/tmp", cfg.ThemePath)
}

package internal

import (
	"fmt"
	"io/fs"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kgaughan/mercury/internal/manifest"
	dflt "github.com/kgaughan/mercury/internal/theme/default"
)

func TestLoadDefaults(t *testing.T) {
	cfg := &Config{}
	reader := strings.NewReader(``)
	if err := cfg.LoadFromReader(reader, "/cfg/mercury.toml"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Theme != dflt.Theme {
		t.Error("got something other than the expected default theme")
	}
	expected := &Config{
		Name:          "Planet",
		Cache:         "/cfg/cache",
		Output:        "/cfg/output",
		ItemsPerPage:  10,
		MaxPages:      5,
		GenerateFeed:  true,
		JobQueueDepth: runtime.NumCPU() * 2,
		Parallelism:   runtime.NumCPU(),
		Feeds:         nil,
		ThemePath:     "",
	}
	if diff := cmp.Diff(expected, cfg, cmpopts.IgnoreInterfaces(struct{ fs.FS }{})); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}
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
	expected := &Config{
		Name:          "Planet",
		Cache:         "/data/cache",
		Output:        "/data/output",
		ItemsPerPage:  10,
		MaxPages:      5,
		GenerateFeed:  true,
		JobQueueDepth: runtime.NumCPU() * 2,
		Parallelism:   runtime.NumCPU(),
		Feeds:         nil,
		ThemePath:     "",
	}
	if diff := cmp.Diff(expected, cfg, cmpopts.IgnoreInterfaces(struct{ fs.FS }{})); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}
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
	expected := &Config{
		Name:          "Planet",
		Cache:         "/cfg/cache",
		Output:        "/cfg/output",
		ItemsPerPage:  10,
		MaxPages:      5,
		GenerateFeed:  true,
		JobQueueDepth: cpuLimit * 2,
		Parallelism:   cpuLimit,
		Feeds:         nil,
		ThemePath:     "",
	}
	if diff := cmp.Diff(expected, cfg, cmpopts.IgnoreInterfaces(struct{ fs.FS }{})); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}
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
	expected := &Config{
		Name:          "Planet",
		Cache:         "/cfg/cache",
		Output:        "/cfg/output",
		ItemsPerPage:  10,
		MaxPages:      5,
		GenerateFeed:  true,
		JobQueueDepth: runtime.NumCPU() * 2,
		Parallelism:   runtime.NumCPU(),
		Feeds: []*manifest.Feed{
			{Name: "", Feed: "https://example.com/feed1.xml", Filters: nil},
			{Name: "", Feed: "https://example.com/feed2.xml", Filters: nil},
		},
		ThemePath: "",
	}
	if diff := cmp.Diff(expected, cfg, cmpopts.IgnoreInterfaces(struct{ fs.FS }{})); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}
}

func TestExternalTheme(t *testing.T) {
	// I need a directory that's guaranteed to exist for this test
	tmp := t.TempDir()
	src := fmt.Sprintf(`theme = "%s"`, tmp)
	cfg := &Config{}
	reader := strings.NewReader(src)
	if err := cfg.LoadFromReader(reader, "/cfg/mercury.toml"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Theme == dflt.Theme {
		t.Error("expected a theme other than the default theme")
	}
	expected := &Config{
		Name:          "Planet",
		Cache:         "/cfg/cache",
		Output:        "/cfg/output",
		ItemsPerPage:  10,
		MaxPages:      5,
		GenerateFeed:  true,
		JobQueueDepth: runtime.NumCPU() * 2,
		Parallelism:   runtime.NumCPU(),
		Feeds:         nil,
		Theme:         nil,
		ThemePath:     tmp,
	}
	if diff := cmp.Diff(expected, cfg, cmpopts.IgnoreInterfaces(struct{ fs.FS }{})); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}
}

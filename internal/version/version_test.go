package version

import (
	"runtime"
	"testing"
)

func TestInfoReturnsGoVersion(t *testing.T) {
	info := Info()
	if info["go_version"] != runtime.Version() {
		t.Errorf("go_version = %q, want %q", info["go_version"], runtime.Version())
	}
}

func TestInfoFallbackWhenLdflagsEmpty(t *testing.T) {
	// When ldflags are not set, Version/Commit/Date are empty strings.
	// Info() should still return non-empty values (from buildinfo or "unknown").
	info := Info()
	for _, key := range []string{"version", "commit", "date", "go_version"} {
		if info[key] == "" {
			t.Errorf("Info()[%q] is empty, expected a value", key)
		}
	}
}

func TestInfoUsesLdflagsWhenSet(t *testing.T) {
	// Temporarily set the package vars to simulate ldflags injection.
	origVersion, origCommit, origDate := Version, Commit, Date
	Version = "1.2.3"
	Commit = "abc1234"
	Date = "2026-01-01T00:00:00Z"
	defer func() {
		Version, Commit, Date = origVersion, origCommit, origDate
	}()

	info := Info()
	if info["version"] != "1.2.3" {
		t.Errorf("version = %q, want %q", info["version"], "1.2.3")
	}
	if info["commit"] != "abc1234" {
		t.Errorf("commit = %q, want %q", info["commit"], "abc1234")
	}
	if info["date"] != "2026-01-01T00:00:00Z" {
		t.Errorf("date = %q, want %q", info["date"], "2026-01-01T00:00:00Z")
	}
}

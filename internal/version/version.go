package version

import (
	"runtime"
	"runtime/debug"
)

// Version, Commit, and Date are set by ldflags at build time.
var (
	Version string
	Commit  string
	Date    string
)

// Info returns version metadata as a string map.
// It prefers ldflags values and falls back to debug.ReadBuildInfo().
func Info() map[string]string {
	v := Version
	c := Commit
	d := Date

	if v == "" || c == "" || d == "" {
		if bi, ok := debug.ReadBuildInfo(); ok {
			if v == "" {
				v = bi.Main.Version
			}
			for _, s := range bi.Settings {
				switch s.Key {
				case "vcs.revision":
					if c == "" {
						c = s.Value
						if len(c) > 7 {
							c = c[:7]
						}
					}
				case "vcs.time":
					if d == "" {
						d = s.Value
					}
				}
			}
		}
	}

	if v == "" {
		v = "unknown"
	}
	if c == "" {
		c = "unknown"
	}
	if d == "" {
		d = "unknown"
	}

	return map[string]string{
		"version":    v,
		"commit":     c,
		"date":       d,
		"go_version": runtime.Version(),
	}
}

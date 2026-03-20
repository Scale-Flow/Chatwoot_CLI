package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestVersionCommandOutput(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// Reset args to just "version"
	rootCmd.SetArgs([]string{"version"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON output: %v\nraw: %s", err, buf.String())
	}

	if raw["ok"] != true {
		t.Errorf("ok = %v, want true", raw["ok"])
	}
	data, ok := raw["data"].(map[string]any)
	if !ok {
		t.Fatal("data is not an object")
	}
	for _, key := range []string{"version", "commit", "date", "go_version"} {
		if _, exists := data[key]; !exists {
			t.Errorf("data.%s is missing", key)
		}
	}
}

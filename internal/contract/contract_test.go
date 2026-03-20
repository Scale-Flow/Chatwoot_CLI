package contract

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestSuccessSingleResource(t *testing.T) {
	resp := Success(map[string]any{"id": 42, "name": "Support Inbox"})

	if !resp.OK {
		t.Fatal("expected ok to be true")
	}
	if resp.Error != nil {
		t.Fatal("expected error to be nil")
	}

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if raw["ok"] != true {
		t.Errorf("ok = %v, want true", raw["ok"])
	}
	data, ok := raw["data"].(map[string]any)
	if !ok {
		t.Fatal("data is not an object")
	}
	if data["id"] != float64(42) {
		t.Errorf("data.id = %v, want 42", data["id"])
	}
	if _, exists := raw["meta"]; exists {
		t.Error("meta should be absent for single resource")
	}
	if _, exists := raw["warnings"]; exists {
		t.Error("warnings should be absent when empty")
	}
	if _, exists := raw["error"]; exists {
		t.Error("error should be absent on success")
	}
}

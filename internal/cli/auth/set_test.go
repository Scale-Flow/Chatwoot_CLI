package auth

import (
	"bytes"
	"encoding/json"
	"testing"

	keyring "github.com/zalando/go-keyring"
)

func TestAuthSetKeychain(t *testing.T) {
	keyring.MockInit()

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"auth", "set", "--mode", "application", "--token", "sk-test-123"})
	err := Cmd.Root().Execute()
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	var resp map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		t.Fatalf("JSON decode error: %v\nOutput: %s", err, stdout.String())
	}
	if resp["ok"] != true {
		t.Errorf("ok = %v, want true", resp["ok"])
	}
	data := resp["data"].(map[string]any)
	if data["mode"] != "application" {
		t.Errorf("mode = %v, want application", data["mode"])
	}
}

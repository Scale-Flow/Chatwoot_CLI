package auth

import (
	"bytes"
	"encoding/json"
	"testing"

	keyring "github.com/zalando/go-keyring"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
)

func TestAuthStatus(t *testing.T) {
	keyring.MockInit()
	ks := credentials.NewKeychainStore()
	_ = ks.Set("default", credentials.ModeApplication, "sk-test")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"auth", "status"})
	err := Cmd.Root().Execute()
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	var resp map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		t.Fatalf("JSON decode: %v\nOutput: %s", err, stdout.String())
	}
	if resp["ok"] != true {
		t.Errorf("ok = %v, want true", resp["ok"])
	}
	data := resp["data"].(map[string]any)
	creds := data["credentials"].(map[string]any)
	app := creds["application"].(map[string]any)
	if app["status"] != "configured" {
		t.Errorf("application status = %v, want configured", app["status"])
	}
}

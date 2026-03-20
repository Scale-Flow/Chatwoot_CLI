package auth

import (
	"bytes"
	"encoding/json"
	"testing"

	keyring "github.com/zalando/go-keyring"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
)

func TestAuthClearSingleMode(t *testing.T) {
	keyring.MockInit()
	ks := credentials.NewKeychainStore()
	_ = ks.Set("default", credentials.ModeApplication, "sk-test")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"auth", "clear", "--mode", "application"})
	err := Cmd.Root().Execute()
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	var resp map[string]any
	json.Unmarshal(stdout.Bytes(), &resp)
	if resp["ok"] != true {
		t.Errorf("ok = %v, want true", resp["ok"])
	}
	data := resp["data"].(map[string]any)
	cleared := data["cleared"].([]any)
	if len(cleared) != 1 || cleared[0] != "application" {
		t.Errorf("cleared = %v, want [application]", cleared)
	}

	// Verify credential was actually removed
	_, err = ks.Get("default", credentials.ModeApplication)
	if err != credentials.ErrNotFound {
		t.Errorf("credential still exists after clear")
	}
}

func TestAuthClearAll(t *testing.T) {
	keyring.MockInit()
	ks := credentials.NewKeychainStore()
	_ = ks.Set("default", credentials.ModeApplication, "sk-test")
	_ = ks.Set("default", credentials.ModePlatform, "pk-test")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"auth", "clear", "--all"})
	err := Cmd.Root().Execute()
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	var resp map[string]any
	json.Unmarshal(stdout.Bytes(), &resp)
	data := resp["data"].(map[string]any)
	cleared := data["cleared"].([]any)
	if len(cleared) != 2 {
		t.Errorf("cleared = %v, want 2 items", cleared)
	}
}

func TestAuthClearRequiresModeOrAll(t *testing.T) {
	keyring.MockInit()

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"auth", "clear"})
	err := Cmd.Root().Execute()
	if err == nil {
		t.Fatal("expected error when neither --mode nor --all provided")
	}
}

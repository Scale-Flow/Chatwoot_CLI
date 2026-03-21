package account

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":   1,
			"name": "Acme Corp",
		})
	}))
	defer srv.Close()

	t.Setenv("CHATWOOT_BASE_URL", srv.URL)
	t.Setenv("CHATWOOT_ACCESS_TOKEN", "sk-test")
	t.Setenv("CHATWOOT_ACCOUNT_ID", "1")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"application", "account", "get"})
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
}

func TestAccountUpdate(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", r.Method)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":   1,
			"name": "New Name",
		})
	}))
	defer srv.Close()

	t.Setenv("CHATWOOT_BASE_URL", srv.URL)
	t.Setenv("CHATWOOT_ACCESS_TOKEN", "sk-test")
	t.Setenv("CHATWOOT_ACCOUNT_ID", "1")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"application", "account", "update", "--name", "New Name"})
	err := Cmd.Root().Execute()
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if gotBody["name"] != "New Name" {
		t.Errorf("body name = %v, want New Name", gotBody["name"])
	}

	var resp map[string]any
	json.Unmarshal(stdout.Bytes(), &resp)
	if resp["ok"] != true {
		t.Errorf("ok = %v, want true", resp["ok"])
	}
}

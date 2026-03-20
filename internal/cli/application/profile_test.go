package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProfileGet(t *testing.T) {
	// Set up mock API server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/profile" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Header.Get("api_access_token") == "" {
			t.Error("missing auth header")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":    1,
			"name":  "Test Agent",
			"email": "agent@test.com",
		})
	}))
	defer srv.Close()

	// Set env vars for base URL and token, bypassing config file
	t.Setenv("CHATWOOT_BASE_URL", srv.URL)
	t.Setenv("CHATWOOT_ACCESS_TOKEN", "sk-test-token")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"application", "profile", "get"})
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
	if data["name"] != "Test Agent" {
		t.Errorf("name = %v, want Test Agent", data["name"])
	}
}

func TestProfileUpdateRequiresFlag(t *testing.T) {
	Cmd.SetOut(&bytes.Buffer{})
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"application", "profile", "update"})
	err := Cmd.Root().Execute()
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}

func TestProfileUpdate(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":    1,
			"name":  "New Name",
			"email": "agent@test.com",
		})
	}))
	defer srv.Close()

	t.Setenv("CHATWOOT_BASE_URL", srv.URL)
	t.Setenv("CHATWOOT_ACCESS_TOKEN", "sk-test-token")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"application", "profile", "update", "--name", "New Name"})
	err := Cmd.Root().Execute()
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if gotMethod != "PATCH" {
		t.Errorf("method = %q, want PATCH", gotMethod)
	}
	if gotBody["name"] != "New Name" {
		t.Errorf("body name = %v, want New Name", gotBody["name"])
	}
	// email should not be in body (not set)
	if _, exists := gotBody["email"]; exists {
		t.Error("email should not be in body when not set")
	}
}

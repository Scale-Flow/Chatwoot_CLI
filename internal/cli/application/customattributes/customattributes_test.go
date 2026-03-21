package customattributes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCustomAttributesList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": 1, "attribute_display_name": "Region", "attribute_key": "region"},
			},
		})
	}))
	defer srv.Close()

	t.Setenv("CHATWOOT_BASE_URL", srv.URL)
	t.Setenv("CHATWOOT_ACCESS_TOKEN", "sk-test")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"application", "custom-attributes", "list"})
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
	data, ok := resp["data"].([]any)
	if !ok || len(data) != 1 {
		t.Errorf("data length = %v, want 1", len(data))
	}
}

func TestCustomAttributesCreate(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":                     1,
			"attribute_display_name": "Region",
			"attribute_key":          "region",
		})
	}))
	defer srv.Close()

	t.Setenv("CHATWOOT_BASE_URL", srv.URL)
	t.Setenv("CHATWOOT_ACCESS_TOKEN", "sk-test")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{
		"application", "custom-attributes", "create",
		"--attribute-key", "region",
		"--attribute-model", "contact",
		"--attribute-type", "text",
		"--name", "Region",
	})
	err := Cmd.Root().Execute()
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if gotBody["attribute_key"] != "region" {
		t.Errorf("body attribute_key = %v, want region", gotBody["attribute_key"])
	}

	var resp map[string]any
	json.Unmarshal(stdout.Bytes(), &resp)
	if resp["ok"] != true {
		t.Errorf("ok = %v, want true", resp["ok"])
	}
}

package automationrules

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAutomationRulesList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"payload": []map[string]any{
				{"id": 1, "name": "Auto-assign", "event_name": "conversation_created"},
			},
		})
	}))
	defer srv.Close()

	t.Setenv("CHATWOOT_BASE_URL", srv.URL)
	t.Setenv("CHATWOOT_ACCESS_TOKEN", "sk-test")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"application", "automation-rules", "list"})
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

func TestAutomationRulesCreate(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":         1,
			"name":       "Auto-assign",
			"event_name": "conversation_created",
			"conditions": []any{},
			"actions":    []any{},
		})
	}))
	defer srv.Close()

	t.Setenv("CHATWOOT_BASE_URL", srv.URL)
	t.Setenv("CHATWOOT_ACCESS_TOKEN", "sk-test")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{
		"application", "automation-rules", "create",
		"--name", "Auto-assign",
		"--event-name", "conversation_created",
		"--conditions", `[{"attribute_key":"status","filter_operator":"equal_to","values":["open"]}]`,
		"--actions", `[{"action_name":"assign_team","action_params":[1]}]`,
	})
	err := Cmd.Root().Execute()
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if gotBody["name"] != "Auto-assign" {
		t.Errorf("body name = %v, want Auto-assign", gotBody["name"])
	}
	if gotBody["event_name"] != "conversation_created" {
		t.Errorf("body event_name = %v, want conversation_created", gotBody["event_name"])
	}
	if gotBody["conditions"] == nil {
		t.Errorf("body conditions is nil, want array")
	}
	if gotBody["actions"] == nil {
		t.Errorf("body actions is nil, want array")
	}

	var resp map[string]any
	json.Unmarshal(stdout.Bytes(), &resp)
	if resp["ok"] != true {
		t.Errorf("ok = %v, want true", resp["ok"])
	}
}

func TestAutomationRulesDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	t.Setenv("CHATWOOT_BASE_URL", srv.URL)
	t.Setenv("CHATWOOT_ACCESS_TOKEN", "sk-test")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"application", "automation-rules", "delete", "--id", "1"})
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
	if data["deleted"] != true {
		t.Errorf("deleted = %v, want true", data["deleted"])
	}
}

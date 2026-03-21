package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListAutomationRules(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/automation_rules" {
			t.Errorf("path = %q, want /api/v1/accounts/1/automation_rules", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"payload": []map[string]any{
				{"id": 1, "name": "Auto-assign", "event_name": "conversation_created"},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	rules, err := client.ListAutomationRules(context.Background())
	if err != nil {
		t.Fatalf("ListAutomationRules error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("len = %d, want 1", len(rules))
	}
	if rules[0].ID != 1 {
		t.Errorf("ID = %d, want 1", rules[0].ID)
	}
	if rules[0].Name != "Auto-assign" {
		t.Errorf("Name = %q, want Auto-assign", rules[0].Name)
	}
	if rules[0].EventName != "conversation_created" {
		t.Errorf("EventName = %q, want conversation_created", rules[0].EventName)
	}
}

func TestGetAutomationRule(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/automation_rules/7" {
			t.Errorf("path = %q, want /api/v1/accounts/1/automation_rules/7", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":         7,
			"name":       "Auto-resolve",
			"event_name": "conversation_updated",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	rule, err := client.GetAutomationRule(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetAutomationRule error: %v", err)
	}
	if rule.ID != 7 {
		t.Errorf("ID = %d, want 7", rule.ID)
	}
	if rule.Name != "Auto-resolve" {
		t.Errorf("Name = %q, want Auto-resolve", rule.Name)
	}
}

func TestCreateAutomationRule(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/automation_rules" {
			t.Errorf("path = %q, want /api/v1/accounts/1/automation_rules", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":         10,
			"name":       "Auto-assign",
			"event_name": "conversation_created",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	rule, err := client.CreateAutomationRule(context.Background(), CreateAutomationRuleOpts{
		Name:      "Auto-assign",
		EventName: "conversation_created",
	})
	if err != nil {
		t.Fatalf("CreateAutomationRule error: %v", err)
	}
	if rule.ID != 10 {
		t.Errorf("ID = %d, want 10", rule.ID)
	}
	if gotBody["name"] != "Auto-assign" {
		t.Errorf("body name = %v, want Auto-assign", gotBody["name"])
	}
	if gotBody["event_name"] != "conversation_created" {
		t.Errorf("body event_name = %v, want conversation_created", gotBody["event_name"])
	}
}

func TestDeleteAutomationRule(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/automation_rules/4" {
			t.Errorf("path = %q, want /api/v1/accounts/1/automation_rules/4", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	err := client.DeleteAutomationRule(context.Background(), 4)
	if err != nil {
		t.Fatalf("DeleteAutomationRule error: %v", err)
	}
}

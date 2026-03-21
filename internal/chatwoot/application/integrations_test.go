package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListIntegrationApps(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/integrations/apps" {
			t.Errorf("path = %q, want /api/v1/accounts/1/integrations/apps", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"payload": []map[string]any{
				{"id": "slack", "name": "Slack"},
				{"id": "dialogflow", "name": "Dialogflow"},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	apps, err := client.ListIntegrationApps(context.Background())
	if err != nil {
		t.Fatalf("ListIntegrationApps error: %v", err)
	}
	if len(apps) != 2 {
		t.Fatalf("len = %d, want 2", len(apps))
	}
	if apps[0].ID != "slack" {
		t.Errorf("ID = %q, want slack", apps[0].ID)
	}
	if apps[0].Name != "Slack" {
		t.Errorf("Name = %q, want Slack", apps[0].Name)
	}
}

func TestCreateIntegrationHook(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/integrations/hooks" {
			t.Errorf("path = %q, want /api/v1/accounts/1/integrations/hooks", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":     10,
			"app_id": "slack",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	hook, err := client.CreateIntegrationHook(context.Background(), CreateIntegrationHookOpts{
		AppID: "slack",
	})
	if err != nil {
		t.Fatalf("CreateIntegrationHook error: %v", err)
	}
	if hook.ID != 10 {
		t.Errorf("ID = %d, want 10", hook.ID)
	}
	if gotBody["app_id"] != "slack" {
		t.Errorf("body app_id = %v, want slack", gotBody["app_id"])
	}
}

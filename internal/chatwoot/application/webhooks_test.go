package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListWebhooks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/webhooks" {
			t.Errorf("path = %q, want /api/v1/accounts/1/webhooks", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"payload": map[string]any{
				"webhooks": []map[string]any{
					{"id": 1, "url": "https://example.com/hook", "subscriptions": []string{"message_created"}},
				},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	hooks, err := client.ListWebhooks(context.Background())
	if err != nil {
		t.Fatalf("ListWebhooks error: %v", err)
	}
	if len(hooks) != 1 {
		t.Fatalf("len = %d, want 1", len(hooks))
	}
	if hooks[0].ID != 1 {
		t.Errorf("ID = %d, want 1", hooks[0].ID)
	}
	if hooks[0].URL != "https://example.com/hook" {
		t.Errorf("URL = %q, want https://example.com/hook", hooks[0].URL)
	}
	if len(hooks[0].Subscriptions) != 1 || hooks[0].Subscriptions[0] != "message_created" {
		t.Errorf("Subscriptions = %v, want [message_created]", hooks[0].Subscriptions)
	}
}

func TestCreateWebhook(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/webhooks" {
			t.Errorf("path = %q, want /api/v1/accounts/1/webhooks", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"payload": map[string]any{
				"webhook": map[string]any{
					"id":            5,
					"url":           "https://example.com/hook",
					"subscriptions": []string{"message_created", "conversation_created"},
				},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	hook, err := client.CreateWebhook(context.Background(), CreateWebhookOpts{
		URL:           "https://example.com/hook",
		Subscriptions: []string{"message_created", "conversation_created"},
	})
	if err != nil {
		t.Fatalf("CreateWebhook error: %v", err)
	}
	if hook.ID != 5 {
		t.Errorf("ID = %d, want 5", hook.ID)
	}
	if gotBody["url"] != "https://example.com/hook" {
		t.Errorf("body url = %v, want https://example.com/hook", gotBody["url"])
	}
	if gotBody["subscriptions"] == nil {
		t.Errorf("body subscriptions = nil, want non-nil")
	}
}

func TestDeleteWebhook(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/webhooks/3" {
			t.Errorf("path = %q, want /api/v1/accounts/1/webhooks/3", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	err := client.DeleteWebhook(context.Background(), 3)
	if err != nil {
		t.Fatalf("DeleteWebhook error: %v", err)
	}
}

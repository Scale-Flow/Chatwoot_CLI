package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListCannedResponses(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/canned_responses" {
			t.Errorf("path = %q, want /api/v1/accounts/1/canned_responses", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "short_code": "greeting", "content": "Hello!"},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	canned, err := client.ListCannedResponses(context.Background())
	if err != nil {
		t.Fatalf("ListCannedResponses error: %v", err)
	}
	if len(canned) != 1 {
		t.Errorf("len = %d, want 1", len(canned))
	}
	if canned[0].ID != 1 {
		t.Errorf("ID = %d, want 1", canned[0].ID)
	}
	if canned[0].ShortCode != "greeting" {
		t.Errorf("ShortCode = %q, want greeting", canned[0].ShortCode)
	}
	if canned[0].Content != "Hello!" {
		t.Errorf("Content = %q, want Hello!", canned[0].Content)
	}
}

func TestCreateCannedResponse(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/canned_responses" {
			t.Errorf("path = %q, want /api/v1/accounts/1/canned_responses", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":         5,
			"short_code": "thanks",
			"content":    "Thank you for reaching out!",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	canned, err := client.CreateCannedResponse(context.Background(), CreateCannedResponseOpts{
		ShortCode: "thanks",
		Content:   "Thank you for reaching out!",
	})
	if err != nil {
		t.Fatalf("CreateCannedResponse error: %v", err)
	}
	if canned.ID != 5 {
		t.Errorf("ID = %d, want 5", canned.ID)
	}
	if canned.ShortCode != "thanks" {
		t.Errorf("ShortCode = %q, want thanks", canned.ShortCode)
	}
	if gotBody["short_code"] != "thanks" {
		t.Errorf("body short_code = %v, want thanks", gotBody["short_code"])
	}
	if gotBody["content"] != "Thank you for reaching out!" {
		t.Errorf("body content = %v, want Thank you for reaching out!", gotBody["content"])
	}
}

func TestDeleteCannedResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/canned_responses/3" {
			t.Errorf("path = %q, want /api/v1/accounts/1/canned_responses/3", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	err := client.DeleteCannedResponse(context.Background(), 3)
	if err != nil {
		t.Fatalf("DeleteCannedResponse error: %v", err)
	}
}

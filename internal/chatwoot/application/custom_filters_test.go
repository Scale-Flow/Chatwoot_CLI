package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListCustomFilters(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/custom_filters" {
			t.Errorf("path = %q, want /api/v1/accounts/1/custom_filters", r.URL.Path)
		}
		if r.URL.Query().Get("filter_type") != "conversation" {
			t.Errorf("filter_type = %q, want conversation", r.URL.Query().Get("filter_type"))
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "Open conversations", "filter_type": "conversation", "query": map[string]any{"status": "open"}},
			{"id": 2, "name": "VIP contacts", "filter_type": "contact", "query": map[string]any{"status": "active"}},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	filters, err := client.ListCustomFilters(context.Background(), "conversation")
	if err != nil {
		t.Fatalf("ListCustomFilters error: %v", err)
	}
	if len(filters) != 2 {
		t.Errorf("len = %d, want 2", len(filters))
	}
	if filters[0].Name != "Open conversations" {
		t.Errorf("filters[0].Name = %q, want Open conversations", filters[0].Name)
	}
	if filters[1].ID != 2 {
		t.Errorf("filters[1].ID = %d, want 2", filters[1].ID)
	}
	if filters[0].Type != "conversation" {
		t.Errorf("filters[0].Type = %q, want conversation", filters[0].Type)
	}
	if filters[1].Type != "contact" {
		t.Errorf("filters[1].Type = %q, want contact", filters[1].Type)
	}
}

func TestCreateCustomFilter(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/custom_filters" {
			t.Errorf("path = %q, want /api/v1/accounts/1/custom_filters", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":          30,
			"name":        "Urgent open",
			"filter_type": "conversation",
			"query":       map[string]any{"status": "open"},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	filter, err := client.CreateCustomFilter(context.Background(), CreateCustomFilterOpts{
		Name:  "Urgent open",
		Type:  "conversation",
		Query: map[string]any{"status": "open"},
	})
	if err != nil {
		t.Fatalf("CreateCustomFilter error: %v", err)
	}
	if filter.ID != 30 {
		t.Errorf("ID = %d, want 30", filter.ID)
	}
	if filter.Name != "Urgent open" {
		t.Errorf("Name = %q, want Urgent open", filter.Name)
	}
	if gotBody["name"] != "Urgent open" {
		t.Errorf("body name = %v, want Urgent open", gotBody["name"])
	}
	if gotBody["filter_type"] != "conversation" {
		t.Errorf("body filter_type = %v, want conversation", gotBody["filter_type"])
	}
	if filter.Type != "conversation" {
		t.Errorf("Type = %q, want conversation", filter.Type)
	}
}

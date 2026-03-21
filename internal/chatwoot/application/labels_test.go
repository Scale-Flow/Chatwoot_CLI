package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListLabels(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/labels" {
			t.Errorf("path = %q, want /api/v1/accounts/1/labels", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"payload": []map[string]any{
				{"id": 1, "title": "urgent", "color": "#ff0000"},
				{"id": 2, "title": "billing", "color": "#00ff00"},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	labels, err := client.ListLabels(context.Background())
	if err != nil {
		t.Fatalf("ListLabels error: %v", err)
	}
	if len(labels) != 2 {
		t.Errorf("len = %d, want 2", len(labels))
	}
	if labels[0].Title != "urgent" {
		t.Errorf("labels[0].Title = %q, want urgent", labels[0].Title)
	}
	if labels[1].ID != 2 {
		t.Errorf("labels[1].ID = %d, want 2", labels[1].ID)
	}
}

func TestCreateLabel(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/labels" {
			t.Errorf("path = %q, want /api/v1/accounts/1/labels", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":    10,
			"title": "vip",
			"color": "#0000ff",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	label, err := client.CreateLabel(context.Background(), CreateLabelOpts{
		Title: "vip",
		Color: "#0000ff",
	})
	if err != nil {
		t.Fatalf("CreateLabel error: %v", err)
	}
	if label.ID != 10 {
		t.Errorf("ID = %d, want 10", label.ID)
	}
	if label.Title != "vip" {
		t.Errorf("Title = %q, want vip", label.Title)
	}
	if gotBody["title"] != "vip" {
		t.Errorf("body title = %v, want vip", gotBody["title"])
	}
}

func TestDeleteLabel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/labels/5" {
			t.Errorf("path = %q, want /api/v1/accounts/1/labels/5", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	err := client.DeleteLabel(context.Background(), 5)
	if err != nil {
		t.Fatalf("DeleteLabel error: %v", err)
	}
}

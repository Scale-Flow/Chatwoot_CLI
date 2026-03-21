package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListPortals(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/portals" {
			t.Errorf("path = %q, want /api/v1/accounts/1/portals", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"payload": []map[string]any{
				{"id": 1, "name": "Support Portal", "slug": "support"},
				{"id": 2, "name": "Dev Docs", "slug": "dev-docs"},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	portals, err := client.ListPortals(context.Background())
	if err != nil {
		t.Fatalf("ListPortals error: %v", err)
	}
	if len(portals) != 2 {
		t.Fatalf("len = %d, want 2", len(portals))
	}
	if portals[0].ID != 1 {
		t.Errorf("ID = %d, want 1", portals[0].ID)
	}
	if portals[0].Name != "Support Portal" {
		t.Errorf("Name = %q, want Support Portal", portals[0].Name)
	}
	if portals[0].Slug != "support" {
		t.Errorf("Slug = %q, want support", portals[0].Slug)
	}
}

func TestCreateArticle(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/portals/7/articles" {
			t.Errorf("path = %q, want /api/v1/accounts/1/portals/7/articles", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":      42,
			"title":   "Getting Started",
			"content": "Welcome to the help center.",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	article, err := client.CreateArticle(context.Background(), 7, CreateArticleOpts{
		Title:   "Getting Started",
		Content: "Welcome to the help center.",
	})
	if err != nil {
		t.Fatalf("CreateArticle error: %v", err)
	}
	if article.ID != 42 {
		t.Errorf("ID = %d, want 42", article.ID)
	}
	if article.Title != "Getting Started" {
		t.Errorf("Title = %q, want Getting Started", article.Title)
	}
	if gotBody["title"] != "Getting Started" {
		t.Errorf("body title = %v, want Getting Started", gotBody["title"])
	}
}

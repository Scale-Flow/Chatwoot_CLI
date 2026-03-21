package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestGetAccount(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1" {
			t.Errorf("path = %q, want /api/v1/accounts/1", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":     1,
			"name":   "Acme Corp",
			"locale": "en",
			"domain": "acme.example.com",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	account, err := client.GetAccount(context.Background())
	if err != nil {
		t.Fatalf("GetAccount error: %v", err)
	}
	if account.ID != 1 {
		t.Errorf("ID = %d, want 1", account.ID)
	}
	if account.Name != "Acme Corp" {
		t.Errorf("Name = %q, want Acme Corp", account.Name)
	}
	if account.Domain != "acme.example.com" {
		t.Errorf("Domain = %q, want acme.example.com", account.Domain)
	}
}

func TestUpdateAccount(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1" {
			t.Errorf("path = %q, want /api/v1/accounts/1", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":     1,
			"name":   "Acme Updated",
			"locale": "fr",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	name := "Acme Updated"
	locale := "fr"
	account, err := client.UpdateAccount(context.Background(), UpdateAccountOpts{
		Name:   &name,
		Locale: &locale,
	})
	if err != nil {
		t.Fatalf("UpdateAccount error: %v", err)
	}
	if account.Name != "Acme Updated" {
		t.Errorf("Name = %q, want Acme Updated", account.Name)
	}
	if gotBody["name"] != "Acme Updated" {
		t.Errorf("body name = %v, want Acme Updated", gotBody["name"])
	}
	if gotBody["locale"] != "fr" {
		t.Errorf("body locale = %v, want fr", gotBody["locale"])
	}
}

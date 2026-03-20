package platform

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
		if r.URL.Path != "/platform/api/v1/accounts/1" {
			t.Errorf("path = %q, want /platform/api/v1/accounts/1", r.URL.Path)
		}
		if r.Header.Get("api_access_token") != "pk-test" {
			t.Errorf("auth header missing or wrong")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":   1,
			"name": "Test Account",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	account, err := client.GetAccount(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetAccount error: %v", err)
	}
	if account.ID != 1 {
		t.Errorf("ID = %d, want 1", account.ID)
	}
	if account.Name != "Test Account" {
		t.Errorf("Name = %q, want %q", account.Name, "Test Account")
	}
}

func TestCreateAccount(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/platform/api/v1/accounts" {
			t.Errorf("path = %q", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		json.NewEncoder(w).Encode(map[string]any{
			"id":   99,
			"name": body["name"],
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	account, err := client.CreateAccount(context.Background(), CreateAccountOpts{Name: "New Account"})
	if err != nil {
		t.Fatalf("CreateAccount error: %v", err)
	}
	if account.ID != 99 {
		t.Errorf("ID = %d, want 99", account.ID)
	}
}

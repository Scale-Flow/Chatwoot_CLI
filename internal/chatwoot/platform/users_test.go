package platform

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestCreateUser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/platform/api/v1/users" {
			t.Errorf("path = %q, want /platform/api/v1/users", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] == nil {
			t.Errorf("request body missing name field")
		}
		if body["email"] == nil {
			t.Errorf("request body missing email field")
		}
		if body["password"] == nil {
			t.Errorf("request body missing password field")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":    1,
			"name":  "Alice",
			"email": "alice@test.com",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	user, err := client.CreateUser(context.Background(), CreateUserOpts{
		Name:     "Alice",
		Email:    "alice@test.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("CreateUser error: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("ID = %d, want 1", user.ID)
	}
	if user.Name != "Alice" {
		t.Errorf("Name = %q, want Alice", user.Name)
	}
	if user.Email != "alice@test.com" {
		t.Errorf("Email = %q, want alice@test.com", user.Email)
	}
}

func TestGetUser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/platform/api/v1/users/1" {
			t.Errorf("path = %q, want /platform/api/v1/users/1", r.URL.Path)
		}
		if r.Header.Get("api_access_token") != "pk-test" {
			t.Errorf("auth header missing or wrong")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":    1,
			"name":  "Alice",
			"email": "alice@test.com",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	user, err := client.GetUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetUser error: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("ID = %d, want 1", user.ID)
	}
	if user.Name != "Alice" {
		t.Errorf("Name = %q, want Alice", user.Name)
	}
}

func TestDeleteUser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/platform/api/v1/users/5" {
			t.Errorf("path = %q, want /platform/api/v1/users/5", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	err := client.DeleteUser(context.Background(), 5)
	if err != nil {
		t.Fatalf("DeleteUser error: %v", err)
	}
}

func TestGetUserSSOLink(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/platform/api/v1/users/1/login" {
			t.Errorf("path = %q, want /platform/api/v1/users/1/login", r.URL.Path)
		}
		if r.Header.Get("api_access_token") != "pk-test" {
			t.Errorf("auth header missing or wrong")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"url": "https://chatwoot.example.com/auth/sign_in?sso_token=abc123",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	link, err := client.GetUserSSOLink(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetUserSSOLink error: %v", err)
	}
	if link.URL == "" {
		t.Errorf("URL is empty, want a non-empty SSO URL")
	}
}

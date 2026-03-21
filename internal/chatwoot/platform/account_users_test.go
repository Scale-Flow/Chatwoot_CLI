package platform

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListAccountUsers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/platform/api/v1/accounts/1/account_users" {
			t.Errorf("path = %q, want /platform/api/v1/accounts/1/account_users", r.URL.Path)
		}
		if r.Header.Get("api_access_token") != "pk-test" {
			t.Errorf("auth header missing or wrong")
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"account_id": 1, "user_id": 10, "role": "administrator"},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	users, err := client.ListAccountUsers(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListAccountUsers error: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("len(users) = %d, want 1", len(users))
	}
	if users[0].AccountID != 1 {
		t.Errorf("AccountID = %d, want 1", users[0].AccountID)
	}
	if users[0].UserID != 10 {
		t.Errorf("UserID = %d, want 10", users[0].UserID)
	}
	if users[0].Role != "administrator" {
		t.Errorf("Role = %q, want administrator", users[0].Role)
	}
}

func TestCreateAccountUser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/platform/api/v1/accounts/1/account_users" {
			t.Errorf("path = %q, want /platform/api/v1/accounts/1/account_users", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["user_id"] == nil {
			t.Errorf("request body missing user_id field")
		}
		if body["role"] == nil {
			t.Errorf("request body missing role field")
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"account_id": 1,
			"user_id":    int(body["user_id"].(float64)),
			"role":       body["role"],
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	err := client.CreateAccountUser(context.Background(), 1, CreateAccountUserOpts{
		UserID: 10,
		Role:   "administrator",
	})
	if err != nil {
		t.Fatalf("CreateAccountUser error: %v", err)
	}
}

func TestDeleteAccountUser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/platform/api/v1/accounts/1/account_users" {
			t.Errorf("path = %q, want /platform/api/v1/accounts/1/account_users", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["user_id"] == nil {
			t.Errorf("request body missing user_id field")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	err := client.DeleteAccountUser(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("DeleteAccountUser error: %v", err)
	}
}

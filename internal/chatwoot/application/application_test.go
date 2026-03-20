package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestGetProfile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/profile" {
			t.Errorf("path = %q, want /api/v1/profile", r.URL.Path)
		}
		if r.Header.Get("api_access_token") != "sk-test" {
			t.Errorf("auth header missing or wrong")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":    1,
			"name":  "Test Agent",
			"email": "agent@test.com",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	profile, err := client.GetProfile(context.Background())
	if err != nil {
		t.Fatalf("GetProfile error: %v", err)
	}
	if profile.ID != 1 {
		t.Errorf("ID = %d, want 1", profile.ID)
	}
	if profile.Name != "Test Agent" {
		t.Errorf("Name = %q, want %q", profile.Name, "Test Agent")
	}
}

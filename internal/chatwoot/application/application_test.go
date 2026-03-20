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

func TestListConversations(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1/conversations" {
			t.Errorf("path = %q, want /api/v1/accounts/1/conversations", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("page = %q, want 1", r.URL.Query().Get("page"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"payload": []map[string]any{
					{"id": 1, "status": "open"},
					{"id": 2, "status": "resolved"},
				},
				"meta": map[string]any{
					"all_count": 42,
				},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	convos, err := client.ListConversations(context.Background(), ListConversationsOpts{Page: 1})
	if err != nil {
		t.Fatalf("ListConversations error: %v", err)
	}
	if len(convos) != 2 {
		t.Errorf("len = %d, want 2", len(convos))
	}
	if convos[0].ID != 1 {
		t.Errorf("convos[0].ID = %d, want 1", convos[0].ID)
	}
}

func TestUpdateProfile(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		if r.URL.Path != "/api/v1/profile" {
			t.Errorf("path = %q, want /api/v1/profile", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":    1,
			"name":  "Updated Name",
			"email": "agent@test.com",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	name := "Updated Name"
	profile, err := client.UpdateProfile(context.Background(), UpdateProfileOpts{Name: &name})
	if err != nil {
		t.Fatalf("UpdateProfile error: %v", err)
	}
	if gotMethod != "PATCH" {
		t.Errorf("method = %q, want PATCH", gotMethod)
	}
	if gotBody["name"] != "Updated Name" {
		t.Errorf("body name = %v, want Updated Name", gotBody["name"])
	}
	if profile.Name != "Updated Name" {
		t.Errorf("Name = %q, want %q", profile.Name, "Updated Name")
	}
}

func TestGetConversation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1/conversations/42" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":     42,
			"status": "open",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	convo, err := client.GetConversation(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetConversation error: %v", err)
	}
	if convo.ID != 42 {
		t.Errorf("ID = %d, want 42", convo.ID)
	}
}

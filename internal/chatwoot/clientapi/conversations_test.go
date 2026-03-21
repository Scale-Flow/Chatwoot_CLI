package clientapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListConversations(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/api/v1/inboxes/inbox-abc/contacts/contact-xyz/conversations" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "status": "open"},
			{"id": 2, "status": "resolved"},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "", "")
	client := NewClient(transport, "inbox-abc")

	conversations, err := client.ListConversations(context.Background(), "contact-xyz")
	if err != nil {
		t.Fatalf("ListConversations error: %v", err)
	}
	if len(conversations) != 2 {
		t.Errorf("len = %d, want 2", len(conversations))
	}
}

func TestCreateConversation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/api/v1/inboxes/inbox-abc/contacts/contact-xyz/conversations" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{"id": 5, "status": "open"})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "", "")
	client := NewClient(transport, "inbox-abc")

	conv, err := client.CreateConversation(context.Background(), "contact-xyz")
	if err != nil {
		t.Fatalf("CreateConversation error: %v", err)
	}
	if conv.ID != 5 {
		t.Errorf("ID = %d, want 5", conv.ID)
	}
}

func TestToggleStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/api/v1/inboxes/inbox-abc/contacts/contact-xyz/conversations/5/toggle_status" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{"id": 5, "status": "resolved"})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "", "")
	client := NewClient(transport, "inbox-abc")

	conv, err := client.ToggleStatus(context.Background(), "contact-xyz", 5)
	if err != nil {
		t.Fatalf("ToggleStatus error: %v", err)
	}
	if conv.Status != "resolved" {
		t.Errorf("Status = %q, want resolved", conv.Status)
	}
}

func TestToggleTyping(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/api/v1/inboxes/inbox-abc/contacts/contact-xyz/conversations/5/toggle_typing" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if _, ok := body["typing_status"]; !ok {
			t.Errorf("body missing typing_status field")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "", "")
	client := NewClient(transport, "inbox-abc")

	err := client.ToggleTyping(context.Background(), "contact-xyz", 5, ToggleTypingOpts{TypingStatus: "on"})
	if err != nil {
		t.Fatalf("ToggleTyping error: %v", err)
	}
}

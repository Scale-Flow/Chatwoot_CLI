package clientapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListMessages(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wantPath := "/public/api/v1/inboxes/inbox-abc/contacts/contact-xyz/conversations/5/messages"
		if r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "content": "Hello"},
			{"id": 2, "content": "World"},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "", "")
	client := NewClient(transport, "inbox-abc")

	messages, err := client.ListMessages(context.Background(), "contact-xyz", 5)
	if err != nil {
		t.Fatalf("ListMessages error: %v", err)
	}
	if len(messages) != 2 {
		t.Errorf("len = %d, want 2", len(messages))
	}
}

func TestCreateMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wantPath := "/public/api/v1/inboxes/inbox-abc/contacts/contact-xyz/conversations/5/messages"
		if r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if _, ok := body["content"]; !ok {
			t.Errorf("body missing content field")
		}
		json.NewEncoder(w).Encode(map[string]any{"id": 10, "content": "Hello"})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "", "")
	client := NewClient(transport, "inbox-abc")

	msg, err := client.CreateMessage(context.Background(), "contact-xyz", 5, CreateMessageOpts{Content: "Hello"})
	if err != nil {
		t.Fatalf("CreateMessage error: %v", err)
	}
	if msg.ID != 10 {
		t.Errorf("ID = %d, want 10", msg.ID)
	}
}

package clientapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestCreateContact(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/api/v1/inboxes/inbox-abc/contacts" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"source_id":    "contact-xyz",
			"name":         "Test User",
			"pubsub_token": "token-123",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "", "")
	client := NewClient(transport, "inbox-abc")

	contact, err := client.CreateContact(context.Background(), CreateContactOpts{Name: "Test User"})
	if err != nil {
		t.Fatalf("CreateContact error: %v", err)
	}
	if contact.SourceID != "contact-xyz" {
		t.Errorf("SourceID = %q, want %q", contact.SourceID, "contact-xyz")
	}
}

func TestGetContact(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/api/v1/inboxes/inbox-abc/contacts/contact-xyz" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"source_id": "contact-xyz",
			"name":      "Existing User",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "", "")
	client := NewClient(transport, "inbox-abc")

	contact, err := client.GetContact(context.Background(), "contact-xyz")
	if err != nil {
		t.Fatalf("GetContact error: %v", err)
	}
	if contact.Name != "Existing User" {
		t.Errorf("Name = %q, want %q", contact.Name, "Existing User")
	}
}

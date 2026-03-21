package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestMessageList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/conversations/5/messages" {
			t.Errorf("path = %q, want /api/v1/accounts/1/conversations/5/messages", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"payload": []map[string]any{
				{"id": 1, "content": "Hello", "message_type": 0, "conversation_id": 5},
				{"id": 2, "content": "Hi there", "message_type": 1, "conversation_id": 5},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	messages, err := client.ListMessages(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListMessages error: %v", err)
	}
	if len(messages) != 2 {
		t.Errorf("len = %d, want 2", len(messages))
	}
	if messages[0].ID != 1 {
		t.Errorf("messages[0].ID = %d, want 1", messages[0].ID)
	}
	if messages[0].Content != "Hello" {
		t.Errorf("messages[0].Content = %q, want Hello", messages[0].Content)
	}
	if messages[1].ConversationID != 5 {
		t.Errorf("messages[1].ConversationID = %d, want 5", messages[1].ConversationID)
	}
}

func TestMessageCreate(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/conversations/5/messages" {
			t.Errorf("path = %q, want /api/v1/accounts/1/conversations/5/messages", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":              100,
			"content":         "Test message",
			"message_type":    0,
			"private":         true,
			"conversation_id": 5,
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	msg, err := client.CreateMessage(context.Background(), 5, CreateMessageOpts{
		Content:     "Test message",
		MessageType: "outgoing",
		Private:     true,
	})
	if err != nil {
		t.Fatalf("CreateMessage error: %v", err)
	}
	if msg.ID != 100 {
		t.Errorf("ID = %d, want 100", msg.ID)
	}
	if msg.Content != "Test message" {
		t.Errorf("Content = %q, want Test message", msg.Content)
	}
	if !msg.Private {
		t.Error("Private = false, want true")
	}
	// Verify request body
	if gotBody["content"] != "Test message" {
		t.Errorf("body content = %v, want Test message", gotBody["content"])
	}
	if gotBody["message_type"] != "outgoing" {
		t.Errorf("body message_type = %v, want outgoing", gotBody["message_type"])
	}
	if gotBody["private"] != true {
		t.Errorf("body private = %v, want true", gotBody["private"])
	}
}

func TestMessageDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/conversations/5/messages/100" {
			t.Errorf("path = %q, want /api/v1/accounts/1/conversations/5/messages/100", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	err := client.DeleteMessage(context.Background(), 5, 100)
	if err != nil {
		t.Fatalf("DeleteMessage error: %v", err)
	}
}

package platform

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListAgentBots(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/platform/api/v1/agent_bots" {
			t.Errorf("path = %q, want /platform/api/v1/agent_bots", r.URL.Path)
		}
		if r.Header.Get("api_access_token") != "pk-test" {
			t.Errorf("auth header missing or wrong")
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "Helper Bot"},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	bots, err := client.ListAgentBots(context.Background())
	if err != nil {
		t.Fatalf("ListAgentBots error: %v", err)
	}
	if len(bots) != 1 {
		t.Fatalf("len(bots) = %d, want 1", len(bots))
	}
	if bots[0].ID != 1 {
		t.Errorf("ID = %d, want 1", bots[0].ID)
	}
	if bots[0].Name != "Helper Bot" {
		t.Errorf("Name = %q, want %q", bots[0].Name, "Helper Bot")
	}
}

func TestGetAgentBot(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/platform/api/v1/agent_bots/1" {
			t.Errorf("path = %q, want /platform/api/v1/agent_bots/1", r.URL.Path)
		}
		if r.Header.Get("api_access_token") != "pk-test" {
			t.Errorf("auth header missing or wrong")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":   1,
			"name": "Helper Bot",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	bot, err := client.GetAgentBot(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetAgentBot error: %v", err)
	}
	if bot.ID != 1 {
		t.Errorf("ID = %d, want 1", bot.ID)
	}
	if bot.Name != "Helper Bot" {
		t.Errorf("Name = %q, want %q", bot.Name, "Helper Bot")
	}
}

func TestCreateAgentBot(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/platform/api/v1/agent_bots" {
			t.Errorf("path = %q, want /platform/api/v1/agent_bots", r.URL.Path)
		}
		if r.Header.Get("api_access_token") != "pk-test" {
			t.Errorf("auth header missing or wrong")
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] == nil {
			t.Errorf("request body missing name field")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":   42,
			"name": body["name"],
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	bot, err := client.CreateAgentBot(context.Background(), CreateAgentBotOpts{Name: "New Bot"})
	if err != nil {
		t.Fatalf("CreateAgentBot error: %v", err)
	}
	if bot.ID != 42 {
		t.Errorf("ID = %d, want 42", bot.ID)
	}
	if bot.Name != "New Bot" {
		t.Errorf("Name = %q, want %q", bot.Name, "New Bot")
	}
}

func TestDeleteAgentBot(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/platform/api/v1/agent_bots/7" {
			t.Errorf("path = %q, want /platform/api/v1/agent_bots/7", r.URL.Path)
		}
		if r.Header.Get("api_access_token") != "pk-test" {
			t.Errorf("auth header missing or wrong")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	err := client.DeleteAgentBot(context.Background(), 7)
	if err != nil {
		t.Fatalf("DeleteAgentBot error: %v", err)
	}
}

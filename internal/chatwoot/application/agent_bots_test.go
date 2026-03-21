package application

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
		if r.URL.Path != "/api/v1/accounts/1/agent_bots" {
			t.Errorf("path = %q, want /api/v1/accounts/1/agent_bots", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "SupportBot", "bot_type": "webhook"},
			{"id": 2, "name": "SalesBot", "bot_type": "webhook"},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	bots, err := client.ListAgentBots(context.Background())
	if err != nil {
		t.Fatalf("ListAgentBots error: %v", err)
	}
	if len(bots) != 2 {
		t.Errorf("len = %d, want 2", len(bots))
	}
	if bots[0].Name != "SupportBot" {
		t.Errorf("bots[0].Name = %q, want SupportBot", bots[0].Name)
	}
	if bots[1].ID != 2 {
		t.Errorf("bots[1].ID = %d, want 2", bots[1].ID)
	}
}

func TestCreateAgentBot(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/agent_bots" {
			t.Errorf("path = %q, want /api/v1/accounts/1/agent_bots", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":           40,
			"name":         "TriageBot",
			"bot_type":     "webhook",
			"outgoing_url": "https://bot.example.com/hook",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	bot, err := client.CreateAgentBot(context.Background(), CreateAgentBotOpts{
		Name:        "TriageBot",
		BotType:     "webhook",
		OutgoingURL: "https://bot.example.com/hook",
	})
	if err != nil {
		t.Fatalf("CreateAgentBot error: %v", err)
	}
	if bot.ID != 40 {
		t.Errorf("ID = %d, want 40", bot.ID)
	}
	if bot.Name != "TriageBot" {
		t.Errorf("Name = %q, want TriageBot", bot.Name)
	}
	if gotBody["name"] != "TriageBot" {
		t.Errorf("body name = %v, want TriageBot", gotBody["name"])
	}
	if gotBody["bot_type"] != "webhook" {
		t.Errorf("body bot_type = %v, want webhook", gotBody["bot_type"])
	}
}

package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListAgents(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/agents" {
			t.Errorf("path = %q, want /api/v1/accounts/1/agents", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "Alice", "email": "alice@test.com"},
			{"id": 2, "name": "Bob", "email": "bob@test.com"},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	agents, err := client.ListAgents(context.Background())
	if err != nil {
		t.Fatalf("ListAgents error: %v", err)
	}
	if len(agents) != 2 {
		t.Errorf("len = %d, want 2", len(agents))
	}
	if agents[0].Name != "Alice" {
		t.Errorf("agents[0].Name = %q, want Alice", agents[0].Name)
	}
	if agents[0].Email != "alice@test.com" {
		t.Errorf("agents[0].Email = %q, want alice@test.com", agents[0].Email)
	}
	if agents[1].ID != 2 {
		t.Errorf("agents[1].ID = %d, want 2", agents[1].ID)
	}
	if agents[1].Name != "Bob" {
		t.Errorf("agents[1].Name = %q, want Bob", agents[1].Name)
	}
}

func TestCreateAgent(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/agents" {
			t.Errorf("path = %q, want /api/v1/accounts/1/agents", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":    10,
			"name":  "Carol",
			"email": "carol@test.com",
			"role":  "agent",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	agent, err := client.CreateAgent(context.Background(), CreateAgentOpts{
		Name:  "Carol",
		Email: "carol@test.com",
		Role:  "agent",
	})
	if err != nil {
		t.Fatalf("CreateAgent error: %v", err)
	}
	if agent.ID != 10 {
		t.Errorf("ID = %d, want 10", agent.ID)
	}
	if agent.Name != "Carol" {
		t.Errorf("Name = %q, want Carol", agent.Name)
	}
	if gotBody["name"] != "Carol" {
		t.Errorf("body name = %v, want Carol", gotBody["name"])
	}
	if gotBody["email"] != "carol@test.com" {
		t.Errorf("body email = %v, want carol@test.com", gotBody["email"])
	}
	if gotBody["role"] != "agent" {
		t.Errorf("body role = %v, want agent", gotBody["role"])
	}
}

func TestDeleteAgent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/agents/7" {
			t.Errorf("path = %q, want /api/v1/accounts/1/agents/7", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	err := client.DeleteAgent(context.Background(), 7)
	if err != nil {
		t.Fatalf("DeleteAgent error: %v", err)
	}
}

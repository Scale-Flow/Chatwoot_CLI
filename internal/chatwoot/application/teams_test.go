package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListTeams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/teams" {
			t.Errorf("path = %q, want /api/v1/accounts/1/teams", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "Support", "account_id": 1},
			{"id": 2, "name": "Sales", "account_id": 1},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	teams, err := client.ListTeams(context.Background())
	if err != nil {
		t.Fatalf("ListTeams error: %v", err)
	}
	if len(teams) != 2 {
		t.Errorf("len = %d, want 2", len(teams))
	}
	if teams[0].Name != "Support" {
		t.Errorf("teams[0].Name = %q, want Support", teams[0].Name)
	}
	if teams[1].ID != 2 {
		t.Errorf("teams[1].ID = %d, want 2", teams[1].ID)
	}
}

func TestGetTeam(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/teams/5" {
			t.Errorf("path = %q, want /api/v1/accounts/1/teams/5", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":          5,
			"name":        "Engineering",
			"description": "Engineering team",
			"account_id":  1,
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	team, err := client.GetTeam(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetTeam error: %v", err)
	}
	if team.ID != 5 {
		t.Errorf("ID = %d, want 5", team.ID)
	}
	if team.Name != "Engineering" {
		t.Errorf("Name = %q, want Engineering", team.Name)
	}
	if team.Description != "Engineering team" {
		t.Errorf("Description = %q, want Engineering team", team.Description)
	}
}

func TestCreateTeam(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/teams" {
			t.Errorf("path = %q, want /api/v1/accounts/1/teams", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":         10,
			"name":       "New Team",
			"account_id": 1,
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	team, err := client.CreateTeam(context.Background(), CreateTeamOpts{Name: "New Team"})
	if err != nil {
		t.Fatalf("CreateTeam error: %v", err)
	}
	if team.ID != 10 {
		t.Errorf("ID = %d, want 10", team.ID)
	}
	if gotBody["name"] != "New Team" {
		t.Errorf("body name = %v, want New Team", gotBody["name"])
	}
}

func TestUpdateTeam(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/teams/3" {
			t.Errorf("path = %q, want /api/v1/accounts/1/teams/3", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":         3,
			"name":       "Updated Team",
			"account_id": 1,
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	name := "Updated Team"
	team, err := client.UpdateTeam(context.Background(), 3, UpdateTeamOpts{Name: &name})
	if err != nil {
		t.Fatalf("UpdateTeam error: %v", err)
	}
	if team.ID != 3 {
		t.Errorf("ID = %d, want 3", team.ID)
	}
	if gotBody["name"] != "Updated Team" {
		t.Errorf("body name = %v, want Updated Team", gotBody["name"])
	}
}

func TestDeleteTeam(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/teams/7" {
			t.Errorf("path = %q, want /api/v1/accounts/1/teams/7", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	err := client.DeleteTeam(context.Background(), 7)
	if err != nil {
		t.Fatalf("DeleteTeam error: %v", err)
	}
}

func TestListTeamMembers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/teams/4/team_members" {
			t.Errorf("path = %q, want /api/v1/accounts/1/teams/4/team_members", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "Alice", "email": "alice@test.com", "role": "agent"},
			{"id": 2, "name": "Bob", "email": "bob@test.com", "role": "agent"},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	agents, err := client.ListTeamMembers(context.Background(), 4)
	if err != nil {
		t.Fatalf("ListTeamMembers error: %v", err)
	}
	if len(agents) != 2 {
		t.Errorf("len = %d, want 2", len(agents))
	}
	if agents[0].Name != "Alice" {
		t.Errorf("agents[0].Name = %q, want Alice", agents[0].Name)
	}
	if agents[1].Email != "bob@test.com" {
		t.Errorf("agents[1].Email = %q, want bob@test.com", agents[1].Email)
	}
}

func TestAddTeamMember(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/teams/4/team_members" {
			t.Errorf("path = %q, want /api/v1/accounts/1/teams/4/team_members", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "Alice", "email": "alice@test.com"},
			{"id": 3, "name": "Charlie", "email": "charlie@test.com"},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	agents, err := client.AddTeamMember(context.Background(), 4, []int{1, 3})
	if err != nil {
		t.Fatalf("AddTeamMember error: %v", err)
	}
	if len(agents) != 2 {
		t.Errorf("len = %d, want 2", len(agents))
	}
	userIDs, ok := gotBody["user_ids"].([]any)
	if !ok {
		t.Fatal("body user_ids not an array")
	}
	if len(userIDs) != 2 {
		t.Errorf("user_ids len = %d, want 2", len(userIDs))
	}
	if userIDs[0] != float64(1) {
		t.Errorf("user_ids[0] = %v, want 1", userIDs[0])
	}
	if userIDs[1] != float64(3) {
		t.Errorf("user_ids[1] = %v, want 3", userIDs[1])
	}
}

func TestUpdateTeamMembers(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/teams/4/team_members" {
			t.Errorf("path = %q, want /api/v1/accounts/1/teams/4/team_members", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 2, "name": "Bob", "email": "bob@test.com"},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	agents, err := client.UpdateTeamMembers(context.Background(), 4, []int{2})
	if err != nil {
		t.Fatalf("UpdateTeamMembers error: %v", err)
	}
	if len(agents) != 1 {
		t.Errorf("len = %d, want 1", len(agents))
	}
	if agents[0].Name != "Bob" {
		t.Errorf("agents[0].Name = %q, want Bob", agents[0].Name)
	}
	userIDs, ok := gotBody["user_ids"].([]any)
	if !ok {
		t.Fatal("body user_ids not an array")
	}
	if len(userIDs) != 1 {
		t.Errorf("user_ids len = %d, want 1", len(userIDs))
	}
	if userIDs[0] != float64(2) {
		t.Errorf("user_ids[0] = %v, want 2", userIDs[0])
	}
}

func TestRemoveTeamMember(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/teams/4/team_members" {
			t.Errorf("path = %q, want /api/v1/accounts/1/teams/4/team_members", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	err := client.RemoveTeamMember(context.Background(), 4, []int{1, 2})
	if err != nil {
		t.Fatalf("RemoveTeamMember error: %v", err)
	}
	userIDs, ok := gotBody["user_ids"].([]any)
	if !ok {
		t.Fatal("body user_ids not an array")
	}
	if len(userIDs) != 2 {
		t.Errorf("user_ids len = %d, want 2", len(userIDs))
	}
	if userIDs[0] != float64(1) {
		t.Errorf("user_ids[0] = %v, want 1", userIDs[0])
	}
	if userIDs[1] != float64(2) {
		t.Errorf("user_ids[1] = %v, want 2", userIDs[1])
	}
}

package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListInboxes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/inboxes" {
			t.Errorf("path = %q, want /api/v1/accounts/1/inboxes", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"payload": []map[string]any{
				{"id": 1, "name": "Email Inbox", "channel_type": "Channel::Email"},
				{"id": 2, "name": "Web Chat", "channel_type": "Channel::WebWidget"},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	inboxes, err := client.ListInboxes(context.Background())
	if err != nil {
		t.Fatalf("ListInboxes error: %v", err)
	}
	if len(inboxes) != 2 {
		t.Errorf("len = %d, want 2", len(inboxes))
	}
	if inboxes[0].Name != "Email Inbox" {
		t.Errorf("inboxes[0].Name = %q, want Email Inbox", inboxes[0].Name)
	}
	if inboxes[1].ChannelType != "Channel::WebWidget" {
		t.Errorf("inboxes[1].ChannelType = %q, want Channel::WebWidget", inboxes[1].ChannelType)
	}
}

func TestGetInbox(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/inboxes/5" {
			t.Errorf("path = %q, want /api/v1/accounts/1/inboxes/5", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":                     5,
			"name":                   "Support Inbox",
			"channel_type":           "Channel::Email",
			"enable_auto_assignment": true,
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	inbox, err := client.GetInbox(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetInbox error: %v", err)
	}
	if inbox.ID != 5 {
		t.Errorf("ID = %d, want 5", inbox.ID)
	}
	if inbox.Name != "Support Inbox" {
		t.Errorf("Name = %q, want Support Inbox", inbox.Name)
	}
	if !inbox.EnableAutoAssignment {
		t.Error("EnableAutoAssignment = false, want true")
	}
}

func TestCreateInbox(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/inboxes" {
			t.Errorf("path = %q, want /api/v1/accounts/1/inboxes", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":           10,
			"name":         "New Inbox",
			"channel_type": "Channel::WebWidget",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	inbox, err := client.CreateInbox(context.Background(), CreateInboxOpts{
		Name:    "New Inbox",
		Channel: map[string]string{"type": "web_widget", "website_url": "https://example.com"},
	})
	if err != nil {
		t.Fatalf("CreateInbox error: %v", err)
	}
	if inbox.ID != 10 {
		t.Errorf("ID = %d, want 10", inbox.ID)
	}
	if gotBody["name"] != "New Inbox" {
		t.Errorf("body name = %v, want New Inbox", gotBody["name"])
	}
	if gotBody["channel"] == nil {
		t.Error("body channel is nil")
	}
}

func TestUpdateInbox(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/inboxes/5" {
			t.Errorf("path = %q, want /api/v1/accounts/1/inboxes/5", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":                     5,
			"name":                   "Renamed Inbox",
			"enable_auto_assignment": false,
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	name := "Renamed Inbox"
	autoAssign := false
	inbox, err := client.UpdateInbox(context.Background(), 5, UpdateInboxOpts{
		Name:                 &name,
		EnableAutoAssignment: &autoAssign,
	})
	if err != nil {
		t.Fatalf("UpdateInbox error: %v", err)
	}
	if inbox.ID != 5 {
		t.Errorf("ID = %d, want 5", inbox.ID)
	}
	if inbox.Name != "Renamed Inbox" {
		t.Errorf("Name = %q, want Renamed Inbox", inbox.Name)
	}
	if gotBody["name"] != "Renamed Inbox" {
		t.Errorf("body name = %v, want Renamed Inbox", gotBody["name"])
	}
}

func TestListInboxMembers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/inbox_members/5" {
			t.Errorf("path = %q, want /api/v1/accounts/1/inbox_members/5", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"payload": []map[string]any{
				{"id": 1, "name": "Agent One", "email": "one@test.com", "role": "agent"},
				{"id": 2, "name": "Agent Two", "email": "two@test.com", "role": "agent"},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	agents, err := client.ListInboxMembers(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListInboxMembers error: %v", err)
	}
	if len(agents) != 2 {
		t.Errorf("len = %d, want 2", len(agents))
	}
	if agents[0].Name != "Agent One" {
		t.Errorf("agents[0].Name = %q, want Agent One", agents[0].Name)
	}
	if agents[1].Email != "two@test.com" {
		t.Errorf("agents[1].Email = %q, want two@test.com", agents[1].Email)
	}
}

func TestAddInboxMember(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/inbox_members" {
			t.Errorf("path = %q, want /api/v1/accounts/1/inbox_members", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"payload": []map[string]any{
				{"id": 1, "name": "Agent One", "email": "one@test.com"},
				{"id": 3, "name": "Agent Three", "email": "three@test.com"},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	agents, err := client.AddInboxMember(context.Background(), 5, []int{1, 3})
	if err != nil {
		t.Fatalf("AddInboxMember error: %v", err)
	}
	if len(agents) != 2 {
		t.Errorf("len = %d, want 2", len(agents))
	}
	if gotBody["inbox_id"] != float64(5) {
		t.Errorf("body inbox_id = %v, want 5", gotBody["inbox_id"])
	}
	userIDs, ok := gotBody["user_ids"].([]any)
	if !ok {
		t.Fatal("body user_ids not an array")
	}
	if len(userIDs) != 2 {
		t.Errorf("body user_ids len = %d, want 2", len(userIDs))
	}
}

func TestUpdateInboxMembers(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/inbox_members" {
			t.Errorf("path = %q, want /api/v1/accounts/1/inbox_members", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"payload": []map[string]any{
				{"id": 2, "name": "Agent Two", "email": "two@test.com"},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	agents, err := client.UpdateInboxMembers(context.Background(), 5, []int{2})
	if err != nil {
		t.Fatalf("UpdateInboxMembers error: %v", err)
	}
	if len(agents) != 1 {
		t.Errorf("len = %d, want 1", len(agents))
	}
	if agents[0].Name != "Agent Two" {
		t.Errorf("agents[0].Name = %q, want Agent Two", agents[0].Name)
	}
	if gotBody["inbox_id"] != float64(5) {
		t.Errorf("body inbox_id = %v, want 5", gotBody["inbox_id"])
	}
	userIDs, ok := gotBody["user_ids"].([]any)
	if !ok {
		t.Fatal("body user_ids not an array")
	}
	if len(userIDs) != 1 {
		t.Errorf("body user_ids len = %d, want 1", len(userIDs))
	}
}

func TestRemoveInboxMember(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/inbox_members" {
			t.Errorf("path = %q, want /api/v1/accounts/1/inbox_members", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	err := client.RemoveInboxMember(context.Background(), 5, []int{1, 3})
	if err != nil {
		t.Fatalf("RemoveInboxMember error: %v", err)
	}
	if gotBody["inbox_id"] != float64(5) {
		t.Errorf("body inbox_id = %v, want 5", gotBody["inbox_id"])
	}
	userIDs, ok := gotBody["user_ids"].([]any)
	if !ok {
		t.Fatal("body user_ids not an array")
	}
	if len(userIDs) != 2 {
		t.Errorf("body user_ids len = %d, want 2", len(userIDs))
	}
}

func TestGetInboxAgentBot(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/inboxes/5/agent_bot" {
			t.Errorf("path = %q, want /api/v1/accounts/1/inboxes/5/agent_bot", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":   7,
			"name": "Support Bot",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	bot, err := client.GetInboxAgentBot(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetInboxAgentBot error: %v", err)
	}
	if bot.ID != 7 {
		t.Errorf("ID = %d, want 7", bot.ID)
	}
	if bot.Name != "Support Bot" {
		t.Errorf("Name = %q, want Support Bot", bot.Name)
	}
}

func TestSetInboxAgentBot(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/inboxes/5/set_agent_bot" {
			t.Errorf("path = %q, want /api/v1/accounts/1/inboxes/5/set_agent_bot", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":   9,
			"name": "New Bot",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	bot, err := client.SetInboxAgentBot(context.Background(), 5, 9)
	if err != nil {
		t.Fatalf("SetInboxAgentBot error: %v", err)
	}
	if bot.ID != 9 {
		t.Errorf("ID = %d, want 9", bot.ID)
	}
	if bot.Name != "New Bot" {
		t.Errorf("Name = %q, want New Bot", bot.Name)
	}
	if gotBody["agent_bot"] != float64(9) {
		t.Errorf("body agent_bot = %v, want 9", gotBody["agent_bot"])
	}
}

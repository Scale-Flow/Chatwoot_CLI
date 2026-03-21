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

	convos, pg, err := client.ListConversations(context.Background(), ListConversationsOpts{Page: 1})
	if err != nil {
		t.Fatalf("ListConversations error: %v", err)
	}
	if len(convos) != 2 {
		t.Errorf("len = %d, want 2", len(convos))
	}
	if convos[0].ID != 1 {
		t.Errorf("convos[0].ID = %d, want 1", convos[0].ID)
	}
	if pg == nil {
		t.Fatal("pagination is nil")
	}
	if pg.TotalCount != 42 {
		t.Errorf("TotalCount = %d, want 42", pg.TotalCount)
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

func TestCreateConversation(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		if r.URL.Path != "/api/v1/accounts/1/conversations" {
			t.Errorf("path = %q, want /api/v1/accounts/1/conversations", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":     10,
			"status": "open",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	convo, err := client.CreateConversation(context.Background(), CreateConversationOpts{
		ContactID: 5,
		InboxID:   3,
	})
	if err != nil {
		t.Fatalf("CreateConversation error: %v", err)
	}
	if gotMethod != "POST" {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotBody["contact_id"] != float64(5) {
		t.Errorf("contact_id = %v, want 5", gotBody["contact_id"])
	}
	if convo.ID != 10 {
		t.Errorf("ID = %d, want 10", convo.ID)
	}
}

func TestUpdateConversation(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		if r.URL.Path != "/api/v1/accounts/1/conversations/10" {
			t.Errorf("path = %q, want /api/v1/accounts/1/conversations/10", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":     10,
			"status": "resolved",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	status := "resolved"
	convo, err := client.UpdateConversation(context.Background(), 10, UpdateConversationOpts{
		Status: &status,
	})
	if err != nil {
		t.Fatalf("UpdateConversation error: %v", err)
	}
	if gotMethod != "PATCH" {
		t.Errorf("method = %q, want PATCH", gotMethod)
	}
	if gotBody["status"] != "resolved" {
		t.Errorf("status = %v, want resolved", gotBody["status"])
	}
	if convo.Status != "resolved" {
		t.Errorf("Status = %q, want resolved", convo.Status)
	}
}

func TestFilterConversations(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		if r.URL.Path != "/api/v1/accounts/1/conversations/filter" {
			t.Errorf("path = %q, want /api/v1/accounts/1/conversations/filter", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"payload": []map[string]any{
					{"id": 1, "status": "open"},
				},
				"meta": map[string]any{
					"all_count": 1,
				},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	convos, pg, err := client.FilterConversations(context.Background(), FilterConversationsOpts{
		Page:    1,
		Payload: []any{map[string]any{"attribute_key": "status", "filter_operator": "equal_to", "values": []string{"open"}}},
	})
	if err != nil {
		t.Fatalf("FilterConversations error: %v", err)
	}
	if gotMethod != "POST" {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if len(convos) != 1 {
		t.Errorf("len = %d, want 1", len(convos))
	}
	if pg == nil {
		t.Fatal("pagination is nil")
	}
	if pg.TotalCount != 1 {
		t.Errorf("TotalCount = %d, want 1", pg.TotalCount)
	}
}

func TestGetConversationMeta(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1/conversations/meta" {
			t.Errorf("path = %q, want /api/v1/accounts/1/conversations/meta", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("method = %q, want GET", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"all_count":      100,
			"open_count":     40,
			"resolved_count": 50,
			"pending_count":  5,
			"snoozed_count":  5,
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	meta, err := client.GetConversationMeta(context.Background())
	if err != nil {
		t.Fatalf("GetConversationMeta error: %v", err)
	}
	if meta.AllCount != 100 {
		t.Errorf("AllCount = %d, want 100", meta.AllCount)
	}
	if meta.OpenCount != 40 {
		t.Errorf("OpenCount = %d, want 40", meta.OpenCount)
	}
}

func TestToggleConversationStatus(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		if r.URL.Path != "/api/v1/accounts/1/conversations/10/toggle_status" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"payload": map[string]any{
				"success":         true,
				"conversation_id": 10,
				"current_status":  "resolved",
				"snoozed_until":   nil,
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	result, err := client.ToggleConversationStatus(context.Background(), 10, "resolved")
	if err != nil {
		t.Fatalf("ToggleConversationStatus error: %v", err)
	}
	if gotMethod != "POST" {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotBody["status"] != "resolved" {
		t.Errorf("body status = %v, want resolved", gotBody["status"])
	}
	if !result.Success {
		t.Errorf("Success = false, want true")
	}
	if result.ConversationID != 10 {
		t.Errorf("ConversationID = %d, want 10", result.ConversationID)
	}
	if result.CurrentStatus != "resolved" {
		t.Errorf("CurrentStatus = %q, want resolved", result.CurrentStatus)
	}
}

func TestToggleConversationPriority(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		if r.URL.Path != "/api/v1/accounts/1/conversations/10/toggle_priority" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"payload": map[string]any{
				"success":          true,
				"conversation_id":  10,
				"current_priority": "urgent",
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	result, err := client.ToggleConversationPriority(context.Background(), 10, "urgent")
	if err != nil {
		t.Fatalf("ToggleConversationPriority error: %v", err)
	}
	if gotMethod != "POST" {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotBody["priority"] != "urgent" {
		t.Errorf("body priority = %v, want urgent", gotBody["priority"])
	}
	if !result.Success {
		t.Errorf("Success = false, want true")
	}
	if result.ConversationID != 10 {
		t.Errorf("ConversationID = %d, want 10", result.ConversationID)
	}
	if result.CurrentPriority != "urgent" {
		t.Errorf("CurrentPriority = %q, want urgent", result.CurrentPriority)
	}
}

func TestAssignConversation(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		if r.URL.Path != "/api/v1/accounts/1/conversations/10/assignments" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":     10,
			"status": "open",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	agentID := 7
	convo, err := client.AssignConversation(context.Background(), 10, AssignOpts{
		AgentID: &agentID,
	})
	if err != nil {
		t.Fatalf("AssignConversation error: %v", err)
	}
	if gotMethod != "POST" {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotBody["assignee_id"] != float64(7) {
		t.Errorf("assignee_id = %v, want 7", gotBody["assignee_id"])
	}
	if convo.ID != 10 {
		t.Errorf("ID = %d, want 10", convo.ID)
	}
}

func TestListConversationLabels(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1/conversations/10/labels" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("method = %q, want GET", r.Method)
		}
		json.NewEncoder(w).Encode([]string{"bug", "urgent"})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	labels, err := client.ListConversationLabels(context.Background(), 10)
	if err != nil {
		t.Fatalf("ListConversationLabels error: %v", err)
	}
	if len(labels) != 2 {
		t.Errorf("len = %d, want 2", len(labels))
	}
	if labels[0] != "bug" {
		t.Errorf("labels[0] = %q, want bug", labels[0])
	}
}

func TestSetConversationLabels(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		if r.URL.Path != "/api/v1/accounts/1/conversations/10/labels" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode([]string{"feature", "high-priority"})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	labels, err := client.SetConversationLabels(context.Background(), 10, []string{"feature", "high-priority"})
	if err != nil {
		t.Fatalf("SetConversationLabels error: %v", err)
	}
	if gotMethod != "POST" {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	rawLabels, ok := gotBody["labels"].([]any)
	if !ok {
		t.Fatalf("body labels not an array")
	}
	if len(rawLabels) != 2 {
		t.Errorf("body labels len = %d, want 2", len(rawLabels))
	}
	if len(labels) != 2 {
		t.Errorf("len = %d, want 2", len(labels))
	}
	if labels[0] != "feature" {
		t.Errorf("labels[0] = %q, want feature", labels[0])
	}
}

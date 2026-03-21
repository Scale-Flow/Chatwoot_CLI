package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListContacts(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/contacts" {
			t.Errorf("path = %q, want /api/v1/accounts/1/contacts", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("page = %q, want 1", r.URL.Query().Get("page"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"payload": []map[string]any{
					{"id": 1, "name": "Alice", "email": "alice@test.com"},
					{"id": 2, "name": "Bob", "email": "bob@test.com"},
				},
				"meta": map[string]any{
					"count":        50,
					"current_page": 1,
				},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	contacts, pg, err := client.ListContacts(context.Background(), ListContactsOpts{Page: 1})
	if err != nil {
		t.Fatalf("ListContacts error: %v", err)
	}
	if len(contacts) != 2 {
		t.Errorf("len = %d, want 2", len(contacts))
	}
	if contacts[0].Name != "Alice" {
		t.Errorf("contacts[0].Name = %q, want Alice", contacts[0].Name)
	}
	if pg == nil {
		t.Fatal("pagination is nil")
	}
	if pg.TotalCount != 50 {
		t.Errorf("TotalCount = %d, want 50", pg.TotalCount)
	}
	if pg.Page != 1 {
		t.Errorf("Page = %d, want 1", pg.Page)
	}
}

func TestGetContact(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/contacts/42" {
			t.Errorf("path = %q, want /api/v1/accounts/1/contacts/42", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":    42,
			"name":  "Alice",
			"email": "alice@test.com",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	contact, err := client.GetContact(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetContact error: %v", err)
	}
	if contact.ID != 42 {
		t.Errorf("ID = %d, want 42", contact.ID)
	}
	if contact.Name != "Alice" {
		t.Errorf("Name = %q, want Alice", contact.Name)
	}
}

func TestCreateContact(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/contacts" {
			t.Errorf("path = %q, want /api/v1/accounts/1/contacts", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":    10,
			"name":  "New Contact",
			"email": "new@test.com",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	email := "new@test.com"
	contact, err := client.CreateContact(context.Background(), CreateContactOpts{
		Name:  "New Contact",
		Email: &email,
	})
	if err != nil {
		t.Fatalf("CreateContact error: %v", err)
	}
	if contact.ID != 10 {
		t.Errorf("ID = %d, want 10", contact.ID)
	}
	if gotBody["name"] != "New Contact" {
		t.Errorf("body name = %v, want New Contact", gotBody["name"])
	}
	if gotBody["email"] != "new@test.com" {
		t.Errorf("body email = %v, want new@test.com", gotBody["email"])
	}
}

func TestUpdateContact(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/contacts/5" {
			t.Errorf("path = %q, want /api/v1/accounts/1/contacts/5", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":   5,
			"name": "Updated Name",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	name := "Updated Name"
	contact, err := client.UpdateContact(context.Background(), 5, UpdateContactOpts{Name: &name})
	if err != nil {
		t.Fatalf("UpdateContact error: %v", err)
	}
	if contact.ID != 5 {
		t.Errorf("ID = %d, want 5", contact.ID)
	}
	if gotBody["name"] != "Updated Name" {
		t.Errorf("body name = %v, want Updated Name", gotBody["name"])
	}
}

func TestDeleteContact(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/contacts/7" {
			t.Errorf("path = %q, want /api/v1/accounts/1/contacts/7", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	err := client.DeleteContact(context.Background(), 7)
	if err != nil {
		t.Fatalf("DeleteContact error: %v", err)
	}
}

func TestSearchContacts(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/contacts/search" {
			t.Errorf("path = %q, want /api/v1/accounts/1/contacts/search", r.URL.Path)
		}
		if r.URL.Query().Get("q") != "alice" {
			t.Errorf("q = %q, want alice", r.URL.Query().Get("q"))
		}
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("page = %q, want 1", r.URL.Query().Get("page"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"payload": []map[string]any{
					{"id": 1, "name": "Alice"},
				},
				"meta": map[string]any{
					"count":        1,
					"current_page": 1,
				},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	contacts, pg, err := client.SearchContacts(context.Background(), "alice", 1)
	if err != nil {
		t.Fatalf("SearchContacts error: %v", err)
	}
	if len(contacts) != 1 {
		t.Errorf("len = %d, want 1", len(contacts))
	}
	if contacts[0].Name != "Alice" {
		t.Errorf("contacts[0].Name = %q, want Alice", contacts[0].Name)
	}
	if pg.TotalCount != 1 {
		t.Errorf("TotalCount = %d, want 1", pg.TotalCount)
	}
}

func TestFilterContacts(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/contacts/filter" {
			t.Errorf("path = %q, want /api/v1/accounts/1/contacts/filter", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"payload": []map[string]any{
					{"id": 3, "name": "Charlie"},
				},
				"meta": map[string]any{
					"count":        1,
					"current_page": 1,
				},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	contacts, pg, err := client.FilterContacts(context.Background(), FilterContactsOpts{
		Page:    1,
		Payload: []any{map[string]any{"attribute_key": "email", "filter_operator": "contains", "values": []string{"test"}}},
	})
	if err != nil {
		t.Fatalf("FilterContacts error: %v", err)
	}
	if len(contacts) != 1 {
		t.Errorf("len = %d, want 1", len(contacts))
	}
	if contacts[0].Name != "Charlie" {
		t.Errorf("contacts[0].Name = %q, want Charlie", contacts[0].Name)
	}
	if pg.TotalCount != 1 {
		t.Errorf("TotalCount = %d, want 1", pg.TotalCount)
	}
	// Verify the payload was sent
	if gotBody["payload"] == nil {
		t.Error("body payload is nil")
	}
}

func TestMergeContacts(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/actions/contact_merge" {
			t.Errorf("path = %q, want /api/v1/accounts/1/actions/contact_merge", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":   1,
			"name": "Merged Contact",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	contact, err := client.MergeContacts(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("MergeContacts error: %v", err)
	}
	if contact.ID != 1 {
		t.Errorf("ID = %d, want 1", contact.ID)
	}
	// JSON numbers decode as float64
	if gotBody["base_contact_id"] != float64(1) {
		t.Errorf("base_contact_id = %v, want 1", gotBody["base_contact_id"])
	}
	if gotBody["mergee_contact_id"] != float64(2) {
		t.Errorf("mergee_contact_id = %v, want 2", gotBody["mergee_contact_id"])
	}
}

func TestListContactLabels(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/contacts/5/labels" {
			t.Errorf("path = %q, want /api/v1/accounts/1/contacts/5/labels", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]string{"vip", "enterprise"})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	labels, err := client.ListContactLabels(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListContactLabels error: %v", err)
	}
	if len(labels) != 2 {
		t.Errorf("len = %d, want 2", len(labels))
	}
	if labels[0] != "vip" {
		t.Errorf("labels[0] = %q, want vip", labels[0])
	}
}

func TestSetContactLabels(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/contacts/5/labels" {
			t.Errorf("path = %q, want /api/v1/accounts/1/contacts/5/labels", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode([]string{"vip", "premium"})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	labels, err := client.SetContactLabels(context.Background(), 5, []string{"vip", "premium"})
	if err != nil {
		t.Fatalf("SetContactLabels error: %v", err)
	}
	if len(labels) != 2 {
		t.Errorf("len = %d, want 2", len(labels))
	}
	if labels[1] != "premium" {
		t.Errorf("labels[1] = %q, want premium", labels[1])
	}
	// Verify the body
	bodyLabels, ok := gotBody["labels"].([]any)
	if !ok {
		t.Fatal("body labels not an array")
	}
	if len(bodyLabels) != 2 {
		t.Errorf("body labels len = %d, want 2", len(bodyLabels))
	}
}

func TestListContactConversations(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/contacts/5/conversations" {
			t.Errorf("path = %q, want /api/v1/accounts/1/contacts/5/conversations", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"payload": []map[string]any{
				{"id": 10, "status": "open"},
				{"id": 20, "status": "resolved"},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	convos, err := client.ListContactConversations(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListContactConversations error: %v", err)
	}
	if len(convos) != 2 {
		t.Errorf("len = %d, want 2", len(convos))
	}
	if convos[0].ID != 10 {
		t.Errorf("convos[0].ID = %d, want 10", convos[0].ID)
	}
	if convos[1].Status != "resolved" {
		t.Errorf("convos[1].Status = %q, want resolved", convos[1].Status)
	}
}

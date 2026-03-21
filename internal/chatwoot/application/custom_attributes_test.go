package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListCustomAttributes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/custom_attribute_definitions" {
			t.Errorf("path = %q, want /api/v1/accounts/1/custom_attribute_definitions", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "attribute_display_name": "Priority", "attribute_key": "priority", "attribute_model": "conversation_attribute"},
			{"id": 2, "attribute_display_name": "Region", "attribute_key": "region", "attribute_model": "contact_attribute"},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	attrs, err := client.ListCustomAttributes(context.Background())
	if err != nil {
		t.Fatalf("ListCustomAttributes error: %v", err)
	}
	if len(attrs) != 2 {
		t.Errorf("len = %d, want 2", len(attrs))
	}
	if attrs[0].AttributeDisplayName != "Priority" {
		t.Errorf("attrs[0].AttributeDisplayName = %q, want Priority", attrs[0].AttributeDisplayName)
	}
	if attrs[1].ID != 2 {
		t.Errorf("attrs[1].ID = %d, want 2", attrs[1].ID)
	}
}

func TestCreateCustomAttribute(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/custom_attribute_definitions" {
			t.Errorf("path = %q, want /api/v1/accounts/1/custom_attribute_definitions", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":                     20,
			"attribute_display_name": "Tier",
			"attribute_key":          "tier",
			"attribute_model":        "conversation_attribute",
			"attribute_display_type": "text",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	attr, err := client.CreateCustomAttribute(context.Background(), CreateCustomAttributeOpts{
		AttributeDisplayName: "Tier",
		AttributeKey:         "tier",
		AttributeModel:       "conversation_attribute",
		AttributeDisplayType: "text",
	})
	if err != nil {
		t.Fatalf("CreateCustomAttribute error: %v", err)
	}
	if attr.ID != 20 {
		t.Errorf("ID = %d, want 20", attr.ID)
	}
	if attr.AttributeDisplayName != "Tier" {
		t.Errorf("AttributeDisplayName = %q, want Tier", attr.AttributeDisplayName)
	}
	if gotBody["attribute_display_name"] != "Tier" {
		t.Errorf("body attribute_display_name = %v, want Tier", gotBody["attribute_display_name"])
	}
	if gotBody["attribute_key"] != "tier" {
		t.Errorf("body attribute_key = %v, want tier", gotBody["attribute_key"])
	}
}

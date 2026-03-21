package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListAuditLogs(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/audit_logs" {
			t.Errorf("path = %q, want /api/v1/accounts/1/audit_logs", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("page = %q, want 2", r.URL.Query().Get("page"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"payload": []map[string]any{
				{"id": 101, "action": "update", "auditable_type": "Conversation", "auditable_id": 55},
				{"id": 102, "action": "create", "auditable_type": "Contact", "auditable_id": 12},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	logs, err := client.ListAuditLogs(context.Background(), 2)
	if err != nil {
		t.Fatalf("ListAuditLogs error: %v", err)
	}
	if len(logs) != 2 {
		t.Errorf("len = %d, want 2", len(logs))
	}
	if logs[0].ID != 101 {
		t.Errorf("logs[0].ID = %d, want 101", logs[0].ID)
	}
	if logs[0].Action != "update" {
		t.Errorf("logs[0].Action = %q, want update", logs[0].Action)
	}
	if logs[1].AuditableType != "Contact" {
		t.Errorf("logs[1].AuditableType = %q, want Contact", logs[1].AuditableType)
	}
}

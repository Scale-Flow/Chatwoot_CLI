package auditlogs

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuditLogsList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"payload": []map[string]any{
				{"id": 1, "action": "create", "auditable_type": "Contact"},
			},
		})
	}))
	defer srv.Close()

	t.Setenv("CHATWOOT_BASE_URL", srv.URL)
	t.Setenv("CHATWOOT_ACCESS_TOKEN", "sk-test")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"application", "audit-logs", "list"})
	err := Cmd.Root().Execute()
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	var resp map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		t.Fatalf("JSON decode: %v\nOutput: %s", err, stdout.String())
	}
	if resp["ok"] != true {
		t.Errorf("ok = %v, want true", resp["ok"])
	}
	data, ok := resp["data"].([]any)
	if !ok || len(data) != 1 {
		t.Errorf("data length = %v, want 1", len(data))
	}
}

package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestGetReports(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v2/accounts/1/reports" {
			t.Errorf("path = %q, want /api/v2/accounts/1/reports", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("metric") != "account" {
			t.Errorf("metric = %q, want account", q.Get("metric"))
		}
		if q.Get("type") != "account" {
			t.Errorf("type = %q, want account", q.Get("type"))
		}
		if q.Get("since") != "1700000000" {
			t.Errorf("since = %q, want 1700000000", q.Get("since"))
		}
		if q.Get("until") != "1700086400" {
			t.Errorf("until = %q, want 1700086400", q.Get("until"))
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"timestamp": 1700000000, "value": 42},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	result, err := client.GetReports(context.Background(), ReportOpts{
		Metric: "account",
		Type:   "account",
		Since:  "1700000000",
		Until:  "1700086400",
	})
	if err != nil {
		t.Fatalf("GetReports error: %v", err)
	}
	if result == nil {
		t.Error("result is nil, want non-nil")
	}
}

func TestGetReportSummary(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v2/accounts/1/reports/summary" {
			t.Errorf("path = %q, want /api/v2/accounts/1/reports/summary", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"avg_first_response_time":  159.0,
			"avg_resolution_time":      10572.05,
			"conversations_count":      25,
			"incoming_messages_count":  100,
			"outgoing_messages_count":  90,
			"resolutions_count":        20,
			"reply_time":               26.0,
			"previous": map[string]any{
				"avg_first_response_time":  0.0,
				"avg_resolution_time":      0.0,
				"conversations_count":      0,
				"incoming_messages_count":  0,
				"outgoing_messages_count":  0,
				"resolutions_count":        0,
				"reply_time":               0.0,
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	summary, err := client.GetReportSummary(context.Background(), ReportOpts{
		Type:  "account",
		Since: "1700000000",
		Until: "1700086400",
	})
	if err != nil {
		t.Fatalf("GetReportSummary error: %v", err)
	}
	if summary == nil {
		t.Fatal("summary is nil")
	}
	if summary.AvgFirstResponseTime != 159.0 {
		t.Errorf("AvgFirstResponseTime = %f, want 159.0", summary.AvgFirstResponseTime)
	}
	if summary.AvgResolutionTime != 10572.05 {
		t.Errorf("AvgResolutionTime = %f, want 10572.05", summary.AvgResolutionTime)
	}
	if summary.ConversationsCount != 25 {
		t.Errorf("ConversationsCount = %d, want 25", summary.ConversationsCount)
	}
	if summary.IncomingMessagesCount != 100 {
		t.Errorf("IncomingMessagesCount = %d, want 100", summary.IncomingMessagesCount)
	}
	if summary.OutgoingMessagesCount != 90 {
		t.Errorf("OutgoingMessagesCount = %d, want 90", summary.OutgoingMessagesCount)
	}
	if summary.ResolutionsCount != 20 {
		t.Errorf("ResolutionsCount = %d, want 20", summary.ResolutionsCount)
	}
	if summary.ReplyTime != 26.0 {
		t.Errorf("ReplyTime = %f, want 26.0", summary.ReplyTime)
	}
	if summary.Previous == nil {
		t.Fatal("Previous is nil, want non-nil")
	}
	if summary.Previous.ConversationsCount != 0 {
		t.Errorf("Previous.ConversationsCount = %d, want 0", summary.Previous.ConversationsCount)
	}
}

func TestGetSummaryByAgent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/api/v2/accounts/1/summary_reports/agent") {
			t.Errorf("path = %q, want prefix /api/v2/accounts/1/summary_reports/agent", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{
				"id":                  1,
				"name":                "Alice",
				"conversations_count": 10,
			},
			{
				"id":                  2,
				"name":                "Bob",
				"conversations_count": 5,
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	result, err := client.GetSummaryByAgent(context.Background(), ReportOpts{
		Since: "1700000000",
		Until: "1700086400",
	})
	if err != nil {
		t.Fatalf("GetSummaryByAgent error: %v", err)
	}
	if result == nil {
		t.Error("result is nil, want non-nil")
	}
	items, ok := result.([]any)
	if !ok {
		t.Fatalf("result type = %T, want []any", result)
	}
	if len(items) != 2 {
		t.Errorf("len = %d, want 2", len(items))
	}
}

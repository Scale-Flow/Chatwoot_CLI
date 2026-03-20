package chatwoot

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientGet(t *testing.T) {
	var gotPath, gotHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotHeader = r.Header.Get("api_access_token")
		w.WriteHeader(200)
		w.Write([]byte(`{"id":1}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "sk-test-token", "api_access_token")
	resp, err := c.Do(context.Background(), http.MethodGet, "/api/v1/accounts/1/conversations", nil)
	if err != nil {
		t.Fatalf("Do error: %v", err)
	}
	defer resp.Body.Close()

	if gotPath != "/api/v1/accounts/1/conversations" {
		t.Errorf("path = %q, want %q", gotPath, "/api/v1/accounts/1/conversations")
	}
	if gotHeader != "sk-test-token" {
		t.Errorf("auth header = %q, want %q", gotHeader, "sk-test-token")
	}
	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestClientPost(t *testing.T) {
	var gotMethod, gotContentType string
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotContentType = r.Header.Get("Content-Type")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
		w.Write([]byte(`{"id":1}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "sk-test", "api_access_token")
	body := []byte(`{"content":"hello"}`)
	resp, err := c.Do(context.Background(), http.MethodPost, "/api/v1/test", body)
	if err != nil {
		t.Fatalf("Do error: %v", err)
	}
	defer resp.Body.Close()

	if gotMethod != "POST" {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotContentType != "application/json" {
		t.Errorf("content-type = %q, want application/json", gotContentType)
	}
	if string(gotBody) != `{"content":"hello"}` {
		t.Errorf("body = %q, want %q", string(gotBody), `{"content":"hello"}`)
	}
}

func TestMapHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		wantCode string
	}{
		{401, "unauthorized"},
		{403, "forbidden"},
		{404, "not_found"},
		{422, "validation_error"},
		{429, "rate_limited"},
		{500, "server_error"},
		{502, "server_error"},
		{503, "server_error"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("status_%d", tt.status), func(t *testing.T) {
			code := MapHTTPStatus(tt.status)
			if code != tt.wantCode {
				t.Errorf("MapHTTPStatus(%d) = %q, want %q", tt.status, code, tt.wantCode)
			}
		})
	}
}

func TestAPIErrorFormat(t *testing.T) {
	err := &APIError{StatusCode: 404, Code: "not_found", Message: "Conversation not found"}
	got := err.Error()
	if got != "chatwoot API error 404: not_found \u2014 Conversation not found" {
		t.Errorf("Error() = %q", got)
	}
}

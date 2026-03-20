package chatwoot

import (
	"context"
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

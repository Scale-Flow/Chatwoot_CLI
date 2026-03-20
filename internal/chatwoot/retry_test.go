// internal/chatwoot/retry_test.go
package chatwoot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestRetryOn429(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n < 3 {
			w.WriteHeader(429)
			w.Write([]byte(`{"error":"rate limited"}`))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"id":1}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "sk-test", "api_access_token")
	c.retryMax = 3
	c.retryBaseDelay = 0 // no delay in tests

	resp, err := c.DoWithRetry(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("DoWithRetry error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("attempts = %d, want 3", atomic.LoadInt32(&attempts))
	}
}

func TestRetryExhausted(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
		w.Write([]byte(`{"error":"rate limited"}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "sk-test", "api_access_token")
	c.retryMax = 3
	c.retryBaseDelay = 0

	resp, err := c.DoWithRetry(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("DoWithRetry error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 429 {
		t.Errorf("status = %d, want 429 (exhausted)", resp.StatusCode)
	}
}

func TestRetryOn5xx(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n < 2 {
			w.WriteHeader(503)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "sk-test", "api_access_token")
	c.retryMax = 3
	c.retryBaseDelay = 0

	resp, err := c.DoWithRetry(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("DoWithRetry error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestNoRetryOn4xx(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(404)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "sk-test", "api_access_token")
	c.retryMax = 3
	c.retryBaseDelay = 0

	resp, err := c.DoWithRetry(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("DoWithRetry error: %v", err)
	}
	defer resp.Body.Close()

	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("attempts = %d, want 1 (no retry on 404)", atomic.LoadInt32(&attempts))
	}
}

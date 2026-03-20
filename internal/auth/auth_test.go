package auth

import (
	"testing"

	"github.com/chatwoot/chatwoot-cli/internal/credentials"
)

type mockStore struct {
	tokens map[string]string
}

func (m *mockStore) Get(profile string, mode credentials.AuthMode) (string, error) {
	key := profile + "/" + string(mode)
	if t, ok := m.tokens[key]; ok {
		return t, nil
	}
	return "", credentials.ErrNotFound
}
func (m *mockStore) Set(_ string, _ credentials.AuthMode, _ string) error { return nil }
func (m *mockStore) Delete(_ string, _ credentials.AuthMode) error       { return nil }

func TestResolveApplicationAuth(t *testing.T) {
	store := &mockStore{tokens: map[string]string{
		"work/application": "sk-test",
	}}
	resolver := credentials.NewResolver(store, nil, nil)

	ctx, err := ResolveApplication(resolver, "work")
	if err != nil {
		t.Fatalf("ResolveApplication error: %v", err)
	}
	if ctx.Token != "sk-test" {
		t.Errorf("Token = %q, want %q", ctx.Token, "sk-test")
	}
	if ctx.HeaderName != "api_access_token" {
		t.Errorf("HeaderName = %q, want %q", ctx.HeaderName, "api_access_token")
	}
}

func TestResolvePlatformAuth(t *testing.T) {
	store := &mockStore{tokens: map[string]string{
		"work/platform": "pk-test",
	}}
	resolver := credentials.NewResolver(store, nil, nil)

	ctx, err := ResolvePlatform(resolver, "work")
	if err != nil {
		t.Fatalf("ResolvePlatform error: %v", err)
	}
	if ctx.Token != "pk-test" {
		t.Errorf("Token = %q, want %q", ctx.Token, "pk-test")
	}
}

func TestResolveApplicationAuthMissing(t *testing.T) {
	store := &mockStore{tokens: map[string]string{}}
	resolver := credentials.NewResolver(store, nil, nil)

	_, err := ResolveApplication(resolver, "work")
	if err == nil {
		t.Fatal("expected error for missing credentials")
	}
}

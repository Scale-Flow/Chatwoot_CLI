package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveContextFromConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	err := os.WriteFile(cfgPath, []byte(`
default_profile: work
profiles:
  work:
    base_url: https://app.chatwoot.com
    account_id: 1
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	ctx, err := resolveContextFromPath(cfgPath, "", "", 0)
	if err != nil {
		t.Fatalf("resolveContextFromPath error: %v", err)
	}
	if ctx.ProfileName != "work" {
		t.Errorf("ProfileName = %q, want %q", ctx.ProfileName, "work")
	}
	if ctx.BaseURL != "https://app.chatwoot.com" {
		t.Errorf("BaseURL = %q, want %q", ctx.BaseURL, "https://app.chatwoot.com")
	}
	if ctx.AccountID != 1 {
		t.Errorf("AccountID = %d, want 1", ctx.AccountID)
	}
}

func TestResolveContextMissingConfigWithFlags(t *testing.T) {
	ctx, err := resolveContextFromPath("/nonexistent/config.yaml", "", "https://flag.com", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.BaseURL != "https://flag.com" {
		t.Errorf("BaseURL = %q, want %q", ctx.BaseURL, "https://flag.com")
	}
	if ctx.ProfileName != "default" {
		t.Errorf("ProfileName = %q, want %q", ctx.ProfileName, "default")
	}
}

func TestResolveContextMissingConfigWithEnv(t *testing.T) {
	t.Setenv("CHATWOOT_BASE_URL", "https://env.com")
	t.Setenv("CHATWOOT_ACCOUNT_ID", "42")
	ctx, err := resolveContextFromPath("/nonexistent/config.yaml", "", "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.BaseURL != "https://env.com" {
		t.Errorf("BaseURL = %q, want %q", ctx.BaseURL, "https://env.com")
	}
	if ctx.AccountID != 42 {
		t.Errorf("AccountID = %d, want 42", ctx.AccountID)
	}
}

func TestResolveContextNoBaseURL(t *testing.T) {
	t.Setenv("CHATWOOT_BASE_URL", "")
	_, err := resolveContextFromPath("/nonexistent/config.yaml", "", "", 0)
	if err == nil {
		t.Fatal("expected error for missing base URL")
	}
}

func TestResolveContextFlagOverridesProfile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgPath, []byte(`
default_profile: work
profiles:
  work:
    base_url: https://app.chatwoot.com
    account_id: 1
`), 0644); err != nil {
		t.Fatal(err)
	}

	ctx, err := resolveContextFromPath(cfgPath, "", "https://override.com", 99)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if ctx.BaseURL != "https://override.com" {
		t.Errorf("BaseURL = %q, want %q", ctx.BaseURL, "https://override.com")
	}
	if ctx.AccountID != 99 {
		t.Errorf("AccountID = %d, want 99", ctx.AccountID)
	}
}

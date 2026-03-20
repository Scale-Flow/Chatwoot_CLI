package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	err := os.WriteFile(cfgPath, []byte(`
default_profile: work
profiles:
  work:
    base_url: https://app.chatwoot.com
    account_id: 1
  selfhosted:
    base_url: https://chatwoot.internal.corp
    account_id: 42
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFrom(cfgPath)
	if err != nil {
		t.Fatalf("LoadFrom error: %v", err)
	}

	if cfg.DefaultProfile != "work" {
		t.Errorf("DefaultProfile = %q, want %q", cfg.DefaultProfile, "work")
	}
	if len(cfg.Profiles) != 2 {
		t.Fatalf("len(Profiles) = %d, want 2", len(cfg.Profiles))
	}
	work := cfg.Profiles["work"]
	if work.BaseURL != "https://app.chatwoot.com" {
		t.Errorf("work.BaseURL = %q, want %q", work.BaseURL, "https://app.chatwoot.com")
	}
	if work.AccountID != 1 {
		t.Errorf("work.AccountID = %d, want 1", work.AccountID)
	}
}

func TestLoadFromMissingFileReturnsError(t *testing.T) {
	_, err := LoadFrom("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing config file")
	}
}

func TestResolveProfile(t *testing.T) {
	cfg := &Config{
		DefaultProfile: "work",
		Profiles: map[string]Profile{
			"work":       {BaseURL: "https://app.chatwoot.com", AccountID: 1},
			"selfhosted": {BaseURL: "https://internal.corp", AccountID: 42},
		},
	}

	tests := []struct {
		name        string
		flagProfile string
		envProfile  string
		wantName    string
		wantURL     string
	}{
		{"flag wins", "selfhosted", "work", "selfhosted", "https://internal.corp"},
		{"env wins over default", "", "selfhosted", "selfhosted", "https://internal.corp"},
		{"default used", "", "", "work", "https://app.chatwoot.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envProfile != "" {
				t.Setenv("CHATWOOT_PROFILE", tt.envProfile)
			}
			name, profile, err := cfg.ResolveProfile(tt.flagProfile)
			if err != nil {
				t.Fatalf("ResolveProfile error: %v", err)
			}
			if name != tt.wantName {
				t.Errorf("name = %q, want %q", name, tt.wantName)
			}
			if profile.BaseURL != tt.wantURL {
				t.Errorf("BaseURL = %q, want %q", profile.BaseURL, tt.wantURL)
			}
		})
	}
}

func TestResolveProfileFallsBackToDefault(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"default": {BaseURL: "https://fallback.com", AccountID: 99},
		},
	}
	name, profile, err := cfg.ResolveProfile("")
	if err != nil {
		t.Fatalf("ResolveProfile error: %v", err)
	}
	if name != "default" {
		t.Errorf("name = %q, want %q", name, "default")
	}
	if profile.BaseURL != "https://fallback.com" {
		t.Errorf("BaseURL = %q, want %q", profile.BaseURL, "https://fallback.com")
	}
}

func TestResolveProfileNotFound(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{},
	}
	_, _, err := cfg.ResolveProfile("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent profile")
	}
}

func TestResolveOverrides(t *testing.T) {
	profile := Profile{BaseURL: "https://app.chatwoot.com", AccountID: 1}

	tests := []struct {
		name          string
		flagBaseURL   string
		envBaseURL    string
		flagAccountID int
		envAccountID  string
		wantURL       string
		wantAccount   int
	}{
		{"no overrides", "", "", 0, "", "https://app.chatwoot.com", 1},
		{"flag base URL", "https://custom.com", "", 0, "", "https://custom.com", 1},
		{"env base URL", "", "https://env.com", 0, "", "https://env.com", 1},
		{"flag account ID", "", "", 99, "", "https://app.chatwoot.com", 99},
		{"env account ID", "", "", 0, "55", "https://app.chatwoot.com", 55},
		{"flag beats env", "https://flag.com", "https://env.com", 77, "88", "https://flag.com", 77},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envBaseURL != "" {
				t.Setenv("CHATWOOT_BASE_URL", tt.envBaseURL)
			}
			if tt.envAccountID != "" {
				t.Setenv("CHATWOOT_ACCOUNT_ID", tt.envAccountID)
			}
			resolved := ResolveOverrides(profile, tt.flagBaseURL, tt.flagAccountID)
			if resolved.BaseURL != tt.wantURL {
				t.Errorf("BaseURL = %q, want %q", resolved.BaseURL, tt.wantURL)
			}
			if resolved.AccountID != tt.wantAccount {
				t.Errorf("AccountID = %d, want %d", resolved.AccountID, tt.wantAccount)
			}
		})
	}
}

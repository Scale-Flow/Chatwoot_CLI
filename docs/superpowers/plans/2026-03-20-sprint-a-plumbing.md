# Sprint A: Plumbing (Layers 2 + 3 + 4) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the full internal stack from config/credentials through HTTP transport to typed API family clients, with no Cobra dependency in any layer.

**Architecture:** Three layers stacked in dependency order: config/credentials (resolve runtime context) → transport (execute HTTP requests with retry, auth, error mapping) → API clients (typed methods per family). Each layer depends only on layers below it and the `contract` package from Layer 1.

**Tech Stack:** Go 1.26.1, `spf13/viper` (config), `zalando/go-keyring` (keychain), `net/http` + `net/http/httptest` (transport/testing), `log/slog` (diagnostics), `encoding/json` (serialization)

**Spec:** `docs/superpowers/specs/2026-03-20-design-decisions-and-sprint-structure.md`

---

## File Structure

### Layer 2: Config & Credentials

| File | Responsibility |
|------|---------------|
| `internal/config/config.go` | Viper config loading, profile struct, profile resolution with precedence chain |
| `internal/config/config_test.go` | Config loading, profile resolution, precedence tests |
| `internal/credentials/store.go` | `Store` interface definition |
| `internal/credentials/keychain.go` | Keychain backend via `go-keyring` |
| `internal/credentials/file.go` | File backend (`0600` permissions) |
| `internal/credentials/env.go` | Environment variable backend |
| `internal/credentials/resolver.go` | Multi-backend resolver (env → keychain → file → error) |
| `internal/credentials/credentials_test.go` | Tests for all backends and resolver |
| `internal/auth/auth.go` | Auth mode enum, credential resolution per API family |
| `internal/auth/auth_test.go` | Auth mode tests |

### Layer 3: Transport

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/client.go` | `Client` struct, `Doer` interface, request builder, auth header injection, error mapping |
| `internal/chatwoot/retry.go` | Retry with exponential backoff + jitter |
| `internal/chatwoot/errors.go` | HTTP status → contract error code mapping |
| `internal/chatwoot/paginate.go` | `ListAll` auto-pagination helper |
| `internal/chatwoot/client_test.go` | Transport tests (auth, error mapping, request building) |
| `internal/chatwoot/retry_test.go` | Retry logic tests |
| `internal/chatwoot/paginate_test.go` | Pagination helper tests |

### Layer 4: API Family Clients

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/application/client.go` | Application API client struct and constructor |
| `internal/chatwoot/application/profile.go` | Profile get/update methods |
| `internal/chatwoot/application/conversations.go` | Conversations list/get methods (proves the pattern) |
| `internal/chatwoot/application/models.go` | Application-specific model types |
| `internal/chatwoot/application/application_test.go` | Application client tests with httptest |
| `internal/chatwoot/platform/client.go` | Platform API client struct and constructor |
| `internal/chatwoot/platform/accounts.go` | Accounts get/create methods (proves the pattern) |
| `internal/chatwoot/platform/models.go` | Platform-specific model types |
| `internal/chatwoot/platform/platform_test.go` | Platform client tests with httptest |
| `internal/chatwoot/clientapi/client.go` | Client API client struct and constructor |
| `internal/chatwoot/clientapi/contacts.go` | Contacts create/get methods (proves the pattern) |
| `internal/chatwoot/clientapi/models.go` | Client API model types |
| `internal/chatwoot/clientapi/clientapi_test.go` | Client API tests with httptest |

### Shared / Support

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/models.go` | Shared types: timestamp helpers, common ID types |
| `internal/testutil/testutil.go` | httptest helpers, mock credential store |

---

## Task 1: Config Package — Profile Struct and Loading

**Files:**
- Modify: `internal/config/config.go` (currently stub)
- Create: `internal/config/config_test.go`

- [ ] **Step 1: Write failing test for config loading from YAML**

```go
// internal/config/config_test.go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/config/ -v -run TestLoadFromFile`
Expected: FAIL — `LoadFrom` not defined

- [ ] **Step 3: Write minimal implementation**

```go
// internal/config/config.go
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Profile holds the configuration for a single named profile.
type Profile struct {
	BaseURL   string `mapstructure:"base_url"`
	AccountID int    `mapstructure:"account_id"`
}

// Config holds the full CLI configuration.
type Config struct {
	DefaultProfile string             `mapstructure:"default_profile"`
	Profiles       map[string]Profile `mapstructure:"profiles"`
}

// LoadFrom reads config from a specific file path.
func LoadFrom(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/config/ -v -run TestLoadFromFile`
Expected: PASS

- [ ] **Step 5: Write failing test for profile resolution precedence**

```go
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
```

- [ ] **Step 6: Run tests to verify they fail**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/config/ -v -run TestResolveProfile`
Expected: FAIL — `ResolveProfile` not defined

- [ ] **Step 7: Implement ResolveProfile**

Add to `internal/config/config.go`:

```go
// ResolveProfile selects the active profile using precedence:
// flag > CHATWOOT_PROFILE env > default_profile > "default".
func (c *Config) ResolveProfile(flagProfile string) (string, Profile, error) {
	name := flagProfile
	if name == "" {
		name = os.Getenv("CHATWOOT_PROFILE")
	}
	if name == "" {
		name = c.DefaultProfile
	}
	if name == "" {
		name = "default"
	}
	profile, ok := c.Profiles[name]
	if !ok {
		return "", Profile{}, fmt.Errorf("profile %q not found in config", name)
	}
	return name, profile, nil
}
```

Add `"os"` to imports.

- [ ] **Step 8: Run all config tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/config/ -v`
Expected: PASS

- [ ] **Step 9: Write failing test for missing config file returns empty config**

```go
func TestLoadFromMissingFileReturnsError(t *testing.T) {
	_, err := LoadFrom("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing config file")
	}
}
```

- [ ] **Step 10: Run test to verify it passes (already handles this case)**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/config/ -v -run TestLoadFromMissing`
Expected: PASS (Viper returns error for missing file)

- [ ] **Step 11: Write failing test for override resolution (base URL and account ID)**

```go
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
```

- [ ] **Step 12: Implement ResolveOverrides**

Add to `internal/config/config.go`:

```go
// ResolveOverrides applies flag and env var overrides to a profile.
// Precedence: flag > env var > profile value.
func ResolveOverrides(p Profile, flagBaseURL string, flagAccountID int) Profile {
	if flagBaseURL != "" {
		p.BaseURL = flagBaseURL
	} else if env := os.Getenv("CHATWOOT_BASE_URL"); env != "" {
		p.BaseURL = env
	}
	if flagAccountID != 0 {
		p.AccountID = flagAccountID
	} else if env := os.Getenv("CHATWOOT_ACCOUNT_ID"); env != "" {
		if id, err := strconv.Atoi(env); err == nil {
			p.AccountID = id
		}
	}
	return p
}
```

Add `"strconv"` to imports.

- [ ] **Step 13: Run all config tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/config/ -v`
Expected: PASS

- [ ] **Step 14: Commit**

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "feat(config): add profile loading and resolution with precedence chain"
```

---

## Task 2: Credential Store Interface and Env Backend

**Files:**
- Create: `internal/credentials/store.go`
- Create: `internal/credentials/env.go`
- Create: `internal/credentials/credentials_test.go`
- Delete content of: `internal/credentials/credentials.go` (replaced by store.go)

- [ ] **Step 1: Write the Store interface and AuthMode type**

```go
// internal/credentials/store.go
package credentials

import "errors"

// ErrNotFound is returned when no credential is found.
var ErrNotFound = errors.New("credential not found")

// AuthMode identifies which API family's credential to resolve.
type AuthMode string

const (
	ModeApplication AuthMode = "application"
	ModePlatform    AuthMode = "platform"
)

// Store is the interface for credential backends.
type Store interface {
	Get(profile string, mode AuthMode) (string, error)
	Set(profile string, mode AuthMode, token string) error
	Delete(profile string, mode AuthMode) error
}
```

- [ ] **Step 2: Write failing test for env backend**

```go
// internal/credentials/credentials_test.go
package credentials

import "testing"

func TestEnvStoreGetApplication(t *testing.T) {
	t.Setenv("CHATWOOT_ACCESS_TOKEN", "sk-test-token")
	store := &EnvStore{}

	token, err := store.Get("any-profile", ModeApplication)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if token != "sk-test-token" {
		t.Errorf("token = %q, want %q", token, "sk-test-token")
	}
}

func TestEnvStoreGetPlatform(t *testing.T) {
	t.Setenv("CHATWOOT_PLATFORM_TOKEN", "pk-test-token")
	store := &EnvStore{}

	token, err := store.Get("any-profile", ModePlatform)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if token != "pk-test-token" {
		t.Errorf("token = %q, want %q", token, "pk-test-token")
	}
}

func TestEnvStoreGetNotSet(t *testing.T) {
	store := &EnvStore{}
	_, err := store.Get("any", ModeApplication)
	if err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/credentials/ -v -run TestEnvStore`
Expected: FAIL — `EnvStore` not defined

- [ ] **Step 4: Implement EnvStore**

```go
// internal/credentials/env.go
package credentials

import "os"

// EnvStore resolves credentials from environment variables.
// It ignores the profile parameter — env vars are global.
type EnvStore struct{}

func (s *EnvStore) Get(_ string, mode AuthMode) (string, error) {
	var envVar string
	switch mode {
	case ModeApplication:
		envVar = "CHATWOOT_ACCESS_TOKEN"
	case ModePlatform:
		envVar = "CHATWOOT_PLATFORM_TOKEN"
	default:
		return "", ErrNotFound
	}
	token := os.Getenv(envVar)
	if token == "" {
		return "", ErrNotFound
	}
	return token, nil
}

func (s *EnvStore) Set(_ string, _ AuthMode, _ string) error {
	return errors.New("cannot set credentials via environment store")
}

func (s *EnvStore) Delete(_ string, _ AuthMode) error {
	return errors.New("cannot delete credentials via environment store")
}
```

Add `"errors"` to imports.

- [ ] **Step 5: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/credentials/ -v`
Expected: PASS

- [ ] **Step 6: Remove old stub file**

Delete `internal/credentials/credentials.go` (the stub with just `package credentials`). The package declaration is now in `store.go`.

- [ ] **Step 7: Verify build**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./...`
Expected: Success

- [ ] **Step 8: Commit**

```bash
git add internal/credentials/
git commit -m "feat(credentials): add Store interface and env var backend"
```

---

## Task 3: Credential File Backend (0600 Permissions)

**Files:**
- Create: `internal/credentials/file.go`
- Modify: `internal/credentials/credentials_test.go`

- [ ] **Step 1: Write failing test for file store round-trip**

Add to `credentials_test.go`:

```go
func TestFileStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.yaml")
	store := NewFileStore(path)

	// Set a credential
	err := store.Set("work", ModeApplication, "sk-secret")
	if err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// Verify file permissions
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("file permissions = %o, want 600", info.Mode().Perm())
	}

	// Get the credential back
	token, err := store.Get("work", ModeApplication)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if token != "sk-secret" {
		t.Errorf("token = %q, want %q", token, "sk-secret")
	}
}

func TestFileStoreGetNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.yaml")
	store := NewFileStore(path)

	_, err := store.Get("work", ModeApplication)
	if err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestFileStoreDelete(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.yaml")
	store := NewFileStore(path)

	_ = store.Set("work", ModeApplication, "sk-secret")
	err := store.Delete("work", ModeApplication)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}

	_, err = store.Get("work", ModeApplication)
	if err != ErrNotFound {
		t.Errorf("after delete: err = %v, want ErrNotFound", err)
	}
}

func TestFileStoreRejectsWidePerm(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.yaml")

	// Write file with wide permissions
	err := os.WriteFile(path, []byte("profiles: {}"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	store := NewFileStore(path)
	_, err = store.Get("work", ModeApplication)
	if err == nil {
		t.Fatal("expected error for wide permissions")
	}
}
```

Add `"os"`, `"path/filepath"` to test imports.

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/credentials/ -v -run TestFileStore`
Expected: FAIL — `NewFileStore` not defined

- [ ] **Step 3: Implement FileStore**

```go
// internal/credentials/file.go
package credentials

import (
	"fmt"
	"os"
	"sync"

	"go.yaml.in/yaml/v3"
)

// fileData is the on-disk credential file structure.
type fileData struct {
	Profiles map[string]fileProfile `yaml:"profiles"`
}

type fileProfile struct {
	ApplicationToken string `yaml:"application_token,omitempty"`
	PlatformToken    string `yaml:"platform_token,omitempty"`
}

// FileStore reads and writes credentials to a YAML file with 0600 permissions.
type FileStore struct {
	path string
	mu   sync.Mutex
}

// NewFileStore creates a file-backed credential store at the given path.
func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

func (s *FileStore) Get(profile string, mode AuthMode) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return "", err
	}

	p, ok := data.Profiles[profile]
	if !ok {
		return "", ErrNotFound
	}

	var token string
	switch mode {
	case ModeApplication:
		token = p.ApplicationToken
	case ModePlatform:
		token = p.PlatformToken
	}
	if token == "" {
		return "", ErrNotFound
	}
	return token, nil
}

func (s *FileStore) Set(profile string, mode AuthMode, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if data.Profiles == nil {
		data.Profiles = make(map[string]fileProfile)
	}

	p := data.Profiles[profile]
	switch mode {
	case ModeApplication:
		p.ApplicationToken = token
	case ModePlatform:
		p.PlatformToken = token
	}
	data.Profiles[profile] = p

	return s.save(data)
}

func (s *FileStore) Delete(profile string, mode AuthMode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return err
	}

	p, ok := data.Profiles[profile]
	if !ok {
		return nil
	}

	switch mode {
	case ModeApplication:
		p.ApplicationToken = ""
	case ModePlatform:
		p.PlatformToken = ""
	}
	data.Profiles[profile] = p

	return s.save(data)
}

func (s *FileStore) load() (fileData, error) {
	var data fileData

	info, err := os.Stat(s.path)
	if os.IsNotExist(err) {
		return data, ErrNotFound
	}
	if err != nil {
		return data, fmt.Errorf("stat credentials file: %w", err)
	}

	// Reject files with permissions wider than 0600
	if info.Mode().Perm()&0077 != 0 {
		return data, fmt.Errorf("credentials file %s has permissions %o, want 0600 or stricter", s.path, info.Mode().Perm())
	}

	raw, err := os.ReadFile(s.path)
	if err != nil {
		return data, fmt.Errorf("read credentials file: %w", err)
	}

	if err := yaml.Unmarshal(raw, &data); err != nil {
		return data, fmt.Errorf("unmarshal credentials: %w", err)
	}
	return data, nil
}

func (s *FileStore) save(data fileData) error {
	raw, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal credentials: %w", err)
	}
	return os.WriteFile(s.path, raw, 0600)
}
```

- [ ] **Step 4: Run all credential tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/credentials/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/credentials/file.go internal/credentials/credentials_test.go
git commit -m "feat(credentials): add file-based credential store with 0600 enforcement"
```

---

## Task 4: Keychain Backend

**Files:**
- Create: `internal/credentials/keychain.go`
- Modify: `internal/credentials/credentials_test.go`

- [ ] **Step 1: Write failing test for keychain store using mock**

Add to `credentials_test.go`:

```go
func TestKeychainStoreRoundTrip(t *testing.T) {
	keyring.MockInit() // go-keyring test mock

	store := NewKeychainStore()

	err := store.Set("work", ModeApplication, "sk-keychain-secret")
	if err != nil {
		t.Fatalf("Set error: %v", err)
	}

	token, err := store.Get("work", ModeApplication)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if token != "sk-keychain-secret" {
		t.Errorf("token = %q, want %q", token, "sk-keychain-secret")
	}
}

func TestKeychainStoreGetNotFound(t *testing.T) {
	keyring.MockInit()
	store := NewKeychainStore()

	_, err := store.Get("nonexistent", ModeApplication)
	if err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestKeychainStoreDelete(t *testing.T) {
	keyring.MockInit()
	store := NewKeychainStore()

	_ = store.Set("work", ModeApplication, "sk-secret")
	err := store.Delete("work", ModeApplication)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}

	_, err = store.Get("work", ModeApplication)
	if err != ErrNotFound {
		t.Errorf("after delete: err = %v, want ErrNotFound", err)
	}
}
```

Add `"github.com/zalando/go-keyring"` as `keyring` to test imports.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/credentials/ -v -run TestKeychainStore`
Expected: FAIL — `NewKeychainStore` not defined

- [ ] **Step 3: Implement KeychainStore**

```go
// internal/credentials/keychain.go
package credentials

import (
	"errors"
	"fmt"

	keyring "github.com/zalando/go-keyring"
)

const keychainService = "chatwoot-cli"

// KeychainStore stores credentials in the OS keychain via go-keyring.
type KeychainStore struct{}

// NewKeychainStore creates a keychain-backed credential store.
func NewKeychainStore() *KeychainStore {
	return &KeychainStore{}
}

func (s *KeychainStore) Get(profile string, mode AuthMode) (string, error) {
	user := fmt.Sprintf("%s/%s", profile, mode)
	token, err := keyring.Get(keychainService, user)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("keychain get: %w", err)
	}
	return token, nil
}

func (s *KeychainStore) Set(profile string, mode AuthMode, token string) error {
	user := fmt.Sprintf("%s/%s", profile, mode)
	if err := keyring.Set(keychainService, user, token); err != nil {
		return fmt.Errorf("keychain set: %w", err)
	}
	return nil
}

func (s *KeychainStore) Delete(profile string, mode AuthMode) error {
	user := fmt.Sprintf("%s/%s", profile, mode)
	if err := keyring.Delete(keychainService, user); err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("keychain delete: %w", err)
	}
	return nil
}
```

- [ ] **Step 4: Run all credential tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/credentials/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/credentials/keychain.go internal/credentials/credentials_test.go
git commit -m "feat(credentials): add keychain backend via go-keyring"
```

---

## Task 5: Credential Resolver (Multi-Backend)

**Files:**
- Create: `internal/credentials/resolver.go`
- Modify: `internal/credentials/credentials_test.go`

- [ ] **Step 1: Write failing test for resolver precedence**

Add to `credentials_test.go`:

```go
func TestResolverEnvWins(t *testing.T) {
	t.Setenv("CHATWOOT_ACCESS_TOKEN", "env-token")
	keyring.MockInit()
	ks := NewKeychainStore()
	_ = ks.Set("work", ModeApplication, "keychain-token")

	dir := t.TempDir()
	fs := NewFileStore(filepath.Join(dir, "creds.yaml"))
	_ = fs.Set("work", ModeApplication, "file-token")

	resolver := NewResolver(&EnvStore{}, ks, fs)
	token, source, err := resolver.Get("work", ModeApplication)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if token != "env-token" {
		t.Errorf("token = %q, want %q", token, "env-token")
	}
	if source != SourceEnv {
		t.Errorf("source = %q, want %q", source, SourceEnv)
	}
}

func TestResolverKeychainFallback(t *testing.T) {
	keyring.MockInit()
	ks := NewKeychainStore()
	_ = ks.Set("work", ModeApplication, "keychain-token")

	resolver := NewResolver(&EnvStore{}, ks, nil)
	token, source, err := resolver.Get("work", ModeApplication)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if token != "keychain-token" {
		t.Errorf("token = %q, want %q", token, "keychain-token")
	}
	if source != SourceKeychain {
		t.Errorf("source = %q, want %q", source, SourceKeychain)
	}
}

func TestResolverFileFallback(t *testing.T) {
	keyring.MockInit()

	dir := t.TempDir()
	fs := NewFileStore(filepath.Join(dir, "creds.yaml"))
	_ = fs.Set("work", ModeApplication, "file-token")

	resolver := NewResolver(&EnvStore{}, NewKeychainStore(), fs)
	token, source, err := resolver.Get("work", ModeApplication)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if token != "file-token" {
		t.Errorf("token = %q, want %q", token, "file-token")
	}
	if source != SourceFile {
		t.Errorf("source = %q, want %q", source, SourceFile)
	}
}

func TestResolverNoneFound(t *testing.T) {
	keyring.MockInit()
	resolver := NewResolver(&EnvStore{}, NewKeychainStore(), nil)
	_, _, err := resolver.Get("work", ModeApplication)
	if err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/credentials/ -v -run TestResolver`
Expected: FAIL — `NewResolver`, `SourceEnv`, etc. not defined

- [ ] **Step 3: Implement Resolver**

```go
// internal/credentials/resolver.go
package credentials

// Source identifies where a credential was resolved from.
type Source string

const (
	SourceEnv      Source = "environment"
	SourceKeychain Source = "keychain"
	SourceFile     Source = "file"
)

// Resolver tries multiple credential backends in priority order.
type Resolver struct {
	env      Store
	keychain Store
	file     Store
}

// NewResolver creates a multi-backend resolver.
// Any backend may be nil (skipped).
func NewResolver(env, keychain, file Store) *Resolver {
	return &Resolver{env: env, keychain: keychain, file: file}
}

// Get resolves a credential by trying: env → keychain → file.
// Returns the token, which source it came from, and any error.
func (r *Resolver) Get(profile string, mode AuthMode) (string, Source, error) {
	backends := []struct {
		store  Store
		source Source
	}{
		{r.env, SourceEnv},
		{r.keychain, SourceKeychain},
		{r.file, SourceFile},
	}
	for _, b := range backends {
		if b.store == nil {
			continue
		}
		token, err := b.store.Get(profile, mode)
		if err == ErrNotFound {
			continue
		}
		if err != nil {
			continue
		}
		return token, b.source, nil
	}
	return "", "", ErrNotFound
}
```

- [ ] **Step 4: Run all credential tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/credentials/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/credentials/resolver.go internal/credentials/credentials_test.go
git commit -m "feat(credentials): add multi-backend resolver with env > keychain > file precedence"
```

---

## Task 6: Auth Mode Modeling

**Files:**
- Modify: `internal/auth/auth.go` (currently stub)
- Create: `internal/auth/auth_test.go`

- [ ] **Step 1: Write failing test for auth context construction**

```go
// internal/auth/auth_test.go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/auth/ -v`
Expected: FAIL — `ResolveApplication` not defined

- [ ] **Step 3: Implement auth resolution**

```go
// internal/auth/auth.go
package auth

import (
	"fmt"

	"github.com/chatwoot/chatwoot-cli/internal/credentials"
)

// TokenAuth holds resolved authentication context for token-based API families.
type TokenAuth struct {
	Token      string
	HeaderName string
	Source     credentials.Source
}

// ClientAuth holds resolved authentication context for the Client API.
type ClientAuth struct {
	InboxIdentifier   string
	ContactIdentifier string
}

// ResolveApplication resolves application API credentials for a profile.
func ResolveApplication(resolver *credentials.Resolver, profile string) (TokenAuth, error) {
	token, source, err := resolver.Get(profile, credentials.ModeApplication)
	if err != nil {
		return TokenAuth{}, fmt.Errorf("no application credentials for profile %q: %w", profile, err)
	}
	return TokenAuth{Token: token, HeaderName: "api_access_token", Source: source}, nil
}

// ResolvePlatform resolves platform API credentials for a profile.
func ResolvePlatform(resolver *credentials.Resolver, profile string) (TokenAuth, error) {
	token, source, err := resolver.Get(profile, credentials.ModePlatform)
	if err != nil {
		return TokenAuth{}, fmt.Errorf("no platform credentials for profile %q: %w", profile, err)
	}
	return TokenAuth{Token: token, HeaderName: "api_access_token", Source: source}, nil
}
```

- [ ] **Step 4: Run all auth tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/auth/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/auth/auth.go internal/auth/auth_test.go
git commit -m "feat(auth): add auth mode modeling and credential resolution per API family"
```

---

## Task 7: HTTP Transport — Client Struct, Doer Interface, Request Builder

**Files:**
- Modify: `internal/chatwoot/client.go` (currently stub)
- Create: `internal/chatwoot/client_test.go`

- [ ] **Step 1: Write failing test for request building with auth header**

```go
// internal/chatwoot/client_test.go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/ -v -run TestClient`
Expected: FAIL — `NewClient` not defined

- [ ] **Step 3: Implement Client**

```go
// internal/chatwoot/client.go
package chatwoot

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

// Doer executes an HTTP request. Matches *http.Client.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client is the shared HTTP transport for all Chatwoot API calls.
type Client struct {
	baseURL    string
	token      string
	headerName string
	http       Doer
}

// NewClient creates a Client with the given base URL and auth credentials.
func NewClient(baseURL, token, headerName string) *Client {
	return &Client{
		baseURL:    baseURL,
		token:      token,
		headerName: headerName,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewClientWithDoer creates a Client with a custom Doer (for testing).
func NewClientWithDoer(baseURL, token, headerName string, doer Doer) *Client {
	return &Client{
		baseURL:    baseURL,
		token:      token,
		headerName: headerName,
		http:       doer,
	}
}

// Do executes an HTTP request against the Chatwoot API.
// Path should start with / (e.g., "/api/v1/accounts/1/conversations").
// Body may be nil for GET/DELETE requests.
func (c *Client) Do(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	url := c.baseURL + path

	var bodyReader *bytes.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	var req *http.Request
	var err error
	if bodyReader != nil {
		req, err = http.NewRequestWithContext(ctx, method, url, bodyReader)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	if c.token != "" && c.headerName != "" {
		req.Header.Set(c.headerName, c.token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.http.Do(req)
}
```

- [ ] **Step 4: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/ -v -run TestClient`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/chatwoot/client.go internal/chatwoot/client_test.go
git commit -m "feat(transport): add HTTP client with auth header injection"
```

---

## Task 8: Transport — Error Mapping

**Files:**
- Create: `internal/chatwoot/errors.go`
- Modify: `internal/chatwoot/client_test.go`

- [ ] **Step 1: Write failing test for HTTP error mapping**

Add to `client_test.go`:

```go
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
```

Add `"fmt"` to test imports.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/ -v -run TestMapHTTPError`
Expected: FAIL — `MapHTTPStatus` not defined

- [ ] **Step 3: Implement error mapping**

```go
// internal/chatwoot/errors.go
package chatwoot

import "github.com/chatwoot/chatwoot-cli/internal/contract"

// MapHTTPStatus maps an HTTP status code to a contract error code.
func MapHTTPStatus(status int) string {
	switch {
	case status == 401:
		return contract.ErrCodeUnauthorized
	case status == 403:
		return contract.ErrCodeForbidden
	case status == 404:
		return contract.ErrCodeNotFound
	case status == 422:
		return contract.ErrCodeValidation
	case status == 429:
		return contract.ErrCodeRateLimited
	case status >= 500:
		return contract.ErrCodeServer
	default:
		return contract.ErrCodeServer
	}
}
```

- [ ] **Step 4: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/ -v -run TestMapHTTPError`
Expected: PASS

- [ ] **Step 5: Write failing test for APIError type**

Add to `client_test.go`:

```go
func TestAPIErrorFormat(t *testing.T) {
	err := &APIError{StatusCode: 404, Code: "not_found", Message: "Conversation not found"}
	got := err.Error()
	if got != "chatwoot API error 404: not_found — Conversation not found" {
		t.Errorf("Error() = %q", got)
	}
}
```

- [ ] **Step 6: Implement APIError**

Add to `internal/chatwoot/errors.go`:

```go
import "fmt"

// APIError represents an error response from the Chatwoot API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Detail     any
}

func (e *APIError) Error() string {
	return fmt.Sprintf("chatwoot API error %d: %s — %s", e.StatusCode, e.Code, e.Message)
}
```

Update imports to include both `"fmt"` and the contract package.

- [ ] **Step 7: Run all transport tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/ -v`
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add internal/chatwoot/errors.go internal/chatwoot/client_test.go
git commit -m "feat(transport): add HTTP status to contract error code mapping"
```

---

## Task 9: Transport — Retry with Exponential Backoff

**Files:**
- Create: `internal/chatwoot/retry.go`
- Create: `internal/chatwoot/retry_test.go`

- [ ] **Step 1: Write failing test for retry on 429**

```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/ -v -run TestRetry`
Expected: FAIL — `DoWithRetry`, `retryMax`, `retryBaseDelay` not defined

- [ ] **Step 3: Add retry fields to Client and update NewClient**

Modify `internal/chatwoot/client.go` to add retry fields to the `Client` struct and set defaults in `NewClient`:

```go
// Add fields to Client struct:
retryMax       int
retryBaseDelay time.Duration
```

Set defaults in `NewClient`:
```go
retryMax:       3,
retryBaseDelay: time.Second,
```

- [ ] **Step 4: Implement DoWithRetry**

```go
// internal/chatwoot/retry.go
package chatwoot

import (
	"context"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

// DoWithRetry executes a request with retry on 429 and 5xx responses.
func (c *Client) DoWithRetry(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.retryMax; attempt++ {
		if attempt > 0 {
			delay := c.backoffDelay(attempt, resp)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, err = c.Do(ctx, method, path, body)
		if err != nil {
			return nil, err
		}

		if !shouldRetry(resp.StatusCode) {
			return resp, nil
		}

		// Close body before retry to prevent resource leak
		if attempt < c.retryMax {
			resp.Body.Close()
		}
	}

	return resp, nil
}

func shouldRetry(status int) bool {
	return status == 429 || status >= 500
}

func (c *Client) backoffDelay(attempt int, resp *http.Response) time.Duration {
	// Respect Retry-After header if present
	if resp != nil {
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, err := strconv.Atoi(ra); err == nil {
				return time.Duration(secs) * time.Second
			}
		}
	}

	// Exponential backoff with jitter
	base := c.retryBaseDelay
	for i := 1; i < attempt; i++ {
		base *= 2
	}
	// Add jitter: 0-25% of base
	jitter := time.Duration(rand.Int64N(int64(base / 4)))
	return base + jitter
}
```

- [ ] **Step 5: Run all retry tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/ -v -run TestRetry`
Expected: PASS

- [ ] **Step 6: Run test for no-retry case**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/ -v -run TestNoRetry`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add internal/chatwoot/retry.go internal/chatwoot/retry_test.go internal/chatwoot/client.go
git commit -m "feat(transport): add retry with exponential backoff on 429 and 5xx"
```

---

## Task 10: Transport — Response Decoding Helper

**Files:**
- Modify: `internal/chatwoot/client.go`
- Modify: `internal/chatwoot/client_test.go`

- [ ] **Step 1: Write failing test for DecodeResponse**

Add to `client_test.go`:

```go
func TestDecodeResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"id":42,"name":"Test"}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "sk-test", "api_access_token")
	resp, _ := c.Do(context.Background(), http.MethodGet, "/test", nil)

	var result struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	err := DecodeResponse(resp, &result)
	if err != nil {
		t.Fatalf("DecodeResponse error: %v", err)
	}
	if result.ID != 42 {
		t.Errorf("ID = %d, want 42", result.ID)
	}
	if result.Name != "Test" {
		t.Errorf("Name = %q, want %q", result.Name, "Test")
	}
}

func TestDecodeResponseError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"message":"not found"}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "sk-test", "api_access_token")
	resp, _ := c.Do(context.Background(), http.MethodGet, "/test", nil)

	var result struct{}
	err := DecodeResponse(resp, &result)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("err type = %T, want *APIError", err)
	}
	if apiErr.Code != "not_found" {
		t.Errorf("Code = %q, want %q", apiErr.Code, "not_found")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/ -v -run TestDecodeResponse`
Expected: FAIL — `DecodeResponse` not defined

- [ ] **Step 3: Implement DecodeResponse**

Add to `internal/chatwoot/client.go`:

```go
// DecodeResponse reads and decodes an HTTP response.
// For 2xx responses, it decodes the body into target.
// For non-2xx responses, it returns an *APIError.
func DecodeResponse(resp *http.Response, target any) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if target != nil && len(body) > 0 {
			if err := json.Unmarshal(body, target); err != nil {
				return fmt.Errorf("decode response: %w", err)
			}
		}
		return nil
	}

	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Code:       MapHTTPStatus(resp.StatusCode),
	}

	// Try to parse error message from response body
	var errBody struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	if json.Unmarshal(body, &errBody) == nil {
		if errBody.Message != "" {
			apiErr.Message = errBody.Message
		} else if errBody.Error != "" {
			apiErr.Message = errBody.Error
		}
	}
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(resp.StatusCode)
	}

	return apiErr
}
```

Add `"encoding/json"` and `"io"` to imports.

- [ ] **Step 4: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/ -v -run TestDecodeResponse`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/chatwoot/client.go internal/chatwoot/client_test.go
git commit -m "feat(transport): add response decoding with error mapping"
```

---

## Task 11: Pagination Helper

**Files:**
- Create: `internal/chatwoot/paginate.go`
- Create: `internal/chatwoot/paginate_test.go`

- [ ] **Step 1: Write failing test for ListAll**

```go
// internal/chatwoot/paginate_test.go
package chatwoot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chatwoot/chatwoot-cli/internal/contract"
)

func TestListAll(t *testing.T) {
	// Server returns 3 pages of 2 items each
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		var items []map[string]any
		switch page {
		case "", "1":
			items = []map[string]any{{"id": 1}, {"id": 2}}
		case "2":
			items = []map[string]any{{"id": 3}, {"id": 4}}
		case "3":
			items = []map[string]any{{"id": 5}, {"id": 6}}
		default:
			items = []map[string]any{}
		}

		resp := map[string]any{
			"data": items,
			"meta": map[string]any{
				"page":        page,
				"total_pages": 3,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	type Item struct {
		ID int `json:"id"`
	}

	fetcher := func(ctx context.Context, page int) ([]Item, *contract.Pagination, error) {
		c := NewClient(srv.URL, "sk-test", "api_access_token")
		resp, err := c.Do(ctx, http.MethodGet, fmt.Sprintf("/items?page=%d", page), nil)
		if err != nil {
			return nil, nil, err
		}
		defer resp.Body.Close()

		var body struct {
			Data []Item         `json:"data"`
			Meta map[string]any `json:"meta"`
		}
		json.NewDecoder(resp.Body).Decode(&body)
		pag := &contract.Pagination{Page: page, TotalPages: 3, PerPage: 2, TotalCount: 6}
		return body.Data, pag, nil
	}

	items, pag, err := ListAll(context.Background(), fetcher)
	if err != nil {
		t.Fatalf("ListAll error: %v", err)
	}
	if len(items) != 6 {
		t.Errorf("len(items) = %d, want 6", len(items))
	}
	if pag.TotalCount != 6 {
		t.Errorf("TotalCount = %d, want 6", pag.TotalCount)
	}
}

func TestListAllSinglePage(t *testing.T) {
	type Item struct{ ID int }

	fetcher := func(ctx context.Context, page int) ([]Item, *contract.Pagination, error) {
		return []Item{{ID: 1}}, &contract.Pagination{Page: 1, TotalPages: 1, PerPage: 25, TotalCount: 1}, nil
	}

	items, _, err := ListAll(context.Background(), fetcher)
	if err != nil {
		t.Fatalf("ListAll error: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("len(items) = %d, want 1", len(items))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/ -v -run TestListAll`
Expected: FAIL — `ListAll` not defined

- [ ] **Step 3: Implement ListAll**

```go
// internal/chatwoot/paginate.go
package chatwoot

import (
	"context"

	"github.com/chatwoot/chatwoot-cli/internal/contract"
)

// PageFetcher fetches a single page of results.
type PageFetcher[T any] func(ctx context.Context, page int) ([]T, *contract.Pagination, error)

// ListAll fetches all pages by calling the fetcher repeatedly.
// It returns the merged results and final pagination metadata.
func ListAll[T any](ctx context.Context, fetch PageFetcher[T]) ([]T, *contract.Pagination, error) {
	var all []T
	var lastPag *contract.Pagination

	for page := 1; ; page++ {
		items, pag, err := fetch(ctx, page)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, items...)
		lastPag = pag

		if pag == nil || page >= pag.TotalPages {
			break
		}
	}

	if lastPag != nil {
		lastPag.TotalCount = len(all)
	}

	return all, lastPag, nil
}
```

- [ ] **Step 4: Run pagination tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/ -v -run TestListAll`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/chatwoot/paginate.go internal/chatwoot/paginate_test.go
git commit -m "feat(transport): add generic ListAll pagination helper"
```

---

## Task 12: Shared Model Types

**Files:**
- Create: `internal/chatwoot/models.go`

- [ ] **Step 1: Create shared model types**

```go
// internal/chatwoot/models.go
package chatwoot

// Timestamp is a string alias for Chatwoot API timestamps (ISO 8601).
type Timestamp = string
```

This is intentionally minimal. We add shared types here as they emerge in the
API client tasks. No test needed for type aliases.

- [ ] **Step 2: Verify build**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./...`
Expected: Success

- [ ] **Step 3: Commit**

```bash
git add internal/chatwoot/models.go
git commit -m "feat(transport): add shared model types placeholder"
```

---

## Task 13: Test Utilities

**Files:**
- Modify: `internal/testutil/testutil.go` (currently stub)

- [ ] **Step 1: Implement test helpers**

```go
// internal/testutil/testutil.go
package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockServer creates an httptest.Server that responds with the given status and body
// for any request. Returns the server (caller must defer Close).
func MockServer(t *testing.T, status int, body any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
	}))
}

// MockServerFunc creates an httptest.Server with a custom handler.
func MockServerFunc(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}
```

- [ ] **Step 2: Verify build**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./...`
Expected: Success

- [ ] **Step 3: Commit**

```bash
git add internal/testutil/testutil.go
git commit -m "feat(testutil): add httptest helpers"
```

---

## Task 14: Application API Client — Struct, Constructor, Profile Methods

**Files:**
- Modify: `internal/chatwoot/application/application.go` (currently stub)
- Create: `internal/chatwoot/application/client.go`
- Create: `internal/chatwoot/application/models.go`
- Create: `internal/chatwoot/application/profile.go`
- Create: `internal/chatwoot/application/application_test.go`

- [ ] **Step 1: Write failing test for GetProfile**

```go
// internal/chatwoot/application/application_test.go
package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestGetProfile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/profile" {
			t.Errorf("path = %q, want /api/v1/profile", r.URL.Path)
		}
		if r.Header.Get("api_access_token") != "sk-test" {
			t.Errorf("auth header missing or wrong")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":    1,
			"name":  "Test Agent",
			"email": "agent@test.com",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	profile, err := client.GetProfile(context.Background())
	if err != nil {
		t.Fatalf("GetProfile error: %v", err)
	}
	if profile.ID != 1 {
		t.Errorf("ID = %d, want 1", profile.ID)
	}
	if profile.Name != "Test Agent" {
		t.Errorf("Name = %q, want %q", profile.Name, "Test Agent")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestGetProfile`
Expected: FAIL — `NewClient`, `GetProfile` not defined

- [ ] **Step 3: Implement Client struct**

```go
// internal/chatwoot/application/client.go
package application

import (
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// Client is the Application API client for /api/v1 and /api/v2.
type Client struct {
	transport *chatwoot.Client
	accountID int
}

// NewClient creates an Application API client.
func NewClient(transport *chatwoot.Client, accountID int) *Client {
	return &Client{transport: transport, accountID: accountID}
}
```

- [ ] **Step 4: Implement models**

```go
// internal/chatwoot/application/models.go
package application

import chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"

// Profile represents the authenticated user's profile.
type Profile struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	AccountID int              `json:"account_id,omitempty"`
	Role      string           `json:"role,omitempty"`
	CreatedAt chatwoot.Timestamp `json:"created_at,omitempty"`
}

// Conversation represents a Chatwoot conversation.
type Conversation struct {
	ID               int              `json:"id"`
	AccountID        int              `json:"account_id"`
	InboxID          int              `json:"inbox_id"`
	Status           string           `json:"status"`
	Priority         string           `json:"priority,omitempty"`
	UnreadCount      int              `json:"unread_count,omitempty"`
	CreatedAt        chatwoot.Timestamp `json:"created_at,omitempty"`
}
```

- [ ] **Step 5: Implement GetProfile**

```go
// internal/chatwoot/application/profile.go
package application

import (
	"context"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// GetProfile fetches the authenticated user's profile.
func (c *Client) GetProfile(ctx context.Context) (*Profile, error) {
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, "/api/v1/profile", nil)
	if err != nil {
		return nil, err
	}
	var profile Profile
	if err := chatwoot.DecodeResponse(resp, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}
```

- [ ] **Step 6: Remove old stub file**

Delete `internal/chatwoot/application/application.go` (the stub). Package declaration is now in `client.go`.

- [ ] **Step 7: Run test**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestGetProfile`
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add internal/chatwoot/application/
git commit -m "feat(application): add Application API client with GetProfile"
```

---

## Task 15: Application API Client — Conversations (Proves the Pattern)

**Files:**
- Create: `internal/chatwoot/application/conversations.go`
- Modify: `internal/chatwoot/application/application_test.go`

- [ ] **Step 1: Write failing test for ListConversations**

Add to `application_test.go`:

```go
func TestListConversations(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1/conversations" {
			t.Errorf("path = %q, want /api/v1/accounts/1/conversations", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("page = %q, want 1", r.URL.Query().Get("page"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"payload": []map[string]any{
					{"id": 1, "status": "open"},
					{"id": 2, "status": "resolved"},
				},
				"meta": map[string]any{
					"all_count": 42,
				},
			},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	convos, err := client.ListConversations(context.Background(), ListConversationsOpts{Page: 1})
	if err != nil {
		t.Fatalf("ListConversations error: %v", err)
	}
	if len(convos) != 2 {
		t.Errorf("len = %d, want 2", len(convos))
	}
	if convos[0].ID != 1 {
		t.Errorf("convos[0].ID = %d, want 1", convos[0].ID)
	}
}

func TestGetConversation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/1/conversations/42" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":     42,
			"status": "open",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	convo, err := client.GetConversation(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetConversation error: %v", err)
	}
	if convo.ID != 42 {
		t.Errorf("ID = %d, want 42", convo.ID)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run "TestListConversations|TestGetConversation"`
Expected: FAIL

- [ ] **Step 3: Implement conversations methods**

```go
// internal/chatwoot/application/conversations.go
package application

import (
	"context"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListConversationsOpts holds parameters for listing conversations.
type ListConversationsOpts struct {
	Page    int
	Status  string
	InboxID int
}

// ListConversations returns conversations for the configured account.
func (c *Client) ListConversations(ctx context.Context, opts ListConversationsOpts) ([]Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations?page=%d", c.accountID, opts.Page)
	if opts.Status != "" {
		path += "&status=" + opts.Status
	}
	if opts.InboxID != 0 {
		path += fmt.Sprintf("&inbox_id=%d", opts.InboxID)
	}

	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Chatwoot wraps conversations in data.payload
	var body struct {
		Data struct {
			Payload []Conversation `json:"payload"`
		} `json:"data"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, err
	}
	return body.Data.Payload, nil
}

// GetConversation fetches a single conversation by ID.
func (c *Client) GetConversation(ctx context.Context, id int) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d", c.accountID, id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var convo Conversation
	if err := chatwoot.DecodeResponse(resp, &convo); err != nil {
		return nil, err
	}
	return &convo, nil
}
```

- [ ] **Step 4: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/chatwoot/application/
git commit -m "feat(application): add ListConversations and GetConversation"
```

---

## Task 16: Platform API Client

**Files:**
- Create: `internal/chatwoot/platform/client.go`
- Create: `internal/chatwoot/platform/accounts.go`
- Create: `internal/chatwoot/platform/models.go`
- Create: `internal/chatwoot/platform/platform_test.go`
- Delete: `internal/chatwoot/platform/platform.go` (stub)

- [ ] **Step 1: Write failing test for GetAccount**

```go
// internal/chatwoot/platform/platform_test.go
package platform

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestGetAccount(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform/api/v1/accounts/1" {
			t.Errorf("path = %q, want /platform/api/v1/accounts/1", r.URL.Path)
		}
		if r.Header.Get("api_access_token") != "pk-test" {
			t.Errorf("auth header missing or wrong")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":   1,
			"name": "Test Account",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	account, err := client.GetAccount(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetAccount error: %v", err)
	}
	if account.ID != 1 {
		t.Errorf("ID = %d, want 1", account.ID)
	}
	if account.Name != "Test Account" {
		t.Errorf("Name = %q, want %q", account.Name, "Test Account")
	}
}

func TestCreateAccount(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/platform/api/v1/accounts" {
			t.Errorf("path = %q", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		json.NewEncoder(w).Encode(map[string]any{
			"id":   99,
			"name": body["name"],
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")
	client := NewClient(transport)

	account, err := client.CreateAccount(context.Background(), CreateAccountOpts{Name: "New Account"})
	if err != nil {
		t.Fatalf("CreateAccount error: %v", err)
	}
	if account.ID != 99 {
		t.Errorf("ID = %d, want 99", account.ID)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/platform/ -v`
Expected: FAIL

- [ ] **Step 3: Implement platform client**

```go
// internal/chatwoot/platform/client.go
package platform

import chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"

// Client is the Platform API client for /platform/api/v1.
type Client struct {
	transport *chatwoot.Client
}

// NewClient creates a Platform API client.
func NewClient(transport *chatwoot.Client) *Client {
	return &Client{transport: transport}
}
```

```go
// internal/chatwoot/platform/models.go
package platform

// Account represents a Chatwoot platform account.
type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CreateAccountOpts holds parameters for creating an account.
type CreateAccountOpts struct {
	Name string `json:"name"`
}
```

```go
// internal/chatwoot/platform/accounts.go
package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// GetAccount fetches an account by ID.
func (c *Client) GetAccount(ctx context.Context, id int) (*Account, error) {
	path := fmt.Sprintf("/platform/api/v1/accounts/%d", id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var account Account
	if err := chatwoot.DecodeResponse(resp, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// CreateAccount creates a new account.
func (c *Client) CreateAccount(ctx context.Context, opts CreateAccountOpts) (*Account, error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create account: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, "/platform/api/v1/accounts", body)
	if err != nil {
		return nil, err
	}
	var account Account
	if err := chatwoot.DecodeResponse(resp, &account); err != nil {
		return nil, err
	}
	return &account, nil
}
```

- [ ] **Step 4: Delete old stub and run tests**

Delete `internal/chatwoot/platform/platform.go`.

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/platform/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/chatwoot/platform/
git commit -m "feat(platform): add Platform API client with GetAccount and CreateAccount"
```

---

## Task 17: Client (Public) API Client

**Files:**
- Create: `internal/chatwoot/clientapi/client.go`
- Create: `internal/chatwoot/clientapi/contacts.go`
- Create: `internal/chatwoot/clientapi/models.go`
- Create: `internal/chatwoot/clientapi/clientapi_test.go`
- Delete: `internal/chatwoot/clientapi/clientapi.go` (stub)

- [ ] **Step 1: Write failing test for CreateContact**

```go
// internal/chatwoot/clientapi/clientapi_test.go
package clientapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestCreateContact(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/api/v1/inboxes/inbox-abc/contacts" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"source_id":          "contact-xyz",
			"name":               "Test User",
			"pubsub_token":       "token-123",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "", "")
	client := NewClient(transport, "inbox-abc")

	contact, err := client.CreateContact(context.Background(), CreateContactOpts{Name: "Test User"})
	if err != nil {
		t.Fatalf("CreateContact error: %v", err)
	}
	if contact.SourceID != "contact-xyz" {
		t.Errorf("SourceID = %q, want %q", contact.SourceID, "contact-xyz")
	}
}

func TestGetContact(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/api/v1/inboxes/inbox-abc/contacts/contact-xyz" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"source_id": "contact-xyz",
			"name":      "Existing User",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "", "")
	client := NewClient(transport, "inbox-abc")

	contact, err := client.GetContact(context.Background(), "contact-xyz")
	if err != nil {
		t.Fatalf("GetContact error: %v", err)
	}
	if contact.Name != "Existing User" {
		t.Errorf("Name = %q, want %q", contact.Name, "Existing User")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/clientapi/ -v`
Expected: FAIL

- [ ] **Step 3: Implement client API client**

```go
// internal/chatwoot/clientapi/client.go
package clientapi

import chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"

// Client is the public Client API client for /public/api/v1.
// Auth is identifier-based (inbox + contact), not token-based.
type Client struct {
	transport       *chatwoot.Client
	inboxIdentifier string
}

// NewClient creates a Client API client for the given inbox.
func NewClient(transport *chatwoot.Client, inboxIdentifier string) *Client {
	return &Client{transport: transport, inboxIdentifier: inboxIdentifier}
}
```

```go
// internal/chatwoot/clientapi/models.go
package clientapi

// Contact represents a public API contact.
type Contact struct {
	SourceID    string `json:"source_id"`
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	PubsubToken string `json:"pubsub_token,omitempty"`
}

// CreateContactOpts holds parameters for creating a contact.
type CreateContactOpts struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone_number,omitempty"`
}
```

```go
// internal/chatwoot/clientapi/contacts.go
package clientapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// CreateContact creates a new contact in the inbox.
func (c *Client) CreateContact(ctx context.Context, opts CreateContactOpts) (*Contact, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts", c.inboxIdentifier)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create contact: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var contact Contact
	if err := chatwoot.DecodeResponse(resp, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

// GetContact fetches a contact by identifier.
func (c *Client) GetContact(ctx context.Context, contactIdentifier string) (*Contact, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s", c.inboxIdentifier, contactIdentifier)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var contact Contact
	if err := chatwoot.DecodeResponse(resp, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}
```

- [ ] **Step 4: Delete old stub and run tests**

Delete `internal/chatwoot/clientapi/clientapi.go`.

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/clientapi/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/chatwoot/clientapi/
git commit -m "feat(clientapi): add Client API client with CreateContact and GetContact"
```

---

## Task 18: Clean Up Reports Stub and Run Full Test Suite

**Files:**
- Verify: `internal/chatwoot/reports/reports.go` (keep as stub — Sprint D scope)
- No changes needed to reports

- [ ] **Step 1: Run full test suite**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./... -v`
Expected: All tests PASS

- [ ] **Step 2: Run vet and build**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go vet ./... && go build ./cmd/chatwoot/`
Expected: No errors

- [ ] **Step 3: Verify no Cobra imports in non-CLI packages**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && grep -r "spf13/cobra" internal/config/ internal/credentials/ internal/auth/ internal/chatwoot/`
Expected: No matches — these packages must not depend on Cobra

- [ ] **Step 4: Commit any cleanup**

If any cleanup was needed, commit it. Otherwise, skip this step.

---

## Sprint A Exit Criteria Checklist

- [ ] `go test ./...` passes with all tests green
- [ ] `go vet ./...` clean
- [ ] `go build ./cmd/chatwoot/` succeeds
- [ ] Config resolution follows: flag > env > config > default
- [ ] Credential resolution follows: env > keychain > file (0600) > error
- [ ] Keychain store uses `chatwoot-cli/<profile>/<mode>` key structure
- [ ] File store enforces 0600 permissions
- [ ] Transport retries on 429 and 5xx with exponential backoff
- [ ] Transport maps HTTP status to contract error codes
- [ ] Application client proves pattern with GetProfile + ListConversations + GetConversation
- [ ] Platform client proves pattern with GetAccount + CreateAccount
- [ ] Client API proves pattern with CreateContact + GetContact
- [ ] ListAll pagination helper works with generics
- [ ] No Cobra dependency in config, credentials, auth, or chatwoot packages

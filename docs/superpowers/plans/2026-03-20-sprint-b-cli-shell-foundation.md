# Sprint B: CLI Shell + Foundation Commands (Layers 5 + 6) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Wire Sprint A plumbing into Cobra commands and ship foundation commands (auth set/status/clear, application profile get/update) that prove the full pipeline end-to-end.

**Architecture:** Per-command context resolution helpers (not PersistentPreRun) load config and credentials on demand. Auth commands use a lightweight profile-name-only path. Application commands use the full pipeline: ResolveContext → ResolveAuth → transport Client → API client → contract envelope.

**Tech Stack:** Go 1.26.1, `spf13/cobra` (CLI), `spf13/viper` (config), `zalando/go-keyring` (keychain), `log/slog` (diagnostics), `encoding/json` (serialization), `net/http/httptest` (testing)

**Spec:** `docs/superpowers/specs/2026-03-20-sprint-b-cli-shell-foundation.md`

---

## File Structure

### Layer 5: CLI Shell

| File | Responsibility |
|------|---------------|
| `internal/cli/root.go` | MODIFY: add `--profile`, `--base-url`, `--account-id` persistent flags; wire slog level in PersistentPreRun; register auth and application groups |
| `internal/cli/context.go` | NEW: `RuntimeContext` struct, `ResolveContext`, `ResolveAuth`, `ResolveProfileName` helpers, flag accessors |
| `internal/cli/context_test.go` | NEW: context resolution tests |

### Layer 5+6: Auth Commands

| File | Responsibility |
|------|---------------|
| `internal/cli/auth/auth.go` | NEW: `authCmd` group command definition, exported for root registration |
| `internal/cli/auth/set.go` | NEW: `auth set` handler — store credential |
| `internal/cli/auth/set_test.go` | NEW: auth set tests |
| `internal/cli/auth/status.go` | NEW: `auth status` handler — inspect credential state |
| `internal/cli/auth/status_test.go` | NEW: auth status tests |
| `internal/cli/auth/clear.go` | NEW: `auth clear` handler — remove credential |
| `internal/cli/auth/clear_test.go` | NEW: auth clear tests |

### Layer 5+6: Application Commands

| File | Responsibility |
|------|---------------|
| `internal/cli/application/application.go` | NEW: `applicationCmd` group with `profile` subgroup, exported for root registration |
| `internal/cli/application/profile.go` | NEW: `profile get` and `profile update` handlers |
| `internal/cli/application/profile_test.go` | NEW: profile command tests with httptest |

### API Client Addition

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/application/profile.go` | MODIFY: add `UpdateProfile` method |
| `internal/chatwoot/application/models.go` | MODIFY: add `UpdateProfileOpts` type |
| `internal/chatwoot/application/application_test.go` | MODIFY: add `TestUpdateProfile` test |

---

## Task 1: Root Command — New Global Flags and Slog Wiring

**Files:**
- Modify: `internal/cli/root.go`

- [ ] **Step 1: Add new flag variables and registration**

Add to the `var` block in `internal/cli/root.go`:

```go
var (
	prettyFlag    bool
	verboseFlag   bool
	profileFlag   string
	baseURLFlag   string
	accountIDFlag int
)
```

Add to `init()`:

```go
rootCmd.PersistentFlags().StringVar(&profileFlag, "profile", "", "Select named profile")
rootCmd.PersistentFlags().StringVar(&baseURLFlag, "base-url", "", "Override base URL")
rootCmd.PersistentFlags().IntVar(&accountIDFlag, "account-id", 0, "Override account ID")
```

- [ ] **Step 2: Wire slog level in PersistentPreRun**

Add to the end of the existing `PersistentPreRun` function (after the `--verbose` env check):

```go
if verboseFlag {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
}
```

Add `"log/slog"` to imports.

- [ ] **Step 3: Verify build**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./cmd/chatwoot/`
Expected: Success

- [ ] **Step 4: Verify existing tests still pass**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/ -v`
Expected: PASS (existing version tests unaffected)

- [ ] **Step 5: Commit**

```bash
git add internal/cli/root.go
git commit -m "feat(cli): add --profile, --base-url, --account-id flags and wire slog verbose"
```

---

## Task 2: Runtime Context Resolution

**Files:**
- Create: `internal/cli/context.go`
- Create: `internal/cli/context_test.go`

- [ ] **Step 1: Write failing test for ResolveContext with config file**

```go
// internal/cli/context_test.go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/ -v -run TestResolveContext`
Expected: FAIL — `resolveContextFromPath` not defined

- [ ] **Step 3: Implement RuntimeContext and resolveContextFromPath**

```go
// internal/cli/context.go
package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chatwoot/chatwoot-cli/internal/auth"
	"github.com/chatwoot/chatwoot-cli/internal/config"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

// RuntimeContext holds the resolved runtime configuration for a command.
type RuntimeContext struct {
	ProfileName string
	BaseURL     string
	AccountID   int
}

// ResolveContext resolves the full runtime context from flags, env, and config.
func ResolveContext(cmd *cobra.Command) (*RuntimeContext, error) {
	flagProfile, _ := cmd.Flags().GetString("profile")
	flagBaseURL, _ := cmd.Flags().GetString("base-url")
	flagAccountID, _ := cmd.Flags().GetInt("account-id")

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	cfgPath := filepath.Join(cfgDir, "chatwoot-cli", "config.yaml")

	return resolveContextFromPath(cfgPath, flagProfile, flagBaseURL, flagAccountID)
}

// resolveContextFromPath is the testable core of ResolveContext.
func resolveContextFromPath(cfgPath, flagProfile, flagBaseURL string, flagAccountID int) (*RuntimeContext, error) {
	cfg, err := config.LoadFrom(cfgPath)
	if err != nil {
		// Config file missing is OK — proceed with empty config.
		// Viper wraps the underlying error, so check with errors.As.
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && errors.Is(pathErr.Err, os.ErrNotExist) {
			cfg = &config.Config{}
		} else if _, statErr := os.Stat(cfgPath); errors.Is(statErr, os.ErrNotExist) {
			// Fallback: file genuinely doesn't exist
			cfg = &config.Config{}
		} else {
			return nil, fmt.Errorf("config_error: %w", err)
		}
	}

	profileName, profile, err := cfg.ResolveProfile(flagProfile)
	if err != nil {
		// If config is empty and profile not found, use empty profile
		if cfg.Profiles == nil || len(cfg.Profiles) == 0 {
			profileName = flagProfile
			if profileName == "" {
				profileName = os.Getenv("CHATWOOT_PROFILE")
			}
			if profileName == "" {
				profileName = "default"
			}
			profile = config.Profile{}
		} else {
			return nil, fmt.Errorf("config_error: %w", err)
		}
	}

	resolved := config.ResolveOverrides(profile, flagBaseURL, flagAccountID)

	if resolved.BaseURL == "" {
		return nil, fmt.Errorf("no base URL configured — set base_url in profile or use --base-url flag")
	}

	return &RuntimeContext{
		ProfileName: profileName,
		BaseURL:     resolved.BaseURL,
		AccountID:   resolved.AccountID,
	}, nil
}

// ResolveProfileName resolves just the profile name (lightweight path for auth commands).
func ResolveProfileName(cmd *cobra.Command) string {
	flagProfile, _ := cmd.Flags().GetString("profile")

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	cfgPath := filepath.Join(cfgDir, "chatwoot-cli", "config.yaml")

	cfg, err := config.LoadFrom(cfgPath)
	if err != nil {
		cfg = &config.Config{}
	}

	name, _, err := cfg.ResolveProfile(flagProfile)
	if err != nil {
		// Fall back to flag, env, or "default"
		name = flagProfile
		if name == "" {
			name = os.Getenv("CHATWOOT_PROFILE")
		}
		if name == "" {
			name = "default"
		}
	}
	return name
}

// ResolveAuth resolves credentials for the given profile and auth mode.
func ResolveAuth(profileName string, mode credentials.AuthMode) (auth.TokenAuth, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	credPath := filepath.Join(cfgDir, "chatwoot-cli", "credentials.yaml")

	resolver := credentials.NewResolver(
		&credentials.EnvStore{},
		credentials.NewKeychainStore(),
		credentials.NewFileStore(credPath),
	)

	switch mode {
	case credentials.ModeApplication:
		return auth.ResolveApplication(resolver, profileName)
	case credentials.ModePlatform:
		return auth.ResolvePlatform(resolver, profileName)
	default:
		return auth.TokenAuth{}, fmt.Errorf("unknown auth mode: %s", mode)
	}
}

// WriteError writes an error envelope to stdout and returns an error for Cobra.
func WriteError(cmd *cobra.Command, code, message string) error {
	resp := contract.Err(code, message)
	_ = contract.Write(cmd.OutOrStdout(), resp, prettyFlag)
	return fmt.Errorf(message)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/ -v -run TestResolveContextFromConfig`
Expected: PASS

- [ ] **Step 5: Write failing test for missing config with flags**

```go
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
```

- [ ] **Step 6: Run test to verify it passes**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/ -v -run TestResolveContextMissing`
Expected: PASS

- [ ] **Step 7: Write failing test for missing config with env vars**

```go
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
```

- [ ] **Step 8: Run test to verify it passes**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/ -v -run TestResolveContextMissingConfigWithEnv`
Expected: PASS

- [ ] **Step 9: Write failing test for no BaseURL**

```go
func TestResolveContextNoBaseURL(t *testing.T) {
	t.Setenv("CHATWOOT_BASE_URL", "")
	_, err := resolveContextFromPath("/nonexistent/config.yaml", "", "", 0)
	if err == nil {
		t.Fatal("expected error for missing base URL")
	}
}
```

- [ ] **Step 10: Run test to verify it passes**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/ -v -run TestResolveContextNoBaseURL`
Expected: PASS

- [ ] **Step 11: Write failing test for flag overrides profile**

```go
func TestResolveContextFlagOverridesProfile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	os.WriteFile(cfgPath, []byte(`
default_profile: work
profiles:
  work:
    base_url: https://app.chatwoot.com
    account_id: 1
`), 0644)

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
```

- [ ] **Step 12: Run all context tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/ -v -run TestResolveContext`
Expected: PASS

- [ ] **Step 13: Commit**

```bash
git add internal/cli/context.go internal/cli/context_test.go
git commit -m "feat(cli): add runtime context resolution with config/flag/env precedence"
```

---

## Task 3: Auth Group Command and Auth Set

**Files:**
- Create: `internal/cli/auth/auth.go`
- Create: `internal/cli/auth/set.go`
- Create: `internal/cli/auth/set_test.go`
- Modify: `internal/cli/root.go` (register auth group)

- [ ] **Step 1: Create auth group command**

```go
// internal/cli/auth/auth.go
package auth

import "github.com/spf13/cobra"

// Cmd is the auth command group.
var Cmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication credentials",
}
```

- [ ] **Step 2: Write failing test for auth set**

```go
// internal/cli/auth/set_test.go
package auth

import (
	"bytes"
	"encoding/json"
	"testing"

	keyring "github.com/zalando/go-keyring"
)

func TestAuthSetKeychain(t *testing.T) {
	keyring.MockInit()

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"auth", "set", "--mode", "application", "--token", "sk-test-123"})
	err := Cmd.Root().Execute()
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	var resp map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		t.Fatalf("JSON decode error: %v\nOutput: %s", err, stdout.String())
	}
	if resp["ok"] != true {
		t.Errorf("ok = %v, want true", resp["ok"])
	}
	data := resp["data"].(map[string]any)
	if data["mode"] != "application" {
		t.Errorf("mode = %v, want application", data["mode"])
	}
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/auth/ -v -run TestAuthSet`
Expected: FAIL — set command not defined

- [ ] **Step 4: Implement auth set command**

```go
// internal/cli/auth/set.go
package auth

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Store a credential for the active profile",
	RunE:  runSet,
}

func init() {
	setCmd.Flags().String("mode", "", "Auth mode: application or platform (required)")
	setCmd.Flags().String("token", "", "Token value (required)")
	_ = setCmd.MarkFlagRequired("mode")
	_ = setCmd.MarkFlagRequired("token")
	Cmd.AddCommand(setCmd)
}

func runSet(cmd *cobra.Command, args []string) error {
	modeStr, _ := cmd.Flags().GetString("mode")
	token, _ := cmd.Flags().GetString("token")

	mode := credentials.AuthMode(modeStr)
	if mode != credentials.ModeApplication && mode != credentials.ModePlatform {
		return fmt.Errorf("invalid mode %q: must be \"application\" or \"platform\"", modeStr)
	}

	profileName := resolveProfileNameForAuth(cmd)

	source, err := storeCredential(profileName, mode, token)
	if err != nil {
		resp := contract.Err(contract.ErrCodeAuth, err.Error())
		return contract.Write(cmd.OutOrStdout(), resp, prettyFromRoot(cmd))
	}

	resp := contract.Success(map[string]string{
		"profile": profileName,
		"mode":    string(mode),
		"source":  string(source),
	})
	return contract.Write(cmd.OutOrStdout(), resp, prettyFromRoot(cmd))
}

func storeCredential(profile string, mode credentials.AuthMode, token string) (credentials.Source, error) {
	ks := credentials.NewKeychainStore()
	if err := ks.Set(profile, mode, token); err == nil {
		return credentials.SourceKeychain, nil
	} else {
		slog.Warn("keychain unavailable, using file-based credential storage", "error", err)
	}

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	credPath := filepath.Join(cfgDir, "chatwoot-cli", "credentials.yaml")

	fs := credentials.NewFileStore(credPath)
	if err := fs.Set(profile, mode, token); err != nil {
		return "", fmt.Errorf("store credential: %w", err)
	}
	return credentials.SourceFile, nil
}

// resolveProfileNameForAuth resolves just the profile name for auth commands.
// Does not validate BaseURL (auth commands don't need it).
func resolveProfileNameForAuth(cmd *cobra.Command) string {
	flagProfile, _ := cmd.Flags().GetString("profile")
	if flagProfile != "" {
		return flagProfile
	}
	if env := os.Getenv("CHATWOOT_PROFILE"); env != "" {
		return env
	}

	cfgDir, _ := os.UserConfigDir()
	if cfgDir == "" {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	cfgPath := filepath.Join(cfgDir, "chatwoot-cli", "config.yaml")

	cfg, err := config.LoadFrom(cfgPath)
	if err != nil {
		return "default"
	}
	if cfg.DefaultProfile != "" {
		return cfg.DefaultProfile
	}
	return "default"
}

func prettyFromRoot(cmd *cobra.Command) bool {
	pretty, _ := cmd.Root().PersistentFlags().GetBool("pretty")
	return pretty
}
```

Add the config import: `"github.com/chatwoot/chatwoot-cli/internal/config"`

- [ ] **Step 5: Register auth group in root.go**

Add to `init()` in `internal/cli/root.go`:

```go
rootCmd.AddCommand(cliauth.Cmd)
```

Add import: `cliauth "github.com/chatwoot/chatwoot-cli/internal/cli/auth"`

- [ ] **Step 6: Run test to verify it passes**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/auth/ -v -run TestAuthSet`
Expected: PASS

- [ ] **Step 7: Verify build and all tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./cmd/chatwoot/ && go test ./... 2>&1 | grep -E "^(ok|FAIL)"`
Expected: All pass

- [ ] **Step 8: Commit**

```bash
git add internal/cli/auth/ internal/cli/root.go
git commit -m "feat(cli): add auth group and auth set command"
```

---

## Task 4: Auth Status Command

**Files:**
- Create: `internal/cli/auth/status.go`
- Create: `internal/cli/auth/status_test.go`

- [ ] **Step 1: Write failing test for auth status**

```go
// internal/cli/auth/status_test.go
package auth

import (
	"bytes"
	"encoding/json"
	"testing"

	keyring "github.com/zalando/go-keyring"
)

func TestAuthStatus(t *testing.T) {
	keyring.MockInit()
	ks := credentials.NewKeychainStore()
	_ = ks.Set("default", credentials.ModeApplication, "sk-test")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"auth", "status"})
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
	data := resp["data"].(map[string]any)
	creds := data["credentials"].(map[string]any)
	app := creds["application"].(map[string]any)
	if app["status"] != "configured" {
		t.Errorf("application status = %v, want configured", app["status"])
	}
}
```

Add `"github.com/chatwoot/chatwoot-cli/internal/credentials"` to test imports.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/auth/ -v -run TestAuthStatus`
Expected: FAIL

- [ ] **Step 3: Implement auth status command**

```go
// internal/cli/auth/status.go
package auth

import (
	"os"
	"path/filepath"

	"github.com/chatwoot/chatwoot-cli/internal/config"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show credential status for the active profile",
	RunE:  runStatus,
}

func init() {
	Cmd.AddCommand(statusCmd)
}

type credentialStatus struct {
	Status string `json:"status"`
	Source string `json:"source,omitempty"`
}

type statusData struct {
	Profile     string                       `json:"profile"`
	BaseURL     string                       `json:"base_url,omitempty"`
	AccountID   int                          `json:"account_id,omitempty"`
	Credentials map[string]credentialStatus   `json:"credentials"`
}

func runStatus(cmd *cobra.Command, args []string) error {
	profileName := resolveProfileNameForAuth(cmd)

	// Try to load config for display (non-fatal if missing)
	var baseURL string
	var accountID int
	cfgDir, _ := os.UserConfigDir()
	if cfgDir == "" {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	cfgPath := filepath.Join(cfgDir, "chatwoot-cli", "config.yaml")
	if cfg, err := config.LoadFrom(cfgPath); err == nil {
		if _, profile, err := cfg.ResolveProfile(profileName); err == nil {
			baseURL = profile.BaseURL
			accountID = profile.AccountID
		}
	}

	// Override with env/flags
	resolved := config.ResolveOverrides(config.Profile{BaseURL: baseURL, AccountID: accountID}, "", 0)
	baseURL = resolved.BaseURL
	accountID = resolved.AccountID

	// Probe credentials
	credPath := filepath.Join(cfgDir, "chatwoot-cli", "credentials.yaml")
	resolver := credentials.NewResolver(
		&credentials.EnvStore{},
		credentials.NewKeychainStore(),
		credentials.NewFileStore(credPath),
	)

	creds := map[string]credentialStatus{
		"application": probeCredential(resolver, profileName, credentials.ModeApplication),
		"platform":    probeCredential(resolver, profileName, credentials.ModePlatform),
	}

	result := statusData{
		Profile:     profileName,
		BaseURL:     baseURL,
		AccountID:   accountID,
		Credentials: creds,
	}

	resp := contract.Success(result)
	return contract.Write(cmd.OutOrStdout(), resp, prettyFromRoot(cmd))
}

func probeCredential(resolver *credentials.Resolver, profile string, mode credentials.AuthMode) credentialStatus {
	_, source, err := resolver.Get(profile, mode)
	if err != nil {
		return credentialStatus{Status: "not_configured"}
	}
	return credentialStatus{Status: "configured", Source: string(source)}
}
```

- [ ] **Step 4: Run test**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/auth/ -v -run TestAuthStatus`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/cli/auth/status.go internal/cli/auth/status_test.go
git commit -m "feat(cli): add auth status command"
```

---

## Task 5: Auth Clear Command

**Files:**
- Create: `internal/cli/auth/clear.go`
- Create: `internal/cli/auth/clear_test.go`

- [ ] **Step 1: Write failing test for auth clear**

```go
// internal/cli/auth/clear_test.go
package auth

import (
	"bytes"
	"encoding/json"
	"testing"

	keyring "github.com/zalando/go-keyring"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
)

func TestAuthClearSingleMode(t *testing.T) {
	keyring.MockInit()
	ks := credentials.NewKeychainStore()
	_ = ks.Set("default", credentials.ModeApplication, "sk-test")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"auth", "clear", "--mode", "application"})
	err := Cmd.Root().Execute()
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	var resp map[string]any
	json.Unmarshal(stdout.Bytes(), &resp)
	if resp["ok"] != true {
		t.Errorf("ok = %v, want true", resp["ok"])
	}
	data := resp["data"].(map[string]any)
	cleared := data["cleared"].([]any)
	if len(cleared) != 1 || cleared[0] != "application" {
		t.Errorf("cleared = %v, want [application]", cleared)
	}

	// Verify credential was actually removed
	_, err = ks.Get("default", credentials.ModeApplication)
	if err != credentials.ErrNotFound {
		t.Errorf("credential still exists after clear")
	}
}

func TestAuthClearAll(t *testing.T) {
	keyring.MockInit()
	ks := credentials.NewKeychainStore()
	_ = ks.Set("default", credentials.ModeApplication, "sk-test")
	_ = ks.Set("default", credentials.ModePlatform, "pk-test")

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"auth", "clear", "--all"})
	err := Cmd.Root().Execute()
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	var resp map[string]any
	json.Unmarshal(stdout.Bytes(), &resp)
	data := resp["data"].(map[string]any)
	cleared := data["cleared"].([]any)
	if len(cleared) != 2 {
		t.Errorf("cleared = %v, want 2 items", cleared)
	}
}

func TestAuthClearRequiresModeOrAll(t *testing.T) {
	keyring.MockInit()

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"auth", "clear"})
	err := Cmd.Root().Execute()
	if err == nil {
		t.Fatal("expected error when neither --mode nor --all provided")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/auth/ -v -run TestAuthClear`
Expected: FAIL

- [ ] **Step 3: Implement auth clear command**

```go
// internal/cli/auth/clear.go
package auth

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove credentials for the active profile",
	RunE:  runClear,
}

func init() {
	clearCmd.Flags().String("mode", "", "Auth mode to clear: application or platform")
	clearCmd.Flags().Bool("all", false, "Clear all credential modes")
	Cmd.AddCommand(clearCmd)
}

func runClear(cmd *cobra.Command, args []string) error {
	modeStr, _ := cmd.Flags().GetString("mode")
	all, _ := cmd.Flags().GetBool("all")

	if modeStr == "" && !all {
		return fmt.Errorf("requires --mode or --all flag")
	}

	profileName := resolveProfileNameForAuth(cmd)

	var modes []credentials.AuthMode
	if all {
		modes = []credentials.AuthMode{credentials.ModeApplication, credentials.ModePlatform}
	} else {
		mode := credentials.AuthMode(modeStr)
		if mode != credentials.ModeApplication && mode != credentials.ModePlatform {
			return fmt.Errorf("invalid mode %q: must be \"application\" or \"platform\"", modeStr)
		}
		modes = []credentials.AuthMode{mode}
	}

	var cleared []string
	for _, mode := range modes {
		if err := deleteCredential(profileName, mode); err != nil {
			slog.Warn("failed to delete credential", "mode", mode, "error", err)
		}
		cleared = append(cleared, string(mode))
	}

	resp := contract.Success(map[string]any{
		"profile": profileName,
		"cleared": cleared,
	})
	return contract.Write(cmd.OutOrStdout(), resp, prettyFromRoot(cmd))
}

func deleteCredential(profile string, mode credentials.AuthMode) error {
	ks := credentials.NewKeychainStore()
	if err := ks.Delete(profile, mode); err == nil {
		return nil
	}

	cfgDir, _ := os.UserConfigDir()
	if cfgDir == "" {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	credPath := filepath.Join(cfgDir, "chatwoot-cli", "credentials.yaml")
	fs := credentials.NewFileStore(credPath)
	return fs.Delete(profile, mode)
}
```

- [ ] **Step 4: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/auth/ -v -run TestAuthClear`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/cli/auth/clear.go internal/cli/auth/clear_test.go
git commit -m "feat(cli): add auth clear command"
```

---

## Task 6: UpdateProfile API Client Method

**Files:**
- Modify: `internal/chatwoot/application/models.go`
- Modify: `internal/chatwoot/application/profile.go`
- Modify: `internal/chatwoot/application/application_test.go`

- [ ] **Step 1: Add UpdateProfileOpts to models.go**

Add to `internal/chatwoot/application/models.go`:

```go
// UpdateProfileOpts holds optional fields for updating a profile.
// Pointer fields are used so only explicitly set values are serialized.
type UpdateProfileOpts struct {
	Name         *string `json:"name,omitempty"`
	Email        *string `json:"email,omitempty"`
	Availability *string `json:"availability,omitempty"`
}
```

- [ ] **Step 2: Write failing test for UpdateProfile**

Add to `internal/chatwoot/application/application_test.go`:

```go
func TestUpdateProfile(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		if r.URL.Path != "/api/v1/profile" {
			t.Errorf("path = %q, want /api/v1/profile", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":    1,
			"name":  "Updated Name",
			"email": "agent@test.com",
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	name := "Updated Name"
	profile, err := client.UpdateProfile(context.Background(), UpdateProfileOpts{Name: &name})
	if err != nil {
		t.Fatalf("UpdateProfile error: %v", err)
	}
	if gotMethod != "PATCH" {
		t.Errorf("method = %q, want PATCH", gotMethod)
	}
	if gotBody["name"] != "Updated Name" {
		t.Errorf("body name = %v, want Updated Name", gotBody["name"])
	}
	if profile.Name != "Updated Name" {
		t.Errorf("Name = %q, want %q", profile.Name, "Updated Name")
	}
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestUpdateProfile`
Expected: FAIL — `UpdateProfile` not defined

- [ ] **Step 4: Implement UpdateProfile**

Add to `internal/chatwoot/application/profile.go`:

```go
// UpdateProfile updates the authenticated user's profile.
func (c *Client) UpdateProfile(ctx context.Context, opts UpdateProfileOpts) (*Profile, error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update profile: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, "/api/v1/profile", body)
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

Add `"encoding/json"` and `"fmt"` to profile.go imports.

- [ ] **Step 5: Run test**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestUpdateProfile`
Expected: PASS

- [ ] **Step 6: Run all application tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v`
Expected: All PASS

- [ ] **Step 7: Commit**

```bash
git add internal/chatwoot/application/
git commit -m "feat(application): add UpdateProfile method"
```

---

## Task 7: Application Group and Profile Commands

**Files:**
- Create: `internal/cli/application/application.go`
- Create: `internal/cli/application/profile.go`
- Create: `internal/cli/application/profile_test.go`
- Modify: `internal/cli/root.go` (register application group)

- [ ] **Step 1: Create application group with profile subgroup**

```go
// internal/cli/application/application.go
package application

import "github.com/spf13/cobra"

// Cmd is the application command group.
var Cmd = &cobra.Command{
	Use:   "application",
	Short: "Application API commands (agent/admin)",
}

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage authenticated user profile",
}

func init() {
	Cmd.AddCommand(profileCmd)
}
```

- [ ] **Step 2: Write failing test for profile get**

```go
// internal/cli/application/profile_test.go
package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	keyring "github.com/zalando/go-keyring"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
)

func TestProfileGet(t *testing.T) {
	// Set up mock API server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/profile" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Header.Get("api_access_token") == "" {
			t.Error("missing auth header")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":    1,
			"name":  "Test Agent",
			"email": "agent@test.com",
		})
	}))
	defer srv.Close()

	// Set up config pointing to test server
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "chatwoot-cli", "config.yaml")
	os.MkdirAll(filepath.Dir(cfgPath), 0755)
	os.WriteFile(cfgPath, []byte(`
default_profile: work
profiles:
  work:
    base_url: `+srv.URL+`
    account_id: 1
`), 0644)

	// Set up mock keychain with test token
	keyring.MockInit()
	ks := credentials.NewKeychainStore()
	_ = ks.Set("work", credentials.ModeApplication, "sk-test-token")

	// Override config dir
	t.Setenv("XDG_CONFIG_HOME", dir)

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"application", "profile", "get"})
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
	data := resp["data"].(map[string]any)
	if data["name"] != "Test Agent" {
		t.Errorf("name = %v, want Test Agent", data["name"])
	}
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/ -v -run TestProfileGet`
Expected: FAIL

- [ ] **Step 4: Implement profile get and update commands**

```go
// internal/cli/application/profile.go
package application

import (
	"context"
	"fmt"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var profileGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get authenticated user profile",
	RunE:  runProfileGet,
}

var profileUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update authenticated user profile",
	RunE:  runProfileUpdate,
}

func init() {
	profileUpdateCmd.Flags().String("name", "", "Display name")
	profileUpdateCmd.Flags().String("email", "", "Email address")
	profileUpdateCmd.Flags().String("availability", "", "Availability: online, offline, busy")
	profileCmd.AddCommand(profileGetCmd)
	profileCmd.AddCommand(profileUpdateCmd)
}

func runProfileGet(cmd *cobra.Command, args []string) error {
	rctx, err := cli.ResolveContext(cmd)
	if err != nil {
		return cli.WriteError(cmd, contract.ErrCodeConfig, err.Error())
	}

	tokenAuth, err := cli.ResolveAuth(rctx.ProfileName, credentials.ModeApplication)
	if err != nil {
		return cli.WriteError(cmd, contract.ErrCodeAuth, err.Error())
	}

	transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
	client := appapi.NewClient(transport, rctx.AccountID)

	profile, err := client.GetProfile(context.Background())
	if err != nil {
		return cli.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(profile)
	return contract.Write(cmd.OutOrStdout(), resp, cli.Pretty())
}

func runProfileUpdate(cmd *cobra.Command, args []string) error {
	// Validate at least one flag is set
	nameChanged := cmd.Flags().Changed("name")
	emailChanged := cmd.Flags().Changed("email")
	availChanged := cmd.Flags().Changed("availability")
	if !nameChanged && !emailChanged && !availChanged {
		return fmt.Errorf("requires at least one of --name, --email, or --availability")
	}

	rctx, err := cli.ResolveContext(cmd)
	if err != nil {
		return cli.WriteError(cmd, contract.ErrCodeConfig, err.Error())
	}

	tokenAuth, err := cli.ResolveAuth(rctx.ProfileName, credentials.ModeApplication)
	if err != nil {
		return cli.WriteError(cmd, contract.ErrCodeAuth, err.Error())
	}

	transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
	client := appapi.NewClient(transport, rctx.AccountID)

	opts := appapi.UpdateProfileOpts{}
	if nameChanged {
		v, _ := cmd.Flags().GetString("name")
		opts.Name = &v
	}
	if emailChanged {
		v, _ := cmd.Flags().GetString("email")
		opts.Email = &v
	}
	if availChanged {
		v, _ := cmd.Flags().GetString("availability")
		opts.Availability = &v
	}

	profile, err := client.UpdateProfile(context.Background(), opts)
	if err != nil {
		return cli.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(profile)
	return contract.Write(cmd.OutOrStdout(), resp, cli.Pretty())
}
```

- [ ] **Step 5: Register application group in root.go**

Add to `init()` in `internal/cli/root.go`:

```go
rootCmd.AddCommand(cliapp.Cmd)
```

Add import: `cliapp "github.com/chatwoot/chatwoot-cli/internal/cli/application"`

- [ ] **Step 6: Run test**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/ -v -run TestProfileGet`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add internal/cli/application/ internal/cli/root.go
git commit -m "feat(cli): add application profile get and update commands"
```

---

## Task 8: Profile Update Test

**Files:**
- Modify: `internal/cli/application/profile_test.go`

- [ ] **Step 1: Write test for profile update**

Add to `profile_test.go`:

```go
func TestProfileUpdate(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{
			"id":    1,
			"name":  "New Name",
			"email": "agent@test.com",
		})
	}))
	defer srv.Close()

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "chatwoot-cli", "config.yaml")
	os.MkdirAll(filepath.Dir(cfgPath), 0755)
	os.WriteFile(cfgPath, []byte(`
default_profile: work
profiles:
  work:
    base_url: `+srv.URL+`
    account_id: 1
`), 0644)

	keyring.MockInit()
	ks := credentials.NewKeychainStore()
	_ = ks.Set("work", credentials.ModeApplication, "sk-test-token")
	t.Setenv("XDG_CONFIG_HOME", dir)

	var stdout bytes.Buffer
	Cmd.SetOut(&stdout)
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"application", "profile", "update", "--name", "New Name"})
	err := Cmd.Root().Execute()
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if gotMethod != "PATCH" {
		t.Errorf("method = %q, want PATCH", gotMethod)
	}
	if gotBody["name"] != "New Name" {
		t.Errorf("body name = %v, want New Name", gotBody["name"])
	}
	// email should not be in body (not set)
	if _, exists := gotBody["email"]; exists {
		t.Error("email should not be in body when not set")
	}
}

func TestProfileUpdateRequiresFlag(t *testing.T) {
	Cmd.SetOut(&bytes.Buffer{})
	Cmd.SetErr(&bytes.Buffer{})
	Cmd.Root().SetArgs([]string{"application", "profile", "update"})
	err := Cmd.Root().Execute()
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}
}
```

- [ ] **Step 2: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/ -v`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add internal/cli/application/profile_test.go
git commit -m "test(cli): add profile update command tests"
```

---

## Task 9: Full Test Suite and Exit Criteria Verification

**Files:**
- No new files

- [ ] **Step 1: Run full test suite**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./... -v 2>&1 | tail -40`
Expected: All tests PASS

- [ ] **Step 2: Run vet and build**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go vet ./... && go build ./cmd/chatwoot/`
Expected: No errors

- [ ] **Step 3: Verify command tree**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot --help`
Expected: Shows `auth`, `application`, and `version` commands

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot auth --help`
Expected: Shows `set`, `status`, `clear` subcommands

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application profile --help`
Expected: Shows `get` and `update` subcommands

- [ ] **Step 4: Verify version still works**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot version`
Expected: JSON envelope with version info

- [ ] **Step 5: Verify global flags exist**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot --help | grep -E "profile|base-url|account-id"`
Expected: All three flags listed

- [ ] **Step 6: Verify exit code 2 on usage error**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot auth set; echo $?`
Expected: Exit code 2 (missing required flags)

- [ ] **Step 7: Verify no business logic in CLI handlers**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && wc -l internal/cli/auth/*.go internal/cli/application/*.go`
Expected: Each handler file is small (under 100 lines), logic lives in internal packages

---

## Sprint B Exit Criteria Checklist

- [ ] `go test ./...` passes with all tests green
- [ ] `go vet ./...` clean
- [ ] `go build ./cmd/chatwoot/` succeeds
- [ ] `chatwoot auth set` stores credentials (keychain with file fallback)
- [ ] `chatwoot auth status` reports credential state without revealing tokens
- [ ] `chatwoot auth clear` removes credentials
- [ ] `chatwoot application profile get` exercises full pipeline end-to-end
- [ ] `chatwoot application profile update` sends PATCH with only provided fields
- [ ] `--profile` flag selects correct config and credentials
- [ ] `--base-url` and `--account-id` flags override profile values
- [ ] Foundation commands produce valid JSON envelopes on stdout
- [ ] Usage errors return exit code 2
- [ ] API/config/auth errors return exit code 1 with error envelope
- [ ] `--verbose` sets slog to Debug level
- [ ] No business logic in command handler files (thin handlers only)

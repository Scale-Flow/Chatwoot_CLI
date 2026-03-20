# Sprint B: CLI Shell + Foundation Commands (Layers 5 + 6)

## Scope

Wire the Sprint A plumbing (config, credentials, auth, transport, API clients)
into a Cobra-based command system and ship the first foundation commands. This
sprint proves the full pipeline end-to-end: command invocation → context
resolution → auth → transport → API client → JSON envelope on stdout.

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Context resolution | Per-command helper, not PersistentPreRun | Commands like `version` don't need config/auth; avoids special-case exclusions |
| `--verbose` behavior | Wire slog level to Debug, no transport logging yet | Transport logging is a separate concern; Sprint B focuses on command wiring |
| Token input | `--token` flag only, no stdin | Sufficient for agents and scripts; stdin adds TTY complexity |
| Profile command placement | `chatwoot application profile get/update` | Consistent with API family grouping principle; no top-level exceptions |
| Auth command placement | Top-level `chatwoot auth set/status/clear` | Auth is a CLI-local concern, not an API family |

## Command Tree

```
chatwoot
├── version                          # existing
├── auth
│   ├── set                          # store credential
│   ├── status                       # inspect credential state
│   └── clear                        # remove credential
└── application
    └── profile
        ├── get                      # GET /api/v1/profile
        └── update                   # PATCH /api/v1/profile
```

## Package Structure

### New Files

```
internal/
  cli/
    root.go                          # MODIFY: add --profile, --base-url, --account-id flags; wire slog level
    context.go                       # NEW: RuntimeContext struct, ResolveContext, ResolveAuth helpers
    context_test.go                  # NEW: context resolution tests
    auth/
      auth.go                        # NEW: auth group command definition
      set.go                         # NEW: auth set handler
      set_test.go                    # NEW
      status.go                      # NEW: auth status handler
      status_test.go                 # NEW
      clear.go                       # NEW: auth clear handler
      clear_test.go                  # NEW
    application/
      application.go                 # NEW: application group command definition
      profile.go                     # NEW: profile get/update handlers
      profile_test.go                # NEW
  chatwoot/application/
    profile.go                       # MODIFY: add UpdateProfile method
```

### Unchanged Files

```
internal/cli/version.go
internal/cli/version_test.go
cmd/chatwoot/main.go
```

## Execution Pipeline

Every authenticated command follows this pipeline:

```
Command handler
  → ResolveContext(cmd)
    → read --profile / --base-url / --account-id flags
    → load config from ~/.config/chatwoot-cli/config.yaml
    → config.ResolveProfile(profileFlag)
    → config.ResolveOverrides(profile, baseURLFlag, accountIDFlag)
    → validate BaseURL is non-empty
    → return RuntimeContext
  → ResolveAuth(profileName, mode)
    → build credentials.Resolver(EnvStore, KeychainStore, FileStore)
    → resolve token for the command's API family
    → return auth.TokenAuth
  → construct chatwoot.Client(baseURL, token, headerName)
  → call API client method
  → wrap result in contract.Success or contract.Err
  → contract.Write to stdout
```

Auth commands use a lighter path: they need profile resolution but not a full
API client or transport.

## Root Command Changes

### New Global Flags

| Flag | Type | Default | Env fallback | Description |
|------|------|---------|-------------|-------------|
| `--profile` | string | "" | `CHATWOOT_PROFILE` | Select named profile |
| `--base-url` | string | "" | `CHATWOOT_BASE_URL` | Override base URL |
| `--account-id` | int | 0 | `CHATWOOT_ACCOUNT_ID` | Override account ID |

Env fallbacks for these flags are handled by `ResolveContext`, not
`PersistentPreRun`. The existing `--pretty` and `--verbose` env handling stays
in `PersistentPreRun` since those are rendering concerns, not context
resolution.

### Verbose Wiring

`PersistentPreRun` is extended (after the existing `--pretty`/`--verbose` env
var checks) to call `slog.SetDefault` with a `slog.NewTextHandler(os.Stderr,
&slog.HandlerOptions{Level: slog.LevelDebug})` when `verboseFlag` is true. No
transport logging is added in this sprint — just the level switch so it's
ready when transport logging lands.

### Command Registration

```go
// in root.go init()
rootCmd.AddCommand(authCmd)         // from cli/auth/
rootCmd.AddCommand(applicationCmd)  // from cli/application/
```

## Runtime Context Resolution (`context.go`)

### RuntimeContext

```go
type RuntimeContext struct {
    ProfileName string
    BaseURL     string
    AccountID   int
}
```

### ResolveContext(cmd *cobra.Command) (*RuntimeContext, error)

1. Read flag values from `cmd.Flags()` for `--profile`, `--base-url`,
   `--account-id`
2. Load config from `os.UserConfigDir()` + `/chatwoot-cli/config.yaml`
   - File missing: proceed with empty `Config{}` (flags/env can provide
     everything)
   - File malformed: return `config_error`
3. Call `config.ResolveProfile(profileFlag)`:
   - If config file was loaded: resolve normally (flag > env > default_profile
     > "default")
   - If config was empty AND `ResolveProfile` returns "profile not found":
     proceed with an empty `Profile{}` — the user may be providing everything
     via flags/env. The profile name defaults to "default".
4. `config.ResolveOverrides(profile, baseURLFlag, accountIDFlag)` — flag > env
   > profile values
5. Validate: BaseURL empty after all resolution → return `config_error` ("no
   base URL configured — set base_url in profile or use --base-url flag")
6. Return `RuntimeContext`

AccountID is not validated here. It is only required for application API
commands that include it in URL paths. Those commands check it themselves.

This design ensures a zero-config setup works when flags or env vars supply
the necessary values (e.g., `chatwoot auth set --base-url https://app.chatwoot.com --token sk-xxxx`).

### ResolveAuth(profileName string, mode credentials.AuthMode) (auth.TokenAuth, error)

1. Build `credentials.Resolver` with `EnvStore`, `KeychainStore`, `FileStore`
   - FileStore path: `os.UserConfigDir()` + `/chatwoot-cli/credentials.yaml`
2. Call `auth.ResolveApplication` or `auth.ResolvePlatform` depending on mode
3. Return `auth.TokenAuth` or error with `auth_error` code

### Flag Accessor Functions

`context.go` exports functions for subcommand packages to read flag values:

```go
func ProfileFlag(cmd *cobra.Command) string
func BaseURLFlag(cmd *cobra.Command) string
func AccountIDFlag(cmd *cobra.Command) int
```

The existing `cli.Pretty()` function in `root.go` continues to serve as the
accessor for the pretty flag. No duplicate accessor is needed.

## Auth Commands

### `chatwoot auth set`

```
chatwoot auth set --mode application --token sk-xxxx
chatwoot auth set --mode platform --token pk-xxxx
chatwoot auth set --profile selfhosted --mode application --token sk-xxxx
```

**Flags:**

| Flag | Type | Required | Values |
|------|------|----------|--------|
| `--mode` | string | yes | `application`, `platform` |
| `--token` | string | yes | token value |

**Behavior:**

1. Resolve profile name only — call `config.ResolveProfile(profileFlag)`
   directly, not `ResolveContext` (which validates BaseURL and would fail when
   no config exists yet). If config file is absent, default to profile name
   "default".
2. Try `KeychainStore.Set(profile, mode, token)`
3. If keychain fails: warn on stderr ("keychain unavailable, using file-based
   credential storage"), fall back to `FileStore.Set()` with 0600
4. Output success envelope

**Output:**

```json
{
  "ok": true,
  "data": {
    "profile": "work",
    "mode": "application",
    "source": "keychain"
  }
}
```

### `chatwoot auth status`

```
chatwoot auth status
chatwoot auth status --profile selfhosted
```

No extra flags beyond global `--profile`.

**Behavior:**

1. Resolve full context (profile, BaseURL, AccountID) for display. If config
   file is missing, show empty/zero values rather than erroring.
2. For each mode (application, platform): probe `Resolver.Get()` to check
   existence and source
3. Report without revealing tokens

**Output:**

```json
{
  "ok": true,
  "data": {
    "profile": "work",
    "base_url": "https://app.chatwoot.com",
    "account_id": 1,
    "credentials": {
      "application": {"status": "configured", "source": "keychain"},
      "platform": {"status": "not_configured"}
    }
  }
}
```

Note: This per-mode object format (with `status` and `source` per credential)
supersedes the flatter format in the parent design decisions spec. The richer
format lets callers know where each credential came from without a separate
top-level `source` field that would be ambiguous with multiple credential
types.

### `chatwoot auth clear`

```
chatwoot auth clear --mode application
chatwoot auth clear --all
chatwoot auth clear --profile selfhosted --all
```

**Flags:**

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--mode` | string | no | `application` or `platform` |
| `--all` | bool | no | Clear all modes |

One of `--mode` or `--all` is required. The command's `RunE` validates this
manually and returns a usage error if neither is provided (Cobra does not have
a built-in at-least-one constraint for flags).

**Behavior:**

1. Resolve profile name (same lightweight path as `auth set`)
2. Delete from keychain; if keychain unavailable, delete from file store
3. If the credential doesn't exist, succeed silently (idempotent delete — same
   behavior as the existing `KeychainStore.Delete` which returns nil on
   not-found)
4. Output success envelope confirming what was cleared

**Output:**

```json
{
  "ok": true,
  "data": {
    "profile": "work",
    "cleared": ["application"]
  }
}
```

## Application Profile Commands

### `chatwoot application profile get`

```
chatwoot application profile get
```

No extra flags. Uses resolved profile context.

**Behavior:**

1. `ResolveContext(cmd)` — needs BaseURL (AccountID not required for this
   endpoint)
2. `ResolveAuth(profileName, ModeApplication)` — needs application token
3. Build `chatwoot.Client`, build `application.Client`
4. Call `client.GetProfile(ctx)` — `GET /api/v1/profile`
5. Wrap in `contract.Success()`, write to stdout

This is the first command exercising the full pipeline end-to-end.

### `chatwoot application profile update`

```
chatwoot application profile update --name "New Name"
chatwoot application profile update --email "new@example.com"
chatwoot application profile update --availability online
```

**Flags:**

| Flag | Type | Required | Values |
|------|------|----------|--------|
| `--name` | string | no | display name |
| `--email` | string | no | email address |
| `--availability` | string | no | `online`, `offline`, `busy` |

At least one flag must be provided. The command's `RunE` validates this
manually and returns a usage error if none are set. Only explicitly set flags
are included in the request body.

**Behavior:**

1. Same context/auth resolution as `get`
2. Build JSON body from provided flags
3. Call `client.UpdateProfile(ctx, opts)` — `PATCH /api/v1/profile`
4. Return updated profile in success envelope

### New API Client Method

Add to `internal/chatwoot/application/profile.go`:

```go
type UpdateProfileOpts struct {
    Name         *string `json:"name,omitempty"`
    Email        *string `json:"email,omitempty"`
    Availability *string `json:"availability,omitempty"`
}

func (c *Client) UpdateProfile(ctx context.Context, opts UpdateProfileOpts) (*Profile, error)
```

Uses pointer fields so only explicitly set values are serialized.

## Error Handling

### Exit Codes

| Scenario | Exit code |
|----------|-----------|
| Success | 0 |
| API error (401, 404, etc.) | 1 |
| Config error (missing BaseURL, malformed config) | 1 |
| Auth error (no credentials) | 1 |
| Usage error (bad flags, missing required flags) | 2 |

All exit-1 errors produce a JSON error envelope on stdout with the appropriate
error code (`config_error`, `auth_error`, `not_found`, etc.).

### Error Envelope Examples

Missing credentials:

```json
{
  "ok": false,
  "error": {
    "code": "auth_error",
    "message": "no application credentials for profile \"work\": credential not found"
  }
}
```

No base URL configured:

```json
{
  "ok": false,
  "error": {
    "code": "config_error",
    "message": "no base URL configured — set base_url in profile or use --base-url flag"
  }
}
```

## Testing Strategy

### Unit Tests: Context Resolution (`context_test.go`)

- Config file exists → resolves profile correctly
- Config file missing → proceeds with empty config, flags/env still work
- No BaseURL resolvable → returns config_error
- Flag > env > config precedence verified
- AccountID not validated at context level

### Command Handler Tests (per command)

- Each command test uses `httptest.NewServer` for API calls (profile
  get/update)
- Auth commands test against mock credential stores using
  `keyring.MockInit()` — no real keychain
- Tests execute Cobra commands programmatically: set args, capture
  stdout/stderr, assert JSON envelope shape
- Pattern: `rootCmd.SetArgs([]string{"auth", "status"})` with captured output

### Integration-Style Tests (full pipeline)

- Complete path: command → context resolution → auth → transport → mock server
  → JSON output
- Verify exit codes: 0 success, 1 API error, 2 usage error
- Verify `--pretty` indents output
- Verify error envelopes have correct `code` field

### Not Tested in Sprint B

- Real keychain integration (use `keyring.MockInit()`)
- Real Chatwoot API calls
- Transport logging / `--verbose` output content
- Stdin token input

## Exit Criteria

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

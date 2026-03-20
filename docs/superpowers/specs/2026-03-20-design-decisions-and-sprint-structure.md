# Design Spec: Design Decisions and Sprint Structure

## Scope

This spec settles the six open design decisions from the roadmap and defines how
the remaining layers are grouped into implementation sprints. It covers
decisions that affect the entire CLI, from config and credentials through
transport, API clients, and command surface.

## Decisions Settled

| # | Decision | Choice | Rationale |
|---|----------|--------|-----------|
| 1 | Output envelope contract | Already frozen in Layer 1 | Boolean-discriminated `ok` field, contract tests lock it down |
| 2 | Account context strategy | Config default with `--account-id` flag override | AWS `--region` pattern; profile stores default, flag overrides |
| 3 | Reports namespace placement | Stay under `application` | Consistent with API family grouping principle |
| 4 | Message command placement | Dedicated `application messages` group with `--conversation-id` flag | Flat for agent ergonomics; parent relationship explicit via flag |
| 5 | Runtime context and credential model | Single config file, named profiles, keychain + 0600 file fallback, env var escape hatch | See detailed design below |
| 6 | Rate limiting and pagination | Single-page default + `--all` auto-paginate; retry with exponential backoff on 429 | Safe default, agent opt-in for full fetches |

## Decision 2: Account Context Strategy

Application API endpoints require `account_id` in the URL path. The CLI uses a
config-default-with-override pattern:

- Each profile stores an optional `account_id`
- Application commands use the profile's `account_id` automatically
- `--account-id` flag overrides the profile default
- `CHATWOOT_ACCOUNT_ID` env var overrides the profile default
- If no account ID is resolvable, commands return `config_error`

Precedence: `--account-id` flag > `CHATWOOT_ACCOUNT_ID` env > profile config > error.

## Decision 3: Reports Namespace Placement

Reports stay under `chatwoot application reports`. This preserves the principle
that top-level command groups map to API families. Reports span `/api/v1` and
`/api/v2` but both are within the application API family.

## Decision 4: Message Command Placement

Messages get a dedicated command group at `chatwoot application messages` rather
than nesting under `chatwoot application conversations messages`. The
`--conversation-id` flag makes the parent relationship explicit.

This keeps the command tree flat for the primary AI agent use case where message
operations are frequent. The API-level nesting
(`/conversations/{id}/messages`) is handled internally by the client package.

## Decision 5: Runtime Context and Credential Model

### Config File

**Location:** `~/.config/chatwoot-cli/config.yaml` (XDG-compliant)

**Schema:**

```yaml
default_profile: work

profiles:
  work:
    base_url: https://app.chatwoot.com
    account_id: 1
  selfhosted:
    base_url: https://chatwoot.internal.corp
    account_id: 42
```

**Profile fields:**

- `base_url` — Chatwoot instance URL (required per profile)
- `account_id` — default account ID for application commands (optional)

No auth mode field in config. The CLI infers which auth mode to use from which
command family is invoked (`application` commands use the application token,
`platform` commands use the platform token).

### Profile Resolution

Precedence for selecting the active profile:

1. `--profile` flag
2. `CHATWOOT_PROFILE` env var
3. `default_profile` in config.yaml
4. Falls back to profile named `default`

### Global Flags

New global flags added to the root command (beyond existing `--pretty` and
`--verbose`):

| Flag | Type | Default | Env var | Description |
|------|------|---------|---------|-------------|
| `--profile` | string | (see resolution) | `CHATWOOT_PROFILE` | Select named profile |
| `--base-url` | string | (from profile) | `CHATWOOT_BASE_URL` | Override base URL |
| `--account-id` | int | (from profile) | `CHATWOOT_ACCOUNT_ID` | Override account ID |

### Credential Storage

**One profile = one server, multiple auth slots.** A profile can hold
credentials for both application and platform API families. The CLI picks the
right credential based on which command family is invoked.

Client API credentials (`inbox_identifier` + `contact_identifier`) are
session-scoped and passed via flags or env vars per command, not stored in
profiles.

**Storage backends (resolution order):**

1. **Env vars** (if set, skip everything else)
2. **OS keychain** via `go-keyring`
3. **Credential file** at `~/.config/chatwoot-cli/credentials.yaml` with `0600` permissions
4. **Error** with `auth_error` code

**Keychain storage structure:**

Each credential is stored via `go-keyring` with:

- Service: `chatwoot-cli`
- User: `<profile>/<auth-type>`

Examples for a profile named `work`:

- Application token: `Set("chatwoot-cli", "work/application", "<token>")`
- Platform token: `Set("chatwoot-cli", "work/platform", "<token>")`

**Env var escape hatch (CI/automation):**

| Env var | Overrides |
|---------|-----------|
| `CHATWOOT_ACCESS_TOKEN` | Application token (skips keychain) |
| `CHATWOOT_PLATFORM_TOKEN` | Platform token (skips keychain) |
| `CHATWOOT_BASE_URL` | Base URL (skips profile) |
| `CHATWOOT_ACCOUNT_ID` | Account ID (skips profile) |

When any token env var is set, it bypasses the entire profile/keychain system.
A CI pipeline can work with zero config files.

**Credential file (headless fallback):**

When the OS keychain is unavailable (headless Linux, Docker, WSL without
desktop session), credentials can be stored in a permission-protected file:

```yaml
profiles:
  work:
    application_token: sk-xxxx
    platform_token: pk-xxxx
```

Safety measures:

- Created with `0600` (`-rw-------`) permissions, owner-only read/write
- `auth set` in headless mode writes to this file and sets `0600`
- On startup, if the file exists with permissions wider than `0600`, warn on stderr and refuse to read it (same approach as SSH with `~/.ssh/id_rsa`)
- Only created in the XDG config path, never inside a project directory

**Keychain unavailable behavior:**

When `go-keyring` fails (no D-Bus, no keychain daemon):

- Warn once on stderr: `"warning: keychain unavailable, using file-based credential storage"`
- Fall back to credential file (env vars still take highest priority if set)
- `auth set` writes to the credential file with `0600` permissions
- `auth status` reports `"source": "file"` or `"source": "environment"`

**Redaction policy:** Tokens are never logged, even with `--verbose`. Verbose
output shows `"Authorization: <redacted>"` for request headers.

### Auth Commands

**`chatwoot auth set`**

Stores a credential in the keychain (or credential file in headless mode).

```bash
chatwoot auth set --mode application --token sk-xxxx
chatwoot auth set --mode platform --token pk-xxxx
chatwoot auth set --profile selfhosted --mode application --token sk-xxxx
```

Output: `{"ok": true, "data": {"profile": "work", "mode": "application", "status": "stored"}}`

**`chatwoot auth status`**

Shows which credentials exist for the active profile without revealing tokens.

```bash
chatwoot auth status
chatwoot auth status --profile selfhosted
```

Output:

```json
{
  "ok": true,
  "data": {
    "profile": "work",
    "base_url": "https://app.chatwoot.com",
    "account_id": 1,
    "credentials": {
      "application": "configured",
      "platform": "not_configured"
    },
    "source": "keychain"
  }
}
```

**`chatwoot auth clear`**

Removes credentials from keychain (or credential file) for the active profile.

```bash
chatwoot auth clear --mode application
chatwoot auth clear --all
chatwoot auth clear --profile selfhosted --all
```

Output: success envelope confirming what was cleared.

## Decision 6: Rate Limiting and Pagination

### Pagination

List commands return a single page by default. The caller controls pagination
with `--page` and `--per-page` flags. Pagination metadata is always present in
the response envelope's `meta.pagination` field.

The `--all` flag opts into auto-pagination: the CLI fetches all pages
sequentially and returns a single merged collection. Rate limiting and retries
apply per-page.

### Rate Limit Handling

HTTP 429 responses trigger automatic retry with exponential backoff:

- Maximum 3 retry attempts
- Backoff schedule: 1s, 2s, 4s (with jitter)
- If `Retry-After` header is present, respect it instead of calculated backoff
- If retries are exhausted, return `rate_limited` error envelope
- `--verbose` logs retry attempts and wait times to stderr

HTTP 5xx responses follow the same retry policy. All other errors fail
immediately.

## Transport Layer

### HTTP Client Wrapper (`internal/chatwoot/client.go`)

A shared `Client` struct used by all API family packages:

- Builds requests with correct base path per API family (`/api/v1`,
  `/platform/api/v1`, `/public/api/v1`)
- Injects auth headers (`api_access_token` for application/platform, path
  params for client API)
- Retries with exponential backoff on 429 and 5xx (3 attempts, 1s → 2s → 4s
  with jitter)
- Logs request/response metadata to stderr via `slog` when `--verbose` is
  active
- Redacts `Authorization` header values in all logs
- Normalizes HTTP errors to contract error codes
- Default timeout: 30s (configurable via config)

### Error Mapping

| HTTP Status | Contract Error Code |
|-------------|-------------------|
| 401 | `unauthorized` |
| 403 | `forbidden` |
| 404 | `not_found` |
| 422 | `validation_error` |
| 429 | `rate_limited` (after retries exhausted) |
| 5xx | `server_error` (after retries exhausted) |
| Connection/timeout | `network_error` |

### Pagination Helper

A `ListAll` function that auto-paginates when `--all` is used. The caller
passes a page-fetcher function; the helper handles the loop, collecting results
into a single merged collection. Respects rate limiting (retries apply
per-page).

### Testing Interface

The transport exposes a `Doer` interface:

```go
type Doer interface {
    Do(*http.Request) (*http.Response, error)
}
```

API client tests use `httptest.NewServer` without touching retry/auth logic.

## API Family Clients

### Design Pattern

Three typed client packages, one per API family. Each follows the same pattern:

- Wraps the transport for its base path
- Methods map to API resources
- Builds requests, calls transport, decodes responses, returns typed data or error
- No Cobra dependency
- No JSON envelope wrapping (that is the CLI layer's responsibility)

### `internal/chatwoot/application/`

Application API client for `/api/v1` and `/api/v2`:

```go
type Client struct { /* transport, accountID */ }

func (c *Client) ListConversations(ctx context.Context, opts ListConversationsOpts) ([]Conversation, *contract.Pagination, error)
func (c *Client) GetConversation(ctx context.Context, id int) (*Conversation, error)
func (c *Client) CreateMessage(ctx context.Context, conversationID int, opts CreateMessageOpts) (*Message, error)
```

### `internal/chatwoot/platform/`

Platform API client for `/platform/api/v1`:

```go
type Client struct { /* transport */ }

func (c *Client) CreateAccount(ctx context.Context, opts CreateAccountOpts) (*Account, error)
func (c *Client) GetUser(ctx context.Context, id int) (*User, error)
```

### `internal/chatwoot/clientapi/`

Client (public) API client for `/public/api/v1`. Auth is identifier-based:

```go
type Client struct { /* transport, inboxIdentifier */ }

func (c *Client) CreateContact(ctx context.Context, opts CreateContactOpts) (*Contact, error)
func (c *Client) ListMessages(ctx context.Context, conversationID int) ([]Message, error)
```

### Shared Model Types

Common types (pagination, identifiers, timestamps) live in
`internal/chatwoot/` at the family root. Resource-specific types (Conversation,
Contact, Message) live in their respective family package.

### Testing Approach

Each client package has tests using `httptest.NewServer` that verify:

- Correct URL path construction (including account ID interpolation)
- Correct HTTP method and headers
- Request body serialization
- Response deserialization
- Error mapping for non-2xx responses

## Sprint Structure

### Sprint A: Plumbing (Layers 2 + 3 + 4)

The full internal stack from config through API clients. Nothing user-facing,
fully testable.

**Scope:**

- Config loading, profile resolution, Viper integration
- Credential store (keychain + `0600` file fallback + env vars)
- HTTP transport with retry, backoff, auth injection, error mapping
- Pagination helper with `--all` support
- Application, Platform, and Client API client packages
- Shared model types
- Full test coverage at every layer

**Exit criteria:**

- `go test ./...` passes
- All internal packages have test coverage
- No Cobra dependency in client code
- Config resolution follows documented precedence chain
- Credentials round-trip through keychain and file fallback
- Transport retries and error mapping are tested with `httptest`

### Sprint B: CLI Shell + Foundation Commands (Layers 5 + 6)

Wire the plumbing into Cobra and prove it end-to-end.

**Scope:**

- Root command with new global flags (`--profile`, `--base-url`, `--account-id`)
- Command registration pattern and shared execution pipeline (resolve context →
  call client → render envelope)
- `chatwoot auth set`, `auth status`, `auth clear`
- `chatwoot profile get`, `profile update`
- Integration tests for the full path from command invocation through transport
  to JSON output

**Exit criteria:**

- Foundation commands produce valid JSON envelopes
- Auth round-trips through keychain
- `--verbose` shows request traces on stderr
- `--profile` flag selects correct config and credentials
- Usage errors return exit code 2
- API errors return exit code 1 with error envelope

### Sprint C: Core Application Commands (Layer 7)

First real user value.

**Scope:**

- `chatwoot application contacts *`
- `chatwoot application conversations *`
- `chatwoot application messages *`
- `chatwoot application inboxes list|get|members *`

**Exit criteria:**

- Core support workflows from the use-case doc are executable end-to-end
- Read and mutation paths behave consistently
- `--all` pagination works across list commands
- Error envelopes are consistent across all commands

### Sprint D onward: Layers 8-11

Individually as defined in the roadmap:

- Sprint D: Layer 8 (Operational Admin Surface)
- Sprint E: Layer 9 (Platform, Client, and Survey Surfaces)
- Sprint F: Layer 10 (Experimental and Source-Derived Surfaces)
- Sprint G: Layer 11 (Packaging, Release, and Maintenance)

## Out of Scope

- Interactive prompts (beyond potential `auth set` token input in future)
- OAuth/SSO login flows
- Credential helper command (`credential_command` plugin) — may be added later
- `--output` global flag (deferred; JSON is the only output format for now)
- `auth set --mode client` (client API credentials are session-scoped flags, not stored; the command structure doc's `--mode client` example is superseded by this decision)
- TUI or GUI
- Daemon mode

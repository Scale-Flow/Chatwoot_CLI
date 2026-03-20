# Design Spec: Layer 0 (Repository Foundation) + Layer 1 (Shared Contracts)

## Scope

This spec covers the first two layers of the Chatwoot CLI roadmap:

- **Layer 0:** Turn the repo from documentation-only into a buildable Go CLI project
- **Layer 1:** Freeze the machine-facing JSON envelope contract and rendering pipeline

Everything downstream (config, transport, API clients, commands) depends on these layers stabilizing first.

## Decisions Made

| Decision | Choice | Rationale |
|---|---|---|
| Go version | 1.26.1 | Latest stable; pinned in `go.mod` |
| Task runner | Taskfile (go-task) | Cross-platform, YAML-based, zero binary dependencies beyond `task` |
| Linting | `golangci-lint` minimal config | vet, gofmt, errcheck, staticcheck — value without noise |
| Version injection | `ldflags` + `debug/buildinfo` fallback | Clean release versions; `go install` still reports something useful |
| Envelope style | Boolean-discriminated (`ok` field) | Single discriminator, trivially parseable by AI agents |

## Layer 0: Repository Foundation

### Go Module

- Module path: `github.com/chatwoot/chatwoot-cli`
- Go version: 1.26.1 pinned in `go.mod`
- Initial dependencies: `cobra`, `viper`, `go-keyring`

Note: the module path `github.com/chatwoot/chatwoot-cli` is used consistently throughout this spec, including in ldflags references. If a different path is chosen at init time, all ldflags paths must update to match.

### Directory Structure

```
cmd/chatwoot/main.go             # entry point, wires root command
internal/
  cli/
    root.go                      # Cobra root command, global flags
    version.go                   # version command handler
  contract/
    contract.go                  # envelope types
    render.go                    # renderer functions
    contract_test.go             # contract-locking tests
  config/                        # (placeholder) Viper config
  credentials/                   # (placeholder) keychain integration
  auth/                          # (placeholder) auth mode modeling
  chatwoot/
    client.go                    # (placeholder) HTTP transport
    application/                 # (placeholder)
    platform/                    # (placeholder)
    clientapi/                   # (placeholder)
    reports/                     # (placeholder)
  testutil/                      # (placeholder) test helpers
  version/
    version.go                   # version variables + buildinfo fallback
```

Placeholder packages contain a single `.go` file named `<package>.go` (e.g., `config/config.go`, `credentials/credentials.go`, `auth/auth.go`) with only the `package` declaration. No dead code — just the package line.

### Taskfile

`Taskfile.yml` with these targets:

| Target | Command |
|---|---|
| `build` | `go build -ldflags "..." ./cmd/chatwoot/` |
| `test` | `go test ./...` |
| `lint` | `golangci-lint run ./...` |
| `fmt` | `gofmt -w .` |
| `vet` | `go vet ./...` |

The `build` target injects version, commit, and date via `-ldflags`:

```
-X github.com/chatwoot/chatwoot-cli/internal/version.Version={{.GIT_TAG}}
-X github.com/chatwoot/chatwoot-cli/internal/version.Commit={{.GIT_COMMIT}}
-X github.com/chatwoot/chatwoot-cli/internal/version.Date={{.BUILD_DATE}}
```

Taskfile dynamic variables:

```yaml
vars:
  GIT_TAG:
    sh: git describe --tags --always --dirty 2>/dev/null || echo "dev"
  GIT_COMMIT:
    sh: git rev-parse --short HEAD
  BUILD_DATE:
    sh: date -u '+%Y-%m-%dT%H:%M:%SZ'
```

### Linting Configuration

`.golangci.yml` enables a minimal set:

```yaml
linters:
  enable:
    - errcheck
    - staticcheck
    - govet
    - gofmt
```

### Version Package (`internal/version/`)

Exports three variables:

```go
var (
    Version string  // set by ldflags, e.g. "0.1.0"
    Commit  string  // set by ldflags, e.g. "a1b2c3d"
    Date    string  // set by ldflags, e.g. "2026-03-19T12:00:00Z"
)
```

`Info() map[string]string` returns version info with these keys:

- `version` — from ldflags, or module version from `debug.ReadBuildInfo()`
- `commit` — from ldflags, or VCS revision from `debug.ReadBuildInfo()`
- `date` — from ldflags, or VCS time from `debug.ReadBuildInfo()`
- `go_version` — always from `runtime.Version()`

When `Version` is empty (no ldflags), the function falls back to `debug.ReadBuildInfo()` to extract the module version and VCS revision. If both are unavailable, values default to `"unknown"`.

### `chatwoot version` Command

The first working command. It calls `version.Info()`, wraps the result in a success envelope via the contract renderer, and writes to stdout:

```json
{
  "ok": true,
  "data": {
    "version": "0.1.0",
    "commit": "a1b2c3d",
    "date": "2026-03-19T12:00:00Z",
    "go_version": "go1.26.1"
  }
}
```

### `.gitignore`

```
/chatwoot
*.exe
/dist/
```

### Layer 0 Exit Criteria

- `task build` produces a working `chatwoot` binary
- `task test` runs and passes
- `task lint` passes with zero warnings
- `./chatwoot version` outputs a valid JSON success envelope
- All planned package directories exist with valid Go package declarations

## Layer 1: Shared Contracts and Rendering

### Envelope Schema

All stdout output from every command uses a single envelope structure.

#### Success (single resource)

```json
{
  "ok": true,
  "data": {
    "id": 42,
    "name": "Support Inbox"
  }
}
```

#### Success (collection)

```json
{
  "ok": true,
  "data": [
    {"id": 1, "name": "Alice"},
    {"id": 2, "name": "Bob"}
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "per_page": 25,
      "total_count": 142,
      "total_pages": 6
    }
  }
}
```

#### Success (empty collection)

```json
{
  "ok": true,
  "data": [],
  "meta": {
    "pagination": {
      "page": 1,
      "per_page": 25,
      "total_count": 0,
      "total_pages": 0
    }
  }
}
```

#### Success with warnings

```json
{
  "ok": true,
  "data": {"id": 42},
  "warnings": [
    {"code": "deprecated_endpoint", "message": "This endpoint will be removed in Chatwoot 5.0"}
  ]
}
```

#### Error

```json
{
  "ok": false,
  "error": {
    "code": "not_found",
    "message": "Conversation 12345 not found"
  }
}
```

#### Error with detail

```json
{
  "ok": false,
  "error": {
    "code": "validation_error",
    "message": "Invalid request parameters",
    "detail": {
      "fields": {"email": "is not a valid email address"}
    }
  }
}
```

### Envelope Invariants

These are enforced by constructor functions, validated by `Write`, and locked down by contract tests:

1. `ok` is always present and is a boolean
2. When `ok: true`: `data` is present, `error` is absent
3. When `ok: false`: `error` is present, `data` is absent
4. `data` for collections is always a JSON array (empty `[]`, never `null`)
5. `data` for single resources is always a JSON object
6. `meta` is only present when there is pagination or other metadata
7. `warnings` is only present when non-empty
8. All field names use `snake_case`

**Enforcement:** The constructor functions (`Success`, `Err`, etc.) are the only intended way to build a `Response`. They enforce mutual exclusivity of `data` and `error`. `Write` validates the envelope before serializing and returns an error if both `Data` and `Error` are non-nil, or if `OK` disagrees with which fields are set.

### Envelope Types

```go
package contract

type Response struct {
    OK       bool         `json:"ok"`
    Data     any          `json:"data,omitempty"`
    Meta     *Meta        `json:"meta,omitempty"`
    Warnings []Warning    `json:"warnings,omitempty"`
    Error    *ErrorDetail `json:"error,omitempty"`
}

type Meta struct {
    Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
    Page       int `json:"page"`
    PerPage    int `json:"per_page"`
    TotalCount int `json:"total_count"`
    TotalPages int `json:"total_pages"`
}

type ErrorDetail struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Detail  any    `json:"detail,omitempty"`
}

type Warning struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

### Renderer Functions

```go
package contract

func Success(data any) Response
func SuccessList(data any, meta Meta) Response
func Err(code string, message string) Response
func ErrWithDetail(code string, message string, detail any) Response

func Write(w io.Writer, resp Response, pretty bool) error
```

`Write` serializes the response as JSON to the given writer:
- Validates envelope invariants before serializing (returns error on violation)
- Compact by default (no indentation)
- Pretty mode (2-space indent) when `--pretty` flag is active
- Always appends a trailing newline
- Deterministic field order (follows Go struct field declaration order)

`SuccessList` normalizes nil slices to empty slices (`[]`) before storing in `Data`, preventing `"data": null` in JSON output. Callers should use `SuccessList` for collection responses and `Success` for single resources.

Command handlers call `Success(data)` or `Err(code, msg)` to build the envelope, then `Write(os.Stdout, resp, prettyFlag)` to output it. Tests pass a `bytes.Buffer` instead of `os.Stdout`.

### Error Code Conventions

Standardized codes mapped from HTTP status and common CLI failure modes:

| Code | Source |
|---|---|
| `unauthorized` | HTTP 401 |
| `forbidden` | HTTP 403 |
| `not_found` | HTTP 404 |
| `validation_error` | HTTP 422 |
| `rate_limited` | HTTP 429 |
| `server_error` | HTTP 5xx |
| `network_error` | Connection/timeout failure |
| `config_error` | Missing config, bad profile |
| `auth_error` | No credentials configured |

These codes are string constants in the `contract` package. They're used by the transport layer (Layer 3) to normalize API errors, and by command handlers for local validation failures.

### Exit Code Policy

| Exit code | Meaning |
|---|---|
| 0 | Success (`ok: true`) |
| 1 | API or runtime error (`ok: false`) |
| 2 | Usage error (bad flags, missing required args) |

Cobra handles exit code 2 natively for usage errors. Exit codes supplement the JSON contract — agents should parse the JSON, but scripts can use exit codes for quick pass/fail.

**`main.go` exit code wiring:** `cmd/chatwoot/main.go` calls `cli.Execute()`, which returns an exit code. If the command wrote an error envelope, `Execute` returns 1. If Cobra detected a usage error, it returns 2. `main.go` calls `os.Exit(code)` — it does not use `log.Fatal` or any other mechanism that would print to stderr.

### Global Flags

These flags are defined as persistent flags on the root Cobra command (available to all subcommands):

| Flag | Type | Default | Env var | Description |
|---|---|---|---|---|
| `--pretty` | bool | `false` | `CHATWOOT_PRETTY` | Indent JSON output with 2 spaces |
| `--verbose` | bool | `false` | `CHATWOOT_VERBOSE` | Enable structured diagnostic logging on stderr |

Additional global flags (`--profile`, `--base-url`, `--output`, `--account-id`) are deferred to Layer 2+.

### Stderr Policy

- Stderr is for diagnostics only, never JSON
- Default: silent on success
- `--verbose` enables structured log lines via `log/slog`
- Verbose output includes: config resolution, request timing, retry decisions
- Stderr must never contain tokens or secrets (redaction is a Layer 3 concern but the policy is established here)

### Contract Tests

`internal/contract/contract_test.go` uses table-driven tests to lock down:

1. Success envelope shape — single resource, collection, empty collection
2. Error envelope shape — with and without detail
3. Warning presence and absence
4. Pagination metadata shape
5. Pretty vs compact output
6. JSON field naming (`snake_case`)
7. Nil-safety — `data: []` not `data: null` for empty collections
8. Invariant enforcement — `ok: true` cannot have `error`, `ok: false` cannot have `data`

These tests are the frozen contract. If they break, it signals a deliberate schema change that must be reviewed.

### Layer 1 Exit Criteria

- All envelope types compile and have godoc
- Renderer functions produce correct JSON for all cases
- Contract tests pass and cover every invariant listed above
- `chatwoot version` uses the renderer and produces a valid success envelope
- Pretty output works with `--pretty` flag
- Error envelopes work with standardized error codes

## Out of Scope

These are deferred to later layers:

- Configuration and profiles (Layer 2)
- Credential storage (Layer 2)
- HTTP transport and retry (Layer 3)
- API clients (Layer 4)
- Command surface beyond `version` (Layers 5+)
- Pagination auto-fetching behavior (Layer 3)
- Rate limit handling (Layer 3)

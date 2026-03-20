# Layer 0+1: Repository Foundation & Shared Contracts Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Bootstrap the Chatwoot CLI Go project and freeze the JSON envelope contract that all future commands depend on.

**Architecture:** Layered Go CLI — `cmd/chatwoot/main.go` wires a Cobra root command, which delegates to thin command handlers that build `contract.Response` envelopes and render them to stdout as JSON. The `internal/version` package provides build metadata with ldflags injection and `debug/buildinfo` fallback.

**Tech Stack:** Go 1.26.1, Cobra, Taskfile (go-task), golangci-lint

**Spec:** `docs/superpowers/specs/2026-03-19-layer-0-1-foundation-contracts-design.md`

---

## File Map

| File | Responsibility |
|---|---|
| `go.mod` | Module definition, Go version pin, dependencies |
| `Taskfile.yml` | Build, test, lint, fmt, vet targets with ldflags injection |
| `.golangci.yml` | Minimal linter config (errcheck, staticcheck, govet, gofmt) |
| `.gitignore` | Ignore built binaries and dist/ |
| `cmd/chatwoot/main.go` | Entry point — wires root command, handles exit codes |
| `internal/version/version.go` | Version variables, `Info()` with ldflags + buildinfo fallback |
| `internal/version/version_test.go` | Tests for `Info()` fallback behavior |
| `internal/contract/contract.go` | Envelope types: `Response`, `Meta`, `Pagination`, `ErrorDetail`, `Warning` |
| `internal/contract/errors.go` | Error code string constants |
| `internal/contract/render.go` | Constructor functions (`Success`, `SuccessList`, `Err`, `ErrWithDetail`) and `Write` |
| `internal/contract/contract_test.go` | Contract-locking tests for all envelope shapes and invariants |
| `internal/cli/root.go` | Cobra root command, `--pretty` and `--verbose` persistent flags, `Execute()` |
| `internal/cli/version.go` | Version command handler |
| `internal/cli/version_test.go` | Test that version command outputs valid envelope |
| `internal/config/config.go` | Placeholder package declaration |
| `internal/credentials/credentials.go` | Placeholder package declaration |
| `internal/auth/auth.go` | Placeholder package declaration |
| `internal/chatwoot/client.go` | Placeholder package declaration |
| `internal/chatwoot/application/application.go` | Placeholder package declaration |
| `internal/chatwoot/platform/platform.go` | Placeholder package declaration |
| `internal/chatwoot/clientapi/clientapi.go` | Placeholder package declaration |
| `internal/chatwoot/reports/reports.go` | Placeholder package declaration |
| `internal/testutil/testutil.go` | Placeholder package declaration |

---

### Task 1: Initialize Go Module and Project Scaffolding

**Files:**
- Create: `go.mod`
- Create: `.gitignore`
- Create: `cmd/chatwoot/main.go` (minimal — just `package main` and empty `func main()`)
- Create: all placeholder packages listed in File Map

- [ ] **Step 1: Initialize Go module**

Run:
```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI
go mod init github.com/chatwoot/chatwoot-cli
```

Then edit `go.mod` to set Go version to 1.26.1.

- [ ] **Step 2: Create .gitignore**

```
/chatwoot
*.exe
/dist/
```

- [ ] **Step 3: Create minimal main.go**

Create `cmd/chatwoot/main.go`:

```go
package main

func main() {
}
```

- [ ] **Step 4: Create all placeholder packages**

Each file contains only the package declaration line. Create these files:

- `internal/config/config.go` → `package config`
- `internal/credentials/credentials.go` → `package credentials`
- `internal/auth/auth.go` → `package auth`
- `internal/chatwoot/client.go` → `package chatwoot`
- `internal/chatwoot/application/application.go` → `package application`
- `internal/chatwoot/platform/platform.go` → `package platform`
- `internal/chatwoot/clientapi/clientapi.go` → `package clientapi`
- `internal/chatwoot/reports/reports.go` → `package reports`
- `internal/testutil/testutil.go` → `package testutil`

- [ ] **Step 5: Verify the project builds**

Run:
```bash
go build ./...
```
Expected: no errors, no output.

- [ ] **Step 6: Commit**

```bash
git add go.mod .gitignore cmd/ internal/
git commit -m "chore: initialize Go module and project scaffolding"
```

---

### Task 2: Add Taskfile and Linting Configuration

**Files:**
- Create: `Taskfile.yml`
- Create: `.golangci.yml`

- [ ] **Step 1: Create Taskfile.yml**

```yaml
version: '3'

vars:
  GIT_TAG:
    sh: git describe --tags --always --dirty 2>/dev/null || echo "dev"
  GIT_COMMIT:
    sh: git rev-parse --short HEAD
  BUILD_DATE:
    sh: date -u '+%Y-%m-%dT%H:%M:%SZ'
  LDFLAGS: >-
    -X github.com/chatwoot/chatwoot-cli/internal/version.Version={{.GIT_TAG}}
    -X github.com/chatwoot/chatwoot-cli/internal/version.Commit={{.GIT_COMMIT}}
    -X github.com/chatwoot/chatwoot-cli/internal/version.Date={{.BUILD_DATE}}

tasks:
  build:
    desc: Build the chatwoot binary
    cmds:
      - go build -ldflags "{{.LDFLAGS}}" -o chatwoot ./cmd/chatwoot/

  test:
    desc: Run all tests
    cmds:
      - go test ./...

  lint:
    desc: Run linters
    cmds:
      - golangci-lint run ./...

  fmt:
    desc: Format all Go files
    cmds:
      - gofmt -w .

  vet:
    desc: Run go vet
    cmds:
      - go vet ./...
```

- [ ] **Step 2: Create .golangci.yml**

```yaml
linters:
  enable:
    - errcheck
    - staticcheck
    - govet
    - gofmt
```

- [ ] **Step 3: Verify task build works**

Run:
```bash
task build
```
Expected: produces a `chatwoot` binary in the project root (does nothing when run, since main is empty).

- [ ] **Step 4: Verify task lint passes**

Run:
```bash
task lint
```
Expected: no errors, no warnings.

- [ ] **Step 5: Commit**

```bash
git add Taskfile.yml .golangci.yml
git commit -m "chore: add Taskfile and golangci-lint configuration"
```

---

### Task 3: Implement Version Package

**Files:**
- Create: `internal/version/version.go`
- Create: `internal/version/version_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/version/version_test.go`:

```go
package version

import (
	"runtime"
	"testing"
)

func TestInfoReturnsGoVersion(t *testing.T) {
	info := Info()
	if info["go_version"] != runtime.Version() {
		t.Errorf("go_version = %q, want %q", info["go_version"], runtime.Version())
	}
}

func TestInfoFallbackWhenLdflagsEmpty(t *testing.T) {
	// When ldflags are not set, Version/Commit/Date are empty strings.
	// Info() should still return non-empty values (from buildinfo or "unknown").
	info := Info()
	for _, key := range []string{"version", "commit", "date", "go_version"} {
		if info[key] == "" {
			t.Errorf("Info()[%q] is empty, expected a value", key)
		}
	}
}

func TestInfoUsesLdflagsWhenSet(t *testing.T) {
	// Temporarily set the package vars to simulate ldflags injection.
	origVersion, origCommit, origDate := Version, Commit, Date
	Version = "1.2.3"
	Commit = "abc1234"
	Date = "2026-01-01T00:00:00Z"
	defer func() {
		Version, Commit, Date = origVersion, origCommit, origDate
	}()

	info := Info()
	if info["version"] != "1.2.3" {
		t.Errorf("version = %q, want %q", info["version"], "1.2.3")
	}
	if info["commit"] != "abc1234" {
		t.Errorf("commit = %q, want %q", info["commit"], "abc1234")
	}
	if info["date"] != "2026-01-01T00:00:00Z" {
		t.Errorf("date = %q, want %q", info["date"], "2026-01-01T00:00:00Z")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:
```bash
go test ./internal/version/ -v
```
Expected: FAIL — `Info` is not defined.

- [ ] **Step 3: Write the implementation**

Create `internal/version/version.go`:

```go
package version

import (
	"runtime"
	"runtime/debug"
)

// Version, Commit, and Date are set by ldflags at build time.
var (
	Version string
	Commit  string
	Date    string
)

// Info returns version metadata as a string map.
// It prefers ldflags values and falls back to debug.ReadBuildInfo().
func Info() map[string]string {
	v := Version
	c := Commit
	d := Date

	if v == "" || c == "" || d == "" {
		if bi, ok := debug.ReadBuildInfo(); ok {
			if v == "" {
				v = bi.Main.Version
			}
			for _, s := range bi.Settings {
				switch s.Key {
				case "vcs.revision":
					if c == "" {
						c = s.Value
						if len(c) > 7 {
							c = c[:7]
						}
					}
				case "vcs.time":
					if d == "" {
						d = s.Value
					}
				}
			}
		}
	}

	if v == "" {
		v = "unknown"
	}
	if c == "" {
		c = "unknown"
	}
	if d == "" {
		d = "unknown"
	}

	return map[string]string{
		"version":    v,
		"commit":     c,
		"date":       d,
		"go_version": runtime.Version(),
	}
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run:
```bash
go test ./internal/version/ -v
```
Expected: all 3 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/version/
git commit -m "feat(version): add version package with ldflags and buildinfo fallback"
```

---

### Task 4: Implement Contract Types and Error Codes

**Files:**
- Create: `internal/contract/contract.go`
- Create: `internal/contract/errors.go`

- [ ] **Step 1: Create contract types**

Create `internal/contract/contract.go`:

```go
// Package contract defines the JSON envelope types for all CLI stdout output.
package contract

// Response is the top-level envelope for all stdout output.
// Use the constructor functions (Success, SuccessList, Err, ErrWithDetail)
// to build responses — do not construct Response literals directly.
type Response struct {
	OK       bool         `json:"ok"`
	Data     any          `json:"data,omitempty"`
	Meta     *Meta        `json:"meta,omitempty"`
	Warnings []Warning    `json:"warnings,omitempty"`
	Error    *ErrorDetail `json:"error,omitempty"`
}

// Meta holds response metadata such as pagination.
type Meta struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination describes the pagination state for collection responses.
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

// ErrorDetail describes an error in the response envelope.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  any    `json:"detail,omitempty"`
}

// Warning describes a non-fatal warning attached to a success response.
type Warning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
```

- [ ] **Step 2: Create error code constants**

Create `internal/contract/errors.go`:

```go
package contract

// Standard error codes used across all commands.
const (
	ErrCodeUnauthorized    = "unauthorized"
	ErrCodeForbidden       = "forbidden"
	ErrCodeNotFound        = "not_found"
	ErrCodeValidation      = "validation_error"
	ErrCodeRateLimited     = "rate_limited"
	ErrCodeServer          = "server_error"
	ErrCodeNetwork         = "network_error"
	ErrCodeConfig          = "config_error"
	ErrCodeAuth            = "auth_error"
)
```

- [ ] **Step 3: Verify it compiles**

Run:
```bash
go build ./internal/contract/
```
Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add internal/contract/contract.go internal/contract/errors.go
git commit -m "feat(contract): add envelope types and error code constants"
```

---

### Task 5: Implement Renderer Functions

**Files:**
- Create: `internal/contract/render.go`

- [ ] **Step 1: Write the failing test for Success**

Create `internal/contract/contract_test.go`:

```go
package contract

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestSuccessSingleResource(t *testing.T) {
	resp := Success(map[string]any{"id": 42, "name": "Support Inbox"})

	if !resp.OK {
		t.Fatal("expected ok to be true")
	}
	if resp.Error != nil {
		t.Fatal("expected error to be nil")
	}

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if raw["ok"] != true {
		t.Errorf("ok = %v, want true", raw["ok"])
	}
	data, ok := raw["data"].(map[string]any)
	if !ok {
		t.Fatal("data is not an object")
	}
	if data["id"] != float64(42) {
		t.Errorf("data.id = %v, want 42", data["id"])
	}
	if _, exists := raw["meta"]; exists {
		t.Error("meta should be absent for single resource")
	}
	if _, exists := raw["warnings"]; exists {
		t.Error("warnings should be absent when empty")
	}
	if _, exists := raw["error"]; exists {
		t.Error("error should be absent on success")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:
```bash
go test ./internal/contract/ -run TestSuccessSingleResource -v
```
Expected: FAIL — `Success` is not defined.

- [ ] **Step 3: Write the renderer implementation**

Create `internal/contract/render.go`:

```go
package contract

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

// Success builds a success envelope for a single resource.
func Success(data any) Response {
	return Response{OK: true, Data: data}
}

// SuccessList builds a success envelope for a collection response.
// It normalizes nil slices to empty slices to ensure JSON output is []
// rather than null.
func SuccessList(data any, meta Meta) Response {
	return Response{OK: true, Data: normalizeSlice(data), Meta: &meta}
}

// Err builds an error envelope.
func Err(code string, message string) Response {
	return Response{OK: false, Error: &ErrorDetail{Code: code, Message: message}}
}

// ErrWithDetail builds an error envelope with additional detail.
func ErrWithDetail(code string, message string, detail any) Response {
	return Response{OK: false, Error: &ErrorDetail{Code: code, Message: message, Detail: detail}}
}

// Write serializes a Response as JSON to w.
// It validates envelope invariants before serializing.
// When pretty is true, output is indented with 2 spaces.
// A trailing newline is always appended.
func Write(w io.Writer, resp Response, pretty bool) error {
	if err := validate(resp); err != nil {
		return err
	}

	var out []byte
	var err error
	if pretty {
		out, err = json.MarshalIndent(resp, "", "  ")
	} else {
		out, err = json.Marshal(resp)
	}
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}
	out = append(out, '\n')
	_, err = w.Write(out)
	return err
}

func validate(resp Response) error {
	if resp.OK && resp.Error != nil {
		return errors.New("contract violation: ok is true but error is set")
	}
	if !resp.OK && resp.Data != nil {
		return errors.New("contract violation: ok is false but data is set")
	}
	if !resp.OK && resp.Error == nil {
		return errors.New("contract violation: ok is false but error is nil")
	}
	return nil
}

// normalizeSlice ensures a nil slice becomes an empty slice so that
// json.Marshal produces [] instead of null.
func normalizeSlice(data any) any {
	if data == nil {
		return []any{}
	}
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice && v.IsNil() {
		return reflect.MakeSlice(v.Type(), 0, 0).Interface()
	}
	return data
}
```

- [ ] **Step 4: Run test to verify it passes**

Run:
```bash
go test ./internal/contract/ -run TestSuccessSingleResource -v
```
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/contract/render.go internal/contract/contract_test.go
git commit -m "feat(contract): add renderer functions with envelope validation"
```

---

### Task 6: Complete Contract Tests

**Files:**
- Modify: `internal/contract/contract_test.go`

- [ ] **Step 1: Add remaining contract tests**

Append the following tests to `internal/contract/contract_test.go`:

```go
func TestSuccessCollection(t *testing.T) {
	items := []map[string]any{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}
	meta := Meta{
		Pagination: &Pagination{
			Page: 1, PerPage: 25, TotalCount: 142, TotalPages: 6,
		},
	}
	resp := SuccessList(items, meta)

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	data, ok := raw["data"].([]any)
	if !ok {
		t.Fatal("data is not an array")
	}
	if len(data) != 2 {
		t.Errorf("data length = %d, want 2", len(data))
	}

	metaRaw, ok := raw["meta"].(map[string]any)
	if !ok {
		t.Fatal("meta is not an object")
	}
	pag, ok := metaRaw["pagination"].(map[string]any)
	if !ok {
		t.Fatal("pagination is not an object")
	}
	if pag["total_count"] != float64(142) {
		t.Errorf("total_count = %v, want 142", pag["total_count"])
	}
	if pag["per_page"] != float64(25) {
		t.Errorf("per_page = %v, want 25", pag["per_page"])
	}
}

func TestSuccessEmptyCollection(t *testing.T) {
	var items []map[string]any // nil slice
	meta := Meta{
		Pagination: &Pagination{
			Page: 1, PerPage: 25, TotalCount: 0, TotalPages: 0,
		},
	}
	resp := SuccessList(items, meta)

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	data, ok := raw["data"].([]any)
	if !ok {
		t.Fatalf("data is not an array, got %T", raw["data"])
	}
	if len(data) != 0 {
		t.Errorf("data length = %d, want 0", len(data))
	}
}

func TestSuccessWithWarnings(t *testing.T) {
	resp := Success(map[string]any{"id": 42})
	resp.Warnings = []Warning{
		{Code: "deprecated_endpoint", Message: "This endpoint will be removed"},
	}

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	warnings, ok := raw["warnings"].([]any)
	if !ok {
		t.Fatal("warnings is not an array")
	}
	if len(warnings) != 1 {
		t.Errorf("warnings length = %d, want 1", len(warnings))
	}
}

func TestSuccessWithoutWarningsOmitsField(t *testing.T) {
	resp := Success(map[string]any{"id": 42})

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	if bytes.Contains(buf.Bytes(), []byte(`"warnings"`)) {
		t.Error("warnings field should be absent when empty")
	}
}

func TestErrorEnvelope(t *testing.T) {
	resp := Err(ErrCodeNotFound, "Conversation 12345 not found")

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if raw["ok"] != false {
		t.Errorf("ok = %v, want false", raw["ok"])
	}
	if _, exists := raw["data"]; exists {
		t.Error("data should be absent on error")
	}

	errObj, ok := raw["error"].(map[string]any)
	if !ok {
		t.Fatal("error is not an object")
	}
	if errObj["code"] != "not_found" {
		t.Errorf("error.code = %v, want not_found", errObj["code"])
	}
	if errObj["message"] != "Conversation 12345 not found" {
		t.Errorf("error.message = %v, want expected message", errObj["message"])
	}
}

func TestErrorWithDetail(t *testing.T) {
	detail := map[string]any{"fields": map[string]any{"email": "is not valid"}}
	resp := ErrWithDetail(ErrCodeValidation, "Invalid request", detail)

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	errObj := raw["error"].(map[string]any)
	if errObj["detail"] == nil {
		t.Error("error.detail should be present")
	}
}

func TestPrettyOutput(t *testing.T) {
	resp := Success(map[string]any{"id": 1})

	var compact bytes.Buffer
	if err := Write(&compact, resp, false); err != nil {
		t.Fatalf("Write compact error: %v", err)
	}

	var pretty bytes.Buffer
	if err := Write(&pretty, resp, true); err != nil {
		t.Fatalf("Write pretty error: %v", err)
	}

	if bytes.Contains(compact.Bytes(), []byte("  ")) {
		t.Error("compact output should not contain indentation")
	}
	if !bytes.Contains(pretty.Bytes(), []byte("  ")) {
		t.Error("pretty output should contain indentation")
	}
	if pretty.Len() <= compact.Len() {
		t.Error("pretty output should be longer than compact")
	}
}

func TestTrailingNewline(t *testing.T) {
	resp := Success(map[string]any{"id": 1})

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	if buf.Bytes()[buf.Len()-1] != '\n' {
		t.Error("output should end with a newline")
	}
}

func TestValidationRejectsOkTrueWithError(t *testing.T) {
	resp := Response{OK: true, Data: "x", Error: &ErrorDetail{Code: "bad", Message: "bad"}}
	var buf bytes.Buffer
	err := Write(&buf, resp, false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidationRejectsOkFalseWithData(t *testing.T) {
	resp := Response{OK: false, Data: "x", Error: &ErrorDetail{Code: "bad", Message: "bad"}}
	var buf bytes.Buffer
	err := Write(&buf, resp, false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidationRejectsOkFalseWithoutError(t *testing.T) {
	resp := Response{OK: false}
	var buf bytes.Buffer
	err := Write(&buf, resp, false)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestSnakeCaseFieldNames(t *testing.T) {
	meta := Meta{
		Pagination: &Pagination{
			Page: 1, PerPage: 10, TotalCount: 50, TotalPages: 5,
		},
	}
	resp := SuccessList([]any{}, meta)

	var buf bytes.Buffer
	if err := Write(&buf, resp, false); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	output := buf.String()
	for _, field := range []string{"per_page", "total_count", "total_pages"} {
		if !bytes.Contains(buf.Bytes(), []byte(field)) {
			t.Errorf("expected snake_case field %q in output: %s", field, output)
		}
	}
	for _, bad := range []string{"perPage", "totalCount", "totalPages", "PerPage", "TotalCount"} {
		if bytes.Contains(buf.Bytes(), []byte(bad)) {
			t.Errorf("unexpected camelCase field %q in output: %s", bad, output)
		}
	}
}
```

- [ ] **Step 2: Run all contract tests**

Run:
```bash
go test ./internal/contract/ -v
```
Expected: all tests PASS.

- [ ] **Step 3: Commit**

```bash
git add internal/contract/contract_test.go
git commit -m "test(contract): add full contract-locking test suite"
```

---

### Task 7: Implement Cobra Root Command with Global Flags

**Files:**
- Create: `internal/cli/root.go`

- [ ] **Step 1: Add dependencies**

Run:
```bash
go get github.com/spf13/cobra github.com/spf13/viper github.com/zalando/go-keyring
```

Note: `viper` and `go-keyring` are not used yet but are listed as initial dependencies in the spec. Adding them now avoids a separate dependency-management step in Layer 2.

- [ ] **Step 2: Create root command**

Create `internal/cli/root.go`:

```go
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	prettyFlag  bool
	verboseFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "chatwoot",
	Short: "Chatwoot CLI — machine-friendly command-line interface for Chatwoot",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Apply env var defaults only when the flag was not explicitly set.
		if !cmd.Flags().Changed("pretty") {
			env := os.Getenv("CHATWOOT_PRETTY")
			if env == "1" || env == "true" {
				prettyFlag = true
			}
		}
		if !cmd.Flags().Changed("verbose") {
			env := os.Getenv("CHATWOOT_VERBOSE")
			if env == "1" || env == "true" {
				verboseFlag = true
			}
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&prettyFlag, "pretty", false, "Indent JSON output with 2 spaces")
	rootCmd.PersistentFlags().BoolVar(&verboseFlag, "verbose", false, "Enable diagnostic logging on stderr")
}

// Execute runs the root command and returns an exit code.
// 0 = success, 1 = runtime/API error, 2 = usage error.
func Execute() int {
	err := rootCmd.Execute()
	if err == nil {
		return 0
	}

	// Detect usage errors (unknown command, bad flags, missing args).
	// Cobra sets SilenceUsage to prevent printing usage on RunE errors,
	// but for unknown-command/flag errors Cobra itself returns an error
	// before RunE executes. We detect this by checking if the error
	// message starts with "unknown command" or "unknown flag".
	errMsg := err.Error()
	if isUsageError(errMsg) {
		fmt.Fprintln(os.Stderr, "Error:", errMsg)
		return 2
	}

	return 1
}

func isUsageError(msg string) bool {
	prefixes := []string{
		"unknown command",
		"unknown flag",
		"unknown shorthand flag",
		"required flag",
		"accepts ",
	}
	for _, p := range prefixes {
		if len(msg) >= len(p) && msg[:len(p)] == p {
			return true
		}
	}
	return false
}

// Pretty returns whether pretty-printing is enabled.
func Pretty() bool {
	return prettyFlag
}
```

- [ ] **Step 3: Verify it compiles**

Run:
```bash
go build ./internal/cli/
```
Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add internal/cli/root.go go.mod go.sum
git commit -m "feat(cli): add Cobra root command with --pretty and --verbose flags"
```

---

### Task 8: Implement Version Command

**Files:**
- Create: `internal/cli/version.go`
- Create: `internal/cli/version_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/cli/version_test.go`:

```go
package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestVersionCommandOutput(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// Reset args to just "version"
	rootCmd.SetArgs([]string{"version"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("invalid JSON output: %v\nraw: %s", err, buf.String())
	}

	if raw["ok"] != true {
		t.Errorf("ok = %v, want true", raw["ok"])
	}
	data, ok := raw["data"].(map[string]any)
	if !ok {
		t.Fatal("data is not an object")
	}
	for _, key := range []string{"version", "commit", "date", "go_version"} {
		if _, exists := data[key]; !exists {
			t.Errorf("data.%s is missing", key)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:
```bash
go test ./internal/cli/ -run TestVersionCommandOutput -v
```
Expected: FAIL — no "version" subcommand registered.

- [ ] **Step 3: Write the version command**

Create `internal/cli/version.go`:

```go
package cli

import (
	"os"

	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build and version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp := contract.Success(version.Info())
		return contract.Write(cmd.OutOrStdout(), resp, prettyFlag)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run:
```bash
go test ./internal/cli/ -run TestVersionCommandOutput -v
```
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/cli/version.go internal/cli/version_test.go
git commit -m "feat(cli): add version command with JSON envelope output"
```

---

### Task 9: Wire main.go Entry Point

**Files:**
- Modify: `cmd/chatwoot/main.go`

- [ ] **Step 1: Update main.go**

Replace `cmd/chatwoot/main.go` with:

```go
package main

import (
	"os"

	"github.com/chatwoot/chatwoot-cli/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
```

Note: Importing the package normally causes all `init()` functions to run, which registers subcommands with the root command.

- [ ] **Step 2: Build and test the binary**

Run:
```bash
task build && ./chatwoot version
```
Expected: JSON output with `ok: true` and version data.

- [ ] **Step 3: Test pretty flag**

Run:
```bash
./chatwoot version --pretty
```
Expected: indented JSON output.

- [ ] **Step 4: Test exit code on success**

Run:
```bash
./chatwoot version; echo "exit: $?"
```
Expected: `exit: 0`

- [ ] **Step 5: Commit**

```bash
git add cmd/chatwoot/main.go
git commit -m "feat: wire main.go entry point with exit code handling"
```

---

### Task 10: Final Verification — All Exit Criteria

**Files:** none (verification only)

- [ ] **Step 1: Run full test suite**

Run:
```bash
task test
```
Expected: all tests pass.

- [ ] **Step 2: Run linter**

Run:
```bash
task lint
```
Expected: no errors, no warnings.

- [ ] **Step 3: Build binary**

Run:
```bash
task build
```
Expected: produces `chatwoot` binary.

- [ ] **Step 4: Run version command**

Run:
```bash
./chatwoot version
```
Expected: valid JSON success envelope with version, commit, date, go_version.

- [ ] **Step 5: Run version with --pretty**

Run:
```bash
./chatwoot version --pretty
```
Expected: indented JSON output.

- [ ] **Step 6: Verify all packages exist**

Run:
```bash
go build ./...
```
Expected: no errors — all placeholder packages compile.

- [ ] **Step 7: Commit any remaining changes and tag**

```bash
git add -A
git status
# If there are changes:
git commit -m "chore: final cleanup for Layer 0+1 completion"
```

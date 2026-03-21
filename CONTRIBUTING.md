# Contributing to chatwoot-cli

## Prerequisites

You need:

- Go 1.26 or later
- [go-task](https://taskfile.dev) (optional; plain `go` commands work fine)

## Getting started

```bash
git clone https://github.com/chatwoot/chatwoot-cli.git
cd chatwoot-cli
go mod download
task build   # or: go build ./cmd/chatwoot/
task test    # or: go test ./...
```

The `task build` target injects version info through ldflags. A plain `go build` produces a working binary without version metadata.

## Project layout

```
cmd/chatwoot/        Entry point (main.go)
internal/
  cli/               Cobra command handlers (thin wrappers)
  contract/          JSON envelope types for stdout
  config/            Viper config, profile resolution
  credentials/       OS keychain integration, env fallbacks
  chatwoot/          HTTP transport and API clients
    reports/         Reports API client
  version/           Build version info
```

Command handlers resolve context, call a client method, and render output. Business logic belongs in the client packages under `internal/chatwoot/`.

## Running tests

```bash
task test
```

API client tests use `net/http/httptest` servers. Domain logic tests use small interfaces backed by in-memory doubles. Avoid mocking frameworks.

When you add a new API endpoint, write at least one `httptest`-based round-trip test that covers the success path and one error path.

## Code style

Format and lint before pushing:

```bash
task fmt    # gofmt -w .
task vet    # go vet ./...
task lint   # golangci-lint run ./...
```

Follow standard Go conventions. Keep exported identifiers minimal. Prefer returning errors over panicking.

stdout emits JSON envelopes only. Diagnostics and logs go to stderr via `log/slog`.

## Commit messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): subject
```

Common types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`.

Scope matches the package or area you changed:

```
feat(reports): add agent summary endpoint
fix(config): fall back to env var when keychain is locked
test(contract): cover error envelope rendering
```

Write the subject in imperative mood ("add X", not "added X"). Keep it under 72 characters.

## Pull requests

Fork the repo, create a feature branch, and open a PR against `main`.

Each PR should:

- Solve one problem
- Include tests for new behavior
- Pass `task lint` and `task test`
- Reference a related issue number when one exists

Describe *what* changed and *why* in the PR body. Reviewers can read the diff for the *how*.

## Reporting bugs

Open a GitHub issue with:

- The `chatwoot version` output
- The command you ran
- The JSON output or error message you received
- What you expected to happen

Include the `--verbose` stderr output when the problem involves API communication.

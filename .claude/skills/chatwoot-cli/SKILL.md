---
name: chatwoot-cli
description: Operational knowledge for the `chatwoot` CLI tool — a machine-friendly, JSON-only command-line interface to the Chatwoot API. Use this skill whenever the user wants to interact with Chatwoot through the command line, manage conversations, contacts, messages, inboxes, agents, teams, reports, or any other Chatwoot API resource. Also trigger when the user mentions chatwoot commands, asks how to script against Chatwoot, wants to automate Chatwoot workflows, pipe Chatwoot data through jq, or build AI agent integrations that talk to a Chatwoot instance. If the user mentions "chatwoot" in a CLI or automation context, use this skill.
---

# Chatwoot CLI

`chatwoot` is a machine-friendly CLI for the Chatwoot API. Every command writes structured JSON to stdout and diagnostics to stderr. No interactive prompts. Built for scripts and AI agents.

## Installation

```bash
go install github.com/chatwoot/chatwoot-cli/cmd/chatwoot@latest
```

## JSON output envelope

Every command returns one of two shapes on stdout:

```json
{"ok": true, "data": {...}, "meta": {"pagination": {"page": 1, "per_page": 15, "total_count": 42, "total_pages": 3}}}
```

```json
{"ok": false, "error": {"code": "not_found", "message": "Conversation not found", "detail": null}}
```

Always check `ok` before reading `data`. Pagination metadata appears in `meta.pagination` on list responses.

## Global flags

| Flag | Type | Purpose |
|------|------|---------|
| `--base-url` | string | Override Chatwoot instance URL |
| `--account-id` | int | Override account ID |
| `--profile` | string | Select a named profile |
| `--pretty` | bool | Indent JSON output |
| `--verbose` | bool | Diagnostic logs to stderr |

## Authentication

Three auth modes correspond to three API families:

**Application API** (agent/admin) — user access token:
```bash
chatwoot auth set --mode application --token <token>
# or env: CHATWOOT_ACCESS_TOKEN
```

**Platform API** (self-hosted admin) — platform app token:
```bash
chatwoot auth set --mode platform --token <token>
# or env: CHATWOOT_PLATFORM_TOKEN
```

**Client API** (public end-user) — no token needed, uses `--inbox-id` and `--contact-id` flags.

Credentials store in the OS keychain per profile. Env vars (`CHATWOOT_ACCESS_TOKEN`, `CHATWOOT_PLATFORM_TOKEN`) serve as fallback for CI.

Check auth status:
```bash
chatwoot auth status
```

## Configuration

Config file: `~/.config/chatwoot-cli/config.yaml`

```yaml
default_profile: production
profiles:
  production:
    base_url: https://app.chatwoot.com
    account_id: 1
  staging:
    base_url: https://staging.example.com
    account_id: 2
```

Resolution precedence: flags > env vars (`CHATWOOT_BASE_URL`, `CHATWOOT_ACCOUNT_ID`, `CHATWOOT_PROFILE`) > config file > defaults.

## Command tree

### `chatwoot application` — agent/admin API (`/api/v1`, `/api/v2`)

Requires application auth (user access token).

| Resource | Subcommands |
|----------|-------------|
| `contacts` | `list`, `get`, `create`, `update`, `delete`, `search`, `filter`, `merge`, `conversations list`, `labels` |
| `conversations` | `list`, `get`, `create`, `update`, `filter`, `meta`, `toggle-status`, `toggle-priority`, `assignments`, `labels` |
| `messages` | `list`, `create`, `delete` (requires `--conversation-id`) |
| `agents` | `list`, `create`, `update`, `delete` |
| `inboxes` | `list`, `get`, `create`, `update`, `members`, `agent-bot` |
| `teams` | `list`, `get`, `create`, `update`, `delete` |
| `labels` | Account label management |
| `webhooks` | `list`, `create`, `update`, `delete` |
| `reports` | `account`, `conversations`, `events`, `summary`, `summary-by-agent`, `summary-by-channel`, `summary-by-inbox`, `summary-by-team`, `first-response-distribution`, `inbox-label-matrix`, `outgoing-messages` |
| `help-center` | `portals`, `categories`, `articles` (each with subcommands) |
| `canned-responses` | Saved reply templates |
| `custom-attributes` | Custom attribute definitions |
| `custom-filters` | Saved filter definitions |
| `automation-rules` | Automation management |
| `integrations` | Integration management |
| `audit-logs` | Read-only audit log access |
| `agent-bots` | Account-level bot management |
| `account` | Account settings |
| `profile` | Authenticated user profile |

### `chatwoot platform` — self-hosted admin API (`/platform/api/v1`)

Requires platform auth (platform app token).

| Resource | Subcommands |
|----------|-------------|
| `accounts` | `get`, `create`, `update`, `delete` |
| `users` | `get`, `create`, `update`, `delete`, `login` (SSO link) |
| `account-users` | Manage user-account associations |
| `agent-bots` | Platform-level bot management |

### `chatwoot client` — public end-user API (`/public/api/v1`)

All commands require `--inbox-id` and `--contact-id`. No access token needed.

| Resource | Subcommands |
|----------|-------------|
| `contacts` | `get`, `create`, `update` |
| `conversations` | `list`, `get`, `create`, `toggle-status`, `toggle-typing`, `update-last-seen` |
| `messages` | `list`, `create`, `update` |

### Utility commands

| Command | Purpose |
|---------|---------|
| `auth` | `set`, `status`, `clear` |
| `version` | Build info as JSON |
| `completion` | `bash`, `zsh`, `fish`, `powershell` |

## Usage patterns

### Parse the envelope first

Always check `ok` before accessing `data`. Handle errors from `error.code` and `error.message`.

### Use `--pretty` for display, skip it for piping

```bash
# Human-readable
chatwoot application contacts list --pretty

# Piped to jq
chatwoot application contacts list | jq '.data[].email'
```

### Paginate with `--page` or `--all`

List commands return paginated results. Use `--page` and `--per-page` to step through, or `--all` to fetch everything at once:

```bash
# Check pagination info
chatwoot application contacts list --page 1 | jq '.meta.pagination'

# Fetch all pages in one call
chatwoot application conversations list --status open --all
```

### Set up auth before running commands

If a command fails with an auth error, check `chatwoot auth status` and run `chatwoot auth set` with the appropriate mode and token.

### Use profiles for multiple instances

Set up profiles in the config file rather than passing `--base-url` on every command:

```bash
chatwoot --profile staging application contacts list
```

### Prefer specific commands

- `contacts search --query` for text search
- `contacts filter` for structured filters
- `contacts list` for browsing

### Contact lookup workflow

```bash
# Search for a contact
chatwoot application contacts search --query "Sarah Chen" | jq '.data[] | {id, name, email}'

# List their conversations (note: --contact-id, not --id)
chatwoot application contacts conversations list --contact-id 42

# Read messages from one of those conversations
chatwoot application messages list --conversation-id 187
```

### Conversation workflow

The typical flow: list/filter conversations, get a specific one, list its messages, create a reply.

```bash
# Find open conversations (--status is on `list`, not `filter`)
chatwoot application conversations list --status open | jq '.data[].id'

# Fetch all pages at once
chatwoot application conversations list --status open --all

# Get conversation details
chatwoot application conversations get --id 42

# List messages
chatwoot application messages list --conversation-id 42

# Send a reply (outgoing by default)
chatwoot application messages create --conversation-id 42 --content "Thanks for reaching out"

# Send a private note (internal, not visible to contact)
chatwoot application messages create --conversation-id 42 --content "Internal note" --private

# Resolve the conversation
chatwoot application conversations toggle-status --id 42 --status resolved
```

**Important:** `conversations list` accepts `--status`, `--inbox-id`, `--page`, `--per-page`, and `--all`. `conversations filter` takes `--payload` (a JSON array) for advanced structured filtering — it does NOT accept `--status` directly.

### Reports need date ranges

Most reports commands require `--since` and `--until` date parameters.

```bash
chatwoot application reports summary --since 2026-01-01 --until 2026-01-31
```

### Client API is for end-user integrations

The client API powers chatbots and widgets. No access token — just inbox and contact identifiers:

```bash
chatwoot client conversations list --inbox-id abc123 --contact-id def456
```

## Key flag reference

These are the exact flags for commonly-used commands, verified against the actual CLI binary. Getting these wrong causes errors, so refer here rather than guessing.

### conversations list
`--status` (open, resolved, pending, snoozed), `--inbox-id`, `--page`, `--per-page`, `--all`

### conversations filter
`--payload` (JSON array — structured filter), `--page`. Does NOT accept `--status`.

### conversations get / toggle-status / toggle-priority
`--id` (conversation ID). `toggle-status` also takes `--status`.

### contacts search
`--query`, `--page`, `--per-page`

### contacts conversations list
`--contact-id` (not `--id`). This is a subcommand: `contacts conversations list --contact-id 5`.

### contacts get
`--id` (contact ID)

### messages list / create / delete
`--conversation-id` (required on all). `create` adds `--content`, `--message-type` (incoming/outgoing, default outgoing), `--private` (send as internal note).

### reports
Most accept `--since`, `--until`, `--type` (account, agent, inbox, team, label).

## Error handling

Failed commands produce:
- Non-zero exit code
- JSON error envelope on stdout: `{"ok": false, "error": {"code": "...", "message": "..."}}`
- With `--verbose`, detailed diagnostics on stderr

Common error codes: `unauthorized`, `not_found`, `validation_error`, `rate_limited`, `server_error`.

When scripting, check the exit code or parse `ok` from the JSON. For `rate_limited`, back off and retry. For `unauthorized`, re-check `chatwoot auth status`.

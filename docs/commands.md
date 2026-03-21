# Command Reference

## Output Format

All commands write JSON to stdout inside a standard envelope:

```json
{"ok": true, "data": {...}, "meta": {"pagination": {...}}}
```

```json
{"ok": false, "error": {"code": "...", "message": "...", "detail": ...}}
```

Pass `--pretty` for indented output. Diagnostics land on stderr when you pass `--verbose`.

## Global Flags

| Flag | Type | Description |
|------|------|-------------|
| `--base-url` | string | Override the Chatwoot instance URL |
| `--account-id` | int | Override the account ID |
| `--profile` | string | Select a named profile |
| `--pretty` | bool | Indent JSON output |
| `--verbose` | bool | Write diagnostic logs to stderr |

---

## Auth Commands

Manage credentials stored in the OS keychain.

### `chatwoot auth set`

Store a credential for the active profile.

```
chatwoot auth set --mode application --token <user-access-token>
chatwoot auth set --mode platform --token <platform-app-token>
```

### `chatwoot auth status`

Show credential status for the active profile.

```
chatwoot auth status
chatwoot auth status --profile staging
```

### `chatwoot auth clear`

Remove stored credentials for the active profile.

```
chatwoot auth clear
```

---

## Application Commands

Target the agent/admin API (`/api/v1`, `/api/v2`). Requires a user access token set via `auth set --mode application`.

```
chatwoot application <resource> <action> [flags]
```

### contacts

| Action | Description |
|--------|-------------|
| `list` | List contacts with pagination |
| `get` | Fetch a single contact by ID |
| `create` | Create a contact |
| `update` | Update a contact |
| `delete` | Delete a contact |
| `search` | Search contacts by query string |
| `filter` | Filter contacts using structured filter payload |
| `merge` | Merge two contact records |
| `conversations` | List conversations for a contact |
| `labels` | Manage labels on a contact |

Examples:

```bash
# Search contacts by email
chatwoot application contacts search --query "alice@example.com" --pretty

# Create a contact
chatwoot application contacts create --name "Alice Wu" --email "alice@example.com"
```

### conversations

| Action | Description |
|--------|-------------|
| `list` | List conversations with pagination |
| `get` | Fetch a conversation by ID |
| `create` | Create a conversation |
| `update` | Update conversation attributes |
| `filter` | Filter conversations with structured payload |
| `meta` | Fetch conversation metadata (counts by status) |
| `toggle-status` | Change conversation status (open, resolved, pending) |
| `toggle-priority` | Set priority (urgent, high, medium, low, none) |
| `assignments` | Assign conversation to agent or team |
| `labels` | Manage labels on a conversation |

Examples:

```bash
# List open conversations assigned to you
chatwoot application conversations list --status open --assignee-type me

# Resolve a conversation
chatwoot application conversations toggle-status --id 42 --status resolved
```

### messages

Scoped to a conversation via `--conversation-id`.

| Action | Description |
|--------|-------------|
| `list` | List messages in a conversation |
| `create` | Send a message |
| `delete` | Delete a message |

Examples:

```bash
# Send a reply
chatwoot application messages create --conversation-id 42 --content "We shipped the fix."

# List recent messages
chatwoot application messages list --conversation-id 42
```

### agents

| Action | Description |
|--------|-------------|
| `list` | List agents in the account |
| `create` | Add an agent |
| `update` | Update agent details |
| `delete` | Remove an agent |

### inboxes

| Action | Description |
|--------|-------------|
| `list` | List inboxes |
| `get` | Fetch inbox by ID |
| `create` | Create an inbox |
| `update` | Update inbox settings |
| `members` | Manage inbox agent membership |
| `agent-bot` | Get or set the agent bot for an inbox |

### webhooks

| Action | Description |
|--------|-------------|
| `list` | List webhooks |
| `create` | Create a webhook |
| `update` | Update a webhook |
| `delete` | Delete a webhook |

### reports

| Action | Description |
|--------|-------------|
| `account` | Account-level report metrics |
| `conversations` | Conversation volume report |
| `events` | Event timeline data |
| `summary` | Account summary stats |
| `summary-by-agent` | Summary broken down by agent |
| `summary-by-channel` | Summary broken down by channel |
| `summary-by-inbox` | Summary broken down by inbox |
| `summary-by-team` | Summary broken down by team |
| `first-response-distribution` | First response time distribution |
| `inbox-label-matrix` | Inbox vs. label cross-tabulation |
| `outgoing-messages` | Outgoing message metrics |

### Other Application Resources

These resources follow the same `list`, `create`, `update`, `delete` pattern unless noted:

| Resource | Notes |
|----------|-------|
| `account` | Account settings |
| `agent-bots` | Account-level agent bots |
| `audit-logs` | Read-only audit log access |
| `automation-rules` | Automation rule management |
| `canned-responses` | Saved reply templates |
| `custom-attributes` | Custom attribute definitions |
| `custom-filters` | Saved filter definitions |
| `help-center` | Subresources: `portals`, `categories`, `articles` |
| `integrations` | Integration management |
| `labels` | Account label definitions |
| `profile` | Authenticated user profile |
| `teams` | Team management |

---

## Platform Commands

Target the self-hosted admin API (`/platform/api/v1`). Requires a platform app token set via `auth set --mode platform`.

```
chatwoot platform <resource> <action> [flags]
```

### accounts

| Action | Description |
|--------|-------------|
| `get` | Fetch account details |
| `create` | Provision a new account |
| `update` | Update account settings |
| `delete` | Delete an account |

### users

| Action | Description |
|--------|-------------|
| `get` | Fetch user details |
| `create` | Create a user |
| `update` | Update user details |
| `delete` | Delete a user |
| `login` | Generate an SSO login link |

### account-users

Manage the association between users and accounts.

### agent-bots

Platform-level agent bot management (separate from account-level bots under `application`).

---

## Client Commands

Target the public end-user API (`/public/api/v1`). No token required. Authenticate with `--inbox-id` and `--contact-id` flags.

```
chatwoot client <resource> <action> --inbox-id <id> --contact-id <id> [flags]
```

### contacts

| Action | Description |
|--------|-------------|
| `get` | Fetch contact details |
| `create` | Create a public contact |
| `update` | Update contact details |

### conversations

| Action | Description |
|--------|-------------|
| `list` | List conversations for the contact |
| `get` | Fetch a conversation |
| `create` | Start a conversation |
| `toggle-status` | Change conversation status |
| `toggle-typing` | Send typing indicator |
| `update-last-seen` | Mark conversation as seen |

### messages

| Action | Description |
|--------|-------------|
| `list` | List messages in a conversation |
| `create` | Send a message |
| `update` | Update a message |

---

## Utility Commands

### `chatwoot version`

Prints build and version info as JSON.

```bash
chatwoot version
```

### `chatwoot completion`

Generates shell completion scripts.

```bash
chatwoot completion bash
chatwoot completion zsh
chatwoot completion fish
chatwoot completion powershell
```

Source the output in your shell config to enable tab completion.

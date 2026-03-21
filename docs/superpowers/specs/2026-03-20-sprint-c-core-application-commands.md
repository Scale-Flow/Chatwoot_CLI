# Design Spec: Sprint C — Core Application Commands (Layer 7)

## Scope

This spec defines the full core application command surface for the Chatwoot
CLI: contacts, conversations, messages, and inboxes. It covers API client
methods, CLI command packages, model types, shared helper extraction, testing
strategy, and existing code changes.

## Goals

- Implement the complete stable command surface from the command structure doc
  for the four core resource families
- Enable the key agent workflows: triage, summarize, reply, resolve, escalate
- Establish reusable patterns for all future command families
- Extract shared CLI helpers to eliminate import-cycle duplication from Sprint B

## Prerequisites

Sprint A (plumbing) and Sprint B (CLI shell + foundation commands) are
complete. The following exist:

- Config, credentials, transport, retry, error mapping
- Application, Platform, and Client API client packages
- CLI shell with auth and profile commands
- JSON envelope contracts and rendering
- `chatwoot.ListAll` pagination helper

## Shared Helper Extraction: `cmdutil`

### Problem

Sprint B revealed an import cycle: `internal/cli` imports
`internal/cli/application` (to register commands), so `internal/cli/application`
cannot import `internal/cli` (to use ResolveContext, etc.). Sprint B worked
around this by inlining helpers. With 4 new command packages, inlining creates
unacceptable duplication.

### Solution

Extract shared CLI helpers into `internal/cli/cmdutil/`:

```
internal/cli/cmdutil/
  context.go    — RuntimeContext, ResolveContext, ResolveAuth, WriteError
  flags.go      — PaginationFlags, Pretty() accessor
```

**`context.go` contents:**

- `RuntimeContext` struct (ProfileName, BaseURL, AccountID)
- `ResolveContext(cmd) (*RuntimeContext, error)` — reads --profile/--base-url/--account-id flags, loads config, resolves profile, validates BaseURL
- `resolveContextFromPath(cfgPath, flagProfile, flagBaseURL, flagAccountID) (*RuntimeContext, error)` — testable core
- `ResolveAuth(profileName, mode) (auth.TokenAuth, error)` — resolves credentials from env/keychain/file
- `WriteError(cmd, code, message) error` — writes error envelope to stdout and returns error for Cobra

**`flags.go` contents:**

- `Pretty(cmd) bool` — reads --pretty from root persistent flags
- `PaginationFlags` struct with `Page int`, `PerPage int`, `All bool`
- `AddPaginationFlags(cmd)` — registers --page, --per-page, --all on a command
- `GetPaginationFlags(cmd) PaginationFlags` — reads pagination flag values

**Migration:**

- `internal/cli/context.go`: Remove functions that move to cmdutil. Keep
  `ResolveProfileName` if still referenced by auth, or remove if auth uses its
  own local copy.
- `internal/cli/application/profile.go`: Replace inlined helpers with cmdutil
  imports.
- `internal/cli/auth/set.go`: Keep local `resolveProfileNameForAuth` and
  `prettyFromRoot` — auth commands use a lightweight path that doesn't need
  the full context pipeline.

### Pagination Convention

All list commands register three flags via `cmdutil.AddPaginationFlags`:

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--page` | int | 1 | Page number |
| `--per-page` | int | 25 | Items per page |
| `--all` | bool | false | Auto-paginate and return all results |

When `--all` is set, the command uses the existing `chatwoot.ListAll` helper
to fetch all pages sequentially. The response envelope always includes
`meta.pagination`.

## API Client Methods

All methods live in `internal/chatwoot/application/`. Each resource gets its
own file. Methods follow the existing pattern: build path, call
`transport.DoWithRetry`, decode response, return typed data.

### `contacts.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListContacts(ctx, ListContactsOpts)` | GET | `/api/v1/accounts/{id}/contacts` | `([]Contact, *Pagination, error)` |
| `GetContact(ctx, id)` | GET | `/api/v1/accounts/{id}/contacts/{cid}` | `(*Contact, error)` |
| `CreateContact(ctx, CreateContactOpts)` | POST | `/api/v1/accounts/{id}/contacts` | `(*Contact, error)` |
| `UpdateContact(ctx, id, UpdateContactOpts)` | PUT | `/api/v1/accounts/{id}/contacts/{cid}` | `(*Contact, error)` |
| `DeleteContact(ctx, id)` | DELETE | `/api/v1/accounts/{id}/contacts/{cid}` | `error` |
| `SearchContacts(ctx, query, page)` | GET | `/api/v1/accounts/{id}/contacts/search` | `([]Contact, *Pagination, error)` |
| `FilterContacts(ctx, FilterContactsOpts)` | POST | `/api/v1/accounts/{id}/contacts/filter` | `([]Contact, *Pagination, error)` |
| `MergeContacts(ctx, baseID, mergeID)` | POST | `/api/v1/accounts/{id}/actions/contact_merge` | `(*Contact, error)` |
| `ListContactLabels(ctx, id)` | GET | `/api/v1/accounts/{id}/contacts/{cid}/labels` | `([]string, error)` |
| `SetContactLabels(ctx, id, labels)` | POST | `/api/v1/accounts/{id}/contacts/{cid}/labels` | `([]string, error)` |
| `ListContactConversations(ctx, id)` | GET | `/api/v1/accounts/{id}/contacts/{cid}/conversations` | `([]Conversation, error)` |

### `conversations.go` (extend existing)

Existing methods: `ListConversations`, `GetConversation`.

`ListConversations` signature changes to return pagination:
`([]Conversation, *Pagination, error)`.

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `CreateConversation(ctx, CreateConversationOpts)` | POST | `/api/v1/accounts/{id}/conversations` | `(*Conversation, error)` |
| `UpdateConversation(ctx, id, UpdateConversationOpts)` | PATCH | `/api/v1/accounts/{id}/conversations/{cid}` | `(*Conversation, error)` |
| `FilterConversations(ctx, FilterConversationsOpts)` | POST | `/api/v1/accounts/{id}/conversations/filter` | `([]Conversation, *Pagination, error)` |
| `GetConversationMeta(ctx)` | GET | `/api/v1/accounts/{id}/conversations/meta` | `(*ConversationMeta, error)` |
| `ToggleConversationStatus(ctx, id, status)` | POST | `/api/v1/accounts/{id}/conversations/{cid}/toggle_status` | `(*Conversation, error)` |
| `ToggleConversationPriority(ctx, id, priority)` | POST | `/api/v1/accounts/{id}/conversations/{cid}/toggle_priority` | `(*Conversation, error)` |
| `AssignConversation(ctx, id, AssignOpts)` | POST | `/api/v1/accounts/{id}/conversations/{cid}/assignments` | `(*Conversation, error)` |
| `ListConversationLabels(ctx, id)` | GET | `/api/v1/accounts/{id}/conversations/{cid}/labels` | `([]string, error)` |
| `SetConversationLabels(ctx, id, labels)` | POST | `/api/v1/accounts/{id}/conversations/{cid}/labels` | `([]string, error)` |

### `messages.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListMessages(ctx, conversationID)` | GET | `/api/v1/accounts/{id}/conversations/{cid}/messages` | `([]Message, error)` |
| `CreateMessage(ctx, conversationID, CreateMessageOpts)` | POST | `/api/v1/accounts/{id}/conversations/{cid}/messages` | `(*Message, error)` |
| `DeleteMessage(ctx, conversationID, messageID)` | DELETE | `/api/v1/accounts/{id}/conversations/{cid}/messages/{mid}` | `error` |

### `inboxes.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListInboxes(ctx)` | GET | `/api/v1/accounts/{id}/inboxes` | `([]Inbox, error)` |
| `GetInbox(ctx, id)` | GET | `/api/v1/accounts/{id}/inboxes/{iid}` | `(*Inbox, error)` |
| `CreateInbox(ctx, CreateInboxOpts)` | POST | `/api/v1/accounts/{id}/inboxes` | `(*Inbox, error)` |
| `UpdateInbox(ctx, id, UpdateInboxOpts)` | PATCH | `/api/v1/accounts/{id}/inboxes/{iid}` | `(*Inbox, error)` |
| `ListInboxMembers(ctx, inboxID)` | GET | `/api/v1/accounts/{id}/inbox_members/{iid}` | `([]Agent, error)` |
| `AddInboxMember(ctx, inboxID, agentIDs)` | POST | `/api/v1/accounts/{id}/inbox_members` | `([]Agent, error)` |
| `UpdateInboxMembers(ctx, inboxID, agentIDs)` | PATCH | `/api/v1/accounts/{id}/inbox_members` | `([]Agent, error)` |
| `RemoveInboxMember(ctx, inboxID, agentIDs)` | DELETE | `/api/v1/accounts/{id}/inbox_members` | `error` |
| `GetInboxAgentBot(ctx, inboxID)` | GET | `/api/v1/accounts/{id}/inboxes/{iid}/agent_bot` | `(*AgentBot, error)` |
| `SetInboxAgentBot(ctx, inboxID, agentBotID)` | POST | `/api/v1/accounts/{id}/inboxes/{iid}/set_agent_bot` | `(*AgentBot, error)` |

## Model Types

All types in `internal/chatwoot/application/models.go`.

### Core Types

```go
type Contact struct {
    ID        int              `json:"id"`
    Name      string           `json:"name"`
    Email     string           `json:"email,omitempty"`
    Phone     string           `json:"phone_number,omitempty"`
    AccountID int              `json:"account_id"`
    CreatedAt chatwoot.Timestamp `json:"created_at,omitempty"`
}

type Message struct {
    ID             int              `json:"id"`
    Content        string           `json:"content,omitempty"`
    MessageType    int              `json:"message_type"`
    ContentType    string           `json:"content_type,omitempty"`
    Private        bool             `json:"private"`
    ConversationID int              `json:"conversation_id"`
    CreatedAt      chatwoot.Timestamp `json:"created_at,omitempty"`
}

type Inbox struct {
    ID                   int    `json:"id"`
    Name                 string `json:"name"`
    ChannelType          string `json:"channel_type,omitempty"`
    AvatarURL            string `json:"avatar_url,omitempty"`
    EnableAutoAssignment bool   `json:"enable_auto_assignment"`
}

type Agent struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    Role      string `json:"role,omitempty"`
}

type AgentBot struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type ConversationMeta struct {
    AllCount      int `json:"all_count"`
    OpenCount     int `json:"open_count"`
    ResolvedCount int `json:"resolved_count"`
    PendingCount  int `json:"pending_count"`
    SnoozedCount  int `json:"snoozed_count"`
}
```

### Opts Types (mutations)

All use pointer fields so only explicitly set values are serialized.

```go
type ListContactsOpts struct {
    Page    int
    PerPage int
}

type CreateContactOpts struct {
    Name  string  `json:"name"`
    Email *string `json:"email,omitempty"`
    Phone *string `json:"phone_number,omitempty"`
}

type UpdateContactOpts struct {
    Name  *string `json:"name,omitempty"`
    Email *string `json:"email,omitempty"`
    Phone *string `json:"phone_number,omitempty"`
}

type FilterContactsOpts struct {
    Page    int   `json:"page,omitempty"`
    Payload []any `json:"payload"`
}

type CreateConversationOpts struct {
    ContactID int    `json:"contact_id"`
    InboxID   int    `json:"inbox_id"`
    Status    string `json:"status,omitempty"`
    Message   *struct {
        Content string `json:"content"`
    } `json:"message,omitempty"`
}

type UpdateConversationOpts struct {
    Status   *string `json:"status,omitempty"`
    Priority *string `json:"priority,omitempty"`
}

type FilterConversationsOpts struct {
    Page    int   `json:"page,omitempty"`
    Payload []any `json:"payload"`
}

type AssignOpts struct {
    AgentID *int `json:"assignee_id,omitempty"`
    TeamID  *int `json:"team_id,omitempty"`
}

type CreateMessageOpts struct {
    Content     string `json:"content"`
    MessageType string `json:"message_type,omitempty"`
    Private     bool   `json:"private,omitempty"`
}

type CreateInboxOpts struct {
    Name    string `json:"name"`
    Channel any    `json:"channel"`
}

type UpdateInboxOpts struct {
    Name                 *string `json:"name,omitempty"`
    EnableAutoAssignment *bool   `json:"enable_auto_assignment,omitempty"`
}
```

## CLI Command Packages

### Package Structure

```
internal/cli/application/
  contacts/
    contacts.go          — Cmd group, register subcommands
    list.go              — contacts list (--page, --per-page, --all)
    get.go               — contacts get --id
    create.go            — contacts create --name [--email, --phone]
    update.go            — contacts update --id [--name, --email, --phone]
    delete.go            — contacts delete --id
    search.go            — contacts search --query [--page, --per-page]
    filter.go            — contacts filter --payload (JSON)
    merge.go             — contacts merge --base-id --merge-id
    labels.go            — contacts labels list --id / contacts labels set --id --labels
    conversations.go     — contacts conversations list --id
    contacts_test.go
    testroot_test.go
  conversations/
    conversations.go     — Cmd group
    list.go              — conversations list (--status, --inbox-id, --page, --per-page, --all)
    get.go               — conversations get --id
    create.go            — conversations create (--contact-id, --inbox-id, etc.)
    update.go            — conversations update --id (--status, --priority)
    filter.go            — conversations filter --payload (JSON)
    meta.go              — conversations meta
    toggle_status.go     — conversations toggle-status --id --status
    toggle_priority.go   — conversations toggle-priority --id --priority
    assignments.go       — conversations assignments create --id (--agent-id, --team-id)
    labels.go            — conversations labels list --id / conversations labels set --id --labels
    conversations_test.go
    testroot_test.go
  messages/
    messages.go          — Cmd group
    list.go              — messages list --conversation-id
    create.go            — messages create --conversation-id --content [--message-type, --private]
    delete.go            — messages delete --conversation-id --id
    messages_test.go
    testroot_test.go
  inboxes/
    inboxes.go           — Cmd group
    list.go              — inboxes list
    get.go               — inboxes get --id
    create.go            — inboxes create --name --channel (JSON)
    update.go            — inboxes update --id [--name, --enable-auto-assignment]
    members.go           — inboxes members list/add/update/delete --inbox-id [--agent-ids]
    agent_bot.go         — inboxes agent-bot get/set --inbox-id [--agent-bot-id]
    inboxes_test.go
    testroot_test.go
```

### Registration

`internal/cli/application/application.go` registers all subgroups:

```go
import (
    "github.com/chatwoot/chatwoot-cli/internal/cli/application/contacts"
    "github.com/chatwoot/chatwoot-cli/internal/cli/application/conversations"
    "github.com/chatwoot/chatwoot-cli/internal/cli/application/inboxes"
    "github.com/chatwoot/chatwoot-cli/internal/cli/application/messages"
)

func init() {
    Cmd.AddCommand(profileCmd)
    Cmd.AddCommand(contacts.Cmd)
    Cmd.AddCommand(conversations.Cmd)
    Cmd.AddCommand(messages.Cmd)
    Cmd.AddCommand(inboxes.Cmd)
}
```

### Command Handler Pattern

Every command handler follows the same thin pipeline:

```go
func runContactsGet(cmd *cobra.Command, args []string) error {
    rctx, err := cmdutil.ResolveContext(cmd)
    if err != nil {
        return cmdutil.WriteError(cmd, contract.ErrCodeConfig, err.Error())
    }

    tokenAuth, err := cmdutil.ResolveAuth(rctx.ProfileName, credentials.ModeApplication)
    if err != nil {
        return cmdutil.WriteError(cmd, contract.ErrCodeAuth, err.Error())
    }

    transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
    client := appapi.NewClient(transport, rctx.AccountID)

    id, _ := cmd.Flags().GetInt("id")
    contact, err := client.GetContact(context.Background(), id)
    if err != nil {
        return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
    }

    resp := contract.Success(contact)
    return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
```

List commands with pagination:

```go
func runContactsList(cmd *cobra.Command, args []string) error {
    // ... resolve context, auth, build client ...

    pf := cmdutil.GetPaginationFlags(cmd)

    if pf.All {
        // Use chatwoot.ListAll for auto-pagination
        allContacts, err := chatwoot.ListAll(ctx, func(page int) ([]Contact, *Pagination, error) {
            return client.ListContacts(ctx, ListContactsOpts{Page: page, PerPage: pf.PerPage})
        })
        // ... render with total pagination meta
    } else {
        contacts, pagination, err := client.ListContacts(ctx, ListContactsOpts{Page: pf.Page, PerPage: pf.PerPage})
        // ... render with pagination meta
    }
}
```

### Flag Conventions

| Flag | Type | Used by | Description |
|------|------|---------|-------------|
| `--id` | int | get, update, delete commands | Primary resource ID |
| `--conversation-id` | int | messages commands | Parent conversation ID |
| `--inbox-id` | int | conversations list, inbox members | Parent inbox ID |
| `--contact-id` | int | conversations create | Contact for new conversation |
| `--query` | string | search commands | Search query string |
| `--status` | string | conversations list, toggle-status | Conversation status filter/value |
| `--priority` | string | toggle-priority | Priority value |
| `--labels` | string | labels set | Comma-separated label list |
| `--payload` | string | filter commands | JSON filter body |
| `--content` | string | messages create | Message text |
| `--message-type` | string | messages create | outgoing, incoming, or activity |
| `--private` | bool | messages create | Internal note flag |
| `--name` | string | contacts/inboxes create/update | Display name |
| `--email` | string | contacts create/update | Email address |
| `--phone` | string | contacts create/update | Phone number |
| `--channel` | string | inboxes create | JSON channel configuration |
| `--agent-id` | int | assignments create | Agent to assign |
| `--team-id` | int | assignments create | Team to assign |
| `--agent-ids` | string | inbox members add/update/delete | Comma-separated agent IDs |
| `--agent-bot-id` | int | inboxes agent-bot set | Agent bot to assign |
| `--base-id` | int | contacts merge | Base contact ID |
| `--merge-id` | int | contacts merge | Contact ID to merge into base |
| `--enable-auto-assignment` | bool | inboxes update | Toggle auto-assignment |

### Labels Convention

The `--labels` flag accepts comma-separated values. The Chatwoot labels
endpoint is a **replace-all** operation (not additive). The `set` verb name
reflects this to prevent accidental data loss. Example:

```bash
chatwoot application contacts labels set --id 42 --labels "vip,enterprise,priority"
```

### Filter Convention

The `--payload` flag accepts a raw JSON string matching the Chatwoot filter
format. Example:

```bash
chatwoot application contacts filter --payload '[{"attribute_key":"email","filter_operator":"contains","values":["@example.com"]}]'
```

## Existing Code Changes

### `internal/chatwoot/application/conversations.go`

Update `ListConversations` to return pagination metadata:

```go
// Before:
func (c *Client) ListConversations(ctx, opts) ([]Conversation, error)

// After:
func (c *Client) ListConversations(ctx, opts) ([]Conversation, *chatwoot.Pagination, error)
```

Parse the `meta` field from the Chatwoot API response to extract pagination.
Update existing tests and the CLI `conversations list` caller (if any) to
handle the new signature.

### `internal/cli/application/application.go`

Add imports for the 4 new subpackages and register their `Cmd` vars.

### `internal/cli/application/profile.go`

Replace inlined `resolveContext`, `resolveAuth`, `writeError`, `prettyFromRoot`
with `cmdutil.ResolveContext`, `cmdutil.ResolveAuth`, `cmdutil.WriteError`,
`cmdutil.Pretty`.

### `internal/cli/context.go`

Functions that moved to cmdutil should be removed to avoid duplication. If
`ResolveProfileName` is still referenced, keep it; otherwise remove.

## Testing Strategy

### API Client Tests

Each new resource file in `internal/chatwoot/application/` gets a test file
using `httptest.NewServer`. Tests verify:

- Correct URL path (including account ID interpolation)
- Correct HTTP method
- Request body serialization (for create/update/filter)
- Response deserialization into typed models
- Query parameter construction (page, status, inbox_id, etc.)

One test per method minimum. Filter and search methods get additional tests
for parameter encoding.

### CLI Command Tests

Each command package gets a test file following the Sprint B pattern:

- `testroot_test.go` wires a synthetic root with persistent flags
  (`--pretty`, `--profile`, `--base-url`, `--account-id`)
- Env vars (`CHATWOOT_BASE_URL`, `CHATWOOT_ACCESS_TOKEN`) bypass
  config/keychain
- `httptest.NewServer` mocks the API
- Tests verify: JSON envelope on stdout (`ok: true`), correct data shape,
  correct HTTP method/path/body sent to server

**Coverage priority:**

- Every list command: verify pagination flags flow through
- Every get command: verify ID flag maps to URL path
- Every create/update: verify flag-to-request-body mapping
- Every delete: verify correct method + path
- `contacts merge`: verify both IDs in request body
- `contacts labels set` / `conversations labels set`: verify replace-all semantics
- `messages create`: verify `--conversation-id` in path and `--content` in body
- At least one `--all` auto-pagination test (on `contacts list` or
  `conversations list`)

**Not tested at CLI layer:** Retry logic, auth header injection, error
mapping — covered by transport layer tests from Sprint A.

### cmdutil Tests

- `resolveContextFromPath` tests (migrated from `cli/context_test.go`)
- `PaginationFlags` registration and reading
- `Pretty` flag reading

## Exit Criteria

- `go test ./...` passes with all tests green
- `go vet ./...` clean
- `go build ./cmd/chatwoot/` succeeds
- `chatwoot application contacts list` returns JSON envelope with pagination
- `chatwoot application conversations list --status open` filters correctly
- `chatwoot application messages create --conversation-id N --content "..."` sends POST
- `chatwoot application inboxes list` returns inbox data
- `--all` flag auto-paginates on at least contacts and conversations list
- `--labels` comma-separated values work on label set commands
- `--payload` JSON filter works on filter commands
- All commands produce valid JSON envelopes on stdout
- All commands use the cmdutil pipeline (no inlined helpers)
- Error envelopes are consistent across all commands
- Core support workflows from the use-case doc are executable end-to-end

## Out of Scope

- Contacts: import, export, notes, avatar delete, custom attributes destroy
- Conversations: search (source-derived, no documented API endpoint — use
  filter instead), toggle-typing-status, custom-attributes, unread, transcript,
  mute/unmute, update-last-seen, attachments, inbox-assistant, participants,
  draft-messages, direct-uploads, reporting-events
- Inboxes: delete, assignable-agents, campaigns, avatar delete,
  sync-templates, health, register-webhook, csat-template, assignment-policy
- Messages: update, translate, retry
- Interactive prompts
- `--output` format flag (JSON only)

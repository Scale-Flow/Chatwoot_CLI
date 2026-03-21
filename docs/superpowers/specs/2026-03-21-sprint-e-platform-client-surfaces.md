# Sprint E: Platform and Client API Surfaces (Layer 9)

## Scope

Sprint E implements the complete Platform and Client API command surfaces — the two non-application API families. This covers 29 CLI commands backed by 27 API methods across 7 resource families.

**Not in scope:** Survey API (no endpoints exist in the Chatwoot API spec), experimental/source-derived commands (Sprint F).

## Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Client API auth model | Flags/env only, not stored in profiles | Client API contacts are ephemeral per-session. Explicit flags are safer for AI agents managing multiple contacts concurrently. |
| Platform `--account-id` | Reuse global flag from profile resolution | Account-scoped platform commands use `rctx.AccountID` like application commands. No reason to force explicit IDs. |
| SSO login endpoint | Include | One trivial GET method. Completes the users resource family. |
| toggle_typing / update_last_seen | Include | Trivial to implement. Complete Client API surface avoids gaps. |
| Sprint structure | Platform first, Client second | Platform reuses existing token-auth pattern. Client introduces new identifier-based auth — better to isolate that as a distinct phase. |

## Phase 1: Platform API

### Auth Model

Platform commands use `credentials.ModePlatform` which resolves `CHATWOOT_PLATFORM_TOKEN` env var or stored platform token from keychain/file. The token is sent as the `api_access_token` header, same as application commands.

Handler pipeline:
```go
rctx, err := cmdutil.ResolveContext(cmd)
tokenAuth, err := cmdutil.ResolveAuth(rctx.ProfileName, credentials.ModePlatform)
transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
client := platform.NewClient(transport)
```

### Platform Client Design

The existing `platform.Client` has no `accountID` field — unlike `application.Client`. This is correct: platform resources split into global (users, agent-bots) and account-scoped (accounts, account-users), with account IDs passed as method parameters.

```go
type Client struct {
    transport *chatwoot.Client
}
```

### Model Types

Add to `internal/chatwoot/platform/models.go`:

```go
// Existing types: Account, CreateAccountOpts

// --- New types ---

type UpdateAccountOpts struct {
    Name *string `json:"name,omitempty"`
}

type User struct {
    ID            int    `json:"id"`
    Name          string `json:"name"`
    Email         string `json:"email"`
    Type          string `json:"type,omitempty"`
    Confirmed     bool   `json:"confirmed,omitempty"`
    CustomAttributes any `json:"custom_attributes,omitempty"`
}

type CreateUserOpts struct {
    Name             string `json:"name"`
    Email            string `json:"email"`
    Password         string `json:"password"`
    CustomAttributes any    `json:"custom_attributes,omitempty"`
}

type UpdateUserOpts struct {
    Name             *string `json:"name,omitempty"`
    Email            *string `json:"email,omitempty"`
    Password         *string `json:"password,omitempty"`
    CustomAttributes any     `json:"custom_attributes,omitempty"`
}

type SSOLink struct {
    URL string `json:"url"`
}

type AccountUser struct {
    AccountID int    `json:"account_id"`
    UserID    int    `json:"user_id"`
    Role      string `json:"role,omitempty"`
}

type CreateAccountUserOpts struct {
    UserID int    `json:"user_id"`
    Role   string `json:"role,omitempty"`
}

type AgentBot struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    OutgoingURL string `json:"outgoing_url,omitempty"`
    BotType     string `json:"bot_type,omitempty"`
    BotConfig   any    `json:"bot_config,omitempty"`
    AccessToken string `json:"access_token,omitempty"`
}

type CreateAgentBotOpts struct {
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    OutgoingURL string `json:"outgoing_url,omitempty"`
    BotType     string `json:"bot_type,omitempty"`
    BotConfig   any    `json:"bot_config,omitempty"`
}

type UpdateAgentBotOpts struct {
    Name        *string `json:"name,omitempty"`
    Description *string `json:"description,omitempty"`
    OutgoingURL *string `json:"outgoing_url,omitempty"`
    BotType     *string `json:"bot_type,omitempty"`
    BotConfig   any     `json:"bot_config,omitempty"`
}
```

### API Methods

**Accounts** — modify `internal/chatwoot/platform/accounts.go` (2 new methods, 2 existing):
- Existing: `GetAccount(ctx, id)`, `CreateAccount(ctx, CreateAccountOpts)`
- `UpdateAccount(ctx, id int, UpdateAccountOpts) (*Account, error)` — PATCH `/platform/api/v1/accounts/{id}`
- `DeleteAccount(ctx, id int) error` — DELETE `/platform/api/v1/accounts/{id}`

**Account Users** — new `internal/chatwoot/platform/account_users.go` (3 methods):
- `ListAccountUsers(ctx, accountID int) ([]AccountUser, error)` — GET `/platform/api/v1/accounts/{id}/account_users`
- `CreateAccountUser(ctx, accountID int, CreateAccountUserOpts) error` — POST `/platform/api/v1/accounts/{id}/account_users`
- `DeleteAccountUser(ctx, accountID int, userID int) error` — DELETE `/platform/api/v1/accounts/{id}/account_users` (body: `{"user_id": N}`)

**Agent Bots** — new `internal/chatwoot/platform/agent_bots.go` (5 methods):
- `ListAgentBots(ctx) ([]AgentBot, error)` — GET `/platform/api/v1/agent_bots`
- `GetAgentBot(ctx, id int) (*AgentBot, error)` — GET `/platform/api/v1/agent_bots/{id}`
- `CreateAgentBot(ctx, CreateAgentBotOpts) (*AgentBot, error)` — POST `/platform/api/v1/agent_bots`
- `UpdateAgentBot(ctx, id int, UpdateAgentBotOpts) (*AgentBot, error)` — PATCH `/platform/api/v1/agent_bots/{id}`
- `DeleteAgentBot(ctx, id int) error` — DELETE `/platform/api/v1/agent_bots/{id}`

**Users** — new `internal/chatwoot/platform/users.go` (5 methods):
- `CreateUser(ctx, CreateUserOpts) (*User, error)` — POST `/platform/api/v1/users`
- `GetUser(ctx, id int) (*User, error)` — GET `/platform/api/v1/users/{id}`
- `UpdateUser(ctx, id int, UpdateUserOpts) (*User, error)` — PATCH `/platform/api/v1/users/{id}`
- `DeleteUser(ctx, id int) error` — DELETE `/platform/api/v1/users/{id}`
- `GetUserSSOLink(ctx, id int) (*SSOLink, error)` — GET `/platform/api/v1/users/{id}/login`

### Response Decode Patterns

Most platform endpoints return direct JSON objects — decode into the target type. Key exceptions to verify during implementation:
- `ListAccountUsers` and `ListAgentBots` — check whether response is a plain array or uses a payload wrapper. Default assumption: plain array (consistent with Chatwoot's pattern for small collections).
- `CreateAccountUser` and `DeleteAccountUser` — may return empty body or simple acknowledgment. Decode with `nil` if no useful response body.

### Platform CLI Commands

#### Package Layout

```
internal/cli/platform/
  platform.go                    # Cmd group
  accounts/
    accounts.go                  # Cmd subgroup
    create.go, get.go, update.go, delete.go
    accounts_test.go, testroot_test.go
  accountusers/
    accountusers.go              # Cmd subgroup (Use: "account-users")
    list.go, create.go, delete.go
    accountusers_test.go, testroot_test.go
  agentbots/
    agentbots.go                 # Cmd subgroup (Use: "agent-bots")
    list.go, get.go, create.go, update.go, delete.go
    agentbots_test.go, testroot_test.go
  users/
    users.go                     # Cmd subgroup
    create.go, get.go, update.go, delete.go, login.go
    users_test.go, testroot_test.go
```

#### Flag Design

**Account-scoped commands** (`accounts get/update/delete`, all `account-users`):
Use `rctx.AccountID` from the global `--account-id` flag / profile. For `accounts create`, no account ID is needed (creating a new one).

**Global resource commands** (`users`, `agent-bots`):
Use `--id` for specific resource operations. No account ID needed.

**Flags per command:**

| Command | Required Flags | Optional Flags |
|---------|---------------|----------------|
| `accounts create` | `--name` | — |
| `accounts get` | (uses `--account-id` from context) | — |
| `accounts update` | (uses `--account-id` from context) | `--name` (at least one) |
| `accounts delete` | (uses `--account-id` from context) | — |
| `account-users list` | (uses `--account-id` from context) | — |
| `account-users create` | `--user-id` | `--role` |
| `account-users delete` | `--user-id` | — |
| `agent-bots list` | — | — |
| `agent-bots get` | `--id` | — |
| `agent-bots create` | `--name` | `--description`, `--outgoing-url`, `--bot-type`, `--bot-config` (JSON) |
| `agent-bots update` | `--id` | `--name`, `--description`, `--outgoing-url`, `--bot-type`, `--bot-config` (at least one) |
| `agent-bots delete` | `--id` | — |
| `users create` | `--name`, `--email`, `--password` | `--custom-attributes` (JSON) |
| `users get` | `--id` | — |
| `users update` | `--id` | `--name`, `--email`, `--password`, `--custom-attributes` (at least one) |
| `users delete` | `--id` | — |
| `users login` | `--id` | — |

#### testroot_test.go Pattern

Platform test roots wire the `platform` group under the root:
```go
package <pkgname>

import "github.com/spf13/cobra"

func init() {
    root := &cobra.Command{Use: "chatwoot", SilenceUsage: true, SilenceErrors: true}
    root.PersistentFlags().Bool("pretty", false, "Indent JSON output")
    root.PersistentFlags().String("profile", "", "Select named profile")
    root.PersistentFlags().String("base-url", "", "Override base URL")
    root.PersistentFlags().Int("account-id", 0, "Override account ID")
    platformCmd := &cobra.Command{Use: "platform"}
    platformCmd.AddCommand(Cmd)
    root.AddCommand(platformCmd)
}
```

#### Registration

`internal/cli/platform/platform.go`:
```go
package platform

import (
    "github.com/chatwoot/chatwoot-cli/internal/cli/platform/accounts"
    "github.com/chatwoot/chatwoot-cli/internal/cli/platform/accountusers"
    "github.com/chatwoot/chatwoot-cli/internal/cli/platform/agentbots"
    "github.com/chatwoot/chatwoot-cli/internal/cli/platform/users"
    "github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
    Use:   "platform",
    Short: "Platform API commands (self-hosted admin)",
}

func init() {
    Cmd.AddCommand(accounts.Cmd)
    Cmd.AddCommand(accountusers.Cmd)
    Cmd.AddCommand(agentbots.Cmd)
    Cmd.AddCommand(users.Cmd)
}
```

Then add to `internal/cli/root.go`:
```go
rootCmd.AddCommand(cliplatform.Cmd)
```

---

## Phase 2: Client API

### Auth Model

Client API commands do not use token auth. They use identifier-based auth passed via flags or env vars:

- `--inbox-id` flag or `CHATWOOT_INBOX_IDENTIFIER` env var (required for all client commands)
- `--contact-id` flag or `CHATWOOT_CONTACT_IDENTIFIER` env var (required for most commands, optional for `contacts create`)

These are persistent flags on the `client` command group, inherited by all subcommands.

#### New cmdutil Helper

Add `ResolveClientAuth` to `internal/cli/cmdutil/context.go`:

```go
type ClientAuthContext struct {
    InboxIdentifier   string
    ContactIdentifier string
}

func ResolveClientAuth(cmd *cobra.Command) (*ClientAuthContext, error) {
    inboxID, _ := cmd.Flags().GetString("inbox-id")
    if inboxID == "" {
        inboxID = os.Getenv("CHATWOOT_INBOX_IDENTIFIER")
    }
    if inboxID == "" {
        return nil, fmt.Errorf("--inbox-id flag or CHATWOOT_INBOX_IDENTIFIER env var required")
    }
    contactID, _ := cmd.Flags().GetString("contact-id")
    if contactID == "" {
        contactID = os.Getenv("CHATWOOT_CONTACT_IDENTIFIER")
    }
    return &ClientAuthContext{InboxIdentifier: inboxID, ContactIdentifier: contactID}, nil
}
```

#### Handler Pipeline

```go
rctx, err := cmdutil.ResolveContext(cmd)
clientAuth, err := cmdutil.ResolveClientAuth(cmd)
transport := chatwoot.NewClient(rctx.BaseURL, "", "")
client := clientapi.NewClient(transport, clientAuth.InboxIdentifier)
// ... call client methods with clientAuth.ContactIdentifier
```

### Client API Client Design

The existing `clientapi.Client` stores `inboxIdentifier`. The `contactIdentifier` is passed per-method because a single inbox can have many contacts.

```go
type Client struct {
    transport       *chatwoot.Client
    inboxIdentifier string
}
```

### Model Types

Add to `internal/chatwoot/clientapi/models.go`:

```go
// Existing types: Contact, CreateContactOpts

// --- New types ---

type UpdateContactOpts struct {
    Name  *string `json:"name,omitempty"`
    Email *string `json:"email,omitempty"`
    Phone *string `json:"phone_number,omitempty"`
}

type Conversation struct {
    ID                int    `json:"id"`
    InboxID           int    `json:"inbox_id,omitempty"`
    Status            string `json:"status,omitempty"`
    AgentID           int    `json:"agent_id,omitempty"`
    ContactLastSeenAt string `json:"contact_last_seen_at,omitempty"`
}

type Message struct {
    ID          int    `json:"id"`
    Content     string `json:"content,omitempty"`
    MessageType string `json:"message_type,omitempty"`
    ContentType string `json:"content_type,omitempty"`
    CreatedAt   string `json:"created_at,omitempty"`
}

type CreateMessageOpts struct {
    Content     string `json:"content"`
    MessageType string `json:"message_type,omitempty"`
}

type UpdateMessageOpts struct {
    Content string `json:"content"`
}

type ToggleTypingOpts struct {
    TypingStatus string `json:"typing_status"`
}
```

### API Methods

**Contacts** — modify `internal/chatwoot/clientapi/contacts.go` (1 new, 2 existing):
- Existing: `CreateContact(ctx, CreateContactOpts)`, `GetContact(ctx, contactIdentifier)`
- `UpdateContact(ctx, contactIdentifier string, UpdateContactOpts) (*Contact, error)` — PATCH `/public/api/v1/inboxes/{inbox}/contacts/{contact}`

**Conversations** — new `internal/chatwoot/clientapi/conversations.go` (6 methods):
- `ListConversations(ctx, contactIdentifier string) ([]Conversation, error)` — GET `.../contacts/{contact}/conversations`
- `GetConversation(ctx, contactIdentifier string, conversationID int) (*Conversation, error)` — GET `.../conversations/{id}`
- `CreateConversation(ctx, contactIdentifier string) (*Conversation, error)` — POST `.../contacts/{contact}/conversations`
- `ToggleStatus(ctx, contactIdentifier string, conversationID int) (*Conversation, error)` — POST `.../conversations/{id}/toggle_status`
- `ToggleTyping(ctx, contactIdentifier string, conversationID int, ToggleTypingOpts) error` — POST `.../conversations/{id}/toggle_typing`
- `UpdateLastSeen(ctx, contactIdentifier string, conversationID int) error` — POST `.../conversations/{id}/update_last_seen`

**Messages** — new `internal/chatwoot/clientapi/messages.go` (3 methods):
- `ListMessages(ctx, contactIdentifier string, conversationID int) ([]Message, error)` — GET `.../conversations/{id}/messages`
- `CreateMessage(ctx, contactIdentifier string, conversationID int, CreateMessageOpts) (*Message, error)` — POST `.../conversations/{id}/messages`
- `UpdateMessage(ctx, contactIdentifier string, conversationID int, messageID int, UpdateMessageOpts) (*Message, error)` — PATCH `.../conversations/{id}/messages/{mid}`

All paths are anchored at `/public/api/v1/inboxes/{inbox_identifier}/contacts/{contact_identifier}/...`.

### Client CLI Commands

#### Package Layout

```
internal/cli/client/
  client.go                       # Cmd group with persistent --inbox-id, --contact-id
  contacts/
    contacts.go                   # Cmd subgroup
    create.go, get.go, update.go
    contacts_test.go, testroot_test.go
  conversations/
    conversations.go              # Cmd subgroup
    list.go, get.go, create.go, toggle_status.go, toggle_typing.go, update_last_seen.go
    conversations_test.go, testroot_test.go
  messages/
    messages.go                   # Cmd subgroup
    list.go, create.go, update.go
    messages_test.go, testroot_test.go
```

#### Flag Design

**Persistent flags on `client` Cmd:**
- `--inbox-id` (string) — env fallback `CHATWOOT_INBOX_IDENTIFIER`
- `--contact-id` (string) — env fallback `CHATWOOT_CONTACT_IDENTIFIER`

**Flags per command:**

| Command | Required Flags | Optional Flags |
|---------|---------------|----------------|
| `contacts create` | — | `--name`, `--email`, `--phone` |
| `contacts get` | `--contact-id` (persistent) | — |
| `contacts update` | `--contact-id` (persistent) | `--name`, `--email`, `--phone` (at least one) |
| `conversations list` | `--contact-id` (persistent) | — |
| `conversations get` | `--contact-id` (persistent), `--conversation-id` | — |
| `conversations create` | `--contact-id` (persistent) | — |
| `conversations toggle-status` | `--contact-id` (persistent), `--conversation-id` | — |
| `conversations toggle-typing` | `--contact-id` (persistent), `--conversation-id`, `--status` ("on"/"off") | — |
| `conversations update-last-seen` | `--contact-id` (persistent), `--conversation-id` | — |
| `messages list` | `--contact-id` (persistent), `--conversation-id` | — |
| `messages create` | `--contact-id` (persistent), `--conversation-id`, `--content` | `--type` |
| `messages update` | `--contact-id` (persistent), `--conversation-id`, `--message-id`, `--content` | — |

#### testroot_test.go Pattern

Client test roots wire the `client` group under root and include the persistent flags:
```go
package <pkgname>

import "github.com/spf13/cobra"

func init() {
    root := &cobra.Command{Use: "chatwoot", SilenceUsage: true, SilenceErrors: true}
    root.PersistentFlags().Bool("pretty", false, "Indent JSON output")
    root.PersistentFlags().String("profile", "", "Select named profile")
    root.PersistentFlags().String("base-url", "", "Override base URL")
    root.PersistentFlags().Int("account-id", 0, "Override account ID")
    clientCmd := &cobra.Command{Use: "client"}
    clientCmd.PersistentFlags().String("inbox-id", "", "Inbox identifier")
    clientCmd.PersistentFlags().String("contact-id", "", "Contact identifier")
    clientCmd.AddCommand(Cmd)
    root.AddCommand(clientCmd)
}
```

#### Registration

`internal/cli/client/client.go`:
```go
package client

import (
    "github.com/chatwoot/chatwoot-cli/internal/cli/client/contacts"
    "github.com/chatwoot/chatwoot-cli/internal/cli/client/conversations"
    "github.com/chatwoot/chatwoot-cli/internal/cli/client/messages"
    "github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
    Use:   "client",
    Short: "Client API commands (public end-user)",
}

func init() {
    Cmd.PersistentFlags().String("inbox-id", "", "Inbox identifier")
    Cmd.PersistentFlags().String("contact-id", "", "Contact identifier")
    Cmd.AddCommand(contacts.Cmd)
    Cmd.AddCommand(conversations.Cmd)
    Cmd.AddCommand(messages.Cmd)
}
```

Then add to `internal/cli/root.go`:
```go
rootCmd.AddCommand(cliclient.Cmd)
```

---

## Task Breakdown

### Phase 1: Platform

1. **Add platform model types** — modify `platform/models.go`
2. **Complete accounts API client** — add update + delete to `platform/accounts.go` + tests
3. **Account-users API client** — new `platform/account_users.go` + tests
4. **Agent-bots API client** — new `platform/agent_bots.go` + tests
5. **Users API client** — new `platform/users.go` + tests
6. **Platform CLI commands** — 4 packages: accounts, accountusers, agentbots, users
7. **Register platform group** — add `platform.go`, update `root.go`
8. **Platform verification** — full test suite, help output, build

### Phase 2: Client

9. **Add client model types** — modify `clientapi/models.go`
10. **Add ResolveClientAuth helper** — modify `cmdutil/context.go`
11. **Complete contacts API client** — add update to `clientapi/contacts.go` + tests
12. **Conversations API client** — new `clientapi/conversations.go` + tests
13. **Messages API client** — new `clientapi/messages.go` + tests
14. **Client CLI commands** — 3 packages: contacts, conversations, messages
15. **Register client group, full verification** — add `client.go`, update `root.go`, verify everything

---

## Exit Criteria

- `go test ./...` passes with all tests green
- `go vet ./...` clean
- `go build ./cmd/chatwoot/` succeeds
- `chatwoot platform --help` shows accounts, account-users, agent-bots, users
- `chatwoot platform accounts create --name X` produces JSON envelope
- `chatwoot platform users login --id 1` produces JSON envelope with URL
- `chatwoot client --help` shows contacts, conversations, messages
- `chatwoot client contacts create --inbox-id X` produces JSON envelope
- `chatwoot client messages create --inbox-id X --contact-id Y --conversation-id 1 --content "hello"` produces JSON envelope
- Client commands fail with clear error when `--inbox-id` is missing
- Platform commands use `CHATWOOT_PLATFORM_TOKEN` / stored platform credential
- Existing application and auth commands still work
- All commands produce valid JSON envelopes on stdout
- All commands use the cmdutil pipeline (no business logic in handlers)

---

## Reference Patterns

Subagents implementing this plan should reference:
- **Platform API client pattern:** `internal/chatwoot/platform/accounts.go` (existing)
- **Client API client pattern:** `internal/chatwoot/clientapi/contacts.go` (existing)
- **Application CLI handler pattern:** `internal/cli/application/contacts/list.go`
- **CLI test pattern:** `internal/cli/application/contacts/contacts_test.go`
- **Members pattern (for account-users):** `internal/cli/application/inboxes/members.go`
- **testroot_test.go pattern:** `internal/cli/application/contacts/testroot_test.go`

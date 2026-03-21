# Sprint E: Platform and Client API Surfaces (Layer 9) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the complete Platform and Client API command surfaces — 7 resource families with ~27 API methods and ~29 CLI commands enabling self-hosted admin operations and end-user chat flows.

**Architecture:** Two phases. Phase 1 builds Platform API (token auth, same pattern as application). Phase 2 builds Client API (identifier-based auth, new `ResolveClientAuth` helper). Each phase follows models → API client → CLI commands → registration.

**Tech Stack:** Go 1.26.1, `spf13/cobra` (CLI), `spf13/viper` (config), `zalando/go-keyring` (keychain), `log/slog` (diagnostics), `encoding/json` (serialization), `net/http/httptest` (testing)

**Spec:** `docs/superpowers/specs/2026-03-21-sprint-e-platform-client-surfaces.md`

**Reference patterns for subagents:**
- Platform API client: `internal/chatwoot/platform/accounts.go` (existing GET/POST pattern)
- Platform API test: `internal/chatwoot/platform/platform_test.go` (httptest with `pk-test` token, `api_access_token` header)
- Client API client: `internal/chatwoot/clientapi/contacts.go` (identifier-based paths, no auth header)
- Client API test: `internal/chatwoot/clientapi/clientapi_test.go` (httptest with empty token)
- CLI handler: `internal/cli/application/contacts/list.go` (cmdutil pipeline)
- CLI test: `internal/cli/application/contacts/contacts_test.go` (httptest, env vars, JSON verification)
- testroot_test.go: `internal/cli/application/contacts/testroot_test.go` (same boilerplate per package)

---

## File Structure

### Platform — Model Types

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/platform/models.go` | MODIFY: add User, AccountUser, AgentBot, SSOLink types and all opts types |

### Platform — API Clients

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/platform/accounts.go` | MODIFY: add UpdateAccount, DeleteAccount (2 new methods to existing 2) |
| `internal/chatwoot/platform/platform_test.go` | MODIFY: add tests for new account methods |
| `internal/chatwoot/platform/account_users.go` | NEW: 3 account-user API methods |
| `internal/chatwoot/platform/account_users_test.go` | NEW: account-user API tests |
| `internal/chatwoot/platform/agent_bots.go` | NEW: 5 agent-bot API methods |
| `internal/chatwoot/platform/agent_bots_test.go` | NEW: agent-bot API tests |
| `internal/chatwoot/platform/users.go` | NEW: 5 user API methods |
| `internal/chatwoot/platform/users_test.go` | NEW: user API tests |

### Platform — CLI Commands

| File | Responsibility |
|------|---------------|
| `internal/cli/platform/platform.go` | NEW: Cmd group, registers 4 subpackages |
| `internal/cli/platform/accounts/accounts.go` | NEW: Cmd group |
| `internal/cli/platform/accounts/create.go` | NEW: accounts create |
| `internal/cli/platform/accounts/get.go` | NEW: accounts get |
| `internal/cli/platform/accounts/update.go` | NEW: accounts update |
| `internal/cli/platform/accounts/delete.go` | NEW: accounts delete |
| `internal/cli/platform/accounts/accounts_test.go` | NEW: tests |
| `internal/cli/platform/accounts/testroot_test.go` | NEW: test helper |
| `internal/cli/platform/accountusers/accountusers.go` | NEW: Cmd group (Use: "account-users") |
| `internal/cli/platform/accountusers/list.go` | NEW: account-users list |
| `internal/cli/platform/accountusers/create.go` | NEW: account-users create |
| `internal/cli/platform/accountusers/delete.go` | NEW: account-users delete |
| `internal/cli/platform/accountusers/accountusers_test.go` | NEW: tests |
| `internal/cli/platform/accountusers/testroot_test.go` | NEW: test helper |
| `internal/cli/platform/agentbots/agentbots.go` | NEW: Cmd group (Use: "agent-bots") |
| `internal/cli/platform/agentbots/list.go` | NEW: agent-bots list |
| `internal/cli/platform/agentbots/get.go` | NEW: agent-bots get |
| `internal/cli/platform/agentbots/create.go` | NEW: agent-bots create |
| `internal/cli/platform/agentbots/update.go` | NEW: agent-bots update |
| `internal/cli/platform/agentbots/delete.go` | NEW: agent-bots delete |
| `internal/cli/platform/agentbots/agentbots_test.go` | NEW: tests |
| `internal/cli/platform/agentbots/testroot_test.go` | NEW: test helper |
| `internal/cli/platform/users/users.go` | NEW: Cmd group |
| `internal/cli/platform/users/create.go` | NEW: users create |
| `internal/cli/platform/users/get.go` | NEW: users get |
| `internal/cli/platform/users/update.go` | NEW: users update |
| `internal/cli/platform/users/delete.go` | NEW: users delete |
| `internal/cli/platform/users/login.go` | NEW: users login |
| `internal/cli/platform/users/users_test.go` | NEW: tests |
| `internal/cli/platform/users/testroot_test.go` | NEW: test helper |

### Client — Model Types and Infrastructure

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/clientapi/models.go` | MODIFY: add Conversation, Message types and opts types |
| `internal/cli/cmdutil/context.go` | MODIFY: add ResolveClientAuth helper |

### Client — API Clients

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/clientapi/contacts.go` | MODIFY: add UpdateContact (1 new method to existing 2) |
| `internal/chatwoot/clientapi/clientapi_test.go` | MODIFY: add test for UpdateContact |
| `internal/chatwoot/clientapi/conversations.go` | NEW: 6 conversation API methods |
| `internal/chatwoot/clientapi/conversations_test.go` | NEW: conversation API tests |
| `internal/chatwoot/clientapi/messages.go` | NEW: 3 message API methods |
| `internal/chatwoot/clientapi/messages_test.go` | NEW: message API tests |

### Client — CLI Commands

| File | Responsibility |
|------|---------------|
| `internal/cli/client/client.go` | NEW: Cmd group with persistent --inbox-id, --contact-id flags |
| `internal/cli/client/contacts/contacts.go` | NEW: Cmd group |
| `internal/cli/client/contacts/create.go` | NEW: contacts create |
| `internal/cli/client/contacts/get.go` | NEW: contacts get |
| `internal/cli/client/contacts/update.go` | NEW: contacts update |
| `internal/cli/client/contacts/contacts_test.go` | NEW: tests |
| `internal/cli/client/contacts/testroot_test.go` | NEW: test helper |
| `internal/cli/client/conversations/conversations.go` | NEW: Cmd group |
| `internal/cli/client/conversations/list.go` | NEW: conversations list |
| `internal/cli/client/conversations/get.go` | NEW: conversations get |
| `internal/cli/client/conversations/create.go` | NEW: conversations create |
| `internal/cli/client/conversations/toggle_status.go` | NEW: conversations toggle-status |
| `internal/cli/client/conversations/toggle_typing.go` | NEW: conversations toggle-typing |
| `internal/cli/client/conversations/update_last_seen.go` | NEW: conversations update-last-seen |
| `internal/cli/client/conversations/conversations_test.go` | NEW: tests |
| `internal/cli/client/conversations/testroot_test.go` | NEW: test helper |
| `internal/cli/client/messages/messages.go` | NEW: Cmd group |
| `internal/cli/client/messages/list.go` | NEW: messages list |
| `internal/cli/client/messages/create.go` | NEW: messages create |
| `internal/cli/client/messages/update.go` | NEW: messages update |
| `internal/cli/client/messages/messages_test.go` | NEW: tests |
| `internal/cli/client/messages/testroot_test.go` | NEW: test helper |

### Registration

| File | Responsibility |
|------|---------------|
| `internal/cli/root.go` | MODIFY: register platform and client command groups |

---

## Task 1: Add Platform Model Types

**Files:**
- Modify: `internal/chatwoot/platform/models.go`

- [ ] **Step 1: Add all new types to models.go**

Add the following types after the existing `CreateAccountOpts`:

```go
type UpdateAccountOpts struct {
	Name *string `json:"name,omitempty"`
}

// User represents a Chatwoot platform user.
type User struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Type             string `json:"type,omitempty"`
	Confirmed        bool   `json:"confirmed,omitempty"`
	CustomAttributes any    `json:"custom_attributes,omitempty"`
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

// SSOLink holds the SSO login URL for a user.
type SSOLink struct {
	URL string `json:"url"`
}

// AccountUser represents an account-user association.
type AccountUser struct {
	AccountID int    `json:"account_id"`
	UserID    int    `json:"user_id"`
	Role      string `json:"role,omitempty"`
}

type CreateAccountUserOpts struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role,omitempty"`
}

// AgentBot represents a platform-scoped agent bot.
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

- [ ] **Step 2: Verify build**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./cmd/chatwoot/`
Expected: Success

- [ ] **Step 3: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/platform/models.go
git commit -m "feat(platform): add Sprint E model types for users, account-users, agent-bots"
```

---

## Task 2: Complete Accounts API Client

**Files:**
- Modify: `internal/chatwoot/platform/accounts.go`
- Modify: `internal/chatwoot/platform/platform_test.go`

- [ ] **Step 1: Add UpdateAccount and DeleteAccount to accounts.go**

Add after existing `CreateAccount`:

**HTTP methods and paths:**
- `UpdateAccount(ctx, id int, UpdateAccountOpts) (*Account, error)` — PATCH `/platform/api/v1/accounts/{id}` — marshal opts, decode into `Account`
- `DeleteAccount(ctx, id int) error` — DELETE `/platform/api/v1/accounts/{id}` — decode with `nil`

Follow the exact pattern from the existing `CreateAccount` and `GetAccount` methods in the same file.

- [ ] **Step 2: Add tests for UpdateAccount and DeleteAccount**

Add to `platform_test.go`:
- `TestUpdateAccount` — verify PATCH method, path `/platform/api/v1/accounts/1`, request body has name, auth header `api_access_token` present
- `TestDeleteAccount` — verify DELETE method, path `/platform/api/v1/accounts/5`

- [ ] **Step 3: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/platform/ -v`
Expected: All PASS (4 tests: existing 2 + new 2)

- [ ] **Step 4: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/platform/accounts.go internal/chatwoot/platform/platform_test.go
git commit -m "feat(platform): add update and delete account methods"
```

---

## Task 3: Account Users API Client

**Files:**
- Create: `internal/chatwoot/platform/account_users.go`
- Create: `internal/chatwoot/platform/account_users_test.go`

- [ ] **Step 1: Write failing test for ListAccountUsers**

Create `account_users_test.go`. The httptest server returns a plain JSON array `[{"account_id": 1, "user_id": 10, "role": "administrator"}]`. Verify GET method and path `/platform/api/v1/accounts/1/account_users`.

Use the existing test pattern: `chatwoot.NewClient(srv.URL, "pk-test", "api_access_token")`, `NewClient(transport)`.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/platform/ -v -run TestListAccountUsers`
Expected: FAIL — `ListAccountUsers` not defined

- [ ] **Step 3: Implement all 3 account-user methods**

Create `account_users.go`:

**HTTP methods and paths:**
- `ListAccountUsers(ctx, accountID int) ([]AccountUser, error)` — GET `/platform/api/v1/accounts/{id}/account_users` — decode into `[]AccountUser`
- `CreateAccountUser(ctx, accountID int, CreateAccountUserOpts) error` — POST `/platform/api/v1/accounts/{id}/account_users` — marshal opts, decode with `nil`
- `DeleteAccountUser(ctx, accountID int, userID int) error` — DELETE `/platform/api/v1/accounts/{id}/account_users` — body: `{"user_id": N}`, decode with `nil`

Note: `DeleteAccountUser` sends the user_id in the request body, not in the URL path.

- [ ] **Step 4: Add tests for CreateAccountUser and DeleteAccountUser**

- `TestCreateAccountUser` — verify POST, path, request body has `user_id` and `role`
- `TestDeleteAccountUser` — verify DELETE, path, request body has `user_id`

- [ ] **Step 5: Run all account-user tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/platform/ -v -run "TestListAccountUsers|TestCreateAccountUser|TestDeleteAccountUser"`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/platform/account_users.go internal/chatwoot/platform/account_users_test.go
git commit -m "feat(platform): add account-users API client with 3 methods"
```

---

## Task 4: Agent Bots API Client

**Files:**
- Create: `internal/chatwoot/platform/agent_bots.go`
- Create: `internal/chatwoot/platform/agent_bots_test.go`

- [ ] **Step 1: Write failing test for ListAgentBots**

Create `agent_bots_test.go`. Server returns `[{"id": 1, "name": "Helper Bot"}]`. Verify GET and path `/platform/api/v1/agent_bots` (no account ID — global resource).

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/platform/ -v -run TestListAgentBots`
Expected: FAIL

- [ ] **Step 3: Implement all 5 agent-bot methods**

Create `agent_bots.go`:

**HTTP methods and paths (all global — no account ID):**
- `ListAgentBots(ctx) ([]AgentBot, error)` — GET `/platform/api/v1/agent_bots` — decode into `[]AgentBot`
- `GetAgentBot(ctx, id int) (*AgentBot, error)` — GET `/platform/api/v1/agent_bots/{id}` — decode into `AgentBot`
- `CreateAgentBot(ctx, CreateAgentBotOpts) (*AgentBot, error)` — POST `/platform/api/v1/agent_bots` — marshal opts, decode into `AgentBot`
- `UpdateAgentBot(ctx, id int, UpdateAgentBotOpts) (*AgentBot, error)` — PATCH `/platform/api/v1/agent_bots/{id}` — marshal opts, decode into `AgentBot`
- `DeleteAgentBot(ctx, id int) error` — DELETE `/platform/api/v1/agent_bots/{id}` — decode with `nil`

- [ ] **Step 4: Add tests for GetAgentBot, CreateAgentBot, DeleteAgentBot**

- [ ] **Step 5: Run all agent-bot tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/platform/ -v -run "TestListAgentBots|TestGetAgentBot|TestCreateAgentBot|TestDeleteAgentBot"`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/platform/agent_bots.go internal/chatwoot/platform/agent_bots_test.go
git commit -m "feat(platform): add agent-bots API client with 5 methods"
```

---

## Task 5: Users API Client

**Files:**
- Create: `internal/chatwoot/platform/users.go`
- Create: `internal/chatwoot/platform/users_test.go`

- [ ] **Step 1: Write failing test for CreateUser**

Create `users_test.go`. Server verifies POST to `/platform/api/v1/users`, request body has `name`, `email`, `password`. Returns `{"id": 1, "name": "Alice", "email": "alice@test.com"}`.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/platform/ -v -run TestCreateUser`
Expected: FAIL

- [ ] **Step 3: Implement all 5 user methods**

Create `users.go`:

**HTTP methods and paths (all global — no account ID):**
- `CreateUser(ctx, CreateUserOpts) (*User, error)` — POST `/platform/api/v1/users` — marshal opts, decode into `User`
- `GetUser(ctx, id int) (*User, error)` — GET `/platform/api/v1/users/{id}` — decode into `User`
- `UpdateUser(ctx, id int, UpdateUserOpts) (*User, error)` — PATCH `/platform/api/v1/users/{id}` — marshal opts, decode into `User`
- `DeleteUser(ctx, id int) error` — DELETE `/platform/api/v1/users/{id}` — decode with `nil`
- `GetUserSSOLink(ctx, id int) (*SSOLink, error)` — GET `/platform/api/v1/users/{id}/login` — decode into `SSOLink`

- [ ] **Step 4: Add tests for GetUser, DeleteUser, GetUserSSOLink**

- `TestGetUser` — verify GET, path `/platform/api/v1/users/1`
- `TestDeleteUser` — verify DELETE, path `/platform/api/v1/users/5`
- `TestGetUserSSOLink` — verify GET, path `/platform/api/v1/users/1/login`, response has `url` field

- [ ] **Step 5: Run all user tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/platform/ -v -run "TestCreateUser|TestGetUser|TestDeleteUser|TestGetUserSSOLink"`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/platform/users.go internal/chatwoot/platform/users_test.go
git commit -m "feat(platform): add users API client with 5 methods"
```

---

## Task 6: Platform CLI Commands

**Files:**
- Create: `internal/cli/platform/platform.go`
- Create: all files in `internal/cli/platform/accounts/`
- Create: all files in `internal/cli/platform/accountusers/`
- Create: all files in `internal/cli/platform/agentbots/`
- Create: all files in `internal/cli/platform/users/`

### Platform CLI Handler Pattern

All platform CLI handlers follow this pipeline (different from application — uses `credentials.ModePlatform` and `platform.NewClient` with no accountID):

```go
func runSomething(cmd *cobra.Command, args []string) error {
    rctx, err := cmdutil.ResolveContext(cmd)
    if err != nil {
        return cmdutil.WriteError(cmd, contract.ErrCodeConfig, err.Error())
    }
    tokenAuth, err := cmdutil.ResolveAuth(rctx.ProfileName, credentials.ModePlatform)
    if err != nil {
        return cmdutil.WriteError(cmd, contract.ErrCodeAuth, err.Error())
    }
    transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
    client := platapi.NewClient(transport)
    // ... get flags, call client, write response
}
```

**Standard imports for platform CLI handlers:**
```go
import (
    "context"
    platapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/platform"
    chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
    "github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
    "github.com/chatwoot/chatwoot-cli/internal/contract"
    "github.com/chatwoot/chatwoot-cli/internal/credentials"
    "github.com/spf13/cobra"
)
```

**Platform testroot_test.go pattern** (note: `platformCmd` instead of `appCmd`):
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

**Platform test pattern** (uses `CHATWOOT_PLATFORM_TOKEN` instead of `CHATWOOT_ACCESS_TOKEN`):
```go
func TestSomething(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(/* response */)
    }))
    defer srv.Close()
    t.Setenv("CHATWOOT_BASE_URL", srv.URL)
    t.Setenv("CHATWOOT_PLATFORM_TOKEN", "pk-test")
    var stdout bytes.Buffer
    Cmd.SetOut(&stdout)
    Cmd.SetErr(&bytes.Buffer{})
    Cmd.Root().SetArgs([]string{"platform", "<resource>", "<action>", ...flags})
    err := Cmd.Root().Execute()
    // ... verify
}
```

### Accounts Package

- [ ] **Step 1: Create testroot_test.go and accounts.go**

`accounts.go`:
```go
package accounts

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{Use: "accounts", Short: "Manage platform accounts"}
```

- [ ] **Step 2: Write failing test for accounts get**

`TestAccountsGet` — server returns `{"id": 1, "name": "Test Account"}`, args: `platform accounts get` (uses `--account-id` from context).

**IMPORTANT:** For `accounts get`, the account ID comes from `rctx.AccountID` (global flag / profile), NOT from an `--id` flag. Same for update and delete.

- [ ] **Step 3: Implement create.go, get.go, update.go, delete.go**

Flags:
- create: `--name` (required) — calls `client.CreateAccount(ctx, opts)` → `contract.Success(account)`. Note: create does NOT use `rctx.AccountID` — it creates a new account.
- get: no resource-specific flags — uses `rctx.AccountID` → `client.GetAccount(ctx, rctx.AccountID)` → `contract.Success(account)`
- update: `--name` (at least one Changed) — uses `rctx.AccountID` → `client.UpdateAccount(ctx, rctx.AccountID, opts)` → `contract.Success(account)`
- delete: no resource-specific flags — uses `rctx.AccountID` → `client.DeleteAccount(ctx, rctx.AccountID)` → `contract.Success(map[string]any{"deleted": true, "id": rctx.AccountID})`

- [ ] **Step 4: Add tests for create and delete**

- [ ] **Step 5: Run all accounts CLI tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/platform/accounts/ -v`
Expected: All PASS

### Account Users Package

- [ ] **Step 6: Create testroot_test.go and accountusers.go**

`accountusers.go` — package `accountusers`, `Use: "account-users"`, `Short: "Manage account users"`

- [ ] **Step 7: Write failing test for account-users list**

`TestAccountUsersList` — server returns `[{"account_id": 1, "user_id": 10, "role": "administrator"}]`, args: `platform account-users list`

- [ ] **Step 8: Implement list.go, create.go, delete.go**

All use `rctx.AccountID` from context.

Flags:
- list: no extra flags → `client.ListAccountUsers(ctx, rctx.AccountID)` → `contract.SuccessList`
- create: `--user-id` (required int), `--role` (optional string) → `client.CreateAccountUser(ctx, rctx.AccountID, opts)` → `contract.Success(map[string]any{"created": true})`
- delete: `--user-id` (required int) → `client.DeleteAccountUser(ctx, rctx.AccountID, userID)` → `contract.Success(map[string]any{"deleted": true})`

- [ ] **Step 9: Add tests for create and delete**

- [ ] **Step 10: Run all account-users CLI tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/platform/accountusers/ -v`
Expected: All PASS

### Agent Bots Package

- [ ] **Step 11: Create testroot_test.go and agentbots.go**

`agentbots.go` — package `agentbots`, `Use: "agent-bots"`, `Short: "Manage platform agent bots"`

- [ ] **Step 12: Write failing test for agent-bots list**

`TestAgentBotsList` — server returns `[{"id": 1, "name": "Helper Bot"}]`, args: `platform agent-bots list`

- [ ] **Step 13: Implement list.go, get.go, create.go, update.go, delete.go**

No `--account-id` needed — these are global resources.

Flags:
- list: no flags → `client.ListAgentBots(ctx)` → `contract.SuccessList`
- get: `--id` (required) → `client.GetAgentBot(ctx, id)` → `contract.Success`
- create: `--name` (required), `--description` (optional), `--outgoing-url` (optional), `--bot-type` (optional), `--bot-config` (optional JSON) → `client.CreateAgentBot(ctx, opts)` → `contract.Success`
- update: `--id` (required), `--name`, `--description`, `--outgoing-url`, `--bot-type`, `--bot-config` (at least one Changed) → `client.UpdateAgentBot(ctx, id, opts)` → `contract.Success`
- delete: `--id` (required) → `client.DeleteAgentBot(ctx, id)` → `contract.Success(map[string]any{"deleted": true, "id": id})`

For `--bot-config`, parse with `json.Unmarshal` into `any` when flag is Changed.

- [ ] **Step 14: Add tests for create and delete**

- [ ] **Step 15: Run all agent-bots CLI tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/platform/agentbots/ -v`
Expected: All PASS

### Users Package

- [ ] **Step 16: Create testroot_test.go and users.go**

`users.go` — package `users`, `Use: "users"`, `Short: "Manage platform users"`

- [ ] **Step 17: Write failing test for users create**

`TestUsersCreate` — server verifies POST, body has name/email/password, returns `{"id": 1, "name": "Alice", "email": "alice@test.com"}`, args: `platform users create --name Alice --email alice@test.com --password secret123`

- [ ] **Step 18: Implement create.go, get.go, update.go, delete.go, login.go**

Flags:
- create: `--name` (required), `--email` (required), `--password` (required), `--custom-attributes` (optional JSON) → `client.CreateUser(ctx, opts)` → `contract.Success`
- get: `--id` (required) → `client.GetUser(ctx, id)` → `contract.Success`
- update: `--id` (required), `--name`, `--email`, `--password`, `--custom-attributes` (at least one Changed) → `client.UpdateUser(ctx, id, opts)` → `contract.Success`
- delete: `--id` (required) → `client.DeleteUser(ctx, id)` → `contract.Success(map[string]any{"deleted": true, "id": id})`
- login: `--id` (required) → `client.GetUserSSOLink(ctx, id)` → `contract.Success(link)` — returns `{"ok": true, "data": {"url": "https://..."}}`

For `--custom-attributes`, parse with `json.Unmarshal` into `any` when flag is Changed.

- [ ] **Step 19: Add tests for get, delete, and login**

- `TestUsersLogin` — verify GET to `/platform/api/v1/users/1/login`, response contains url field

- [ ] **Step 20: Run all users CLI tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/platform/users/ -v`
Expected: All PASS

### Commit All Platform CLI

- [ ] **Step 21: Commit all 4 packages**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/cli/platform/
git commit -m "feat(cli): add platform command groups — accounts, account-users, agent-bots, users"
```

---

## Task 7: Register Platform Group

**Files:**
- Create: `internal/cli/platform/platform.go`
- Modify: `internal/cli/root.go`

- [ ] **Step 1: Create platform.go group command**

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

- [ ] **Step 2: Register in root.go**

Add import `cliplatform "github.com/chatwoot/chatwoot-cli/internal/cli/platform"` and `rootCmd.AddCommand(cliplatform.Cmd)` in the `init()` function.

- [ ] **Step 3: Build and verify**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./cmd/chatwoot/ && ./chatwoot platform --help`
Expected: Shows accounts, account-users, agent-bots, users

- [ ] **Step 4: Run full test suite**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./... 2>&1 | grep -E "^(ok|FAIL)"`
Expected: All pass

- [ ] **Step 5: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/cli/platform/platform.go internal/cli/root.go
git commit -m "feat(cli): register platform command group with 4 resource families"
```

---

## Task 8: Platform Verification

**Files:**
- No new files

- [ ] **Step 1: Verify platform command tree**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot platform --help`
Expected: Shows accounts, account-users, agent-bots, users

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot platform accounts --help`
Expected: Shows create, get, update, delete

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot platform account-users --help`
Expected: Shows list, create, delete

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot platform agent-bots --help`
Expected: Shows list, get, create, update, delete

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot platform users --help`
Expected: Shows create, get, update, delete, login

- [ ] **Step 2: Verify existing commands still work**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application --help`
Expected: Shows all Sprint C+D commands

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot version`
Expected: JSON envelope with version info

---

## Task 9: Add Client Model Types

**Files:**
- Modify: `internal/chatwoot/clientapi/models.go`

- [ ] **Step 1: Add all new types to models.go**

Add after existing `CreateContactOpts`:

```go
type UpdateContactOpts struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Phone *string `json:"phone_number,omitempty"`
}

// Conversation represents a client API conversation.
type Conversation struct {
	ID                int    `json:"id"`
	InboxID           int    `json:"inbox_id,omitempty"`
	Status            string `json:"status,omitempty"`
	AgentID           int    `json:"agent_id,omitempty"`
	ContactLastSeenAt string `json:"contact_last_seen_at,omitempty"`
}

// Message represents a client API message.
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

- [ ] **Step 2: Verify build**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./cmd/chatwoot/`
Expected: Success

- [ ] **Step 3: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/clientapi/models.go
git commit -m "feat(clientapi): add Sprint E model types for conversations and messages"
```

---

## Task 10: Add ResolveClientAuth Helper

**Files:**
- Modify: `internal/cli/cmdutil/context.go`

- [ ] **Step 1: Add ClientAuthContext and ResolveClientAuth**

Add to `internal/cli/cmdutil/context.go`:

```go
// ClientAuthContext holds identifier-based auth for client API commands.
type ClientAuthContext struct {
	InboxIdentifier   string
	ContactIdentifier string
}

// ResolveClientAuth resolves client API identifiers from flags and env vars.
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

- [ ] **Step 2: Verify build**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./cmd/chatwoot/`
Expected: Success

- [ ] **Step 3: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/cli/cmdutil/context.go
git commit -m "feat(cmdutil): add ResolveClientAuth for client API identifier-based auth"
```

---

## Task 11: Complete Client Contacts API + Conversations API + Messages API

**Files:**
- Modify: `internal/chatwoot/clientapi/contacts.go`
- Modify: `internal/chatwoot/clientapi/clientapi_test.go`
- Create: `internal/chatwoot/clientapi/conversations.go`
- Create: `internal/chatwoot/clientapi/conversations_test.go`
- Create: `internal/chatwoot/clientapi/messages.go`
- Create: `internal/chatwoot/clientapi/messages_test.go`

All client API paths are anchored at `/public/api/v1/inboxes/{inbox_identifier}/...`.

### Contacts

- [ ] **Step 1: Add UpdateContact to contacts.go**

`UpdateContact(ctx, contactIdentifier string, UpdateContactOpts) (*Contact, error)` — PATCH `/public/api/v1/inboxes/{inbox}/contacts/{contact}` — marshal opts, decode into `Contact`

- [ ] **Step 2: Add TestUpdateContact to clientapi_test.go**

Verify PATCH method, path includes both inbox and contact identifiers, request body has name field.

### Conversations

- [ ] **Step 3: Create conversations.go with 6 methods**

**HTTP methods and paths (all under `.../contacts/{contact_identifier}/conversations`):**
- `ListConversations(ctx, contactIdentifier string) ([]Conversation, error)` — GET `.../conversations` — decode into `[]Conversation`
- `GetConversation(ctx, contactIdentifier string, conversationID int) (*Conversation, error)` — GET `.../conversations/{id}` — decode into `Conversation`
- `CreateConversation(ctx, contactIdentifier string) (*Conversation, error)` — POST `.../conversations` — no body needed, decode into `Conversation`
- `ToggleStatus(ctx, contactIdentifier string, conversationID int) (*Conversation, error)` — POST `.../conversations/{id}/toggle_status` — no body, decode into `Conversation`
- `ToggleTyping(ctx, contactIdentifier string, conversationID int, ToggleTypingOpts) error` — POST `.../conversations/{id}/toggle_typing` — marshal opts, decode with `nil`
- `UpdateLastSeen(ctx, contactIdentifier string, conversationID int) error` — POST `.../conversations/{id}/update_last_seen` — no body, decode with `nil`

Helper for building the base path:
```go
func (c *Client) conversationBasePath(contactIdentifier string) string {
    return fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s/conversations", c.inboxIdentifier, contactIdentifier)
}
```

- [ ] **Step 4: Create conversations_test.go**

Tests:
- `TestListConversations` — verify GET, path `.../contacts/contact-xyz/conversations`
- `TestCreateConversation` — verify POST, same base path
- `TestToggleStatus` — verify POST, path `.../conversations/5/toggle_status`
- `TestToggleTyping` — verify POST, body has `typing_status` field

### Messages

- [ ] **Step 5: Create messages.go with 3 methods**

**HTTP methods and paths (all under `.../conversations/{conversation_id}/messages`):**
- `ListMessages(ctx, contactIdentifier string, conversationID int) ([]Message, error)` — GET — decode into `[]Message`
- `CreateMessage(ctx, contactIdentifier string, conversationID int, CreateMessageOpts) (*Message, error)` — POST — marshal opts, decode into `Message`
- `UpdateMessage(ctx, contactIdentifier string, conversationID int, messageID int, UpdateMessageOpts) (*Message, error)` — PATCH `.../messages/{mid}` — marshal opts, decode into `Message`

Helper:
```go
func (c *Client) messageBasePath(contactIdentifier string, conversationID int) string {
    return fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s/conversations/%d/messages", c.inboxIdentifier, contactIdentifier, conversationID)
}
```

- [ ] **Step 6: Create messages_test.go**

Tests:
- `TestListMessages` — verify GET, correct deeply nested path
- `TestCreateMessage` — verify POST, body has `content`

### Run and Commit

- [ ] **Step 7: Run all client API tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/clientapi/ -v`
Expected: All PASS

- [ ] **Step 8: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/clientapi/
git commit -m "feat(clientapi): add conversations and messages API clients, complete contacts"
```

---

## Task 12: Client CLI Commands

**Files:**
- Create: all files in `internal/cli/client/contacts/`
- Create: all files in `internal/cli/client/conversations/`
- Create: all files in `internal/cli/client/messages/`

### Client CLI Handler Pattern

Client CLI handlers use a different pipeline — no token auth, uses identifiers:

```go
func runSomething(cmd *cobra.Command, args []string) error {
    rctx, err := cmdutil.ResolveContext(cmd)
    if err != nil {
        return cmdutil.WriteError(cmd, contract.ErrCodeConfig, err.Error())
    }
    clientAuth, err := cmdutil.ResolveClientAuth(cmd)
    if err != nil {
        return cmdutil.WriteError(cmd, contract.ErrCodeAuth, err.Error())
    }
    transport := chatwoot.NewClient(rctx.BaseURL, "", "")
    client := clientapi.NewClient(transport, clientAuth.InboxIdentifier)
    // ... call client methods with clientAuth.ContactIdentifier
}
```

**Standard imports for client CLI handlers:**
```go
import (
    "context"
    clientapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/clientapi"
    chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
    "github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
    "github.com/chatwoot/chatwoot-cli/internal/contract"
    "github.com/spf13/cobra"
)
```

Note: No `credentials` import needed — client API doesn't use the credential system.

**Client testroot_test.go pattern** (includes persistent `--inbox-id` and `--contact-id` flags on the client group):
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

**Client test pattern** (uses `CHATWOOT_INBOX_IDENTIFIER` and `CHATWOOT_CONTACT_IDENTIFIER` env vars):
```go
func TestSomething(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(/* response */)
    }))
    defer srv.Close()
    t.Setenv("CHATWOOT_BASE_URL", srv.URL)
    t.Setenv("CHATWOOT_INBOX_IDENTIFIER", "inbox-abc")
    t.Setenv("CHATWOOT_CONTACT_IDENTIFIER", "contact-xyz")
    var stdout bytes.Buffer
    Cmd.SetOut(&stdout)
    Cmd.SetErr(&bytes.Buffer{})
    Cmd.Root().SetArgs([]string{"client", "<resource>", "<action>", ...flags})
    err := Cmd.Root().Execute()
    // ... verify
}
```

### Contacts Package

- [ ] **Step 1: Create testroot_test.go and contacts.go**

`contacts.go` — package `contacts`, `Use: "contacts"`, `Short: "Manage client contacts"`

- [ ] **Step 2: Implement create.go, get.go, update.go**

Flags:
- create: `--name` (optional), `--email` (optional), `--phone` (optional) — does NOT require `--contact-id` (creates a new contact). Uses `clientAuth.InboxIdentifier` only.
- get: uses `clientAuth.ContactIdentifier` (from persistent `--contact-id` flag) → `client.GetContact(ctx, clientAuth.ContactIdentifier)` → `contract.Success`
- update: `--name`, `--email`, `--phone` (at least one Changed) → `client.UpdateContact(ctx, clientAuth.ContactIdentifier, opts)` → `contract.Success`

For get and update, if `clientAuth.ContactIdentifier` is empty, return error: `"--contact-id flag or CHATWOOT_CONTACT_IDENTIFIER env var required for this command"`.

- [ ] **Step 3: Write contacts_test.go**

Tests:
- `TestContactsCreate` — server returns contact with source_id, args: `client contacts create --name "Test User"`. Env: set `CHATWOOT_INBOX_IDENTIFIER` but NOT `CHATWOOT_CONTACT_IDENTIFIER`.
- `TestContactsGet` — args: `client contacts get`. Env: set both inbox and contact identifiers.

- [ ] **Step 4: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/client/contacts/ -v`
Expected: All PASS

### Conversations Package

- [ ] **Step 5: Create testroot_test.go and conversations.go**

`conversations.go` — package `conversations`, `Use: "conversations"`, `Short: "Manage client conversations"`

- [ ] **Step 6: Implement list.go, get.go, create.go, toggle_status.go, toggle_typing.go, update_last_seen.go**

All require `clientAuth.ContactIdentifier` (validate non-empty at start of handler).

Flags:
- list: no extra flags → `client.ListConversations(ctx, contactID)` → `contract.SuccessList`
- get: `--conversation-id` (required int) → `client.GetConversation(ctx, contactID, convID)` → `contract.Success`
- create: no extra flags → `client.CreateConversation(ctx, contactID)` → `contract.Success`
- toggle-status: `--conversation-id` (required) → `client.ToggleStatus(ctx, contactID, convID)` → `contract.Success`
- toggle-typing: `--conversation-id` (required), `--status` (required: "on" or "off") → build `ToggleTypingOpts{TypingStatus: status}` → `client.ToggleTyping(ctx, contactID, convID, opts)` → `contract.Success(map[string]any{"typing_status": status})`
- update-last-seen: `--conversation-id` (required) → `client.UpdateLastSeen(ctx, contactID, convID)` → `contract.Success(map[string]any{"updated": true})`

- [ ] **Step 7: Write conversations_test.go**

Tests:
- `TestConversationsList` — server returns `[{"id": 1, "status": "open"}]`
- `TestConversationsCreate` — verify POST
- `TestConversationsToggleTyping` — args include `--conversation-id 1 --status on`, verify body has `typing_status`

- [ ] **Step 8: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/client/conversations/ -v`
Expected: All PASS

### Messages Package

- [ ] **Step 9: Create testroot_test.go and messages.go**

`messages.go` — package `messages`, `Use: "messages"`, `Short: "Manage client messages"`

- [ ] **Step 10: Implement list.go, create.go, update.go**

All require `clientAuth.ContactIdentifier` (validate non-empty).

Flags:
- list: `--conversation-id` (required int) → `client.ListMessages(ctx, contactID, convID)` → `contract.SuccessList`
- create: `--conversation-id` (required), `--content` (required), `--type` (optional, maps to MessageType) → build `CreateMessageOpts` → `client.CreateMessage(ctx, contactID, convID, opts)` → `contract.Success`
- update: `--conversation-id` (required), `--message-id` (required int), `--content` (required) → build `UpdateMessageOpts{Content: content}` → `client.UpdateMessage(ctx, contactID, convID, msgID, opts)` → `contract.Success`

- [ ] **Step 11: Write messages_test.go**

Tests:
- `TestMessagesList` — args include `--conversation-id 1`, server returns `[{"id": 1, "content": "Hello"}]`
- `TestMessagesCreate` — args include `--conversation-id 1 --content "Hello"`, verify POST body has `content`

- [ ] **Step 12: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/client/messages/ -v`
Expected: All PASS

### Commit All Client CLI

- [ ] **Step 13: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/cli/client/
git commit -m "feat(cli): add client command groups — contacts, conversations, messages"
```

---

## Task 13: Register Client Group

**Files:**
- Create: `internal/cli/client/client.go`
- Modify: `internal/cli/root.go`

- [ ] **Step 1: Create client.go group command**

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

- [ ] **Step 2: Register in root.go**

Add import `cliclient "github.com/chatwoot/chatwoot-cli/internal/cli/client"` and `rootCmd.AddCommand(cliclient.Cmd)`.

- [ ] **Step 3: Build and verify**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./cmd/chatwoot/ && ./chatwoot client --help`
Expected: Shows contacts, conversations, messages plus persistent flags --inbox-id and --contact-id

- [ ] **Step 4: Run full test suite**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./... 2>&1 | grep -E "^(ok|FAIL)"`
Expected: All pass

- [ ] **Step 5: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/cli/client/client.go internal/cli/root.go
git commit -m "feat(cli): register client command group with 3 resource families"
```

---

## Task 14: Full Test Suite and Exit Criteria Verification

**Files:**
- No new files

- [ ] **Step 1: Run full test suite**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./... -v 2>&1 | tail -80`
Expected: All tests PASS

- [ ] **Step 2: Run vet and build**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go vet ./... && go build ./cmd/chatwoot/`
Expected: No errors

- [ ] **Step 3: Verify platform command tree**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot platform --help`
Expected: Shows accounts, account-users, agent-bots, users

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot platform users --help`
Expected: Shows create, get, update, delete, login

- [ ] **Step 4: Verify client command tree**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot client --help`
Expected: Shows contacts, conversations, messages plus persistent --inbox-id, --contact-id flags

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot client conversations --help`
Expected: Shows list, get, create, toggle-status, toggle-typing, update-last-seen

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot client messages --help`
Expected: Shows list, create, update

- [ ] **Step 5: Verify existing commands still work**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application --help`
Expected: Shows all Sprint C+D commands

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot version`
Expected: JSON envelope with version info

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot auth --help`
Expected: Shows set, status, clear

---

## Sprint E Exit Criteria Checklist

- [ ] `go test ./...` passes with all tests green
- [ ] `go vet ./...` clean
- [ ] `go build ./cmd/chatwoot/` succeeds
- [ ] `chatwoot platform --help` shows accounts, account-users, agent-bots, users
- [ ] `chatwoot platform accounts create --name X` produces JSON envelope
- [ ] `chatwoot platform account-users list` uses `--account-id` from profile
- [ ] `chatwoot platform agent-bots list` works without `--account-id`
- [ ] `chatwoot platform users login --id 1` returns SSO URL in JSON envelope
- [ ] `chatwoot client --help` shows contacts, conversations, messages
- [ ] `chatwoot client contacts create --inbox-id X` works without `--contact-id`
- [ ] `chatwoot client messages create --inbox-id X --contact-id Y --conversation-id 1 --content "hello"` produces JSON envelope
- [ ] Client commands fail with clear error when `--inbox-id` is missing
- [ ] Platform commands use `CHATWOOT_PLATFORM_TOKEN` auth
- [ ] Existing application and auth commands still work
- [ ] All commands produce valid JSON envelopes on stdout
- [ ] All commands use the cmdutil pipeline (no business logic in handlers)

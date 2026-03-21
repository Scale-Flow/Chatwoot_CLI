# Sprint C: Core Application Commands (Layer 7) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the complete stable CLI command surface for contacts, conversations, messages, and inboxes — enabling core agent workflows (triage, summarize, reply, resolve).

**Architecture:** Extract shared CLI helpers into `internal/cli/cmdutil/` to eliminate import-cycle duplication. Add ~30 API client methods across 4 resource files. Wire ~35 CLI commands across 4 new command packages, each following the thin handler pattern: ResolveContext → ResolveAuth → transport → API client → contract envelope.

**Tech Stack:** Go 1.26.1, `spf13/cobra` (CLI), `spf13/viper` (config), `zalando/go-keyring` (keychain), `log/slog` (diagnostics), `encoding/json` (serialization), `net/http/httptest` (testing)

**Spec:** `docs/superpowers/specs/2026-03-20-sprint-c-core-application-commands.md`

---

## File Structure

### Shared Helpers

| File | Responsibility |
|------|---------------|
| `internal/cli/cmdutil/context.go` | NEW: RuntimeContext, ResolveContext, resolveContextFromPath, ResolveAuth, WriteError |
| `internal/cli/cmdutil/context_test.go` | NEW: migrated context resolution tests |
| `internal/cli/cmdutil/flags.go` | NEW: PaginationFlags, AddPaginationFlags, GetPaginationFlags, Pretty |
| `internal/cli/cmdutil/flags_test.go` | NEW: pagination and pretty flag tests |

### Migration (existing files)

| File | Responsibility |
|------|---------------|
| `internal/cli/context.go` | MODIFY: remove functions that moved to cmdutil |
| `internal/cli/context_test.go` | MODIFY: remove tests that moved to cmdutil |
| `internal/cli/application/profile.go` | MODIFY: replace inlined helpers with cmdutil imports |
| `internal/cli/application/profile_test.go` | MODIFY: update if needed for cmdutil changes |

### Model Types

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/application/models.go` | MODIFY: add Contact, Message, Inbox, Agent, AgentBot, ConversationMeta, and all Opts types |

### API Client — Contacts

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/application/contacts.go` | NEW: 11 contact API methods |
| `internal/chatwoot/application/contacts_test.go` | NEW: contact API tests |

### API Client — Conversations

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/application/conversations.go` | MODIFY: update ListConversations signature, add 9 new methods |
| `internal/chatwoot/application/application_test.go` | MODIFY: update existing tests for new ListConversations signature, add new conversation tests |

### API Client — Messages

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/application/messages.go` | NEW: 3 message API methods |
| `internal/chatwoot/application/messages_test.go` | NEW: message API tests |

### API Client — Inboxes

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/application/inboxes.go` | NEW: 10 inbox API methods |
| `internal/chatwoot/application/inboxes_test.go` | NEW: inbox API tests |

### CLI — Contacts Commands

| File | Responsibility |
|------|---------------|
| `internal/cli/application/contacts/contacts.go` | NEW: Cmd group, register subcommands |
| `internal/cli/application/contacts/list.go` | NEW: contacts list |
| `internal/cli/application/contacts/get.go` | NEW: contacts get |
| `internal/cli/application/contacts/create.go` | NEW: contacts create |
| `internal/cli/application/contacts/update.go` | NEW: contacts update |
| `internal/cli/application/contacts/delete.go` | NEW: contacts delete |
| `internal/cli/application/contacts/search.go` | NEW: contacts search |
| `internal/cli/application/contacts/filter.go` | NEW: contacts filter |
| `internal/cli/application/contacts/merge.go` | NEW: contacts merge |
| `internal/cli/application/contacts/labels.go` | NEW: contacts labels list/set |
| `internal/cli/application/contacts/conversations.go` | NEW: contacts conversations list |
| `internal/cli/application/contacts/contacts_test.go` | NEW: tests |
| `internal/cli/application/contacts/testroot_test.go` | NEW: test helper |

### CLI — Conversations Commands

| File | Responsibility |
|------|---------------|
| `internal/cli/application/conversations/conversations.go` | NEW: Cmd group |
| `internal/cli/application/conversations/list.go` | NEW: conversations list |
| `internal/cli/application/conversations/get.go` | NEW: conversations get |
| `internal/cli/application/conversations/create.go` | NEW: conversations create |
| `internal/cli/application/conversations/update.go` | NEW: conversations update |
| `internal/cli/application/conversations/filter.go` | NEW: conversations filter |
| `internal/cli/application/conversations/meta.go` | NEW: conversations meta |
| `internal/cli/application/conversations/toggle_status.go` | NEW: conversations toggle-status |
| `internal/cli/application/conversations/toggle_priority.go` | NEW: conversations toggle-priority |
| `internal/cli/application/conversations/assignments.go` | NEW: conversations assignments create |
| `internal/cli/application/conversations/labels.go` | NEW: conversations labels list/set |
| `internal/cli/application/conversations/conversations_test.go` | NEW: tests |
| `internal/cli/application/conversations/testroot_test.go` | NEW: test helper |

### CLI — Messages Commands

| File | Responsibility |
|------|---------------|
| `internal/cli/application/messages/messages.go` | NEW: Cmd group |
| `internal/cli/application/messages/list.go` | NEW: messages list |
| `internal/cli/application/messages/create.go` | NEW: messages create |
| `internal/cli/application/messages/delete.go` | NEW: messages delete |
| `internal/cli/application/messages/messages_test.go` | NEW: tests |
| `internal/cli/application/messages/testroot_test.go` | NEW: test helper |

### CLI — Inboxes Commands

| File | Responsibility |
|------|---------------|
| `internal/cli/application/inboxes/inboxes.go` | NEW: Cmd group |
| `internal/cli/application/inboxes/list.go` | NEW: inboxes list |
| `internal/cli/application/inboxes/get.go` | NEW: inboxes get |
| `internal/cli/application/inboxes/create.go` | NEW: inboxes create |
| `internal/cli/application/inboxes/update.go` | NEW: inboxes update |
| `internal/cli/application/inboxes/members.go` | NEW: inboxes members list/add/update/delete |
| `internal/cli/application/inboxes/agent_bot.go` | NEW: inboxes agent-bot get/set |
| `internal/cli/application/inboxes/inboxes_test.go` | NEW: tests |
| `internal/cli/application/inboxes/testroot_test.go` | NEW: test helper |

### Registration

| File | Responsibility |
|------|---------------|
| `internal/cli/application/application.go` | MODIFY: register contacts, conversations, messages, inboxes subgroups |

---

## Task 1: Extract cmdutil Package

**Files:**
- Create: `internal/cli/cmdutil/context.go`
- Create: `internal/cli/cmdutil/context_test.go`
- Create: `internal/cli/cmdutil/flags.go`
- Create: `internal/cli/cmdutil/flags_test.go`
- Modify: `internal/cli/context.go`
- Modify: `internal/cli/context_test.go`

- [ ] **Step 1: Create `cmdutil/context.go` by copying from `cli/context.go`**

Copy the following from `internal/cli/context.go` into `internal/cli/cmdutil/context.go`, changing package to `cmdutil`:

- `RuntimeContext` struct
- `ResolveContext` function
- `resolveContextFromPath` function
- `ResolveAuth` function
- `WriteError` function

All imports stay the same except the package declaration changes to `package cmdutil`.

Note: `WriteError` currently references `prettyFlag` (a package-level var in `cli`). Change it to read the flag from the command root:

```go
func WriteError(cmd *cobra.Command, code, message string) error {
	pretty, _ := cmd.Root().PersistentFlags().GetBool("pretty")
	resp := contract.Err(code, message)
	_ = contract.Write(cmd.OutOrStdout(), resp, pretty)
	return errors.New(message)
}
```

- [ ] **Step 2: Create `cmdutil/flags.go`**

```go
package cmdutil

import "github.com/spf13/cobra"

// PaginationFlags holds pagination flag values.
type PaginationFlags struct {
	Page    int
	PerPage int
	All     bool
}

// AddPaginationFlags registers --page, --per-page, and --all flags on cmd.
func AddPaginationFlags(cmd *cobra.Command) {
	cmd.Flags().Int("page", 1, "Page number")
	cmd.Flags().Int("per-page", 25, "Items per page")
	cmd.Flags().Bool("all", false, "Fetch all pages")
}

// GetPaginationFlags reads pagination flag values from cmd.
func GetPaginationFlags(cmd *cobra.Command) PaginationFlags {
	page, _ := cmd.Flags().GetInt("page")
	perPage, _ := cmd.Flags().GetInt("per-page")
	all, _ := cmd.Flags().GetBool("all")
	return PaginationFlags{Page: page, PerPage: perPage, All: all}
}

// Pretty reads the --pretty flag from the root command.
func Pretty(cmd *cobra.Command) bool {
	pretty, _ := cmd.Root().PersistentFlags().GetBool("pretty")
	return pretty
}
```

- [ ] **Step 3: Create `cmdutil/context_test.go`**

Move the 5 `TestResolveContext*` tests from `internal/cli/context_test.go` into `internal/cli/cmdutil/context_test.go`, changing package to `cmdutil`. The tests use `resolveContextFromPath` which is unexported — since it's in the same package, this works.

- [ ] **Step 4: Create `cmdutil/flags_test.go`**

```go
package cmdutil

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestAddAndGetPaginationFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	AddPaginationFlags(cmd)

	// Defaults
	pf := GetPaginationFlags(cmd)
	if pf.Page != 1 {
		t.Errorf("Page = %d, want 1", pf.Page)
	}
	if pf.PerPage != 25 {
		t.Errorf("PerPage = %d, want 25", pf.PerPage)
	}
	if pf.All {
		t.Error("All = true, want false")
	}
}

func TestPrettyFromRoot(t *testing.T) {
	root := &cobra.Command{Use: "root"}
	root.PersistentFlags().Bool("pretty", false, "")
	child := &cobra.Command{Use: "child"}
	root.AddCommand(child)

	if Pretty(child) {
		t.Error("Pretty = true, want false")
	}
}
```

- [ ] **Step 5: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/cmdutil/ -v`
Expected: PASS

- [ ] **Step 6: Strip `cli/context.go` down to `ResolveProfileName` only**

Remove from `internal/cli/context.go`: `RuntimeContext`, `ResolveContext`, `resolveContextFromPath`, `ResolveAuth`, `WriteError`, and their associated imports. Keep only `ResolveProfileName` (still used by auth indirectly through its own local copy, but this is for reference/future use).

If nothing in `internal/cli/` references these anymore, the file can be reduced to just the `ResolveProfileName` function or deleted entirely. Check with `go build ./...`.

- [ ] **Step 7: Strip `cli/context_test.go`**

Remove the 5 `TestResolveContext*` tests that were migrated. If no tests remain, delete the file.

- [ ] **Step 8: Verify build and all tests pass**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./cmd/chatwoot/ && go test ./... 2>&1 | grep -E "^(ok|FAIL)"`
Expected: All pass

- [ ] **Step 9: Commit**

```bash
git add internal/cli/cmdutil/ internal/cli/context.go internal/cli/context_test.go
git commit -m "refactor(cli): extract cmdutil package for shared CLI helpers"
```

---

## Task 2: Migrate profile.go to cmdutil

**Files:**
- Modify: `internal/cli/application/profile.go`
- Modify: `internal/cli/application/profile_test.go` (if needed)

- [ ] **Step 1: Replace inlined helpers in profile.go with cmdutil imports**

In `internal/cli/application/profile.go`:

1. Replace import of `innerauth`, `config`, `credentials`, and the inlined helper functions with a single import of `"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"`.
2. Replace all calls to local `resolveContext(cmd)` with `cmdutil.ResolveContext(cmd)`.
3. Replace all calls to local `resolveAuth(...)` with `cmdutil.ResolveAuth(...)`.
4. Replace all calls to local `writeError(...)` with `cmdutil.WriteError(...)`.
5. Replace all calls to local `prettyFromRoot(cmd)` with `cmdutil.Pretty(cmd)`.
6. Delete the inlined helper functions (`resolveContext`, `resolveContextFromPath`, `resolveAuth`, `writeError`, `prettyFromRoot`) from `profile.go`.

- [ ] **Step 2: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/ -v`
Expected: PASS

- [ ] **Step 3: Run full test suite**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./cmd/chatwoot/ && go test ./... 2>&1 | grep -E "^(ok|FAIL)"`
Expected: All pass

- [ ] **Step 4: Commit**

```bash
git add internal/cli/application/profile.go
git commit -m "refactor(cli): migrate profile commands to cmdutil helpers"
```

---

## Task 3: Add Model Types

**Files:**
- Modify: `internal/chatwoot/application/models.go`

- [ ] **Step 1: Add all new types to models.go**

Add to `internal/chatwoot/application/models.go`:

```go
// Contact represents a Chatwoot contact.
type Contact struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Email     string           `json:"email,omitempty"`
	Phone     string           `json:"phone_number,omitempty"`
	AccountID int              `json:"account_id"`
	CreatedAt chatwoot.Timestamp `json:"created_at,omitempty"`
}

// Message represents a Chatwoot message.
type Message struct {
	ID             int              `json:"id"`
	Content        string           `json:"content,omitempty"`
	MessageType    int              `json:"message_type"`
	ContentType    string           `json:"content_type,omitempty"`
	Private        bool             `json:"private"`
	ConversationID int              `json:"conversation_id"`
	CreatedAt      chatwoot.Timestamp `json:"created_at,omitempty"`
}

// Inbox represents a Chatwoot inbox.
type Inbox struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	ChannelType          string `json:"channel_type,omitempty"`
	AvatarURL            string `json:"avatar_url,omitempty"`
	EnableAutoAssignment bool   `json:"enable_auto_assignment"`
}

// Agent represents a Chatwoot agent.
type Agent struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role,omitempty"`
}

// AgentBot represents a Chatwoot agent bot.
type AgentBot struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ConversationMeta holds conversation count metadata.
type ConversationMeta struct {
	AllCount      int `json:"all_count"`
	OpenCount     int `json:"open_count"`
	ResolvedCount int `json:"resolved_count"`
	PendingCount  int `json:"pending_count"`
	SnoozedCount  int `json:"snoozed_count"`
}

// --- Opts types for mutations ---

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

type ListConversationsOpts struct {
	Page    int
	PerPage int
	Status  string
	InboxID int
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

Note: The existing `ListConversationsOpts` in `conversations.go` should be removed from there since it now belongs in `models.go`. Check that the existing type matches the one above (it has `Page`, `Status`, `InboxID` — add `PerPage` if not already there).

- [ ] **Step 2: Verify build**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./cmd/chatwoot/`
Expected: Success

- [ ] **Step 3: Commit**

```bash
git add internal/chatwoot/application/models.go
git commit -m "feat(application): add model types for contacts, messages, inboxes, and opts"
```

---

## Task 4: Contacts API Client

**Files:**
- Create: `internal/chatwoot/application/contacts.go`
- Create: `internal/chatwoot/application/contacts_test.go`

- [ ] **Step 1: Write failing test for ListContacts**

Create `internal/chatwoot/application/contacts_test.go` with `TestListContacts` that:
- Sets up httptest server returning a JSON payload matching Chatwoot's contacts list format (`{"data": {"payload": [...]}, "meta": {...}}`)
- Creates transport + client
- Calls `ListContacts` with page 1
- Asserts correct URL path `/api/v1/accounts/1/contacts?page=1`
- Asserts returned contacts match mock data

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestListContacts`
Expected: FAIL — `ListContacts` not defined

- [ ] **Step 3: Implement all 11 contacts API methods**

Create `internal/chatwoot/application/contacts.go` with all methods from the spec:
- `ListContacts`, `GetContact`, `CreateContact`, `UpdateContact`, `DeleteContact`
- `SearchContacts`, `FilterContacts`, `MergeContacts`
- `ListContactLabels`, `SetContactLabels`, `ListContactConversations`

Each method follows the same pattern: build path with `fmt.Sprintf`, call `c.transport.DoWithRetry`, decode with `chatwoot.DecodeResponse`.

For list methods that return pagination, parse the Chatwoot response wrapper:
```go
var body struct {
	Data struct {
		Payload []Contact `json:"payload"`
		Meta    struct {
			// Chatwoot nests pagination fields directly
		} `json:"meta"`
	} `json:"data"`
	// Some endpoints put meta at top level
}
```

Note: Chatwoot API response format varies by endpoint. Contacts list wraps data in `{"data": {"payload": [...], "meta": {...}}}`. Check actual response format and decode accordingly.

- [ ] **Step 4: Add remaining tests**

Add tests for `GetContact`, `CreateContact`, `UpdateContact`, `DeleteContact`, `SearchContacts`, `FilterContacts`, `MergeContacts`, `ListContactLabels`, `SetContactLabels`, `ListContactConversations`. Each test verifies:
- Correct HTTP method
- Correct URL path (including account ID and resource ID)
- Request body (for POST/PUT methods)
- Response deserialization

- [ ] **Step 5: Run all contacts tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestContact`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add internal/chatwoot/application/contacts.go internal/chatwoot/application/contacts_test.go
git commit -m "feat(application): add contacts API client with 11 methods"
```

---

## Task 5: Conversations API Client Updates

**Files:**
- Modify: `internal/chatwoot/application/conversations.go`
- Modify: `internal/chatwoot/application/application_test.go`

- [ ] **Step 1: Update ListConversations to return pagination**

Change signature from `([]Conversation, error)` to `([]Conversation, *contract.Pagination, error)`. Parse the `meta` field from the Chatwoot response. Add the `contract` import.

- [ ] **Step 2: Update existing tests for new signature**

Update `TestListConversations` in `application_test.go` to handle the third return value.

- [ ] **Step 3: Write failing test for CreateConversation**

Add `TestCreateConversation` to `application_test.go`.

- [ ] **Step 4: Implement all 9 new conversation methods**

Add to `conversations.go`:
- `CreateConversation`, `UpdateConversation`, `FilterConversations`
- `GetConversationMeta`, `ToggleConversationStatus`, `ToggleConversationPriority`
- `AssignConversation`, `ListConversationLabels`, `SetConversationLabels`

- [ ] **Step 5: Add remaining tests**

Add tests for all 9 new methods.

- [ ] **Step 6: Run all conversation tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run "TestConversation|TestCreate|TestUpdate|TestFilter|TestToggle|TestAssign|TestLabel"`
Expected: All PASS

- [ ] **Step 7: Commit**

```bash
git add internal/chatwoot/application/conversations.go internal/chatwoot/application/application_test.go
git commit -m "feat(application): update ListConversations with pagination, add 9 conversation methods"
```

---

## Task 6: Messages API Client

**Files:**
- Create: `internal/chatwoot/application/messages.go`
- Create: `internal/chatwoot/application/messages_test.go`

- [ ] **Step 1: Write failing test for ListMessages**

Create `messages_test.go` with `TestListMessages`.

- [ ] **Step 2: Implement all 3 message methods**

Create `messages.go` with `ListMessages`, `CreateMessage`, `DeleteMessage`.

- [ ] **Step 3: Add remaining tests**

Add `TestCreateMessage`, `TestDeleteMessage`.

- [ ] **Step 4: Run all message tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestMessage`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add internal/chatwoot/application/messages.go internal/chatwoot/application/messages_test.go
git commit -m "feat(application): add messages API client with list, create, delete"
```

---

## Task 7: Inboxes API Client

**Files:**
- Create: `internal/chatwoot/application/inboxes.go`
- Create: `internal/chatwoot/application/inboxes_test.go`

- [ ] **Step 1: Write failing test for ListInboxes**

Create `inboxes_test.go` with `TestListInboxes`.

- [ ] **Step 2: Implement all 10 inbox methods**

Create `inboxes.go` with:
- `ListInboxes`, `GetInbox`, `CreateInbox`, `UpdateInbox`
- `ListInboxMembers`, `AddInboxMember`, `UpdateInboxMembers`, `RemoveInboxMember`
- `GetInboxAgentBot`, `SetInboxAgentBot`

Note: Member methods use the `/inbox_members` path (not nested under `/inboxes/{id}/members`). `AddInboxMember`, `UpdateInboxMembers`, `RemoveInboxMember` require `inbox_id` in the request body.

- [ ] **Step 3: Add remaining tests**

Add tests for all 10 methods.

- [ ] **Step 4: Run all inbox tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run "TestInbox|TestMember|TestAgentBot"`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add internal/chatwoot/application/inboxes.go internal/chatwoot/application/inboxes_test.go
git commit -m "feat(application): add inboxes API client with 10 methods"
```

---

## Task 8: Contacts CLI Commands

**Files:**
- Create: `internal/cli/application/contacts/contacts.go`
- Create: `internal/cli/application/contacts/list.go`
- Create: `internal/cli/application/contacts/get.go`
- Create: `internal/cli/application/contacts/create.go`
- Create: `internal/cli/application/contacts/update.go`
- Create: `internal/cli/application/contacts/delete.go`
- Create: `internal/cli/application/contacts/search.go`
- Create: `internal/cli/application/contacts/filter.go`
- Create: `internal/cli/application/contacts/merge.go`
- Create: `internal/cli/application/contacts/labels.go`
- Create: `internal/cli/application/contacts/conversations.go`
- Create: `internal/cli/application/contacts/contacts_test.go`
- Create: `internal/cli/application/contacts/testroot_test.go`

- [ ] **Step 1: Create testroot_test.go**

```go
package contacts

import "github.com/spf13/cobra"

func init() {
	root := &cobra.Command{
		Use:           "chatwoot",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().Bool("pretty", false, "Indent JSON output")
	root.PersistentFlags().String("profile", "", "Select named profile")
	root.PersistentFlags().String("base-url", "", "Override base URL")
	root.PersistentFlags().Int("account-id", 0, "Override account ID")
	appCmd := &cobra.Command{Use: "application"}
	appCmd.AddCommand(Cmd)
	root.AddCommand(appCmd)
}
```

- [ ] **Step 2: Create contacts.go group command**

```go
package contacts

import "github.com/spf13/cobra"

// Cmd is the contacts command group.
var Cmd = &cobra.Command{
	Use:   "contacts",
	Short: "Manage contacts",
}

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Manage contact labels",
}

var contactConversationsCmd = &cobra.Command{
	Use:   "conversations",
	Short: "List contact conversations",
}

func init() {
	Cmd.AddCommand(labelsCmd)
	Cmd.AddCommand(contactConversationsCmd)
}
```

- [ ] **Step 3: Write failing test for contacts list**

In `contacts_test.go`, write `TestContactsList` that:
- Sets up httptest server returning contacts JSON
- Sets `CHATWOOT_BASE_URL` and `CHATWOOT_ACCESS_TOKEN` env vars
- Executes `application contacts list`
- Verifies JSON envelope with `ok: true` and data array

- [ ] **Step 4: Implement list.go**

Handler follows the cmdutil pipeline pattern:
- `cmdutil.ResolveContext(cmd)` → `cmdutil.ResolveAuth(...)` → `chatwoot.NewClient(...)` → `appapi.NewClient(...)` → `client.ListContacts(...)` → `contract.SuccessList(...)` → `contract.Write(...)`
- Register pagination flags via `cmdutil.AddPaginationFlags(listCmd)` in `init()`
- Support `--all` flag using `chatwoot.ListAll`

- [ ] **Step 5: Implement get.go, create.go, update.go, delete.go**

Each follows the same thin handler pattern. Required flags:
- get: `--id` (required)
- create: `--name` (required), `--email`, `--phone`
- update: `--id` (required), `--name`, `--email`, `--phone` (at least one)
- delete: `--id` (required)

- [ ] **Step 6: Implement search.go, filter.go, merge.go**

- search: `--query` (required), `--page`, `--per-page`
- filter: `--payload` (required, JSON string), `--page`
- merge: `--base-id` (required), `--merge-id` (required)

For filter, parse `--payload` with `json.Unmarshal` into `[]any`.

- [ ] **Step 7: Implement labels.go and conversations.go**

- `labels list`: `--id` (required) → `client.ListContactLabels`
- `labels set`: `--id` (required), `--labels` (required, comma-separated) → split by comma, call `client.SetContactLabels`
- `conversations list`: `--id` (required) → `client.ListContactConversations`

Register subcommands on `labelsCmd` and `contactConversationsCmd`.

- [ ] **Step 8: Add tests for get, create, delete, search, merge, labels set**

Add tests to `contacts_test.go` for key commands. Each test:
- Sets up httptest server
- Uses env vars for auth
- Executes command via `Cmd.Root().Execute()`
- Verifies JSON output and HTTP request

- [ ] **Step 9: Run all contacts CLI tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/contacts/ -v`
Expected: All PASS

- [ ] **Step 10: Commit**

```bash
git add internal/cli/application/contacts/
git commit -m "feat(cli): add contacts command group with 11 commands"
```

---

## Task 9: Conversations CLI Commands

**Files:**
- Create: `internal/cli/application/conversations/conversations.go`
- Create: `internal/cli/application/conversations/list.go`
- Create: `internal/cli/application/conversations/get.go`
- Create: `internal/cli/application/conversations/create.go`
- Create: `internal/cli/application/conversations/update.go`
- Create: `internal/cli/application/conversations/filter.go`
- Create: `internal/cli/application/conversations/meta.go`
- Create: `internal/cli/application/conversations/toggle_status.go`
- Create: `internal/cli/application/conversations/toggle_priority.go`
- Create: `internal/cli/application/conversations/assignments.go`
- Create: `internal/cli/application/conversations/labels.go`
- Create: `internal/cli/application/conversations/conversations_test.go`
- Create: `internal/cli/application/conversations/testroot_test.go`

- [ ] **Step 1: Create testroot_test.go**

Same pattern as contacts testroot but with `Cmd` from this package.

- [ ] **Step 2: Create conversations.go group command**

Exports `Cmd`. Contains `labelsCmd` and `assignmentsCmd` subgroups. Registers all subcommands.

- [ ] **Step 3: Write failing test for conversations list**

Test verifies `--status open` flag flows through to the API request query parameter.

- [ ] **Step 4: Implement list.go**

List command with:
- `--status` (optional: open, resolved, pending, snoozed)
- `--inbox-id` (optional)
- Pagination flags via `cmdutil.AddPaginationFlags`
- `--all` support via `chatwoot.ListAll`

- [ ] **Step 5: Implement get.go, create.go, update.go**

- get: `--id` (required)
- create: `--contact-id` (required), `--inbox-id` (required), `--status` (optional), `--content` (optional, initial message)
- update: `--id` (required), `--status`, `--priority` (at least one)

- [ ] **Step 6: Implement filter.go, meta.go**

- filter: `--payload` (required JSON), `--page`
- meta: no flags, returns conversation counts

- [ ] **Step 7: Implement toggle_status.go, toggle_priority.go, assignments.go**

- toggle-status: `--id` (required), `--status` (required: open, resolved, pending, snoozed)
- toggle-priority: `--id` (required), `--priority` (required: urgent, high, medium, low, none)
- assignments create: `--id` (required), `--agent-id` and/or `--team-id`

- [ ] **Step 8: Implement labels.go**

- labels list: `--id` (required)
- labels set: `--id` (required), `--labels` (required, comma-separated)

- [ ] **Step 9: Add tests for key commands**

Test list (with --status), get, create, toggle-status, assignments create, labels set.

- [ ] **Step 10: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/conversations/ -v`
Expected: All PASS

- [ ] **Step 11: Commit**

```bash
git add internal/cli/application/conversations/
git commit -m "feat(cli): add conversations command group with 11 commands"
```

---

## Task 10: Messages CLI Commands

**Files:**
- Create: `internal/cli/application/messages/messages.go`
- Create: `internal/cli/application/messages/list.go`
- Create: `internal/cli/application/messages/create.go`
- Create: `internal/cli/application/messages/delete.go`
- Create: `internal/cli/application/messages/messages_test.go`
- Create: `internal/cli/application/messages/testroot_test.go`

- [ ] **Step 1: Create testroot_test.go and messages.go group**

- [ ] **Step 2: Write failing test for messages create**

Test verifies: `--conversation-id` maps to URL path, `--content` appears in request body, response is JSON envelope.

- [ ] **Step 3: Implement list.go, create.go, delete.go**

All commands require `--conversation-id` flag.
- list: just conversation-id
- create: `--conversation-id` (required), `--content` (required), `--message-type` (optional, default "outgoing"), `--private` (optional bool)
- delete: `--conversation-id` (required), `--id` (required)

- [ ] **Step 4: Add tests for list and delete**

- [ ] **Step 5: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/messages/ -v`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add internal/cli/application/messages/
git commit -m "feat(cli): add messages command group with list, create, delete"
```

---

## Task 11: Inboxes CLI Commands

**Files:**
- Create: `internal/cli/application/inboxes/inboxes.go`
- Create: `internal/cli/application/inboxes/list.go`
- Create: `internal/cli/application/inboxes/get.go`
- Create: `internal/cli/application/inboxes/create.go`
- Create: `internal/cli/application/inboxes/update.go`
- Create: `internal/cli/application/inboxes/members.go`
- Create: `internal/cli/application/inboxes/agent_bot.go`
- Create: `internal/cli/application/inboxes/inboxes_test.go`
- Create: `internal/cli/application/inboxes/testroot_test.go`

- [ ] **Step 1: Create testroot_test.go and inboxes.go group**

Group command with `membersCmd` and `agentBotCmd` subgroups.

- [ ] **Step 2: Write failing test for inboxes list**

- [ ] **Step 3: Implement list.go, get.go, create.go, update.go**

- list: no required flags (returns all inboxes for account)
- get: `--id` (required)
- create: `--name` (required), `--channel` (required, JSON string parsed into `any`)
- update: `--id` (required), `--name`, `--enable-auto-assignment` (at least one)

- [ ] **Step 4: Implement members.go**

Contains 4 subcommands registered on `membersCmd`:
- `members list`: `--inbox-id` (required)
- `members add`: `--inbox-id` (required), `--agent-ids` (required, comma-separated → parsed to `[]int`)
- `members update`: `--inbox-id` (required), `--agent-ids` (required)
- `members delete`: `--inbox-id` (required), `--agent-ids` (required)

- [ ] **Step 5: Implement agent_bot.go**

Contains 2 subcommands registered on `agentBotCmd`:
- `agent-bot get`: `--inbox-id` (required)
- `agent-bot set`: `--inbox-id` (required), `--agent-bot-id` (required)

- [ ] **Step 6: Add tests for list, get, members add**

- [ ] **Step 7: Run tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/inboxes/ -v`
Expected: All PASS

- [ ] **Step 8: Commit**

```bash
git add internal/cli/application/inboxes/
git commit -m "feat(cli): add inboxes command group with 10 commands"
```

---

## Task 12: Register Subgroups and Wire Everything

**Files:**
- Modify: `internal/cli/application/application.go`

- [ ] **Step 1: Register all 4 new subgroups**

Add to `internal/cli/application/application.go`:

```go
import (
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/contacts"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/conversations"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/inboxes"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/messages"
)
```

In `init()`:
```go
Cmd.AddCommand(contacts.Cmd)
Cmd.AddCommand(conversations.Cmd)
Cmd.AddCommand(messages.Cmd)
Cmd.AddCommand(inboxes.Cmd)
```

- [ ] **Step 2: Build and verify**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./cmd/chatwoot/`
Expected: Success

- [ ] **Step 3: Run full test suite**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./... 2>&1 | grep -E "^(ok|FAIL)"`
Expected: All pass

- [ ] **Step 4: Commit**

```bash
git add internal/cli/application/application.go
git commit -m "feat(cli): register contacts, conversations, messages, inboxes command groups"
```

---

## Task 13: Full Test Suite and Exit Criteria Verification

**Files:**
- No new files

- [ ] **Step 1: Run full test suite**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./... -v 2>&1 | tail -60`
Expected: All tests PASS

- [ ] **Step 2: Run vet and build**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go vet ./... && go build ./cmd/chatwoot/`
Expected: No errors

- [ ] **Step 3: Verify command tree**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application --help`
Expected: Shows contacts, conversations, inboxes, messages, profile subcommands

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application contacts --help`
Expected: Shows list, get, create, update, delete, search, filter, merge, labels, conversations

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application conversations --help`
Expected: Shows list, get, create, update, filter, meta, toggle-status, toggle-priority, assignments, labels

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application messages --help`
Expected: Shows list, create, delete

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application inboxes --help`
Expected: Shows list, get, create, update, members, agent-bot

- [ ] **Step 4: Verify global flags propagate**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application contacts list --help | grep -E "page|per-page|all"`
Expected: All three pagination flags listed

- [ ] **Step 5: Verify existing commands still work**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot version`
Expected: JSON envelope with version info

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot auth --help`
Expected: Shows set, status, clear

---

## Sprint C Exit Criteria Checklist

- [ ] `go test ./...` passes with all tests green
- [ ] `go vet ./...` clean
- [ ] `go build ./cmd/chatwoot/` succeeds
- [ ] `chatwoot application contacts list` returns JSON envelope with pagination
- [ ] `chatwoot application conversations list --status open` filters correctly
- [ ] `chatwoot application messages create --conversation-id N --content "..."` sends POST
- [ ] `chatwoot application inboxes list` returns inbox data
- [ ] `--all` flag auto-paginates on contacts and/or conversations list
- [ ] `--labels` comma-separated values work on label set commands
- [ ] `--payload` JSON filter works on filter commands
- [ ] All commands produce valid JSON envelopes on stdout
- [ ] All commands use the cmdutil pipeline (no inlined helpers)
- [ ] Error envelopes are consistent across all commands
- [ ] No business logic in command handler files (thin handlers only)
- [ ] Existing auth and profile commands still work

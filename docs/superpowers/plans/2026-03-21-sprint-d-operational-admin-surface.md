# Sprint D: Operational Admin Surface (Layer 8) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the complete Layer 8 operational admin surface — 14 resource families with ~70 API methods and ~69 CLI commands enabling team management, reporting, automation, and knowledge base operations.

**Architecture:** One file per resource in `internal/chatwoot/application/` for API clients, one package per resource in `internal/cli/application/` for CLI commands. All follow the Sprint C thin handler pattern: ResolveContext → ResolveAuth → transport → API client → contract envelope.

**Tech Stack:** Go 1.26.1, `spf13/cobra` (CLI), `spf13/viper` (config), `zalando/go-keyring` (keychain), `log/slog` (diagnostics), `encoding/json` (serialization), `net/http/httptest` (testing)

**Spec:** `docs/superpowers/specs/2026-03-21-sprint-d-operational-admin-surface.md`

**Reference patterns for subagents:**
- API client: `internal/chatwoot/application/contacts.go` (list with pagination, CRUD, labels, merge)
- CLI handler: `internal/cli/application/contacts/list.go` (cmdutil pipeline, pagination flags)
- CLI test: `internal/cli/application/contacts/contacts_test.go` (httptest, env vars, JSON verification)
- Members pattern: `internal/cli/application/inboxes/members.go` (nested member CRUD with comma-separated IDs)
- testroot_test.go: `internal/cli/application/contacts/testroot_test.go` (same boilerplate per package)

---

## File Structure

### Model Types

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/application/models.go` | MODIFY: add ~30 new core types and opts types for Sprint D resources |

### API Clients — Use-Case-Critical

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/application/teams.go` | NEW: 9 team + team member API methods |
| `internal/chatwoot/application/teams_test.go` | NEW: team API tests |
| `internal/chatwoot/application/agents.go` | NEW: 4 agent API methods |
| `internal/chatwoot/application/agents_test.go` | NEW: agent API tests |
| `internal/chatwoot/application/canned_responses.go` | NEW: 4 canned response API methods |
| `internal/chatwoot/application/canned_responses_test.go` | NEW: canned response API tests |
| `internal/chatwoot/application/reports.go` | NEW: 12 report API methods |
| `internal/chatwoot/application/reports_test.go` | NEW: report API tests |
| `internal/chatwoot/application/webhooks.go` | NEW: 4 webhook API methods |
| `internal/chatwoot/application/webhooks_test.go` | NEW: webhook API tests |
| `internal/chatwoot/application/automation_rules.go` | NEW: 5 automation rule API methods |
| `internal/chatwoot/application/automation_rules_test.go` | NEW: automation rule API tests |

### API Clients — Simple CRUD Batch

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/application/labels.go` | NEW: 5 label API methods |
| `internal/chatwoot/application/labels_test.go` | NEW: label API tests |
| `internal/chatwoot/application/custom_attributes.go` | NEW: 5 custom attribute API methods |
| `internal/chatwoot/application/custom_attributes_test.go` | NEW: custom attribute API tests |
| `internal/chatwoot/application/custom_filters.go` | NEW: 5 custom filter API methods |
| `internal/chatwoot/application/custom_filters_test.go` | NEW: custom filter API tests |
| `internal/chatwoot/application/account.go` | NEW: 2 account API methods |
| `internal/chatwoot/application/account_test.go` | NEW: account API tests |
| `internal/chatwoot/application/agent_bots.go` | NEW: 5 account agent bot API methods |
| `internal/chatwoot/application/agent_bots_test.go` | NEW: agent bot API tests |
| `internal/chatwoot/application/audit_logs.go` | NEW: 1 audit log API method |
| `internal/chatwoot/application/audit_logs_test.go` | NEW: audit log API tests |

### API Clients — Integrations + Help Center

| File | Responsibility |
|------|---------------|
| `internal/chatwoot/application/integrations.go` | NEW: 4 integration API methods |
| `internal/chatwoot/application/integrations_test.go` | NEW: integration API tests |
| `internal/chatwoot/application/help_center.go` | NEW: 5 help center API methods |
| `internal/chatwoot/application/help_center_test.go` | NEW: help center API tests |

### CLI — Teams Commands

| File | Responsibility |
|------|---------------|
| `internal/cli/application/teams/teams.go` | NEW: Cmd group, membersCmd subgroup |
| `internal/cli/application/teams/list.go` | NEW: teams list |
| `internal/cli/application/teams/get.go` | NEW: teams get |
| `internal/cli/application/teams/create.go` | NEW: teams create |
| `internal/cli/application/teams/update.go` | NEW: teams update |
| `internal/cli/application/teams/delete.go` | NEW: teams delete |
| `internal/cli/application/teams/members.go` | NEW: teams members list/add/update/delete |
| `internal/cli/application/teams/teams_test.go` | NEW: tests |
| `internal/cli/application/teams/testroot_test.go` | NEW: test helper |

### CLI — Agents Commands

| File | Responsibility |
|------|---------------|
| `internal/cli/application/agents/agents.go` | NEW: Cmd group |
| `internal/cli/application/agents/list.go` | NEW: agents list |
| `internal/cli/application/agents/create.go` | NEW: agents create |
| `internal/cli/application/agents/update.go` | NEW: agents update |
| `internal/cli/application/agents/delete.go` | NEW: agents delete |
| `internal/cli/application/agents/agents_test.go` | NEW: tests |
| `internal/cli/application/agents/testroot_test.go` | NEW: test helper |

### CLI — Canned Responses Commands

| File | Responsibility |
|------|---------------|
| `internal/cli/application/cannedresponses/cannedresponses.go` | NEW: Cmd group (Use: "canned-responses") |
| `internal/cli/application/cannedresponses/list.go` | NEW: canned-responses list |
| `internal/cli/application/cannedresponses/create.go` | NEW: canned-responses create |
| `internal/cli/application/cannedresponses/update.go` | NEW: canned-responses update |
| `internal/cli/application/cannedresponses/delete.go` | NEW: canned-responses delete |
| `internal/cli/application/cannedresponses/cannedresponses_test.go` | NEW: tests |
| `internal/cli/application/cannedresponses/testroot_test.go` | NEW: test helper |

### CLI — Reports Commands

| File | Responsibility |
|------|---------------|
| `internal/cli/application/reports/reports.go` | NEW: Cmd group |
| `internal/cli/application/reports/account.go` | NEW: reports account |
| `internal/cli/application/reports/summary.go` | NEW: reports summary |
| `internal/cli/application/reports/conversations.go` | NEW: reports conversations |
| `internal/cli/application/reports/first_response.go` | NEW: reports first-response-distribution |
| `internal/cli/application/reports/inbox_label.go` | NEW: reports inbox-label-matrix |
| `internal/cli/application/reports/outgoing.go` | NEW: reports outgoing-messages |
| `internal/cli/application/reports/summary_agent.go` | NEW: reports summary-by-agent |
| `internal/cli/application/reports/summary_channel.go` | NEW: reports summary-by-channel |
| `internal/cli/application/reports/summary_inbox.go` | NEW: reports summary-by-inbox |
| `internal/cli/application/reports/summary_team.go` | NEW: reports summary-by-team |
| `internal/cli/application/reports/events.go` | NEW: reports events |
| `internal/cli/application/reports/reports_test.go` | NEW: tests |
| `internal/cli/application/reports/testroot_test.go` | NEW: test helper |

### CLI — Webhooks Commands

| File | Responsibility |
|------|---------------|
| `internal/cli/application/webhooks/webhooks.go` | NEW: Cmd group |
| `internal/cli/application/webhooks/list.go` | NEW: webhooks list |
| `internal/cli/application/webhooks/create.go` | NEW: webhooks create |
| `internal/cli/application/webhooks/update.go` | NEW: webhooks update |
| `internal/cli/application/webhooks/delete.go` | NEW: webhooks delete |
| `internal/cli/application/webhooks/webhooks_test.go` | NEW: tests |
| `internal/cli/application/webhooks/testroot_test.go` | NEW: test helper |

### CLI — Automation Rules Commands

| File | Responsibility |
|------|---------------|
| `internal/cli/application/automationrules/automationrules.go` | NEW: Cmd group (Use: "automation-rules") |
| `internal/cli/application/automationrules/list.go` | NEW: automation-rules list |
| `internal/cli/application/automationrules/get.go` | NEW: automation-rules get |
| `internal/cli/application/automationrules/create.go` | NEW: automation-rules create |
| `internal/cli/application/automationrules/update.go` | NEW: automation-rules update |
| `internal/cli/application/automationrules/delete.go` | NEW: automation-rules delete |
| `internal/cli/application/automationrules/automationrules_test.go` | NEW: tests |
| `internal/cli/application/automationrules/testroot_test.go` | NEW: test helper |

### CLI — Simple CRUD Batch

| File | Responsibility |
|------|---------------|
| `internal/cli/application/labels/labels.go` | NEW: Cmd group |
| `internal/cli/application/labels/list.go` | NEW: labels list |
| `internal/cli/application/labels/get.go` | NEW: labels get |
| `internal/cli/application/labels/create.go` | NEW: labels create |
| `internal/cli/application/labels/update.go` | NEW: labels update |
| `internal/cli/application/labels/delete.go` | NEW: labels delete |
| `internal/cli/application/labels/labels_test.go` | NEW: tests |
| `internal/cli/application/labels/testroot_test.go` | NEW: test helper |
| `internal/cli/application/customattributes/customattributes.go` | NEW: Cmd group (Use: "custom-attributes") |
| `internal/cli/application/customattributes/list.go` | NEW: custom-attributes list |
| `internal/cli/application/customattributes/get.go` | NEW: custom-attributes get |
| `internal/cli/application/customattributes/create.go` | NEW: custom-attributes create |
| `internal/cli/application/customattributes/update.go` | NEW: custom-attributes update |
| `internal/cli/application/customattributes/delete.go` | NEW: custom-attributes delete |
| `internal/cli/application/customattributes/customattributes_test.go` | NEW: tests |
| `internal/cli/application/customattributes/testroot_test.go` | NEW: test helper |
| `internal/cli/application/customfilters/customfilters.go` | NEW: Cmd group (Use: "custom-filters") |
| `internal/cli/application/customfilters/list.go` | NEW: custom-filters list |
| `internal/cli/application/customfilters/get.go` | NEW: custom-filters get |
| `internal/cli/application/customfilters/create.go` | NEW: custom-filters create |
| `internal/cli/application/customfilters/update.go` | NEW: custom-filters update |
| `internal/cli/application/customfilters/delete.go` | NEW: custom-filters delete |
| `internal/cli/application/customfilters/customfilters_test.go` | NEW: tests |
| `internal/cli/application/customfilters/testroot_test.go` | NEW: test helper |
| `internal/cli/application/account/account.go` | NEW: Cmd group |
| `internal/cli/application/account/get.go` | NEW: account get |
| `internal/cli/application/account/update.go` | NEW: account update |
| `internal/cli/application/account/account_test.go` | NEW: tests |
| `internal/cli/application/account/testroot_test.go` | NEW: test helper |
| `internal/cli/application/agentbots/agentbots.go` | NEW: Cmd group (Use: "agent-bots") |
| `internal/cli/application/agentbots/list.go` | NEW: agent-bots list |
| `internal/cli/application/agentbots/get.go` | NEW: agent-bots get |
| `internal/cli/application/agentbots/create.go` | NEW: agent-bots create |
| `internal/cli/application/agentbots/update.go` | NEW: agent-bots update |
| `internal/cli/application/agentbots/delete.go` | NEW: agent-bots delete |
| `internal/cli/application/agentbots/agentbots_test.go` | NEW: tests |
| `internal/cli/application/agentbots/testroot_test.go` | NEW: test helper |
| `internal/cli/application/auditlogs/auditlogs.go` | NEW: Cmd group (Use: "audit-logs") |
| `internal/cli/application/auditlogs/list.go` | NEW: audit-logs list |
| `internal/cli/application/auditlogs/auditlogs_test.go` | NEW: tests |
| `internal/cli/application/auditlogs/testroot_test.go` | NEW: test helper |

### CLI — Integrations + Help Center

| File | Responsibility |
|------|---------------|
| `internal/cli/application/integrations/integrations.go` | NEW: Cmd group with appsCmd and hooksCmd subgroups |
| `internal/cli/application/integrations/apps.go` | NEW: integrations apps list |
| `internal/cli/application/integrations/hooks.go` | NEW: integrations hooks create/update/delete |
| `internal/cli/application/integrations/integrations_test.go` | NEW: tests |
| `internal/cli/application/integrations/testroot_test.go` | NEW: test helper |
| `internal/cli/application/helpcenter/helpcenter.go` | NEW: Cmd group (Use: "help-center") with portalsCmd, articlesCmd, categoriesCmd |
| `internal/cli/application/helpcenter/portals.go` | NEW: portals list/create/update |
| `internal/cli/application/helpcenter/articles.go` | NEW: articles create |
| `internal/cli/application/helpcenter/categories.go` | NEW: categories create |
| `internal/cli/application/helpcenter/helpcenter_test.go` | NEW: tests |
| `internal/cli/application/helpcenter/testroot_test.go` | NEW: test helper |

### Registration

| File | Responsibility |
|------|---------------|
| `internal/cli/application/application.go` | MODIFY: register all 14 new command subgroups |

---

## Task 1: Add All Model Types

**Files:**
- Modify: `internal/chatwoot/application/models.go`

- [ ] **Step 1: Add all core types to models.go**

Add the following types after the existing types in `models.go`:

```go
// Team represents a Chatwoot team.
type Team struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description,omitempty"`
	AllowAutoAssign bool   `json:"allow_auto_assign"`
	AccountID       int    `json:"account_id"`
	IsMember        bool   `json:"is_member,omitempty"`
}

// CannedResponse represents a saved reply template.
type CannedResponse struct {
	ID        int    `json:"id"`
	ShortCode string `json:"short_code"`
	Content   string `json:"content"`
	AccountID int    `json:"account_id,omitempty"`
}

// Webhook represents a webhook subscription.
type Webhook struct {
	ID            int      `json:"id"`
	URL           string   `json:"url"`
	Subscriptions []string `json:"subscriptions"`
	AccountID     int      `json:"account_id,omitempty"`
}

// AutomationRule represents an automation rule.
type AutomationRule struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	EventName   string `json:"event_name"`
	Conditions  any    `json:"conditions"`
	Actions     any    `json:"actions"`
	AccountID   int    `json:"account_id,omitempty"`
}

// Label represents an account-level label definition.
type Label struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description,omitempty"`
	Color         string `json:"color,omitempty"`
	ShowOnSidebar bool   `json:"show_on_sidebar,omitempty"`
}

// CustomAttribute represents a custom attribute definition.
type CustomAttribute struct {
	ID                   int    `json:"id"`
	AttributeDisplayName string `json:"attribute_display_name"`
	AttributeDisplayType string `json:"attribute_display_type"`
	AttributeDescription string `json:"attribute_description,omitempty"`
	AttributeKey         string `json:"attribute_key"`
	AttributeModel       string `json:"attribute_model"`
}

// CustomFilter represents a saved filter query.
type CustomFilter struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Query any    `json:"query"`
}

// AccountInfo represents account details.
type AccountInfo struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Locale           string `json:"locale,omitempty"`
	Domain           string `json:"domain,omitempty"`
	CustomAttributes any    `json:"custom_attributes,omitempty"`
}

// AccountAgentBot represents an account-scoped agent bot.
type AccountAgentBot struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	BotType     string `json:"bot_type,omitempty"`
	OutgoingURL string `json:"outgoing_url,omitempty"`
	BotConfig   any    `json:"bot_config,omitempty"`
	AccountID   int    `json:"account_id,omitempty"`
}

// AuditLog represents an audit log entry.
type AuditLog struct {
	ID            int              `json:"id"`
	Action        string           `json:"action"`
	AuditableType string           `json:"auditable_type"`
	AuditableID   int              `json:"auditable_id"`
	UserID        int              `json:"user_id,omitempty"`
	CreatedAt     chatwoot.Timestamp `json:"created_at,omitempty"`
}

// Integration represents an available integration app.
type Integration struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Hooks []any  `json:"hooks,omitempty"`
}

// IntegrationHook represents an activated integration hook.
type IntegrationHook struct {
	ID       int    `json:"id"`
	AppID    string `json:"app_id"`
	InboxID  int    `json:"inbox_id,omitempty"`
	Status   int    `json:"status,omitempty"`
	Settings any    `json:"settings,omitempty"`
}

// Portal represents a help center portal.
type Portal struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug,omitempty"`
	Color        string `json:"color,omitempty"`
	HeaderText   string `json:"header_text,omitempty"`
	CustomDomain string `json:"custom_domain,omitempty"`
	Archived     bool   `json:"archived,omitempty"`
}

// Article represents a help center article.
type Article struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug,omitempty"`
	Content     string `json:"content,omitempty"`
	Description string `json:"description,omitempty"`
	Status      int    `json:"status,omitempty"`
	CategoryID  int    `json:"category_id,omitempty"`
	AuthorID    int    `json:"author_id,omitempty"`
}

// Category represents a help center category.
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Locale      string `json:"locale,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Position    int    `json:"position,omitempty"`
}

// ReportSummary holds account report summary data.
type ReportSummary struct {
	AvgFirstResponseTime  string `json:"avg_first_response_time"`
	AvgResolutionTime     string `json:"avg_resolution_time"`
	ConversationsCount    int    `json:"conversations_count"`
	IncomingMessagesCount int    `json:"incoming_messages_count"`
	OutgoingMessagesCount int    `json:"outgoing_messages_count"`
	ResolutionsCount      int    `json:"resolutions_count"`
	Previous              any    `json:"previous,omitempty"`
}
```

- [ ] **Step 2: Add all opts types to models.go**

Add the following opts types after the existing opts section:

```go
// --- Sprint D Opts types ---

type CreateTeamOpts struct {
	Name            string `json:"name"`
	Description     string `json:"description,omitempty"`
	AllowAutoAssign *bool  `json:"allow_auto_assign,omitempty"`
}

type UpdateTeamOpts struct {
	Name            *string `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	AllowAutoAssign *bool   `json:"allow_auto_assign,omitempty"`
}

type CreateAgentOpts struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role,omitempty"`
}

type UpdateAgentOpts struct {
	Name *string `json:"name,omitempty"`
	Role *string `json:"role,omitempty"`
}

type CreateCannedResponseOpts struct {
	ShortCode string `json:"short_code"`
	Content   string `json:"content"`
}

type UpdateCannedResponseOpts struct {
	ShortCode *string `json:"short_code,omitempty"`
	Content   *string `json:"content,omitempty"`
}

type ReportOpts struct {
	Metric string
	Type   string
	ID     string
	Since  string
	Until  string
}

type CreateWebhookOpts struct {
	URL           string   `json:"url"`
	Subscriptions []string `json:"subscriptions,omitempty"`
}

type UpdateWebhookOpts struct {
	URL           *string  `json:"url,omitempty"`
	Subscriptions []string `json:"subscriptions,omitempty"`
}

type CreateAutomationRuleOpts struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	EventName   string `json:"event_name"`
	Conditions  any    `json:"conditions"`
	Actions     any    `json:"actions"`
}

type UpdateAutomationRuleOpts struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	EventName   *string `json:"event_name,omitempty"`
	Conditions  any     `json:"conditions,omitempty"`
	Actions     any     `json:"actions,omitempty"`
}

type CreateLabelOpts struct {
	Title         string `json:"title"`
	Description   string `json:"description,omitempty"`
	Color         string `json:"color,omitempty"`
	ShowOnSidebar *bool  `json:"show_on_sidebar,omitempty"`
}

type UpdateLabelOpts struct {
	Title         *string `json:"title,omitempty"`
	Description   *string `json:"description,omitempty"`
	Color         *string `json:"color,omitempty"`
	ShowOnSidebar *bool   `json:"show_on_sidebar,omitempty"`
}

type CreateCustomAttributeOpts struct {
	AttributeDisplayName string `json:"attribute_display_name"`
	AttributeKey         string `json:"attribute_key"`
	AttributeModel       string `json:"attribute_model"`
	AttributeDisplayType string `json:"attribute_display_type"`
	AttributeDescription string `json:"attribute_description,omitempty"`
}

type UpdateCustomAttributeOpts struct {
	AttributeDisplayName *string `json:"attribute_display_name,omitempty"`
	AttributeDescription *string `json:"attribute_description,omitempty"`
}

type CreateCustomFilterOpts struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Query any    `json:"query"`
}

type UpdateCustomFilterOpts struct {
	Name  *string `json:"name,omitempty"`
	Query any     `json:"query,omitempty"`
}

type UpdateAccountOpts struct {
	Name   *string `json:"name,omitempty"`
	Locale *string `json:"locale,omitempty"`
	Domain *string `json:"domain,omitempty"`
}

type CreateAgentBotOpts struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	BotType     string `json:"bot_type,omitempty"`
	OutgoingURL string `json:"outgoing_url,omitempty"`
	BotConfig   any    `json:"bot_config,omitempty"`
}

type UpdateAgentBotOpts struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	BotType     *string `json:"bot_type,omitempty"`
	OutgoingURL *string `json:"outgoing_url,omitempty"`
	BotConfig   any     `json:"bot_config,omitempty"`
}

type CreateIntegrationHookOpts struct {
	AppID    string `json:"app_id"`
	InboxID  int    `json:"inbox_id,omitempty"`
	Settings any    `json:"settings,omitempty"`
}

type UpdateIntegrationHookOpts struct {
	Settings any `json:"settings,omitempty"`
}

type CreatePortalOpts struct {
	Name         string `json:"name"`
	Slug         string `json:"slug,omitempty"`
	Color        string `json:"color,omitempty"`
	HeaderText   string `json:"header_text,omitempty"`
	CustomDomain string `json:"custom_domain,omitempty"`
}

type UpdatePortalOpts struct {
	Name         *string `json:"name,omitempty"`
	Slug         *string `json:"slug,omitempty"`
	Color        *string `json:"color,omitempty"`
	HeaderText   *string `json:"header_text,omitempty"`
	CustomDomain *string `json:"custom_domain,omitempty"`
	Archived     *bool   `json:"archived,omitempty"`
}

type CreateArticleOpts struct {
	Title       string `json:"title"`
	Content     string `json:"content,omitempty"`
	Description string `json:"description,omitempty"`
	Status      int    `json:"status,omitempty"`
	CategoryID  int    `json:"category_id,omitempty"`
	AuthorID    int    `json:"author_id,omitempty"`
}

type CreateCategoryOpts struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Locale      string `json:"locale,omitempty"`
	Position    int    `json:"position,omitempty"`
}
```

- [ ] **Step 3: Verify build**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./cmd/chatwoot/`
Expected: Success

- [ ] **Step 4: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/application/models.go
git commit -m "feat(application): add Sprint D model types for 14 resource families"
```

---

## Task 2: Teams API Client

**Files:**
- Create: `internal/chatwoot/application/teams.go`
- Create: `internal/chatwoot/application/teams_test.go`

- [ ] **Step 1: Write failing test for ListTeams**

Create `internal/chatwoot/application/teams_test.go`:

```go
package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func TestListTeams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/teams" {
			t.Errorf("path = %q, want /api/v1/accounts/1/teams", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "Support", "account_id": 1},
			{"id": 2, "name": "Sales", "account_id": 1},
		})
	}))
	defer srv.Close()

	transport := chatwoot.NewClient(srv.URL, "sk-test", "api_access_token")
	client := NewClient(transport, 1)

	teams, err := client.ListTeams(context.Background())
	if err != nil {
		t.Fatalf("ListTeams error: %v", err)
	}
	if len(teams) != 2 {
		t.Errorf("len = %d, want 2", len(teams))
	}
	if teams[0].Name != "Support" {
		t.Errorf("teams[0].Name = %q, want Support", teams[0].Name)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestListTeams`
Expected: FAIL — `ListTeams` not defined

- [ ] **Step 3: Implement all 9 team API methods in teams.go**

Create `internal/chatwoot/application/teams.go` with all methods from the spec:

**Response decode patterns:**
- `ListTeams`: Direct decode into `[]Team` (no wrapper — Chatwoot returns a plain JSON array for teams list)
- `GetTeam`, `CreateTeam`, `UpdateTeam`: Direct decode into `Team`
- `DeleteTeam`: Decode with `nil` target
- `ListTeamMembers`, `AddTeamMember`, `UpdateTeamMembers`: Direct decode into `[]Agent` (plain array)
- `RemoveTeamMember`: Decode with `nil` target

**HTTP methods and paths:**
- `ListTeams(ctx)` — GET `/api/v1/accounts/{id}/teams`
- `GetTeam(ctx, teamID)` — GET `/api/v1/accounts/{id}/teams/{tid}`
- `CreateTeam(ctx, CreateTeamOpts)` — POST `/api/v1/accounts/{id}/teams` (marshal opts to JSON body)
- `UpdateTeam(ctx, teamID, UpdateTeamOpts)` — PATCH `/api/v1/accounts/{id}/teams/{tid}` (marshal opts to JSON body)
- `DeleteTeam(ctx, teamID)` — DELETE `/api/v1/accounts/{id}/teams/{tid}`
- `ListTeamMembers(ctx, teamID)` — GET `/api/v1/accounts/{id}/teams/{tid}/team_members`
- `AddTeamMember(ctx, teamID, agentIDs)` — POST `/api/v1/accounts/{id}/teams/{tid}/team_members` (body: `{"user_ids": [1, 2, 3]}`)
- `UpdateTeamMembers(ctx, teamID, agentIDs)` — PATCH `/api/v1/accounts/{id}/teams/{tid}/team_members` (body: `{"user_ids": [1, 2, 3]}`)
- `RemoveTeamMember(ctx, teamID, agentIDs)` — DELETE `/api/v1/accounts/{id}/teams/{tid}/team_members` (body: `{"user_ids": [1, 2, 3]}`)

- [ ] **Step 4: Add remaining tests**

Add tests for `GetTeam`, `CreateTeam`, `DeleteTeam`, `ListTeamMembers`, `AddTeamMember`, `RemoveTeamMember`. Each verifies correct HTTP method, URL path, request body (for mutations), and response deserialization.

- [ ] **Step 5: Run all team tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestTeam`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/application/teams.go internal/chatwoot/application/teams_test.go
git commit -m "feat(application): add teams API client with 9 methods"
```

---

## Task 3: Agents API Client

**Files:**
- Create: `internal/chatwoot/application/agents.go`
- Create: `internal/chatwoot/application/agents_test.go`

- [ ] **Step 1: Write failing test for ListAgents**

Create `internal/chatwoot/application/agents_test.go` with `TestListAgents`. The httptest server returns a plain JSON array `[{"id": 1, "name": "Alice", "email": "alice@test.com"}, ...]`.

**Response decode pattern:** `ListAgents` returns a plain JSON array — decode directly into `[]Agent` (no wrapper).

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestListAgents`
Expected: FAIL

- [ ] **Step 3: Implement all 4 agent methods**

Create `internal/chatwoot/application/agents.go`:

**HTTP methods and paths:**
- `ListAgents(ctx)` — GET `/api/v1/accounts/{id}/agents` — decode into `[]Agent`
- `CreateAgent(ctx, CreateAgentOpts)` — POST `/api/v1/accounts/{id}/agents` — marshal opts, decode into `Agent`
- `UpdateAgent(ctx, agentID, UpdateAgentOpts)` — PATCH `/api/v1/accounts/{id}/agents/{aid}` — marshal opts, decode into `Agent`
- `DeleteAgent(ctx, agentID)` — DELETE `/api/v1/accounts/{id}/agents/{aid}` — decode with `nil`

- [ ] **Step 4: Add tests for CreateAgent and DeleteAgent**

- [ ] **Step 5: Run all agent tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestAgent`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/application/agents.go internal/chatwoot/application/agents_test.go
git commit -m "feat(application): add agents API client with 4 methods"
```

---

## Task 4: Canned Responses API Client

**Files:**
- Create: `internal/chatwoot/application/canned_responses.go`
- Create: `internal/chatwoot/application/canned_responses_test.go`

- [ ] **Step 1: Write failing test for ListCannedResponses**

Create `internal/chatwoot/application/canned_responses_test.go` with `TestListCannedResponses`. The httptest server returns a plain JSON array `[{"id": 1, "short_code": "greeting", "content": "Hello!"}, ...]`.

**Response decode pattern:** Direct decode into `[]CannedResponse` (plain array, no wrapper).

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestListCannedResponses`
Expected: FAIL

- [ ] **Step 3: Implement all 4 canned response methods**

**HTTP methods and paths:**
- `ListCannedResponses(ctx)` — GET `/api/v1/accounts/{id}/canned_responses` — decode into `[]CannedResponse`
- `CreateCannedResponse(ctx, CreateCannedResponseOpts)` — POST `/api/v1/accounts/{id}/canned_responses` — marshal opts, decode into `CannedResponse`
- `UpdateCannedResponse(ctx, cannedID, UpdateCannedResponseOpts)` — PATCH `/api/v1/accounts/{id}/canned_responses/{cid}` — marshal opts, decode into `CannedResponse`
- `DeleteCannedResponse(ctx, cannedID)` — DELETE `/api/v1/accounts/{id}/canned_responses/{cid}` — decode with `nil`

- [ ] **Step 4: Add tests for CreateCannedResponse and DeleteCannedResponse**

- [ ] **Step 5: Run all canned response tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestCannedResponse`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/application/canned_responses.go internal/chatwoot/application/canned_responses_test.go
git commit -m "feat(application): add canned responses API client with 4 methods"
```

---

## Task 5: Reports API Client

**Files:**
- Create: `internal/chatwoot/application/reports.go`
- Create: `internal/chatwoot/application/reports_test.go`

- [ ] **Step 1: Write failing test for GetReports**

Create `internal/chatwoot/application/reports_test.go` with `TestGetReports`. The httptest server verifies query parameters (`metric`, `type`, `since`, `until`) and returns a JSON array (raw data).

**Key implementation note:** Report methods build query strings from `ReportOpts` fields. The `ReportOpts` struct has no JSON tags because values are passed as URL query parameters, not request body.

Query parameter construction pattern for all report methods:
```go
func (c *Client) buildReportQuery(opts ReportOpts) string {
	params := url.Values{}
	if opts.Metric != "" {
		params.Set("metric", opts.Metric)
	}
	if opts.Type != "" {
		params.Set("type", opts.Type)
	}
	if opts.ID != "" {
		params.Set("id", opts.ID)
	}
	if opts.Since != "" {
		params.Set("since", opts.Since)
	}
	if opts.Until != "" {
		params.Set("until", opts.Until)
	}
	if len(params) > 0 {
		return "?" + params.Encode()
	}
	return ""
}
```

**Response decode patterns:** Most report endpoints return varied JSON shapes. Use `json.RawMessage` to decode, then convert to `any`:
```go
var raw json.RawMessage
if err := chatwoot.DecodeResponse(resp, &raw); err != nil {
	return nil, err
}
var result any
if err := json.Unmarshal(raw, &result); err != nil {
	return nil, fmt.Errorf("decode report: %w", err)
}
return result, nil
```

`GetReportSummary` is the exception — it decodes into `*ReportSummary` directly.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestGetReports`
Expected: FAIL

- [ ] **Step 3: Implement all 12 report methods**

Create `internal/chatwoot/application/reports.go`:

**HTTP methods and paths (all GET):**
- `GetReports(ctx, ReportOpts)` — `/api/v2/accounts/{id}/reports` — returns `(any, error)`
- `GetReportSummary(ctx, ReportOpts)` — `/api/v2/accounts/{id}/reports/summary` — returns `(*ReportSummary, error)`
- `GetConversationMetrics(ctx, ReportOpts)` — `/api/v2/accounts/{id}/reports/conversations` — returns `(any, error)`
- `GetAgentConversationMetrics(ctx, ReportOpts)` — `/api/v2/accounts/{id}/reports/conversations/` (trailing slash) — returns `(any, error)`
- `GetFirstResponseTimeDistribution(ctx, ReportOpts)` — `/api/v2/accounts/{id}/reports/first_response_time_distribution` — returns `(any, error)`
- `GetInboxLabelMatrix(ctx, ReportOpts)` — `/api/v2/accounts/{id}/reports/inbox_label_matrix` — returns `(any, error)`
- `GetOutgoingMessagesCount(ctx, ReportOpts)` — `/api/v2/accounts/{id}/reports/outgoing_messages_count` — returns `(any, error)`
- `GetSummaryByAgent(ctx, ReportOpts)` — `/api/v2/accounts/{id}/summary_reports/agent` — returns `(any, error)`
- `GetSummaryByChannel(ctx, ReportOpts)` — `/api/v2/accounts/{id}/summary_reports/channel` — returns `(any, error)`
- `GetSummaryByInbox(ctx, ReportOpts)` — `/api/v2/accounts/{id}/summary_reports/inbox` — returns `(any, error)`
- `GetSummaryByTeam(ctx, ReportOpts)` — `/api/v2/accounts/{id}/summary_reports/team` — returns `(any, error)`
- `GetReportingEvents(ctx, ReportOpts)` — `/api/v1/accounts/{id}/reporting_events` — returns `(any, error)`

Add a helper method `buildReportQuery(opts ReportOpts) string` used by all methods.

- [ ] **Step 4: Add tests for GetReportSummary and GetSummaryByAgent**

Test `GetReportSummary` verifies:
- Correct path `/api/v2/accounts/1/reports/summary`
- Query params include `type`, `since`, `until`
- Response decodes into `ReportSummary` struct

Test `GetSummaryByAgent` verifies:
- Correct path `/api/v2/accounts/1/summary_reports/agent`
- Returns decoded `any` (can be a map or array)

- [ ] **Step 5: Run all report tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestReport`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/application/reports.go internal/chatwoot/application/reports_test.go
git commit -m "feat(application): add reports API client with 12 methods"
```

---

## Task 6: Webhooks API Client

**Files:**
- Create: `internal/chatwoot/application/webhooks.go`
- Create: `internal/chatwoot/application/webhooks_test.go`

- [ ] **Step 1: Write failing test for ListWebhooks**

Create `internal/chatwoot/application/webhooks_test.go` with `TestListWebhooks`. The httptest server returns:
```json
{"payload": [{"id": 1, "url": "https://example.com/hook", "subscriptions": ["message_created"]}]}
```

**Response decode pattern:** `ListWebhooks` uses a payload wrapper: `var body struct { Payload []Webhook \x60json:"payload"\x60 }`.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestListWebhooks`
Expected: FAIL

- [ ] **Step 3: Implement all 4 webhook methods**

**HTTP methods and paths:**
- `ListWebhooks(ctx)` — GET `/api/v1/accounts/{id}/webhooks` — decode with payload wrapper: `var body struct { Payload []Webhook \x60json:"payload"\x60 }`
- `CreateWebhook(ctx, CreateWebhookOpts)` — POST `/api/v1/accounts/{id}/webhooks` — marshal opts, decode into `Webhook`
- `UpdateWebhook(ctx, webhookID, UpdateWebhookOpts)` — PATCH `/api/v1/accounts/{id}/webhooks/{wid}` — marshal opts, decode into `Webhook`
- `DeleteWebhook(ctx, webhookID)` — DELETE `/api/v1/accounts/{id}/webhooks/{wid}` — decode with `nil`

- [ ] **Step 4: Add tests for CreateWebhook and DeleteWebhook**

Test `CreateWebhook` verifies the request body includes `url` and `subscriptions` fields.

- [ ] **Step 5: Run all webhook tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestWebhook`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/application/webhooks.go internal/chatwoot/application/webhooks_test.go
git commit -m "feat(application): add webhooks API client with 4 methods"
```

---

## Task 7: Automation Rules API Client

**Files:**
- Create: `internal/chatwoot/application/automation_rules.go`
- Create: `internal/chatwoot/application/automation_rules_test.go`

- [ ] **Step 1: Write failing test for ListAutomationRules**

Create `internal/chatwoot/application/automation_rules_test.go` with `TestListAutomationRules`. The httptest server returns:
```json
{"payload": [{"id": 1, "name": "Auto-assign", "event_name": "conversation_created"}]}
```

**Response decode pattern:** `ListAutomationRules` uses a payload wrapper: `var body struct { Payload []AutomationRule \x60json:"payload"\x60 }`.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestListAutomationRules`
Expected: FAIL

- [ ] **Step 3: Implement all 5 automation rule methods**

**HTTP methods and paths:**
- `ListAutomationRules(ctx)` — GET `/api/v1/accounts/{id}/automation_rules` — decode with payload wrapper
- `GetAutomationRule(ctx, ruleID)` — GET `/api/v1/accounts/{id}/automation_rules/{rid}` — decode into `AutomationRule`
- `CreateAutomationRule(ctx, CreateAutomationRuleOpts)` — POST `/api/v1/accounts/{id}/automation_rules` — marshal opts, decode into `AutomationRule`
- `UpdateAutomationRule(ctx, ruleID, UpdateAutomationRuleOpts)` — PATCH `/api/v1/accounts/{id}/automation_rules/{rid}` — marshal opts, decode into `AutomationRule`
- `DeleteAutomationRule(ctx, ruleID)` — DELETE `/api/v1/accounts/{id}/automation_rules/{rid}` — decode with `nil`

- [ ] **Step 4: Add tests for GetAutomationRule, CreateAutomationRule, DeleteAutomationRule**

- [ ] **Step 5: Run all automation rule tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run TestAutomationRule`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/application/automation_rules.go internal/chatwoot/application/automation_rules_test.go
git commit -m "feat(application): add automation rules API client with 5 methods"
```

---

## Task 8: Simple CRUD API Clients Batch

**Files:**
- Create: `internal/chatwoot/application/labels.go`
- Create: `internal/chatwoot/application/labels_test.go`
- Create: `internal/chatwoot/application/custom_attributes.go`
- Create: `internal/chatwoot/application/custom_attributes_test.go`
- Create: `internal/chatwoot/application/custom_filters.go`
- Create: `internal/chatwoot/application/custom_filters_test.go`
- Create: `internal/chatwoot/application/account.go`
- Create: `internal/chatwoot/application/account_test.go`
- Create: `internal/chatwoot/application/agent_bots.go`
- Create: `internal/chatwoot/application/agent_bots_test.go`
- Create: `internal/chatwoot/application/audit_logs.go`
- Create: `internal/chatwoot/application/audit_logs_test.go`

All six resources follow the same CRUD pattern as contacts.go. Each resource gets a separate client file and test file.

### Labels (`labels.go`)

- [ ] **Step 1: Implement labels.go with 5 methods**

**HTTP methods and paths:**
- `ListLabels(ctx)` — GET `/api/v1/accounts/{id}/labels` — decode with payload wrapper: `var body struct { Payload []Label \x60json:"payload"\x60 }`
- `GetLabel(ctx, labelID)` — GET `/api/v1/accounts/{id}/labels/{lid}` — decode into `Label`
- `CreateLabel(ctx, CreateLabelOpts)` — POST `/api/v1/accounts/{id}/labels` — marshal opts, decode into `Label`
- `UpdateLabel(ctx, labelID, UpdateLabelOpts)` — PATCH `/api/v1/accounts/{id}/labels/{lid}` — marshal opts, decode into `Label`
- `DeleteLabel(ctx, labelID)` — DELETE `/api/v1/accounts/{id}/labels/{lid}` — decode with `nil`

- [ ] **Step 2: Write labels_test.go with TestListLabels, TestCreateLabel, TestDeleteLabel**

### Custom Attributes (`custom_attributes.go`)

- [ ] **Step 3: Implement custom_attributes.go with 5 methods**

**HTTP methods and paths:**
- `ListCustomAttributes(ctx)` — GET `/api/v1/accounts/{id}/custom_attribute_definitions` — decode with payload wrapper: `var body struct { Data []CustomAttribute \x60json:"data"\x60 }`
- `GetCustomAttribute(ctx, attrID)` — GET `/api/v1/accounts/{id}/custom_attribute_definitions/{aid}` — decode into `CustomAttribute`
- `CreateCustomAttribute(ctx, CreateCustomAttributeOpts)` — POST `/api/v1/accounts/{id}/custom_attribute_definitions` — marshal opts, decode into `CustomAttribute`
- `UpdateCustomAttribute(ctx, attrID, UpdateCustomAttributeOpts)` — PATCH `/api/v1/accounts/{id}/custom_attribute_definitions/{aid}` — marshal opts, decode into `CustomAttribute`
- `DeleteCustomAttribute(ctx, attrID)` — DELETE `/api/v1/accounts/{id}/custom_attribute_definitions/{aid}` — decode with `nil`

**Note:** The API path uses `custom_attribute_definitions` (not `custom_attributes`).

- [ ] **Step 4: Write custom_attributes_test.go with TestListCustomAttributes, TestCreateCustomAttribute**

### Custom Filters (`custom_filters.go`)

- [ ] **Step 5: Implement custom_filters.go with 5 methods**

**HTTP methods and paths:**
- `ListCustomFilters(ctx, filterType)` — GET `/api/v1/accounts/{id}/custom_filters?filter_type={filterType}` — decode directly into `[]CustomFilter`
- `GetCustomFilter(ctx, filterID)` — GET `/api/v1/accounts/{id}/custom_filters/{fid}` — decode into `CustomFilter`
- `CreateCustomFilter(ctx, CreateCustomFilterOpts)` — POST `/api/v1/accounts/{id}/custom_filters` — marshal opts, decode into `CustomFilter`
- `UpdateCustomFilter(ctx, filterID, UpdateCustomFilterOpts)` — PATCH `/api/v1/accounts/{id}/custom_filters/{fid}` — marshal opts, decode into `CustomFilter`
- `DeleteCustomFilter(ctx, filterID)` — DELETE `/api/v1/accounts/{id}/custom_filters/{fid}` — decode with `nil`

**Note:** `ListCustomFilters` takes a `filterType` string parameter appended as query param. The response is a plain JSON array (no wrapper).

- [ ] **Step 6: Write custom_filters_test.go with TestListCustomFilters, TestCreateCustomFilter**

### Account (`account.go`)

- [ ] **Step 7: Implement account.go with 2 methods**

**HTTP methods and paths:**
- `GetAccount(ctx)` — GET `/api/v1/accounts/{id}` — decode into `AccountInfo`
- `UpdateAccount(ctx, UpdateAccountOpts)` — PATCH `/api/v1/accounts/{id}` — marshal opts, decode into `AccountInfo`

**Note:** No resource ID parameter — uses `c.accountID` from the client.

- [ ] **Step 8: Write account_test.go with TestGetAccount, TestUpdateAccount**

### Agent Bots (`agent_bots.go`)

- [ ] **Step 9: Implement agent_bots.go with 5 methods**

**HTTP methods and paths:**
- `ListAgentBots(ctx)` — GET `/api/v1/accounts/{id}/agent_bots` — decode directly into `[]AccountAgentBot`
- `GetAgentBot(ctx, botID)` — GET `/api/v1/accounts/{id}/agent_bots/{bid}` — decode into `AccountAgentBot`
- `CreateAgentBot(ctx, CreateAgentBotOpts)` — POST `/api/v1/accounts/{id}/agent_bots` — marshal opts, decode into `AccountAgentBot`
- `UpdateAgentBot(ctx, botID, UpdateAgentBotOpts)` — PATCH `/api/v1/accounts/{id}/agent_bots/{bid}` — marshal opts, decode into `AccountAgentBot`
- `DeleteAgentBot(ctx, botID)` — DELETE `/api/v1/accounts/{id}/agent_bots/{bid}` — decode with `nil`

**Note:** Uses `AccountAgentBot` type (not the existing `AgentBot` which is a simpler struct for inbox bot assignment).

- [ ] **Step 10: Write agent_bots_test.go with TestListAgentBots, TestCreateAgentBot**

### Audit Logs (`audit_logs.go`)

- [ ] **Step 11: Implement audit_logs.go with 1 method**

**HTTP methods and paths:**
- `ListAuditLogs(ctx, page)` — GET `/api/v1/accounts/{id}/audit_logs?page={page}` — decode with payload wrapper: `var body struct { Payload []AuditLog \x60json:"payload"\x60 }`

- [ ] **Step 12: Write audit_logs_test.go with TestListAuditLogs**

### Run and Commit

- [ ] **Step 13: Run all tests for the batch**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run "TestLabel|TestCustomAttribute|TestCustomFilter|TestAccount|TestAgentBot|TestAuditLog"`
Expected: All PASS

- [ ] **Step 14: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/application/labels.go internal/chatwoot/application/labels_test.go internal/chatwoot/application/custom_attributes.go internal/chatwoot/application/custom_attributes_test.go internal/chatwoot/application/custom_filters.go internal/chatwoot/application/custom_filters_test.go internal/chatwoot/application/account.go internal/chatwoot/application/account_test.go internal/chatwoot/application/agent_bots.go internal/chatwoot/application/agent_bots_test.go internal/chatwoot/application/audit_logs.go internal/chatwoot/application/audit_logs_test.go
git commit -m "feat(application): add labels, custom attributes, custom filters, account, agent bots, audit logs API clients"
```

---

## Task 9: Integrations + Help Center API Clients

**Files:**
- Create: `internal/chatwoot/application/integrations.go`
- Create: `internal/chatwoot/application/integrations_test.go`
- Create: `internal/chatwoot/application/help_center.go`
- Create: `internal/chatwoot/application/help_center_test.go`

### Integrations (`integrations.go`)

- [ ] **Step 1: Implement integrations.go with 4 methods**

**HTTP methods and paths:**
- `ListIntegrationApps(ctx)` — GET `/api/v1/accounts/{id}/integrations/apps` — decode with payload wrapper: `var body struct { Payload []Integration \x60json:"payload"\x60 }`
- `CreateIntegrationHook(ctx, CreateIntegrationHookOpts)` — POST `/api/v1/accounts/{id}/integrations/hooks` — marshal opts, decode into `IntegrationHook`
- `UpdateIntegrationHook(ctx, hookID, UpdateIntegrationHookOpts)` — PATCH `/api/v1/accounts/{id}/integrations/hooks/{hid}` — marshal opts, decode into `IntegrationHook`
- `DeleteIntegrationHook(ctx, hookID)` — DELETE `/api/v1/accounts/{id}/integrations/hooks/{hid}` — decode with `nil`

- [ ] **Step 2: Write integrations_test.go with TestListIntegrationApps, TestCreateIntegrationHook**

### Help Center (`help_center.go`)

- [ ] **Step 3: Implement help_center.go with 5 methods**

**HTTP methods and paths:**
- `ListPortals(ctx)` — GET `/api/v1/accounts/{id}/portals` — decode with payload wrapper: `var body struct { Payload []Portal \x60json:"payload"\x60 }`
- `CreatePortal(ctx, CreatePortalOpts)` — POST `/api/v1/accounts/{id}/portals` — marshal opts, decode into `Portal`
- `UpdatePortal(ctx, portalID, UpdatePortalOpts)` — PATCH `/api/v1/accounts/{id}/portals/{pid}` — marshal opts, decode into `Portal`
- `CreateArticle(ctx, portalID, CreateArticleOpts)` — POST `/api/v1/accounts/{id}/portals/{pid}/articles` — marshal opts, decode into `Article`
- `CreateCategory(ctx, portalID, CreateCategoryOpts)` — POST `/api/v1/accounts/{id}/portals/{pid}/categories` — marshal opts, decode into `Category`

**Note:** `CreateArticle` and `CreateCategory` take a `portalID` parameter that goes into the URL path.

- [ ] **Step 4: Write help_center_test.go with TestListPortals, TestCreateArticle**

### Run and Commit

- [ ] **Step 5: Run all tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/chatwoot/application/ -v -run "TestIntegration|TestPortal|TestArticle"`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/chatwoot/application/integrations.go internal/chatwoot/application/integrations_test.go internal/chatwoot/application/help_center.go internal/chatwoot/application/help_center_test.go
git commit -m "feat(application): add integrations and help center API clients"
```

---

## Task 10: Teams CLI Commands

**Files:**
- Create: `internal/cli/application/teams/teams.go`
- Create: `internal/cli/application/teams/list.go`
- Create: `internal/cli/application/teams/get.go`
- Create: `internal/cli/application/teams/create.go`
- Create: `internal/cli/application/teams/update.go`
- Create: `internal/cli/application/teams/delete.go`
- Create: `internal/cli/application/teams/members.go`
- Create: `internal/cli/application/teams/teams_test.go`
- Create: `internal/cli/application/teams/testroot_test.go`

- [ ] **Step 1: Create testroot_test.go**

```go
package teams

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

- [ ] **Step 2: Create teams.go group command**

```go
package teams

import "github.com/spf13/cobra"

// Cmd is the teams command group.
var Cmd = &cobra.Command{
	Use:   "teams",
	Short: "Manage teams",
}

var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage team members",
}

func init() {
	Cmd.AddCommand(membersCmd)
}
```

- [ ] **Step 3: Write failing test for teams list**

In `teams_test.go`, write `TestTeamsList` that sets up an httptest server returning `[{"id": 1, "name": "Support", "account_id": 1}]`, uses env vars, executes `application teams list`, and verifies JSON envelope.

- [ ] **Step 4: Implement list.go**

The handler follows the standard cmdutil pipeline. `ListTeams` does not take pagination — it returns all teams. No pagination flags needed.

```go
// Handler pattern:
// cmdutil.ResolveContext → cmdutil.ResolveAuth → chatwoot.NewClient → appapi.NewClient
// → client.ListTeams(ctx) → contract.SuccessList(teams, contract.Meta{}) → contract.Write
```

- [ ] **Step 5: Implement get.go, create.go, update.go, delete.go**

Flags:
- get: `--id` (required)
- create: `--name` (required), `--description` (optional), `--allow-auto-assign` (optional bool)
- update: `--id` (required), `--name`, `--description`, `--allow-auto-assign` (at least one)
- delete: `--id` (required)

For `delete`, return `contract.Success(map[string]any{"deleted": true})`.

- [ ] **Step 6: Implement members.go**

Follow the exact pattern from `internal/cli/application/inboxes/members.go`. Four subcommands on `membersCmd`:
- `members list`: `--team-id` (required) — calls `client.ListTeamMembers`
- `members add`: `--team-id` (required), `--agent-ids` (required, comma-separated) — calls `client.AddTeamMember`
- `members update`: `--team-id` (required), `--agent-ids` (required) — calls `client.UpdateTeamMembers`
- `members delete`: `--team-id` (required), `--agent-ids` (required) — calls `client.RemoveTeamMember`

Include `parseIntList` helper (same as in inboxes/members.go) for parsing comma-separated agent IDs.

- [ ] **Step 7: Add tests for get, create, delete, members list, members add**

- [ ] **Step 8: Run all teams CLI tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/teams/ -v`
Expected: All PASS

- [ ] **Step 9: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/cli/application/teams/
git commit -m "feat(cli): add teams command group with 9 commands"
```

---

## Task 11: Agents CLI Commands

**Files:**
- Create: `internal/cli/application/agents/agents.go`
- Create: `internal/cli/application/agents/list.go`
- Create: `internal/cli/application/agents/create.go`
- Create: `internal/cli/application/agents/update.go`
- Create: `internal/cli/application/agents/delete.go`
- Create: `internal/cli/application/agents/agents_test.go`
- Create: `internal/cli/application/agents/testroot_test.go`

- [ ] **Step 1: Create testroot_test.go**

Same boilerplate as teams, just change package to `agents`.

- [ ] **Step 2: Create agents.go group command**

```go
package agents

import "github.com/spf13/cobra"

// Cmd is the agents command group.
var Cmd = &cobra.Command{
	Use:   "agents",
	Short: "Manage agents",
}
```

- [ ] **Step 3: Write failing test for agents list**

`TestAgentsList` — httptest server returns `[{"id": 1, "name": "Alice", "email": "alice@test.com"}]`.

- [ ] **Step 4: Implement list.go, create.go, update.go, delete.go**

Flags:
- list: no required flags (returns all agents)
- create: `--name` (required), `--email` (required), `--role` (optional, default "agent")
- update: `--id` (required), `--name`, `--role` (at least one)
- delete: `--id` (required)

- [ ] **Step 5: Add tests for create and delete**

- [ ] **Step 6: Run all agents CLI tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/agents/ -v`
Expected: All PASS

- [ ] **Step 7: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/cli/application/agents/
git commit -m "feat(cli): add agents command group with 4 commands"
```

---

## Task 12: Canned Responses CLI Commands

**Files:**
- Create: `internal/cli/application/cannedresponses/cannedresponses.go`
- Create: `internal/cli/application/cannedresponses/list.go`
- Create: `internal/cli/application/cannedresponses/create.go`
- Create: `internal/cli/application/cannedresponses/update.go`
- Create: `internal/cli/application/cannedresponses/delete.go`
- Create: `internal/cli/application/cannedresponses/cannedresponses_test.go`
- Create: `internal/cli/application/cannedresponses/testroot_test.go`

- [ ] **Step 1: Create testroot_test.go**

Same boilerplate, package `cannedresponses`.

- [ ] **Step 2: Create cannedresponses.go group command**

```go
package cannedresponses

import "github.com/spf13/cobra"

// Cmd is the canned-responses command group.
var Cmd = &cobra.Command{
	Use:   "canned-responses",
	Short: "Manage canned responses",
}
```

**Note:** Go package is `cannedresponses`, CLI use string is `canned-responses`.

- [ ] **Step 3: Write failing test for canned-responses list**

`TestCannedResponsesList` — httptest server returns `[{"id": 1, "short_code": "greeting", "content": "Hello!"}]`.

- [ ] **Step 4: Implement list.go, create.go, update.go, delete.go**

Flags:
- list: no required flags
- create: `--short-code` (required), `--content` (required)
- update: `--id` (required), `--short-code`, `--content` (at least one)
- delete: `--id` (required)

- [ ] **Step 5: Add tests for create and delete**

- [ ] **Step 6: Run all canned responses CLI tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/cannedresponses/ -v`
Expected: All PASS

- [ ] **Step 7: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/cli/application/cannedresponses/
git commit -m "feat(cli): add canned-responses command group with 4 commands"
```

---

## Task 13: Reports CLI Commands

**Files:**
- Create: `internal/cli/application/reports/reports.go`
- Create: `internal/cli/application/reports/account.go`
- Create: `internal/cli/application/reports/summary.go`
- Create: `internal/cli/application/reports/conversations.go`
- Create: `internal/cli/application/reports/first_response.go`
- Create: `internal/cli/application/reports/inbox_label.go`
- Create: `internal/cli/application/reports/outgoing.go`
- Create: `internal/cli/application/reports/summary_agent.go`
- Create: `internal/cli/application/reports/summary_channel.go`
- Create: `internal/cli/application/reports/summary_inbox.go`
- Create: `internal/cli/application/reports/summary_team.go`
- Create: `internal/cli/application/reports/events.go`
- Create: `internal/cli/application/reports/reports_test.go`
- Create: `internal/cli/application/reports/testroot_test.go`

- [ ] **Step 1: Create testroot_test.go**

Same boilerplate, package `reports`.

- [ ] **Step 2: Create reports.go group command**

```go
package reports

import "github.com/spf13/cobra"

// Cmd is the reports command group.
var Cmd = &cobra.Command{
	Use:   "reports",
	Short: "View reports and metrics",
}
```

- [ ] **Step 3: Write failing test for reports summary**

`TestReportsSummary` — httptest server verifies query params include `type`, `since`, `until` and returns `{"avg_first_response_time": "5m", "conversations_count": 100, ...}`.

- [ ] **Step 4: Implement all 11 report subcommands**

Each subcommand has common flags: `--since`, `--until`, `--type`, `--metric`, `--id`. Not all are required for every subcommand.

**Subcommand → API method mapping:**

| Subcommand (Use) | Flags | API Method |
|---|---|---|
| `account` | `--metric` (required), `--type` (required), `--since`, `--until`, `--id` | `GetReports` |
| `summary` | `--type` (required), `--since`, `--until`, `--id` | `GetReportSummary` |
| `conversations` | `--type` (required: "account" or "agent"), `--since`, `--until` | If `--type=account`: `GetConversationMetrics`; if `--type=agent`: `GetAgentConversationMetrics` |
| `first-response-distribution` | `--since`, `--until` | `GetFirstResponseTimeDistribution` |
| `inbox-label-matrix` | `--since`, `--until` | `GetInboxLabelMatrix` |
| `outgoing-messages` | `--since`, `--until` | `GetOutgoingMessagesCount` |
| `summary-by-agent` | `--since`, `--until` | `GetSummaryByAgent` |
| `summary-by-channel` | `--since`, `--until` | `GetSummaryByChannel` |
| `summary-by-inbox` | `--since`, `--until` | `GetSummaryByInbox` |
| `summary-by-team` | `--since`, `--until` | `GetSummaryByTeam` |
| `events` | `--since`, `--until` | `GetReportingEvents` |

Each handler builds `appapi.ReportOpts` from flags, calls the appropriate method, wraps the result with `contract.Success(result)` (single object, not list), and writes.

**Pattern for a report handler:**
```go
func runSummary(cmd *cobra.Command, args []string) error {
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

	typ, _ := cmd.Flags().GetString("type")
	since, _ := cmd.Flags().GetString("since")
	until, _ := cmd.Flags().GetString("until")
	id, _ := cmd.Flags().GetString("id")

	opts := appapi.ReportOpts{Type: typ, Since: since, Until: until, ID: id}
	result, err := client.GetReportSummary(context.Background(), opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}
	resp := contract.Success(result)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
```

- [ ] **Step 5: Add tests for account, conversations (both type variants), and events**

Test `TestReportsSummary` verifies `--since` and `--until` flow through to query parameters.
Test `TestReportsConversationsAgent` verifies `--type agent` dispatches to the agent variant.

- [ ] **Step 6: Run all reports CLI tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/reports/ -v`
Expected: All PASS

- [ ] **Step 7: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/cli/application/reports/
git commit -m "feat(cli): add reports command group with 11 commands"
```

---

## Task 14: Webhooks CLI Commands

**Files:**
- Create: `internal/cli/application/webhooks/webhooks.go`
- Create: `internal/cli/application/webhooks/list.go`
- Create: `internal/cli/application/webhooks/create.go`
- Create: `internal/cli/application/webhooks/update.go`
- Create: `internal/cli/application/webhooks/delete.go`
- Create: `internal/cli/application/webhooks/webhooks_test.go`
- Create: `internal/cli/application/webhooks/testroot_test.go`

- [ ] **Step 1: Create testroot_test.go**

Same boilerplate, package `webhooks`.

- [ ] **Step 2: Create webhooks.go group command**

```go
package webhooks

import "github.com/spf13/cobra"

// Cmd is the webhooks command group.
var Cmd = &cobra.Command{
	Use:   "webhooks",
	Short: "Manage webhooks",
}
```

- [ ] **Step 3: Write failing test for webhooks list**

`TestWebhooksList` — httptest server returns `{"payload": [{"id": 1, "url": "https://example.com", "subscriptions": ["message_created"]}]}`.

- [ ] **Step 4: Implement list.go, create.go, update.go, delete.go**

Flags:
- list: no required flags
- create: `--url` (required), `--subscriptions` (optional, comma-separated event types)
- update: `--id` (required), `--url`, `--subscriptions` (at least one)
- delete: `--id` (required)

For create/update, split `--subscriptions` string by comma into `[]string`.

- [ ] **Step 5: Add tests for create and delete**

Test `TestWebhooksCreate` verifies request body includes `url` and `subscriptions`.

- [ ] **Step 6: Run all webhooks CLI tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/webhooks/ -v`
Expected: All PASS

- [ ] **Step 7: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/cli/application/webhooks/
git commit -m "feat(cli): add webhooks command group with 4 commands"
```

---

## Task 15: Automation Rules CLI Commands

**Files:**
- Create: `internal/cli/application/automationrules/automationrules.go`
- Create: `internal/cli/application/automationrules/list.go`
- Create: `internal/cli/application/automationrules/get.go`
- Create: `internal/cli/application/automationrules/create.go`
- Create: `internal/cli/application/automationrules/update.go`
- Create: `internal/cli/application/automationrules/delete.go`
- Create: `internal/cli/application/automationrules/automationrules_test.go`
- Create: `internal/cli/application/automationrules/testroot_test.go`

- [ ] **Step 1: Create testroot_test.go**

Same boilerplate, package `automationrules`.

- [ ] **Step 2: Create automationrules.go group command**

```go
package automationrules

import "github.com/spf13/cobra"

// Cmd is the automation-rules command group.
var Cmd = &cobra.Command{
	Use:   "automation-rules",
	Short: "Manage automation rules",
}
```

**Note:** Go package is `automationrules`, CLI use string is `automation-rules`.

- [ ] **Step 3: Write failing test for automation-rules list**

`TestAutomationRulesList` — httptest server returns `{"payload": [{"id": 1, "name": "Auto-assign", "event_name": "conversation_created"}]}`.

- [ ] **Step 4: Implement list.go, get.go, create.go, update.go, delete.go**

Flags:
- list: no required flags
- get: `--id` (required)
- create: `--name` (required), `--event-name` (required), `--conditions` (required, JSON string), `--actions` (required, JSON string), `--description` (optional)
- update: `--id` (required), `--name`, `--event-name`, `--conditions`, `--actions`, `--description` (at least one)
- delete: `--id` (required)

For create/update, parse `--conditions` and `--actions` with `json.Unmarshal` into `any`:
```go
conditionsStr, _ := cmd.Flags().GetString("conditions")
var conditions any
if err := json.Unmarshal([]byte(conditionsStr), &conditions); err != nil {
	return cmdutil.WriteError(cmd, contract.ErrCodeValidation, "invalid conditions JSON: "+err.Error())
}
```

- [ ] **Step 5: Add tests for create and delete**

Test `TestAutomationRulesCreate` verifies conditions and actions JSON flow through to request body.

- [ ] **Step 6: Run all automation rules CLI tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/automationrules/ -v`
Expected: All PASS

- [ ] **Step 7: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/cli/application/automationrules/
git commit -m "feat(cli): add automation-rules command group with 5 commands"
```

---

## Task 16: Simple CRUD CLI Commands Batch

**Files (6 packages):**

### Labels
- Create: `internal/cli/application/labels/labels.go`
- Create: `internal/cli/application/labels/list.go`
- Create: `internal/cli/application/labels/get.go`
- Create: `internal/cli/application/labels/create.go`
- Create: `internal/cli/application/labels/update.go`
- Create: `internal/cli/application/labels/delete.go`
- Create: `internal/cli/application/labels/labels_test.go`
- Create: `internal/cli/application/labels/testroot_test.go`

### Custom Attributes
- Create: `internal/cli/application/customattributes/customattributes.go`
- Create: `internal/cli/application/customattributes/list.go`
- Create: `internal/cli/application/customattributes/get.go`
- Create: `internal/cli/application/customattributes/create.go`
- Create: `internal/cli/application/customattributes/update.go`
- Create: `internal/cli/application/customattributes/delete.go`
- Create: `internal/cli/application/customattributes/customattributes_test.go`
- Create: `internal/cli/application/customattributes/testroot_test.go`

### Custom Filters
- Create: `internal/cli/application/customfilters/customfilters.go`
- Create: `internal/cli/application/customfilters/list.go`
- Create: `internal/cli/application/customfilters/get.go`
- Create: `internal/cli/application/customfilters/create.go`
- Create: `internal/cli/application/customfilters/update.go`
- Create: `internal/cli/application/customfilters/delete.go`
- Create: `internal/cli/application/customfilters/customfilters_test.go`
- Create: `internal/cli/application/customfilters/testroot_test.go`

### Account
- Create: `internal/cli/application/account/account.go`
- Create: `internal/cli/application/account/get.go`
- Create: `internal/cli/application/account/update.go`
- Create: `internal/cli/application/account/account_test.go`
- Create: `internal/cli/application/account/testroot_test.go`

### Agent Bots
- Create: `internal/cli/application/agentbots/agentbots.go`
- Create: `internal/cli/application/agentbots/list.go`
- Create: `internal/cli/application/agentbots/get.go`
- Create: `internal/cli/application/agentbots/create.go`
- Create: `internal/cli/application/agentbots/update.go`
- Create: `internal/cli/application/agentbots/delete.go`
- Create: `internal/cli/application/agentbots/agentbots_test.go`
- Create: `internal/cli/application/agentbots/testroot_test.go`

### Audit Logs
- Create: `internal/cli/application/auditlogs/auditlogs.go`
- Create: `internal/cli/application/auditlogs/list.go`
- Create: `internal/cli/application/auditlogs/auditlogs_test.go`
- Create: `internal/cli/application/auditlogs/testroot_test.go`

### Labels Package

- [ ] **Step 1: Create labels testroot_test.go and labels.go**

```go
package labels

import "github.com/spf13/cobra"

// Cmd is the labels command group.
var Cmd = &cobra.Command{
	Use:   "labels",
	Short: "Manage account labels",
}
```

- [ ] **Step 2: Implement list.go, get.go, create.go, update.go, delete.go**

Flags:
- list: no required flags
- get: `--id` (required)
- create: `--title` (required), `--description` (optional), `--color` (optional), `--show-on-sidebar` (optional bool)
- update: `--id` (required), `--title`, `--description`, `--color`, `--show-on-sidebar` (at least one)
- delete: `--id` (required)

- [ ] **Step 3: Write labels_test.go with TestLabelsList and TestLabelsCreate**

### Custom Attributes Package

- [ ] **Step 4: Create customattributes testroot_test.go and customattributes.go**

```go
package customattributes

import "github.com/spf13/cobra"

// Cmd is the custom-attributes command group.
var Cmd = &cobra.Command{
	Use:   "custom-attributes",
	Short: "Manage custom attribute definitions",
}
```

- [ ] **Step 5: Implement list.go, get.go, create.go, update.go, delete.go**

Flags:
- list: no required flags
- get: `--id` (required)
- create: `--attribute-key` (required), `--attribute-model` (required: "contact" or "conversation"), `--attribute-type` (required, maps to `attribute_display_type`), `--name` (required, maps to `attribute_display_name`), `--description` (optional)
- update: `--id` (required), `--name`, `--description` (at least one)
- delete: `--id` (required)

- [ ] **Step 6: Write customattributes_test.go with TestCustomAttributesList and TestCustomAttributesCreate**

### Custom Filters Package

- [ ] **Step 7: Create customfilters testroot_test.go and customfilters.go**

```go
package customfilters

import "github.com/spf13/cobra"

// Cmd is the custom-filters command group.
var Cmd = &cobra.Command{
	Use:   "custom-filters",
	Short: "Manage custom filters",
}
```

- [ ] **Step 8: Implement list.go, get.go, create.go, update.go, delete.go**

Flags:
- list: `--filter-type` (optional: "conversation", "contact", "report")
- get: `--id` (required)
- create: `--name` (required), `--type` (required: "conversation", "contact", "report"), `--query` (required, JSON string)
- update: `--id` (required), `--name`, `--query` (at least one)
- delete: `--id` (required)

For `list`, pass `--filter-type` value to `client.ListCustomFilters(ctx, filterType)`.
For `create`/`update`, parse `--query` with `json.Unmarshal` into `any`.

- [ ] **Step 9: Write customfilters_test.go with TestCustomFiltersList and TestCustomFiltersCreate**

### Account Package

- [ ] **Step 10: Create account testroot_test.go and account.go**

```go
package account

import "github.com/spf13/cobra"

// Cmd is the account command group.
var Cmd = &cobra.Command{
	Use:   "account",
	Short: "Manage account settings",
}
```

- [ ] **Step 11: Implement get.go and update.go**

Flags:
- get: no required flags (uses account ID from runtime context)
- update: `--name`, `--locale`, `--domain` (at least one)

**Note:** No `--id` flag. The account ID comes from `rctx.AccountID` via the context resolver, same as every other command.

- [ ] **Step 12: Write account_test.go with TestAccountGet and TestAccountUpdate**

### Agent Bots Package

- [ ] **Step 13: Create agentbots testroot_test.go and agentbots.go**

```go
package agentbots

import "github.com/spf13/cobra"

// Cmd is the agent-bots command group.
var Cmd = &cobra.Command{
	Use:   "agent-bots",
	Short: "Manage account agent bots",
}
```

- [ ] **Step 14: Implement list.go, get.go, create.go, update.go, delete.go**

Flags:
- list: no required flags
- get: `--id` (required)
- create: `--name` (required), `--description` (optional), `--bot-type` (optional), `--outgoing-url` (optional), `--bot-config` (optional, JSON string)
- update: `--id` (required), `--name`, `--description`, `--bot-type`, `--outgoing-url`, `--bot-config` (at least one)
- delete: `--id` (required)

For `--bot-config`, parse with `json.Unmarshal` into `any`.

- [ ] **Step 15: Write agentbots_test.go with TestAgentBotsList and TestAgentBotsCreate**

### Audit Logs Package

- [ ] **Step 16: Create auditlogs testroot_test.go and auditlogs.go**

```go
package auditlogs

import "github.com/spf13/cobra"

// Cmd is the audit-logs command group.
var Cmd = &cobra.Command{
	Use:   "audit-logs",
	Short: "View audit logs",
}
```

- [ ] **Step 17: Implement list.go with pagination**

Flags: `--page` (optional, default 1)

The handler calls `client.ListAuditLogs(ctx, page)` and wraps with `contract.SuccessList`.

- [ ] **Step 18: Write auditlogs_test.go with TestAuditLogsList**

### Run and Commit

- [ ] **Step 19: Run all tests for the batch**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/labels/ ./internal/cli/application/customattributes/ ./internal/cli/application/customfilters/ ./internal/cli/application/account/ ./internal/cli/application/agentbots/ ./internal/cli/application/auditlogs/ -v`
Expected: All PASS

- [ ] **Step 20: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/cli/application/labels/ internal/cli/application/customattributes/ internal/cli/application/customfilters/ internal/cli/application/account/ internal/cli/application/agentbots/ internal/cli/application/auditlogs/
git commit -m "feat(cli): add labels, custom-attributes, custom-filters, account, agent-bots, audit-logs command groups"
```

---

## Task 17: Integrations + Help Center CLI Commands

**Files:**

### Integrations
- Create: `internal/cli/application/integrations/integrations.go`
- Create: `internal/cli/application/integrations/apps.go`
- Create: `internal/cli/application/integrations/hooks.go`
- Create: `internal/cli/application/integrations/integrations_test.go`
- Create: `internal/cli/application/integrations/testroot_test.go`

### Help Center
- Create: `internal/cli/application/helpcenter/helpcenter.go`
- Create: `internal/cli/application/helpcenter/portals.go`
- Create: `internal/cli/application/helpcenter/articles.go`
- Create: `internal/cli/application/helpcenter/categories.go`
- Create: `internal/cli/application/helpcenter/helpcenter_test.go`
- Create: `internal/cli/application/helpcenter/testroot_test.go`

### Integrations Package

- [ ] **Step 1: Create integrations testroot_test.go and integrations.go**

```go
package integrations

import "github.com/spf13/cobra"

// Cmd is the integrations command group.
var Cmd = &cobra.Command{
	Use:   "integrations",
	Short: "Manage integrations",
}

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Manage integration apps",
}

var hooksCmd = &cobra.Command{
	Use:   "hooks",
	Short: "Manage integration hooks",
}

func init() {
	Cmd.AddCommand(appsCmd)
	Cmd.AddCommand(hooksCmd)
}
```

- [ ] **Step 2: Implement apps.go**

Single subcommand on `appsCmd`:
- `apps list`: no required flags — calls `client.ListIntegrationApps`

- [ ] **Step 3: Implement hooks.go**

Three subcommands on `hooksCmd`:
- `hooks create`: `--app-id` (required), `--inbox-id` (optional), `--settings` (optional, JSON string)
- `hooks update`: `--hook-id` (required), `--settings` (optional, JSON string)
- `hooks delete`: `--hook-id` (required)

For `--settings`, parse with `json.Unmarshal` into `any`.

- [ ] **Step 4: Write integrations_test.go with TestIntegrationsAppsList and TestIntegrationsHooksCreate**

### Help Center Package

- [ ] **Step 5: Create helpcenter testroot_test.go and helpcenter.go**

```go
package helpcenter

import "github.com/spf13/cobra"

// Cmd is the help-center command group.
var Cmd = &cobra.Command{
	Use:   "help-center",
	Short: "Manage help center",
}

var portalsCmd = &cobra.Command{
	Use:   "portals",
	Short: "Manage portals",
}

var articlesCmd = &cobra.Command{
	Use:   "articles",
	Short: "Manage articles",
}

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "Manage categories",
}

func init() {
	Cmd.AddCommand(portalsCmd)
	Cmd.AddCommand(articlesCmd)
	Cmd.AddCommand(categoriesCmd)
}
```

**Note:** Go package is `helpcenter`, CLI use string is `help-center`.

- [ ] **Step 6: Implement portals.go**

Three subcommands on `portalsCmd`:
- `portals list`: no required flags — calls `client.ListPortals`
- `portals create`: `--name` (required), `--slug` (optional), `--color` (optional), `--header-text` (optional), `--custom-domain` (optional)
- `portals update`: `--id` (required), `--name`, `--slug`, `--color`, `--header-text`, `--custom-domain`, `--archived` (at least one)

- [ ] **Step 7: Implement articles.go**

Single subcommand on `articlesCmd`:
- `articles create`: `--portal-id` (required), `--title` (required), `--content` (optional), `--description` (optional), `--status` (optional int), `--category-id` (optional), `--author-id` (optional)

- [ ] **Step 8: Implement categories.go**

Single subcommand on `categoriesCmd`:
- `categories create`: `--portal-id` (required), `--name` (required), `--description` (optional), `--locale` (optional), `--position` (optional int)

### Run and Commit

- [ ] **Step 9: Write helpcenter_test.go with TestPortalsList and TestArticlesCreate**

- [ ] **Step 10: Run all tests**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./internal/cli/application/integrations/ ./internal/cli/application/helpcenter/ -v`
Expected: All PASS

- [ ] **Step 11: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/cli/application/integrations/ internal/cli/application/helpcenter/
git commit -m "feat(cli): add integrations and help-center command groups"
```

---

## Task 18: Register All Subgroups

**Files:**
- Modify: `internal/cli/application/application.go`

- [ ] **Step 1: Add all 14 new imports and AddCommand calls**

Update `internal/cli/application/application.go` to register all Sprint D command packages:

```go
package application

import (
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/account"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/agentbots"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/agents"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/auditlogs"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/automationrules"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/cannedresponses"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/contacts"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/conversations"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/customattributes"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/customfilters"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/helpcenter"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/inboxes"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/integrations"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/labels"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/messages"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/reports"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/teams"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/webhooks"
	"github.com/spf13/cobra"
)

// Cmd is the application command group.
var Cmd = &cobra.Command{
	Use:   "application",
	Short: "Application API commands (agent/admin)",
}

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage authenticated user profile",
}

func init() {
	Cmd.AddCommand(profileCmd)
	// Sprint C
	Cmd.AddCommand(contacts.Cmd)
	Cmd.AddCommand(conversations.Cmd)
	Cmd.AddCommand(messages.Cmd)
	Cmd.AddCommand(inboxes.Cmd)
	// Sprint D
	Cmd.AddCommand(teams.Cmd)
	Cmd.AddCommand(agents.Cmd)
	Cmd.AddCommand(cannedresponses.Cmd)
	Cmd.AddCommand(reports.Cmd)
	Cmd.AddCommand(webhooks.Cmd)
	Cmd.AddCommand(automationrules.Cmd)
	Cmd.AddCommand(labels.Cmd)
	Cmd.AddCommand(customattributes.Cmd)
	Cmd.AddCommand(customfilters.Cmd)
	Cmd.AddCommand(account.Cmd)
	Cmd.AddCommand(agentbots.Cmd)
	Cmd.AddCommand(auditlogs.Cmd)
	Cmd.AddCommand(integrations.Cmd)
	Cmd.AddCommand(helpcenter.Cmd)
}
```

- [ ] **Step 2: Build and verify**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go build ./cmd/chatwoot/`
Expected: Success

- [ ] **Step 3: Run full test suite**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./... 2>&1 | grep -E "^(ok|FAIL)"`
Expected: All pass

- [ ] **Step 4: Commit**

```bash
cd /Users/brettmcdowell/Dev/Chatwoot_CLI && git add internal/cli/application/application.go
git commit -m "feat(cli): register all 14 Sprint D command groups in application"
```

---

## Task 19: Full Test Suite and Exit Criteria Verification

**Files:**
- No new files

- [ ] **Step 1: Run full test suite**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go test ./... -v 2>&1 | tail -80`
Expected: All tests PASS

- [ ] **Step 2: Run vet and build**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && go vet ./... && go build ./cmd/chatwoot/`
Expected: No errors

- [ ] **Step 3: Verify all 14 resource families visible**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application --help`
Expected: Shows teams, agents, canned-responses, reports, webhooks, automation-rules, labels, custom-attributes, custom-filters, account, agent-bots, audit-logs, integrations, help-center (plus existing contacts, conversations, messages, inboxes, profile)

- [ ] **Step 4: Verify teams command tree**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application teams --help`
Expected: Shows list, get, create, update, delete, members

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application teams members --help`
Expected: Shows list, add, update, delete

- [ ] **Step 5: Verify reports command tree**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application reports --help`
Expected: Shows account, summary, conversations, first-response-distribution, inbox-label-matrix, outgoing-messages, summary-by-agent, summary-by-channel, summary-by-inbox, summary-by-team, events

- [ ] **Step 6: Verify integrations command tree**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application integrations --help`
Expected: Shows apps, hooks

- [ ] **Step 7: Verify help-center command tree**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application help-center --help`
Expected: Shows portals, articles, categories

- [ ] **Step 8: Verify existing Sprint C commands still work**

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot application contacts --help`
Expected: Shows list, get, create, update, delete, search, filter, merge, labels, conversations

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot version`
Expected: JSON envelope with version info

Run: `cd /Users/brettmcdowell/Dev/Chatwoot_CLI && ./chatwoot auth --help`
Expected: Shows set, status, clear

---

## Sprint D Exit Criteria Checklist

- [ ] `go test ./...` passes with all tests green
- [ ] `go vet ./...` clean
- [ ] `go build ./cmd/chatwoot/` succeeds
- [ ] All 14 resource families visible under `chatwoot application --help`
- [ ] `chatwoot application teams list` returns JSON envelope
- [ ] `chatwoot application teams members list --team-id 1` returns JSON envelope
- [ ] `chatwoot application reports summary --type account --since X --until Y` works
- [ ] `chatwoot application webhooks create --url X --subscriptions Y` works
- [ ] `chatwoot application automation-rules list` returns JSON envelope
- [ ] `chatwoot application canned-responses list` returns JSON envelope
- [ ] `chatwoot application labels list` returns JSON envelope
- [ ] `chatwoot application integrations apps list` returns JSON envelope
- [ ] `chatwoot application help-center portals list` returns JSON envelope
- [ ] Team member operations follow inbox member pattern
- [ ] Existing Sprint C commands still work
- [ ] All commands produce valid JSON envelopes on stdout
- [ ] All commands use the cmdutil pipeline
- [ ] No business logic in command handler files

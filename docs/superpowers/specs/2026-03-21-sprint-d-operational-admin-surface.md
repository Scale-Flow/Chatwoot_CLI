# Design Spec: Sprint D — Operational Admin Surface (Layer 8)

## Scope

This spec defines the operational admin command surface for the Chatwoot CLI:
14 resource families covering account administration, team management,
reporting, automation, and knowledge base operations. It covers API client
methods, CLI command packages, model types, and testing strategy.

## Goals

- Implement the complete Layer 8 admin surface from the roadmap
- Enable operational workflows: backlog summaries (UC8), webhook/automation
  management (UC9), team-based escalation (UC4), canned response lookups (UC3)
- Follow the same patterns established in Sprint C for consistency
- Each resource family is independently implementable by a subagent

## Prerequisites

Sprints A–C are complete. The following exist:

- Config, credentials, transport, retry, error mapping
- Application API client with contacts, conversations, messages, inboxes
- CLI shell with auth, profile, and core commands
- JSON envelope contracts and rendering
- `cmdutil` shared helpers (context resolution, pagination flags)
- `chatwoot.ListAll` pagination helper

## Resource Ordering

Use-case-critical resources first, then admin housekeeping:

| # | Resource | API Methods | CLI Commands | Use Case |
|---|----------|-------------|-------------|----------|
| 1 | Teams (+members) | 9 | 9 | UC4: Escalation targets |
| 2 | Agents | 4 | 4 | UC4: Assignment lookups |
| 3 | Canned Responses | 4 | 4 | UC3: Reply templates |
| 4 | Reports | 12 | 11 | UC8: Backlog summary |
| 5 | Webhooks | 4 | 4 | UC9: Event subscriptions |
| 6 | Automation Rules | 5 | 5 | UC9: Workflow automation |
| 7 | Labels | 5 | 5 | Label definition CRUD |
| 8 | Custom Attributes | 5 | 5 | Attribute definitions |
| 9 | Custom Filters | 5 | 5 | Saved searches |
| 10 | Account | 2 | 2 | Account settings |
| 11 | Agent Bots | 5 | 5 | Bot management |
| 12 | Audit Logs | 1 | 1 | Compliance/history |
| 13 | Integrations | 4 | 4 | Apps + hooks |
| 14 | Help Center | 5 | 5 | Knowledge base |

**Totals:** ~70 API methods, ~69 CLI commands

## API Client Methods

All methods live in `internal/chatwoot/application/`, one file per resource.
Each follows the Sprint C pattern: build path → `transport.DoWithRetry` →
`chatwoot.DecodeResponse` → return typed data.

### `teams.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListTeams(ctx)` | GET | `/api/v1/accounts/{id}/teams` | `([]Team, error)` |
| `GetTeam(ctx, teamID)` | GET | `/api/v1/accounts/{id}/teams/{tid}` | `(*Team, error)` |
| `CreateTeam(ctx, CreateTeamOpts)` | POST | `/api/v1/accounts/{id}/teams` | `(*Team, error)` |
| `UpdateTeam(ctx, teamID, UpdateTeamOpts)` | PATCH | `/api/v1/accounts/{id}/teams/{tid}` | `(*Team, error)` |
| `DeleteTeam(ctx, teamID)` | DELETE | `/api/v1/accounts/{id}/teams/{tid}` | `error` |
| `ListTeamMembers(ctx, teamID)` | GET | `/api/v1/accounts/{id}/teams/{tid}/team_members` | `([]Agent, error)` |
| `AddTeamMember(ctx, teamID, agentIDs)` | POST | `/api/v1/accounts/{id}/teams/{tid}/team_members` | `([]Agent, error)` |
| `UpdateTeamMembers(ctx, teamID, agentIDs)` | PATCH | `/api/v1/accounts/{id}/teams/{tid}/team_members` | `([]Agent, error)` |
| `RemoveTeamMember(ctx, teamID, agentIDs)` | DELETE | `/api/v1/accounts/{id}/teams/{tid}/team_members` | `error` |

Team members body format: `{"user_ids": [1, 2, 3]}`.

### `agents.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListAgents(ctx)` | GET | `/api/v1/accounts/{id}/agents` | `([]Agent, error)` |
| `CreateAgent(ctx, CreateAgentOpts)` | POST | `/api/v1/accounts/{id}/agents` | `(*Agent, error)` |
| `UpdateAgent(ctx, agentID, UpdateAgentOpts)` | PATCH | `/api/v1/accounts/{id}/agents/{aid}` | `(*Agent, error)` |
| `DeleteAgent(ctx, agentID)` | DELETE | `/api/v1/accounts/{id}/agents/{aid}` | `error` |

### `canned_responses.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListCannedResponses(ctx)` | GET | `/api/v1/accounts/{id}/canned_responses` | `([]CannedResponse, error)` |
| `CreateCannedResponse(ctx, CreateCannedResponseOpts)` | POST | `/api/v1/accounts/{id}/canned_responses` | `(*CannedResponse, error)` |
| `UpdateCannedResponse(ctx, id, UpdateCannedResponseOpts)` | PATCH | `/api/v1/accounts/{id}/canned_responses/{cid}` | `(*CannedResponse, error)` |
| `DeleteCannedResponse(ctx, id)` | DELETE | `/api/v1/accounts/{id}/canned_responses/{cid}` | `error` |

### `reports.go` (new file)

Reports are read-only, span `/api/v1` and `/api/v2`, and use date-range
parameters instead of resource IDs.

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `GetReports(ctx, ReportOpts)` | GET | `/api/v2/accounts/{id}/reports` | `(any, error)` |
| `GetReportSummary(ctx, ReportOpts)` | GET | `/api/v2/accounts/{id}/reports/summary` | `(*ReportSummary, error)` |
| `GetConversationMetrics(ctx, ReportOpts)` | GET | `/api/v2/accounts/{id}/reports/conversations` | `(any, error)` |
| `GetAgentConversationMetrics(ctx, ReportOpts)` | GET | `/api/v2/accounts/{id}/reports/conversations/` | `(any, error)` |
| `GetFirstResponseTimeDistribution(ctx, ReportOpts)` | GET | `/api/v2/accounts/{id}/reports/first_response_time_distribution` | `(any, error)` |
| `GetInboxLabelMatrix(ctx, ReportOpts)` | GET | `/api/v2/accounts/{id}/reports/inbox_label_matrix` | `(any, error)` |
| `GetOutgoingMessagesCount(ctx, ReportOpts)` | GET | `/api/v2/accounts/{id}/reports/outgoing_messages_count` | `(any, error)` |
| `GetSummaryByAgent(ctx, ReportOpts)` | GET | `/api/v2/accounts/{id}/summary_reports/agent` | `(any, error)` |
| `GetSummaryByChannel(ctx, ReportOpts)` | GET | `/api/v2/accounts/{id}/summary_reports/channel` | `(any, error)` |
| `GetSummaryByInbox(ctx, ReportOpts)` | GET | `/api/v2/accounts/{id}/summary_reports/inbox` | `(any, error)` |
| `GetSummaryByTeam(ctx, ReportOpts)` | GET | `/api/v2/accounts/{id}/summary_reports/team` | `(any, error)` |
| `GetReportingEvents(ctx, ReportOpts)` | GET | `/api/v1/accounts/{id}/reporting_events` | `(any, error)` |

Note: Most report endpoints return varied response shapes. Using `any` as
return type — the CLI passes the raw decoded JSON through to the contract
envelope. `ReportSummary` is typed because it has a known stable shape.

`GetConversationMetrics` and `GetAgentConversationMetrics` map to two
variants of the same base path. The CLI `reports conversations` command
dispatches to one or the other based on `--type`: `account` uses
`GetConversationMetrics`, `agent` uses `GetAgentConversationMetrics`.
This means 12 API methods map to 11 CLI commands.

`ReportOpts` query parameters:

| Field | Type | Description |
|-------|------|-------------|
| `Metric` | string | `conversations_count`, `incoming_messages_count`, `outgoing_messages_count`, `avg_first_response_time`, `avg_resolution_time`, `resolutions_count` |
| `Type` | string | `account`, `agent`, `inbox`, `label`, `team` |
| `ID` | string | Specific resource ID for agent/inbox/label filtering |
| `Since` | string | Timestamp for range start |
| `Until` | string | Timestamp for range end |

### `webhooks.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListWebhooks(ctx)` | GET | `/api/v1/accounts/{id}/webhooks` | `([]Webhook, error)` |
| `CreateWebhook(ctx, CreateWebhookOpts)` | POST | `/api/v1/accounts/{id}/webhooks` | `(*Webhook, error)` |
| `UpdateWebhook(ctx, webhookID, UpdateWebhookOpts)` | PATCH | `/api/v1/accounts/{id}/webhooks/{wid}` | `(*Webhook, error)` |
| `DeleteWebhook(ctx, webhookID)` | DELETE | `/api/v1/accounts/{id}/webhooks/{wid}` | `error` |

Valid subscription events: `conversation_created`, `conversation_status_changed`,
`conversation_updated`, `contact_created`, `contact_updated`,
`message_created`, `message_updated`, `webwidget_triggered`.

### `automation_rules.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListAutomationRules(ctx)` | GET | `/api/v1/accounts/{id}/automation_rules` | `([]AutomationRule, error)` |
| `GetAutomationRule(ctx, ruleID)` | GET | `/api/v1/accounts/{id}/automation_rules/{rid}` | `(*AutomationRule, error)` |
| `CreateAutomationRule(ctx, CreateAutomationRuleOpts)` | POST | `/api/v1/accounts/{id}/automation_rules` | `(*AutomationRule, error)` |
| `UpdateAutomationRule(ctx, ruleID, UpdateAutomationRuleOpts)` | PATCH | `/api/v1/accounts/{id}/automation_rules/{rid}` | `(*AutomationRule, error)` |
| `DeleteAutomationRule(ctx, ruleID)` | DELETE | `/api/v1/accounts/{id}/automation_rules/{rid}` | `error` |

### `labels.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListLabels(ctx)` | GET | `/api/v1/accounts/{id}/labels` | `([]Label, error)` |
| `GetLabel(ctx, labelID)` | GET | `/api/v1/accounts/{id}/labels/{lid}` | `(*Label, error)` |
| `CreateLabel(ctx, CreateLabelOpts)` | POST | `/api/v1/accounts/{id}/labels` | `(*Label, error)` |
| `UpdateLabel(ctx, labelID, UpdateLabelOpts)` | PATCH | `/api/v1/accounts/{id}/labels/{lid}` | `(*Label, error)` |
| `DeleteLabel(ctx, labelID)` | DELETE | `/api/v1/accounts/{id}/labels/{lid}` | `error` |

Note: These are account-level label definitions. Contact/conversation label
assignment (replace-all) was implemented in Sprint C.

### `custom_attributes.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListCustomAttributes(ctx)` | GET | `/api/v1/accounts/{id}/custom_attribute_definitions` | `([]CustomAttribute, error)` |
| `GetCustomAttribute(ctx, attrID)` | GET | `/api/v1/accounts/{id}/custom_attribute_definitions/{aid}` | `(*CustomAttribute, error)` |
| `CreateCustomAttribute(ctx, CreateCustomAttributeOpts)` | POST | `/api/v1/accounts/{id}/custom_attribute_definitions` | `(*CustomAttribute, error)` |
| `UpdateCustomAttribute(ctx, attrID, UpdateCustomAttributeOpts)` | PATCH | `/api/v1/accounts/{id}/custom_attribute_definitions/{aid}` | `(*CustomAttribute, error)` |
| `DeleteCustomAttribute(ctx, attrID)` | DELETE | `/api/v1/accounts/{id}/custom_attribute_definitions/{aid}` | `error` |

### `custom_filters.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListCustomFilters(ctx, filterType)` | GET | `/api/v1/accounts/{id}/custom_filters?filter_type=T` | `([]CustomFilter, error)` |
| `GetCustomFilter(ctx, filterID)` | GET | `/api/v1/accounts/{id}/custom_filters/{fid}` | `(*CustomFilter, error)` |
| `CreateCustomFilter(ctx, CreateCustomFilterOpts)` | POST | `/api/v1/accounts/{id}/custom_filters` | `(*CustomFilter, error)` |
| `UpdateCustomFilter(ctx, filterID, UpdateCustomFilterOpts)` | PATCH | `/api/v1/accounts/{id}/custom_filters/{fid}` | `(*CustomFilter, error)` |
| `DeleteCustomFilter(ctx, filterID)` | DELETE | `/api/v1/accounts/{id}/custom_filters/{fid}` | `error` |

### `account.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `GetAccount(ctx)` | GET | `/api/v1/accounts/{id}` | `(*AccountInfo, error)` |
| `UpdateAccount(ctx, UpdateAccountOpts)` | PATCH | `/api/v1/accounts/{id}` | `(*AccountInfo, error)` |

### `agent_bots.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListAgentBots(ctx)` | GET | `/api/v1/accounts/{id}/agent_bots` | `([]AccountAgentBot, error)` |
| `GetAgentBot(ctx, botID)` | GET | `/api/v1/accounts/{id}/agent_bots/{bid}` | `(*AccountAgentBot, error)` |
| `CreateAgentBot(ctx, CreateAgentBotOpts)` | POST | `/api/v1/accounts/{id}/agent_bots` | `(*AccountAgentBot, error)` |
| `UpdateAgentBot(ctx, botID, UpdateAgentBotOpts)` | PATCH | `/api/v1/accounts/{id}/agent_bots/{bid}` | `(*AccountAgentBot, error)` |
| `DeleteAgentBot(ctx, botID)` | DELETE | `/api/v1/accounts/{id}/agent_bots/{bid}` | `error` |

Note: Named `AccountAgentBot` to distinguish from the existing `AgentBot`
type (used for inbox agent bot assignment in Sprint C). The existing `AgentBot`
struct only has `ID` and `Name`; the account-scoped version adds more fields.

### `audit_logs.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListAuditLogs(ctx, page)` | GET | `/api/v1/accounts/{id}/audit_logs?page=N` | `([]AuditLog, error)` |

### `integrations.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListIntegrationApps(ctx)` | GET | `/api/v1/accounts/{id}/integrations/apps` | `([]Integration, error)` |
| `CreateIntegrationHook(ctx, CreateIntegrationHookOpts)` | POST | `/api/v1/accounts/{id}/integrations/hooks` | `(*IntegrationHook, error)` |
| `UpdateIntegrationHook(ctx, hookID, UpdateIntegrationHookOpts)` | PATCH | `/api/v1/accounts/{id}/integrations/hooks/{hid}` | `(*IntegrationHook, error)` |
| `DeleteIntegrationHook(ctx, hookID)` | DELETE | `/api/v1/accounts/{id}/integrations/hooks/{hid}` | `error` |

### `help_center.go` (new file)

| Method | HTTP | Path | Returns |
|--------|------|------|---------|
| `ListPortals(ctx)` | GET | `/api/v1/accounts/{id}/portals` | `([]Portal, error)` |
| `CreatePortal(ctx, CreatePortalOpts)` | POST | `/api/v1/accounts/{id}/portals` | `(*Portal, error)` |
| `UpdatePortal(ctx, portalID, UpdatePortalOpts)` | PATCH | `/api/v1/accounts/{id}/portals/{pid}` | `(*Portal, error)` |
| `CreateArticle(ctx, portalID, CreateArticleOpts)` | POST | `/api/v1/accounts/{id}/portals/{pid}/articles` | `(*Article, error)` |
| `CreateCategory(ctx, portalID, CreateCategoryOpts)` | POST | `/api/v1/accounts/{id}/portals/{pid}/categories` | `(*Category, error)` |

## Model Types

All types in `internal/chatwoot/application/models.go`.

### Core Types

```go
// Team represents a Chatwoot team.
type Team struct {
    ID               int    `json:"id"`
    Name             string `json:"name"`
    Description      string `json:"description,omitempty"`
    AllowAutoAssign  bool   `json:"allow_auto_assign"`
    AccountID        int    `json:"account_id"`
    IsMember         bool   `json:"is_member,omitempty"`
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
    ID             int    `json:"id"`
    Title          string `json:"title"`
    Description    string `json:"description,omitempty"`
    Color          string `json:"color,omitempty"`
    ShowOnSidebar  bool   `json:"show_on_sidebar,omitempty"`
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
    AvgFirstResponseTime string `json:"avg_first_response_time"`
    AvgResolutionTime    string `json:"avg_resolution_time"`
    ConversationsCount   int    `json:"conversations_count"`
    IncomingMessagesCount int   `json:"incoming_messages_count"`
    OutgoingMessagesCount int   `json:"outgoing_messages_count"`
    ResolutionsCount     int    `json:"resolutions_count"`
    Previous             any    `json:"previous,omitempty"`
}
```

### Opts Types

```go
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

## CLI Command Packages

Every resource gets its own package under `internal/cli/application/`. Each
follows the Sprint C thin handler pattern: `cmdutil.ResolveContext` →
`cmdutil.ResolveAuth` → transport → API client → contract envelope.

### Standard CRUD Pattern

Most resources follow this exact structure (using `labels` as example):

```
internal/cli/application/labels/
  labels.go          — Cmd group, init registers subcommands
  list.go            — labels list
  get.go             — labels get --id
  create.go          — labels create --title [--description] [--color]
  update.go          — labels update --id [--title] [--description] [--color]
  delete.go          — labels delete --id
  labels_test.go     — tests
  testroot_test.go   — test helper
```

Resources using this pattern: agents, agent-bots, labels, custom-attributes,
custom-filters, canned-responses, webhooks, automation-rules.

### Teams (CRUD + nested members)

Same as standard CRUD plus `members.go` with 4 subcommands on `membersCmd`.
Follows the exact pattern from `inboxes/members.go`.

### Account (singleton, no ID flag)

Only `get` and `update`. No `--id` flag — uses account ID from runtime context.

### Reports (read-only, date-range)

Each report type is its own subcommand with common flags:

| Flag | Type | Description |
|------|------|-------------|
| `--type` | string | Required for some: account, agent, inbox, label, team |
| `--metric` | string | Required for `reports account`: conversations_count, etc. |
| `--since` | string | Timestamp range start |
| `--until` | string | Timestamp range end |
| `--id` | string | Resource ID for agent/inbox/label filtering |

### Integrations (apps + hooks split)

```
internal/cli/application/integrations/
  integrations.go    — Cmd with appsCmd and hooksCmd subgroups
  apps.go            — apps list
  hooks.go           — hooks create/update/delete
```

### Help Center (portal hierarchy)

```
internal/cli/application/helpcenter/
  helpcenter.go      — Cmd with portalsCmd, articlesCmd, categoriesCmd
  portals.go         — portals list/create/update
  articles.go        — articles create --portal-id --title
  categories.go      — categories create --portal-id --name
```

Note: Go package name is `helpcenter` (no hyphen). The CLI command use string
is `help-center` (with hyphen).

### Audit Logs (read-only list)

```
internal/cli/application/auditlogs/
  auditlogs.go       — Cmd (Use: "audit-logs")
  list.go            — audit-logs list (with pagination)
```

Note: Go package name is `auditlogs` (no hyphen). The CLI command use string
is `audit-logs` (with hyphen).

### Flag Conventions

Follow Sprint C conventions. Additional flags for Sprint D:

| Flag | Type | Used by | Description |
|------|------|---------|-------------|
| `--team-id` | int | teams members commands | Parent team ID |
| `--agent-ids` | string | teams members add/update/delete | Comma-separated agent IDs |
| `--short-code` | string | canned-responses create | Template shortcut |
| `--content` | string | canned-responses create | Template body |
| `--url` | string | webhooks create | Webhook URL |
| `--subscriptions` | string | webhooks create | Comma-separated event types |
| `--event-name` | string | automation-rules create | Trigger event |
| `--conditions` | string | automation-rules create | JSON conditions array |
| `--actions` | string | automation-rules create | JSON actions array |
| `--title` | string | labels create | Label title |
| `--color` | string | labels create | Label color hex |
| `--attribute-key` | string | custom-attributes create | Attribute system key |
| `--attribute-model` | string | custom-attributes create | contact or conversation |
| `--attribute-type` | string | custom-attributes create | Data type |
| `--filter-type` | string | custom-filters list | conversation, contact, or report |
| `--query` | string | custom-filters create | JSON query object |
| `--since` | string | reports commands | Date range start |
| `--until` | string | reports commands | Date range end |
| `--metric` | string | reports account | Metric to query |
| `--portal-id` | int | articles/categories create | Parent portal |
| `--app-id` | string | integration hooks create | Integration app ID (e.g. "slack") |
| `--hook-id` | int | integration hooks update/delete | Hook ID |
| `--bot-type` | string | agent-bots create | Bot classification |
| `--bot-config` | string | agent-bots create | JSON bot configuration |
| `--locale` | string | account update, categories create | Language/region |

### Registration

All 14 new command packages registered in
`internal/cli/application/application.go` via `Cmd.AddCommand(...)`.

## Testing Strategy

### API Client Tests

Same pattern as Sprint C: `httptest.NewServer`, verify method/path/body/response.

Coverage priority per resource:
- List method test (verify path, decode response)
- One mutation test (verify method, path, request body)
- Delete test (verify method, path)

### CLI Command Tests

Same pattern as Sprint C: `testroot_test.go`, env vars for auth bypass, httptest
server, verify JSON envelope output.

Coverage priority:
- Every list command
- One create per resource
- Reports: verify `--since`/`--until` flow through to query params

## Task Structure for Implementation

| Task | Content | Complexity |
|------|---------|-----------|
| 1 | Add all model types (~30 structs + opts) | Medium |
| 2 | Teams API client + tests (9 methods) | Medium |
| 3 | Agents API client + tests (4 methods) | Small |
| 4 | Canned Responses API client + tests (4 methods) | Small |
| 5 | Reports API client + tests (12 methods) | Large |
| 6 | Webhooks API client + tests (4 methods) | Small |
| 7 | Automation Rules API client + tests (5 methods) | Small |
| 8 | Simple CRUD API clients batch: Labels, Custom Attributes, Custom Filters, Account, Agent Bots, Audit Logs (26 methods) | Medium |
| 9 | Integrations + Help Center API clients + tests (9 methods) | Small-Medium |
| 10 | Teams CLI commands + tests (9 commands) | Medium |
| 11 | Agents CLI commands + tests (4 commands) | Small |
| 12 | Canned Responses CLI commands + tests (4 commands) | Small |
| 13 | Reports CLI commands + tests (11 commands) | Large |
| 14 | Webhooks CLI commands + tests (4 commands) | Small |
| 15 | Automation Rules CLI commands + tests (5 commands) | Small |
| 16 | Simple CRUD CLI batch: Labels, Custom Attributes, Custom Filters, Account, Agent Bots, Audit Logs | Medium |
| 17 | Integrations + Help Center CLI commands + tests | Small-Medium |
| 18 | Register all subgroups in application.go | Trivial |
| 19 | Full test suite + exit criteria verification | Trivial |

## Exit Criteria

- `go test ./...` all green
- `go vet ./...` clean
- `go build ./cmd/chatwoot/` succeeds
- All 14 resource families visible under `chatwoot application --help`
- `chatwoot application teams list` returns JSON envelope
- `chatwoot application reports summary --type account --since X --until Y` works
- `chatwoot application webhooks create --url X --subscriptions Y` works
- `chatwoot application automation-rules list` returns JSON envelope
- Team member operations follow inbox member pattern
- Existing Sprint C commands still work
- All commands produce valid JSON envelopes on stdout
- All commands use the cmdutil pipeline
- No business logic in command handler files

## Out of Scope

- Platform-scoped agent bots (`/platform/api/v1/agent_bots`) — Sprint E
- Report data visualization or formatting beyond JSON
- Webhook event delivery testing
- Help center article updates/deletes (not in planned command structure)
- Category updates/deletes (not in planned command structure)
- Portal deletion (not in Chatwoot API)
- Custom attribute value assignment on contacts/conversations (would extend
  existing Sprint C commands — separate concern)

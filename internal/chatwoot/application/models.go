package application

import chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"

// Profile represents a Chatwoot agent/admin user profile.
type Profile struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	AccountID int              `json:"account_id,omitempty"`
	Role      string           `json:"role,omitempty"`
	CreatedAt chatwoot.Timestamp `json:"created_at,omitempty"`
}

// UpdateProfileOpts holds optional fields for updating a profile.
// Pointer fields are used so only explicitly set values are serialized.
type UpdateProfileOpts struct {
	Name         *string `json:"name,omitempty"`
	Email        *string `json:"email,omitempty"`
	Availability *string `json:"availability,omitempty"`
}

// Conversation represents a Chatwoot conversation.
type Conversation struct {
	ID          int              `json:"id"`
	AccountID   int              `json:"account_id"`
	InboxID     int              `json:"inbox_id"`
	Status      string           `json:"status"`
	Priority    string           `json:"priority,omitempty"`
	UnreadCount int              `json:"unread_count,omitempty"`
	CreatedAt   chatwoot.Timestamp `json:"created_at,omitempty"`
}

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

// AssignmentResponse represents the API response from conversation assignment.
// The API returns the assigned agent as a flat object, not a conversation.
type AssignmentResponse struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	AvailableName      string `json:"available_name,omitempty"`
	Email              string `json:"email,omitempty"`
	Role               string `json:"role,omitempty"`
	AccountID          int    `json:"account_id,omitempty"`
	AvailabilityStatus string `json:"availability_status,omitempty"`
	AutoOffline        bool   `json:"auto_offline,omitempty"`
	Confirmed          bool   `json:"confirmed,omitempty"`
	Thumbnail          string `json:"thumbnail,omitempty"`
}

// StatusToggleResponse represents the API response from toggle_status.
type StatusToggleResponse struct {
	Success        bool    `json:"success"`
	ConversationID int     `json:"conversation_id"`
	CurrentStatus  string  `json:"current_status"`
	SnoozedUntil   *string `json:"snoozed_until,omitempty"`
}

// PriorityToggleResponse represents the API response from toggle_priority.
type PriorityToggleResponse struct {
	Success        bool   `json:"success"`
	ConversationID int    `json:"conversation_id"`
	CurrentPriority string `json:"current_priority"`
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
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"filter_type"`
	Query     any    `json:"query"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
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
	ID            int                `json:"id"`
	Action        string             `json:"action"`
	AuditableType string             `json:"auditable_type"`
	AuditableID   int                `json:"auditable_id"`
	UserID        int                `json:"user_id,omitempty"`
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
// The Chatwoot API returns time fields as floats (seconds), not strings.
type ReportSummary struct {
	AvgFirstResponseTime  float64         `json:"avg_first_response_time"`
	AvgResolutionTime     float64         `json:"avg_resolution_time"`
	ConversationsCount    int             `json:"conversations_count"`
	IncomingMessagesCount int             `json:"incoming_messages_count"`
	OutgoingMessagesCount int             `json:"outgoing_messages_count"`
	ResolutionsCount      int             `json:"resolutions_count"`
	ReplyTime             float64         `json:"reply_time"`
	Previous              *ReportSummary  `json:"previous,omitempty"`
}

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
	Type  string `json:"filter_type"`
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

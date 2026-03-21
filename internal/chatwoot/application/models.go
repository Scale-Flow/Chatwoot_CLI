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

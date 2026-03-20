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

package clientapi

type Contact struct {
	SourceID    string `json:"source_id"`
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	PubsubToken string `json:"pubsub_token,omitempty"`
}

type CreateContactOpts struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone_number,omitempty"`
}

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

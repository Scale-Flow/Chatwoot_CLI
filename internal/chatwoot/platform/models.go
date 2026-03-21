package platform

type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreateAccountOpts struct {
	Name string `json:"name"`
}

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

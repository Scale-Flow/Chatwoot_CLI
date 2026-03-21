// internal/chatwoot/application/inboxes.go
package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListInboxes returns all inboxes for the account.
func (c *Client) ListInboxes(ctx context.Context) ([]Inbox, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/inboxes", c.accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var body struct {
		Payload []Inbox `json:"payload"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, err
	}
	return body.Payload, nil
}

// GetInbox returns a single inbox by ID.
func (c *Client) GetInbox(ctx context.Context, id int) (*Inbox, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/inboxes/%d", c.accountID, id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var inbox Inbox
	if err := chatwoot.DecodeResponse(resp, &inbox); err != nil {
		return nil, err
	}
	return &inbox, nil
}

// CreateInbox creates a new inbox.
func (c *Client) CreateInbox(ctx context.Context, opts CreateInboxOpts) (*Inbox, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/inboxes", c.accountID)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create inbox: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var inbox Inbox
	if err := chatwoot.DecodeResponse(resp, &inbox); err != nil {
		return nil, err
	}
	return &inbox, nil
}

// UpdateInbox updates an existing inbox.
func (c *Client) UpdateInbox(ctx context.Context, id int, opts UpdateInboxOpts) (*Inbox, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/inboxes/%d", c.accountID, id)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update inbox: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	var inbox Inbox
	if err := chatwoot.DecodeResponse(resp, &inbox); err != nil {
		return nil, err
	}
	return &inbox, nil
}

// ListInboxMembers returns the agents assigned to an inbox.
func (c *Client) ListInboxMembers(ctx context.Context, inboxID int) ([]Agent, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/inbox_members/%d", c.accountID, inboxID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var body struct {
		Payload []Agent `json:"payload"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, err
	}
	return body.Payload, nil
}

// inboxMemberPayload is the request body for inbox member mutations.
type inboxMemberPayload struct {
	InboxID int   `json:"inbox_id"`
	UserIDs []int `json:"user_ids"`
}

// AddInboxMember adds agents to an inbox.
func (c *Client) AddInboxMember(ctx context.Context, inboxID int, agentIDs []int) ([]Agent, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/inbox_members", c.accountID)
	body, err := json.Marshal(inboxMemberPayload{InboxID: inboxID, UserIDs: agentIDs})
	if err != nil {
		return nil, fmt.Errorf("marshal add inbox member: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var respBody struct {
		Payload []Agent `json:"payload"`
	}
	if err := chatwoot.DecodeResponse(resp, &respBody); err != nil {
		return nil, err
	}
	return respBody.Payload, nil
}

// UpdateInboxMembers replaces the agent list for an inbox.
func (c *Client) UpdateInboxMembers(ctx context.Context, inboxID int, agentIDs []int) ([]Agent, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/inbox_members", c.accountID)
	body, err := json.Marshal(inboxMemberPayload{InboxID: inboxID, UserIDs: agentIDs})
	if err != nil {
		return nil, fmt.Errorf("marshal update inbox members: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	var respBody struct {
		Payload []Agent `json:"payload"`
	}
	if err := chatwoot.DecodeResponse(resp, &respBody); err != nil {
		return nil, err
	}
	return respBody.Payload, nil
}

// RemoveInboxMember removes agents from an inbox (DELETE with body).
func (c *Client) RemoveInboxMember(ctx context.Context, inboxID int, agentIDs []int) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/inbox_members", c.accountID)
	body, err := json.Marshal(inboxMemberPayload{InboxID: inboxID, UserIDs: agentIDs})
	if err != nil {
		return fmt.Errorf("marshal remove inbox member: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, body)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

// GetInboxAgentBot returns the agent bot attached to an inbox.
func (c *Client) GetInboxAgentBot(ctx context.Context, inboxID int) (*AgentBot, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/inboxes/%d/agent_bot", c.accountID, inboxID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var bot AgentBot
	if err := chatwoot.DecodeResponse(resp, &bot); err != nil {
		return nil, err
	}
	return &bot, nil
}

// SetInboxAgentBot attaches an agent bot to an inbox.
func (c *Client) SetInboxAgentBot(ctx context.Context, inboxID int, agentBotID int) (*AgentBot, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/inboxes/%d/set_agent_bot", c.accountID, inboxID)
	payload := map[string]int{"agent_bot": agentBotID}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal set agent bot: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var bot AgentBot
	if err := chatwoot.DecodeResponse(resp, &bot); err != nil {
		return nil, err
	}
	return &bot, nil
}

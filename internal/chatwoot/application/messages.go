// internal/chatwoot/application/messages.go
package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListMessages returns messages for a conversation.
func (c *Client) ListMessages(ctx context.Context, conversationID int) ([]Message, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/messages", c.accountID, conversationID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var body struct {
		Payload []Message `json:"payload"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, err
	}
	return body.Payload, nil
}

// CreateMessage creates a new message in a conversation.
func (c *Client) CreateMessage(ctx context.Context, conversationID int, opts CreateMessageOpts) (*Message, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/messages", c.accountID, conversationID)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create message: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var msg Message
	if err := chatwoot.DecodeResponse(resp, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// DeleteMessage deletes a message from a conversation.
func (c *Client) DeleteMessage(ctx context.Context, conversationID, messageID int) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/messages/%d", c.accountID, conversationID, messageID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

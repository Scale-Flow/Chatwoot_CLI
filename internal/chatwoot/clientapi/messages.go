package clientapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func (c *Client) messageBasePath(contactIdentifier string, conversationID int) string {
	return fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s/conversations/%d/messages", c.inboxIdentifier, contactIdentifier, conversationID)
}

func (c *Client) ListMessages(ctx context.Context, contactIdentifier string, conversationID int) ([]Message, error) {
	path := c.messageBasePath(contactIdentifier, conversationID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var messages []Message
	if err := chatwoot.DecodeResponse(resp, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (c *Client) CreateMessage(ctx context.Context, contactIdentifier string, conversationID int, opts CreateMessageOpts) (*Message, error) {
	path := c.messageBasePath(contactIdentifier, conversationID)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create message: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var message Message
	if err := chatwoot.DecodeResponse(resp, &message); err != nil {
		return nil, err
	}
	return &message, nil
}

func (c *Client) UpdateMessage(ctx context.Context, contactIdentifier string, conversationID int, messageID int, opts UpdateMessageOpts) (*Message, error) {
	path := fmt.Sprintf("%s/%d", c.messageBasePath(contactIdentifier, conversationID), messageID)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update message: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	var message Message
	if err := chatwoot.DecodeResponse(resp, &message); err != nil {
		return nil, err
	}
	return &message, nil
}

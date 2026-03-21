package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListWebhooks returns all webhooks for the account.
func (c *Client) ListWebhooks(ctx context.Context) ([]Webhook, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/webhooks", c.accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var body struct {
		Payload []Webhook `json:"payload"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, err
	}
	return body.Payload, nil
}

// CreateWebhook creates a new webhook subscription.
func (c *Client) CreateWebhook(ctx context.Context, opts CreateWebhookOpts) (*Webhook, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/webhooks", c.accountID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create webhook: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, b)
	if err != nil {
		return nil, err
	}
	var webhook Webhook
	if err := chatwoot.DecodeResponse(resp, &webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}

// UpdateWebhook updates an existing webhook.
func (c *Client) UpdateWebhook(ctx context.Context, webhookID int, opts UpdateWebhookOpts) (*Webhook, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/webhooks/%d", c.accountID, webhookID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update webhook: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, b)
	if err != nil {
		return nil, err
	}
	var webhook Webhook
	if err := chatwoot.DecodeResponse(resp, &webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}

// DeleteWebhook deletes a webhook by ID.
func (c *Client) DeleteWebhook(ctx context.Context, webhookID int) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/webhooks/%d", c.accountID, webhookID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

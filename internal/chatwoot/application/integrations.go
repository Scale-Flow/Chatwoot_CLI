package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListIntegrationApps returns all available integration apps for the account.
func (c *Client) ListIntegrationApps(ctx context.Context) ([]Integration, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/integrations/apps", c.accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var body struct {
		Payload []Integration `json:"payload"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, err
	}
	return body.Payload, nil
}

// CreateIntegrationHook creates a new integration hook.
func (c *Client) CreateIntegrationHook(ctx context.Context, opts CreateIntegrationHookOpts) (*IntegrationHook, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/integrations/hooks", c.accountID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create integration hook: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, b)
	if err != nil {
		return nil, err
	}
	var hook IntegrationHook
	if err := chatwoot.DecodeResponse(resp, &hook); err != nil {
		return nil, err
	}
	return &hook, nil
}

// UpdateIntegrationHook updates an existing integration hook.
func (c *Client) UpdateIntegrationHook(ctx context.Context, hookID int, opts UpdateIntegrationHookOpts) (*IntegrationHook, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/integrations/hooks/%d", c.accountID, hookID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update integration hook: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, b)
	if err != nil {
		return nil, err
	}
	var hook IntegrationHook
	if err := chatwoot.DecodeResponse(resp, &hook); err != nil {
		return nil, err
	}
	return &hook, nil
}

// DeleteIntegrationHook deletes an integration hook by ID.
func (c *Client) DeleteIntegrationHook(ctx context.Context, hookID int) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/integrations/hooks/%d", c.accountID, hookID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

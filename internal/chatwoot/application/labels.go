package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListLabels returns all labels for the account.
func (c *Client) ListLabels(ctx context.Context) ([]Label, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/labels", c.accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var body struct {
		Payload []Label `json:"payload"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, err
	}
	return body.Payload, nil
}

// GetLabel returns a single label by ID.
func (c *Client) GetLabel(ctx context.Context, labelID int) (*Label, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/labels/%d", c.accountID, labelID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var label Label
	if err := chatwoot.DecodeResponse(resp, &label); err != nil {
		return nil, err
	}
	return &label, nil
}

// CreateLabel creates a new label.
func (c *Client) CreateLabel(ctx context.Context, opts CreateLabelOpts) (*Label, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/labels", c.accountID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create label: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, b)
	if err != nil {
		return nil, err
	}
	var label Label
	if err := chatwoot.DecodeResponse(resp, &label); err != nil {
		return nil, err
	}
	return &label, nil
}

// UpdateLabel updates an existing label.
func (c *Client) UpdateLabel(ctx context.Context, labelID int, opts UpdateLabelOpts) (*Label, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/labels/%d", c.accountID, labelID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update label: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, b)
	if err != nil {
		return nil, err
	}
	var label Label
	if err := chatwoot.DecodeResponse(resp, &label); err != nil {
		return nil, err
	}
	return &label, nil
}

// DeleteLabel deletes a label by ID.
func (c *Client) DeleteLabel(ctx context.Context, labelID int) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/labels/%d", c.accountID, labelID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

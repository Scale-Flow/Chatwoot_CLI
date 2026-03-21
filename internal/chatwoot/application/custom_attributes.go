package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListCustomAttributes returns all custom attribute definitions for the account.
func (c *Client) ListCustomAttributes(ctx context.Context) ([]CustomAttribute, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_attribute_definitions", c.accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var body struct {
		Data []CustomAttribute `json:"data"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, err
	}
	return body.Data, nil
}

// GetCustomAttribute returns a single custom attribute definition by ID.
func (c *Client) GetCustomAttribute(ctx context.Context, attrID int) (*CustomAttribute, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_attribute_definitions/%d", c.accountID, attrID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var attr CustomAttribute
	if err := chatwoot.DecodeResponse(resp, &attr); err != nil {
		return nil, err
	}
	return &attr, nil
}

// CreateCustomAttribute creates a new custom attribute definition.
func (c *Client) CreateCustomAttribute(ctx context.Context, opts CreateCustomAttributeOpts) (*CustomAttribute, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_attribute_definitions", c.accountID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create custom attribute: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, b)
	if err != nil {
		return nil, err
	}
	var attr CustomAttribute
	if err := chatwoot.DecodeResponse(resp, &attr); err != nil {
		return nil, err
	}
	return &attr, nil
}

// UpdateCustomAttribute updates an existing custom attribute definition.
func (c *Client) UpdateCustomAttribute(ctx context.Context, attrID int, opts UpdateCustomAttributeOpts) (*CustomAttribute, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_attribute_definitions/%d", c.accountID, attrID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update custom attribute: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, b)
	if err != nil {
		return nil, err
	}
	var attr CustomAttribute
	if err := chatwoot.DecodeResponse(resp, &attr); err != nil {
		return nil, err
	}
	return &attr, nil
}

// DeleteCustomAttribute deletes a custom attribute definition by ID.
func (c *Client) DeleteCustomAttribute(ctx context.Context, attrID int) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_attribute_definitions/%d", c.accountID, attrID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

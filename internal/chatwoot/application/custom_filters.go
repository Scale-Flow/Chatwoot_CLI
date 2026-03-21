package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListCustomFilters returns all custom filters for the account.
// If filterType is non-empty, it is appended as the filter_type query parameter.
func (c *Client) ListCustomFilters(ctx context.Context, filterType string) ([]CustomFilter, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_filters", c.accountID)
	if filterType != "" {
		path = path + "?filter_type=" + filterType
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var filters []CustomFilter
	if err := chatwoot.DecodeResponse(resp, &filters); err != nil {
		return nil, err
	}
	return filters, nil
}

// GetCustomFilter returns a single custom filter by ID.
func (c *Client) GetCustomFilter(ctx context.Context, filterID int) (*CustomFilter, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_filters/%d", c.accountID, filterID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var filter CustomFilter
	if err := chatwoot.DecodeResponse(resp, &filter); err != nil {
		return nil, err
	}
	return &filter, nil
}

// CreateCustomFilter creates a new custom filter.
func (c *Client) CreateCustomFilter(ctx context.Context, opts CreateCustomFilterOpts) (*CustomFilter, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_filters", c.accountID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create custom filter: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, b)
	if err != nil {
		return nil, err
	}
	var filter CustomFilter
	if err := chatwoot.DecodeResponse(resp, &filter); err != nil {
		return nil, err
	}
	return &filter, nil
}

// UpdateCustomFilter updates an existing custom filter.
func (c *Client) UpdateCustomFilter(ctx context.Context, filterID int, opts UpdateCustomFilterOpts) (*CustomFilter, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_filters/%d", c.accountID, filterID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update custom filter: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, b)
	if err != nil {
		return nil, err
	}
	var filter CustomFilter
	if err := chatwoot.DecodeResponse(resp, &filter); err != nil {
		return nil, err
	}
	return &filter, nil
}

// DeleteCustomFilter deletes a custom filter by ID.
func (c *Client) DeleteCustomFilter(ctx context.Context, filterID int) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_filters/%d", c.accountID, filterID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

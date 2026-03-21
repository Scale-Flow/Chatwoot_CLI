package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListCannedResponses returns all canned responses for the account.
func (c *Client) ListCannedResponses(ctx context.Context) ([]CannedResponse, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/canned_responses", c.accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var canned []CannedResponse
	if err := chatwoot.DecodeResponse(resp, &canned); err != nil {
		return nil, err
	}
	return canned, nil
}

// CreateCannedResponse creates a new canned response.
func (c *Client) CreateCannedResponse(ctx context.Context, opts CreateCannedResponseOpts) (*CannedResponse, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/canned_responses", c.accountID)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create canned response: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var canned CannedResponse
	if err := chatwoot.DecodeResponse(resp, &canned); err != nil {
		return nil, err
	}
	return &canned, nil
}

// UpdateCannedResponse updates an existing canned response.
func (c *Client) UpdateCannedResponse(ctx context.Context, cannedID int, opts UpdateCannedResponseOpts) (*CannedResponse, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/canned_responses/%d", c.accountID, cannedID)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update canned response: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	var canned CannedResponse
	if err := chatwoot.DecodeResponse(resp, &canned); err != nil {
		return nil, err
	}
	return &canned, nil
}

// DeleteCannedResponse deletes a canned response by ID.
func (c *Client) DeleteCannedResponse(ctx context.Context, cannedID int) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/canned_responses/%d", c.accountID, cannedID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

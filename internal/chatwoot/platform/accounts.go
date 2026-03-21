package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func (c *Client) GetAccount(ctx context.Context, id int) (*Account, error) {
	path := fmt.Sprintf("/platform/api/v1/accounts/%d", id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var account Account
	if err := chatwoot.DecodeResponse(resp, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

func (c *Client) CreateAccount(ctx context.Context, opts CreateAccountOpts) (*Account, error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create account: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, "/platform/api/v1/accounts", body)
	if err != nil {
		return nil, err
	}
	var account Account
	if err := chatwoot.DecodeResponse(resp, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

func (c *Client) UpdateAccount(ctx context.Context, id int, opts UpdateAccountOpts) (*Account, error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update account: %w", err)
	}
	path := fmt.Sprintf("/platform/api/v1/accounts/%d", id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	var account Account
	if err := chatwoot.DecodeResponse(resp, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

func (c *Client) DeleteAccount(ctx context.Context, id int) error {
	path := fmt.Sprintf("/platform/api/v1/accounts/%d", id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

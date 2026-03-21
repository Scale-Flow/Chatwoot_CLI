package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// GetAccount returns the account details.
func (c *Client) GetAccount(ctx context.Context) (*AccountInfo, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d", c.accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var account AccountInfo
	if err := chatwoot.DecodeResponse(resp, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// UpdateAccount updates the account details.
func (c *Client) UpdateAccount(ctx context.Context, opts UpdateAccountOpts) (*AccountInfo, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d", c.accountID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update account: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, b)
	if err != nil {
		return nil, err
	}
	var account AccountInfo
	if err := chatwoot.DecodeResponse(resp, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

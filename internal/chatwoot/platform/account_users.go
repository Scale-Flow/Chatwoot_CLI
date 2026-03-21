package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func (c *Client) ListAccountUsers(ctx context.Context, accountID int) ([]AccountUser, error) {
	path := fmt.Sprintf("/platform/api/v1/accounts/%d/account_users", accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var users []AccountUser
	if err := chatwoot.DecodeResponse(resp, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (c *Client) CreateAccountUser(ctx context.Context, accountID int, opts CreateAccountUserOpts) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return fmt.Errorf("marshal create account user: %w", err)
	}
	path := fmt.Sprintf("/platform/api/v1/accounts/%d/account_users", accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

func (c *Client) DeleteAccountUser(ctx context.Context, accountID int, userID int) error {
	body, err := json.Marshal(map[string]int{"user_id": userID})
	if err != nil {
		return fmt.Errorf("marshal delete account user: %w", err)
	}
	path := fmt.Sprintf("/platform/api/v1/accounts/%d/account_users", accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, body)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func (c *Client) CreateUser(ctx context.Context, opts CreateUserOpts) (*User, error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create user: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, "/platform/api/v1/users", body)
	if err != nil {
		return nil, err
	}
	var user User
	if err := chatwoot.DecodeResponse(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) GetUser(ctx context.Context, id int) (*User, error) {
	path := fmt.Sprintf("/platform/api/v1/users/%d", id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var user User
	if err := chatwoot.DecodeResponse(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) UpdateUser(ctx context.Context, id int, opts UpdateUserOpts) (*User, error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update user: %w", err)
	}
	path := fmt.Sprintf("/platform/api/v1/users/%d", id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	var user User
	if err := chatwoot.DecodeResponse(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) DeleteUser(ctx context.Context, id int) error {
	path := fmt.Sprintf("/platform/api/v1/users/%d", id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

func (c *Client) GetUserSSOLink(ctx context.Context, id int) (*SSOLink, error) {
	path := fmt.Sprintf("/platform/api/v1/users/%d/login", id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var link SSOLink
	if err := chatwoot.DecodeResponse(resp, &link); err != nil {
		return nil, err
	}
	return &link, nil
}

package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// UpdateProfile updates the authenticated user's profile.
func (c *Client) UpdateProfile(ctx context.Context, opts UpdateProfileOpts) (*Profile, error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update profile: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, "/api/v1/profile", body)
	if err != nil {
		return nil, err
	}
	var profile Profile
	if err := chatwoot.DecodeResponse(resp, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// GetProfile retrieves the authenticated user's profile.
func (c *Client) GetProfile(ctx context.Context) (*Profile, error) {
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, "/api/v1/profile", nil)
	if err != nil {
		return nil, err
	}
	var profile Profile
	if err := chatwoot.DecodeResponse(resp, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

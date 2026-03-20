package application

import (
	"context"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

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

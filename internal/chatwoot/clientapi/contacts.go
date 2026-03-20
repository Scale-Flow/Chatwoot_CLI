package clientapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func (c *Client) CreateContact(ctx context.Context, opts CreateContactOpts) (*Contact, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts", c.inboxIdentifier)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create contact: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var contact Contact
	if err := chatwoot.DecodeResponse(resp, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

func (c *Client) GetContact(ctx context.Context, contactIdentifier string) (*Contact, error) {
	path := fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s", c.inboxIdentifier, contactIdentifier)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var contact Contact
	if err := chatwoot.DecodeResponse(resp, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

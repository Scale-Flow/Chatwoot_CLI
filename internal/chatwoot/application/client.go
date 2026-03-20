// Package application provides a typed client for the Chatwoot Application API.
package application

import chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"

// Client wraps the shared transport for Application API calls.
type Client struct {
	transport *chatwoot.Client
	accountID int
}

// NewClient creates an Application API client.
func NewClient(transport *chatwoot.Client, accountID int) *Client {
	return &Client{transport: transport, accountID: accountID}
}

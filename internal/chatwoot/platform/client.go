package platform

import chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"

type Client struct {
	transport *chatwoot.Client
}

func NewClient(transport *chatwoot.Client) *Client {
	return &Client{transport: transport}
}

package clientapi

import chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"

type Client struct {
	transport       *chatwoot.Client
	inboxIdentifier string
}

func NewClient(transport *chatwoot.Client, inboxIdentifier string) *Client {
	return &Client{transport: transport, inboxIdentifier: inboxIdentifier}
}

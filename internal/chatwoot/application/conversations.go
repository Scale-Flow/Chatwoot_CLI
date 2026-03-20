// internal/chatwoot/application/conversations.go
package application

import (
	"context"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListConversationsOpts holds optional filters for listing conversations.
type ListConversationsOpts struct {
	Page    int
	Status  string
	InboxID int
}

// ListConversations returns a page of conversations for the account.
func (c *Client) ListConversations(ctx context.Context, opts ListConversationsOpts) ([]Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations?page=%d", c.accountID, opts.Page)
	if opts.Status != "" {
		path += "&status=" + opts.Status
	}
	if opts.InboxID != 0 {
		path += fmt.Sprintf("&inbox_id=%d", opts.InboxID)
	}

	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var body struct {
		Data struct {
			Payload []Conversation `json:"payload"`
		} `json:"data"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, err
	}
	return body.Data.Payload, nil
}

// GetConversation returns a single conversation by ID.
func (c *Client) GetConversation(ctx context.Context, id int) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d", c.accountID, id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var convo Conversation
	if err := chatwoot.DecodeResponse(resp, &convo); err != nil {
		return nil, err
	}
	return &convo, nil
}

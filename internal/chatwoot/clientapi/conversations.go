package clientapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func (c *Client) conversationBasePath(contactIdentifier string) string {
	return fmt.Sprintf("/public/api/v1/inboxes/%s/contacts/%s/conversations", c.inboxIdentifier, contactIdentifier)
}

func (c *Client) ListConversations(ctx context.Context, contactIdentifier string) ([]Conversation, error) {
	path := c.conversationBasePath(contactIdentifier)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var conversations []Conversation
	if err := chatwoot.DecodeResponse(resp, &conversations); err != nil {
		return nil, err
	}
	return conversations, nil
}

func (c *Client) GetConversation(ctx context.Context, contactIdentifier string, conversationID int) (*Conversation, error) {
	path := fmt.Sprintf("%s/%d", c.conversationBasePath(contactIdentifier), conversationID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var conversation Conversation
	if err := chatwoot.DecodeResponse(resp, &conversation); err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (c *Client) CreateConversation(ctx context.Context, contactIdentifier string) (*Conversation, error) {
	path := c.conversationBasePath(contactIdentifier)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}
	var conversation Conversation
	if err := chatwoot.DecodeResponse(resp, &conversation); err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (c *Client) ToggleStatus(ctx context.Context, contactIdentifier string, conversationID int) (*Conversation, error) {
	path := fmt.Sprintf("%s/%d/toggle_status", c.conversationBasePath(contactIdentifier), conversationID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}
	var conversation Conversation
	if err := chatwoot.DecodeResponse(resp, &conversation); err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (c *Client) ToggleTyping(ctx context.Context, contactIdentifier string, conversationID int, opts ToggleTypingOpts) error {
	path := fmt.Sprintf("%s/%d/toggle_typing", c.conversationBasePath(contactIdentifier), conversationID)
	body, err := json.Marshal(opts)
	if err != nil {
		return fmt.Errorf("marshal toggle typing: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

func (c *Client) UpdateLastSeen(ctx context.Context, contactIdentifier string, conversationID int) error {
	path := fmt.Sprintf("%s/%d/update_last_seen", c.conversationBasePath(contactIdentifier), conversationID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

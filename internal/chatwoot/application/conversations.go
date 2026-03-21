// internal/chatwoot/application/conversations.go
package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
)

// ListConversations returns a page of conversations for the account.
func (c *Client) ListConversations(ctx context.Context, opts ListConversationsOpts) ([]Conversation, *contract.Pagination, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations?page=%d", c.accountID, opts.Page)
	if opts.PerPage != 0 {
		path += fmt.Sprintf("&per_page=%d", opts.PerPage)
	}
	if opts.Status != "" {
		path += "&status=" + opts.Status
	}
	if opts.InboxID != 0 {
		path += fmt.Sprintf("&inbox_id=%d", opts.InboxID)
	}

	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body struct {
		Data struct {
			Payload []Conversation `json:"payload"`
			Meta    struct {
				AllCount    int `json:"all_count"`
				CurrentPage int `json:"current_page"`
			} `json:"meta"`
		} `json:"data"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, nil, err
	}

	pg := &contract.Pagination{
		Page:       opts.Page,
		TotalCount: body.Data.Meta.AllCount,
	}
	if body.Data.Meta.CurrentPage != 0 {
		pg.Page = body.Data.Meta.CurrentPage
	}

	return body.Data.Payload, pg, nil
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

// CreateConversation creates a new conversation.
func (c *Client) CreateConversation(ctx context.Context, opts CreateConversationOpts) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations", c.accountID)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create conversation: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var convo Conversation
	if err := chatwoot.DecodeResponse(resp, &convo); err != nil {
		return nil, err
	}
	return &convo, nil
}

// UpdateConversation updates an existing conversation.
func (c *Client) UpdateConversation(ctx context.Context, id int, opts UpdateConversationOpts) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d", c.accountID, id)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update conversation: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	var convo Conversation
	if err := chatwoot.DecodeResponse(resp, &convo); err != nil {
		return nil, err
	}
	return &convo, nil
}

// FilterConversations filters conversations using a payload of filter rules.
func (c *Client) FilterConversations(ctx context.Context, opts FilterConversationsOpts) ([]Conversation, *contract.Pagination, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/filter", c.accountID)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal filter conversations: %w", err)
	}

	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, nil, err
	}

	var respBody struct {
		Data struct {
			Payload []Conversation `json:"payload"`
			Meta    struct {
				AllCount    int `json:"all_count"`
				CurrentPage int `json:"current_page"`
			} `json:"meta"`
		} `json:"data"`
	}
	if err := chatwoot.DecodeResponse(resp, &respBody); err != nil {
		return nil, nil, err
	}

	pg := &contract.Pagination{
		Page:       opts.Page,
		TotalCount: respBody.Data.Meta.AllCount,
	}
	if respBody.Data.Meta.CurrentPage != 0 {
		pg.Page = respBody.Data.Meta.CurrentPage
	}

	return respBody.Data.Payload, pg, nil
}

// GetConversationMeta returns conversation count metadata for the account.
func (c *Client) GetConversationMeta(ctx context.Context) (*ConversationMeta, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/meta", c.accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var meta ConversationMeta
	if err := chatwoot.DecodeResponse(resp, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

// ToggleConversationStatus toggles a conversation's status.
func (c *Client) ToggleConversationStatus(ctx context.Context, id int, status string) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/toggle_status", c.accountID, id)
	payload := map[string]string{"status": status}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal toggle status: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var convo Conversation
	if err := chatwoot.DecodeResponse(resp, &convo); err != nil {
		return nil, err
	}
	return &convo, nil
}

// ToggleConversationPriority toggles a conversation's priority.
func (c *Client) ToggleConversationPriority(ctx context.Context, id int, priority string) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/toggle_priority", c.accountID, id)
	payload := map[string]string{"priority": priority}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal toggle priority: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var convo Conversation
	if err := chatwoot.DecodeResponse(resp, &convo); err != nil {
		return nil, err
	}
	return &convo, nil
}

// AssignConversation assigns a conversation to an agent and/or team.
func (c *Client) AssignConversation(ctx context.Context, id int, opts AssignOpts) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/assignments", c.accountID, id)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal assign conversation: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var convo Conversation
	if err := chatwoot.DecodeResponse(resp, &convo); err != nil {
		return nil, err
	}
	return &convo, nil
}

// ListConversationLabels returns the labels for a conversation.
func (c *Client) ListConversationLabels(ctx context.Context, id int) ([]string, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/labels", c.accountID, id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var labels []string
	if err := chatwoot.DecodeResponse(resp, &labels); err != nil {
		return nil, err
	}
	return labels, nil
}

// SetConversationLabels replaces all labels on a conversation.
func (c *Client) SetConversationLabels(ctx context.Context, id int, labels []string) ([]string, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/labels", c.accountID, id)
	payload := map[string][]string{"labels": labels}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal set labels: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var result []string
	if err := chatwoot.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

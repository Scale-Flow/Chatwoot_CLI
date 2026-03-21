package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListAgentBots returns all account-scoped agent bots.
func (c *Client) ListAgentBots(ctx context.Context) ([]AccountAgentBot, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/agent_bots", c.accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var bots []AccountAgentBot
	if err := chatwoot.DecodeResponse(resp, &bots); err != nil {
		return nil, err
	}
	return bots, nil
}

// GetAgentBot returns a single account-scoped agent bot by ID.
func (c *Client) GetAgentBot(ctx context.Context, botID int) (*AccountAgentBot, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/agent_bots/%d", c.accountID, botID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var bot AccountAgentBot
	if err := chatwoot.DecodeResponse(resp, &bot); err != nil {
		return nil, err
	}
	return &bot, nil
}

// CreateAgentBot creates a new account-scoped agent bot.
func (c *Client) CreateAgentBot(ctx context.Context, opts CreateAgentBotOpts) (*AccountAgentBot, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/agent_bots", c.accountID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create agent bot: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, b)
	if err != nil {
		return nil, err
	}
	var bot AccountAgentBot
	if err := chatwoot.DecodeResponse(resp, &bot); err != nil {
		return nil, err
	}
	return &bot, nil
}

// UpdateAgentBot updates an existing account-scoped agent bot.
func (c *Client) UpdateAgentBot(ctx context.Context, botID int, opts UpdateAgentBotOpts) (*AccountAgentBot, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/agent_bots/%d", c.accountID, botID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update agent bot: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, b)
	if err != nil {
		return nil, err
	}
	var bot AccountAgentBot
	if err := chatwoot.DecodeResponse(resp, &bot); err != nil {
		return nil, err
	}
	return &bot, nil
}

// DeleteAgentBot deletes an account-scoped agent bot by ID.
func (c *Client) DeleteAgentBot(ctx context.Context, botID int) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/agent_bots/%d", c.accountID, botID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

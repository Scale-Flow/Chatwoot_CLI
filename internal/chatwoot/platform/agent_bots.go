package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

func (c *Client) ListAgentBots(ctx context.Context) ([]AgentBot, error) {
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, "/platform/api/v1/agent_bots", nil)
	if err != nil {
		return nil, err
	}
	var bots []AgentBot
	if err := chatwoot.DecodeResponse(resp, &bots); err != nil {
		return nil, err
	}
	return bots, nil
}

func (c *Client) GetAgentBot(ctx context.Context, id int) (*AgentBot, error) {
	path := fmt.Sprintf("/platform/api/v1/agent_bots/%d", id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var bot AgentBot
	if err := chatwoot.DecodeResponse(resp, &bot); err != nil {
		return nil, err
	}
	return &bot, nil
}

func (c *Client) CreateAgentBot(ctx context.Context, opts CreateAgentBotOpts) (*AgentBot, error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create agent bot: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, "/platform/api/v1/agent_bots", body)
	if err != nil {
		return nil, err
	}
	var bot AgentBot
	if err := chatwoot.DecodeResponse(resp, &bot); err != nil {
		return nil, err
	}
	return &bot, nil
}

func (c *Client) UpdateAgentBot(ctx context.Context, id int, opts UpdateAgentBotOpts) (*AgentBot, error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update agent bot: %w", err)
	}
	path := fmt.Sprintf("/platform/api/v1/agent_bots/%d", id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	var bot AgentBot
	if err := chatwoot.DecodeResponse(resp, &bot); err != nil {
		return nil, err
	}
	return &bot, nil
}

func (c *Client) DeleteAgentBot(ctx context.Context, id int) error {
	path := fmt.Sprintf("/platform/api/v1/agent_bots/%d", id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

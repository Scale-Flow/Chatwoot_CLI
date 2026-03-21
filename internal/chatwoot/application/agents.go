package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListAgents returns all agents for the account.
func (c *Client) ListAgents(ctx context.Context) ([]Agent, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/agents", c.accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var agents []Agent
	if err := chatwoot.DecodeResponse(resp, &agents); err != nil {
		return nil, err
	}
	return agents, nil
}

// CreateAgent creates a new agent.
func (c *Client) CreateAgent(ctx context.Context, opts CreateAgentOpts) (*Agent, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/agents", c.accountID)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create agent: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var agent Agent
	if err := chatwoot.DecodeResponse(resp, &agent); err != nil {
		return nil, err
	}
	return &agent, nil
}

// UpdateAgent updates an existing agent.
func (c *Client) UpdateAgent(ctx context.Context, agentID int, opts UpdateAgentOpts) (*Agent, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/agents/%d", c.accountID, agentID)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update agent: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	var agent Agent
	if err := chatwoot.DecodeResponse(resp, &agent); err != nil {
		return nil, err
	}
	return &agent, nil
}

// DeleteAgent deletes an agent by ID.
func (c *Client) DeleteAgent(ctx context.Context, agentID int) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/agents/%d", c.accountID, agentID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

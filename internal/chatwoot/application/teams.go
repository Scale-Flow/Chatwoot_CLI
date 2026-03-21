package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListTeams returns all teams for the account.
func (c *Client) ListTeams(ctx context.Context) ([]Team, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/teams", c.accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var teams []Team
	if err := chatwoot.DecodeResponse(resp, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

// GetTeam returns a single team by ID.
func (c *Client) GetTeam(ctx context.Context, teamID int) (*Team, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/teams/%d", c.accountID, teamID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var team Team
	if err := chatwoot.DecodeResponse(resp, &team); err != nil {
		return nil, err
	}
	return &team, nil
}

// CreateTeam creates a new team.
func (c *Client) CreateTeam(ctx context.Context, opts CreateTeamOpts) (*Team, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/teams", c.accountID)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create team: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var team Team
	if err := chatwoot.DecodeResponse(resp, &team); err != nil {
		return nil, err
	}
	return &team, nil
}

// UpdateTeam updates an existing team.
func (c *Client) UpdateTeam(ctx context.Context, teamID int, opts UpdateTeamOpts) (*Team, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/teams/%d", c.accountID, teamID)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update team: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	var team Team
	if err := chatwoot.DecodeResponse(resp, &team); err != nil {
		return nil, err
	}
	return &team, nil
}

// DeleteTeam deletes a team by ID.
func (c *Client) DeleteTeam(ctx context.Context, teamID int) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/teams/%d", c.accountID, teamID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

// ListTeamMembers returns all agents who are members of a team.
func (c *Client) ListTeamMembers(ctx context.Context, teamID int) ([]Agent, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/teams/%d/team_members", c.accountID, teamID)
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

// AddTeamMember adds agents to a team by agent IDs.
func (c *Client) AddTeamMember(ctx context.Context, teamID int, agentIDs []int) ([]Agent, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/teams/%d/team_members", c.accountID, teamID)
	payload := map[string][]int{"user_ids": agentIDs}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal add team member: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var agents []Agent
	if err := chatwoot.DecodeResponse(resp, &agents); err != nil {
		return nil, err
	}
	return agents, nil
}

// UpdateTeamMembers replaces the full member list of a team.
func (c *Client) UpdateTeamMembers(ctx context.Context, teamID int, agentIDs []int) ([]Agent, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/teams/%d/team_members", c.accountID, teamID)
	payload := map[string][]int{"user_ids": agentIDs}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal update team members: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	var agents []Agent
	if err := chatwoot.DecodeResponse(resp, &agents); err != nil {
		return nil, err
	}
	return agents, nil
}

// RemoveTeamMember removes agents from a team by agent IDs.
func (c *Client) RemoveTeamMember(ctx context.Context, teamID int, agentIDs []int) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/teams/%d/team_members", c.accountID, teamID)
	payload := map[string][]int{"user_ids": agentIDs}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal remove team member: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, body)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

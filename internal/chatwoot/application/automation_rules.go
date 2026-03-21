package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListAutomationRules returns all automation rules for the account.
func (c *Client) ListAutomationRules(ctx context.Context) ([]AutomationRule, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/automation_rules", c.accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var body struct {
		Payload []AutomationRule `json:"payload"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, err
	}
	return body.Payload, nil
}

// GetAutomationRule returns a single automation rule by ID.
func (c *Client) GetAutomationRule(ctx context.Context, ruleID int) (*AutomationRule, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/automation_rules/%d", c.accountID, ruleID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var rule AutomationRule
	if err := chatwoot.DecodeResponse(resp, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// CreateAutomationRule creates a new automation rule.
func (c *Client) CreateAutomationRule(ctx context.Context, opts CreateAutomationRuleOpts) (*AutomationRule, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/automation_rules", c.accountID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create automation rule: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, b)
	if err != nil {
		return nil, err
	}
	var rule AutomationRule
	if err := chatwoot.DecodeResponse(resp, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// UpdateAutomationRule updates an existing automation rule.
func (c *Client) UpdateAutomationRule(ctx context.Context, ruleID int, opts UpdateAutomationRuleOpts) (*AutomationRule, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/automation_rules/%d", c.accountID, ruleID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update automation rule: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, b)
	if err != nil {
		return nil, err
	}
	var rule AutomationRule
	if err := chatwoot.DecodeResponse(resp, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// DeleteAutomationRule deletes an automation rule by ID.
func (c *Client) DeleteAutomationRule(ctx context.Context, ruleID int) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/automation_rules/%d", c.accountID, ruleID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

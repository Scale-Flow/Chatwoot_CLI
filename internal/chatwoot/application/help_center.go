package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListPortals returns all help center portals for the account.
func (c *Client) ListPortals(ctx context.Context) ([]Portal, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/portals", c.accountID)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var body struct {
		Payload []Portal `json:"payload"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, err
	}
	return body.Payload, nil
}

// CreatePortal creates a new help center portal.
func (c *Client) CreatePortal(ctx context.Context, opts CreatePortalOpts) (*Portal, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/portals", c.accountID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create portal: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, b)
	if err != nil {
		return nil, err
	}
	var portal Portal
	if err := chatwoot.DecodeResponse(resp, &portal); err != nil {
		return nil, err
	}
	return &portal, nil
}

// UpdatePortal updates an existing help center portal.
func (c *Client) UpdatePortal(ctx context.Context, portalID int, opts UpdatePortalOpts) (*Portal, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/portals/%d", c.accountID, portalID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update portal: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPatch, path, b)
	if err != nil {
		return nil, err
	}
	var portal Portal
	if err := chatwoot.DecodeResponse(resp, &portal); err != nil {
		return nil, err
	}
	return &portal, nil
}

// CreateArticle creates a new article in the specified portal.
func (c *Client) CreateArticle(ctx context.Context, portalID int, opts CreateArticleOpts) (*Article, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/portals/%d/articles", c.accountID, portalID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create article: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, b)
	if err != nil {
		return nil, err
	}
	var article Article
	if err := chatwoot.DecodeResponse(resp, &article); err != nil {
		return nil, err
	}
	return &article, nil
}

// CreateCategory creates a new category in the specified portal.
func (c *Client) CreateCategory(ctx context.Context, portalID int, opts CreateCategoryOpts) (*Category, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/portals/%d/categories", c.accountID, portalID)
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create category: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, b)
	if err != nil {
		return nil, err
	}
	var category Category
	if err := chatwoot.DecodeResponse(resp, &category); err != nil {
		return nil, err
	}
	return &category, nil
}

// internal/chatwoot/application/contacts.go
package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
)

// ListContacts returns a page of contacts for the account.
func (c *Client) ListContacts(ctx context.Context, opts ListContactsOpts) ([]Contact, *contract.Pagination, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts?page=%d", c.accountID, opts.Page)
	if opts.PerPage != 0 {
		path += fmt.Sprintf("&per_page=%d", opts.PerPage)
	}

	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body struct {
		Data struct {
			Payload []Contact `json:"payload"`
			Meta    struct {
				Count       int `json:"count"`
				CurrentPage int `json:"current_page"`
			} `json:"meta"`
		} `json:"data"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, nil, err
	}

	pg := &contract.Pagination{
		Page:       opts.Page,
		TotalCount: body.Data.Meta.Count,
	}
	if body.Data.Meta.CurrentPage != 0 {
		pg.Page = body.Data.Meta.CurrentPage
	}

	return body.Data.Payload, pg, nil
}

// GetContact returns a single contact by ID.
func (c *Client) GetContact(ctx context.Context, id int) (*Contact, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts/%d", c.accountID, id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var contact Contact
	if err := chatwoot.DecodeResponse(resp, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

// CreateContact creates a new contact.
func (c *Client) CreateContact(ctx context.Context, opts CreateContactOpts) (*Contact, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts", c.accountID)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal create contact: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var contact Contact
	if err := chatwoot.DecodeResponse(resp, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

// UpdateContact updates an existing contact.
func (c *Client) UpdateContact(ctx context.Context, id int, opts UpdateContactOpts) (*Contact, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts/%d", c.accountID, id)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("marshal update contact: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	var contact Contact
	if err := chatwoot.DecodeResponse(resp, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

// DeleteContact deletes a contact by ID.
func (c *Client) DeleteContact(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts/%d", c.accountID, id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return chatwoot.DecodeResponse(resp, nil)
}

// SearchContacts searches contacts by query string.
func (c *Client) SearchContacts(ctx context.Context, query string, page int) ([]Contact, *contract.Pagination, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts/search?q=%s&page=%d", c.accountID, query, page)

	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body struct {
		Data struct {
			Payload []Contact `json:"payload"`
			Meta    struct {
				Count       int `json:"count"`
				CurrentPage int `json:"current_page"`
			} `json:"meta"`
		} `json:"data"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, nil, err
	}

	pg := &contract.Pagination{
		Page:       page,
		TotalCount: body.Data.Meta.Count,
	}
	if body.Data.Meta.CurrentPage != 0 {
		pg.Page = body.Data.Meta.CurrentPage
	}

	return body.Data.Payload, pg, nil
}

// FilterContacts filters contacts using a payload of filter rules.
func (c *Client) FilterContacts(ctx context.Context, opts FilterContactsOpts) ([]Contact, *contract.Pagination, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts/filter", c.accountID)
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal filter contacts: %w", err)
	}

	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, nil, err
	}

	var respBody struct {
		Data struct {
			Payload []Contact `json:"payload"`
			Meta    struct {
				Count       int `json:"count"`
				CurrentPage int `json:"current_page"`
			} `json:"meta"`
		} `json:"data"`
	}
	if err := chatwoot.DecodeResponse(resp, &respBody); err != nil {
		return nil, nil, err
	}

	pg := &contract.Pagination{
		Page:       opts.Page,
		TotalCount: respBody.Data.Meta.Count,
	}
	if respBody.Data.Meta.CurrentPage != 0 {
		pg.Page = respBody.Data.Meta.CurrentPage
	}

	return respBody.Data.Payload, pg, nil
}

// MergeContacts merges two contacts, keeping the base and discarding the mergee.
func (c *Client) MergeContacts(ctx context.Context, baseID, mergeID int) (*Contact, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/actions/contact_merge", c.accountID)
	payload := map[string]int{
		"base_contact_id":   baseID,
		"mergee_contact_id": mergeID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal merge contacts: %w", err)
	}
	resp, err := c.transport.DoWithRetry(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	var contact Contact
	if err := chatwoot.DecodeResponse(resp, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

// ListContactLabels returns the labels for a contact.
func (c *Client) ListContactLabels(ctx context.Context, id int) ([]string, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts/%d/labels", c.accountID, id)
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

// SetContactLabels replaces all labels on a contact.
func (c *Client) SetContactLabels(ctx context.Context, id int, labels []string) ([]string, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts/%d/labels", c.accountID, id)
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

// ListContactConversations returns conversations for a contact.
func (c *Client) ListContactConversations(ctx context.Context, id int) ([]Conversation, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/contacts/%d/conversations", c.accountID, id)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var body struct {
		Payload []Conversation `json:"payload"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, err
	}
	return body.Payload, nil
}

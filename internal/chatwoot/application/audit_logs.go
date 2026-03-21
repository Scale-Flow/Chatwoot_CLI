package application

import (
	"context"
	"fmt"
	"net/http"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// ListAuditLogs returns a paginated list of audit logs for the account.
func (c *Client) ListAuditLogs(ctx context.Context, page int) ([]AuditLog, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/audit_logs?page=%d", c.accountID, page)
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var body struct {
		Payload []AuditLog `json:"payload"`
	}
	if err := chatwoot.DecodeResponse(resp, &body); err != nil {
		return nil, err
	}
	return body.Payload, nil
}

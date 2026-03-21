package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
)

// buildReportQuery constructs a URL query string from ReportOpts.
func (c *Client) buildReportQuery(opts ReportOpts) string {
	params := url.Values{}
	if opts.Metric != "" {
		params.Set("metric", opts.Metric)
	}
	if opts.Type != "" {
		params.Set("type", opts.Type)
	}
	if opts.ID != "" {
		params.Set("id", opts.ID)
	}
	if opts.Since != "" {
		params.Set("since", opts.Since)
	}
	if opts.Until != "" {
		params.Set("until", opts.Until)
	}
	if len(params) > 0 {
		return "?" + params.Encode()
	}
	return ""
}

// getReport is a shared helper for all any-returning GET report endpoints.
func (c *Client) getReport(ctx context.Context, path string) (any, error) {
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var raw json.RawMessage
	if err := chatwoot.DecodeResponse(resp, &raw); err != nil {
		return nil, err
	}
	var result any
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("decode report: %w", err)
	}
	return result, nil
}

// GetReports returns report data for the account.
func (c *Client) GetReports(ctx context.Context, opts ReportOpts) (any, error) {
	path := fmt.Sprintf("/api/v2/accounts/%d/reports%s", c.accountID, c.buildReportQuery(opts))
	return c.getReport(ctx, path)
}

// GetReportSummary returns an account report summary.
func (c *Client) GetReportSummary(ctx context.Context, opts ReportOpts) (*ReportSummary, error) {
	path := fmt.Sprintf("/api/v2/accounts/%d/reports/summary%s", c.accountID, c.buildReportQuery(opts))
	resp, err := c.transport.DoWithRetry(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var summary ReportSummary
	if err := chatwoot.DecodeResponse(resp, &summary); err != nil {
		return nil, err
	}
	return &summary, nil
}

// GetConversationMetrics returns conversation metrics.
func (c *Client) GetConversationMetrics(ctx context.Context, opts ReportOpts) (any, error) {
	path := fmt.Sprintf("/api/v2/accounts/%d/reports/conversations%s", c.accountID, c.buildReportQuery(opts))
	return c.getReport(ctx, path)
}

// GetAgentConversationMetrics returns per-agent conversation metrics.
// NOTE: trailing slash is intentional per the API spec.
func (c *Client) GetAgentConversationMetrics(ctx context.Context, opts ReportOpts) (any, error) {
	path := fmt.Sprintf("/api/v2/accounts/%d/reports/conversations/%s", c.accountID, c.buildReportQuery(opts))
	return c.getReport(ctx, path)
}

// GetFirstResponseTimeDistribution returns first response time distribution data.
func (c *Client) GetFirstResponseTimeDistribution(ctx context.Context, opts ReportOpts) (any, error) {
	path := fmt.Sprintf("/api/v2/accounts/%d/reports/first_response_time_distribution%s", c.accountID, c.buildReportQuery(opts))
	return c.getReport(ctx, path)
}

// GetInboxLabelMatrix returns the inbox/label matrix report.
func (c *Client) GetInboxLabelMatrix(ctx context.Context, opts ReportOpts) (any, error) {
	path := fmt.Sprintf("/api/v2/accounts/%d/reports/inbox_label_matrix%s", c.accountID, c.buildReportQuery(opts))
	return c.getReport(ctx, path)
}

// GetOutgoingMessagesCount returns outgoing message counts.
func (c *Client) GetOutgoingMessagesCount(ctx context.Context, opts ReportOpts) (any, error) {
	path := fmt.Sprintf("/api/v2/accounts/%d/reports/outgoing_messages_count%s", c.accountID, c.buildReportQuery(opts))
	return c.getReport(ctx, path)
}

// GetSummaryByAgent returns summary reports grouped by agent.
func (c *Client) GetSummaryByAgent(ctx context.Context, opts ReportOpts) (any, error) {
	path := fmt.Sprintf("/api/v2/accounts/%d/summary_reports/agent%s", c.accountID, c.buildReportQuery(opts))
	return c.getReport(ctx, path)
}

// GetSummaryByChannel returns summary reports grouped by channel.
func (c *Client) GetSummaryByChannel(ctx context.Context, opts ReportOpts) (any, error) {
	path := fmt.Sprintf("/api/v2/accounts/%d/summary_reports/channel%s", c.accountID, c.buildReportQuery(opts))
	return c.getReport(ctx, path)
}

// GetSummaryByInbox returns summary reports grouped by inbox.
func (c *Client) GetSummaryByInbox(ctx context.Context, opts ReportOpts) (any, error) {
	path := fmt.Sprintf("/api/v2/accounts/%d/summary_reports/inbox%s", c.accountID, c.buildReportQuery(opts))
	return c.getReport(ctx, path)
}

// GetSummaryByTeam returns summary reports grouped by team.
func (c *Client) GetSummaryByTeam(ctx context.Context, opts ReportOpts) (any, error) {
	path := fmt.Sprintf("/api/v2/accounts/%d/summary_reports/team%s", c.accountID, c.buildReportQuery(opts))
	return c.getReport(ctx, path)
}

// GetReportingEvents returns reporting events. Uses v1 API.
func (c *Client) GetReportingEvents(ctx context.Context, opts ReportOpts) (any, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/reporting_events%s", c.accountID, c.buildReportQuery(opts))
	return c.getReport(ctx, path)
}

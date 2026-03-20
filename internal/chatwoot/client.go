package chatwoot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Doer executes an HTTP request. Matches *http.Client.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client is the shared HTTP transport for all Chatwoot API calls.
type Client struct {
	baseURL        string
	token          string
	headerName     string
	http           Doer
	retryMax       int
	retryBaseDelay time.Duration
}

// NewClient creates a Client with the given base URL and auth credentials.
func NewClient(baseURL, token, headerName string) *Client {
	return &Client{
		baseURL:    baseURL,
		token:      token,
		headerName: headerName,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
		retryMax:       3,
		retryBaseDelay: time.Second,
	}
}

// NewClientWithDoer creates a Client with a custom Doer (for testing).
func NewClientWithDoer(baseURL, token, headerName string, doer Doer) *Client {
	return &Client{
		baseURL:        baseURL,
		token:          token,
		headerName:     headerName,
		http:           doer,
		retryMax:       3,
		retryBaseDelay: time.Second,
	}
}

// Do executes an HTTP request against the Chatwoot API.
func (c *Client) Do(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	url := c.baseURL + path

	var bodyReader *bytes.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	var req *http.Request
	var err error
	if bodyReader != nil {
		req, err = http.NewRequestWithContext(ctx, method, url, bodyReader)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	if c.token != "" && c.headerName != "" {
		req.Header.Set(c.headerName, c.token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.http.Do(req)
}

// DecodeResponse reads and decodes an HTTP response.
// For 2xx responses, it decodes the body into target.
// For non-2xx responses, it returns an *APIError.
func DecodeResponse(resp *http.Response, target any) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if target != nil && len(body) > 0 {
			if err := json.Unmarshal(body, target); err != nil {
				return fmt.Errorf("decode response: %w", err)
			}
		}
		return nil
	}

	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Code:       MapHTTPStatus(resp.StatusCode),
	}

	var errBody struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	if json.Unmarshal(body, &errBody) == nil {
		if errBody.Message != "" {
			apiErr.Message = errBody.Message
		} else if errBody.Error != "" {
			apiErr.Message = errBody.Error
		}
	}
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(resp.StatusCode)
	}

	return apiErr
}

// internal/chatwoot/retry.go
package chatwoot

import (
	"context"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

// DoWithRetry executes a request with retry on 429 and 5xx responses.
func (c *Client) DoWithRetry(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.retryMax; attempt++ {
		if attempt > 0 {
			delay := c.backoffDelay(attempt, resp)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, err = c.Do(ctx, method, path, body)
		if err != nil {
			return nil, err
		}

		if !shouldRetry(resp.StatusCode) {
			return resp, nil
		}

		// Close body before retry to prevent resource leak
		if attempt < c.retryMax {
			resp.Body.Close()
		}
	}

	return resp, nil
}

func shouldRetry(status int) bool {
	return status == 429 || status >= 500
}

func (c *Client) backoffDelay(attempt int, resp *http.Response) time.Duration {
	// Respect Retry-After header if present
	if resp != nil {
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, err := strconv.Atoi(ra); err == nil {
				return time.Duration(secs) * time.Second
			}
		}
	}

	// Exponential backoff with jitter
	base := c.retryBaseDelay
	for i := 1; i < attempt; i++ {
		base *= 2
	}
	// Add jitter: 0-25% of base
	if base > 0 {
		jitter := time.Duration(rand.Int64N(int64(base / 4)))
		return base + jitter
	}
	return 0
}

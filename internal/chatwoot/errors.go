package chatwoot

import (
	"fmt"

	"github.com/chatwoot/chatwoot-cli/internal/contract"
)

// APIError represents an error response from the Chatwoot API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Detail     any
}

func (e *APIError) Error() string {
	return fmt.Sprintf("chatwoot API error %d: %s \u2014 %s", e.StatusCode, e.Code, e.Message)
}

// MapHTTPStatus maps an HTTP status code to a contract error code.
func MapHTTPStatus(status int) string {
	switch {
	case status == 401:
		return contract.ErrCodeUnauthorized
	case status == 403:
		return contract.ErrCodeForbidden
	case status == 404:
		return contract.ErrCodeNotFound
	case status == 422:
		return contract.ErrCodeValidation
	case status == 429:
		return contract.ErrCodeRateLimited
	case status >= 500:
		return contract.ErrCodeServer
	default:
		return contract.ErrCodeServer
	}
}

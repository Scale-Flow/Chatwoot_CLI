package contract

// Standard error codes used across all commands.
const (
	ErrCodeUnauthorized = "unauthorized"
	ErrCodeForbidden    = "forbidden"
	ErrCodeNotFound     = "not_found"
	ErrCodeValidation   = "validation_error"
	ErrCodeRateLimited  = "rate_limited"
	ErrCodeServer       = "server_error"
	ErrCodeNetwork      = "network_error"
	ErrCodeConfig       = "config_error"
	ErrCodeAuth         = "auth_error"
)

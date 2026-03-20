// Package contract defines the JSON envelope types for all CLI stdout output.
package contract

// Response is the top-level envelope for all stdout output.
// Use the constructor functions (Success, SuccessList, Err, ErrWithDetail)
// to build responses — do not construct Response literals directly.
type Response struct {
	OK       bool         `json:"ok"`
	Data     any          `json:"data,omitempty"`
	Meta     *Meta        `json:"meta,omitempty"`
	Warnings []Warning    `json:"warnings,omitempty"`
	Error    *ErrorDetail `json:"error,omitempty"`
}

// Meta holds response metadata such as pagination.
type Meta struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination describes the pagination state for collection responses.
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

// ErrorDetail describes an error in the response envelope.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  any    `json:"detail,omitempty"`
}

// Warning describes a non-fatal warning attached to a success response.
type Warning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

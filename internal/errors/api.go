package errors

import (
	"fmt"
	"net/http"
)

// APIError represents a structured API error with an HTTP status code.
type APIError struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	RequestID  string `json:"-"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// ErrorResponse is the JSON envelope returned to clients.
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

// ErrorBody contains the user-facing error details.
type ErrorBody struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

// ToResponse converts an APIError into an ErrorResponse for serialization.
func (e *APIError) ToResponse() ErrorResponse {
	return ErrorResponse{
		Error: ErrorBody{
			Type:      e.Type,
			Message:   e.Message,
			RequestID: e.RequestID,
		},
	}
}

// NewUnauthorized returns a 401 APIError.
func NewUnauthorized(message string) *APIError {
	return &APIError{
		Type:       "UNAUTHORIZED",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewNotFound returns a 404 APIError.
func NewNotFound(message string) *APIError {
	return &APIError{
		Type:       "NOT_FOUND",
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

// NewBadRequest returns a 400 APIError.
func NewBadRequest(message string) *APIError {
	return &APIError{
		Type:       "BAD_REQUEST",
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewConflict returns a 409 APIError.
func NewConflict(message string) *APIError {
	return &APIError{
		Type:       "CONFLICT",
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// NewTooManyRequests returns a 429 APIError with a Retry-After hint.
func NewTooManyRequests(retryAfter string) *APIError {
	return &APIError{
		Type:       "TOO_MANY_REQUESTS",
		Message:    fmt.Sprintf("rate limit exceeded, retry after %s", retryAfter),
		StatusCode: http.StatusTooManyRequests,
	}
}

// NewInternal returns a 500 APIError.
func NewInternal(message string) *APIError {
	return &APIError{
		Type:       "INTERNAL_SERVER_ERROR",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

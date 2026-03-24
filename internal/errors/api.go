package errors

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	RequestID  string `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func (e *APIError) ToResponse() ErrorResponse {
	return ErrorResponse{
		Error: ErrorBody{
			Type:      e.Type,
			Message:   e.Message,
			RequestID: e.RequestID,
		},
	}
}

func NewUnauthorized(message string) *APIError {
	return &APIError{
		Type:       "UNAUTHORIZED",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewNotFound(message string) *APIError {
	return &APIError{
		Type:       "NOT_FOUND",
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func NewBadRequest(message string) *APIError {
	return &APIError{
		Type:       "BAD_REQUEST",
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func NewConflict(message string) *APIError {
	return &APIError{
		Type:       "CONFLICT",
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func NewTooManyRequests(retryAfter string) *APIError {
	return &APIError{
		Type:       "TOO_MANY_REQUESTS",
		Message:    fmt.Sprintf("rate limit exceeded, retry after %s", retryAfter),
		StatusCode: http.StatusTooManyRequests,
	}
}

func NewInternal(message string) *APIError {
	return &APIError{
		Type:       "INTERNAL_SERVER_ERROR",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

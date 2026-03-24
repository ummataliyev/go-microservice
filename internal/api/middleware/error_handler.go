package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	domainerrors "github.com/thealish/go-microservice/internal/domain/errors"
)

// ErrorHandler is a Fiber custom error handler that converts errors into structured
// JSON error responses.
func ErrorHandler(c *fiber.Ctx, err error) error {
	requestID, _ := c.Locals("request_id").(string)

	// Check for domain APIError.
	if apiErr, ok := err.(*domainerrors.APIError); ok {
		apiErr.RequestID = requestID
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	// Check for Fiber error.
	if fiberErr, ok := err.(*fiber.Error); ok {
		apiErr := mapFiberError(fiberErr)
		apiErr.RequestID = requestID
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	// Default: 500 Internal Server Error.
	apiErr := &domainerrors.APIError{
		Type:       "INTERNAL_ERROR",
		Message:    "internal server error",
		StatusCode: http.StatusInternalServerError,
		RequestID:  requestID,
	}
	return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
}

func mapFiberError(fiberErr *fiber.Error) *domainerrors.APIError {
	switch fiberErr.Code {
	case http.StatusNotFound:
		return &domainerrors.APIError{
			Type:       "NOT_FOUND",
			Message:    fiberErr.Message,
			StatusCode: http.StatusNotFound,
		}
	case http.StatusBadRequest:
		return &domainerrors.APIError{
			Type:       "BAD_REQUEST",
			Message:    fiberErr.Message,
			StatusCode: http.StatusBadRequest,
		}
	case http.StatusUnauthorized:
		return &domainerrors.APIError{
			Type:       "UNAUTHORIZED",
			Message:    fiberErr.Message,
			StatusCode: http.StatusUnauthorized,
		}
	case http.StatusForbidden:
		return &domainerrors.APIError{
			Type:       "FORBIDDEN",
			Message:    fiberErr.Message,
			StatusCode: http.StatusForbidden,
		}
	case http.StatusConflict:
		return &domainerrors.APIError{
			Type:       "CONFLICT",
			Message:    fiberErr.Message,
			StatusCode: http.StatusConflict,
		}
	case http.StatusTooManyRequests:
		return &domainerrors.APIError{
			Type:       "TOO_MANY_REQUESTS",
			Message:    fiberErr.Message,
			StatusCode: http.StatusTooManyRequests,
		}
	default:
		return &domainerrors.APIError{
			Type:       "INTERNAL_ERROR",
			Message:    fiberErr.Message,
			StatusCode: fiberErr.Code,
		}
	}
}

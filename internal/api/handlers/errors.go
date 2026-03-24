package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	svcerrors "go-microservice/internal/errors"
)

// handleServiceError is the shared error mapper used by all handlers.
func handleServiceError(c *fiber.Ctx, err error) error {
	var apiErr *svcerrors.APIError

	switch {
	case errors.Is(err, svcerrors.ErrUserAlreadyExists):
		apiErr = svcerrors.NewConflict(err.Error())
	case errors.Is(err, svcerrors.ErrInvalidCredentials):
		apiErr = svcerrors.NewUnauthorized(err.Error())
	case errors.Is(err, svcerrors.ErrLoginLocked):
		apiErr = svcerrors.NewTooManyRequests("30m")
	case errors.Is(err, svcerrors.ErrInvalidToken):
		apiErr = svcerrors.NewUnauthorized(err.Error())
	case errors.Is(err, svcerrors.ErrInvalidTokenType):
		apiErr = svcerrors.NewUnauthorized(err.Error())
	case errors.Is(err, svcerrors.ErrUserNotFound):
		apiErr = svcerrors.NewNotFound(err.Error())
	default:
		apiErr = svcerrors.NewInternal("internal server error")
	}

	return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
}

package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	apierrors "go-microservice/internal/errors"
)

var validate = validator.New()

func validateBody(c *fiber.Ctx, out interface{}) error {
	if err := c.BodyParser(out); err != nil {
		return apierrors.NewBadRequest("invalid request body")
	}
	if err := validate.Struct(out); err != nil {
		return apierrors.NewBadRequest(formatValidationError(err))
	}
	return nil
}

func formatValidationError(err error) string {
	if ve, ok := err.(validator.ValidationErrors); ok {
		// Return first validation error in a clean format.
		fe := ve[0]
		switch fe.Tag() {
		case "required":
			return fe.Field() + " is required"
		case "email":
			return fe.Field() + " must be a valid email"
		case "min":
			return fe.Field() + " must be at least " + fe.Param() + " characters"
		default:
			return fe.Field() + " is invalid"
		}
	}
	return "validation failed"
}

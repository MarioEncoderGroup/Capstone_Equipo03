package middleware

import (
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
	"github.com/gofiber/fiber/v2"
)

// ValidatorMiddleware adds the validator instance to the Fiber context
func ValidatorMiddleware(validator *validatorapi.XValidator) fiber.Handler {
	// Register custom validations for MisViaticos
	validator.RegisterCustomValidations()
	
	return func(c *fiber.Ctx) error {
		c.Locals("validator", validator)
		return c.Next()
	}
}
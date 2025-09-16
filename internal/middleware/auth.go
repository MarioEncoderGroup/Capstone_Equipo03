package middleware

import (
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT tokens and user authentication for MisViaticos
func AuthMiddleware(dbControl *postgresql.PostgresqlClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// For now, return a placeholder response
		// This will be fully implemented when we create the auth module
		
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Authentication required",
				"message": "Authorization header is missing",
				"code":    "AUTH_REQUIRED",
			})
		}

		// TODO: Implement full JWT validation
		// 1. Extract token from "Bearer <token>" format
		// 2. Validate JWT signature
		// 3. Check token expiration
		// 4. Get user from database
		// 5. Check if user is active
		// 6. Add user info to context

		// Placeholder for development - remove in production
		c.Locals("user", fiber.Map{
			"id":       "placeholder-user-id",
			"email":    "dev@misviaticos.cl",
			"tenantId": "1",
		})

		return c.Next()
	}
}
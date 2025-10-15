package middleware

import (
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/gofiber/fiber/v2"
)

// DatabaseControlMiddleware adds the control database connection to the context
func DatabaseControlMiddleware(db *postgresql.PostgresqlClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("db_control", db)
		return c.Next()
	}
}

// DatabaseTenant1Middleware adds the default tenant database connection to the context
func DatabaseTenant1Middleware(db *postgresql.PostgresqlClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("db_tenant", db)
		return c.Next()
	}
}

// DatabaseTenantMiddleware adds a specific tenant database connection to the context
// This would be used when the tenant is determined from the request (JWT, header, etc.)
func DatabaseTenantMiddleware(getTenantDB func(tenantID string) *postgresql.PostgresqlClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get tenant ID from JWT token or header
		tenantID := c.Get("X-Tenant-ID")
		if tenantID == "" {
			// Try to get from JWT claims (would be implemented in auth middleware)
			if claims := c.Locals("user"); claims != nil {
				// Extract tenant ID from claims
				// tenantID = claims.(JWTClaims).TenantID
			}
		}

		if tenantID != "" {
			if tenantDB := getTenantDB(tenantID); tenantDB != nil {
				c.Locals("db_tenant", tenantDB)
			}
		}

		return c.Next()
	}
}
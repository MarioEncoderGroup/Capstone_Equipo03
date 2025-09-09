package routes

import (
	"github.com/JoseLuis21/mv-backend/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes defines all public routes for MisViaticos API
func PublicRoutes(app *fiber.App) *fiber.App {
	// Welcome endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to MisViaticos API",
			"version": "1.0.0",
			"docs":    "/api/v1/docs",
		})
	})

	// API v1 public routes
	public := app.Group("/api/v1")

	// Authentication routes
	auth := public.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Get("/verify-email/:token", controllers.VerifyUserEmail)
	auth.Post("/forgot-password", controllers.ForgotPassword)
	auth.Post("/reset-password", controllers.ResetPassword)
	auth.Post("/refresh-token", controllers.RefreshToken)

	// Public information routes
	info := public.Group("/info")
	info.Get("/currencies", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"currencies": []fiber.Map{
				{"code": "CLP", "name": "Chilean Peso", "symbol": "$"},
				{"code": "USD", "name": "US Dollar", "symbol": "$"},
				{"code": "EUR", "name": "Euro", "symbol": "€"},
			},
		})
	})

	// Health check
	public.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"service":   "misviaticos-api",
			"timestamp": "2025-01-01T00:00:00Z",
		})
	})

	return app
}
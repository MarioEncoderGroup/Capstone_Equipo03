package routes

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/JoseLuis21/mv-backend/internal/controllers"
)

// PublicRoutes defines all public routes for MisViaticos API
func PublicRoutes(app *fiber.App, regionController *controllers.RegionController, communeController *controllers.CommuneController) *fiber.App {
	// Welcome endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to MisViaticos API",
			"version": "1.0.0",
		})
	})

	// API v1 public routes
	public := app.Group("/api/v1")

	// Authentication routes - DEPRECATED
	// These routes are now handled by AuthRoutes() function
	// Legacy routes commented out as they're now managed by dedicated AuthRoutes
	/*
	auth := public.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Get("/verify-email/:token", controllers.VerifyUserEmail)
	auth.Post("/forgot-password", controllers.ForgotPassword)
	auth.Post("/reset-password", controllers.ResetPassword)
	auth.Post("/refresh-token", controllers.RefreshToken)
	*/

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
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "1.0.0",
		})
	})

	// Regiones y Comunas - Datos públicos para formularios
	public.Get("/regions", regionController.GetAllRegions)
	public.Get("/regions/:id", regionController.GetRegionByID)
	public.Get("/communes", communeController.GetCommunes)
	public.Get("/communes/:id", communeController.GetCommuneByID)

	return app
}
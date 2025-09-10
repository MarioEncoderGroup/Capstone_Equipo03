package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/JoseLuis21/mv-backend/internal/controllers"
)

// AuthRoutes configura las rutas de autenticación usando el AuthController
func AuthRoutes(app *fiber.App, authController *controllers.AuthController) {
	// API v1 public routes
	public := app.Group("/api/v1")

	// Authentication routes group
	auth := public.Group("/auth")

	// Register endpoints - Usa el AuthController real
	auth.Post("/register", authController.Register)
	auth.Post("/verify-email", authController.VerifyUserEmail)
	auth.Post("/resend-verification", authController.ResendEmailVerification)
	auth.Get("/health", authController.HealthCheck)

	// Legacy endpoints - mantener por compatibilidad mientras migramos
	auth.Post("/login", controllers.Login)
	auth.Post("/forgot-password", controllers.ForgotPassword)
	auth.Post("/reset-password", controllers.ResetPassword)
	auth.Post("/refresh-token", controllers.RefreshToken)

	// Endpoint alternativo para verificación via GET (links de email)
	auth.Get("/verify-email/:token", func(c *fiber.Ctx) error {
		token := c.Params("token")
		if token == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Token requerido",
				"error":   "Token no proporcionado en la URL",
			})
		}

		// Crear request body para el método POST
		c.Request().SetBody([]byte(`{"token":"` + token + `"}`))
		c.Request().Header.SetContentType("application/json")

		// Llamar al método POST del controller
		return authController.VerifyUserEmail(c)
	})
}
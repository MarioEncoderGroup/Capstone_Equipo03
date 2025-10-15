package routes

import (
	"github.com/JoseLuis21/mv-backend/internal/controllers"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/JoseLuis21/mv-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// OCRRoutes configura las rutas de OCR usando el OCRController
func OCRRoutes(app *fiber.App, dbControl *postgresql.PostgresqlClient, ocrController *controllers.OCRController) {
	// Si el OCR controller no está configurado (faltan credenciales), no crear rutas
	if ocrController == nil {
		return
	}

	// Crear middleware de autenticación
	authMiddleware := middleware.AuthMiddleware(dbControl)

	// API v1 private routes
	private := app.Group("/api/v1", authMiddleware)

	// OCR routes group
	ocr := private.Group("/ocr")

	// POST /api/v1/ocr/analyze - Analizar comprobante con OCR
	ocr.Post("/analyze", ocrController.AnalyzeReceipt)
}

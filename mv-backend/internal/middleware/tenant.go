package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
)

// RequireTenantMiddleware valida que el usuario tenga un tenant seleccionado
// Este middleware es crítico para la seguridad multi-tenant
func RequireTenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Locals("tenantID")
		
		// Validar que el tenant esté presente
		if tenantID == nil || tenantID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
				Success: false,
				Message: "Debe seleccionar un tenant para acceder a este recurso",
				Error:   "TENANT_NOT_SELECTED",
			})
		}

		// Validar que sea un string válido
		if str, ok := tenantID.(string); !ok || str == "" {
			return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
				Success: false,
				Message: "Tenant ID inválido",
				Error:   "INVALID_TENANT_ID",
			})
		}

		// Continuar con el siguiente handler
		return c.Next()
	}
}

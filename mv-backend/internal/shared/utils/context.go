package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetUserIDFromContext extrae el userID del contexto de Fiber
// Este userID es establecido por el AuthMiddleware después de validar el JWT
func GetUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userIDInterface := c.Locals("userID")
	if userIDInterface == nil {
		return uuid.Nil, fmt.Errorf("user ID not found in context - authentication required")
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user ID format in context")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID UUID format: %w", err)
	}

	return userID, nil
}

// GetTenantIDFromContext extrae el tenantID del contexto de Fiber (opcional)
// Retorna nil si no hay tenant seleccionado
func GetTenantIDFromContext(c *fiber.Ctx) (*uuid.UUID, error) {
	tenantIDInterface := c.Locals("tenantID")
	if tenantIDInterface == nil {
		// No hay tenant seleccionado, esto es válido
		return nil, nil
	}

	tenantIDStr, ok := tenantIDInterface.(string)
	if !ok {
		return nil, fmt.Errorf("invalid tenant ID format in context")
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID UUID format: %w", err)
	}

	return &tenantID, nil
}

// RequireTenantID extrae el tenantID y retorna error si no está presente
func RequireTenantID(c *fiber.Ctx) (uuid.UUID, error) {
	tenantIDPtr, err := GetTenantIDFromContext(c)
	if err != nil {
		return uuid.Nil, err
	}

	if tenantIDPtr == nil {
		return uuid.Nil, fmt.Errorf("tenant selection required - no tenant in context")
	}

	return *tenantIDPtr, nil
}

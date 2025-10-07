package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
)

// User representa la información del usuario autenticado
type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	TenantID *uuid.UUID `json:"tenant_id,omitempty"`
}

// GetUserFromContext extrae la información del usuario desde el contexto de Fiber
func GetUserFromContext(c *fiber.Ctx) (*User, error) {
	userLocal := c.Locals("user")
	if userLocal == nil {
		return nil, sharedErrors.NewAuthError("Usuario no encontrado en el contexto", "USER_NOT_FOUND")
	}

	// Si es un mapa (como el placeholder del middleware de desarrollo)
	if userMap, ok := userLocal.(fiber.Map); ok {
		user := &User{}

		// Extraer ID
		if idStr, ok := userMap["id"].(string); ok {
			if idStr != "placeholder-user-id" {
				id, err := uuid.Parse(idStr)
				if err != nil {
					return nil, fmt.Errorf("ID de usuario inválido: %w", err)
				}
				user.ID = id
			} else {
				// Para desarrollo, usar un UUID fijo
				user.ID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
			}
		}

		// Extraer email
		if email, ok := userMap["email"].(string); ok {
			user.Email = email
		}

		// Extraer full name si existe
		if fullName, ok := userMap["full_name"].(string); ok {
			user.FullName = fullName
		} else if name, ok := userMap["name"].(string); ok {
			user.FullName = name
		}

		// Extraer tenant ID si existe
		if tenantIDStr, ok := userMap["tenantId"].(string); ok {
			if tenantIDStr != "" && tenantIDStr != "1" {
				tenantID, err := uuid.Parse(tenantIDStr)
				if err == nil {
					user.TenantID = &tenantID
				}
			}
		}

		return user, nil
	}

	// Si es directamente una estructura User
	if user, ok := userLocal.(*User); ok {
		return user, nil
	}

	return nil, sharedErrors.NewAuthError("Formato de usuario inválido en el contexto", "INVALID_USER_FORMAT")
}

// GetTenantIDFromContext extrae el tenant ID desde el contexto o headers
func GetTenantIDFromContext(c *fiber.Ctx) (*uuid.UUID, error) {
	// Intentar obtener del usuario autenticado primero
	user, err := GetUserFromContext(c)
	if err == nil && user.TenantID != nil {
		return user.TenantID, nil
	}

	// Intentar obtener del header X-Tenant-ID
	if tenantIDStr := c.Get("X-Tenant-ID"); tenantIDStr != "" {
		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			return nil, fmt.Errorf("tenant ID inválido en header: %w", err)
		}
		return &tenantID, nil
	}

	// Intentar obtener de parámetros de ruta
	if tenantIDParam := c.Params("tenantId"); tenantIDParam != "" {
		tenantID, err := uuid.Parse(tenantIDParam)
		if err != nil {
			return nil, fmt.Errorf("tenant ID inválido en parámetro: %w", err)
		}
		return &tenantID, nil
	}

	return nil, nil // No hay tenant ID especificado
}

// HandleAuthError maneja errores de autenticación de manera consistente
func HandleAuthError(c *fiber.Ctx, err error) error {
	if authErr, ok := err.(*sharedErrors.AuthError); ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": authErr.Message,
			"error":   authErr.Code,
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"success": false,
		"message": "Error de autenticación",
		"error":   err.Error(),
	})
}

// SetUserInContext establece la información del usuario en el contexto
func SetUserInContext(c *fiber.Ctx, user *User) {
	c.Locals("user", user)
}

// RequireAuthentication es un helper que verifica que hay un usuario autenticado
func RequireAuthentication(c *fiber.Ctx) (*User, error) {
	user, err := GetUserFromContext(c)
	if err != nil {
		return nil, sharedErrors.NewAuthError("Autenticación requerida", "AUTHENTICATION_REQUIRED")
	}
	return user, nil
}
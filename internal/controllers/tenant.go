package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	authPorts "github.com/JoseLuis21/mv-backend/internal/core/auth/ports"
	tenantDomain "github.com/JoseLuis21/mv-backend/internal/core/tenant/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/tenant/ports"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
	"github.com/google/uuid"
)

// TenantController maneja las operaciones de tenant
type TenantController struct {
	tenantService ports.TenantService
	authService   authPorts.AuthService
	validator     *validatorapi.Validator
}

// NewTenantController crea una nueva instancia del controller de tenant
func NewTenantController(tenantService ports.TenantService, authService authPorts.AuthService, validator *validatorapi.Validator) *TenantController {
	return &TenantController{
		tenantService: tenantService,
		authService:   authService,
		validator:     validator,
	}
}


// GetTenantsByUser obtiene todos los tenants del usuario autenticado
func (tc *TenantController) GetTenantsByUser(c *fiber.Ctx) error {
	// 1. Obtener userID del contexto (middleware de autenticación)
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
			Success: false,
			Message: "Usuario no autenticado",
			Error:   "UNAUTHORIZED",
		})
	}

	// 2. Usar contexto de Fiber
	// 3. Obtener tenants del usuario
	tenants, err := tc.tenantService.GetTenantsByUser(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo tenants",
			Error:   err.Error(),
		})
	}

	// 4. Respuesta exitosa
	return c.JSON(types.APIResponse{
		Success: true,
		Message: "Tenants obtenidos exitosamente",
		Data:    tenants,
	})
}

// SelectTenant selecciona un tenant específico para el usuario
func (tc *TenantController) SelectTenant(c *fiber.Ctx) error {
	// 1. Obtener userID del contexto
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
			Success: false,
			Message: "Usuario no autenticado",
			Error:   "UNAUTHORIZED",
		})
	}

	// 2. Obtener tenantID de los parámetros
	tenantIDStr := c.Params("tenantId")
	if tenantIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID del tenant requerido",
			Error:   "TENANT_ID_REQUIRED",
		})
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID del tenant inválido",
			Error:   "INVALID_TENANT_ID",
		})
	}

	// 3. Usar contexto de Fiber
	// 4. Obtener información del tenant desde TenantService
	tenant, err := tc.tenantService.GetTenantProfile(c.Context(), tenantID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(types.APIResponse{
			Success: false,
			Message: "Tenant no encontrado",
			Error:   err.Error(),
		})
	}

	// 5. Seleccionar tenant usando AuthService (genera nuevos tokens con tenant_id)
	response, err := tc.authService.SelectTenant(c.Context(), tenant, userID)
	if err != nil {
		if appErr, ok := sharedErrors.IsAppError(err); ok {
			return c.Status(appErr.HTTPCode).JSON(types.APIResponse{
				Success: false,
				Message: appErr.Message,
				Error:   appErr.Code,
			})
		}
		
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error seleccionando tenant",
			Error:   err.Error(),
		})
	}

	// 6. Respuesta exitosa
	return c.JSON(types.APIResponse{
		Success: true,
		Message: "Tenant seleccionado exitosamente",
		Data:    response,
	})
}

// GetTenantProfile obtiene el perfil del tenant seleccionado
func (tc *TenantController) GetTenantProfile(c *fiber.Ctx) error {
	// 1. Obtener tenantID de los parámetros
	tenantIDStr := c.Params("tenantId")
	if tenantIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID del tenant requerido",
			Error:   "TENANT_ID_REQUIRED",
		})
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID del tenant inválido",
			Error:   "INVALID_TENANT_ID",
		})
	}

	// 2. Usar contexto de Fiber
	// 3. Obtener perfil del tenant
	tenant, err := tc.tenantService.GetTenantProfile(c.Context(), tenantID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(types.APIResponse{
			Success: false,
			Message: "Tenant no encontrado",
			Error:   err.Error(),
		})
	}

	// 4. Respuesta exitosa
	return c.JSON(types.APIResponse{
		Success: true,
		Message: "Perfil del tenant obtenido exitosamente",
		Data:    tenant,
	})
}

// UpdateTenantProfile actualiza el perfil del tenant
func (tc *TenantController) UpdateTenantProfile(c *fiber.Ctx) error {
	// 1. Obtener tenantID de los parámetros
	tenantIDStr := c.Params("tenantId")
	if tenantIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID del tenant requerido",
			Error:   "TENANT_ID_REQUIRED",
		})
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID del tenant inválido",
			Error:   "INVALID_TENANT_ID",
		})
	}

	// 2. Parsear datos de actualización
	var updateData tenantDomain.Tenant
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error parseando datos",
			Error:   "Formato de datos inválido",
		})
	}

	// 3. Validar datos
	if errors := tc.validator.ValidateStruct(updateData); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Errores de validación",
			Data:    errors,
		})
	}

	// 4. Asegurar que el ID coincida
	updateData.ID = tenantID

	// 5. Usar contexto de Fiber
	// 6. Actualizar tenant
	if err := tc.tenantService.UpdateTenant(c.Context(), &updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error actualizando tenant",
			Error:   err.Error(),
		})
	}

	// 7. Respuesta exitosa
	return c.JSON(types.APIResponse{
		Success: true,
		Message: "Perfil del tenant actualizado exitosamente",
	})
}

// HealthCheck endpoint de health check para el módulo de tenant
func (tc *TenantController) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(types.APIResponse{
		Success: true,
		Message: "Tenant module is healthy",
		Data: fiber.Map{
			"timestamp": time.Now().Format(time.RFC3339),
			"module":    "tenant",
			"version":   "1.0.0",
		},
	})
}

// getUserIDFromContext obtiene el ID del usuario desde el contexto (helper function)
// En un entorno real, esto vendría del middleware de autenticación
func getUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	// Por ahora retornamos un UUID mock - esto debe ser implementado con el middleware de JWT
	// userIDStr := c.Locals("userID")
	// if userIDStr == nil {
	//     return uuid.Nil, fmt.Errorf("user ID not found in context")
	// }
	// return uuid.Parse(userIDStr.(string))
	
	// Placeholder para desarrollo - debe ser reemplazado por la lógica real
	return uuid.New(), nil
}
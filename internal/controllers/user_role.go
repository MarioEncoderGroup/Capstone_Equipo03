package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	userRoleDomain "github.com/JoseLuis21/mv-backend/internal/core/user_role/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/user_role/ports"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
)

// UserRoleController maneja las operaciones de relaciones usuario-rol
type UserRoleController struct {
	userRoleService ports.UserRoleService
	validator       *validatorapi.Validator
}

// NewUserRoleController crea una nueva instancia del controller de user_role
func NewUserRoleController(userRoleService ports.UserRoleService, validator *validatorapi.Validator) *UserRoleController {
	return &UserRoleController{
		userRoleService: userRoleService,
		validator:       validator,
	}
}

// CreateUserRole maneja POST /user-roles - Asigna un rol a un usuario
func (urc *UserRoleController) CreateUserRole(c *fiber.Ctx) error {
	// Parsear request body
	var req userRoleDomain.CreateUserRoleDto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error parseando datos de entrada",
			Error:   "Formato de datos inválido",
		})
	}

	// Extraer TenantID del contexto si está disponible
	if tenantIDStr := c.Locals("tenantId"); tenantIDStr != nil && tenantIDStr != "" {
		tenantID, err := uuid.Parse(tenantIDStr.(string))
		if err == nil {
			req.TenantID = &tenantID
		}
	}

	// Validar estructura de datos
	if errors := urc.validator.ValidateStruct(req); len(errors) > 0 {
		var validationErrors []types.ValidationErrorResponse
		for _, err := range errors {
			validationErrors = append(validationErrors, types.ValidationErrorResponse{
				Field:   err.Field,
				Message: err.Message,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Data:    validationErrors,
		})
	}

	// Crear relación usuario-rol a través del servicio
	result, err := urc.userRoleService.CreateUserRole(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error asignando rol al usuario",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(types.APIResponse{
		Success: true,
		Message: result.Message,
		Data:    result.UserRole,
	})
}

// DeleteUserRole maneja DELETE /user-roles/:userId/:roleId - Elimina un rol de un usuario
func (urc *UserRoleController) DeleteUserRole(c *fiber.Ctx) error {
	// Parsear IDs
	userIDStr := c.Params("userId")
	roleIDStr := c.Params("roleId")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de usuario inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de rol inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Extraer TenantID del contexto si está disponible
	var tenantID *uuid.UUID
	if tenantIDStr := c.Locals("tenantId"); tenantIDStr != nil && tenantIDStr != "" {
		tid, err := uuid.Parse(tenantIDStr.(string))
		if err == nil {
			tenantID = &tid
		}
	}

	// Eliminar relación usuario-rol a través del servicio
	if err := urc.userRoleService.DeleteUserRole(c.Context(), userID, roleID, tenantID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error eliminando rol del usuario",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Rol eliminado del usuario exitosamente",
	})
}

// GetUserRoles maneja GET /user-roles/user/:userId - Obtiene todos los roles de un usuario
func (urc *UserRoleController) GetUserRoles(c *fiber.Ctx) error {
	// Parsear ID del usuario
	userIDStr := c.Params("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de usuario inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Extraer TenantID del contexto si está disponible
	var tenantID *uuid.UUID
	if tenantIDStr := c.Locals("tenantId"); tenantIDStr != nil && tenantIDStr != "" {
		tid, err := uuid.Parse(tenantIDStr.(string))
		if err == nil {
			tenantID = &tid
		}
	}

	// Obtener roles del usuario desde el servicio
	userRoles, err := urc.userRoleService.GetUserRoles(c.Context(), userID, tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo roles del usuario",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Roles del usuario obtenidos exitosamente",
		Data:    userRoles,
	})
}

// GetRoleUsers maneja GET /user-roles/role/:roleId - Obtiene todos los usuarios de un rol
func (urc *UserRoleController) GetRoleUsers(c *fiber.Ctx) error {
	// Parsear ID del rol
	roleIDStr := c.Params("roleId")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de rol inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Extraer TenantID del contexto si está disponible
	var tenantID *uuid.UUID
	if tenantIDStr := c.Locals("tenantId"); tenantIDStr != nil && tenantIDStr != "" {
		tid, err := uuid.Parse(tenantIDStr.(string))
		if err == nil {
			tenantID = &tid
		}
	}

	// Obtener usuarios del rol desde el servicio
	roleUsers, err := urc.userRoleService.GetRoleUsers(c.Context(), roleID, tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo usuarios del rol",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Usuarios del rol obtenidos exitosamente",
		Data:    roleUsers,
	})
}

// SyncUserRoles maneja PUT /user-roles/sync/user/:userId - Sincroniza roles de un usuario
func (urc *UserRoleController) SyncUserRoles(c *fiber.Ctx) error {
	// Parsear ID del usuario
	userIDStr := c.Params("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de usuario inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Parsear request body
	var reqBody struct {
		RoleIDs []string `json:"role_ids" validate:"required"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error parseando datos de entrada",
			Error:   "Formato de datos inválido",
		})
	}

	// Convertir strings a UUIDs
	var roleIDs []uuid.UUID
	for _, roleIDStr := range reqBody.RoleIDs {
		roleID, err := uuid.Parse(roleIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
				Success: false,
				Message: "ID de rol inválido",
				Error:   "Todos los IDs deben ser UUIDs válidos",
			})
		}
		roleIDs = append(roleIDs, roleID)
	}

	// Extraer TenantID del contexto si está disponible
	var tenantID *uuid.UUID
	if tenantIDStr := c.Locals("tenantId"); tenantIDStr != nil && tenantIDStr != "" {
		tid, err := uuid.Parse(tenantIDStr.(string))
		if err == nil {
			tenantID = &tid
		}
	}

	// Crear DTO de sincronización
	req := &userRoleDomain.SyncUserRolesDto{
		UserID:   userID,
		RoleIDs:  roleIDs,
		TenantID: tenantID,
	}

	// Validar estructura de datos
	if errors := urc.validator.ValidateStruct(req); len(errors) > 0 {
		var validationErrors []types.ValidationErrorResponse
		for _, err := range errors {
			validationErrors = append(validationErrors, types.ValidationErrorResponse{
				Field:   err.Field,
				Message: err.Message,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Data:    validationErrors,
		})
	}

	// Sincronizar roles del usuario a través del servicio
	if err := urc.userRoleService.SyncUserRoles(c.Context(), req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error sincronizando roles del usuario",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Roles del usuario sincronizados exitosamente",
	})
}

// SyncRoleUsers maneja PUT /user-roles/sync/role/:roleId - Sincroniza usuarios de un rol
func (urc *UserRoleController) SyncRoleUsers(c *fiber.Ctx) error {
	// Parsear ID del rol
	roleIDStr := c.Params("roleId")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de rol inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Parsear request body
	var reqBody struct {
		UserIDs []string `json:"user_ids" validate:"required"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error parseando datos de entrada",
			Error:   "Formato de datos inválido",
		})
	}

	// Convertir strings a UUIDs
	var userIDs []uuid.UUID
	for _, userIDStr := range reqBody.UserIDs {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
				Success: false,
				Message: "ID de usuario inválido",
				Error:   "Todos los IDs deben ser UUIDs válidos",
			})
		}
		userIDs = append(userIDs, userID)
	}

	// Extraer TenantID del contexto si está disponible
	var tenantID *uuid.UUID
	if tenantIDStr := c.Locals("tenantId"); tenantIDStr != nil && tenantIDStr != "" {
		tid, err := uuid.Parse(tenantIDStr.(string))
		if err == nil {
			tenantID = &tid
		}
	}

	// Crear DTO de sincronización
	req := &userRoleDomain.SyncRoleUsersDto{
		RoleID:   roleID,
		UserIDs:  userIDs,
		TenantID: tenantID,
	}

	// Validar estructura de datos
	if errors := urc.validator.ValidateStruct(req); len(errors) > 0 {
		var validationErrors []types.ValidationErrorResponse
		for _, err := range errors {
			validationErrors = append(validationErrors, types.ValidationErrorResponse{
				Field:   err.Field,
				Message: err.Message,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Data:    validationErrors,
		})
	}

	// Sincronizar usuarios del rol a través del servicio
	if err := urc.userRoleService.SyncRoleUsers(c.Context(), req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error sincronizando usuarios del rol",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Usuarios del rol sincronizados exitosamente",
	})
}

// CheckUserHasRole maneja GET /user-roles/check/:userId/:roleId - Verifica si un usuario tiene un rol
func (urc *UserRoleController) CheckUserHasRole(c *fiber.Ctx) error {
	// Parsear IDs
	userIDStr := c.Params("userId")
	roleIDStr := c.Params("roleId")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de usuario inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de rol inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Extraer TenantID del contexto si está disponible
	var tenantID *uuid.UUID
	if tenantIDStr := c.Locals("tenantId"); tenantIDStr != nil && tenantIDStr != "" {
		tid, err := uuid.Parse(tenantIDStr.(string))
		if err == nil {
			tenantID = &tid
		}
	}

	// Verificar si el usuario tiene el rol
	hasRole, err := urc.userRoleService.UserHasRole(c.Context(), userID, roleID, tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error verificando rol del usuario",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Verificación de rol completada",
		Data: fiber.Map{
			"user_id":  userID,
			"role_id":  roleID,
			"has_role": hasRole,
		},
	})
}

// RemoveUserFromAllRoles maneja DELETE /user-roles/user/:userId/all - Elimina un usuario de todos sus roles
func (urc *UserRoleController) RemoveUserFromAllRoles(c *fiber.Ctx) error {
	// Parsear ID del usuario
	userIDStr := c.Params("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de usuario inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Extraer TenantID del contexto si está disponible
	var tenantID *uuid.UUID
	if tenantIDStr := c.Locals("tenantId"); tenantIDStr != nil && tenantIDStr != "" {
		tid, err := uuid.Parse(tenantIDStr.(string))
		if err == nil {
			tenantID = &tid
		}
	}

	// Eliminar usuario de todos los roles a través del servicio
	if err := urc.userRoleService.RemoveUserFromAllRoles(c.Context(), userID, tenantID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error eliminando usuario de todos los roles",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Usuario eliminado de todos los roles exitosamente",
	})
}

// RemoveAllUsersFromRole maneja DELETE /user-roles/role/:roleId/all - Elimina todos los usuarios de un rol
func (urc *UserRoleController) RemoveAllUsersFromRole(c *fiber.Ctx) error {
	// Parsear ID del rol
	roleIDStr := c.Params("roleId")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de rol inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Extraer TenantID del contexto si está disponible
	var tenantID *uuid.UUID
	if tenantIDStr := c.Locals("tenantId"); tenantIDStr != nil && tenantIDStr != "" {
		tid, err := uuid.Parse(tenantIDStr.(string))
		if err == nil {
			tenantID = &tid
		}
	}

	// Eliminar todos los usuarios del rol a través del servicio
	if err := urc.userRoleService.RemoveAllUsersFromRole(c.Context(), roleID, tenantID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error eliminando todos los usuarios del rol",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Todos los usuarios eliminados del rol exitosamente",
	})
}
package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	rolePermissionDomain "github.com/JoseLuis21/mv-backend/internal/core/role_permission/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/role_permission/ports"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
)

// RolePermissionController maneja las operaciones de relaciones rol-permiso
type RolePermissionController struct {
	rolePermissionService ports.RolePermissionService
	validator             *validatorapi.Validator
}

// NewRolePermissionController crea una nueva instancia del controller de role_permission
func NewRolePermissionController(rolePermissionService ports.RolePermissionService, validator *validatorapi.Validator) *RolePermissionController {
	return &RolePermissionController{
		rolePermissionService: rolePermissionService,
		validator:             validator,
	}
}

// CreateRolePermission maneja POST /role-permissions/assign - Asigna un permiso a un rol
func (rpc *RolePermissionController) CreateRolePermission(c *fiber.Ctx) error {
	// Parsear request body
	var req rolePermissionDomain.CreateRolePermissionDto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error parseando datos de entrada",
			Error:   "Formato de datos inválido",
		})
	}

	// Validar estructura de datos
	if errors := rpc.validator.ValidateStruct(req); len(errors) > 0 {
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

	// Crear relación rol-permiso a través del servicio
	result, err := rpc.rolePermissionService.CreateRolePermission(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error asignando permiso al rol",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(types.APIResponse{
		Success: true,
		Message: result.Message,
		Data:    result.RolePermission,
	})
}

// DeleteRolePermission maneja DELETE /role-permissions/unassign - Elimina un permiso de un rol
func (rpc *RolePermissionController) DeleteRolePermission(c *fiber.Ctx) error {
	// Parsear query parameters
	roleIDStr := c.Query("role_id")
	permissionIDStr := c.Query("permission_id")

	if roleIDStr == "" || permissionIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "role_id y permission_id son requeridos",
			Error:   "Parámetros faltantes",
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

	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de permiso inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Eliminar relación rol-permiso a través del servicio
	if err := rpc.rolePermissionService.DeleteRolePermission(c.Context(), roleID, permissionID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error eliminando permiso del rol",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Permiso eliminado del rol exitosamente",
	})
}

// GetRolePermissions maneja GET /role-permissions/roles/:roleID - Obtiene todos los permisos de un rol
func (rpc *RolePermissionController) GetRolePermissions(c *fiber.Ctx) error {
	// Parsear ID del rol
	roleIDStr := c.Params("roleID")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de rol inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Obtener permisos del rol desde el servicio
	permissions, err := rpc.rolePermissionService.GetAllPermissionsByRoleID(c.Context(), roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo permisos del rol",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Permisos del rol obtenidos exitosamente",
		Data:    permissions,
	})
}

// GetPermissionRoles maneja GET /role-permissions/permissions/:permissionID - Obtiene todos los roles de un permiso
func (rpc *RolePermissionController) GetPermissionRoles(c *fiber.Ctx) error {
	// Parsear ID del permiso
	permissionIDStr := c.Params("permissionID")
	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de permiso inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Obtener roles del permiso desde el servicio
	rolePermissions, err := rpc.rolePermissionService.GetPermissionRoles(c.Context(), permissionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo roles del permiso",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Roles del permiso obtenidos exitosamente",
		Data:    rolePermissions,
	})
}

// SyncRolePermissions maneja PUT /role-permissions/roles/:roleID/sync - Sincroniza permisos de un rol
func (rpc *RolePermissionController) SyncRolePermissions(c *fiber.Ctx) error {
	// Parsear ID del rol
	roleIDStr := c.Params("roleID")
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
		PermissionIDs []string `json:"permission_ids" validate:"required"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error parseando datos de entrada",
			Error:   "Formato de datos inválido",
		})
	}

	// Convertir strings a UUIDs
	var permissionIDs []uuid.UUID
	for _, permissionIDStr := range reqBody.PermissionIDs {
		permissionID, err := uuid.Parse(permissionIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
				Success: false,
				Message: "ID de permiso inválido",
				Error:   "Todos los IDs deben ser UUIDs válidos",
			})
		}
		permissionIDs = append(permissionIDs, permissionID)
	}

	// Crear DTO de sincronización
	req := &rolePermissionDomain.SyncRolePermissionsDto{
		RoleID:        roleID,
		PermissionIDs: permissionIDs,
	}

	// Validar estructura de datos
	if errors := rpc.validator.ValidateStruct(req); len(errors) > 0 {
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

	// Sincronizar permisos del rol a través del servicio
	if err := rpc.rolePermissionService.SyncRolePermissions(c.Context(), req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error sincronizando permisos del rol",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Permisos del rol sincronizados exitosamente",
	})
}

// CheckRoleHasPermission maneja GET /role-permissions/check - Verifica si un rol tiene un permiso
func (rpc *RolePermissionController) CheckRoleHasPermission(c *fiber.Ctx) error {
	// Parsear query parameters
	roleIDStr := c.Query("role_id")
	permissionIDStr := c.Query("permission_id")

	if roleIDStr == "" || permissionIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "role_id y permission_id son requeridos",
			Error:   "Parámetros faltantes",
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

	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de permiso inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Verificar si el rol tiene el permiso
	hasPermission, err := rpc.rolePermissionService.RoleHasPermission(c.Context(), roleID, permissionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error verificando permiso del rol",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Verificación de permiso completada",
		Data: fiber.Map{
			"role_id":        roleID,
			"permission_id":  permissionID,
			"has_permission": hasPermission,
		},
	})
}

// RemoveAllPermissionsFromRole maneja DELETE /role-permissions/roles/:roleID/all - Elimina todos los permisos de un rol
func (rpc *RolePermissionController) RemoveAllPermissionsFromRole(c *fiber.Ctx) error {
	// Parsear ID del rol
	roleIDStr := c.Params("roleID")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de rol inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Eliminar todos los permisos del rol a través del servicio
	if err := rpc.rolePermissionService.RemoveAllPermissionsFromRole(c.Context(), roleID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error eliminando todos los permisos del rol",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Todos los permisos eliminados del rol exitosamente",
	})
}

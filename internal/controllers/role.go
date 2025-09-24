package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	roleDomain "github.com/JoseLuis21/mv-backend/internal/core/role/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/role/ports"
	"github.com/JoseLuis21/mv-backend/internal/shared/pagination"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
)

// RoleController maneja las operaciones CRUD de roles
type RoleController struct {
	roleService ports.RoleService
	validator   *validatorapi.Validator
}

// NewRoleController crea una nueva instancia del controller de roles
func NewRoleController(roleService ports.RoleService, validator *validatorapi.Validator) *RoleController {
	return &RoleController{
		roleService: roleService,
		validator:   validator,
	}
}

// GetRoles maneja GET /roles - Lista roles con paginación y filtros
func (rc *RoleController) GetRoles(c *fiber.Ctx) error {
	// Parsear parámetros de paginación
	paginationReq, err := pagination.ParsePaginationFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Parámetros de paginación inválidos",
			Error:   err.Error(),
		})
	}

	// Extraer TenantID del contexto (puede ser nil para roles globales)
	var tenantID *uuid.UUID
	if tenantIDStr := c.Locals("tenantId"); tenantIDStr != nil {
		if tenantIDStr != "" {
			tid, err := uuid.Parse(tenantIDStr.(string))
			if err == nil {
				tenantID = &tid
			}
		}
	}

	// Crear filtro de roles
	filter := &roleDomain.RoleFilterRequest{
		TenantID: tenantID,
		Name:     paginationReq.Search,
		Page:     paginationReq.Page,
		Limit:    paginationReq.PageSize,
	}

	// Obtener roles desde el servicio
	roleListResponse, err := rc.roleService.GetRoles(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo roles",
			Error:   err.Error(),
		})
	}

	// Calcular información de paginación
	paginationInfo := paginationReq.CalculatePagination(int64(roleListResponse.Total))

	return c.Status(fiber.StatusOK).JSON(types.PaginatedAPIResponse{
		Success:    true,
		Message:    "Roles obtenidos exitosamente",
		Data:       roleListResponse.Roles,
		Pagination: paginationInfo,
	})
}

// GetRoleByID maneja GET /roles/:id - Obtiene un rol específico
func (rc *RoleController) GetRoleByID(c *fiber.Ctx) error {
	// Parsear ID del rol
	idStr := c.Params("id")
	roleID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de rol inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Obtener rol desde el servicio
	role, err := rc.roleService.GetRoleByID(c.Context(), roleID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(types.APIResponse{
			Success: false,
			Message: "Rol no encontrado",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Rol obtenido exitosamente",
		Data:    role,
	})
}

// CreateRole maneja POST /roles - Crea un nuevo rol
func (rc *RoleController) CreateRole(c *fiber.Ctx) error {
	// Parsear request body
	var req roleDomain.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error parseando datos de entrada",
			Error:   "Formato de datos inválido",
		})
	}

	// Extraer TenantID del contexto para roles de tenant
	if tenantIDStr := c.Locals("tenantId"); tenantIDStr != nil && tenantIDStr != "" {
		tenantID, err := uuid.Parse(tenantIDStr.(string))
		if err == nil {
			req.TenantID = &tenantID
		}
	}

	// Validar estructura de datos
	if errors := rc.validator.ValidateStruct(req); len(errors) > 0 {
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

	// Crear rol a través del servicio
	role, err := rc.roleService.CreateRole(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error creando rol",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(types.APIResponse{
		Success: true,
		Message: "Rol creado exitosamente",
		Data:    role,
	})
}

// UpdateRole maneja PUT /roles/:id - Actualiza un rol
func (rc *RoleController) UpdateRole(c *fiber.Ctx) error {
	// Parsear ID del rol
	idStr := c.Params("id")
	roleID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de rol inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Parsear request body
	var req roleDomain.UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error parseando datos de entrada",
			Error:   "Formato de datos inválido",
		})
	}

	// Validar estructura de datos
	if errors := rc.validator.ValidateStruct(req); len(errors) > 0 {
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

	// Actualizar rol a través del servicio
	role, err := rc.roleService.UpdateRole(c.Context(), roleID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error actualizando rol",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Rol actualizado exitosamente",
		Data:    role,
	})
}

// DeleteRole maneja DELETE /roles/:id - Elimina (soft delete) un rol
func (rc *RoleController) DeleteRole(c *fiber.Ctx) error {
	// Parsear ID del rol
	idStr := c.Params("id")
	roleID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de rol inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Eliminar rol a través del servicio
	if err := rc.roleService.DeleteRole(c.Context(), roleID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error eliminando rol",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Rol eliminado exitosamente",
	})
}

// GetGlobalRoles maneja GET /roles/global - Obtiene todos los roles globales del sistema
func (rc *RoleController) GetGlobalRoles(c *fiber.Ctx) error {
	// Obtener roles globales desde el servicio
	roles, err := rc.roleService.GetGlobalRoles(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo roles globales",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Roles globales obtenidos exitosamente",
		Data:    roles,
	})
}

// GetTenantRoles maneja GET /roles/tenant/:tenantId - Obtiene roles de un tenant específico
func (rc *RoleController) GetTenantRoles(c *fiber.Ctx) error {
	// Parsear ID del tenant
	tenantIDStr := c.Params("tenantId")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de tenant inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Obtener roles del tenant desde el servicio
	roles, err := rc.roleService.GetTenantRoles(c.Context(), tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo roles del tenant",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Roles del tenant obtenidos exitosamente",
		Data:    roles,
	})
}

// GetUserRoles maneja GET /roles/user/:userId - Obtiene roles de un usuario específico
func (rc *RoleController) GetUserRoles(c *fiber.Ctx) error {
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

	// Extraer TenantID del contexto (puede ser nil para roles globales)
	var tenantID *uuid.UUID
	if tenantIDStr := c.Locals("tenantId"); tenantIDStr != nil {
		if tenantIDStr != "" {
			tid, err := uuid.Parse(tenantIDStr.(string))
			if err == nil {
				tenantID = &tid
			}
		}
	}

	// Obtener roles del usuario desde el servicio
	roles, err := rc.roleService.GetUserRoles(c.Context(), userID, tenantID)
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
		Data:    roles,
	})
}

// InitializeSystemRoles maneja POST /roles/initialize-system - Inicializa roles del sistema
func (rc *RoleController) InitializeSystemRoles(c *fiber.Ctx) error {
	// Inicializar roles del sistema
	if err := rc.roleService.InitializeSystemRoles(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error inicializando roles del sistema",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Roles del sistema inicializados exitosamente",
	})
}

// InitializeTenantRoles maneja POST /roles/initialize-tenant/:tenantId - Inicializa roles de un tenant
func (rc *RoleController) InitializeTenantRoles(c *fiber.Ctx) error {
	// Parsear ID del tenant
	tenantIDStr := c.Params("tenantId")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de tenant inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Inicializar roles del tenant
	if err := rc.roleService.InitializeTenantRoles(c.Context(), tenantID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error inicializando roles del tenant",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Roles del tenant inicializados exitosamente",
	})
}
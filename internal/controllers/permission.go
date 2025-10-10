package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	permissionDomain "github.com/JoseLuis21/mv-backend/internal/core/permission/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/permission/ports"
	"github.com/JoseLuis21/mv-backend/internal/shared/pagination"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
)

// PermissionController maneja las operaciones CRUD de permisos
type PermissionController struct {
	permissionService ports.PermissionService
	validator         *validatorapi.Validator
}

// NewPermissionController crea una nueva instancia del controller de permisos
func NewPermissionController(permissionService ports.PermissionService, validator *validatorapi.Validator) *PermissionController {
	return &PermissionController{
		permissionService: permissionService,
		validator:         validator,
	}
}

// GetPermissions maneja GET /permissions - Lista permisos con paginación y filtros
func (pc *PermissionController) GetPermissions(c *fiber.Ctx) error {
	// Parsear parámetros de paginación
	paginationReq, err := pagination.ParsePaginationFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Parámetros de paginación inválidos",
			Error:   err.Error(),
		})
	}

	// Crear filtro de permisos
	filter := &permissionDomain.PermissionFilterRequest{
		Name:    paginationReq.Search,
		Section: c.Query("section"),
		Page:    paginationReq.Page,
		Limit:   paginationReq.PageSize,
	}

	// Obtener permisos desde el servicio
	permissionListResponse, err := pc.permissionService.GetPermissions(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo permisos",
			Error:   err.Error(),
		})
	}

	// Calcular información de paginación
	paginationInfo := paginationReq.CalculatePagination(int64(permissionListResponse.Total))

	return c.Status(fiber.StatusOK).JSON(types.PaginatedAPIResponse{
		Success:    true,
		Message:    "Permisos obtenidos exitosamente",
		Data:       permissionListResponse.Permissions,
		Pagination: paginationInfo,
	})
}

// GetPermissionByID maneja GET /permissions/:id - Obtiene un permiso específico
func (pc *PermissionController) GetPermissionByID(c *fiber.Ctx) error {
	// Parsear ID del permiso
	idStr := c.Params("id")
	permissionID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de permiso inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Obtener permiso desde el servicio
	permission, err := pc.permissionService.GetPermissionByID(c.Context(), permissionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(types.APIResponse{
			Success: false,
			Message: "Permiso no encontrado",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Permiso obtenido exitosamente",
		Data:    permission,
	})
}

// CreatePermission maneja POST /permissions - Crea un nuevo permiso
func (pc *PermissionController) CreatePermission(c *fiber.Ctx) error {
	// Parsear request body
	var req permissionDomain.CreatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error parseando datos de entrada",
			Error:   "Formato de datos inválido",
		})
	}

	// Validar estructura de datos
	if errors := pc.validator.ValidateStruct(req); len(errors) > 0 {
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

	// Crear permiso a través del servicio
	permission, err := pc.permissionService.CreatePermission(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error creando permiso",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(types.APIResponse{
		Success: true,
		Message: "Permiso creado exitosamente",
		Data:    permission,
	})
}

// UpdatePermission maneja PUT /permissions/:id - Actualiza un permiso
func (pc *PermissionController) UpdatePermission(c *fiber.Ctx) error {
	// Parsear ID del permiso
	idStr := c.Params("id")
	permissionID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de permiso inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Parsear request body
	var req permissionDomain.UpdatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error parseando datos de entrada",
			Error:   "Formato de datos inválido",
		})
	}

	// Validar estructura de datos
	if errors := pc.validator.ValidateStruct(req); len(errors) > 0 {
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

	// Actualizar permiso a través del servicio
	permission, err := pc.permissionService.UpdatePermission(c.Context(), permissionID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error actualizando permiso",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Permiso actualizado exitosamente",
		Data:    permission,
	})
}

// DeletePermission maneja DELETE /permissions/:id - Elimina (soft delete) un permiso
func (pc *PermissionController) DeletePermission(c *fiber.Ctx) error {
	// Parsear ID del permiso
	idStr := c.Params("id")
	permissionID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de permiso inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Eliminar permiso a través del servicio
	if err := pc.permissionService.DeletePermission(c.Context(), permissionID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error eliminando permiso",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Permiso eliminado exitosamente",
	})
}

// GetPermissionsBySection maneja GET /permissions/section/:section - Obtiene permisos por sección
func (pc *PermissionController) GetPermissionsBySection(c *fiber.Ctx) error {
	// Obtener sección del parámetro
	section := c.Params("section")
	if section == "" {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Sección requerida",
			Error:   "Debe especificar una sección válida",
		})
	}

	// Obtener permisos por sección desde el servicio
	permissions, err := pc.permissionService.GetPermissionsBySection(c.Context(), section)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo permisos por sección",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Permisos obtenidos exitosamente",
		Data:    permissions,
	})
}

// GetPermissionsGrouped maneja GET /permissions/grouped - Obtiene permisos agrupados por sección
func (pc *PermissionController) GetPermissionsGrouped(c *fiber.Ctx) error {
	// Obtener permisos agrupados desde el servicio
	groupedPermissions, err := pc.permissionService.GetPermissionsGroupedBySection(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo permisos agrupados",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Permisos agrupados obtenidos exitosamente",
		Data:    groupedPermissions,
	})
}

// GetRolePermissions maneja GET /permissions/role/:roleId - Obtiene permisos de un rol específico
func (pc *PermissionController) GetRolePermissions(c *fiber.Ctx) error {
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

	// Obtener permisos del rol desde el servicio
	permissions, err := pc.permissionService.GetRolePermissions(c.Context(), roleID)
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

// GetAvailableSections maneja GET /permissions/sections - Obtiene todas las secciones disponibles
func (pc *PermissionController) GetAvailableSections(c *fiber.Ctx) error {
	// Obtener secciones disponibles desde el servicio
	sections, err := pc.permissionService.GetAvailableSections(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo secciones disponibles",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Secciones disponibles obtenidas exitosamente",
		Data:    sections,
	})
}

// InitializeSystemPermissions maneja POST /permissions/initialize-system - Inicializa permisos del sistema
func (pc *PermissionController) InitializeSystemPermissions(c *fiber.Ctx) error {
	// Inicializar permisos del sistema
	if err := pc.permissionService.InitializeSystemPermissions(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error inicializando permisos del sistema",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Permisos del sistema inicializados exitosamente",
	})
}
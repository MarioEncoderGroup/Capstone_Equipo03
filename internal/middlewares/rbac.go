package middlewares

import (
	"context"
	"fmt"
	"strings"

	permissionPorts "github.com/JoseLuis21/mv-backend/internal/core/permission/ports"
	rolePorts "github.com/JoseLuis21/mv-backend/internal/core/role/ports"
	sharedAuth "github.com/JoseLuis21/mv-backend/internal/shared/auth"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
	sharedTypes "github.com/JoseLuis21/mv-backend/internal/shared/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RBACMiddleware maneja el control de acceso basado en roles y permisos
type RBACMiddleware struct {
	roleService       rolePorts.RoleService
	permissionService permissionPorts.PermissionService
}

// NewRBACMiddleware crea una nueva instancia del middleware RBAC
func NewRBACMiddleware(roleService rolePorts.RoleService, permissionService permissionPorts.PermissionService) *RBACMiddleware {
	return &RBACMiddleware{
		roleService:       roleService,
		permissionService: permissionService,
	}
}

// RequirePermission valida que el usuario tenga el permiso específico
func (m *RBACMiddleware) RequirePermission(permissionName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		// 1. PRIMERO: Intentar obtener permisos del JWT (guardados en contexto por AuthMiddleware)
		permissionsInterface := c.Locals("permissions")
		if permissionsInterface != nil {
			// Los permisos están en el JWT - validación rápida sin DB
			if userPermissions, ok := permissionsInterface.([]string); ok {
				// Verificar si tiene el permiso requerido (case-insensitive)
				for _, userPermission := range userPermissions {
					if strings.EqualFold(userPermission, permissionName) {
						return c.Next()
					}
				}

				// Si no tiene el permiso, denegar acceso
				return c.Status(fiber.StatusForbidden).JSON(sharedTypes.APIResponse{
					Success: false,
					Message: "Acceso denegado: permiso insuficiente",
					Error:   fmt.Sprintf("Requiere el permiso: %s. Permisos actuales: %s", permissionName, strings.Join(userPermissions, ", ")),
				})
			}
		}

		// 2. FALLBACK: Si no hay permisos en el JWT (token viejo o sin tenant), consultar BD
		user, err := sharedAuth.GetUserFromContext(c)
		if err != nil {
			return sharedAuth.HandleAuthError(c, sharedErrors.NewAuthError("Usuario no autenticado", "UNAUTHENTICATED"))
		}

		// Obtener tenant del contexto si existe
		var tenantID *uuid.UUID
		if tenantIDStr := c.Get("X-Tenant-ID"); tenantIDStr != "" {
			tid, err := uuid.Parse(tenantIDStr)
			if err == nil {
				tenantID = &tid
			}
		}

		// Verificar si el usuario tiene el permiso (desde BD)
		hasPermission, err := m.checkUserPermission(ctx, user.ID, permissionName, tenantID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(sharedTypes.APIResponse{
				Success: false,
				Message: "Error verificando permisos",
				Error:   err.Error(),
			})
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(sharedTypes.APIResponse{
				Success: false,
				Message: "Acceso denegado: permiso insuficiente",
				Error:   fmt.Sprintf("Requiere el permiso: %s", permissionName),
			})
		}

		return c.Next()
	}
}

// RequireRole valida que el usuario tenga uno de los roles especificados
func (m *RBACMiddleware) RequireRole(roleNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		// 1. PRIMERO: Intentar obtener roles del JWT (guardados en contexto por AuthMiddleware)
		rolesInterface := c.Locals("roles")
		if rolesInterface != nil {
			// Los roles están en el JWT - validación rápida sin DB
			if userRoles, ok := rolesInterface.([]string); ok {
				// Verificar si tiene alguno de los roles requeridos (case-insensitive)
				for _, userRole := range userRoles {
					for _, requiredRole := range roleNames {
						if strings.EqualFold(userRole, requiredRole) {
							return c.Next()
						}
					}
				}

				// Si no tiene el rol, denegar acceso
				return c.Status(fiber.StatusForbidden).JSON(sharedTypes.APIResponse{
					Success: false,
					Message: "Acceso denegado: rol insuficiente",
					Error:   fmt.Sprintf("Requiere uno de los roles: %s. Roles actuales: %s", strings.Join(roleNames, ", "), strings.Join(userRoles, ", ")),
				})
			}
		}

		// 2. FALLBACK: Si no hay roles en el JWT (token viejo o sin tenant), consultar BD
		user, err := sharedAuth.GetUserFromContext(c)
		if err != nil {
			return sharedAuth.HandleAuthError(c, sharedErrors.NewAuthError("Usuario no autenticado", "UNAUTHENTICATED"))
		}

		// Obtener tenant del contexto si existe
		var tenantID *uuid.UUID
		if tenantIDStr := c.Get("X-Tenant-ID"); tenantIDStr != "" {
			tid, err := uuid.Parse(tenantIDStr)
			if err == nil {
				tenantID = &tid
			}
		}

		// Verificar si el usuario tiene alguno de los roles (desde BD)
		hasRole, err := m.checkUserRoles(ctx, user.ID, roleNames, tenantID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(sharedTypes.APIResponse{
				Success: false,
				Message: "Error verificando roles",
				Error:   err.Error(),
			})
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(sharedTypes.APIResponse{
				Success: false,
				Message: "Acceso denegado: rol insuficiente",
				Error:   fmt.Sprintf("Requiere uno de los roles: %s", strings.Join(roleNames, ", ")),
			})
		}

		return c.Next()
	}
}

// RequireAdminRole valida que el usuario sea administrador
func (m *RBACMiddleware) RequireAdminRole() fiber.Handler {
	return m.RequireRole("administrator")
}

// RequireSystemAdmin valida que el usuario sea administrador del sistema (global)
func (m *RBACMiddleware) RequireSystemAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		// Obtener usuario autenticado del contexto
		user, err := sharedAuth.GetUserFromContext(c)
		if err != nil {
			return sharedAuth.HandleAuthError(c, sharedErrors.NewAuthError("Usuario no autenticado", "UNAUTHENTICATED"))
		}

		// Verificar si el usuario tiene roles de sistema (tenant_id = null)
		userRoles, err := m.roleService.GetUserRoles(ctx, user.ID, nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(sharedTypes.APIResponse{
				Success: false,
				Message: "Error verificando roles de sistema",
				Error:   err.Error(),
			})
		}

		// Buscar rol de administrador del sistema
		isSystemAdmin := false
		for _, role := range userRoles {
			if role.Name == "administrator" {
				isSystemAdmin = true
				break
			}
		}

		if !isSystemAdmin {
			return c.Status(fiber.StatusForbidden).JSON(sharedTypes.APIResponse{
				Success: false,
				Message: "Acceso denegado: requiere permisos de administrador del sistema",
				Error:   "SYSTEM_ADMIN_REQUIRED",
			})
		}

		return c.Next()
	}
}

// RequireTenantAccess valida que el usuario tenga acceso al tenant especificado
func (m *RBACMiddleware) RequireTenantAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obtener usuario autenticado del contexto
		user, err := sharedAuth.GetUserFromContext(c)
		if err != nil {
			return sharedAuth.HandleAuthError(c, sharedErrors.NewAuthError("Usuario no autenticado", "UNAUTHENTICATED"))
		}

		// Obtener tenant ID de los parámetros de ruta o headers
		var tenantID uuid.UUID

		// Intentar obtener de parámetros de ruta primero
		if tenantIDParam := c.Params("tenantId"); tenantIDParam != "" {
			tid, err := uuid.Parse(tenantIDParam)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(sharedTypes.APIResponse{
					Success: false,
					Message: "ID de tenant inválido",
					Error:   err.Error(),
				})
			}
			tenantID = tid
		} else if tenantIDHeader := c.Get("X-Tenant-ID"); tenantIDHeader != "" {
			// Intentar obtener del header
			tid, err := uuid.Parse(tenantIDHeader)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(sharedTypes.APIResponse{
					Success: false,
					Message: "ID de tenant inválido en header",
					Error:   err.Error(),
				})
			}
			tenantID = tid
		} else {
			// Si no hay tenant especificado, verificar que tenga permisos globales
			return m.RequireSystemAdmin()(c)
		}

		// Verificar si el usuario tiene acceso al tenant
		// Esto se podría implementar en el servicio de usuarios
		// Por ahora verificamos si tiene algún rol en ese tenant
		ctx := c.Context()
		userRoles, err := m.roleService.GetUserRoles(ctx, user.ID, &tenantID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(sharedTypes.APIResponse{
				Success: false,
				Message: "Error verificando acceso al tenant",
				Error:   err.Error(),
			})
		}

		if len(userRoles) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(sharedTypes.APIResponse{
				Success: false,
				Message: "Acceso denegado: sin permisos en este tenant",
				Error:   "TENANT_ACCESS_DENIED",
			})
		}

		// Agregar tenant ID al contexto para uso posterior
		c.Locals("tenantID", tenantID)

		return c.Next()
	}
}

// checkUserPermission verifica si un usuario tiene un permiso específico
func (m *RBACMiddleware) checkUserPermission(ctx context.Context, userID uuid.UUID, permissionName string, tenantID *uuid.UUID) (bool, error) {
	// Obtener roles del usuario
	userRoles, err := m.roleService.GetUserRoles(ctx, userID, tenantID)
	if err != nil {
		return false, fmt.Errorf("error obteniendo roles del usuario: %w", err)
	}

	// Para cada rol, obtener sus permisos
	for _, role := range userRoles {
		roleID := role.ID

		// Obtener permisos del rol
		rolePermissions, err := m.permissionService.GetRolePermissions(ctx, roleID)
		if err != nil {
			continue
		}

		// Verificar si el permiso existe en este rol
		for _, permission := range rolePermissions {
			if permission.Name == permissionName {
				return true, nil
			}
		}
	}

	// También verificar roles globales si no se encontró en el tenant
	if tenantID != nil {
		globalRoles, err := m.roleService.GetUserRoles(ctx, userID, nil)
		if err != nil {
			return false, fmt.Errorf("error obteniendo roles globales: %w", err)
		}

		for _, role := range globalRoles {
			roleID := role.ID

			rolePermissions, err := m.permissionService.GetRolePermissions(ctx, roleID)
			if err != nil {
				continue
			}

			for _, permission := range rolePermissions {
				if permission.Name == permissionName {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// checkUserRoles verifica si un usuario tiene alguno de los roles especificados
func (m *RBACMiddleware) checkUserRoles(ctx context.Context, userID uuid.UUID, roleNames []string, tenantID *uuid.UUID) (bool, error) {
	// Obtener roles del usuario
	userRoles, err := m.roleService.GetUserRoles(ctx, userID, tenantID)
	if err != nil {
		return false, fmt.Errorf("error obteniendo roles del usuario: %w", err)
	}

	// Verificar si el usuario tiene alguno de los roles requeridos
	for _, userRole := range userRoles {
		for _, requiredRole := range roleNames {
			if userRole.Name == requiredRole {
				return true, nil
			}
		}
	}

	// También verificar roles globales si no se encontró en el tenant
	if tenantID != nil {
		globalRoles, err := m.roleService.GetUserRoles(ctx, userID, nil)
		if err != nil {
			return false, fmt.Errorf("error obteniendo roles globales: %w", err)
		}

		for _, userRole := range globalRoles {
			for _, requiredRole := range roleNames {
				if userRole.Name == requiredRole {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// Optional: RequirePermissionOrRole permite acceso con permiso O rol
func (m *RBACMiddleware) RequirePermissionOrRole(permissionName string, roleNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		user, err := sharedAuth.GetUserFromContext(c)
		if err != nil {
			return sharedAuth.HandleAuthError(c, sharedErrors.NewAuthError("Usuario no autenticado", "UNAUTHENTICATED"))
		}

		var tenantID *uuid.UUID
		if tenantIDStr := c.Get("X-Tenant-ID"); tenantIDStr != "" {
			tid, err := uuid.Parse(tenantIDStr)
			if err == nil {
				tenantID = &tid
			}
		}

		// Verificar permiso
		hasPermission, err := m.checkUserPermission(ctx, user.ID, permissionName, tenantID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(sharedTypes.APIResponse{
				Success: false,
				Message: "Error verificando permisos",
				Error:   err.Error(),
			})
		}

		if hasPermission {
			return c.Next()
		}

		// Si no tiene el permiso, verificar roles
		hasRole, err := m.checkUserRoles(ctx, user.ID, roleNames, tenantID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(sharedTypes.APIResponse{
				Success: false,
				Message: "Error verificando roles",
				Error:   err.Error(),
			})
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(sharedTypes.APIResponse{
				Success: false,
				Message: "Acceso denegado: permiso o rol insuficiente",
				Error:   fmt.Sprintf("Requiere el permiso '%s' o uno de los roles: %s", permissionName, strings.Join(roleNames, ", ")),
			})
		}

		return c.Next()
	}
}

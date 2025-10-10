package ports

import (
	"context"
	"github.com/google/uuid"
	domain_role_permission "github.com/JoseLuis21/mv-backend/internal/core/role_permission/domain"
)

// RolePermissionRepository define las operaciones de acceso a datos para role_permission
type RolePermissionRepository interface {
	// Create crea una nueva relación rol-permiso
	Create(ctx context.Context, rolePermission *domain_role_permission.RolePermission) error

	// Delete elimina una relación rol-permiso específica
	Delete(ctx context.Context, roleID, permissionID uuid.UUID) error

	// GetByRoleID obtiene todos los permisos de un rol específico
	GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]domain_role_permission.RolePermission, error)

	// GetByPermissionID obtiene todos los roles de un permiso específico
	GetByPermissionID(ctx context.Context, permissionID uuid.UUID) ([]domain_role_permission.RolePermission, error)

	// Exists verifica si existe una relación rol-permiso específica
	Exists(ctx context.Context, roleID, permissionID uuid.UUID) (bool, error)

	// SyncRolePermissions sincroniza múltiples permisos a un rol (reemplaza existentes)
	SyncRolePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error

	// DeleteByRoleID elimina todas las relaciones de un rol
	DeleteByRoleID(ctx context.Context, roleID uuid.UUID) error

	// DeleteByPermissionID elimina todas las relaciones de un permiso
	DeleteByPermissionID(ctx context.Context, permissionID uuid.UUID) error

	// GetPermissionsByRoleID obtiene información completa de permisos por rol
	GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]domain_role_permission.Permission, error)

	// GetPermissionsByRoleIDs obtiene permisos por múltiples roles
	GetPermissionsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]domain_role_permission.Permission, error)

	// GetPermissionNamesByRoleIDs obtiene solo los nombres de permisos por múltiples roles
	GetPermissionNamesByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]string, error)
}

// RolePermissionService define la lógica de negocio para role_permission
type RolePermissionService interface {
	// CreateRolePermission crea una nueva relación rol-permiso con validaciones
	CreateRolePermission(ctx context.Context, req *domain_role_permission.CreateRolePermissionDto) (*domain_role_permission.RolePermissionResponseDto, error)

	// DeleteRolePermission elimina una relación rol-permiso específica
	DeleteRolePermission(ctx context.Context, roleID, permissionID uuid.UUID) error

	// GetRolePermissions obtiene todos los permisos de un rol
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]domain_role_permission.RolePermission, error)

	// GetPermissionRoles obtiene todos los roles de un permiso
	GetPermissionRoles(ctx context.Context, permissionID uuid.UUID) ([]domain_role_permission.RolePermission, error)

	// SyncRolePermissions sincroniza permisos a un rol con validaciones
	SyncRolePermissions(ctx context.Context, req *domain_role_permission.SyncRolePermissionsDto) error

	// RoleHasPermission verifica si un rol tiene un permiso específico
	RoleHasPermission(ctx context.Context, roleID, permissionID uuid.UUID) (bool, error)

	// RemoveAllPermissionsFromRole elimina todos los permisos de un rol
	RemoveAllPermissionsFromRole(ctx context.Context, roleID uuid.UUID) error

	// RemoveRoleFromAllPermissions elimina un rol de todos los permisos
	RemoveRoleFromAllPermissions(ctx context.Context, permissionID uuid.UUID) error

	// GetAllPermissionsByRoleID obtiene información completa de permisos por rol
	GetAllPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]domain_role_permission.Permission, error)

	// GetAllPermissionsByRoleIDs obtiene permisos por múltiples roles
	GetAllPermissionsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]domain_role_permission.Permission, error)

	// GetPermissionNamesByRoleIDs obtiene solo nombres de permisos por múltiples roles
	GetPermissionNamesByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]string, error)
}
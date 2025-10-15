package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/permission/domain"
)

// PermissionRepository define el contrato para la persistencia de permisos
// Esta interfaz será implementada por el adaptador de PostgreSQL
// Sigue el patrón Repository para aislar el dominio de la infraestructura
type PermissionRepository interface {
	// Create crea un nuevo permiso en la base de datos control
	Create(ctx context.Context, permission *domain.Permission) error

	// GetByID obtiene un permiso por su ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Permission, error)

	// GetByName obtiene un permiso por su nombre
	GetByName(ctx context.Context, name string) (*domain.Permission, error)

	// GetBySection obtiene todos los permisos de una sección específica
	GetBySection(ctx context.Context, section string) ([]*domain.Permission, error)

	// GetAllPermissions obtiene permisos con filtros de búsqueda
	GetAllPermissions(ctx context.Context, filter *domain.PermissionFilterRequest) ([]*domain.Permission, int, error)

	// GetGroupedBySection obtiene permisos agrupados por sección
	GetGroupedBySection(ctx context.Context) (map[string][]*domain.Permission, error)

	// Update actualiza un permiso existente
	Update(ctx context.Context, permission *domain.Permission) error

	// Delete elimina lógicamente un permiso (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// ExistsByName verifica si existe un permiso con el nombre dado
	ExistsByName(ctx context.Context, name string) (bool, error)

	// IsSystemPermission verifica si un permiso es uno de los permisos predefinidos del sistema
	IsSystemPermission(ctx context.Context, permissionID uuid.UUID) (bool, error)

	// GetRolePermissions obtiene todos los permisos asignados a un rol
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*domain.Permission, error)

	// GetAvailableSections obtiene todas las secciones de permisos disponibles
	GetAvailableSections(ctx context.Context) ([]string, error)
}

// PermissionService define el contrato para el servicio de permisos
// Esta interfaz encapsula la lógica de negocio relacionada con permisos
// y sigue los principios de Clean Architecture
type PermissionService interface {
	// CreatePermission crea un nuevo permiso
	CreatePermission(ctx context.Context, req *domain.CreatePermissionRequest) (*domain.PermissionResponse, error)

	// GetPermissionByID obtiene un permiso por su ID
	GetPermissionByID(ctx context.Context, id uuid.UUID) (*domain.PermissionResponse, error)

	// GetPermissionByName obtiene un permiso por su nombre
	GetPermissionByName(ctx context.Context, name string) (*domain.PermissionResponse, error)

	// GetPermissionsBySection obtiene todos los permisos de una sección específica
	GetPermissionsBySection(ctx context.Context, section string) ([]*domain.PermissionResponse, error)

	// GetPermissions obtiene permisos con filtros de búsqueda y paginación
	GetPermissions(ctx context.Context, filter *domain.PermissionFilterRequest) (*domain.PermissionListResponse, error)

	// GetPermissionsGroupedBySection obtiene permisos agrupados por sección
	GetPermissionsGroupedBySection(ctx context.Context) ([]*domain.PermissionGroupedResponse, error)

	// UpdatePermission actualiza un permiso existente
	UpdatePermission(ctx context.Context, id uuid.UUID, req *domain.UpdatePermissionRequest) (*domain.PermissionResponse, error)

	// DeletePermission elimina lógicamente un permiso
	DeletePermission(ctx context.Context, id uuid.UUID) error

	// ValidatePermissionCreation valida que se puede crear un permiso con los datos dados
	ValidatePermissionCreation(ctx context.Context, req *domain.CreatePermissionRequest) error

	// ValidatePermissionUpdate valida que se puede actualizar un permiso con los datos dados
	ValidatePermissionUpdate(ctx context.Context, id uuid.UUID, req *domain.UpdatePermissionRequest) error

	// InitializeSystemPermissions crea los permisos predefinidos del sistema si no existen
	InitializeSystemPermissions(ctx context.Context) error

	// GetRolePermissions obtiene todos los permisos asignados a un rol
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*domain.PermissionResponse, error)

	// GetAvailableSections obtiene todas las secciones de permisos disponibles
	GetAvailableSections(ctx context.Context) (*domain.PermissionSectionResponse, error)

	// CheckPermissionExistsByID verifica si un permiso existe por ID
	CheckPermissionExistsByID(ctx context.Context, id string) (bool, error)
}
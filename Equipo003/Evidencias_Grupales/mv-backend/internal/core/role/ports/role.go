package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/role/domain"
)

// RoleRepository define el contrato para la persistencia de roles
// Esta interfaz será implementada por el adaptador de PostgreSQL
// Sigue el patrón Repository para aislar el dominio de la infraestructura
type RoleRepository interface {
	// Create crea un nuevo rol en la base de datos (control o tenant según tipo)
	Create(ctx context.Context, role *domain.Role) error

	// GetByID obtiene un rol por su ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error)

	// GetByName obtiene un rol por su nombre
	// Para roles globales tenantID debe ser nil
	// Para roles de tenant, tenantID debe especificarse
	GetByName(ctx context.Context, name string, tenantID *uuid.UUID) (*domain.Role, error)

	// GetGlobalRoles obtiene todos los roles globales del sistema
	GetGlobalRoles(ctx context.Context) ([]*domain.Role, error)

	// GetTenantRoles obtiene todos los roles de un tenant específico
	GetTenantRoles(ctx context.Context, tenantID uuid.UUID) ([]*domain.Role, error)

	// GetAllRoles obtiene roles con filtros de búsqueda
	GetAllRoles(ctx context.Context, filter *domain.RoleFilterRequest) ([]*domain.Role, int, error)

	// Update actualiza un rol existente
	Update(ctx context.Context, role *domain.Role) error

	// Delete elimina lógicamente un rol (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// ExistsByName verifica si existe un rol con el nombre dado
	// Para roles globales tenantID debe ser nil
	// Para roles de tenant, tenantID debe especificarse
	ExistsByName(ctx context.Context, name string, tenantID *uuid.UUID) (bool, error)

	// IsSystemRole verifica si un rol es uno de los roles predefinidos del sistema
	IsSystemRole(ctx context.Context, roleID uuid.UUID) (bool, error)

	// GetUserRoles obtiene todos los roles asignados a un usuario
	GetUserRoles(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) ([]*domain.Role, error)
}

// RoleService define el contrato para el servicio de roles
// Esta interfaz encapsula la lógica de negocio relacionada con roles
// y sigue los principios de Clean Architecture
type RoleService interface {
	// CreateRole crea un nuevo rol
	CreateRole(ctx context.Context, req *domain.CreateRoleRequest) (*domain.RoleResponse, error)

	// GetRoleByID obtiene un rol por su ID
	GetRoleByID(ctx context.Context, id uuid.UUID) (*domain.RoleResponse, error)

	// GetRoleByName obtiene un rol por su nombre
	GetRoleByName(ctx context.Context, name string, tenantID *uuid.UUID) (*domain.RoleResponse, error)

	// GetGlobalRoles obtiene todos los roles globales del sistema
	GetGlobalRoles(ctx context.Context) ([]*domain.RoleResponse, error)

	// GetTenantRoles obtiene todos los roles de un tenant específico
	GetTenantRoles(ctx context.Context, tenantID uuid.UUID) ([]*domain.RoleResponse, error)

	// GetRoles obtiene roles con filtros de búsqueda y paginación
	GetRoles(ctx context.Context, filter *domain.RoleFilterRequest) (*domain.RoleListResponse, error)

	// UpdateRole actualiza un rol existente
	UpdateRole(ctx context.Context, id uuid.UUID, req *domain.UpdateRoleRequest) (*domain.RoleResponse, error)

	// DeleteRole elimina lógicamente un rol
	DeleteRole(ctx context.Context, id uuid.UUID) error

	// ValidateRoleCreation valida que se puede crear un rol con los datos dados
	ValidateRoleCreation(ctx context.Context, req *domain.CreateRoleRequest) error

	// ValidateRoleUpdate valida que se puede actualizar un rol con los datos dados
	ValidateRoleUpdate(ctx context.Context, id uuid.UUID, req *domain.UpdateRoleRequest) error

	// InitializeSystemRoles crea los roles predefinidos del sistema si no existen
	InitializeSystemRoles(ctx context.Context) error

	// InitializeTenantRoles crea los roles por defecto para un nuevo tenant
	InitializeTenantRoles(ctx context.Context, tenantID uuid.UUID) error

	// GetUserRoles obtiene todos los roles asignados a un usuario
	GetUserRoles(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) ([]*domain.RoleResponse, error)

	// CheckRoleExistsByID verifica si un rol existe por ID
	CheckRoleExistsByID(ctx context.Context, id, tenantID string) (bool, error)

	// CheckRoleExistsByName verifica si un rol existe por nombre
	CheckRoleExistsByName(ctx context.Context, name, tenantID string) (bool, error)
}
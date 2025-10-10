package ports

import (
	"context"
	"github.com/google/uuid"
	domain_user_role "github.com/JoseLuis21/mv-backend/internal/core/user_role/domain"
)

// UserRoleRepository define las operaciones de acceso a datos para user_role
type UserRoleRepository interface {
	// Create crea una nueva relación usuario-rol
	Create(ctx context.Context, userRole *domain_user_role.UserRole) error

	// Delete elimina una relación usuario-rol específica
	Delete(ctx context.Context, userID, roleID uuid.UUID, tenantID *uuid.UUID) error

	// GetByUserID obtiene todos los roles de un usuario específico
	GetByUserID(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) ([]domain_user_role.UserRole, error)

	// GetByRoleID obtiene todos los usuarios de un rol específico
	GetByRoleID(ctx context.Context, roleID uuid.UUID, tenantID *uuid.UUID) ([]domain_user_role.UserRole, error)

	// Exists verifica si existe una relación usuario-rol específica
	Exists(ctx context.Context, userID, roleID uuid.UUID, tenantID *uuid.UUID) (bool, error)

	// SyncRoleUsers sincroniza múltiples usuarios a un rol (reemplaza existentes)
	SyncRoleUsers(ctx context.Context, roleID uuid.UUID, userIDs []uuid.UUID, tenantID *uuid.UUID) error

	// SyncUserRoles sincroniza múltiples roles a un usuario (reemplaza existentes)
	SyncUserRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID, tenantID *uuid.UUID) error

	// DeleteByUserID elimina todas las relaciones de un usuario
	DeleteByUserID(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) error

	// DeleteByRoleID elimina todas las relaciones de un rol
	DeleteByRoleID(ctx context.Context, roleID uuid.UUID, tenantID *uuid.UUID) error
}

// UserRoleService define la lógica de negocio para user_role
type UserRoleService interface {
	// CreateUserRole crea una nueva relación usuario-rol con validaciones
	CreateUserRole(ctx context.Context, req *domain_user_role.CreateUserRoleDto) (*domain_user_role.UserRoleResponseDto, error)

	// DeleteUserRole elimina una relación usuario-rol específica
	DeleteUserRole(ctx context.Context, userID, roleID uuid.UUID, tenantID *uuid.UUID) error

	// GetUserRoles obtiene todos los roles de un usuario
	GetUserRoles(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) ([]domain_user_role.UserRole, error)

	// GetRoleUsers obtiene todos los usuarios de un rol
	GetRoleUsers(ctx context.Context, roleID uuid.UUID, tenantID *uuid.UUID) ([]domain_user_role.UserRole, error)

	// SyncRoleUsers sincroniza usuarios a un rol con validaciones
	SyncRoleUsers(ctx context.Context, req *domain_user_role.SyncRoleUsersDto) error

	// SyncUserRoles sincroniza roles a un usuario con validaciones
	SyncUserRoles(ctx context.Context, req *domain_user_role.SyncUserRolesDto) error

	// UserHasRole verifica si un usuario tiene un rol específico
	UserHasRole(ctx context.Context, userID, roleID uuid.UUID, tenantID *uuid.UUID) (bool, error)

	// RemoveUserFromAllRoles elimina un usuario de todos sus roles
	RemoveUserFromAllRoles(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) error

	// RemoveAllUsersFromRole elimina todos los usuarios de un rol
	RemoveAllUsersFromRole(ctx context.Context, roleID uuid.UUID, tenantID *uuid.UUID) error
}
package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/user/domain"
)

// UserRepository define el contrato para la persistencia de usuarios
// Esta interfaz será implementada por el adaptador de PostgreSQL
// Sigue el patrón Repository para aislar el dominio de la infraestructura
type UserRepository interface {
	// Create crea un nuevo usuario en la base de datos control
	Create(ctx context.Context, user *domain.User) error
	
	// GetByID obtiene un usuario por su ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	
	// GetByEmail obtiene un usuario por su email (único)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	
	// GetByUsername obtiene un usuario por su username (único)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	
	// GetByEmailToken obtiene un usuario por su token de verificación de email
	GetByEmailToken(ctx context.Context, token string) (*domain.User, error)
	
	// Update actualiza un usuario existente
	Update(ctx context.Context, user *domain.User) error
	
	// Delete elimina lógicamente un usuario (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error
	
	// ExistsByEmail verifica si existe un usuario con el email dado
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	
	// ExistsByUsername verifica si existe un usuario con el username dado
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	
	// GetUserTenants obtiene todos los tenants asociados a un usuario
	GetUserTenants(ctx context.Context, userID uuid.UUID) ([]*domain.TenantUser, error)
	
	// GetTenantsByUser obtiene todos los tenant_users por userID
	GetTenantsByUser(ctx context.Context, userID uuid.UUID) ([]*domain.TenantUser, error)
	
	// UserHasAccessToTenant verifica si un usuario tiene acceso a un tenant
	UserHasAccessToTenant(ctx context.Context, userID, tenantID uuid.UUID) (bool, error)
	
	// AddUserToTenant asocia un usuario a un tenant
	AddUserToTenant(ctx context.Context, tenantUser *domain.TenantUser) error
	
	// RemoveUserFromTenant desasocia un usuario de un tenant (soft delete)
	RemoveUserFromTenant(ctx context.Context, userID, tenantID uuid.UUID) error
}

// UserService define el contrato para el servicio de usuarios
// Esta interfaz encapsula la lógica de negocio relacionada con usuarios
// y sigue los principios de Clean Architecture
type UserService interface {
	// CreateUser crea un nuevo usuario
	CreateUser(ctx context.Context, user *domain.User) error
	
	// GetUserByID obtiene un usuario por su ID
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	
	// GetUserByEmail obtiene un usuario por su email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	
	// GetUserByEmailToken obtiene un usuario por su token de verificación de email
	GetUserByEmailToken(ctx context.Context, token string) (*domain.User, error)
	
	// UpdateUser actualiza un usuario existente
	UpdateUser(ctx context.Context, user *domain.User) error
	
	// ExistsByEmail verifica si existe un usuario con el email dado
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	
	// GetTenantsByUser obtiene todos los tenant_users por userID
	GetTenantsByUser(ctx context.Context, userID uuid.UUID) ([]*domain.TenantUser, error)
	
	// UserHasAccessToTenant verifica si un usuario tiene acceso a un tenant
	UserHasAccessToTenant(ctx context.Context, userID, tenantID uuid.UUID) (bool, error)
	
	// AddUserToTenant asocia un usuario a un tenant
	AddUserToTenant(ctx context.Context, tenantUser *domain.TenantUser) error
}
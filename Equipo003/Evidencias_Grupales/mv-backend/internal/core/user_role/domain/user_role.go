package domain

import (
	"time"
	"github.com/google/uuid"
)

// UserRole representa la relación entre usuario y rol en el dominio
// Mapea directamente con la tabla 'user_roles' en las bases de datos
type UserRole struct {
	ID       uuid.UUID  `json:"id"`
	UserID   uuid.UUID  `json:"user_id"`
	RoleID   uuid.UUID  `json:"role_id"`
	TenantID *uuid.UUID `json:"tenant_id,omitempty"` // Para roles específicos de tenant
	Created  time.Time  `json:"created"`
	Updated  time.Time  `json:"updated"`
}

// CreateUserRoleDto representa los datos necesarios para crear una relación usuario-rol
type CreateUserRoleDto struct {
	UserID   uuid.UUID  `json:"user_id" validate:"required,uuid"`
	RoleID   uuid.UUID  `json:"role_id" validate:"required,uuid"`
	TenantID *uuid.UUID `json:"tenant_id,omitempty"`
}

// UserRoleResponseDto representa la respuesta al crear una relación usuario-rol
type UserRoleResponseDto struct {
	UserRole *UserRole `json:"user_role"`
	Message  string    `json:"message"`
}

// SyncRoleUsersDto representa los datos para sincronizar usuarios a un rol
type SyncRoleUsersDto struct {
	RoleID   uuid.UUID   `json:"role_id" validate:"required,uuid"`
	UserIDs  []uuid.UUID `json:"user_ids" validate:"required"`
	TenantID *uuid.UUID  `json:"tenant_id,omitempty"`
}

// SyncUserRolesDto representa los datos para sincronizar roles a un usuario
type SyncUserRolesDto struct {
	UserID   uuid.UUID   `json:"user_id" validate:"required,uuid"`
	RoleIDs  []uuid.UUID `json:"role_ids" validate:"required"`
	TenantID *uuid.UUID  `json:"tenant_id,omitempty"`
}

// NewUserRole crea una nueva instancia de UserRole con valores por defecto
func NewUserRole(userID, roleID uuid.UUID, tenantID *uuid.UUID) *UserRole {
	now := time.Now()

	return &UserRole{
		ID:       uuid.New(),
		UserID:   userID,
		RoleID:   roleID,
		TenantID: tenantID,
		Created:  now,
		Updated:  now,
	}
}
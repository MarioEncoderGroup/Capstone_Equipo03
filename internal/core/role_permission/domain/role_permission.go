package domain

import (
	"time"
	"github.com/google/uuid"
)

// RolePermission representa la relación entre rol y permiso en el dominio
// Mapea directamente con la tabla 'role_permissions' en las bases de datos
type RolePermission struct {
	ID           uuid.UUID `json:"id"`
	RoleID       uuid.UUID `json:"role_id"`
	PermissionID uuid.UUID `json:"permission_id"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}

// CreateRolePermissionDto representa los datos necesarios para crear una relación rol-permiso
type CreateRolePermissionDto struct {
	RoleID       uuid.UUID `json:"role_id" validate:"required,uuid"`
	PermissionID uuid.UUID `json:"permission_id" validate:"required,uuid"`
}

// RolePermissionResponseDto representa la respuesta al crear una relación rol-permiso
type RolePermissionResponseDto struct {
	RolePermission *RolePermission `json:"role_permission"`
	Message        string          `json:"message"`
}

// SyncRolePermissionsDto representa los datos para sincronizar permisos a un rol
type SyncRolePermissionsDto struct {
	RoleID        uuid.UUID   `json:"role_id" validate:"required,uuid"`
	PermissionIDs []uuid.UUID `json:"permission_ids" validate:"required"`
}

// Permission representa un permiso con información detallada
// Se utiliza para obtener información completa de permisos asociados a roles
type Permission struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Section     string    `json:"section"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

// NewRolePermission crea una nueva instancia de RolePermission con valores por defecto
func NewRolePermission(roleID, permissionID uuid.UUID) *RolePermission {
	now := time.Now()

	return &RolePermission{
		ID:           uuid.New(),
		RoleID:       roleID,
		PermissionID: permissionID,
		Created:      now,
		Updated:      now,
	}
}
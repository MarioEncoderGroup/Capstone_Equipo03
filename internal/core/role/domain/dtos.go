package domain

import "github.com/google/uuid"

// CreateRoleRequest DTO para crear un nuevo rol
type CreateRoleRequest struct {
	Name        string     `json:"name" validate:"required,min=3,max=50"`
	Description string     `json:"description" validate:"max=500"`
	TenantID    *uuid.UUID `json:"tenant_id,omitempty"`
}

// UpdateRoleRequest DTO para actualizar un rol existente
type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"max=500"`
}

// RoleResponse DTO para respuestas de rol
type RoleResponse struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    *uuid.UUID `json:"tenant_id,omitempty"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	IsGlobal    bool       `json:"is_global"`
	IsSystem    bool       `json:"is_system"`
	Created     string     `json:"created"`     // ISO 8601 format
	Updated     string     `json:"updated"`     // ISO 8601 format
}

// RoleListResponse DTO para listados de roles
type RoleListResponse struct {
	Roles []RoleResponse `json:"roles"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// RoleFilterRequest DTO para filtros de búsqueda
type RoleFilterRequest struct {
	TenantID *uuid.UUID `json:"tenant_id,omitempty"`
	Name     string     `json:"name,omitempty"`
	Page     int        `json:"page" validate:"min=1"`
	Limit    int        `json:"limit" validate:"min=1,max=100"`
}

// ToResponse convierte una entidad Role a RoleResponse
func (r *Role) ToResponse() *RoleResponse {
	return &RoleResponse{
		ID:          r.ID,
		TenantID:    r.TenantID,
		Name:        r.Name,
		Description: r.Description,
		IsGlobal:    r.IsGlobalRole(),
		IsSystem:    r.IsSystemRole(),
		Created:     r.Created.Format("2006-01-02T15:04:05Z07:00"),
		Updated:     r.Updated.Format("2006-01-02T15:04:05Z07:00"),
	}
}
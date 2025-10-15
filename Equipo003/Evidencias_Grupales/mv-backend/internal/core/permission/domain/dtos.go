package domain

import "github.com/google/uuid"

// CreatePermissionRequest DTO para crear un nuevo permiso
type CreatePermissionRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"max=500"`
	Section     string `json:"section" validate:"required,min=3,max=50"`
}

// UpdatePermissionRequest DTO para actualizar un permiso existente
type UpdatePermissionRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"max=500"`
	Section     string `json:"section" validate:"required,min=3,max=50"`
}

// PermissionResponse DTO para respuestas de permiso
type PermissionResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Section     string    `json:"section"`
	IsSystem    bool      `json:"is_system"`
	IsActive    bool      `json:"is_active"`
	Created     string    `json:"created"`     // ISO 8601 format
	Updated     string    `json:"updated"`     // ISO 8601 format
}

// PermissionListResponse DTO para listados de permisos
type PermissionListResponse struct {
	Permissions []PermissionResponse `json:"permissions"`
	Total       int                  `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
}

// PermissionFilterRequest DTO para filtros de búsqueda
type PermissionFilterRequest struct {
	Name    string `json:"name,omitempty"`
	Section string `json:"section,omitempty"`
	Search  string `json:"search,omitempty"`
	SortBy  string `json:"sort_by,omitempty"`
	SortDir string `json:"sort_dir,omitempty"`
	Page    int    `json:"page" validate:"min=1"`
	Limit   int    `json:"limit" validate:"min=1,max=100"`
}

// PermissionGroupedResponse DTO para permisos agrupados por sección
type PermissionGroupedResponse struct {
	Section     string               `json:"section"`
	Permissions []PermissionResponse `json:"permissions"`
}

// PermissionSectionResponse DTO para respuesta de secciones disponibles
type PermissionSectionResponse struct {
	Sections []string `json:"sections"`
}

// ToResponse convierte una entidad Permission a PermissionResponse
func (p *Permission) ToResponse() *PermissionResponse {
	return &PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Section:     p.Section,
		IsSystem:    p.IsSystemPermission(),
		IsActive:    p.IsActive(),
		Created:     p.Created.Format("2006-01-02T15:04:05Z07:00"),
		Updated:     p.Updated.Format("2006-01-02T15:04:05Z07:00"),
	}
}
package domain

import (
	"time"
	"github.com/google/uuid"
)

// Role representa la entidad principal de rol en el dominio
// Mapea directamente con la tabla 'roles' en las bases de datos (control y tenant)
type Role struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    *uuid.UUID `json:"tenant_id"`     // NULL para roles globales
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// RoleType representa los tipos de roles del sistema
type RoleType string

const (
	RoleTypeGlobal RoleType = "global" // Roles del sistema (administrator)
	RoleTypeTenant RoleType = "tenant" // Roles específicos de tenant
)

// Roles predefinidos del sistema
const (
	RoleNameAdministrator    = "administrator"
	RoleNameApprover        = "approver"
	RoleNameExpenseSubmitter = "expense-submitter"
)

// NewRole crea una nueva instancia de rol con valores por defecto
// Para roles globales, tenantID debe ser nil
func NewRole(name, description string, tenantID *uuid.UUID) *Role {
	now := time.Now()

	// Validar descripción como puntero
	var descPtr *string
	if description != "" {
		descPtr = &description
	}

	return &Role{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        name,
		Description: descPtr,
		Created:     now,
		Updated:     now,
	}
}

// IsGlobalRole verifica si el rol es global (del sistema)
func (r *Role) IsGlobalRole() bool {
	return r.TenantID == nil
}

// IsTenantRole verifica si el rol pertenece a un tenant específico
func (r *Role) IsTenantRole() bool {
	return r.TenantID != nil
}

// IsSystemRole verifica si es uno de los roles predefinidos del sistema
func (r *Role) IsSystemRole() bool {
	systemRoles := []string{
		RoleNameAdministrator,
		RoleNameApprover,
		RoleNameExpenseSubmitter,
	}

	for _, sysRole := range systemRoles {
		if r.Name == sysRole {
			return true
		}
	}
	return false
}

// Update actualiza los campos modificables del rol
func (r *Role) Update(name, description string) {
	r.Name = name
	if description != "" {
		r.Description = &description
	} else {
		r.Description = nil
	}
	r.Updated = time.Now()
}

// SoftDelete marca el rol como eliminado
func (r *Role) SoftDelete() {
	now := time.Now()
	r.DeletedAt = &now
	r.Updated = now
}
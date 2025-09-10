package domain

import (
	"time"
	"github.com/google/uuid"
)

// TenantUser representa la relación many-to-many entre usuarios y tenants
// Mapea directamente con la tabla 'tenant_users' del control database
// Un usuario puede pertenecer a múltiples tenants (empresas)
type TenantUser struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Created   time.Time  `json:"created"`
	Updated   time.Time  `json:"updated"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// NewTenantUser crea una nueva relación entre usuario y tenant
func NewTenantUser(tenantID, userID uuid.UUID) *TenantUser {
	now := time.Now()
	return &TenantUser{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Created:  now,
		Updated:  now,
	}
}

// IsActive verifica si la relación está activa (no eliminada)
func (tu *TenantUser) IsActive() bool {
	return tu.DeletedAt == nil
}

// SoftDelete elimina lógicamente la relación
func (tu *TenantUser) SoftDelete() {
	now := time.Now()
	tu.DeletedAt = &now
	tu.Updated = now
}
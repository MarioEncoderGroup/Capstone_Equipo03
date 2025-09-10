package domain

import (
	"time"

	"github.com/google/uuid"
)

// TenantStatus representa los posibles estados de un tenant
// Mapea con el ENUM tenant_status de la BD
type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusInactive  TenantStatus = "inactive"
	TenantStatusSuspended TenantStatus = "suspended"
)

// Tenant representa una empresa/organización en el sistema multi-tenant
// Mapea directamente con la tabla 'tenants' del control database
type Tenant struct {
	ID           uuid.UUID    `json:"id"`
	RUT          string       `json:"rut"`           // RUT del negocio chileno
	BusinessName string       `json:"business_name"` // Nombre de la empresa
	Email        string       `json:"email"`
	Phone        string       `json:"phone"`
	Address      string       `json:"address"`
	Website      string       `json:"website"`
	Logo         *string      `json:"logo"`
	RegionID     string       `json:"region_id"`  // ID de región chilena (AP, TA, etc)
	CommuneID    string       `json:"commune_id"` // ID de comuna chilena
	CountryID    uuid.UUID    `json:"country_id"` // ID del país
	Status       TenantStatus `json:"status"`
	NodeNumber   int          `json:"node_number"` // Número de nodo para distribución
	Slug         string       `json:"slug"`        // Slug del tenant
	TenantName   string       `json:"tenant_name"` // Nombre de la base de datos del tenant
	CreatedBy    uuid.UUID    `json:"created_by"`  // Usuario que creó el tenant
	UpdatedBy    uuid.UUID    `json:"updated_by"`  // Usuario que actualizó el tenant
	Created      time.Time    `json:"created"`
	Updated      time.Time    `json:"updated"`
	DeletedAt    *time.Time   `json:"deleted_at"`
}

// NewTenant crea una nueva instancia de tenant
// Aplica reglas de negocio para crear un tenant válido
func NewTenant(rut, businessName, email, phone, address, website string, regionID, communeID string, countryID uuid.UUID, nodeNumber int, slug string, createdBy uuid.UUID) *Tenant {
	now := time.Now()
	tenantID := uuid.New()

	return &Tenant{
		ID:           tenantID,
		RUT:          rut,
		BusinessName: businessName,
		Email:        email,
		Phone:        phone,
		Address:      address,
		Website:      website,
		RegionID:     regionID,
		CommuneID:    communeID,
		CountryID:    countryID,
		Status:       TenantStatusActive,
		NodeNumber:   nodeNumber,
		Slug:         slug,
		TenantName:   generateTenantDBName(tenantID),
		CreatedBy:    createdBy,
		UpdatedBy:    createdBy,
		Created:      now,
		Updated:      now,
	}
}

// generateTenantDBName genera el nombre de la base de datos del tenant
// Sigue el patrón: misviaticos_tenant_{uuid_sin_guiones}
func generateTenantDBName(tenantID uuid.UUID) string {
	return "misviaticos_tenant_" + tenantID.String()[:8] // Primeros 8 caracteres del UUID
}

// IsActive verifica si el tenant está activo
func (t *Tenant) IsActive() bool {
	return t.Status == TenantStatusActive && t.DeletedAt == nil
}

// Suspend suspende el tenant
func (t *Tenant) Suspend(updatedBy uuid.UUID) {
	t.Status = TenantStatusSuspended
	t.UpdatedBy = updatedBy
	t.Updated = time.Now()
}

// Activate activa el tenant
func (t *Tenant) Activate(updatedBy uuid.UUID) {
	t.Status = TenantStatusActive
	t.UpdatedBy = updatedBy
	t.Updated = time.Now()
}

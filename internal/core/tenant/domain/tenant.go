package domain_tenant

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID           uuid.UUID  `json:"id,omitempty"`
	Rut          string     `json:"rut,omitempty"`
	BusinessName string     `json:"business_name,omitempty"`
	Email        string     `json:"email,omitempty"`
	Phone        string     `json:"phone,omitempty"`
	Address      string     `json:"address,omitempty"`
	Website      *string    `json:"website,omitempty"`
	Logo         *string    `json:"logo,omitempty"`
	RegionID     string     `json:"region_id,omitempty"`
	CommuneID    string     `json:"commune_id,omitempty"`
	CountryID    uuid.UUID  `json:"country_id,omitempty"`
	Status       string     `json:"status,omitempty"`
	NodeNumber   int        `json:"node_number,omitempty"`
	TenantName   string     `json:"tenant_name,omitempty"`
	CreatedBy    uuid.UUID  `json:"created_by,omitempty"`
	UpdatedBy    uuid.UUID  `json:"updated_by,omitempty"`
	Created      time.Time  `json:"created,omitempty"`
	Updated      time.Time  `json:"updated,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	Trial        time.Time  `json:"trial,omitempty"`
}

// TenantStatus representa los estados posibles de un tenant
type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusInactive  TenantStatus = "inactive"
	TenantStatusSuspended TenantStatus = "suspended"
)

// NewTenant crea una nueva instancia de Tenant
func NewTenant(rut, businessName, email, phone, address string, website *string, regionID, communeID string, countryID uuid.UUID, nodeNumber int, tenantName string, createdBy uuid.UUID) *Tenant {
	now := time.Now()
	return &Tenant{
		ID:           uuid.New(),
		Rut:          rut,
		BusinessName: businessName,
		Email:        email,
		Phone:        phone,
		Address:      address,
		Website:      website,
		RegionID:     regionID,
		CommuneID:    communeID,
		CountryID:    countryID,
		Status:       string(TenantStatusActive),
		NodeNumber:   nodeNumber,
		TenantName:   tenantName,
		CreatedBy:    createdBy,
		UpdatedBy:    createdBy,
		Created:      now,
		Updated:      now,
		Trial:        now.AddDate(0, 1, 0), // 1 mes de trial
	}
}

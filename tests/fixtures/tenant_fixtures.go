package fixtures

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/tenant/domain"
)

// TenantFixtures provides predefined tenant data for testing
type TenantFixtures struct{}

// NewTenantFixtures creates a new TenantFixtures instance
func NewTenantFixtures() *TenantFixtures {
	return &TenantFixtures{}
}

// ValidTenant returns a valid tenant for testing
func (f *TenantFixtures) ValidTenant(createdBy uuid.UUID) *domain.Tenant {
	return domain.NewTenant(
		"76.123.456-0",
		"Empresa de Prueba SpA",
		"contacto@empresa.cl",
		"+56987654321",
		"Av. Providencia 123, Santiago",
		"https://empresa.cl",
		"RM",
		"Santiago",
		uuid.New(), // country_id
		1,
		"empresa_prueba",
		createdBy,
	)
}

// TenantWithDifferentRUT returns a tenant with different RUT
func (f *TenantFixtures) TenantWithDifferentRUT(createdBy uuid.UUID) *domain.Tenant {
	return domain.NewTenant(
		"77.456.789-0",
		"Otra Empresa Ltda",
		"info@otraempresa.cl",
		"+56912345678",
		"Av. Las Condes 456, Santiago",
		"https://otraempresa.cl",
		"RM",
		"Las Condes",
		uuid.New(),
		2,
		"otra_empresa",
		createdBy,
	)
}

// ActiveTenant returns an active tenant
func (f *TenantFixtures) ActiveTenant(createdBy uuid.UUID) *domain.Tenant {
	tenant := f.ValidTenant(createdBy)
	tenant.Status = "active"
	return tenant
}

// InactiveTenant returns an inactive tenant
func (f *TenantFixtures) InactiveTenant(createdBy uuid.UUID) *domain.Tenant {
	tenant := f.ValidTenant(createdBy)
	tenant.Status = "inactive"
	return tenant
}

// SuspendedTenant returns a suspended tenant
func (f *TenantFixtures) SuspendedTenant(createdBy uuid.UUID) *domain.Tenant {
	tenant := f.ValidTenant(createdBy)
	tenant.Status = "suspended"
	return tenant
}

// TenantWithLogo returns a tenant with logo URL
func (f *TenantFixtures) TenantWithLogo(createdBy uuid.UUID) *domain.Tenant {
	tenant := f.ValidTenant(createdBy)
	logo := "https://empresa.cl/logo.png"
	tenant.Logo = &logo
	return tenant
}

// TenantsForBulkTesting returns multiple tenants for bulk operations
func (f *TenantFixtures) TenantsForBulkTesting(count int, createdBy uuid.UUID) []*domain.Tenant {
	tenants := make([]*domain.Tenant, count)
	for i := 0; i < count; i++ {
		// Generate valid RUT for each tenant
		rutNumber := 76123456 + i
		verificator := f.calculateRUTVerificator(rutNumber)
		
		tenants[i] = domain.NewTenant(
			fmt.Sprintf("%s-%s", f.formatRUT(rutNumber), verificator),
			fmt.Sprintf("Empresa %d SpA", i+1),
			fmt.Sprintf("contacto%d@empresa%d.cl", i+1, i+1),
			"+56987654321",
			fmt.Sprintf("Av. Providencia %d, Santiago", 100+i),
			fmt.Sprintf("https://empresa%d.cl", i+1),
			"RM",
			"Santiago",
			uuid.New(),
			i+1,
			fmt.Sprintf("empresa_%d", i+1),
			createdBy,
		)
	}
	return tenants
}

// ValidTenantRegisterRequestData returns request data for tenant registration
func (f *TenantFixtures) ValidTenantRegisterRequestData() map[string]interface{} {
	return map[string]interface{}{
		"username":               "adminuser",
		"full_name":              "Administrador Empresa",
		"email":                  "admin@empresa.cl",
		"password":               "password123",
		"phone":                  "+56912345678",
		"identification_number":  "12.345.678-5",
		"create_tenant":          true,
		"tenant_data": map[string]interface{}{
			"rut":          "76.123.456-0",
			"business_name": "Empresa de Prueba SpA",
			"email":        "contacto@empresa.cl",
			"phone":        "+56987654321",
			"address":      "Av. Providencia 123, Santiago",
			"website":      "https://empresa.cl",
			"region_id":    "RM",
			"commune_id":   "Santiago",
			"country_id":   "01234567-89ab-cdef-0123-456789abcdef",
		},
	}
}

// InvalidTenantRegisterRequestData returns invalid request data for testing
func (f *TenantFixtures) InvalidTenantRegisterRequestData() map[string]interface{} {
	return map[string]interface{}{
		"username":               "adminuser",
		"full_name":              "Administrador Empresa",
		"email":                  "admin@empresa.cl",
		"password":               "password123",
		"phone":                  "+56912345678",
		"identification_number":  "12.345.678-5",
		"create_tenant":          true,
		"tenant_data": map[string]interface{}{
			"rut":          "76.123.456-X", // Invalid RUT
			"business_name": "",             // Empty business name
			"email":        "invalid-email", // Invalid email
			"phone":        "invalid",       // Invalid phone
			"address":      "",              // Empty address
			"website":      "not-a-url",     // Invalid URL
			"region_id":    "",              // Empty region
			"commune_id":   "",              // Empty commune
			"country_id":   "invalid-uuid",  // Invalid UUID
		},
	}
}

// ChileanRegions returns common Chilean regions for testing
func (f *TenantFixtures) ChileanRegions() []map[string]string {
	return []map[string]string{
		{"id": "XV", "name": "Arica y Parinacota"},
		{"id": "I", "name": "Tarapacá"},
		{"id": "II", "name": "Antofagasta"},
		{"id": "III", "name": "Atacama"},
		{"id": "IV", "name": "Coquimbo"},
		{"id": "V", "name": "Valparaíso"},
		{"id": "RM", "name": "Metropolitana"},
		{"id": "VI", "name": "O'Higgins"},
		{"id": "VII", "name": "Maule"},
		{"id": "VIII", "name": "Biobío"},
		{"id": "IX", "name": "Araucanía"},
		{"id": "XIV", "name": "Los Ríos"},
		{"id": "X", "name": "Los Lagos"},
		{"id": "XI", "name": "Aysén"},
		{"id": "XII", "name": "Magallanes"},
	}
}

// Helper functions for RUT generation
func (f *TenantFixtures) formatRUT(rut int) string {
	rutStr := fmt.Sprintf("%d", rut)
	if len(rutStr) >= 7 {
		return fmt.Sprintf("%s.%s.%s",
			rutStr[:len(rutStr)-6],
			rutStr[len(rutStr)-6:len(rutStr)-3],
			rutStr[len(rutStr)-3:])
	}
	return rutStr
}

func (f *TenantFixtures) calculateRUTVerificator(rut int) string {
	sum := 0
	multiplier := 2
	
	for rut > 0 {
		sum += (rut % 10) * multiplier
		rut = rut / 10
		multiplier++
		if multiplier > 7 {
			multiplier = 2
		}
	}
	
	remainder := sum % 11
	verificator := 11 - remainder
	
	if verificator == 11 {
		return "0"
	} else if verificator == 10 {
		return "K"
	}
	
	return fmt.Sprintf("%d", verificator)
}
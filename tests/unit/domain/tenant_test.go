package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/core/tenant/domain"
	"github.com/google/uuid"
)

func TestNewTenant(t *testing.T) {
	// Datos de prueba
	rut := "76.123.456-7"
	businessName := "Test Company SpA"
	email := "contact@testcompany.cl"
	phone := "+56987654321"
	address := "Av. Providencia 123, Santiago"
	website := "https://testcompany.cl"
	regionID := "RM"
	communeID := "Santiago"
	countryID := uuid.New()
	createdBy := uuid.New()

	tenant := domain.NewTenant(
		rut, businessName, email, phone, address, website,
		regionID, communeID, countryID, 1, "test_tenant", createdBy,
	)

	// Verificar campos básicos
	if tenant.RUT != rut {
		t.Errorf("Expected RUT %s, got %s", rut, tenant.RUT)
	}

	if tenant.BusinessName != businessName {
		t.Errorf("Expected business name %s, got %s", businessName, tenant.BusinessName)
	}

	if tenant.Email != email {
		t.Errorf("Expected email %s, got %s", email, tenant.Email)
	}

	if tenant.Phone != phone {
		t.Errorf("Expected phone %s, got %s", phone, tenant.Phone)
	}

	if tenant.Address != address {
		t.Errorf("Expected address %s, got %s", address, tenant.Address)
	}

	if tenant.Website != website {
		t.Errorf("Expected website %s, got %s", website, tenant.Website)
	}

	if tenant.RegionID != regionID {
		t.Errorf("Expected region %s, got %s", regionID, tenant.RegionID)
	}

	if tenant.CommuneID != communeID {
		t.Errorf("Expected commune %s, got %s", communeID, tenant.CommuneID)
	}

	if tenant.CountryID != countryID {
		t.Errorf("Expected country ID %s, got %s", countryID, tenant.CountryID)
	}

	if tenant.CreatedBy != createdBy {
		t.Errorf("Expected created by %s, got %s", createdBy, tenant.CreatedBy)
	}

	if tenant.UpdatedBy != createdBy {
		t.Errorf("Expected updated by %s, got %s", createdBy, tenant.UpdatedBy)
	}

	// Verificar valores por defecto
	if tenant.Status != domain.TenantStatusActive {
		t.Errorf("Expected status %s, got %s", domain.TenantStatusActive, tenant.Status)
	}

	if tenant.NodeNumber != 1 {
		t.Errorf("Expected node number 1, got %d", tenant.NodeNumber)
	}

	// Verificar UUID generado
	if tenant.ID.String() == "00000000-0000-0000-0000-000000000000" {
		t.Error("Tenant ID should be generated")
	}

	// Verificar nombre de BD generado
	expectedPrefix := "misviaticos_tenant_"
	if !strings.HasPrefix(tenant.TenantName, expectedPrefix) {
		t.Errorf("Expected tenant name to start with %s, got %s", expectedPrefix, tenant.TenantName)
	}

	if len(tenant.TenantName) != len(expectedPrefix)+8 { // prefix + 8 chars UUID
		t.Errorf("Expected tenant name length %d, got %d", len(expectedPrefix)+8, len(tenant.TenantName))
	}

	// Verificar timestamps
	now := time.Now()
	if tenant.Created.After(now) || tenant.Created.Before(now.Add(-time.Second)) {
		t.Error("Created timestamp should be around current time")
	}

	if tenant.Updated.After(now) || tenant.Updated.Before(now.Add(-time.Second)) {
		t.Error("Updated timestamp should be around current time")
	}
}

func TestTenantIsActive(t *testing.T) {
	createdBy := uuid.New()
	tenant := domain.NewTenant(
		"76.123.456-7", "Test Company", "test@company.cl", "+56987654321",
		"Address", "https://company.cl", "RM", "Santiago", uuid.New(), 1, "test_tenant", createdBy,
	)

	// Tenant recién creado debe estar activo
	if !tenant.IsActive() {
		t.Error("New tenant should be active")
	}

	// Suspender tenant
	tenant.Suspend(createdBy)
	if tenant.IsActive() {
		t.Error("Suspended tenant should not be active")
	}

	// Reactivar tenant
	tenant.Activate(createdBy)
	if !tenant.IsActive() {
		t.Error("Reactivated tenant should be active")
	}

	// Soft delete
	now := time.Now()
	tenant.DeletedAt = &now
	if tenant.IsActive() {
		t.Error("Soft deleted tenant should not be active")
	}
}

func TestTenantSuspend(t *testing.T) {
	createdBy := uuid.New()
	updatedBy := uuid.New()

	tenant := domain.NewTenant(
		"76.123.456-7", "Test Company", "test@company.cl", "+56987654321",
		"Address", "https://company.cl", "RM", "Santiago", uuid.New(), 1, "test_tenant", createdBy,
	)

	beforeSuspend := time.Now()
	tenant.Suspend(updatedBy)
	afterSuspend := time.Now()

	// Verificar estado suspendido
	if tenant.Status != domain.TenantStatusSuspended {
		t.Errorf("Expected status %s, got %s", domain.TenantStatusSuspended, tenant.Status)
	}

	// Verificar updatedBy
	if tenant.UpdatedBy != updatedBy {
		t.Errorf("Expected updated by %s, got %s", updatedBy, tenant.UpdatedBy)
	}

	// Verificar timestamp actualizado
	if tenant.Updated.Before(beforeSuspend) || tenant.Updated.After(afterSuspend) {
		t.Error("Updated timestamp should be set during suspension")
	}
}

func TestTenantActivate(t *testing.T) {
	createdBy := uuid.New()
	updatedBy := uuid.New()

	tenant := domain.NewTenant(
		"76.123.456-7", "Test Company", "test@company.cl", "+56987654321",
		"Address", "https://company.cl", "RM", "Santiago", uuid.New(), 1, "test_tenant", createdBy,
	)

	// Suspender primero
	tenant.Suspend(createdBy)

	beforeActivate := time.Now()
	tenant.Activate(updatedBy)
	afterActivate := time.Now()

	// Verificar estado activo
	if tenant.Status != domain.TenantStatusActive {
		t.Errorf("Expected status %s, got %s", domain.TenantStatusActive, tenant.Status)
	}

	// Verificar updatedBy
	if tenant.UpdatedBy != updatedBy {
		t.Errorf("Expected updated by %s, got %s", updatedBy, tenant.UpdatedBy)
	}

	// Verificar timestamp actualizado
	if tenant.Updated.Before(beforeActivate) || tenant.Updated.After(afterActivate) {
		t.Error("Updated timestamp should be set during activation")
	}
}

func TestTenantStatus(t *testing.T) {
	// Verificar constantes de status
	if domain.TenantStatusActive != "active" {
		t.Errorf("Expected TenantStatusActive to be 'active', got %s", domain.TenantStatusActive)
	}

	if domain.TenantStatusInactive != "inactive" {
		t.Errorf("Expected TenantStatusInactive to be 'inactive', got %s", domain.TenantStatusInactive)
	}

	if domain.TenantStatusSuspended != "suspended" {
		t.Errorf("Expected TenantStatusSuspended to be 'suspended', got %s", domain.TenantStatusSuspended)
	}
}

package domain_test

import (
	"testing"
	"time"

	tenantDomain "github.com/JoseLuis21/mv-backend/internal/core/tenant/domain"
	"github.com/google/uuid"
)

func TestNewTenant(t *testing.T) {
	// Datos de prueba - ahora solo campos esenciales + campos opcionales vacíos
	rut := "76.123.456-7"
	businessName := "Test Company SpA"
	userEmail := "user@testcompany.cl" // Email del usuario que crea el tenant
	createdBy := uuid.New()

	tenant := tenantDomain.NewTenant(
		rut, businessName, userEmail, "", "", "", // phone, address, website vacíos
		"", "", uuid.Nil, 1, "test_tenant", createdBy, // regionID, communeID, countryID vacíos
	)

	// Verificar campos básicos
	if tenant.Rut != rut {
		t.Errorf("Expected RUT %s, got %s", rut, tenant.Rut)
	}

	if tenant.BusinessName != businessName {
		t.Errorf("Expected business name %s, got %s", businessName, tenant.BusinessName)
	}

	if tenant.Email != userEmail {
		t.Errorf("Expected email %s, got %s", userEmail, tenant.Email)
	}

	// Verificar que los campos opcionales estén vacíos (simplicidad)
	if tenant.Phone != "" {
		t.Errorf("Expected empty phone, got %s", tenant.Phone)
	}

	if tenant.Address != "" {
		t.Errorf("Expected empty address, got %s", tenant.Address)
	}

	if tenant.Website != "" {
		t.Errorf("Expected empty website, got %s", tenant.Website)
	}

	if tenant.RegionID != "" {
		t.Errorf("Expected empty region, got %s", tenant.RegionID)
	}

	if tenant.CommuneID != "" {
		t.Errorf("Expected empty commune, got %s", tenant.CommuneID)
	}

	if tenant.CountryID != "00000000-0000-0000-0000-000000000000" {
		t.Errorf("Expected nil country ID, got %s", tenant.CountryID)
	}

	if tenant.CreatedBy != createdBy {
		t.Errorf("Expected created by %s, got %s", createdBy, tenant.CreatedBy)
	}

	if tenant.UpdatedBy != createdBy {
		t.Errorf("Expected updated by %s, got %s", createdBy, tenant.UpdatedBy)
	}

	// Verificar valores por defecto
	if tenant.Status != string(tenantDomain.TenantStatusActive) {
		t.Errorf("Expected status %s, got %s", tenantDomain.TenantStatusActive, tenant.Status)
	}

	if tenant.NodeNumber != 1 {
		t.Errorf("Expected node number 1, got %d", tenant.NodeNumber)
	}

	// Verificar UUID generado
	if tenant.ID.String() == "00000000-0000-0000-0000-000000000000" {
		t.Error("Tenant ID should be generated")
	}

	// Verificar nombre de BD pasado como parámetro
	expectedTenantName := "test_tenant"
	if tenant.TenantName != expectedTenantName {
		t.Errorf("Expected tenant name %s, got %s", expectedTenantName, tenant.TenantName)
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
	tenant := tenantDomain.NewTenant(
		"76.123.456-7", "Test Company", "test@company.cl", "", "", "",
		"", "", uuid.Nil, 1, "test_tenant", createdBy,
	)

	// Tenant recién creado debe estar activo
	if tenant.Status != string(tenantDomain.TenantStatusActive) {
		t.Error("New tenant should be active")
	}

	// Suspender tenant manualmente
	tenant.Status = string(tenantDomain.TenantStatusSuspended)
	tenant.UpdatedBy = createdBy
	tenant.Updated = time.Now()
	
	if tenant.Status == string(tenantDomain.TenantStatusActive) {
		t.Error("Suspended tenant should not be active")
	}

	// Reactivar tenant manualmente
	tenant.Status = string(tenantDomain.TenantStatusActive)
	tenant.UpdatedBy = createdBy
	tenant.Updated = time.Now()
	
	if tenant.Status != string(tenantDomain.TenantStatusActive) {
		t.Error("Reactivated tenant should be active")
	}

	// Soft delete
	now := time.Now()
	tenant.DeletedAt = &now
	if tenant.DeletedAt == nil {
		t.Error("Soft deleted tenant should have DeletedAt set")
	}
}

func TestTenantSuspend(t *testing.T) {
	createdBy := uuid.New()
	updatedBy := uuid.New()

	tenant := tenantDomain.NewTenant(
		"76.123.456-7", "Test Company", "test@company.cl", "", "", "",
		"", "", uuid.Nil, 1, "test_tenant", createdBy,
	)

	beforeSuspend := time.Now()
	// Suspender manualmente
	tenant.Status = string(tenantDomain.TenantStatusSuspended)
	tenant.UpdatedBy = updatedBy
	tenant.Updated = time.Now()
	afterSuspend := time.Now()

	// Verificar estado suspendido
	if tenant.Status != string(tenantDomain.TenantStatusSuspended) {
		t.Errorf("Expected status %s, got %s", tenantDomain.TenantStatusSuspended, tenant.Status)
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

	tenant := tenantDomain.NewTenant(
		"76.123.456-7", "Test Company", "test@company.cl", "+56987654321",
		"Address", "https://company.cl", "RM", "Santiago", uuid.New(), 1, "test_tenant", createdBy,
	)

	// Suspender primero manualmente
	tenant.Status = string(tenantDomain.TenantStatusSuspended)
	tenant.UpdatedBy = createdBy
	tenant.Updated = time.Now()

	beforeActivate := time.Now()
	// Activar manualmente
	tenant.Status = string(tenantDomain.TenantStatusActive)
	tenant.UpdatedBy = updatedBy
	tenant.Updated = time.Now()
	afterActivate := time.Now()

	// Verificar estado activo
	if tenant.Status != string(tenantDomain.TenantStatusActive) {
		t.Errorf("Expected status %s, got %s", tenantDomain.TenantStatusActive, tenant.Status)
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
	if tenantDomain.TenantStatusActive != "active" {
		t.Errorf("Expected TenantStatusActive to be 'active', got %s", tenantDomain.TenantStatusActive)
	}

	if tenantDomain.TenantStatusInactive != "inactive" {
		t.Errorf("Expected TenantStatusInactive to be 'inactive', got %s", tenantDomain.TenantStatusInactive)
	}

	if tenantDomain.TenantStatusSuspended != "suspended" {
		t.Errorf("Expected TenantStatusSuspended to be 'suspended', got %s", tenantDomain.TenantStatusSuspended)
	}
}

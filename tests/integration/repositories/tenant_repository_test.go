package repositories_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/tenant/adapters"
	"github.com/JoseLuis21/mv-backend/internal/core/tenant/domain"
	"github.com/JoseLuis21/mv-backend/tests/helpers"
)

func TestTenantRepository_Create(t *testing.T) {
	// Setup
	client, _, cleanup := helpers.CreateTestDatabase(t)
	defer cleanup()
	helpers.SetupTestTables(t, client)
	
	repo := adapters.NewPostgreSQLTenantRepository(client)
	ctx := context.Background()
	
	// Create test user for created_by and updated_by
	userID := uuid.New()
	
	// Test data
	tenant := domain.NewTenant(
		"76.123.456-7",
		"Empresa de Prueba SpA",
		"contacto@empresa.cl",
		"+56987654321",
		"Av. Providencia 123, Santiago",
		"https://empresa.cl",
		"RM",
		"Santiago",
		uuid.New(), // country_id
		123, // node_number
		"empresa_prueba",
		userID,
	)
	
	// Execute
	err := repo.Create(ctx, tenant)
	
	// Assert
	if err != nil {
		t.Fatalf("Failed to create tenant: %v", err)
	}
	
	// Verify tenant was created
	retrieved, err := repo.GetByID(ctx, tenant.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve created tenant: %v", err)
	}
	
	if retrieved.BusinessName != tenant.BusinessName {
		t.Errorf("Expected business name %s, got %s", tenant.BusinessName, retrieved.BusinessName)
	}
	
	if retrieved.RUT != tenant.RUT {
		t.Errorf("Expected RUT %s, got %s", tenant.RUT, retrieved.RUT)
	}
	
	if retrieved.Email != tenant.Email {
		t.Errorf("Expected email %s, got %s", tenant.Email, retrieved.Email)
	}
}

func TestTenantRepository_GetByRUT(t *testing.T) {
	// Setup
	client, _, cleanup := helpers.CreateTestDatabase(t)
	defer cleanup()
	helpers.SetupTestTables(t, client)
	
	repo := adapters.NewPostgreSQLTenantRepository(client)
	ctx := context.Background()
	
	userID := uuid.New()
	
	// Create test tenant
	tenant := domain.NewTenant(
		"76.123.456-7",
		"Empresa de Prueba SpA",
		"contacto@empresa.cl",
		"+56987654321",
		"Av. Providencia 123, Santiago",
		"https://empresa.cl",
		"RM",
		"Santiago",
		uuid.New(),
		123,
		"empresa_prueba",
		userID,
	)
	
	err := repo.Create(ctx, tenant)
	if err != nil {
		t.Fatalf("Failed to create test tenant: %v", err)
	}
	
	// Test - Find existing tenant
	retrieved, err := repo.GetByRUT(ctx, tenant.RUT)
	if err != nil {
		t.Fatalf("Failed to get tenant by RUT: %v", err)
	}
	
	if retrieved.ID != tenant.ID {
		t.Errorf("Expected ID %s, got %s", tenant.ID, retrieved.ID)
	}
	
	// Test - Non-existing tenant
	_, err = repo.GetByRUT(ctx, "99.999.999-9")
	if err == nil {
		t.Error("Expected error for non-existent RUT, got nil")
	}
}

func TestTenantRepository_ExistsByRUT(t *testing.T) {
	// Setup
	client, _, cleanup := helpers.CreateTestDatabase(t)
	defer cleanup()
	helpers.SetupTestTables(t, client)
	
	repo := adapters.NewPostgreSQLTenantRepository(client)
	ctx := context.Background()
	
	rut := "76.123.456-7"
	
	// Test - RUT doesn't exist initially
	exists, err := repo.ExistsByRUT(ctx, rut)
	if err != nil {
		t.Fatalf("Failed to check RUT existence: %v", err)
	}
	if exists {
		t.Error("RUT should not exist initially")
	}
	
	// Create tenant
	userID := uuid.New()
	tenant := domain.NewTenant(
		rut,
		"Empresa de Prueba SpA",
		"contacto@empresa.cl",
		"+56987654321",
		"Av. Providencia 123, Santiago",
		"https://empresa.cl",
		"RM",
		"Santiago",
		uuid.New(),
		123,
		"empresa_prueba",
		userID,
	)
	
	err = repo.Create(ctx, tenant)
	if err != nil {
		t.Fatalf("Failed to create tenant: %v", err)
	}
	
	// Test - RUT exists now
	exists, err = repo.ExistsByRUT(ctx, rut)
	if err != nil {
		t.Fatalf("Failed to check RUT existence: %v", err)
	}
	if !exists {
		t.Error("RUT should exist after creation")
	}
}

func TestTenantRepository_Update(t *testing.T) {
	// Setup
	client, _, cleanup := helpers.CreateTestDatabase(t)
	defer cleanup()
	helpers.SetupTestTables(t, client)
	
	repo := adapters.NewPostgreSQLTenantRepository(client)
	ctx := context.Background()
	
	userID := uuid.New()
	
	// Create tenant
	tenant := domain.NewTenant(
		"76.123.456-7",
		"Empresa de Prueba SpA",
		"contacto@empresa.cl",
		"+56987654321",
		"Av. Providencia 123, Santiago",
		"https://empresa.cl",
		"RM",
		"Santiago",
		uuid.New(),
		123,
		"empresa_prueba",
		userID,
	)
	
	err := repo.Create(ctx, tenant)
	if err != nil {
		t.Fatalf("Failed to create tenant: %v", err)
	}
	
	// Update tenant
	tenant.BusinessName = "Empresa Actualizada SpA"
	tenant.Email = "nuevo@empresa.cl"
	tenant.Status = "active"
	tenant.UpdatedBy = userID
	
	err = repo.Update(ctx, tenant)
	if err != nil {
		t.Fatalf("Failed to update tenant: %v", err)
	}
	
	// Verify update
	retrieved, err := repo.GetByID(ctx, tenant.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated tenant: %v", err)
	}
	
	if retrieved.BusinessName != "Empresa Actualizada SpA" {
		t.Errorf("Expected business name 'Empresa Actualizada SpA', got %s", retrieved.BusinessName)
	}
	
	if retrieved.Email != "nuevo@empresa.cl" {
		t.Errorf("Expected email 'nuevo@empresa.cl', got %s", retrieved.Email)
	}
	
	if retrieved.Status != "active" {
		t.Errorf("Expected status 'active', got %s", retrieved.Status)
	}
}

func TestTenantRepository_SoftDelete(t *testing.T) {
	// Setup
	client, _, cleanup := helpers.CreateTestDatabase(t)
	defer cleanup()
	helpers.SetupTestTables(t, client)
	
	repo := adapters.NewPostgreSQLTenantRepository(client)
	ctx := context.Background()
	
	userID := uuid.New()
	
	// Create tenant
	tenant := domain.NewTenant(
		"76.123.456-7",
		"Empresa de Prueba SpA",
		"contacto@empresa.cl",
		"+56987654321",
		"Av. Providencia 123, Santiago",
		"https://empresa.cl",
		"RM",
		"Santiago",
		uuid.New(),
		123,
		"empresa_prueba",
		userID,
	)
	
	err := repo.Create(ctx, tenant)
	if err != nil {
		t.Fatalf("Failed to create tenant: %v", err)
	}
	
	// Soft delete
	err = repo.Delete(ctx, tenant.ID)
	if err != nil {
		t.Fatalf("Failed to soft delete tenant: %v", err)
	}
	
	// Verify tenant is not found (soft deleted)
	_, err = repo.GetByID(ctx, tenant.ID)
	if err == nil {
		t.Error("Expected error when getting soft deleted tenant")
	}
	
	// Verify tenant doesn't exist in searches
	exists, err := repo.ExistsByRUT(ctx, tenant.RUT)
	if err != nil {
		t.Fatalf("Failed to check RUT existence: %v", err)
	}
	if exists {
		t.Error("Soft deleted tenant should not exist in RUT check")
	}
}

func TestTenantRepository_GetNextNodeNumber(t *testing.T) {
	// Setup
	client, _, cleanup := helpers.CreateTestDatabase(t)
	defer cleanup()
	helpers.SetupTestTables(t, client)
	
	repo := adapters.NewPostgreSQLTenantRepository(client)
	ctx := context.Background()
	
	// Test initial node number (should start at 1)
	nodeNumber, err := repo.GetNextNodeNumber(ctx)
	if err != nil {
		t.Fatalf("Failed to get next node number: %v", err)
	}
	
	if nodeNumber != 1 {
		t.Errorf("Expected initial node number 1, got %d", nodeNumber)
	}
	
	// Create a tenant
	userID := uuid.New()
	tenant := domain.NewTenant(
		"76.123.456-7",
		"Empresa de Prueba SpA",
		"contacto@empresa.cl",
		"+56987654321",
		"Av. Providencia 123, Santiago",
		"https://empresa.cl",
		"RM",
		"Santiago",
		uuid.New(),
		nodeNumber,
		"empresa_prueba",
		userID,
	)
	
	err = repo.Create(ctx, tenant)
	if err != nil {
		t.Fatalf("Failed to create tenant: %v", err)
	}
	
	// Next node number should be 2
	nextNodeNumber, err := repo.GetNextNodeNumber(ctx)
	if err != nil {
		t.Fatalf("Failed to get next node number: %v", err)
	}
	
	if nextNodeNumber != 2 {
		t.Errorf("Expected next node number 2, got %d", nextNodeNumber)
	}
}

// Benchmarks
func BenchmarkTenantRepository_Create(b *testing.B) {
	if !helpers.CheckPostgreSQLAvailable() {
		b.Skip("PostgreSQL not available")
	}
	
	client, _, cleanup := helpers.CreateTestDatabase(&testing.T{})
	defer cleanup()
	helpers.SetupTestTables(&testing.T{}, client)
	
	repo := adapters.NewPostgreSQLTenantRepository(client)
	ctx := context.Background()
	userID := uuid.New()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tenant := domain.NewTenant(
			fmt.Sprintf("76.123.%03d-%d", i%900+100, (i%9)+1),
			fmt.Sprintf("Empresa %d SpA", i),
			fmt.Sprintf("contacto%d@empresa.cl", i),
			"+56987654321",
			"Av. Providencia 123, Santiago",
			fmt.Sprintf("https://empresa%d.cl", i),
			"RM",
			"Santiago",
			uuid.New(),
			i+1,
			fmt.Sprintf("empresa_%d", i),
			userID,
		)
		repo.Create(ctx, tenant)
	}
}

func BenchmarkTenantRepository_GetByRUT(b *testing.B) {
	if !helpers.CheckPostgreSQLAvailable() {
		b.Skip("PostgreSQL not available")
	}
	
	client, _, cleanup := helpers.CreateTestDatabase(&testing.T{})
	defer cleanup()
	helpers.SetupTestTables(&testing.T{}, client)
	
	repo := adapters.NewPostgreSQLTenantRepository(client)
	ctx := context.Background()
	userID := uuid.New()
	
	// Create test tenant
	tenant := domain.NewTenant(
		"76.123.456-7",
		"Empresa de Prueba SpA",
		"contacto@empresa.cl",
		"+56987654321",
		"Av. Providencia 123, Santiago",
		"https://empresa.cl",
		"RM",
		"Santiago",
		uuid.New(),
		123,
		"empresa_prueba",
		userID,
	)
	repo.Create(ctx, tenant)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.GetByRUT(ctx, "76.123.456-7")
	}
}
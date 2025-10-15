package repositories_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/user/adapters"
	"github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	"github.com/JoseLuis21/mv-backend/tests/helpers"
)

func TestUserRepository_Create(t *testing.T) {
	// Setup
	client, _, cleanup := helpers.CreateTestDatabase(t)
	defer cleanup()
	helpers.SetupTestTables(t, client)
	
	repo := adapters.NewPostgreSQLUserRepository(client)
	ctx := context.Background()
	
	// Test data
	user := domain.NewUser("Juan", "Pérez", "test@example.cl", "+56987654321", "hashedpassword")
	
	// Execute
	err := repo.Create(ctx, user)
	
	// Assert
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Verify user was created
	retrieved, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve created user: %v", err)
	}
	
	if retrieved.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrieved.Username)
	}
	
	if retrieved.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrieved.Email)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	// Setup
	client, _, cleanup := helpers.CreateTestDatabase(t)
	defer cleanup()
	helpers.SetupTestTables(t, client)
	
	repo := adapters.NewPostgreSQLUserRepository(client)
	ctx := context.Background()
	
	// Create test user
	user := domain.NewUser("Juan", "Pérez", "test@example.cl", "+56987654321", "hashedpassword")
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	
	// Test - Find existing user
	retrieved, err := repo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}
	
	if retrieved.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, retrieved.ID)
	}
	
	// Test - Non-existing user
	_, err = repo.GetByEmail(ctx, "nonexistent@example.cl")
	if err == nil {
		t.Error("Expected error for non-existent email, got nil")
	}
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	// Setup
	client, _, cleanup := helpers.CreateTestDatabase(t)
	defer cleanup()
	helpers.SetupTestTables(t, client)
	
	repo := adapters.NewPostgreSQLUserRepository(client)
	ctx := context.Background()
	
	email := "test@example.cl"
	
	// Test - Email doesn't exist initially
	exists, err := repo.ExistsByEmail(ctx, email)
	if err != nil {
		t.Fatalf("Failed to check email existence: %v", err)
	}
	if exists {
		t.Error("Email should not exist initially")
	}
	
	// Create user
	user := domain.NewUser("Juan", "Pérez", email, "+56987654321", "hashedpassword")
	err = repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Test - Email exists now
	exists, err = repo.ExistsByEmail(ctx, email)
	if err != nil {
		t.Fatalf("Failed to check email existence: %v", err)
	}
	if !exists {
		t.Error("Email should exist after creation")
	}
}

func TestUserRepository_Update(t *testing.T) {
	// Setup
	client, _, cleanup := helpers.CreateTestDatabase(t)
	defer cleanup()
	helpers.SetupTestTables(t, client)
	
	repo := adapters.NewPostgreSQLUserRepository(client)
	ctx := context.Background()
	
	// Create user
	user := domain.NewUser("Juan", "Pérez", "test@example.cl", "+56987654321", "hashedpassword")
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Update user
	user.FullName = "Updated Full Name"
	user.IsActive = true
	user.EmailVerified = true
	
	err = repo.Update(ctx, user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	
	// Verify update
	retrieved, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated user: %v", err)
	}
	
	if retrieved.FullName != "Updated Full Name" {
		t.Errorf("Expected full name 'Updated Full Name', got %s", retrieved.FullName)
	}
	
	if !retrieved.IsActive {
		t.Error("User should be active after update")
	}
	
	if !retrieved.EmailVerified {
		t.Error("User should have verified email after update")
	}
}

func TestUserRepository_EmailTokenOperations(t *testing.T) {
	// Setup
	client, _, cleanup := helpers.CreateTestDatabase(t)
	defer cleanup()
	helpers.SetupTestTables(t, client)
	
	repo := adapters.NewPostgreSQLUserRepository(client)
	ctx := context.Background()
	
	// Create user with email token
	user := domain.NewUser("Juan", "Pérez", "test@example.cl", "+56987654321", "hashedpassword")
	token := "verification_token_123"
	expiry := time.Now().Add(24 * time.Hour)
	user.EmailToken = &token
	user.EmailTokenExpires = &expiry
	
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Test GetByEmailToken
	retrieved, err := repo.GetByEmailToken(ctx, token)
	if err != nil {
		t.Fatalf("Failed to get user by email token: %v", err)
	}
	
	if retrieved.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, retrieved.ID)
	}
	
	if retrieved.EmailToken == nil || *retrieved.EmailToken != token {
		t.Errorf("Expected token %s, got %v", token, retrieved.EmailToken)
	}
	
	// Test with non-existent token
	_, err = repo.GetByEmailToken(ctx, "nonexistent_token")
	if err == nil {
		t.Error("Expected error for non-existent token")
	}
}

func TestUserRepository_SoftDelete(t *testing.T) {
	// Setup
	client, _, cleanup := helpers.CreateTestDatabase(t)
	defer cleanup()
	helpers.SetupTestTables(t, client)
	
	repo := adapters.NewPostgreSQLUserRepository(client)
	ctx := context.Background()
	
	// Create user
	user := domain.NewUser("Juan", "Pérez", "test@example.cl", "+56987654321", "hashedpassword")
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Soft delete
	err = repo.Delete(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to soft delete user: %v", err)
	}
	
	// Verify user is not found (soft deleted)
	_, err = repo.GetByID(ctx, user.ID)
	if err == nil {
		t.Error("Expected error when getting soft deleted user")
	}
	
	// Verify user doesn't exist in searches
	exists, err := repo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("Failed to check email existence: %v", err)
	}
	if exists {
		t.Error("Soft deleted user should not exist in email check")
	}
}

func TestUserRepository_TenantOperations(t *testing.T) {
	// Setup
	client, _, cleanup := helpers.CreateTestDatabase(t)
	defer cleanup()
	helpers.SetupTestTables(t, client)
	
	repo := adapters.NewPostgreSQLUserRepository(client)
	ctx := context.Background()
	
	// Create user
	user := domain.NewUser("Juan", "Pérez", "test@example.cl", "+56987654321", "hashedpassword")
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Create tenant association
	tenantID := uuid.New()
	tenantUser := domain.NewTenantUser(tenantID, user.ID)
	
	err = repo.AddUserToTenant(ctx, tenantUser)
	if err != nil {
		t.Fatalf("Failed to add user to tenant: %v", err)
	}
	
	// Get user tenants
	tenants, err := repo.GetUserTenants(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user tenants: %v", err)
	}
	
	if len(tenants) != 1 {
		t.Errorf("Expected 1 tenant, got %d", len(tenants))
	}
	
	if tenants[0].TenantID != tenantID {
		t.Errorf("Expected tenant ID %s, got %s", tenantID, tenants[0].TenantID)
	}
	
	// Remove user from tenant
	err = repo.RemoveUserFromTenant(ctx, user.ID, tenantID)
	if err != nil {
		t.Fatalf("Failed to remove user from tenant: %v", err)
	}
	
	// Verify removal
	tenants, err = repo.GetUserTenants(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user tenants after removal: %v", err)
	}
	
	if len(tenants) != 0 {
		t.Errorf("Expected 0 tenants after removal, got %d", len(tenants))
	}
}

// Benchmarks
func BenchmarkUserRepository_Create(b *testing.B) {
	if !helpers.CheckPostgreSQLAvailable() {
		b.Skip("PostgreSQL not available")
	}
	
	client, _, cleanup := helpers.CreateTestDatabase(&testing.T{})
	defer cleanup()
	helpers.SetupTestTables(&testing.T{}, client)
	
	repo := adapters.NewPostgreSQLUserRepository(client)
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := domain.NewUser(
			"Juan",
			"Pérez",
			fmt.Sprintf("test%d@example.cl", i),
			"+56987654321",
			"hashedpassword",
		)
		repo.Create(ctx, user)
	}
}

func BenchmarkUserRepository_GetByEmail(b *testing.B) {
	if !helpers.CheckPostgreSQLAvailable() {
		b.Skip("PostgreSQL not available")
	}
	
	client, _, cleanup := helpers.CreateTestDatabase(&testing.T{})
	defer cleanup()
	helpers.SetupTestTables(&testing.T{}, client)
	
	repo := adapters.NewPostgreSQLUserRepository(client)
	ctx := context.Background()
	
	// Create test user
	user := domain.NewUser("Juan", "Pérez", "test@example.cl", "+56987654321", "hashedpassword")
	repo.Create(ctx, user)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.GetByEmail(ctx, "test@example.cl")
	}
}
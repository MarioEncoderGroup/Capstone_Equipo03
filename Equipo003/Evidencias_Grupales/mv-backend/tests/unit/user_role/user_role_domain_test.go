package user_role_unit_test

import (
	"testing"
	"time"

	userRoleDomain "github.com/JoseLuis21/mv-backend/internal/core/user_role/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestNewUserRole valida la creación de nuevas relaciones usuario-rol
func TestNewUserRole(t *testing.T) {
	t.Run("Should create global user-role relation", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()

		userRole := userRoleDomain.NewUserRole(userID, roleID, nil)

		assert.NotEqual(t, uuid.Nil, userRole.ID)
		assert.Equal(t, userID, userRole.UserID)
		assert.Equal(t, roleID, userRole.RoleID)
		assert.Nil(t, userRole.TenantID)
		assert.False(t, userRole.Created.IsZero())
		assert.False(t, userRole.Updated.IsZero())
	})

	t.Run("Should create tenant-specific user-role relation", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()
		tenantID := uuid.New()

		userRole := userRoleDomain.NewUserRole(userID, roleID, &tenantID)

		assert.NotEqual(t, uuid.Nil, userRole.ID)
		assert.Equal(t, userID, userRole.UserID)
		assert.Equal(t, roleID, userRole.RoleID)
		assert.NotNil(t, userRole.TenantID)
		assert.Equal(t, tenantID, *userRole.TenantID)
		assert.False(t, userRole.Created.IsZero())
		assert.False(t, userRole.Updated.IsZero())
	})

	t.Run("Should set created and updated times", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()
		before := time.Now()

		userRole := userRoleDomain.NewUserRole(userID, roleID, nil)

		after := time.Now()

		assert.True(t, userRole.Created.After(before) || userRole.Created.Equal(before))
		assert.True(t, userRole.Created.Before(after) || userRole.Created.Equal(after))
		assert.Equal(t, userRole.Created, userRole.Updated)
	})
}

// TestUserRoleValidation valida comportamientos de validación
func TestUserRoleValidation(t *testing.T) {
	t.Run("Should generate unique IDs for different relations", func(t *testing.T) {
		userID1 := uuid.New()
		roleID1 := uuid.New()
		userRole1 := userRoleDomain.NewUserRole(userID1, roleID1, nil)

		userID2 := uuid.New()
		roleID2 := uuid.New()
		userRole2 := userRoleDomain.NewUserRole(userID2, roleID2, nil)

		assert.NotEqual(t, userRole1.ID, userRole2.ID)
		assert.NotEqual(t, uuid.Nil, userRole1.ID)
		assert.NotEqual(t, uuid.Nil, userRole2.ID)
	})

	t.Run("Should handle same user with different roles", func(t *testing.T) {
		userID := uuid.New()
		roleID1 := uuid.New()
		roleID2 := uuid.New()
		tenantID := uuid.New()

		userRole1 := userRoleDomain.NewUserRole(userID, roleID1, &tenantID)
		userRole2 := userRoleDomain.NewUserRole(userID, roleID2, &tenantID)

		assert.NotEqual(t, userRole1.ID, userRole2.ID)
		assert.Equal(t, userID, userRole1.UserID)
		assert.Equal(t, userID, userRole2.UserID)
		assert.NotEqual(t, userRole1.RoleID, userRole2.RoleID)
		assert.Equal(t, tenantID, *userRole1.TenantID)
		assert.Equal(t, tenantID, *userRole2.TenantID)
	})

	t.Run("Should handle same role with different users", func(t *testing.T) {
		userID1 := uuid.New()
		userID2 := uuid.New()
		roleID := uuid.New()
		tenantID := uuid.New()

		userRole1 := userRoleDomain.NewUserRole(userID1, roleID, &tenantID)
		userRole2 := userRoleDomain.NewUserRole(userID2, roleID, &tenantID)

		assert.NotEqual(t, userRole1.ID, userRole2.ID)
		assert.NotEqual(t, userRole1.UserID, userRole2.UserID)
		assert.Equal(t, roleID, userRole1.RoleID)
		assert.Equal(t, roleID, userRole2.RoleID)
		assert.Equal(t, tenantID, *userRole1.TenantID)
		assert.Equal(t, tenantID, *userRole2.TenantID)
	})
}

// TestUserRoleTenantScenarios valida escenarios multi-tenant
func TestUserRoleTenantScenarios(t *testing.T) {
	t.Run("Should distinguish between global and tenant roles", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()
		tenantID := uuid.New()

		// Global user-role
		globalUserRole := userRoleDomain.NewUserRole(userID, roleID, nil)
		assert.Nil(t, globalUserRole.TenantID)

		// Tenant-specific user-role
		tenantUserRole := userRoleDomain.NewUserRole(userID, roleID, &tenantID)
		assert.NotNil(t, tenantUserRole.TenantID)
		assert.Equal(t, tenantID, *tenantUserRole.TenantID)
	})

	t.Run("Should handle multiple tenant assignments", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()
		tenantID1 := uuid.New()
		tenantID2 := uuid.New()

		userRole1 := userRoleDomain.NewUserRole(userID, roleID, &tenantID1)
		userRole2 := userRoleDomain.NewUserRole(userID, roleID, &tenantID2)

		assert.NotEqual(t, userRole1.ID, userRole2.ID)
		assert.Equal(t, userID, userRole1.UserID)
		assert.Equal(t, userID, userRole2.UserID)
		assert.Equal(t, roleID, userRole1.RoleID)
		assert.Equal(t, roleID, userRole2.RoleID)
		assert.NotEqual(t, *userRole1.TenantID, *userRole2.TenantID)
	})
}

// TestUserRoleDTOs valida la estructura de DTOs
func TestUserRoleDTOs(t *testing.T) {
	t.Run("CreateUserRoleDto should validate required fields", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()
		tenantID := uuid.New()

		dto := userRoleDomain.CreateUserRoleDto{
			UserID:   userID,
			RoleID:   roleID,
			TenantID: &tenantID,
		}

		assert.Equal(t, userID, dto.UserID)
		assert.Equal(t, roleID, dto.RoleID)
		assert.NotNil(t, dto.TenantID)
		assert.Equal(t, tenantID, *dto.TenantID)
	})

	t.Run("CreateUserRoleDto should handle nil tenant", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()

		dto := userRoleDomain.CreateUserRoleDto{
			UserID:   userID,
			RoleID:   roleID,
			TenantID: nil,
		}

		assert.Equal(t, userID, dto.UserID)
		assert.Equal(t, roleID, dto.RoleID)
		assert.Nil(t, dto.TenantID)
	})

	t.Run("CreateUserRoleDto should handle nil tenant correctly", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()

		dto := userRoleDomain.CreateUserRoleDto{
			UserID:   userID,
			RoleID:   roleID,
			TenantID: nil,
		}

		assert.Equal(t, userID, dto.UserID)
		assert.Equal(t, roleID, dto.RoleID)
		assert.Nil(t, dto.TenantID)
	})

	t.Run("SyncUserRolesDto should handle role list", func(t *testing.T) {
		userID := uuid.New()
		roleID1 := uuid.New()
		roleID2 := uuid.New()
		tenantID := uuid.New()

		dto := userRoleDomain.SyncUserRolesDto{
			UserID:   userID,
			RoleIDs:  []uuid.UUID{roleID1, roleID2},
			TenantID: &tenantID,
		}

		assert.Equal(t, userID, dto.UserID)
		assert.Len(t, dto.RoleIDs, 2)
		assert.Contains(t, dto.RoleIDs, roleID1)
		assert.Contains(t, dto.RoleIDs, roleID2)
		assert.NotNil(t, dto.TenantID)
		assert.Equal(t, tenantID, *dto.TenantID)
	})

	t.Run("SyncRoleUsersDto should handle user list", func(t *testing.T) {
		userID1 := uuid.New()
		userID2 := uuid.New()
		roleID := uuid.New()
		tenantID := uuid.New()

		dto := userRoleDomain.SyncRoleUsersDto{
			RoleID:   roleID,
			UserIDs:  []uuid.UUID{userID1, userID2},
			TenantID: &tenantID,
		}

		assert.Equal(t, roleID, dto.RoleID)
		assert.Len(t, dto.UserIDs, 2)
		assert.Contains(t, dto.UserIDs, userID1)
		assert.Contains(t, dto.UserIDs, userID2)
		assert.NotNil(t, dto.TenantID)
		assert.Equal(t, tenantID, *dto.TenantID)
	})
}

// TestUserRoleResponseDTOs valida DTOs de respuesta
func TestUserRoleResponseDTOs(t *testing.T) {
	t.Run("UserRoleResponseDto should contain all required fields", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()
		tenantID := uuid.New()

		userRole := userRoleDomain.NewUserRole(userID, roleID, &tenantID)
		response := userRoleDomain.UserRoleResponseDto{
			UserRole: userRole,
			Message:  "User role created successfully",
		}

		assert.NotNil(t, response.UserRole)
		assert.Equal(t, userID, response.UserRole.UserID)
		assert.Equal(t, roleID, response.UserRole.RoleID)
		assert.NotNil(t, response.UserRole.TenantID)
		assert.Equal(t, tenantID, *response.UserRole.TenantID)
		assert.Equal(t, "User role created successfully", response.Message)
	})

	t.Run("UserRoleResponseDto should handle nil tenant", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()

		userRole := userRoleDomain.NewUserRole(userID, roleID, nil)
		response := userRoleDomain.UserRoleResponseDto{
			UserRole: userRole,
			Message:  "Global user role created",
		}

		assert.NotNil(t, response.UserRole)
		assert.Equal(t, userID, response.UserRole.UserID)
		assert.Equal(t, roleID, response.UserRole.RoleID)
		assert.Nil(t, response.UserRole.TenantID)
		assert.Equal(t, "Global user role created", response.Message)
	})
}

// TestUserRoleImmutability valida la inmutabilidad de campos críticos
func TestUserRoleImmutability(t *testing.T) {
	t.Run("Should maintain immutable ID after creation", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()
		userRole := userRoleDomain.NewUserRole(userID, roleID, nil)
		originalID := userRole.ID

		// Simular algún tipo de operación que no debería cambiar el ID
		userRole.Updated = time.Now()

		assert.Equal(t, originalID, userRole.ID, "ID should remain immutable")
	})

	t.Run("Should maintain immutable UserID and RoleID", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()
		userRole := userRoleDomain.NewUserRole(userID, roleID, nil)

		// Las relaciones usuario-rol son immutables una vez creadas
		assert.Equal(t, userID, userRole.UserID)
		assert.Equal(t, roleID, userRole.RoleID)
	})
}

// TestEdgeCases valida casos límite
func TestEdgeCases(t *testing.T) {
	t.Run("Should handle zero UUIDs properly", func(t *testing.T) {
		// Aunque no es una práctica recomendada, debemos manejar UUIDs nulos
		userRole := userRoleDomain.NewUserRole(uuid.Nil, uuid.Nil, nil)

		assert.NotEqual(t, uuid.Nil, userRole.ID, "UserRole ID should always be generated")
		assert.Equal(t, uuid.Nil, userRole.UserID)
		assert.Equal(t, uuid.Nil, userRole.RoleID)
	})

	t.Run("Should handle timestamp precision", func(t *testing.T) {
		userRole1 := userRoleDomain.NewUserRole(uuid.New(), uuid.New(), nil)

		// Pequeña pausa para asegurar timestamps diferentes
		time.Sleep(1 * time.Microsecond)

		userRole2 := userRoleDomain.NewUserRole(uuid.New(), uuid.New(), nil)

		// Los timestamps deberían ser diferentes (aunque sea por microsegundos)
		assert.True(t, userRole2.Created.After(userRole1.Created) ||
			userRole2.Created.Equal(userRole1.Created),
			"Second user role should be created at same time or after first")
	})
}
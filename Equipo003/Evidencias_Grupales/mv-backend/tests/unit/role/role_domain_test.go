package role_unit_test

import (
	"testing"
	"time"

	roleDomain "github.com/JoseLuis21/mv-backend/internal/core/role/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewRole valida la creación de nuevos roles
func TestNewRole(t *testing.T) {
	t.Run("Should create global role", func(t *testing.T) {
		name := "administrator"
		description := "System administrator with full access"

		role := roleDomain.NewRole(name, description, nil)

		assert.NotEqual(t, uuid.Nil, role.ID)
		assert.Nil(t, role.TenantID)
		assert.Equal(t, name, role.Name)
		assert.NotNil(t, role.Description)
		assert.Equal(t, description, *role.Description)
		assert.False(t, role.Created.IsZero())
		assert.False(t, role.Updated.IsZero())
		assert.Nil(t, role.DeletedAt)
	})

	t.Run("Should create tenant role", func(t *testing.T) {
		tenantID := uuid.New()
		name := "manager"
		description := "Department manager role"

		role := roleDomain.NewRole(name, description, &tenantID)

		assert.NotEqual(t, uuid.Nil, role.ID)
		assert.Equal(t, &tenantID, role.TenantID)
		assert.Equal(t, name, role.Name)
		assert.NotNil(t, role.Description)
		assert.Equal(t, description, *role.Description)
		assert.False(t, role.Created.IsZero())
		assert.False(t, role.Updated.IsZero())
		assert.Nil(t, role.DeletedAt)
	})

	t.Run("Should create role with empty description", func(t *testing.T) {
		name := "basic-user"
		description := ""

		role := roleDomain.NewRole(name, description, nil)

		assert.NotEqual(t, uuid.Nil, role.ID)
		assert.Equal(t, name, role.Name)
		assert.Nil(t, role.Description) // Empty description should be nil
	})

	t.Run("Should set created and updated times", func(t *testing.T) {
		before := time.Now()
		role := roleDomain.NewRole("test-role", "Test role", nil)
		after := time.Now()

		assert.True(t, role.Created.After(before) || role.Created.Equal(before))
		assert.True(t, role.Created.Before(after) || role.Created.Equal(after))
		assert.Equal(t, role.Created, role.Updated)
	})
}

// TestRoleSystemMethods valida los métodos relacionados con roles del sistema
func TestRoleSystemMethods(t *testing.T) {
	t.Run("IsGlobalRole should work correctly", func(t *testing.T) {
		// Global role
		globalRole := roleDomain.NewRole("administrator", "Admin role", nil)
		assert.True(t, globalRole.IsGlobalRole())
		assert.False(t, globalRole.IsTenantRole())

		// Tenant role
		tenantID := uuid.New()
		tenantRole := roleDomain.NewRole("manager", "Manager role", &tenantID)
		assert.False(t, tenantRole.IsGlobalRole())
		assert.True(t, tenantRole.IsTenantRole())
	})

	t.Run("IsSystemRole should identify system roles", func(t *testing.T) {
		// System roles
		adminRole := roleDomain.NewRole(roleDomain.RoleNameAdministrator, "Admin", nil)
		approverRole := roleDomain.NewRole(roleDomain.RoleNameApprover, "Approver", nil)
		submitterRole := roleDomain.NewRole(roleDomain.RoleNameExpenseSubmitter, "Submitter", nil)

		assert.True(t, adminRole.IsSystemRole())
		assert.True(t, approverRole.IsSystemRole())
		assert.True(t, submitterRole.IsSystemRole())

		// Custom role
		customRole := roleDomain.NewRole("custom-role", "Custom role", nil)
		assert.False(t, customRole.IsSystemRole())
	})

	t.Run("Should verify system role constants", func(t *testing.T) {
		assert.Equal(t, "administrator", roleDomain.RoleNameAdministrator)
		assert.Equal(t, "approver", roleDomain.RoleNameApprover)
		assert.Equal(t, "expense-submitter", roleDomain.RoleNameExpenseSubmitter)
	})
}

// TestRoleUpdate valida la actualización de roles
func TestRoleUpdate(t *testing.T) {
	t.Run("Should update role with description", func(t *testing.T) {
		role := roleDomain.NewRole("original-name", "Original description", nil)
		originalUpdated := role.Updated

		// Wait a moment to ensure different timestamps
		time.Sleep(1 * time.Millisecond)

		newName := "updated-name"
		newDescription := "Updated description"
		role.Update(newName, newDescription)

		assert.Equal(t, newName, role.Name)
		assert.NotNil(t, role.Description)
		assert.Equal(t, newDescription, *role.Description)
		assert.True(t, role.Updated.After(originalUpdated))
	})

	t.Run("Should update role with empty description", func(t *testing.T) {
		role := roleDomain.NewRole("role-name", "Original description", nil)

		role.Update("updated-name", "")

		assert.Equal(t, "updated-name", role.Name)
		assert.Nil(t, role.Description) // Empty description should set to nil
	})

	t.Run("Should update timestamp on each update", func(t *testing.T) {
		role := roleDomain.NewRole("test-role", "Test", nil)
		firstUpdate := role.Updated

		time.Sleep(1 * time.Millisecond)
		role.Update("updated-role", "Updated")
		secondUpdate := role.Updated

		time.Sleep(1 * time.Millisecond)
		role.Update("final-role", "Final")
		thirdUpdate := role.Updated

		assert.True(t, secondUpdate.After(firstUpdate))
		assert.True(t, thirdUpdate.After(secondUpdate))
	})
}

// TestRoleSoftDelete valida la eliminación lógica de roles
func TestRoleSoftDelete(t *testing.T) {
	t.Run("Should perform soft delete", func(t *testing.T) {
		role := roleDomain.NewRole("test-role", "Test role", nil)
		originalUpdated := role.Updated

		assert.Nil(t, role.DeletedAt)

		// Wait a moment to ensure different timestamps
		time.Sleep(1 * time.Millisecond)

		role.SoftDelete()

		assert.NotNil(t, role.DeletedAt)
		assert.False(t, role.DeletedAt.IsZero())
		assert.True(t, role.Updated.After(originalUpdated))
		assert.Equal(t, *role.DeletedAt, role.Updated)
	})

	t.Run("Should handle multiple soft deletes", func(t *testing.T) {
		role := roleDomain.NewRole("test-role", "Test role", nil)

		role.SoftDelete()
		firstDeleteTime := *role.DeletedAt

		time.Sleep(1 * time.Millisecond)
		role.SoftDelete()
		secondDeleteTime := *role.DeletedAt

		assert.True(t, secondDeleteTime.After(firstDeleteTime))
	})
}

// TestRoleToResponse valida la conversión a DTO de respuesta
func TestRoleToResponse(t *testing.T) {
	t.Run("Should convert global role to response", func(t *testing.T) {
		role := roleDomain.NewRole(roleDomain.RoleNameAdministrator, "System admin", nil)

		response := role.ToResponse()

		assert.Equal(t, role.ID, response.ID)
		assert.Equal(t, role.TenantID, response.TenantID)
		assert.Equal(t, role.Name, response.Name)
		assert.Equal(t, role.Description, response.Description)
		assert.True(t, response.IsGlobal)
		assert.True(t, response.IsSystem)
		assert.NotEmpty(t, response.Created)
		assert.NotEmpty(t, response.Updated)
	})

	t.Run("Should convert tenant role to response", func(t *testing.T) {
		tenantID := uuid.New()
		description := "Custom tenant role"
		role := roleDomain.NewRole("custom-role", description, &tenantID)

		response := role.ToResponse()

		assert.Equal(t, role.ID, response.ID)
		assert.Equal(t, &tenantID, response.TenantID)
		assert.Equal(t, "custom-role", response.Name)
		assert.Equal(t, &description, response.Description)
		assert.False(t, response.IsGlobal)
		assert.False(t, response.IsSystem)
	})

	t.Run("Should format timestamps in ISO 8601", func(t *testing.T) {
		role := roleDomain.NewRole("test-role", "Test", nil)

		response := role.ToResponse()

		// Verify timestamp format (ISO 8601)
		_, err := time.Parse("2006-01-02T15:04:05Z07:00", response.Created)
		require.NoError(t, err, "Created timestamp should be in ISO 8601 format")

		_, err = time.Parse("2006-01-02T15:04:05Z07:00", response.Updated)
		require.NoError(t, err, "Updated timestamp should be in ISO 8601 format")
	})

	t.Run("Should handle nil description", func(t *testing.T) {
		role := roleDomain.NewRole("test-role", "", nil) // Empty description becomes nil

		response := role.ToResponse()

		assert.Nil(t, response.Description)
	})
}

// TestRoleTypes valida las constantes de tipos de rol
func TestRoleTypes(t *testing.T) {
	t.Run("Should have correct role type constants", func(t *testing.T) {
		assert.Equal(t, roleDomain.RoleType("global"), roleDomain.RoleTypeGlobal)
		assert.Equal(t, roleDomain.RoleType("tenant"), roleDomain.RoleTypeTenant)
	})
}

// TestRoleValidation valida comportamientos de validación implícitos
func TestRoleValidation(t *testing.T) {
	t.Run("Should generate unique IDs", func(t *testing.T) {
		role1 := roleDomain.NewRole("role1", "Description 1", nil)
		role2 := roleDomain.NewRole("role2", "Description 2", nil)

		assert.NotEqual(t, role1.ID, role2.ID)
		assert.NotEqual(t, uuid.Nil, role1.ID)
		assert.NotEqual(t, uuid.Nil, role2.ID)
	})

	t.Run("Should handle various tenant ID scenarios", func(t *testing.T) {
		// Nil tenant (global)
		globalRole := roleDomain.NewRole("global", "Global role", nil)
		assert.True(t, globalRole.IsGlobalRole())

		// Valid tenant ID
		tenantID := uuid.New()
		tenantRole := roleDomain.NewRole("tenant", "Tenant role", &tenantID)
		assert.True(t, tenantRole.IsTenantRole())
		assert.Equal(t, tenantID, *tenantRole.TenantID)
	})
}
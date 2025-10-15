package permission_unit_test

import (
	"testing"
	"time"

	permissionDomain "github.com/JoseLuis21/mv-backend/internal/core/permission/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewPermission valida la creación de nuevos permisos
func TestNewPermission(t *testing.T) {
	t.Run("Should create permission with description", func(t *testing.T) {
		name := "create-users"
		description := "Permission to create new users"
		section := "user"

		permission := permissionDomain.NewPermission(name, description, section)

		assert.NotEqual(t, uuid.Nil, permission.ID)
		assert.Equal(t, name, permission.Name)
		assert.NotNil(t, permission.Description)
		assert.Equal(t, description, *permission.Description)
		assert.Equal(t, section, permission.Section)
		assert.False(t, permission.Created.IsZero())
		assert.False(t, permission.Updated.IsZero())
		assert.Nil(t, permission.DeletedAt)
	})

	t.Run("Should create permission with empty description", func(t *testing.T) {
		name := "view-dashboard"
		description := ""
		section := "dashboard"

		permission := permissionDomain.NewPermission(name, description, section)

		assert.NotEqual(t, uuid.Nil, permission.ID)
		assert.Equal(t, name, permission.Name)
		assert.Nil(t, permission.Description) // Empty description should be nil
		assert.Equal(t, section, permission.Section)
	})

	t.Run("Should set created and updated times", func(t *testing.T) {
		before := time.Now()
		permission := permissionDomain.NewPermission("test-permission", "Test", "test")
		after := time.Now()

		assert.True(t, permission.Created.After(before) || permission.Created.Equal(before))
		assert.True(t, permission.Created.Before(after) || permission.Created.Equal(after))
		assert.Equal(t, permission.Created, permission.Updated)
	})
}

// TestPermissionSystemMethods valida los métodos relacionados con permisos del sistema
func TestPermissionSystemMethods(t *testing.T) {
	t.Run("IsSystemPermission should identify system permissions", func(t *testing.T) {
		// System permissions
		listRolePermission := permissionDomain.NewPermission(permissionDomain.PermissionListRole, "List roles", string(permissionDomain.SectionRole))
		createUserPermission := permissionDomain.NewPermission(permissionDomain.PermissionCreateUser, "Create users", string(permissionDomain.SectionUser))
		updatePermissionPermission := permissionDomain.NewPermission(permissionDomain.PermissionUpdatePermission, "Update permissions", string(permissionDomain.SectionPermission))

		assert.True(t, listRolePermission.IsSystemPermission())
		assert.True(t, createUserPermission.IsSystemPermission())
		assert.True(t, updatePermissionPermission.IsSystemPermission())

		// Custom permission
		customPermission := permissionDomain.NewPermission("custom-action", "Custom permission", "custom")
		assert.False(t, customPermission.IsSystemPermission())
	})

	t.Run("Should verify system permission constants", func(t *testing.T) {
		assert.Equal(t, "list-role", permissionDomain.PermissionListRole)
		assert.Equal(t, "create-role", permissionDomain.PermissionCreateRole)
		assert.Equal(t, "update-role", permissionDomain.PermissionUpdateRole)
		assert.Equal(t, "delete-role", permissionDomain.PermissionDeleteRole)

		assert.Equal(t, "list-user", permissionDomain.PermissionListUser)
		assert.Equal(t, "create-user", permissionDomain.PermissionCreateUser)
		assert.Equal(t, "update-user", permissionDomain.PermissionUpdateUser)
		assert.Equal(t, "delete-user", permissionDomain.PermissionDeleteUser)
	})
}

// TestPermissionSectionValidation valida las secciones de permisos
func TestPermissionSectionValidation(t *testing.T) {
	t.Run("IsValidSection should validate permission sections", func(t *testing.T) {
		// Valid sections
		rolePermission := permissionDomain.NewPermission("test", "test", string(permissionDomain.SectionRole))
		userPermission := permissionDomain.NewPermission("test", "test", string(permissionDomain.SectionUser))
		tenantPermission := permissionDomain.NewPermission("test", "test", string(permissionDomain.SectionTenant))

		assert.True(t, rolePermission.IsValidSection())
		assert.True(t, userPermission.IsValidSection())
		assert.True(t, tenantPermission.IsValidSection())

		// Invalid section
		invalidPermission := permissionDomain.NewPermission("test", "test", "invalid-section")
		assert.False(t, invalidPermission.IsValidSection())
	})

	t.Run("Should have correct section constants", func(t *testing.T) {
		assert.Equal(t, "role", string(permissionDomain.SectionRole))
		assert.Equal(t, "permission", string(permissionDomain.SectionPermission))
		assert.Equal(t, "user", string(permissionDomain.SectionUser))
		assert.Equal(t, "tenant", string(permissionDomain.SectionTenant))
		assert.Equal(t, "category", string(permissionDomain.SectionCategory))
	})

	t.Run("GetPermissionsBySection should return correct permissions", func(t *testing.T) {
		rolePermissions := permissionDomain.GetPermissionsBySection(permissionDomain.SectionRole)
		expectedRolePerms := []string{
			permissionDomain.PermissionListRole,
			permissionDomain.PermissionCreateRole,
			permissionDomain.PermissionUpdateRole,
			permissionDomain.PermissionDeleteRole,
		}
		assert.ElementsMatch(t, expectedRolePerms, rolePermissions)

		userPermissions := permissionDomain.GetPermissionsBySection(permissionDomain.SectionUser)
		expectedUserPerms := []string{
			permissionDomain.PermissionListUser,
			permissionDomain.PermissionCreateUser,
			permissionDomain.PermissionUpdateUser,
			permissionDomain.PermissionDeleteUser,
		}
		assert.ElementsMatch(t, expectedUserPerms, userPermissions)

		// Non-existent section
		nonExistentPerms := permissionDomain.GetPermissionsBySection(permissionDomain.PermissionSection("non-existent"))
		assert.Empty(t, nonExistentPerms)
	})
}

// TestPermissionUpdate valida la actualización de permisos
func TestPermissionUpdate(t *testing.T) {
	t.Run("Should update permission with description", func(t *testing.T) {
		permission := permissionDomain.NewPermission("original-permission", "Original description", "original")
		originalUpdated := permission.Updated

		// Wait a moment to ensure different timestamps
		time.Sleep(1 * time.Millisecond)

		newName := "updated-permission"
		newDescription := "Updated description"
		newSection := "updated"
		permission.Update(newName, newDescription, newSection)

		assert.Equal(t, newName, permission.Name)
		assert.NotNil(t, permission.Description)
		assert.Equal(t, newDescription, *permission.Description)
		assert.Equal(t, newSection, permission.Section)
		assert.True(t, permission.Updated.After(originalUpdated))
	})

	t.Run("Should update permission with empty description", func(t *testing.T) {
		permission := permissionDomain.NewPermission("permission-name", "Original description", "original")

		permission.Update("updated-permission", "", "updated")

		assert.Equal(t, "updated-permission", permission.Name)
		assert.Nil(t, permission.Description) // Empty description should set to nil
		assert.Equal(t, "updated", permission.Section)
	})

	t.Run("Should update timestamp on each update", func(t *testing.T) {
		permission := permissionDomain.NewPermission("test-permission", "Test", "test")
		firstUpdate := permission.Updated

		time.Sleep(1 * time.Millisecond)
		permission.Update("updated-permission", "Updated", "test")
		secondUpdate := permission.Updated

		time.Sleep(1 * time.Millisecond)
		permission.Update("final-permission", "Final", "test")
		thirdUpdate := permission.Updated

		assert.True(t, secondUpdate.After(firstUpdate))
		assert.True(t, thirdUpdate.After(secondUpdate))
	})
}

// TestPermissionSoftDelete valida la eliminación lógica de permisos
func TestPermissionSoftDelete(t *testing.T) {
	t.Run("Should perform soft delete", func(t *testing.T) {
		permission := permissionDomain.NewPermission("test-permission", "Test permission", "test")
		originalUpdated := permission.Updated

		assert.Nil(t, permission.DeletedAt)
		assert.True(t, permission.IsActive())

		// Wait a moment to ensure different timestamps
		time.Sleep(1 * time.Millisecond)

		permission.SoftDelete()

		assert.NotNil(t, permission.DeletedAt)
		assert.False(t, permission.DeletedAt.IsZero())
		assert.False(t, permission.IsActive())
		assert.True(t, permission.Updated.After(originalUpdated))
		assert.Equal(t, *permission.DeletedAt, permission.Updated)
	})

	t.Run("Should handle multiple soft deletes", func(t *testing.T) {
		permission := permissionDomain.NewPermission("test-permission", "Test permission", "test")

		permission.SoftDelete()
		firstDeleteTime := *permission.DeletedAt

		time.Sleep(1 * time.Millisecond)
		permission.SoftDelete()
		secondDeleteTime := *permission.DeletedAt

		assert.True(t, secondDeleteTime.After(firstDeleteTime))
	})
}

// TestPermissionToResponse valida la conversión a DTO de respuesta
func TestPermissionToResponse(t *testing.T) {
	t.Run("Should convert system permission to response", func(t *testing.T) {
		permission := permissionDomain.NewPermission(permissionDomain.PermissionCreateUser, "Create users", string(permissionDomain.SectionUser))

		response := permission.ToResponse()

		assert.Equal(t, permission.ID, response.ID)
		assert.Equal(t, permission.Name, response.Name)
		assert.Equal(t, permission.Description, response.Description)
		assert.Equal(t, permission.Section, response.Section)
		assert.True(t, response.IsSystem)
		assert.True(t, response.IsActive)
		assert.NotEmpty(t, response.Created)
		assert.NotEmpty(t, response.Updated)
	})

	t.Run("Should convert custom permission to response", func(t *testing.T) {
		description := "Custom permission"
		permission := permissionDomain.NewPermission("custom-permission", description, "custom")

		response := permission.ToResponse()

		assert.Equal(t, permission.ID, response.ID)
		assert.Equal(t, "custom-permission", response.Name)
		assert.Equal(t, &description, response.Description)
		assert.Equal(t, "custom", response.Section)
		assert.False(t, response.IsSystem)
		assert.True(t, response.IsActive)
	})

	t.Run("Should format timestamps in ISO 8601", func(t *testing.T) {
		permission := permissionDomain.NewPermission("test-permission", "Test", "test")

		response := permission.ToResponse()

		// Verify timestamp format (ISO 8601)
		_, err := time.Parse("2006-01-02T15:04:05Z07:00", response.Created)
		require.NoError(t, err, "Created timestamp should be in ISO 8601 format")

		_, err = time.Parse("2006-01-02T15:04:05Z07:00", response.Updated)
		require.NoError(t, err, "Updated timestamp should be in ISO 8601 format")
	})

	t.Run("Should handle nil description", func(t *testing.T) {
		permission := permissionDomain.NewPermission("test-permission", "", "test") // Empty description becomes nil

		response := permission.ToResponse()

		assert.Nil(t, response.Description)
	})

	t.Run("Should handle soft deleted permission", func(t *testing.T) {
		permission := permissionDomain.NewPermission("test-permission", "Test", "test")
		permission.SoftDelete()

		response := permission.ToResponse()

		assert.False(t, response.IsActive)
	})
}

// TestPermissionValidation valida comportamientos de validación implícitos
func TestPermissionValidation(t *testing.T) {
	t.Run("Should generate unique IDs", func(t *testing.T) {
		permission1 := permissionDomain.NewPermission("permission1", "Description 1", "section1")
		permission2 := permissionDomain.NewPermission("permission2", "Description 2", "section2")

		assert.NotEqual(t, permission1.ID, permission2.ID)
		assert.NotEqual(t, uuid.Nil, permission1.ID)
		assert.NotEqual(t, uuid.Nil, permission2.ID)
	})

	t.Run("Should handle various section scenarios", func(t *testing.T) {
		// Valid predefined sections
		rolePermission := permissionDomain.NewPermission("test", "test", string(permissionDomain.SectionRole))
		assert.True(t, rolePermission.IsValidSection())

		userPermission := permissionDomain.NewPermission("test", "test", string(permissionDomain.SectionUser))
		assert.True(t, userPermission.IsValidSection())

		// Custom section (should be invalid according to validation)
		customPermission := permissionDomain.NewPermission("test", "test", "custom-section")
		assert.False(t, customPermission.IsValidSection())
	})

	t.Run("Should maintain immutable ID after creation", func(t *testing.T) {
		permission := permissionDomain.NewPermission("test-permission", "Test", "test")
		originalID := permission.ID

		permission.Update("updated-permission", "Updated description", "updated")

		assert.Equal(t, originalID, permission.ID, "ID should not change after update")

		permission.SoftDelete()

		assert.Equal(t, originalID, permission.ID, "ID should not change after soft delete")
	})
}
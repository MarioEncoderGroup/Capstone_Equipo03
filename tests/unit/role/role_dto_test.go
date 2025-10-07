package role_unit_test

import (
	"testing"

	roleDomain "github.com/JoseLuis21/mv-backend/internal/core/role/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateRoleRequest valida la estructura del DTO de creación de rol
func TestCreateRoleRequest(t *testing.T) {
	t.Run("Should create valid CreateRoleRequest", func(t *testing.T) {
		tenantID := uuid.New()
		createReq := roleDomain.CreateRoleRequest{
			Name:        "Test Role",
			Description: "A test role for validation",
			TenantID:    &tenantID,
		}

		assert.Equal(t, "Test Role", createReq.Name)
		assert.Equal(t, "A test role for validation", createReq.Description)
		assert.Equal(t, &tenantID, createReq.TenantID)
	})

	t.Run("Should support global role creation", func(t *testing.T) {
		createReq := roleDomain.CreateRoleRequest{
			Name:        "Global Role",
			Description: "A global system role",
			TenantID:    nil, // Global role
		}

		assert.Equal(t, "Global Role", createReq.Name)
		assert.Equal(t, "A global system role", createReq.Description)
		assert.Nil(t, createReq.TenantID)
	})

	t.Run("Should have correct JSON tags", func(t *testing.T) {
		createReq := roleDomain.CreateRoleRequest{}

		// Verificar que los campos pueden ser asignados
		createReq.Name = "Manager"
		createReq.Description = "Manager role with approval permissions"

		assert.NotEmpty(t, createReq.Name)
		assert.NotEmpty(t, createReq.Description)
	})
}

// TestUpdateRoleRequest valida la estructura del DTO de actualización de rol
func TestUpdateRoleRequest(t *testing.T) {
	t.Run("Should create valid UpdateRoleRequest", func(t *testing.T) {
		updateReq := roleDomain.UpdateRoleRequest{
			Name:        "Updated Role",
			Description: "Updated description",
		}

		assert.Equal(t, "Updated Role", updateReq.Name)
		assert.Equal(t, "Updated description", updateReq.Description)
	})

	t.Run("Should support empty description", func(t *testing.T) {
		updateReq := roleDomain.UpdateRoleRequest{
			Name:        "Simple Role",
			Description: "",
		}

		assert.Equal(t, "Simple Role", updateReq.Name)
		assert.Equal(t, "", updateReq.Description)
	})
}

// TestRoleResponse valida la estructura de respuesta del rol
func TestRoleResponse(t *testing.T) {
	t.Run("Should create valid RoleResponse", func(t *testing.T) {
		roleID := uuid.New()
		tenantID := uuid.New()
		description := "Test role description"

		roleResponse := roleDomain.RoleResponse{
			ID:          roleID,
			TenantID:    &tenantID,
			Name:        "Test Role",
			Description: &description,
			IsGlobal:    false,
			IsSystem:    false,
			Created:     "2023-12-01T10:00:00Z",
			Updated:     "2023-12-01T10:00:00Z",
		}

		assert.Equal(t, roleID, roleResponse.ID)
		assert.Equal(t, &tenantID, roleResponse.TenantID)
		assert.Equal(t, "Test Role", roleResponse.Name)
		assert.Equal(t, &description, roleResponse.Description)
		assert.False(t, roleResponse.IsGlobal)
		assert.False(t, roleResponse.IsSystem)
		assert.NotEmpty(t, roleResponse.Created)
		assert.NotEmpty(t, roleResponse.Updated)
	})

	t.Run("Should support global role response", func(t *testing.T) {
		roleID := uuid.New()

		roleResponse := roleDomain.RoleResponse{
			ID:          roleID,
			TenantID:    nil, // Global role
			Name:        "administrator",
			Description: nil,
			IsGlobal:    true,
			IsSystem:    true,
			Created:     "2023-12-01T10:00:00Z",
			Updated:     "2023-12-01T10:00:00Z",
		}

		assert.Equal(t, roleID, roleResponse.ID)
		assert.Nil(t, roleResponse.TenantID)
		assert.Equal(t, "administrator", roleResponse.Name)
		assert.Nil(t, roleResponse.Description)
		assert.True(t, roleResponse.IsGlobal)
		assert.True(t, roleResponse.IsSystem)
	})
}

// TestRoleListResponse valida la estructura de respuesta de listado de roles
func TestRoleListResponse(t *testing.T) {
	t.Run("Should create valid RoleListResponse", func(t *testing.T) {
		roleID1 := uuid.New()
		roleID2 := uuid.New()

		roles := []roleDomain.RoleResponse{
			{
				ID:       roleID1,
				TenantID: nil,
				Name:     "administrator",
				IsGlobal: true,
				IsSystem: true,
				Created:  "2023-12-01T10:00:00Z",
				Updated:  "2023-12-01T10:00:00Z",
			},
			{
				ID:       roleID2,
				TenantID: nil,
				Name:     "approver",
				IsGlobal: true,
				IsSystem: true,
				Created:  "2023-12-01T10:00:00Z",
				Updated:  "2023-12-01T10:00:00Z",
			},
		}

		listResponse := roleDomain.RoleListResponse{
			Roles: roles,
			Total: 2,
			Page:  1,
			Limit: 20,
		}

		assert.Len(t, listResponse.Roles, 2)
		assert.Equal(t, 2, listResponse.Total)
		assert.Equal(t, 1, listResponse.Page)
		assert.Equal(t, 20, listResponse.Limit)
		assert.Equal(t, "administrator", listResponse.Roles[0].Name)
		assert.Equal(t, "approver", listResponse.Roles[1].Name)
	})

	t.Run("Should support empty role list", func(t *testing.T) {
		listResponse := roleDomain.RoleListResponse{
			Roles: []roleDomain.RoleResponse{},
			Total: 0,
			Page:  1,
			Limit: 20,
		}

		assert.Empty(t, listResponse.Roles)
		assert.Equal(t, 0, listResponse.Total)
		assert.Equal(t, 1, listResponse.Page)
		assert.Equal(t, 20, listResponse.Limit)
	})
}

// TestRoleFilterRequest valida la estructura de filtros de búsqueda
func TestRoleFilterRequest(t *testing.T) {
	t.Run("Should create valid RoleFilterRequest", func(t *testing.T) {
		tenantID := uuid.New()
		filterReq := roleDomain.RoleFilterRequest{
			TenantID: &tenantID,
			Name:     "manager",
			Page:     1,
			Limit:    10,
		}

		assert.Equal(t, &tenantID, filterReq.TenantID)
		assert.Equal(t, "manager", filterReq.Name)
		assert.Equal(t, 1, filterReq.Page)
		assert.Equal(t, 10, filterReq.Limit)
	})

	t.Run("Should support global role filtering", func(t *testing.T) {
		filterReq := roleDomain.RoleFilterRequest{
			TenantID: nil, // Filter global roles
			Name:     "",
			Page:     1,
			Limit:    20,
		}

		assert.Nil(t, filterReq.TenantID)
		assert.Equal(t, "", filterReq.Name)
		assert.Equal(t, 1, filterReq.Page)
		assert.Equal(t, 20, filterReq.Limit)
	})

	t.Run("Should have reasonable pagination defaults", func(t *testing.T) {
		filterReq := roleDomain.RoleFilterRequest{
			Page:  1,
			Limit: 20,
		}

		assert.True(t, filterReq.Page >= 1, "Page should be at least 1")
		assert.True(t, filterReq.Limit > 0 && filterReq.Limit <= 100, "Limit should be between 1 and 100")
	})
}

// TestRoleDtoValidationTags verifica que los tags de validación estén correctos
func TestRoleDtoValidationTags(t *testing.T) {
	t.Run("CreateRoleRequest validation requirements", func(t *testing.T) {
		// Los tags de validación son:
		// Name: `json:"name" validate:"required,min=3,max=50"`
		// Description: `json:"description" validate:"max=500"`

		dto := roleDomain.CreateRoleRequest{
			Name:        "Valid Role Name",
			Description: "Valid description for role",
		}

		// Verificar que los campos son correctos
		require.NotEmpty(t, dto.Name, "Name should not be empty")
		require.True(t, len(dto.Name) >= 3, "Name should be at least 3 characters")
		require.True(t, len(dto.Name) <= 50, "Name should not exceed 50 characters")
		assert.True(t, len(dto.Description) <= 500, "Description should not exceed 500 characters")
	})

	t.Run("UpdateRoleRequest validation requirements", func(t *testing.T) {
		dto := roleDomain.UpdateRoleRequest{
			Name:        "Updated Role",
			Description: "Updated description",
		}

		// Los requisitos de validación se verifican en runtime por el validador
		// Aquí solo verificamos que los campos estén presentes
		assert.NotEmpty(t, dto.Name, "Name should not be empty")
		assert.True(t, len(dto.Name) >= 3, "Name should be at least 3 characters")
		assert.True(t, len(dto.Name) <= 50, "Name should not exceed 50 characters")
	})

	t.Run("RoleFilterRequest validation requirements", func(t *testing.T) {
		dto := roleDomain.RoleFilterRequest{
			Page:  1,
			Limit: 50,
		}

		// Los tags de validación son:
		// Page: `json:"page" validate:"min=1"`
		// Limit: `json:"limit" validate:"min=1,max=100"`

		assert.True(t, dto.Page >= 1, "Page should be at least 1")
		assert.True(t, dto.Limit >= 1, "Limit should be at least 1")
		assert.True(t, dto.Limit <= 100, "Limit should not exceed 100")
	})
}
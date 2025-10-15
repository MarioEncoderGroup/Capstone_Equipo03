package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	permissionDomain "github.com/JoseLuis21/mv-backend/internal/core/permission/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/permission/ports"
	permissionServices "github.com/JoseLuis21/mv-backend/internal/core/permission/services"
	"github.com/JoseLuis21/mv-backend/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// PermissionServiceTestSuite agrupa todos los tests del servicio de permisos
type PermissionServiceTestSuite struct {
	suite.Suite
	mockRepo           *mocks.MockPermissionRepository
	permissionService  ports.PermissionService
	ctx                context.Context
}

// SetupTest configura cada test
func (suite *PermissionServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.MockPermissionRepository)
	suite.permissionService = permissionServices.NewPermissionService(suite.mockRepo)
	suite.ctx = context.Background()
}

// TearDownTest limpia después de cada test
func (suite *PermissionServiceTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestCreatePermission valida la creación de permisos
func (suite *PermissionServiceTestSuite) TestCreatePermission() {
	suite.Run("Should create permission successfully", func() {
		req := &permissionDomain.CreatePermissionRequest{
			Name:        "test-permission",
			Description: "Test permission description",
			Section:     "test",
		}

		// Mock: verificar que no existe el permiso
		suite.mockRepo.On("ExistsByName", suite.ctx, req.Name).Return(false, nil)

		// Mock: crear el permiso
		suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*domain.Permission")).
			Return(nil)

		response, err := suite.permissionService.CreatePermission(suite.ctx, req)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), response)
		assert.Equal(suite.T(), req.Name, response.Name)
		assert.Equal(suite.T(), req.Description, *response.Description)
		assert.Equal(suite.T(), req.Section, response.Section)
	})

	suite.Run("Should fail if permission already exists", func() {
		req := &permissionDomain.CreatePermissionRequest{
			Name:        "existing-permission",
			Description: "Existing permission",
			Section:     "test",
		}

		// Mock: el permiso ya existe
		suite.mockRepo.On("ExistsByName", suite.ctx, req.Name).Return(true, nil)

		response, err := suite.permissionService.CreatePermission(suite.ctx, req)

		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), response)
		assert.Contains(suite.T(), err.Error(), "ya existe un permiso con el nombre")
	})

	suite.Run("Should fail if repository error on exists check", func() {
		req := &permissionDomain.CreatePermissionRequest{
			Name:        "test-permission",
			Description: "Test permission",
			Section:     "test",
		}

		repoError := errors.New("database connection error")
		suite.mockRepo.On("ExistsByName", suite.ctx, req.Name).Return(false, repoError)

		response, err := suite.permissionService.CreatePermission(suite.ctx, req)

		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), response)
		assert.Contains(suite.T(), err.Error(), "error verificando existencia de permiso")
	})

	suite.Run("Should fail if repository error on create", func() {
		req := &permissionDomain.CreatePermissionRequest{
			Name:        "test-permission",
			Description: "Test permission",
			Section:     "test",
		}

		suite.mockRepo.On("ExistsByName", suite.ctx, req.Name).Return(false, nil)

		repoError := errors.New("database insert error")
		suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*domain.Permission")).
			Return(repoError)

		response, err := suite.permissionService.CreatePermission(suite.ctx, req)

		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), response)
		assert.Contains(suite.T(), err.Error(), "error creando permiso")
	})
}

// TestGetPermissionByID valida la obtención de permisos por ID
func (suite *PermissionServiceTestSuite) TestGetPermissionByID() {
	suite.Run("Should get permission by ID successfully", func() {
		permissionID := uuid.New()
		expectedPermission := permissionDomain.NewPermission("test-permission", "Test description", "test")
		expectedPermission.ID = permissionID

		suite.mockRepo.On("GetByID", suite.ctx, permissionID).Return(expectedPermission, nil)

		response, err := suite.permissionService.GetPermissionByID(suite.ctx, permissionID)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), response)
		assert.Equal(suite.T(), permissionID, response.ID)
		assert.Equal(suite.T(), "test-permission", response.Name)
	})

	suite.Run("Should fail if permission not found", func() {
		permissionID := uuid.New()
		repoError := errors.New("permission not found")

		suite.mockRepo.On("GetByID", suite.ctx, permissionID).Return(nil, repoError)

		response, err := suite.permissionService.GetPermissionByID(suite.ctx, permissionID)

		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), response)
		assert.Contains(suite.T(), err.Error(), "error obteniendo permiso")
	})
}

// TestGetPermissions valida la obtención de permisos con filtros
func (suite *PermissionServiceTestSuite) TestGetPermissions() {
	suite.Run("Should get permissions with filter successfully", func() {
		filter := &permissionDomain.PermissionFilterRequest{
			Search:  "test",
			Section: "user",
			Page:    1,
			Limit:   10,
		}

		permissions := []*permissionDomain.Permission{
			permissionDomain.NewPermission("test-permission-1", "Test 1", "user"),
			permissionDomain.NewPermission("test-permission-2", "Test 2", "user"),
		}

		suite.mockRepo.On("GetAllPermissions", suite.ctx, filter).
			Return(permissions, 2, nil)

		response, err := suite.permissionService.GetPermissions(suite.ctx, filter)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), response)
		assert.Len(suite.T(), response.Permissions, 2)
		assert.Equal(suite.T(), 2, response.Total)
		assert.Equal(suite.T(), 1, response.Page)
		assert.Equal(suite.T(), 10, response.Limit)
	})

	suite.Run("Should handle empty result", func() {
		filter := &permissionDomain.PermissionFilterRequest{
			Search: "nonexistent",
			Page:   1,
			Limit:  10,
		}

		suite.mockRepo.On("GetAllPermissions", suite.ctx, filter).
			Return([]*permissionDomain.Permission{}, 0, nil)

		response, err := suite.permissionService.GetPermissions(suite.ctx, filter)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), response)
		assert.Empty(suite.T(), response.Permissions)
		assert.Equal(suite.T(), 0, response.Total)
	})
}

// TestUpdatePermission valida la actualización de permisos
func (suite *PermissionServiceTestSuite) TestUpdatePermission() {
	suite.Run("Should update permission successfully", func() {
		permissionID := uuid.New()
		existingPermission := permissionDomain.NewPermission("old-permission", "Old description", "old")
		existingPermission.ID = permissionID

		updateReq := &permissionDomain.UpdatePermissionRequest{
			Name:        "updated-permission",
			Description: "Updated description",
			Section:     "updated",
		}

		// Mock: obtener permiso existente
		suite.mockRepo.On("GetByID", suite.ctx, permissionID).Return(existingPermission, nil)

		// Mock: verificar que no es permiso del sistema
		suite.mockRepo.On("IsSystemPermission", suite.ctx, permissionID).Return(false, nil)

		// Mock: actualizar permiso
		suite.mockRepo.On("Update", suite.ctx, mock.AnythingOfType("*domain.Permission")).
			Return(nil)

		response, err := suite.permissionService.UpdatePermission(suite.ctx, permissionID, updateReq)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), response)
		assert.Equal(suite.T(), "updated-permission", response.Name)
		assert.Equal(suite.T(), "Updated description", *response.Description)
		assert.Equal(suite.T(), "updated", response.Section)
	})

	suite.Run("Should fail if permission not found", func() {
		permissionID := uuid.New()
		updateReq := &permissionDomain.UpdatePermissionRequest{
			Name:        "updated-permission",
			Description: "Updated description",
			Section:     "updated",
		}

		suite.mockRepo.On("GetByID", suite.ctx, permissionID).Return(nil, errors.New("not found"))

		response, err := suite.permissionService.UpdatePermission(suite.ctx, permissionID, updateReq)

		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), response)
	})

	suite.Run("Should fail if trying to update system permission", func() {
		permissionID := uuid.New()
		systemPermission := permissionDomain.NewPermission(permissionDomain.PermissionCreateUser, "Create users", string(permissionDomain.SectionUser))
		systemPermission.ID = permissionID

		updateReq := &permissionDomain.UpdatePermissionRequest{
			Name:        "hacked-permission",
			Description: "Malicious update",
			Section:     "malicious",
		}

		suite.mockRepo.On("GetByID", suite.ctx, permissionID).Return(systemPermission, nil)
		suite.mockRepo.On("IsSystemPermission", suite.ctx, permissionID).Return(true, nil)

		response, err := suite.permissionService.UpdatePermission(suite.ctx, permissionID, updateReq)

		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), response)
		assert.Contains(suite.T(), err.Error(), "los permisos del sistema no pueden ser modificados")
	})
}

// TestDeletePermission valida la eliminación de permisos
func (suite *PermissionServiceTestSuite) TestDeletePermission() {
	suite.Run("Should delete permission successfully", func() {
		permissionID := uuid.New()

		// Mock: verificar que no es permiso del sistema
		suite.mockRepo.On("IsSystemPermission", suite.ctx, permissionID).Return(false, nil)

		// Mock: eliminar permiso
		suite.mockRepo.On("Delete", suite.ctx, permissionID).Return(nil)

		err := suite.permissionService.DeletePermission(suite.ctx, permissionID)

		assert.NoError(suite.T(), err)
	})

	suite.Run("Should fail if trying to delete system permission", func() {
		permissionID := uuid.New()

		suite.mockRepo.On("IsSystemPermission", suite.ctx, permissionID).Return(true, nil)

		err := suite.permissionService.DeletePermission(suite.ctx, permissionID)

		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "los permisos del sistema no pueden ser eliminados")
	})

	suite.Run("Should fail if repository error", func() {
		permissionID := uuid.New()
		repoError := errors.New("database error")

		suite.mockRepo.On("IsSystemPermission", suite.ctx, permissionID).Return(false, nil)
		suite.mockRepo.On("Delete", suite.ctx, permissionID).Return(repoError)

		err := suite.permissionService.DeletePermission(suite.ctx, permissionID)

		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "error eliminando permiso")
	})
}

// TestGetPermissionsGroupedBySection valida la obtención de permisos agrupados
func (suite *PermissionServiceTestSuite) TestGetPermissionsGroupedBySection() {
	suite.Run("Should get permissions grouped by section successfully", func() {
		groupedPermissions := map[string][]*permissionDomain.Permission{
			"user": {
				permissionDomain.NewPermission("create-user", "Create users", "user"),
				permissionDomain.NewPermission("update-user", "Update users", "user"),
			},
			"role": {
				permissionDomain.NewPermission("create-role", "Create roles", "role"),
			},
		}

		suite.mockRepo.On("GetGroupedBySection", suite.ctx).Return(groupedPermissions, nil)

		response, err := suite.permissionService.GetPermissionsGroupedBySection(suite.ctx)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), response)
		assert.Len(suite.T(), response, 2)

		// Verificar estructura de respuesta
		var userSection, roleSection *permissionDomain.PermissionGroupedResponse
		for _, section := range response {
			if section.Section == "user" {
				userSection = section
			} else if section.Section == "role" {
				roleSection = section
			}
		}

		assert.NotNil(suite.T(), userSection)
		assert.Len(suite.T(), userSection.Permissions, 2)
		assert.NotNil(suite.T(), roleSection)
		assert.Len(suite.T(), roleSection.Permissions, 1)
	})
}

// TestGetAvailableSections valida la obtención de secciones disponibles
func (suite *PermissionServiceTestSuite) TestGetAvailableSections() {
	suite.Run("Should get available sections successfully", func() {
		expectedSections := []string{"user", "role", "permission", "tenant"}

		suite.mockRepo.On("GetAvailableSections", suite.ctx).Return(expectedSections, nil)

		response, err := suite.permissionService.GetAvailableSections(suite.ctx)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), response)
		assert.ElementsMatch(suite.T(), expectedSections, response.Sections)
	})
}

// TestValidationMethods valida los métodos de validación
func (suite *PermissionServiceTestSuite) TestValidationMethods() {
	suite.Run("ValidatePermissionCreation should validate required fields", func() {
		// Validación exitosa
		validReq := &permissionDomain.CreatePermissionRequest{
			Name:        "valid-permission",
			Description: "Valid description",
			Section:     "valid",
		}

		err := suite.permissionService.ValidatePermissionCreation(suite.ctx, validReq)
		assert.NoError(suite.T(), err)

		// Nombre vacío
		invalidReq := &permissionDomain.CreatePermissionRequest{
			Name:        "",
			Description: "Valid description",
			Section:     "valid",
		}

		err = suite.permissionService.ValidatePermissionCreation(suite.ctx, invalidReq)
		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "nombre del permiso es requerido")

		// Sección vacía
		invalidReq2 := &permissionDomain.CreatePermissionRequest{
			Name:        "valid-name",
			Description: "Valid description",
			Section:     "",
		}

		err = suite.permissionService.ValidatePermissionCreation(suite.ctx, invalidReq2)
		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "sección del permiso es requerida")
	})
}

// TestCheckPermissionExistsByID valida la verificación de existencia por ID
func (suite *PermissionServiceTestSuite) TestCheckPermissionExistsByID() {
	suite.Run("Should return true if permission exists", func() {
		permissionID := uuid.New()
		permission := permissionDomain.NewPermission("test", "test", "test")

		suite.mockRepo.On("GetByID", suite.ctx, permissionID).Return(permission, nil)

		exists, err := suite.permissionService.CheckPermissionExistsByID(suite.ctx, permissionID.String())

		assert.NoError(suite.T(), err)
		assert.True(suite.T(), exists)
	})

	suite.Run("Should return false if permission does not exist", func() {
		permissionID := uuid.New()

		suite.mockRepo.On("GetByID", suite.ctx, permissionID).Return(nil, errors.New("not found"))

		exists, err := suite.permissionService.CheckPermissionExistsByID(suite.ctx, permissionID.String())

		assert.NoError(suite.T(), err)
		assert.False(suite.T(), exists)
	})

	suite.Run("Should fail with invalid UUID", func() {
		exists, err := suite.permissionService.CheckPermissionExistsByID(suite.ctx, "invalid-uuid")

		assert.Error(suite.T(), err)
		assert.False(suite.T(), exists)
		assert.Contains(suite.T(), err.Error(), "ID de permiso inválido")
	})
}

// Ejecutar la suite de tests
func TestPermissionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionServiceTestSuite))
}
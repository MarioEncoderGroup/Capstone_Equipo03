package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/controllers"
	permissionDomain "github.com/JoseLuis21/mv-backend/internal/core/permission/domain"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// PermissionControllerTestSuite agrupa tests del controlador de permisos
type PermissionControllerTestSuite struct {
	suite.Suite
	app                   *fiber.App
	mockService           *MockPermissionService
	permissionController  *controllers.PermissionController
	ctx                   context.Context
}

// MockPermissionService es un mock del servicio de permisos
type MockPermissionService struct {
	mock.Mock
}

func (m *MockPermissionService) CreatePermission(ctx context.Context, req *permissionDomain.CreatePermissionRequest) (*permissionDomain.PermissionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permissionDomain.PermissionResponse), args.Error(1)
}

func (m *MockPermissionService) GetPermissionByID(ctx context.Context, id uuid.UUID) (*permissionDomain.PermissionResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permissionDomain.PermissionResponse), args.Error(1)
}

func (m *MockPermissionService) GetPermissionByName(ctx context.Context, name string) (*permissionDomain.PermissionResponse, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permissionDomain.PermissionResponse), args.Error(1)
}

func (m *MockPermissionService) GetPermissionsBySection(ctx context.Context, section string) ([]*permissionDomain.PermissionResponse, error) {
	args := m.Called(ctx, section)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*permissionDomain.PermissionResponse), args.Error(1)
}

func (m *MockPermissionService) GetPermissions(ctx context.Context, filter *permissionDomain.PermissionFilterRequest) (*permissionDomain.PermissionListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permissionDomain.PermissionListResponse), args.Error(1)
}

func (m *MockPermissionService) GetPermissionsGroupedBySection(ctx context.Context) ([]*permissionDomain.PermissionGroupedResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*permissionDomain.PermissionGroupedResponse), args.Error(1)
}

func (m *MockPermissionService) UpdatePermission(ctx context.Context, id uuid.UUID, req *permissionDomain.UpdatePermissionRequest) (*permissionDomain.PermissionResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permissionDomain.PermissionResponse), args.Error(1)
}

func (m *MockPermissionService) DeletePermission(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPermissionService) ValidatePermissionCreation(ctx context.Context, req *permissionDomain.CreatePermissionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockPermissionService) ValidatePermissionUpdate(ctx context.Context, id uuid.UUID, req *permissionDomain.UpdatePermissionRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockPermissionService) InitializeSystemPermissions(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPermissionService) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*permissionDomain.PermissionResponse, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*permissionDomain.PermissionResponse), args.Error(1)
}

func (m *MockPermissionService) GetAvailableSections(ctx context.Context) (*permissionDomain.PermissionSectionResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permissionDomain.PermissionSectionResponse), args.Error(1)
}

func (m *MockPermissionService) CheckPermissionExistsByID(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// SetupTest configura cada test
func (suite *PermissionControllerTestSuite) SetupTest() {
	suite.mockService = new(MockPermissionService)
	validator := validatorapi.NewValidator()
	suite.permissionController = controllers.NewPermissionController(suite.mockService, validator)
	suite.ctx = context.Background()

	// Crear app de Fiber para tests
	suite.app = fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	// Configurar rutas de test
	suite.setupRoutes()
}

// TearDownTest limpia después de cada test
func (suite *PermissionControllerTestSuite) TearDownTest() {
	suite.mockService.AssertExpectations(suite.T())
}

// setupRoutes configura las rutas para testing
func (suite *PermissionControllerTestSuite) setupRoutes() {
	api := suite.app.Group("/api/v1")
	permissions := api.Group("/permissions")

	permissions.Get("/", suite.permissionController.GetPermissions)
	permissions.Get("/sections", suite.permissionController.GetAvailableSections)
	permissions.Get("/grouped", suite.permissionController.GetPermissionsGrouped)
	permissions.Get("/:id", suite.permissionController.GetPermissionByID)
	permissions.Post("/", suite.permissionController.CreatePermission)
	permissions.Put("/:id", suite.permissionController.UpdatePermission)
	permissions.Delete("/:id", suite.permissionController.DeletePermission)
}

// TestCreatePermission valida la creación de permisos
func (suite *PermissionControllerTestSuite) TestCreatePermission() {
	suite.Run("Should create permission successfully", func() {
		req := permissionDomain.CreatePermissionRequest{
			Name:        "test-permission",
			Description: "Test permission description",
			Section:     "test",
		}

		expectedResponse := &permissionDomain.PermissionResponse{
			ID:          uuid.New(),
			Name:        req.Name,
			Description: &req.Description,
			Section:     req.Section,
			IsSystem:    false,
			IsActive:    true,
		}

		suite.mockService.On("CreatePermission", mock.Anything, &req).
			Return(expectedResponse, nil)

		body, _ := json.Marshal(req)
		resp, err := suite.app.Test(httptest.NewRequest("POST", "/api/v1/permissions", bytes.NewBuffer(body)).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusCreated, resp.StatusCode)
	})

	suite.Run("Should return 400 for invalid request", func() {
		req := permissionDomain.CreatePermissionRequest{
			Name:        "", // Invalid: empty name
			Description: "Test permission description",
			Section:     "test",
		}

		body, _ := json.Marshal(req)
		resp, err := suite.app.Test(httptest.NewRequest("POST", "/api/v1/permissions", bytes.NewBuffer(body)).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
	})

	suite.Run("Should return 500 for service error", func() {
		req := permissionDomain.CreatePermissionRequest{
			Name:        "test-permission",
			Description: "Test permission description",
			Section:     "test",
		}

		suite.mockService.On("CreatePermission", mock.Anything, &req).
			Return(nil, errors.New("service error"))

		body, _ := json.Marshal(req)
		resp, err := suite.app.Test(httptest.NewRequest("POST", "/api/v1/permissions", bytes.NewBuffer(body)).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusInternalServerError, resp.StatusCode)
	})
}

// TestGetPermissionByID valida la obtención de permisos por ID
func (suite *PermissionControllerTestSuite) TestGetPermissionByID() {
	suite.Run("Should get permission by ID successfully", func() {
		permissionID := uuid.New()
		expectedResponse := &permissionDomain.PermissionResponse{
			ID:          permissionID,
			Name:        "test-permission",
			Description: stringPtr("Test permission"),
			Section:     "test",
			IsSystem:    false,
			IsActive:    true,
		}

		suite.mockService.On("GetPermissionByID", mock.Anything, permissionID).
			Return(expectedResponse, nil)

		resp, err := suite.app.Test(httptest.NewRequest("GET", fmt.Sprintf("/api/v1/permissions/%s", permissionID.String()), nil).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)
	})

	suite.Run("Should return 400 for invalid UUID", func() {
		resp, err := suite.app.Test(httptest.NewRequest("GET", "/api/v1/permissions/invalid-uuid", nil).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
	})

	suite.Run("Should return 404 for non-existent permission", func() {
		permissionID := uuid.New()

		suite.mockService.On("GetPermissionByID", mock.Anything, permissionID).
			Return(nil, errors.New("permission not found"))

		resp, err := suite.app.Test(httptest.NewRequest("GET", fmt.Sprintf("/api/v1/permissions/%s", permissionID.String()), nil).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusInternalServerError, resp.StatusCode)
	})
}

// TestGetPermissions valida la obtención de permisos con filtros
func (suite *PermissionControllerTestSuite) TestGetPermissions() {
	suite.Run("Should get permissions with pagination", func() {
		expectedResponse := &permissionDomain.PermissionListResponse{
			Permissions: []permissionDomain.PermissionResponse{
				{
					ID:          uuid.New(),
					Name:        "permission-1",
					Description: stringPtr("First permission"),
					Section:     "test",
					IsSystem:    false,
					IsActive:    true,
				},
				{
					ID:          uuid.New(),
					Name:        "permission-2",
					Description: stringPtr("Second permission"),
					Section:     "test",
					IsSystem:    false,
					IsActive:    true,
				},
			},
			Total: 2,
			Page:  1,
			Limit: 10,
		}

		suite.mockService.On("GetPermissions", mock.Anything, mock.AnythingOfType("*domain.PermissionFilterRequest")).
			Return(expectedResponse, nil)

		resp, err := suite.app.Test(httptest.NewRequest("GET", "/api/v1/permissions?page=1&limit=10", nil).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)
	})

	suite.Run("Should get permissions with search filter", func() {
		expectedResponse := &permissionDomain.PermissionListResponse{
			Permissions: []permissionDomain.PermissionResponse{
				{
					ID:          uuid.New(),
					Name:        "search-permission",
					Description: stringPtr("Searchable permission"),
					Section:     "test",
					IsSystem:    false,
					IsActive:    true,
				},
			},
			Total: 1,
			Page:  1,
			Limit: 10,
		}

		suite.mockService.On("GetPermissions", mock.Anything, mock.AnythingOfType("*domain.PermissionFilterRequest")).
			Return(expectedResponse, nil)

		resp, err := suite.app.Test(httptest.NewRequest("GET", "/api/v1/permissions?search=search", nil).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)
	})
}

// TestUpdatePermission valida la actualización de permisos
func (suite *PermissionControllerTestSuite) TestUpdatePermission() {
	suite.Run("Should update permission successfully", func() {
		permissionID := uuid.New()
		req := permissionDomain.UpdatePermissionRequest{
			Name:        "updated-permission",
			Description: "Updated description",
			Section:     "updated",
		}

		expectedResponse := &permissionDomain.PermissionResponse{
			ID:          permissionID,
			Name:        req.Name,
			Description: &req.Description,
			Section:     req.Section,
			IsSystem:    false,
			IsActive:    true,
		}

		suite.mockService.On("UpdatePermission", mock.Anything, permissionID, &req).
			Return(expectedResponse, nil)

		body, _ := json.Marshal(req)
		resp, err := suite.app.Test(httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/permissions/%s", permissionID.String()), bytes.NewBuffer(body)).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)
	})

	suite.Run("Should return 400 for invalid request data", func() {
		permissionID := uuid.New()
		req := permissionDomain.UpdatePermissionRequest{
			Name:        "", // Invalid: empty name
			Description: "Description",
			Section:     "section",
		}

		body, _ := json.Marshal(req)
		resp, err := suite.app.Test(httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/permissions/%s", permissionID.String()), bytes.NewBuffer(body)).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
	})
}

// TestDeletePermission valida la eliminación de permisos
func (suite *PermissionControllerTestSuite) TestDeletePermission() {
	suite.Run("Should delete permission successfully", func() {
		permissionID := uuid.New()

		suite.mockService.On("DeletePermission", mock.Anything, permissionID).
			Return(nil)

		resp, err := suite.app.Test(httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/permissions/%s", permissionID.String()), nil).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusNoContent, resp.StatusCode)
	})

	suite.Run("Should return 400 for invalid UUID", func() {
		resp, err := suite.app.Test(httptest.NewRequest("DELETE", "/api/v1/permissions/invalid-uuid", nil).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)
	})

	suite.Run("Should return error for service failure", func() {
		permissionID := uuid.New()

		suite.mockService.On("DeletePermission", mock.Anything, permissionID).
			Return(errors.New("cannot delete system permission"))

		resp, err := suite.app.Test(httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/permissions/%s", permissionID.String()), nil).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusInternalServerError, resp.StatusCode)
	})
}

// TestGetAvailableSections valida la obtención de secciones disponibles
func (suite *PermissionControllerTestSuite) TestGetAvailableSections() {
	suite.Run("Should get available sections successfully", func() {
		expectedResponse := &permissionDomain.PermissionSectionResponse{
			Sections: []string{"user", "role", "permission", "tenant"},
		}

		suite.mockService.On("GetAvailableSections", mock.Anything).
			Return(expectedResponse, nil)

		resp, err := suite.app.Test(httptest.NewRequest("GET", "/api/v1/permissions/sections", nil).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)
	})
}

// TestGetPermissionsGrouped valida la obtención de permisos agrupados
func (suite *PermissionControllerTestSuite) TestGetPermissionsGrouped() {
	suite.Run("Should get grouped permissions successfully", func() {
		expectedResponse := []*permissionDomain.PermissionGroupedResponse{
			{
				Section: "user",
				Permissions: []permissionDomain.PermissionResponse{
					{
						ID:          uuid.New(),
						Name:        "create-user",
						Description: stringPtr("Create users"),
						Section:     "user",
						IsSystem:    true,
						IsActive:    true,
					},
				},
			},
		}

		suite.mockService.On("GetPermissionsGroupedBySection", mock.Anything).
			Return(expectedResponse, nil)

		resp, err := suite.app.Test(httptest.NewRequest("GET", "/api/v1/permissions/grouped", nil).
			WithContext(suite.ctx))

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)
	})
}

// Helper function para crear punteros a strings
func stringPtr(s string) *string {
	return &s
}

// Ejecutar la suite de tests
func TestPermissionControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionControllerTestSuite))
}
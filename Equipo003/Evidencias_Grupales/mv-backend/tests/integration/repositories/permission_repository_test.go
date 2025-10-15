package repositories_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	permissionAdapters "github.com/JoseLuis21/mv-backend/internal/core/permission/adapters"
	permissionDomain "github.com/JoseLuis21/mv-backend/internal/core/permission/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/permission/ports"
	"github.com/JoseLuis21/mv-backend/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// PermissionRepositoryIntegrationTestSuite agrupa tests de integración del repositorio de permisos
type PermissionRepositoryIntegrationTestSuite struct {
	suite.Suite
	repo    ports.PermissionRepository
	ctx     context.Context
	cleanup func()
	dbName  string
}

// SetupSuite configura la suite de tests
func (suite *PermissionRepositoryIntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Crear base de datos temporal para testing
	db, dbName, cleanup := helpers.CreateTestDatabase(suite.T())
	if db == nil {
		suite.T().Skip("Database not available for integration tests")
		return
	}

	suite.cleanup = cleanup
	suite.dbName = dbName

	// Crear esquema de permisos para testing
	suite.createPermissionSchema(db)

	// Inicializar repositorio
	suite.repo = permissionAdapters.NewPgPermissionRepository(db)
}

// TearDownSuite limpia después de la suite
func (suite *PermissionRepositoryIntegrationTestSuite) TearDownSuite() {
	if suite.cleanup != nil {
		suite.cleanup()
	}
}

// createPermissionSchema crea el esquema necesario para tests
func (suite *PermissionRepositoryIntegrationTestSuite) createPermissionSchema(db interface{}) {
	// En un test real, aquí se ejecutarían las migraciones necesarias
	// Por ahora, asumimos que las tablas ya existen o se crean automáticamente
}

// TestCreateAndGetPermission valida la creación y obtención de permisos
func (suite *PermissionRepositoryIntegrationTestSuite) TestCreateAndGetPermission() {
	suite.Run("Should create and retrieve permission by ID", func() {
		permission := permissionDomain.NewPermission("test-create-permission", "Test permission for create", "test")

		// Crear permiso
		err := suite.repo.Create(suite.ctx, permission)
		require.NoError(suite.T(), err)

		// Obtener por ID
		retrieved, err := suite.repo.GetByID(suite.ctx, permission.ID)
		require.NoError(suite.T(), err)
		assert.NotNil(suite.T(), retrieved)
		assert.Equal(suite.T(), permission.ID, retrieved.ID)
		assert.Equal(suite.T(), permission.Name, retrieved.Name)
		assert.Equal(suite.T(), permission.Section, retrieved.Section)
	})

	suite.Run("Should create and retrieve permission by name", func() {
		permission := permissionDomain.NewPermission("test-get-by-name", "Test permission for name retrieval", "test")

		// Crear permiso
		err := suite.repo.Create(suite.ctx, permission)
		require.NoError(suite.T(), err)

		// Obtener por nombre
		retrieved, err := suite.repo.GetByName(suite.ctx, permission.Name)
		require.NoError(suite.T(), err)
		assert.NotNil(suite.T(), retrieved)
		assert.Equal(suite.T(), permission.Name, retrieved.Name)
		assert.Equal(suite.T(), permission.Section, retrieved.Section)
	})

	suite.Run("Should return error for non-existent permission", func() {
		nonExistentID := uuid.New()

		retrieved, err := suite.repo.GetByID(suite.ctx, nonExistentID)
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), retrieved)
	})
}

// TestExistsByName valida la verificación de existencia por nombre
func (suite *PermissionRepositoryIntegrationTestSuite) TestExistsByName() {
	suite.Run("Should return true for existing permission", func() {
		permission := permissionDomain.NewPermission("test-exists-permission", "Test permission for exists", "test")

		// Crear permiso
		err := suite.repo.Create(suite.ctx, permission)
		require.NoError(suite.T(), err)

		// Verificar existencia
		exists, err := suite.repo.ExistsByName(suite.ctx, permission.Name)
		require.NoError(suite.T(), err)
		assert.True(suite.T(), exists)
	})

	suite.Run("Should return false for non-existent permission", func() {
		exists, err := suite.repo.ExistsByName(suite.ctx, "non-existent-permission")
		require.NoError(suite.T(), err)
		assert.False(suite.T(), exists)
	})
}

// TestGetBySection valida la obtención de permisos por sección
func (suite *PermissionRepositoryIntegrationTestSuite) TestGetBySection() {
	suite.Run("Should get permissions by section", func() {
		section := "test-section"

		// Crear múltiples permisos en la misma sección
		permission1 := permissionDomain.NewPermission("test-section-perm-1", "Test permission 1", section)
		permission2 := permissionDomain.NewPermission("test-section-perm-2", "Test permission 2", section)

		err := suite.repo.Create(suite.ctx, permission1)
		require.NoError(suite.T(), err)

		err = suite.repo.Create(suite.ctx, permission2)
		require.NoError(suite.T(), err)

		// Obtener permisos por sección
		permissions, err := suite.repo.GetBySection(suite.ctx, section)
		require.NoError(suite.T(), err)
		assert.GreaterOrEqual(suite.T(), len(permissions), 2)

		// Verificar que todos los permisos pertenecen a la sección
		for _, perm := range permissions {
			assert.Equal(suite.T(), section, perm.Section)
		}
	})

	suite.Run("Should return empty slice for non-existent section", func() {
		permissions, err := suite.repo.GetBySection(suite.ctx, "non-existent-section")
		require.NoError(suite.T(), err)
		assert.Empty(suite.T(), permissions)
	})
}

// TestGetAllPermissions valida la obtención de permisos con filtros
func (suite *PermissionRepositoryIntegrationTestSuite) TestGetAllPermissions() {
	suite.Run("Should get all permissions with pagination", func() {
		// Crear algunos permisos de prueba
		testSection := "pagination-test"
		for i := 0; i < 5; i++ {
			permission := permissionDomain.NewPermission(
				fmt.Sprintf("pagination-perm-%d", i),
				fmt.Sprintf("Pagination test permission %d", i),
				testSection,
			)
			err := suite.repo.Create(suite.ctx, permission)
			require.NoError(suite.T(), err)
		}

		filter := &permissionDomain.PermissionFilterRequest{
			Section: testSection,
			Page:    1,
			Limit:   3,
		}

		permissions, total, err := suite.repo.GetAllPermissions(suite.ctx, filter)
		require.NoError(suite.T(), err)
		assert.LessOrEqual(suite.T(), len(permissions), 3)
		assert.GreaterOrEqual(suite.T(), total, 5)
	})

	suite.Run("Should filter permissions by search term", func() {
		// Crear permisos con nombres específicos para búsqueda
		searchTerm := "searchable"
		permission1 := permissionDomain.NewPermission("searchable-permission-1", "First searchable permission", "search")
		permission2 := permissionDomain.NewPermission("not-matching", "Not matching permission", "search")

		err := suite.repo.Create(suite.ctx, permission1)
		require.NoError(suite.T(), err)

		err = suite.repo.Create(suite.ctx, permission2)
		require.NoError(suite.T(), err)

		filter := &permissionDomain.PermissionFilterRequest{
			Search: searchTerm,
			Page:   1,
			Limit:  10,
		}

		permissions, total, err := suite.repo.GetAllPermissions(suite.ctx, filter)
		require.NoError(suite.T(), err)
		assert.GreaterOrEqual(suite.T(), total, 1)

		// Verificar que los resultados contienen el término de búsqueda
		found := false
		for _, perm := range permissions {
			if strings.Contains(strings.ToLower(perm.Name), strings.ToLower(searchTerm)) {
				found = true
				break
			}
		}
		assert.True(suite.T(), found, "Should find permissions matching search term")
	})
}

// TestUpdatePermission valida la actualización de permisos
func (suite *PermissionRepositoryIntegrationTestSuite) TestUpdatePermission() {
	suite.Run("Should update permission successfully", func() {
		permission := permissionDomain.NewPermission("test-update", "Original description", "original")

		// Crear permiso
		err := suite.repo.Create(suite.ctx, permission)
		require.NoError(suite.T(), err)

		// Actualizar permiso
		permission.Update("test-updated", "Updated description", "updated")
		err = suite.repo.Update(suite.ctx, permission)
		require.NoError(suite.T(), err)

		// Verificar actualización
		retrieved, err := suite.repo.GetByID(suite.ctx, permission.ID)
		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), "test-updated", retrieved.Name)
		assert.Equal(suite.T(), "Updated description", *retrieved.Description)
		assert.Equal(suite.T(), "updated", retrieved.Section)
	})
}

// TestSoftDelete valida la eliminación lógica
func (suite *PermissionRepositoryIntegrationTestSuite) TestSoftDelete() {
	suite.Run("Should soft delete permission", func() {
		permission := permissionDomain.NewPermission("test-delete", "Permission to delete", "test")

		// Crear permiso
		err := suite.repo.Create(suite.ctx, permission)
		require.NoError(suite.T(), err)

		// Verificar que existe
		exists, err := suite.repo.ExistsByName(suite.ctx, permission.Name)
		require.NoError(suite.T(), err)
		assert.True(suite.T(), exists)

		// Eliminar lógicamente
		err = suite.repo.Delete(suite.ctx, permission.ID)
		require.NoError(suite.T(), err)

		// Verificar que ya no existe (soft delete)
		exists, err = suite.repo.ExistsByName(suite.ctx, permission.Name)
		require.NoError(suite.T(), err)
		assert.False(suite.T(), exists)

		// Verificar que GetByID también devuelve error
		_, err = suite.repo.GetByID(suite.ctx, permission.ID)
		assert.Error(suite.T(), err)
	})
}

// TestGetGroupedBySection valida la obtención agrupada por sección
func (suite *PermissionRepositoryIntegrationTestSuite) TestGetGroupedBySection() {
	suite.Run("Should get permissions grouped by section", func() {
		// Crear permisos en diferentes secciones
		sections := []string{"group-test-1", "group-test-2"}

		for _, section := range sections {
			for i := 0; i < 2; i++ {
				permission := permissionDomain.NewPermission(
					fmt.Sprintf("%s-perm-%d", section, i),
					fmt.Sprintf("Permission %d for section %s", i, section),
					section,
				)
				err := suite.repo.Create(suite.ctx, permission)
				require.NoError(suite.T(), err)
			}
		}

		grouped, err := suite.repo.GetGroupedBySection(suite.ctx)
		require.NoError(suite.T(), err)
		assert.NotEmpty(suite.T(), grouped)

		// Verificar que las secciones están presentes
		for _, section := range sections {
			permissions, exists := grouped[section]
			assert.True(suite.T(), exists, "Section %s should exist", section)
			assert.GreaterOrEqual(suite.T(), len(permissions), 2, "Section %s should have at least 2 permissions", section)
		}
	})
}

// TestGetAvailableSections valida la obtención de secciones disponibles
func (suite *PermissionRepositoryIntegrationTestSuite) TestGetAvailableSections() {
	suite.Run("Should get available sections", func() {
		// Crear permisos en secciones únicas
		uniqueSections := []string{"section-1", "section-2", "section-3"}

		for _, section := range uniqueSections {
			permission := permissionDomain.NewPermission(
				fmt.Sprintf("%s-permission", section),
				fmt.Sprintf("Permission for %s", section),
				section,
			)
			err := suite.repo.Create(suite.ctx, permission)
			require.NoError(suite.T(), err)
		}

		sections, err := suite.repo.GetAvailableSections(suite.ctx)
		require.NoError(suite.T(), err)
		assert.NotEmpty(suite.T(), sections)

		// Verificar que nuestras secciones están incluidas
		sectionSet := make(map[string]bool)
		for _, section := range sections {
			sectionSet[section] = true
		}

		for _, expectedSection := range uniqueSections {
			assert.True(suite.T(), sectionSet[expectedSection], "Section %s should be available", expectedSection)
		}
	})
}

// TestIsSystemPermission valida la identificación de permisos del sistema
func (suite *PermissionRepositoryIntegrationTestSuite) TestIsSystemPermission() {
	suite.Run("Should identify system permissions", func() {
		// Crear un permiso del sistema (basado en sección 'system')
		systemPermission := permissionDomain.NewPermission("system-permission", "System permission", "system")
		err := suite.repo.Create(suite.ctx, systemPermission)
		require.NoError(suite.T(), err)

		// Crear un permiso regular
		regularPermission := permissionDomain.NewPermission("regular-permission", "Regular permission", "custom")
		err = suite.repo.Create(suite.ctx, regularPermission)
		require.NoError(suite.T(), err)

		// Verificar identificación de permiso del sistema
		isSystem, err := suite.repo.IsSystemPermission(suite.ctx, systemPermission.ID)
		require.NoError(suite.T(), err)
		assert.True(suite.T(), isSystem)

		// Verificar que el permiso regular no es del sistema
		isSystem, err = suite.repo.IsSystemPermission(suite.ctx, regularPermission.ID)
		require.NoError(suite.T(), err)
		assert.False(suite.T(), isSystem)
	})
}

// Ejecutar la suite de tests
func TestPermissionRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionRepositoryIntegrationTestSuite))
}

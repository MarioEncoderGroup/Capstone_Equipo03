package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/permission/domain"
	"github.com/stretchr/testify/mock"
)

// MockPermissionRepository es un mock del repositorio de permisos para testing
type MockPermissionRepository struct {
	mock.Mock
}

// Create crea un nuevo permiso en la base de datos
func (m *MockPermissionRepository) Create(ctx context.Context, permission *domain.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

// GetByID obtiene un permiso por su ID
func (m *MockPermissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Permission), args.Error(1)
}

// GetByName obtiene un permiso por su nombre
func (m *MockPermissionRepository) GetByName(ctx context.Context, name string) (*domain.Permission, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Permission), args.Error(1)
}

// GetBySection obtiene todos los permisos de una sección específica
func (m *MockPermissionRepository) GetBySection(ctx context.Context, section string) ([]*domain.Permission, error) {
	args := m.Called(ctx, section)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Permission), args.Error(1)
}

// GetAllPermissions obtiene permisos con filtros de búsqueda
func (m *MockPermissionRepository) GetAllPermissions(ctx context.Context, filter *domain.PermissionFilterRequest) ([]*domain.Permission, int, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Permission), args.Int(1), args.Error(2)
}

// GetGroupedBySection obtiene permisos agrupados por sección
func (m *MockPermissionRepository) GetGroupedBySection(ctx context.Context) (map[string][]*domain.Permission, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string][]*domain.Permission), args.Error(1)
}

// Update actualiza un permiso existente
func (m *MockPermissionRepository) Update(ctx context.Context, permission *domain.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

// Delete elimina lógicamente un permiso
func (m *MockPermissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ExistsByName verifica si existe un permiso con el nombre dado
func (m *MockPermissionRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

// IsSystemPermission verifica si un permiso es uno de los permisos predefinidos del sistema
func (m *MockPermissionRepository) IsSystemPermission(ctx context.Context, permissionID uuid.UUID) (bool, error) {
	args := m.Called(ctx, permissionID)
	return args.Bool(0), args.Error(1)
}

// GetRolePermissions obtiene todos los permisos asignados a un rol
func (m *MockPermissionRepository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*domain.Permission, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Permission), args.Error(1)
}

// GetAvailableSections obtiene todas las secciones de permisos disponibles
func (m *MockPermissionRepository) GetAvailableSections(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}
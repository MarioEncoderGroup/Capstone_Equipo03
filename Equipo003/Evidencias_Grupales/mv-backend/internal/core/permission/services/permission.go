package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/permission/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/permission/ports"
)

// permissionService implementa la lógica de negocio para permisos
type permissionService struct {
	permissionRepo ports.PermissionRepository
}

// NewPermissionService crea una nueva instancia del servicio permission
func NewPermissionService(permissionRepo ports.PermissionRepository) ports.PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
	}
}

// CreatePermission crea un nuevo permiso
func (s *permissionService) CreatePermission(ctx context.Context, req *domain.CreatePermissionRequest) (*domain.PermissionResponse, error) {
	// Validar que no existe un permiso con el mismo nombre
	exists, err := s.permissionRepo.ExistsByName(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("error verificando existencia de permiso: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("ya existe un permiso con el nombre %s", req.Name)
	}

	// Crear entidad Permission
	permission := &domain.Permission{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: &req.Description,
		Section:     req.Section,
	}

	// Crear en la base de datos
	if err := s.permissionRepo.Create(ctx, permission); err != nil {
		return nil, fmt.Errorf("error creando permiso: %w", err)
	}

	return permission.ToResponse(), nil
}

// GetPermissionByID obtiene un permiso por su ID
func (s *permissionService) GetPermissionByID(ctx context.Context, id uuid.UUID) (*domain.PermissionResponse, error) {
	permission, err := s.permissionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo permiso: %w", err)
	}

	return permission.ToResponse(), nil
}

// GetPermissionByName obtiene un permiso por su nombre
func (s *permissionService) GetPermissionByName(ctx context.Context, name string) (*domain.PermissionResponse, error) {
	permission, err := s.permissionRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo permiso: %w", err)
	}

	return permission.ToResponse(), nil
}

// GetPermissionsBySection obtiene todos los permisos de una sección específica
func (s *permissionService) GetPermissionsBySection(ctx context.Context, section string) ([]*domain.PermissionResponse, error) {
	permissions, err := s.permissionRepo.GetBySection(ctx, section)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo permisos por sección: %w", err)
	}

	responses := make([]*domain.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		responses[i] = permission.ToResponse()
	}

	return responses, nil
}

// GetPermissions obtiene permisos con filtros de búsqueda y paginación
func (s *permissionService) GetPermissions(ctx context.Context, filter *domain.PermissionFilterRequest) (*domain.PermissionListResponse, error) {
	permissions, total, err := s.permissionRepo.GetAllPermissions(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo permisos: %w", err)
	}

	responses := make([]domain.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		responses[i] = *permission.ToResponse()
	}

	return &domain.PermissionListResponse{
		Permissions: responses,
		Total:       total,
		Page:        filter.Page,
		Limit:       filter.Limit,
	}, nil
}

// GetPermissionsGroupedBySection obtiene permisos agrupados por sección
func (s *permissionService) GetPermissionsGroupedBySection(ctx context.Context) ([]*domain.PermissionGroupedResponse, error) {
	groupedPerms, err := s.permissionRepo.GetGroupedBySection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo permisos agrupados: %w", err)
	}

	var responses []*domain.PermissionGroupedResponse
	for section, permissions := range groupedPerms {
		permResponses := make([]domain.PermissionResponse, len(permissions))
		for i, perm := range permissions {
			permResponses[i] = *perm.ToResponse()
		}

		responses = append(responses, &domain.PermissionGroupedResponse{
			Section:     section,
			Permissions: permResponses,
		})
	}

	return responses, nil
}

// UpdatePermission actualiza un permiso existente
func (s *permissionService) UpdatePermission(ctx context.Context, id uuid.UUID, req *domain.UpdatePermissionRequest) (*domain.PermissionResponse, error) {
	// Obtener permiso existente
	permission, err := s.permissionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo permiso: %w", err)
	}

	// Verificar si es un permiso del sistema
	isSystem, err := s.permissionRepo.IsSystemPermission(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error verificando si es permiso del sistema: %w", err)
	}
	if isSystem {
		return nil, fmt.Errorf("los permisos del sistema no pueden ser modificados")
	}

	// Actualizar campos
	permission.Name = req.Name
	permission.Description = &req.Description
	permission.Section = req.Section

	// Actualizar en la base de datos
	if err := s.permissionRepo.Update(ctx, permission); err != nil {
		return nil, fmt.Errorf("error actualizando permiso: %w", err)
	}

	return permission.ToResponse(), nil
}

// DeletePermission elimina lógicamente un permiso
func (s *permissionService) DeletePermission(ctx context.Context, id uuid.UUID) error {
	// Verificar si es un permiso del sistema
	isSystem, err := s.permissionRepo.IsSystemPermission(ctx, id)
	if err != nil {
		return fmt.Errorf("error verificando si es permiso del sistema: %w", err)
	}
	if isSystem {
		return fmt.Errorf("los permisos del sistema no pueden ser eliminados")
	}

	if err := s.permissionRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error eliminando permiso: %w", err)
	}

	return nil
}

// ValidatePermissionCreation valida que se puede crear un permiso con los datos dados
func (s *permissionService) ValidatePermissionCreation(ctx context.Context, req *domain.CreatePermissionRequest) error {
	if req.Name == "" {
		return fmt.Errorf("el nombre del permiso es requerido")
	}
	if req.Section == "" {
		return fmt.Errorf("la sección del permiso es requerida")
	}
	return nil
}

// ValidatePermissionUpdate valida que se puede actualizar un permiso con los datos dados
func (s *permissionService) ValidatePermissionUpdate(ctx context.Context, id uuid.UUID, req *domain.UpdatePermissionRequest) error {
	if req.Name == "" {
		return fmt.Errorf("el nombre del permiso es requerido")
	}
	if req.Section == "" {
		return fmt.Errorf("la sección del permiso es requerida")
	}
	return nil
}

// InitializeSystemPermissions crea los permisos predefinidos del sistema si no existen
func (s *permissionService) InitializeSystemPermissions(ctx context.Context) error {
	// TODO: Implementar inicialización de permisos del sistema
	// Esta función debería crear todos los permisos predefinidos del sistema
	return nil
}

// GetRolePermissions obtiene todos los permisos asignados a un rol
func (s *permissionService) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*domain.PermissionResponse, error) {
	permissions, err := s.permissionRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo permisos del rol: %w", err)
	}

	responses := make([]*domain.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		responses[i] = permission.ToResponse()
	}

	return responses, nil
}

// GetAvailableSections obtiene todas las secciones de permisos disponibles
func (s *permissionService) GetAvailableSections(ctx context.Context) (*domain.PermissionSectionResponse, error) {
	sections, err := s.permissionRepo.GetAvailableSections(ctx)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo secciones disponibles: %w", err)
	}

	return &domain.PermissionSectionResponse{
		Sections: sections,
	}, nil
}

// CheckPermissionExistsByID verifica si un permiso existe por ID
func (s *permissionService) CheckPermissionExistsByID(ctx context.Context, id string) (bool, error) {
	permissionID, err := uuid.Parse(id)
	if err != nil {
		return false, fmt.Errorf("ID de permiso inválido: %w", err)
	}

	_, err = s.permissionRepo.GetByID(ctx, permissionID)
	if err != nil {
		return false, nil // Si no se encuentra, no existe
	}

	return true, nil
}
package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domain_role_permission "github.com/JoseLuis21/mv-backend/internal/core/role_permission/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/role_permission/ports"
	role_ports "github.com/JoseLuis21/mv-backend/internal/core/role/ports"
	permission_ports "github.com/JoseLuis21/mv-backend/internal/core/permission/ports"
)

// rolePermissionService implementa la lógica de negocio para role_permission
type rolePermissionService struct {
	rolePermissionRepo ports.RolePermissionRepository
	roleService        role_ports.RoleService
	permissionService  permission_ports.PermissionService
}

// NewRolePermissionService crea una nueva instancia del servicio role_permission
func NewRolePermissionService(
	rolePermissionRepo ports.RolePermissionRepository,
	roleService role_ports.RoleService,
	permissionService permission_ports.PermissionService,
) ports.RolePermissionService {
	return &rolePermissionService{
		rolePermissionRepo: rolePermissionRepo,
		roleService:        roleService,
		permissionService:  permissionService,
	}
}

// CreateRolePermission crea una nueva relación rol-permiso con validaciones
func (s *rolePermissionService) CreateRolePermission(ctx context.Context, req *domain_role_permission.CreateRolePermissionDto) (*domain_role_permission.RolePermissionResponseDto, error) {
	// Validar que el rol existe
	roleExists, err := s.roleService.CheckRoleExistsByID(ctx, req.RoleID.String(), "")
	if err != nil {
		return nil, fmt.Errorf("error verificando rol: %w", err)
	}
	if !roleExists {
		return nil, fmt.Errorf("el rol con ID %s no existe", req.RoleID)
	}

	// Validar que el permiso existe
	permissionExists, err := s.permissionService.CheckPermissionExistsByID(ctx, req.PermissionID.String())
	if err != nil {
		return nil, fmt.Errorf("error verificando permiso: %w", err)
	}
	if !permissionExists {
		return nil, fmt.Errorf("el permiso con ID %s no existe", req.PermissionID)
	}

	// Verificar que la relación no existe ya
	exists, err := s.rolePermissionRepo.Exists(ctx, req.RoleID, req.PermissionID)
	if err != nil {
		return nil, fmt.Errorf("error verificando relación existente: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("la relación rol-permiso ya existe")
	}

	// Crear la relación
	rolePermission := domain_role_permission.NewRolePermission(req.RoleID, req.PermissionID)

	if err := s.rolePermissionRepo.Create(ctx, rolePermission); err != nil {
		return nil, fmt.Errorf("error creando relación rol-permiso: %w", err)
	}

	return &domain_role_permission.RolePermissionResponseDto{
		RolePermission: rolePermission,
		Message:        "Relación rol-permiso creada exitosamente",
	}, nil
}

// DeleteRolePermission elimina una relación rol-permiso específica
func (s *rolePermissionService) DeleteRolePermission(ctx context.Context, roleID, permissionID uuid.UUID) error {
	// Verificar que la relación existe
	exists, err := s.rolePermissionRepo.Exists(ctx, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("error verificando relación: %w", err)
	}
	if !exists {
		return fmt.Errorf("la relación rol-permiso no existe")
	}

	if err := s.rolePermissionRepo.Delete(ctx, roleID, permissionID); err != nil {
		return fmt.Errorf("error eliminando relación rol-permiso: %w", err)
	}

	return nil
}

// GetRolePermissions obtiene todos los permisos de un rol
func (s *rolePermissionService) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]domain_role_permission.RolePermission, error) {
	// Validar que el rol existe
	roleExists, err := s.roleService.CheckRoleExistsByID(ctx, roleID.String(), "")
	if err != nil {
		return nil, fmt.Errorf("error verificando rol: %w", err)
	}
	if !roleExists {
		return nil, fmt.Errorf("el rol con ID %s no existe", roleID)
	}

	rolePermissions, err := s.rolePermissionRepo.GetByRoleID(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo permisos del rol: %w", err)
	}

	return rolePermissions, nil
}

// GetPermissionRoles obtiene todos los roles de un permiso
func (s *rolePermissionService) GetPermissionRoles(ctx context.Context, permissionID uuid.UUID) ([]domain_role_permission.RolePermission, error) {
	// Validar que el permiso existe
	permissionExists, err := s.permissionService.CheckPermissionExistsByID(ctx, permissionID.String())
	if err != nil {
		return nil, fmt.Errorf("error verificando permiso: %w", err)
	}
	if !permissionExists {
		return nil, fmt.Errorf("el permiso con ID %s no existe", permissionID)
	}

	rolePermissions, err := s.rolePermissionRepo.GetByPermissionID(ctx, permissionID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo roles del permiso: %w", err)
	}

	return rolePermissions, nil
}

// SyncRolePermissions sincroniza permisos a un rol con validaciones
func (s *rolePermissionService) SyncRolePermissions(ctx context.Context, req *domain_role_permission.SyncRolePermissionsDto) error {
	// Validar que el rol existe
	roleExists, err := s.roleService.CheckRoleExistsByID(ctx, req.RoleID.String(), "")
	if err != nil {
		return fmt.Errorf("error verificando rol: %w", err)
	}
	if !roleExists {
		return fmt.Errorf("el rol con ID %s no existe", req.RoleID)
	}

	// Validar que todos los permisos existen
	for _, permissionID := range req.PermissionIDs {
		permissionExists, err := s.permissionService.CheckPermissionExistsByID(ctx, permissionID.String())
		if err != nil {
			return fmt.Errorf("error verificando permiso %s: %w", permissionID, err)
		}
		if !permissionExists {
			return fmt.Errorf("el permiso con ID %s no existe", permissionID)
		}
	}

	// Realizar sincronización
	if err := s.rolePermissionRepo.SyncRolePermissions(ctx, req.RoleID, req.PermissionIDs); err != nil {
		return fmt.Errorf("error sincronizando permisos al rol: %w", err)
	}

	return nil
}

// RoleHasPermission verifica si un rol tiene un permiso específico
func (s *rolePermissionService) RoleHasPermission(ctx context.Context, roleID, permissionID uuid.UUID) (bool, error) {
	exists, err := s.rolePermissionRepo.Exists(ctx, roleID, permissionID)
	if err != nil {
		return false, fmt.Errorf("error verificando relación rol-permiso: %w", err)
	}

	return exists, nil
}

// RemoveAllPermissionsFromRole elimina todos los permisos de un rol
func (s *rolePermissionService) RemoveAllPermissionsFromRole(ctx context.Context, roleID uuid.UUID) error {
	// Validar que el rol existe
	roleExists, err := s.roleService.CheckRoleExistsByID(ctx, roleID.String(), "")
	if err != nil {
		return fmt.Errorf("error verificando rol: %w", err)
	}
	if !roleExists {
		return fmt.Errorf("el rol con ID %s no existe", roleID)
	}

	if err := s.rolePermissionRepo.DeleteByRoleID(ctx, roleID); err != nil {
		return fmt.Errorf("error eliminando permisos del rol: %w", err)
	}

	return nil
}

// RemoveRoleFromAllPermissions elimina un rol de todos los permisos
func (s *rolePermissionService) RemoveRoleFromAllPermissions(ctx context.Context, permissionID uuid.UUID) error {
	// Validar que el permiso existe
	permissionExists, err := s.permissionService.CheckPermissionExistsByID(ctx, permissionID.String())
	if err != nil {
		return fmt.Errorf("error verificando permiso: %w", err)
	}
	if !permissionExists {
		return fmt.Errorf("el permiso con ID %s no existe", permissionID)
	}

	if err := s.rolePermissionRepo.DeleteByPermissionID(ctx, permissionID); err != nil {
		return fmt.Errorf("error eliminando roles del permiso: %w", err)
	}

	return nil
}

// GetAllPermissionsByRoleID obtiene información completa de permisos por rol
func (s *rolePermissionService) GetAllPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]domain_role_permission.Permission, error) {
	// Validar que el rol existe
	roleExists, err := s.roleService.CheckRoleExistsByID(ctx, roleID.String(), "")
	if err != nil {
		return nil, fmt.Errorf("error verificando rol: %w", err)
	}
	if !roleExists {
		return nil, fmt.Errorf("el rol con ID %s no existe", roleID)
	}

	permissions, err := s.rolePermissionRepo.GetPermissionsByRoleID(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo permisos del rol: %w", err)
	}

	return permissions, nil
}

// GetAllPermissionsByRoleIDs obtiene permisos por múltiples roles
func (s *rolePermissionService) GetAllPermissionsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]domain_role_permission.Permission, error) {
	if len(roleIDs) == 0 {
		return []domain_role_permission.Permission{}, nil
	}

	// Validar que todos los roles existen
	for _, roleID := range roleIDs {
		roleExists, err := s.roleService.CheckRoleExistsByID(ctx, roleID.String(), "")
		if err != nil {
			return nil, fmt.Errorf("error verificando rol %s: %w", roleID, err)
		}
		if !roleExists {
			return nil, fmt.Errorf("el rol con ID %s no existe", roleID)
		}
	}

	permissions, err := s.rolePermissionRepo.GetPermissionsByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo permisos de los roles: %w", err)
	}

	return permissions, nil
}

// GetPermissionNamesByRoleIDs obtiene solo nombres de permisos por múltiples roles
func (s *rolePermissionService) GetPermissionNamesByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]string, error) {
	if len(roleIDs) == 0 {
		return []string{}, nil
	}

	// Validar que todos los roles existen
	for _, roleID := range roleIDs {
		roleExists, err := s.roleService.CheckRoleExistsByID(ctx, roleID.String(), "")
		if err != nil {
			return nil, fmt.Errorf("error verificando rol %s: %w", roleID, err)
		}
		if !roleExists {
			return nil, fmt.Errorf("el rol con ID %s no existe", roleID)
		}
	}

	permissionNames, err := s.rolePermissionRepo.GetPermissionNamesByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo nombres de permisos de los roles: %w", err)
	}

	return permissionNames, nil
}
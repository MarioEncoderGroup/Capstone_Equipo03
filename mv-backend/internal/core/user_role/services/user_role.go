package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domain_user_role "github.com/JoseLuis21/mv-backend/internal/core/user_role/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/user_role/ports"
	user_ports "github.com/JoseLuis21/mv-backend/internal/core/user/ports"
	role_ports "github.com/JoseLuis21/mv-backend/internal/core/role/ports"
)

// userRoleService implementa la lógica de negocio para user_role
type userRoleService struct {
	userRoleRepo ports.UserRoleRepository
	userService  user_ports.UserService
	roleService  role_ports.RoleService
}

// NewUserRoleService crea una nueva instancia del servicio user_role
func NewUserRoleService(
	userRoleRepo ports.UserRoleRepository,
	userService user_ports.UserService,
	roleService role_ports.RoleService,
) ports.UserRoleService {
	return &userRoleService{
		userRoleRepo: userRoleRepo,
		userService:  userService,
		roleService:  roleService,
	}
}

// CreateUserRole crea una nueva relación usuario-rol con validaciones
func (s *userRoleService) CreateUserRole(ctx context.Context, req *domain_user_role.CreateUserRoleDto) (*domain_user_role.UserRoleResponseDto, error) {
	// Validar que el usuario existe
	userExists, err := s.userService.CheckUserExists(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("error verificando usuario: %w", err)
	}
	if !userExists {
		return nil, fmt.Errorf("el usuario con ID %s no existe", req.UserID)
	}

	// Validar que el rol existe
	tenantID := ""
	if req.TenantID != nil {
		tenantID = req.TenantID.String()
	}

	roleExists, err := s.roleService.CheckRoleExistsByID(ctx, req.RoleID.String(), tenantID)
	if err != nil {
		return nil, fmt.Errorf("error verificando rol: %w", err)
	}
	if !roleExists {
		return nil, fmt.Errorf("el rol con ID %s no existe", req.RoleID)
	}

	// Verificar que la relación no existe ya
	exists, err := s.userRoleRepo.Exists(ctx, req.UserID, req.RoleID, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("error verificando relación existente: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("la relación usuario-rol ya existe")
	}

	// Crear la relación
	userRole := domain_user_role.NewUserRole(req.UserID, req.RoleID, req.TenantID)

	if err := s.userRoleRepo.Create(ctx, userRole); err != nil {
		return nil, fmt.Errorf("error creando relación usuario-rol: %w", err)
	}

	return &domain_user_role.UserRoleResponseDto{
		UserRole: userRole,
		Message:  "Relación usuario-rol creada exitosamente",
	}, nil
}

// DeleteUserRole elimina una relación usuario-rol específica
func (s *userRoleService) DeleteUserRole(ctx context.Context, userID, roleID uuid.UUID, tenantID *uuid.UUID) error {
	// Verificar que la relación existe
	exists, err := s.userRoleRepo.Exists(ctx, userID, roleID, tenantID)
	if err != nil {
		return fmt.Errorf("error verificando relación: %w", err)
	}
	if !exists {
		return fmt.Errorf("la relación usuario-rol no existe")
	}

	if err := s.userRoleRepo.Delete(ctx, userID, roleID, tenantID); err != nil {
		return fmt.Errorf("error eliminando relación usuario-rol: %w", err)
	}

	return nil
}

// GetUserRoles obtiene todos los roles de un usuario
func (s *userRoleService) GetUserRoles(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) ([]domain_user_role.UserRole, error) {
	// Validar que el usuario existe
	userExists, err := s.userService.CheckUserExists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error verificando usuario: %w", err)
	}
	if !userExists {
		return nil, fmt.Errorf("el usuario con ID %s no existe", userID)
	}

	userRoles, err := s.userRoleRepo.GetByUserID(ctx, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo roles del usuario: %w", err)
	}

	return userRoles, nil
}

// GetRoleUsers obtiene todos los usuarios de un rol
func (s *userRoleService) GetRoleUsers(ctx context.Context, roleID uuid.UUID, tenantID *uuid.UUID) ([]domain_user_role.UserRole, error) {
	// Validar que el rol existe
	tenantIDStr := ""
	if tenantID != nil {
		tenantIDStr = tenantID.String()
	}

	roleExists, err := s.roleService.CheckRoleExistsByID(ctx, roleID.String(), tenantIDStr)
	if err != nil {
		return nil, fmt.Errorf("error verificando rol: %w", err)
	}
	if !roleExists {
		return nil, fmt.Errorf("el rol con ID %s no existe", roleID)
	}

	userRoles, err := s.userRoleRepo.GetByRoleID(ctx, roleID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo usuarios del rol: %w", err)
	}

	return userRoles, nil
}

// SyncRoleUsers sincroniza usuarios a un rol con validaciones
func (s *userRoleService) SyncRoleUsers(ctx context.Context, req *domain_user_role.SyncRoleUsersDto) error {
	// Validar que el rol existe
	tenantIDStr := ""
	if req.TenantID != nil {
		tenantIDStr = req.TenantID.String()
	}

	roleExists, err := s.roleService.CheckRoleExistsByID(ctx, req.RoleID.String(), tenantIDStr)
	if err != nil {
		return fmt.Errorf("error verificando rol: %w", err)
	}
	if !roleExists {
		return fmt.Errorf("el rol con ID %s no existe", req.RoleID)
	}

	// Validar que todos los usuarios existen
	for _, userID := range req.UserIDs {
		userExists, err := s.userService.CheckUserExists(ctx, userID)
		if err != nil {
			return fmt.Errorf("error verificando usuario %s: %w", userID, err)
		}
		if !userExists {
			return fmt.Errorf("el usuario con ID %s no existe", userID)
		}
	}

	// Realizar sincronización
	if err := s.userRoleRepo.SyncRoleUsers(ctx, req.RoleID, req.UserIDs, req.TenantID); err != nil {
		return fmt.Errorf("error sincronizando usuarios al rol: %w", err)
	}

	return nil
}

// SyncUserRoles sincroniza roles a un usuario con validaciones
func (s *userRoleService) SyncUserRoles(ctx context.Context, req *domain_user_role.SyncUserRolesDto) error {
	// Validar que el usuario existe
	userExists, err := s.userService.CheckUserExists(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("error verificando usuario: %w", err)
	}
	if !userExists {
		return fmt.Errorf("el usuario con ID %s no existe", req.UserID)
	}

	// Validar que todos los roles existen
	tenantIDStr := ""
	if req.TenantID != nil {
		tenantIDStr = req.TenantID.String()
	}

	for _, roleID := range req.RoleIDs {
		roleExists, err := s.roleService.CheckRoleExistsByID(ctx, roleID.String(), tenantIDStr)
		if err != nil {
			return fmt.Errorf("error verificando rol %s: %w", roleID, err)
		}
		if !roleExists {
			return fmt.Errorf("el rol con ID %s no existe", roleID)
		}
	}

	// Realizar sincronización
	if err := s.userRoleRepo.SyncUserRoles(ctx, req.UserID, req.RoleIDs, req.TenantID); err != nil {
		return fmt.Errorf("error sincronizando roles al usuario: %w", err)
	}

	return nil
}

// UserHasRole verifica si un usuario tiene un rol específico
func (s *userRoleService) UserHasRole(ctx context.Context, userID, roleID uuid.UUID, tenantID *uuid.UUID) (bool, error) {
	exists, err := s.userRoleRepo.Exists(ctx, userID, roleID, tenantID)
	if err != nil {
		return false, fmt.Errorf("error verificando relación usuario-rol: %w", err)
	}

	return exists, nil
}

// RemoveUserFromAllRoles elimina un usuario de todos sus roles
func (s *userRoleService) RemoveUserFromAllRoles(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) error {
	// Validar que el usuario existe
	userExists, err := s.userService.CheckUserExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("error verificando usuario: %w", err)
	}
	if !userExists {
		return fmt.Errorf("el usuario con ID %s no existe", userID)
	}

	if err := s.userRoleRepo.DeleteByUserID(ctx, userID, tenantID); err != nil {
		return fmt.Errorf("error eliminando roles del usuario: %w", err)
	}

	return nil
}

// RemoveAllUsersFromRole elimina todos los usuarios de un rol
func (s *userRoleService) RemoveAllUsersFromRole(ctx context.Context, roleID uuid.UUID, tenantID *uuid.UUID) error {
	// Validar que el rol existe
	tenantIDStr := ""
	if tenantID != nil {
		tenantIDStr = tenantID.String()
	}

	roleExists, err := s.roleService.CheckRoleExistsByID(ctx, roleID.String(), tenantIDStr)
	if err != nil {
		return fmt.Errorf("error verificando rol: %w", err)
	}
	if !roleExists {
		return fmt.Errorf("el rol con ID %s no existe", roleID)
	}

	if err := s.userRoleRepo.DeleteByRoleID(ctx, roleID, tenantID); err != nil {
		return fmt.Errorf("error eliminando usuarios del rol: %w", err)
	}

	return nil
}
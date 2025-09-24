package services

import (
	"context"
	"fmt"

	"github.com/JoseLuis21/mv-backend/internal/core/role/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/role/ports"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
	"github.com/google/uuid"
)

// roleService implementa el servicio de roles siguiendo principios de Clean Architecture
type roleService struct {
	roleRepo ports.RoleRepository
}

// NewRoleService crea una nueva instancia del servicio de roles
func NewRoleService(roleRepo ports.RoleRepository) ports.RoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}

// CreateRole crea un nuevo rol
func (s *roleService) CreateRole(ctx context.Context, req *domain.CreateRoleRequest) (*domain.RoleResponse, error) {
	// 1. Validar la solicitud de creación
	if err := s.ValidateRoleCreation(ctx, req); err != nil {
		return nil, err
	}

	// 2. Verificar si ya existe un rol con el mismo nombre
	exists, err := s.roleRepo.ExistsByName(ctx, req.Name, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("error verificando existencia del rol: %w", err)
	}
	if exists {
		if req.TenantID == nil {
			return nil, sharedErrors.NewBusinessError("ROLE_ALREADY_EXISTS",
				"Ya existe un rol global con este nombre", req.Name)
		}
		return nil, sharedErrors.NewBusinessError("ROLE_ALREADY_EXISTS",
			"Ya existe un rol con este nombre en el tenant", req.Name)
	}

	// 3. Crear entidad de rol
	role := domain.NewRole(req.Name, req.Description, req.TenantID)

	// 4. Guardar en la base de datos
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, fmt.Errorf("error creando rol: %w", err)
	}

	// 5. Retornar respuesta
	return role.ToResponse(), nil
}

// GetRoleByID obtiene un rol por su ID
func (s *roleService) GetRoleByID(ctx context.Context, id uuid.UUID) (*domain.RoleResponse, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo rol por ID: %w", err)
	}
	if role == nil {
		return nil, sharedErrors.NewBusinessError("ROLE_NOT_FOUND",
			"Rol no encontrado", id.String())
	}

	return role.ToResponse(), nil
}

// GetRoleByName obtiene un rol por su nombre
func (s *roleService) GetRoleByName(ctx context.Context, name string, tenantID *uuid.UUID) (*domain.RoleResponse, error) {
	role, err := s.roleRepo.GetByName(ctx, name, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo rol por nombre: %w", err)
	}
	if role == nil {
		return nil, sharedErrors.NewBusinessError("ROLE_NOT_FOUND",
			"Rol no encontrado", name)
	}

	return role.ToResponse(), nil
}

// GetGlobalRoles obtiene todos los roles globales del sistema
func (s *roleService) GetGlobalRoles(ctx context.Context) ([]*domain.RoleResponse, error) {
	roles, err := s.roleRepo.GetGlobalRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo roles globales: %w", err)
	}

	responses := make([]*domain.RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = role.ToResponse()
	}

	return responses, nil
}

// GetTenantRoles obtiene todos los roles de un tenant específico
func (s *roleService) GetTenantRoles(ctx context.Context, tenantID uuid.UUID) ([]*domain.RoleResponse, error) {
	roles, err := s.roleRepo.GetTenantRoles(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo roles del tenant: %w", err)
	}

	responses := make([]*domain.RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = role.ToResponse()
	}

	return responses, nil
}

// GetRoles obtiene roles con filtros de búsqueda y paginación
func (s *roleService) GetRoles(ctx context.Context, filter *domain.RoleFilterRequest) (*domain.RoleListResponse, error) {
	// Validar filtros de paginación
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	roles, total, err := s.roleRepo.GetAllRoles(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo roles: %w", err)
	}

	responses := make([]domain.RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = *role.ToResponse()
	}

	return &domain.RoleListResponse{
		Roles: responses,
		Total: total,
		Page:  filter.Page,
		Limit: filter.Limit,
	}, nil
}

// UpdateRole actualiza un rol existente
func (s *roleService) UpdateRole(ctx context.Context, id uuid.UUID, req *domain.UpdateRoleRequest) (*domain.RoleResponse, error) {
	// 1. Validar la solicitud de actualización
	if err := s.ValidateRoleUpdate(ctx, id, req); err != nil {
		return nil, err
	}

	// 2. Obtener el rol existente
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo rol para actualizar: %w", err)
	}
	if role == nil {
		return nil, sharedErrors.NewBusinessError("ROLE_NOT_FOUND",
			"Rol no encontrado", id.String())
	}

	// 3. Verificar que no es un rol del sistema (no se puede modificar)
	if role.IsSystemRole() {
		return nil, sharedErrors.NewBusinessError("SYSTEM_ROLE_READONLY",
			"Los roles del sistema no pueden ser modificados", role.Name)
	}

	// 4. Verificar si el nuevo nombre ya existe (si es diferente al actual)
	if req.Name != role.Name {
		exists, err := s.roleRepo.ExistsByName(ctx, req.Name, role.TenantID)
		if err != nil {
			return nil, fmt.Errorf("error verificando existencia del rol: %w", err)
		}
		if exists {
			return nil, sharedErrors.NewBusinessError("ROLE_ALREADY_EXISTS",
				"Ya existe un rol con este nombre", req.Name)
		}
	}

	// 5. Actualizar el rol usando método de dominio
	role.Update(req.Name, req.Description)

	// 6. Guardar cambios
	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, fmt.Errorf("error actualizando rol: %w", err)
	}

	return role.ToResponse(), nil
}

// DeleteRole elimina lógicamente un rol
func (s *roleService) DeleteRole(ctx context.Context, id uuid.UUID) error {
	// 1. Obtener el rol
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error obteniendo rol para eliminar: %w", err)
	}
	if role == nil {
		return sharedErrors.NewBusinessError("ROLE_NOT_FOUND",
			"Rol no encontrado", id.String())
	}

	// 2. Verificar que no es un rol del sistema (no se puede eliminar)
	if role.IsSystemRole() {
		return sharedErrors.NewBusinessError("SYSTEM_ROLE_READONLY",
			"Los roles del sistema no pueden ser eliminados", role.Name)
	}

	// TODO: Verificar que no hay usuarios asignados a este rol
	// userRoles, err := s.roleRepo.GetUserRoles(ctx, uuid.Nil, role.TenantID)
	// if len(userRoles) > 0 { return error }

	// 3. Eliminar el rol (soft delete)
	if err := s.roleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error eliminando rol: %w", err)
	}

	return nil
}

// ValidateRoleCreation valida que se puede crear un rol con los datos dados
func (s *roleService) ValidateRoleCreation(ctx context.Context, req *domain.CreateRoleRequest) error {
	// 1. Validaciones básicas de campos
	if req.Name == "" {
		return sharedErrors.NewValidationError("El nombre del rol es requerido", "name")
	}
	if len(req.Name) < 3 || len(req.Name) > 50 {
		return sharedErrors.NewValidationError("El nombre debe tener entre 3 y 50 caracteres", "name")
	}
	if len(req.Description) > 500 {
		return sharedErrors.NewValidationError("La descripción no puede exceder 500 caracteres", "description")
	}

	// 2. Validaciones específicas para roles del sistema
	systemRoles := []string{
		domain.RoleNameAdministrator,
		domain.RoleNameApprover,
		domain.RoleNameExpenseSubmitter,
	}

	for _, sysRole := range systemRoles {
		if req.Name == sysRole {
			return sharedErrors.NewBusinessError("SYSTEM_ROLE_READONLY",
				"No se pueden crear roles con nombres reservados del sistema", req.Name)
		}
	}

	// 3. Validaciones de multi-tenancy
	// Los roles globales solo deben ser creados por administradores del sistema
	if req.TenantID == nil {
		// TODO: Verificar permisos de administrador global
		// Por ahora permitimos la creación para pruebas
	}

	return nil
}

// ValidateRoleUpdate valida que se puede actualizar un rol con los datos dados
func (s *roleService) ValidateRoleUpdate(ctx context.Context, id uuid.UUID, req *domain.UpdateRoleRequest) error {
	// 1. Validaciones básicas de campos
	if req.Name == "" {
		return sharedErrors.NewValidationError("El nombre del rol es requerido", "name")
	}
	if len(req.Name) < 3 || len(req.Name) > 50 {
		return sharedErrors.NewValidationError("El nombre debe tener entre 3 y 50 caracteres", "name")
	}
	if len(req.Description) > 500 {
		return sharedErrors.NewValidationError("La descripción no puede exceder 500 caracteres", "description")
	}

	// 2. Validar que el ID no sea nulo
	if id == uuid.Nil {
		return sharedErrors.NewValidationError("ID del rol requerido", "id")
	}

	return nil
}

// InitializeSystemRoles crea los roles predefinidos del sistema si no existen
func (s *roleService) InitializeSystemRoles(ctx context.Context) error {
	systemRoles := []struct {
		name        string
		description string
	}{
		{
			name:        domain.RoleNameAdministrator,
			description: "Administrador del sistema con acceso completo",
		},
		{
			name:        domain.RoleNameApprover,
			description: "Aprobador de gastos y viáticos",
		},
		{
			name:        domain.RoleNameExpenseSubmitter,
			description: "Usuario que puede enviar gastos y solicitar viáticos",
		},
	}

	for _, sysRole := range systemRoles {
		// Verificar si el rol ya existe
		exists, err := s.roleRepo.ExistsByName(ctx, sysRole.name, nil)
		if err != nil {
			return fmt.Errorf("error verificando rol del sistema %s: %w", sysRole.name, err)
		}

		if !exists {
			// Crear el rol del sistema
			role := domain.NewRole(sysRole.name, sysRole.description, nil)
			if err := s.roleRepo.Create(ctx, role); err != nil {
				return fmt.Errorf("error creando rol del sistema %s: %w", sysRole.name, err)
			}
		}
	}

	return nil
}

// InitializeTenantRoles crea los roles por defecto para un nuevo tenant
func (s *roleService) InitializeTenantRoles(ctx context.Context, tenantID uuid.UUID) error {
	tenantRoles := []struct {
		name        string
		description string
	}{
		{
			name:        "Administrador Empresa",
			description: "Administrador de la empresa con acceso completo al tenant",
		},
		{
			name:        "Gerente",
			description: "Gerente con permisos de aprobación y supervisión",
		},
		{
			name:        "Empleado",
			description: "Empleado básico con permisos para crear gastos",
		},
	}

	for _, tenantRole := range tenantRoles {
		// Verificar si el rol ya existe para este tenant
		exists, err := s.roleRepo.ExistsByName(ctx, tenantRole.name, &tenantID)
		if err != nil {
			return fmt.Errorf("error verificando rol del tenant %s: %w", tenantRole.name, err)
		}

		if !exists {
			// Crear el rol del tenant
			role := domain.NewRole(tenantRole.name, tenantRole.description, &tenantID)
			if err := s.roleRepo.Create(ctx, role); err != nil {
				return fmt.Errorf("error creando rol del tenant %s: %w", tenantRole.name, err)
			}
		}
	}

	return nil
}

// GetUserRoles obtiene todos los roles asignados a un usuario
func (s *roleService) GetUserRoles(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) ([]*domain.RoleResponse, error) {
	roles, err := s.roleRepo.GetUserRoles(ctx, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo roles del usuario: %w", err)
	}

	responses := make([]*domain.RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = role.ToResponse()
	}

	return responses, nil
}

// CheckRoleExistsByID verifica si un rol existe por ID
func (s *roleService) CheckRoleExistsByID(ctx context.Context, id, tenantID string) (bool, error) {
	// Validar UUID
	roleID, err := uuid.Parse(id)
	if err != nil {
		return false, fmt.Errorf("ID de rol inválido: %w", err)
	}

	// Obtener el rol
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		// Si no se encuentra, no existe
		return false, nil
	}

	// Verificar tenant si se especifica
	if tenantID != "" {
		tid, err := uuid.Parse(tenantID)
		if err != nil {
			return false, fmt.Errorf("ID de tenant inválido: %w", err)
		}

		if role.TenantID == nil || *role.TenantID != tid {
			return false, nil
		}
	}

	return true, nil
}

// CheckRoleExistsByName verifica si un rol existe por nombre
func (s *roleService) CheckRoleExistsByName(ctx context.Context, name, tenantID string) (bool, error) {
	var tenantIDPtr *uuid.UUID

	if tenantID != "" {
		tid, err := uuid.Parse(tenantID)
		if err != nil {
			return false, fmt.Errorf("ID de tenant inválido: %w", err)
		}
		tenantIDPtr = &tid
	}

	return s.roleRepo.ExistsByName(ctx, name, tenantIDPtr)
}

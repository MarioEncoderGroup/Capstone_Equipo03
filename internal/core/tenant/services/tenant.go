package services

import (
	"context"
	"fmt"

	tenantDomain "github.com/JoseLuis21/mv-backend/internal/core/tenant/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/tenant/ports"
	userDomain "github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	userPorts "github.com/JoseLuis21/mv-backend/internal/core/user/ports"
	"github.com/google/uuid"
)

// tenantService implementa el servicio de tenant
type tenantService struct {
	tenantRepo  ports.TenantRepository
	userService userPorts.UserService
}

// NewTenantService crea una nueva instancia del servicio de tenant
func NewTenantService(
	tenantRepo ports.TenantRepository,
	userService userPorts.UserService,
) ports.TenantService {
	return &tenantService{
		tenantRepo:  tenantRepo,
		userService: userService,
	}
}

// GetTenantsByUser obtiene todos los tenants asociados a un usuario
func (s *tenantService) GetTenantsByUser(ctx context.Context, userID uuid.UUID) ([]*tenantDomain.Tenant, error) {
	// Obtener IDs de tenants del usuario
	tenantUsers, err := s.userService.GetTenantsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo tenants del usuario: %w", err)
	}

	var tenants []*tenantDomain.Tenant
	for _, tenantUser := range tenantUsers {
		tenant, err := s.tenantRepo.GetByID(ctx, tenantUser.TenantID)
		if err != nil {
			continue // Skip tenants that can't be loaded
		}
		if tenant != nil {
			tenants = append(tenants, tenant)
		}
	}

	return tenants, nil
}

// SelectTenant selecciona un tenant para el usuario autenticado
func (s *tenantService) SelectTenant(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*tenantDomain.SelectTenantResponseDto, error) {
	// Verificar que el usuario tenga acceso al tenant
	hasAccess, err := s.userService.UserHasAccessToTenant(ctx, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error verificando acceso al tenant: %w", err)
	}
	if !hasAccess {
		return nil, fmt.Errorf("usuario no tiene acceso al tenant especificado")
	}

	// Obtener el tenant
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo tenant: %w", err)
	}
	if tenant == nil {
		return nil, fmt.Errorf("tenant no encontrado")
	}

	// TODO: Generar nuevo token JWT con contexto del tenant
	accessToken := "jwt_token_with_tenant_context" // Placeholder

	return &tenantDomain.SelectTenantResponseDto{
		Tenant:      tenant,
		AccessToken: accessToken,
	}, nil
}

// GetTenantByID obtiene un tenant por su ID
func (s *tenantService) GetTenantByID(ctx context.Context, tenantID uuid.UUID) (*tenantDomain.Tenant, error) {
	return s.tenantRepo.GetByID(ctx, tenantID)
}

// CreateTenant crea un nuevo tenant
func (s *tenantService) CreateTenant(ctx context.Context, tenant *tenantDomain.Tenant) error {
	return s.tenantRepo.Create(ctx, tenant)
}

// CreateTenantFromDTO crea un nuevo tenant desde un DTO y lo asocia al usuario creador
func (s *tenantService) CreateTenantFromDTO(ctx context.Context, dto *tenantDomain.CreateTenantDTO, userID uuid.UUID) (*tenantDomain.Tenant, error) {
	// 1. Verificar que el RUT no exista
	exists, err := s.tenantRepo.ExistsByRUT(ctx, dto.Rut)
	if err != nil {
		return nil, fmt.Errorf("error verificando RUT: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("ya existe un tenant con el RUT proporcionado")
	}

	countryID, err := uuid.Parse(dto.CountryID)
	if err != nil {
		return nil, fmt.Errorf("country_id inválido: %w", err)
	}

	// 3. Crear entidad de tenant
	tenant := &tenantDomain.Tenant{
		ID:           uuid.New(),
		Rut:          dto.Rut,
		BusinessName: dto.BusinessName,
		Email:        dto.Email,
		Phone:        dto.Phone,
		Address:      dto.Address,
		Website:      dto.Website,
		Logo:         dto.Logo,
		RegionID:     dto.RegionID,
		CommuneID:    dto.CommuneID,
		CountryID:    countryID,
		Status:       string(tenantDomain.TenantStatusActive),
		NodeNumber:   1, // Fijo en 1 según arquitectura
		TenantName:   fmt.Sprintf("misviaticos_tenant_%s", uuid.New().String()[:8]),
		CreatedBy:    userID,
		UpdatedBy:    userID,
	}

	// 4. Crear tenant en la base de datos control
	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, fmt.Errorf("error creando tenant: %w", err)
	}

	// 5. Asociar usuario al tenant (crear entrada en tenant_users)
	tenantUser := &userDomain.TenantUser{
		ID:       uuid.New(),
		TenantID: tenant.ID,
		UserID:   userID,
	}
	if err := s.userService.AddUserToTenant(ctx, tenantUser); err != nil {
		return nil, fmt.Errorf("error asociando usuario al tenant: %w", err)
	}

	// 6. Crear base de datos del tenant dinámicamente
	// TODO: Implementar cuando tengamos el sistema de migraciones para tenants
	// if err := s.tenantRepo.CreateTenantDatabase(ctx, tenant.TenantName); err != nil {
	//     return nil, fmt.Errorf("error creando base de datos del tenant: %w", err)
	// }

	return tenant, nil
}

// UpdateTenant actualiza un tenant existente
func (s *tenantService) UpdateTenant(ctx context.Context, tenant *tenantDomain.Tenant) error {
	return s.tenantRepo.Update(ctx, tenant)
}

// GetTenantProfile obtiene el perfil completo del tenant
func (s *tenantService) GetTenantProfile(ctx context.Context, tenantID uuid.UUID) (*tenantDomain.Tenant, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo perfil del tenant: %w", err)
	}
	if tenant == nil {
		return nil, fmt.Errorf("tenant no encontrado")
	}
	return tenant, nil
}
package services

import (
	"context"
	"fmt"

	tenantDomain "github.com/JoseLuis21/mv-backend/internal/core/tenant/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/tenant/ports"
	userPorts "github.com/JoseLuis21/mv-backend/internal/core/user/ports"
	"github.com/google/uuid"
)

// tenantService implementa el servicio de tenant
type tenantService struct {
	tenantRepo ports.TenantRepository
	userRepo   userPorts.UserRepository
}

// NewTenantService crea una nueva instancia del servicio de tenant
func NewTenantService(
	tenantRepo ports.TenantRepository,
	userRepo userPorts.UserRepository,
) ports.TenantService {
	return &tenantService{
		tenantRepo: tenantRepo,
		userRepo:   userRepo,
	}
}

// GetTenantsByUser obtiene todos los tenants asociados a un usuario
func (s *tenantService) GetTenantsByUser(ctx context.Context, userID uuid.UUID) ([]*tenantDomain.Tenant, error) {
	// Obtener IDs de tenants del usuario
	tenantUsers, err := s.userRepo.GetTenantsByUser(ctx, userID)
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
	hasAccess, err := s.userRepo.UserHasAccessToTenant(ctx, userID, tenantID)
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
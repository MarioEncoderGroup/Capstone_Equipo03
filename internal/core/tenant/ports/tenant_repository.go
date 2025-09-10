package ports

import (
	"context"

	"github.com/JoseLuis21/mv-backend/internal/core/tenant/domain"
	"github.com/google/uuid"
)

// TenantRepository define el contrato para la persistencia de tenants
// Esta interfaz será implementada por el adaptador de PostgreSQL
type TenantRepository interface {
	// Create crea un nuevo tenant en la base de datos control
	Create(ctx context.Context, tenant *domain.Tenant) error

	// GetByID obtiene un tenant por su ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)

	// GetByRUT obtiene un tenant por su RUT (único para empresas chilenas)
	GetByRUT(ctx context.Context, rut string) (*domain.Tenant, error)

	// GetByBusinessName obtiene un tenant por su nombre de negocio
	GetByBusinessName(ctx context.Context, businessName string) (*domain.Tenant, error)

	// Update actualiza un tenant existente
	Update(ctx context.Context, tenant *domain.Tenant) error

	// Delete elimina lógicamente un tenant (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// ExistsByRUT verifica si existe un tenant con el RUT dado
	ExistsByRUT(ctx context.Context, rut string) (bool, error)

	// GetActiveTenantsCount obtiene el número de tenants activos
	GetActiveTenantsCount(ctx context.Context) (int64, error)

	// GetTenantUsers obtiene todos los usuarios asociados a un tenant
	GetTenantUsers(ctx context.Context, tenantID uuid.UUID) ([]uuid.UUID, error)

	// CreateTenantDatabase crea la base de datos específica del tenant
	// Esta operación crea una nueva base de datos con las tablas necesarias
	CreateTenantDatabase(ctx context.Context, tenantName string) error

	// GetNextNodeNumber obtiene el siguiente número de nodo disponible
	GetNextNodeNumber(ctx context.Context) (int, error)
}

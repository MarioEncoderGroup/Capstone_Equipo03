package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/core/tenant/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/tenant/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PostgreSQLTenantRepository implementa TenantRepository usando PostgreSQL
// Conecta con la base de datos misviaticos_control para gestionar tenants
type PostgreSQLTenantRepository struct {
	client *postgresql.PostgresqlClient
}

// NewPostgreSQLTenantRepository crea una nueva instancia del repositorio
func NewPostgreSQLTenantRepository(client *postgresql.PostgresqlClient) ports.TenantRepository {
	return &PostgreSQLTenantRepository{
		client: client,
	}
}

func (r *PostgreSQLTenantRepository) GetNextNodeNumber(ctx context.Context) (int, error) {
	var nodeNumber int
	query := `SELECT COALESCE(MAX(node_number), 0) + 1 FROM tenants`
	err := r.client.QueryRow(ctx, query).Scan(&nodeNumber)
	if err != nil {
		return 0, fmt.Errorf("error obteniendo siguiente número de nodo: %w", err)
	}
	return nodeNumber, nil
}

// Create crea un nuevo tenant en la base de datos control
func (r *PostgreSQLTenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		INSERT INTO tenants (
			id, rut, business_name, email, phone, address, website, logo,
			region_id, commune_id, country_id, status, node_number, tenant_name,
			created_by, updated_by, created, updated
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)`

	err := r.client.Exec(ctx, query,
		tenant.ID,
		tenant.RUT,
		tenant.BusinessName,
		tenant.Email,
		tenant.Phone,
		tenant.Address,
		tenant.Website,
		tenant.Logo,
		tenant.RegionID,
		tenant.CommuneID,
		tenant.CountryID,
		string(tenant.Status),
		tenant.NodeNumber,
		tenant.TenantName,
		tenant.CreatedBy,
		tenant.UpdatedBy,
		tenant.Created,
		tenant.Updated,
	)

	if err != nil {
		return fmt.Errorf("error creando tenant: %w", err)
	}

	return nil
}

// GetByID obtiene un tenant por su ID
func (r *PostgreSQLTenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	query := `
		SELECT id, rut, business_name, email, phone, address, website, logo,
			   region_id, commune_id, country_id, status, node_number, tenant_name,
			   created_by, updated_by, created, updated, deleted_at
		FROM tenants 
		WHERE id = $1 AND deleted_at IS NULL`

	tenant, err := r.scanTenant(r.client.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("tenant no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo tenant: %w", err)
	}

	return tenant, nil
}

// GetByRUT obtiene un tenant por su RUT
func (r *PostgreSQLTenantRepository) GetByRUT(ctx context.Context, rut string) (*domain.Tenant, error) {
	query := `
		SELECT id, rut, business_name, email, phone, address, website, logo,
			   region_id, commune_id, country_id, status, node_number, tenant_name,
			   created_by, updated_by, created, updated, deleted_at
		FROM tenants 
		WHERE rut = $1 AND deleted_at IS NULL`

	tenant, err := r.scanTenant(r.client.QueryRow(ctx, query, rut))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("tenant no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo tenant: %w", err)
	}

	return tenant, nil
}

// GetByBusinessName obtiene un tenant por su nombre de negocio
func (r *PostgreSQLTenantRepository) GetByBusinessName(ctx context.Context, businessName string) (*domain.Tenant, error) {
	query := `
		SELECT id, rut, business_name, email, phone, address, website, logo,
			   region_id, commune_id, country_id, status, node_number, tenant_name,
			   created_by, updated_by, created, updated, deleted_at
		FROM tenants 
		WHERE business_name = $1 AND deleted_at IS NULL`

	tenant, err := r.scanTenant(r.client.QueryRow(ctx, query, businessName))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("tenant no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo tenant: %w", err)
	}

	return tenant, nil
}

// Update actualiza un tenant existente
func (r *PostgreSQLTenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		UPDATE tenants SET
			rut = $2, business_name = $3, email = $4, phone = $5, address = $6,
			website = $7, logo = $8, region_id = $9, commune_id = $10, country_id = $11,
			status = $12, node_number = $13, tenant_name = $14, updated_by = $15, updated = $16
		WHERE id = $1 AND deleted_at IS NULL`

	tenant.Updated = time.Now()

	err := r.client.Exec(ctx, query,
		tenant.ID,
		tenant.RUT,
		tenant.BusinessName,
		tenant.Email,
		tenant.Phone,
		tenant.Address,
		tenant.Website,
		tenant.Logo,
		tenant.RegionID,
		tenant.CommuneID,
		tenant.CountryID,
		string(tenant.Status),
		tenant.NodeNumber,
		tenant.TenantName,
		tenant.UpdatedBy,
		tenant.Updated,
	)

	if err != nil {
		return fmt.Errorf("error actualizando tenant: %w", err)
	}

	return nil
}

// Delete elimina lógicamente un tenant (soft delete)
func (r *PostgreSQLTenantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE tenants 
		SET deleted_at = $2, updated = $2 
		WHERE id = $1 AND deleted_at IS NULL`

	now := time.Now()
	err := r.client.Exec(ctx, query, id, now)
	if err != nil {
		return fmt.Errorf("error eliminando tenant: %w", err)
	}

	return nil
}

// ExistsByRUT verifica si existe un tenant con el RUT dado
func (r *PostgreSQLTenantRepository) ExistsByRUT(ctx context.Context, rut string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tenants WHERE rut = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.client.QueryRow(ctx, query, rut).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error verificando RUT: %w", err)
	}

	return exists, nil
}

// GetActiveTenantsCount obtiene el número de tenants activos
func (r *PostgreSQLTenantRepository) GetActiveTenantsCount(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM tenants WHERE status = 'active' AND deleted_at IS NULL`

	var count int64
	err := r.client.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error contando tenants activos: %w", err)
	}

	return count, nil
}

// GetTenantUsers obtiene todos los usuarios asociados a un tenant
func (r *PostgreSQLTenantRepository) GetTenantUsers(ctx context.Context, tenantID uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT user_id
		FROM tenant_users 
		WHERE tenant_id = $1 AND deleted_at IS NULL`

	rows, err := r.client.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo usuarios del tenant: %w", err)
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("error escaneando user_id: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

// CreateTenantDatabase crea la base de datos específica del tenant
// Esta operación requiere permisos administrativos en PostgreSQL
func (r *PostgreSQLTenantRepository) CreateTenantDatabase(ctx context.Context, tenantName string) error {
	// IMPORTANTE: Esta operación requiere una conexión con permisos de superusuario
	// En un entorno de producción, esto podría manejarse de manera diferente

	// 1. Crear la base de datos
	createDBQuery := fmt.Sprintf(`CREATE DATABASE "%s"`, tenantName)

	if err := r.client.Exec(ctx, createDBQuery); err != nil {
		return fmt.Errorf("error creando base de datos del tenant: %w", err)
	}

	// TODO: 2. Aplicar migraciones específicas del tenant
	// Esto requeriría una conexión separada a la nueva base de datos
	// y la ejecución de las migraciones de tenant

	// Por ahora, solo creamos la base de datos vacía
	// En una implementación completa, aquí se ejecutarían las migraciones
	// que se encuentran en db/migrations-tenants/ (si existen)

	return nil
}

// scanTenant escanea una fila de tenant desde la base de datos
func (r *PostgreSQLTenantRepository) scanTenant(row pgx.Row) (*domain.Tenant, error) {
	tenant := &domain.Tenant{}
	var deletedAt sql.NullTime
	var statusStr string

	err := row.Scan(
		&tenant.ID,
		&tenant.RUT,
		&tenant.BusinessName,
		&tenant.Email,
		&tenant.Phone,
		&tenant.Address,
		&tenant.Website,
		&tenant.Logo,
		&tenant.RegionID,
		&tenant.CommuneID,
		&tenant.CountryID,
		&statusStr,
		&tenant.NodeNumber,
		&tenant.TenantName,
		&tenant.CreatedBy,
		&tenant.UpdatedBy,
		&tenant.Created,
		&tenant.Updated,
		&deletedAt,
	)

	if err != nil {
		return nil, err
	}

	// Convertir string a TenantStatus
	tenant.Status = domain.TenantStatus(statusStr)

	if deletedAt.Valid {
		tenant.DeletedAt = &deletedAt.Time
	}

	return tenant, nil
}

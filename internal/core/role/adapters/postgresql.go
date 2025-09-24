package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/JoseLuis21/mv-backend/internal/core/role/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/role/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
)

// PostgreSQLRoleRepository implementa RoleRepository usando PostgreSQL
// Conecta con la base de datos misviaticos_control para roles globales
// y bases de datos tenant específicas para roles de tenant
type PostgreSQLRoleRepository struct {
	client *postgresql.PostgresqlClient
}

// NewPostgreSQLRoleRepository crea una nueva instancia del repositorio
func NewPostgreSQLRoleRepository(client *postgresql.PostgresqlClient) ports.RoleRepository {
	return &PostgreSQLRoleRepository{
		client: client,
	}
}

// Create crea un nuevo rol en la base de datos
func (r *PostgreSQLRoleRepository) Create(ctx context.Context, role *domain.Role) error {
	query := `
		INSERT INTO roles (
			id, tenant_id, name, description, created, updated
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)`

	err := r.client.Exec(ctx, query,
		role.ID,
		role.TenantID,
		role.Name,
		role.Description,
		role.Created,
		role.Updated,
	)

	if err != nil {
		return fmt.Errorf("error creando rol: %w", err)
	}

	return nil
}

// GetByID obtiene un rol por su ID
func (r *PostgreSQLRoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	query := `
		SELECT id, tenant_id, name, description, created, updated, deleted_at
		FROM roles
		WHERE id = $1 AND deleted_at IS NULL`

	role, err := r.scanRole(r.client.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sharedErrors.NewBusinessError("ROLE_NOT_FOUND",
				"Rol no encontrado", id.String())
		}
		return nil, fmt.Errorf("error obteniendo rol por ID: %w", err)
	}

	return role, nil
}

// GetByName obtiene un rol por su nombre
func (r *PostgreSQLRoleRepository) GetByName(ctx context.Context, name string, tenantID *uuid.UUID) (*domain.Role, error) {
	var query string
	var args []interface{}

	if tenantID == nil {
		// Buscar rol global
		query = `
			SELECT id, tenant_id, name, description, created, updated, deleted_at
			FROM roles
			WHERE name = $1 AND tenant_id IS NULL AND deleted_at IS NULL`
		args = []interface{}{name}
	} else {
		// Buscar rol de tenant específico
		query = `
			SELECT id, tenant_id, name, description, created, updated, deleted_at
			FROM roles
			WHERE name = $1 AND tenant_id = $2 AND deleted_at IS NULL`
		args = []interface{}{name, *tenantID}
	}

	role, err := r.scanRole(r.client.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sharedErrors.NewBusinessError("ROLE_NOT_FOUND",
				"Rol no encontrado", name)
		}
		return nil, fmt.Errorf("error obteniendo rol por nombre: %w", err)
	}

	return role, nil
}

// GetGlobalRoles obtiene todos los roles globales del sistema
func (r *PostgreSQLRoleRepository) GetGlobalRoles(ctx context.Context) ([]*domain.Role, error) {
	query := `
		SELECT id, tenant_id, name, description, created, updated, deleted_at
		FROM roles
		WHERE tenant_id IS NULL AND deleted_at IS NULL
		ORDER BY name ASC`

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo roles globales: %w", err)
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		role, err := r.scanRoleFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("error escaneando rol global: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// GetTenantRoles obtiene todos los roles de un tenant específico
func (r *PostgreSQLRoleRepository) GetTenantRoles(ctx context.Context, tenantID uuid.UUID) ([]*domain.Role, error) {
	query := `
		SELECT id, tenant_id, name, description, created, updated, deleted_at
		FROM roles
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY name ASC`

	rows, err := r.client.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo roles del tenant: %w", err)
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		role, err := r.scanRoleFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("error escaneando rol del tenant: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// GetAllRoles obtiene roles con filtros de búsqueda
func (r *PostgreSQLRoleRepository) GetAllRoles(ctx context.Context, filter *domain.RoleFilterRequest) ([]*domain.Role, int, error) {
	// Construir WHERE clause
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIndex := 1

	if filter.TenantID != nil {
		whereClause += fmt.Sprintf(" AND tenant_id = $%d", argIndex)
		args = append(args, *filter.TenantID)
		argIndex++
	}

	if filter.Name != "" {
		whereClause += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Name+"%")
		argIndex++
	}

	// Query para contar total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM roles %s", whereClause)
	var total int
	err := r.client.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error contando roles: %w", err)
	}

	// Query para obtener roles con paginación
	offset := (filter.Page - 1) * filter.Limit
	dataQuery := fmt.Sprintf(`
		SELECT id, tenant_id, name, description, created, updated, deleted_at
		FROM roles
		%s
		ORDER BY name ASC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, filter.Limit, offset)

	rows, err := r.client.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error obteniendo roles: %w", err)
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		role, err := r.scanRoleFromRows(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("error escaneando rol: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, total, nil
}

// Update actualiza un rol existente
func (r *PostgreSQLRoleRepository) Update(ctx context.Context, role *domain.Role) error {
	query := `
		UPDATE roles SET
			name = $2, description = $3, updated = $4
		WHERE id = $1 AND deleted_at IS NULL`

	role.Updated = time.Now()

	err := r.client.Exec(ctx, query,
		role.ID,
		role.Name,
		role.Description,
		role.Updated,
	)

	if err != nil {
		return fmt.Errorf("error actualizando rol: %w", err)
	}

	return nil
}

// Delete elimina lógicamente un rol (soft delete)
func (r *PostgreSQLRoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE roles
		SET deleted_at = $2, updated = $2
		WHERE id = $1 AND deleted_at IS NULL`

	now := time.Now()
	err := r.client.Exec(ctx, query, id, now)
	if err != nil {
		return fmt.Errorf("error eliminando rol: %w", err)
	}

	return nil
}

// ExistsByName verifica si existe un rol con el nombre dado
func (r *PostgreSQLRoleRepository) ExistsByName(ctx context.Context, name string, tenantID *uuid.UUID) (bool, error) {
	var query string
	var args []interface{}

	if tenantID == nil {
		// Verificar rol global
		query = `SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1 AND tenant_id IS NULL AND deleted_at IS NULL)`
		args = []interface{}{name}
	} else {
		// Verificar rol de tenant específico
		query = `SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1 AND tenant_id = $2 AND deleted_at IS NULL)`
		args = []interface{}{name, *tenantID}
	}

	var exists bool
	err := r.client.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error verificando existencia del rol: %w", err)
	}

	return exists, nil
}

// IsSystemRole verifica si un rol es uno de los roles predefinidos del sistema
func (r *PostgreSQLRoleRepository) IsSystemRole(ctx context.Context, roleID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM roles
			WHERE id = $1 AND tenant_id IS NULL AND name IN ($2, $3, $4) AND deleted_at IS NULL
		)`

	var exists bool
	err := r.client.QueryRow(ctx, query, roleID,
		domain.RoleNameAdministrator,
		domain.RoleNameApprover,
		domain.RoleNameExpenseSubmitter).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error verificando si es rol del sistema: %w", err)
	}

	return exists, nil
}

// GetUserRoles obtiene todos los roles asignados a un usuario
func (r *PostgreSQLRoleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) ([]*domain.Role, error) {
	var query string
	var args []interface{}

	if tenantID == nil {
		// Obtener todos los roles del usuario (globales y de todos los tenants)
		query = `
			SELECT r.id, r.tenant_id, r.name, r.description, r.created, r.updated, r.deleted_at
			FROM roles r
			INNER JOIN user_roles ur ON r.id = ur.role_id
			WHERE ur.user_id = $1 AND r.deleted_at IS NULL AND ur.deleted_at IS NULL
			ORDER BY r.name ASC`
		args = []interface{}{userID}
	} else {
		// Obtener roles del usuario para un tenant específico (incluye roles globales)
		query = `
			SELECT r.id, r.tenant_id, r.name, r.description, r.created, r.updated, r.deleted_at
			FROM roles r
			INNER JOIN user_roles ur ON r.id = ur.role_id
			WHERE ur.user_id = $1 AND (r.tenant_id = $2 OR r.tenant_id IS NULL)
			AND r.deleted_at IS NULL AND ur.deleted_at IS NULL
			ORDER BY r.name ASC`
		args = []interface{}{userID, *tenantID}
	}

	rows, err := r.client.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo roles del usuario: %w", err)
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		role, err := r.scanRoleFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("error escaneando rol del usuario: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// scanRole escanea una fila de rol desde la base de datos (pgx.Row)
func (r *PostgreSQLRoleRepository) scanRole(row pgx.Row) (*domain.Role, error) {
	role := &domain.Role{}
	var deletedAt sql.NullTime

	err := row.Scan(
		&role.ID,
		&role.TenantID,
		&role.Name,
		&role.Description,
		&role.Created,
		&role.Updated,
		&deletedAt,
	)

	if err != nil {
		return nil, err
	}

	if deletedAt.Valid {
		role.DeletedAt = &deletedAt.Time
	}

	return role, nil
}

// scanRoleFromRows escanea una fila de rol desde múltiples filas (pgx.Rows)
func (r *PostgreSQLRoleRepository) scanRoleFromRows(rows pgx.Rows) (*domain.Role, error) {
	role := &domain.Role{}
	var deletedAt sql.NullTime

	err := rows.Scan(
		&role.ID,
		&role.TenantID,
		&role.Name,
		&role.Description,
		&role.Created,
		&role.Updated,
		&deletedAt,
	)

	if err != nil {
		return nil, err
	}

	if deletedAt.Valid {
		role.DeletedAt = &deletedAt.Time
	}

	return role, nil
}
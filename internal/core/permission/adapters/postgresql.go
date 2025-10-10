package adapters

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/permission/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/permission/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
)

// pgPermissionRepository implementa PermissionRepository para PostgreSQL
type pgPermissionRepository struct {
	db *postgresql.PostgresqlClient
}

// NewPgPermissionRepository crea una nueva instancia del repositorio PostgreSQL
func NewPgPermissionRepository(db *postgresql.PostgresqlClient) ports.PermissionRepository {
	return &pgPermissionRepository{db: db}
}

// Create crea un nuevo permiso en la base de datos
func (r *pgPermissionRepository) Create(ctx context.Context, permission *domain.Permission) error {
	query := `
		INSERT INTO permissions (id, name, description, section, created, updated)
		VALUES ($1, $2, $3, $4, NOW(), NOW())`

	err := r.db.Exec(ctx, query,
		permission.ID,
		permission.Name,
		permission.Description,
		permission.Section,
	)

	if err != nil {
		return fmt.Errorf("error creando permiso: %w", err)
	}

	return nil
}

// GetByID obtiene un permiso por su ID
func (r *pgPermissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Permission, error) {
	query := `
		SELECT id, name, description, section, created, updated
		FROM permissions
		WHERE id = $1 AND deleted_at IS NULL`

	var permission domain.Permission
	var description *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&permission.ID,
		&permission.Name,
		&description,
		&permission.Section,
		&permission.Created,
		&permission.Updated,
	)

	if err != nil {
		return nil, fmt.Errorf("error obteniendo permiso por ID: %w", err)
	}

	permission.Description = description
	return &permission, nil
}

// GetByName obtiene un permiso por su nombre
func (r *pgPermissionRepository) GetByName(ctx context.Context, name string) (*domain.Permission, error) {
	query := `
		SELECT id, name, description, section, created, updated
		FROM permissions
		WHERE name = $1 AND deleted_at IS NULL`

	var permission domain.Permission
	var description *string

	err := r.db.QueryRow(ctx, query, name).Scan(
		&permission.ID,
		&permission.Name,
		&description,
		&permission.Section,
		&permission.Created,
		&permission.Updated,
	)

	if err != nil {
		return nil, fmt.Errorf("error obteniendo permiso por nombre: %w", err)
	}

	permission.Description = description
	return &permission, nil
}

// GetBySection obtiene todos los permisos de una sección específica
func (r *pgPermissionRepository) GetBySection(ctx context.Context, section string) ([]*domain.Permission, error) {
	query := `
		SELECT id, name, description, section, created, updated
		FROM permissions
		WHERE section = $1 AND deleted_at IS NULL
		ORDER BY name ASC`

	rows, err := r.db.Query(ctx, query, section)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo permisos por sección: %w", err)
	}
	defer rows.Close()

	var permissions []*domain.Permission

	for rows.Next() {
		var permission domain.Permission
		var description *string

		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&description,
			&permission.Section,
			&permission.Created,
			&permission.Updated,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando permiso: %w", err)
		}

		permission.Description = description
		permissions = append(permissions, &permission)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando permisos: %w", err)
	}

	return permissions, nil
}

// GetAllPermissions obtiene permisos con filtros de búsqueda y paginación
func (r *pgPermissionRepository) GetAllPermissions(ctx context.Context, filter *domain.PermissionFilterRequest) ([]*domain.Permission, int, error) {
	// Construir la consulta base
	baseQuery := `FROM permissions WHERE deleted_at IS NULL`
	var args []interface{}
	argCount := 0

	// Aplicar filtros
	if filter.Search != "" {
		argCount++
		baseQuery += fmt.Sprintf(" AND (name ILIKE $%d OR section ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filter.Search+"%")
	}

	if filter.Section != "" {
		argCount++
		baseQuery += fmt.Sprintf(" AND section = $%d", argCount)
		args = append(args, filter.Section)
	}

	// Contar total de registros
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error contando permisos: %w", err)
	}

	// Obtener registros con paginación
	selectQuery := `
		SELECT id, name, description, section, created, updated ` +
		baseQuery +
		` ORDER BY ` + r.buildSortClause(filter.SortBy, filter.SortDir) +
		fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount+1, argCount+2)

	args = append(args, filter.Limit, (filter.Page-1)*filter.Limit)

	rows, err := r.db.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error obteniendo permisos: %w", err)
	}
	defer rows.Close()

	var permissions []*domain.Permission

	for rows.Next() {
		var permission domain.Permission
		var description *string

		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&description,
			&permission.Section,
			&permission.Created,
			&permission.Updated,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error escaneando permiso: %w", err)
		}

		permission.Description = description
		permissions = append(permissions, &permission)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterando permisos: %w", err)
	}

	return permissions, total, nil
}

// GetGroupedBySection obtiene permisos agrupados por sección
func (r *pgPermissionRepository) GetGroupedBySection(ctx context.Context) (map[string][]*domain.Permission, error) {
	query := `
		SELECT id, name, description, section, created, updated
		FROM permissions
		WHERE deleted_at IS NULL
		ORDER BY section ASC, name ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo permisos agrupados: %w", err)
	}
	defer rows.Close()

	grouped := make(map[string][]*domain.Permission)

	for rows.Next() {
		var permission domain.Permission
		var description *string

		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&description,
			&permission.Section,
			&permission.Created,
			&permission.Updated,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando permiso: %w", err)
		}

		permission.Description = description
		grouped[permission.Section] = append(grouped[permission.Section], &permission)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando permisos agrupados: %w", err)
	}

	return grouped, nil
}

// Update actualiza un permiso existente
func (r *pgPermissionRepository) Update(ctx context.Context, permission *domain.Permission) error {
	query := `
		UPDATE permissions
		SET name = $2, description = $3, section = $4, updated = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.Exec(ctx, query,
		permission.ID,
		permission.Name,
		permission.Description,
		permission.Section,
	)

	if err != nil {
		return fmt.Errorf("error actualizando permiso: %w", err)
	}

	return nil
}

// Delete elimina lógicamente un permiso
func (r *pgPermissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE permissions
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error eliminando permiso: %w", err)
	}

	return nil
}

// ExistsByName verifica si existe un permiso con el nombre dado
func (r *pgPermissionRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM permissions
			WHERE name = $1 AND deleted_at IS NULL
		)`

	var exists bool
	err := r.db.QueryRow(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error verificando existencia de permiso: %w", err)
	}

	return exists, nil
}

// IsSystemPermission verifica si un permiso es del sistema (no puede ser modificado)
func (r *pgPermissionRepository) IsSystemPermission(ctx context.Context, id uuid.UUID) (bool, error) {
	// Los permisos del sistema tienen secciones específicas o nombres específicos
	query := `
		SELECT EXISTS(
			SELECT 1 FROM permissions
			WHERE id = $1
			AND (section IN ('system', 'core') OR name LIKE 'system_%')
			AND deleted_at IS NULL
		)`

	var isSystem bool
	err := r.db.QueryRow(ctx, query, id).Scan(&isSystem)
	if err != nil {
		return false, fmt.Errorf("error verificando si es permiso del sistema: %w", err)
	}

	return isSystem, nil
}

// GetRolePermissions obtiene todos los permisos asignados a un rol
func (r *pgPermissionRepository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*domain.Permission, error) {
	query := `
		SELECT p.id, p.name, p.description, p.section, p.created, p.updated
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1 AND p.deleted_at IS NULL AND rp.deleted_at IS NULL
		ORDER BY p.section ASC, p.name ASC`

	rows, err := r.db.Query(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo permisos del rol: %w", err)
	}
	defer rows.Close()

	var permissions []*domain.Permission

	for rows.Next() {
		var permission domain.Permission
		var description *string

		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&description,
			&permission.Section,
			&permission.Created,
			&permission.Updated,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando permiso del rol: %w", err)
		}

		permission.Description = description
		permissions = append(permissions, &permission)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando permisos del rol: %w", err)
	}

	return permissions, nil
}

// GetAvailableSections obtiene todas las secciones de permisos disponibles
func (r *pgPermissionRepository) GetAvailableSections(ctx context.Context) ([]string, error) {
	query := `
		SELECT DISTINCT section
		FROM permissions
		WHERE deleted_at IS NULL
		ORDER BY section ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo secciones disponibles: %w", err)
	}
	defer rows.Close()

	var sections []string

	for rows.Next() {
		var section string
		err := rows.Scan(&section)
		if err != nil {
			return nil, fmt.Errorf("error escaneando sección: %w", err)
		}

		sections = append(sections, section)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando secciones: %w", err)
	}

	return sections, nil
}

// buildSortClause construye la cláusula ORDER BY con validación de seguridad
func (r *pgPermissionRepository) buildSortClause(sortBy, sortDir string) string {
	// Campos permitidos para ordenamiento
	allowedSortFields := map[string]string{
		"name":    "name",
		"section": "section",
		"created": "created",
		"updated": "updated",
	}

	// Validar campo de ordenamiento
	sortField, exists := allowedSortFields[sortBy]
	if !exists {
		sortField = "created" // campo por defecto
	}

	// Validar dirección de ordenamiento
	sortDirection := "DESC"
	if strings.ToUpper(sortDir) == "ASC" {
		sortDirection = "ASC"
	}

	return fmt.Sprintf("%s %s", sortField, sortDirection)
}
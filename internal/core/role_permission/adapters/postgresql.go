package adapters

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domain_role_permission "github.com/JoseLuis21/mv-backend/internal/core/role_permission/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/role_permission/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
)

// pgRolePermissionRepository implementa RolePermissionRepository para PostgreSQL
type pgRolePermissionRepository struct {
	db *postgresql.PostgresqlClient
}

// NewPgRolePermissionRepository crea una nueva instancia del repositorio PostgreSQL
func NewPgRolePermissionRepository(db *postgresql.PostgresqlClient) ports.RolePermissionRepository {
	return &pgRolePermissionRepository{db: db}
}

// Create crea una nueva relación rol-permiso
func (r *pgRolePermissionRepository) Create(ctx context.Context, rolePermission *domain_role_permission.RolePermission) error {
	query := `
		INSERT INTO role_permissions (id, role_id, permission_id, created, updated)
		VALUES ($1, $2, $3, $4, $5)`

	err := r.db.Exec(ctx, query,
		rolePermission.ID,
		rolePermission.RoleID,
		rolePermission.PermissionID,
		rolePermission.Created,
		rolePermission.Updated,
	)

	if err != nil {
		return fmt.Errorf("error creando relación rol-permiso: %w", err)
	}

	return nil
}

// Delete elimina una relación rol-permiso específica
func (r *pgRolePermissionRepository) Delete(ctx context.Context, roleID, permissionID uuid.UUID) error {
	query := `
		DELETE FROM role_permissions
		WHERE role_id = $1 AND permission_id = $2`

	err := r.db.Exec(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("error eliminando relación rol-permiso: %w", err)
	}

	return nil
}

// GetByRoleID obtiene todos los permisos de un rol específico
func (r *pgRolePermissionRepository) GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]domain_role_permission.RolePermission, error) {
	query := `
		SELECT id, role_id, permission_id, created, updated
		FROM role_permissions
		WHERE role_id = $1
		ORDER BY created DESC`

	rows, err := r.db.Query(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("error consultando permisos del rol: %w", err)
	}
	defer rows.Close()

	var rolePermissions []domain_role_permission.RolePermission

	for rows.Next() {
		var rolePermission domain_role_permission.RolePermission

		err := rows.Scan(
			&rolePermission.ID,
			&rolePermission.RoleID,
			&rolePermission.PermissionID,
			&rolePermission.Created,
			&rolePermission.Updated,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando permiso del rol: %w", err)
		}

		rolePermissions = append(rolePermissions, rolePermission)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando permisos del rol: %w", err)
	}

	return rolePermissions, nil
}

// GetByPermissionID obtiene todos los roles de un permiso específico
func (r *pgRolePermissionRepository) GetByPermissionID(ctx context.Context, permissionID uuid.UUID) ([]domain_role_permission.RolePermission, error) {
	query := `
		SELECT id, role_id, permission_id, created, updated
		FROM role_permissions
		WHERE permission_id = $1
		ORDER BY created DESC`

	rows, err := r.db.Query(ctx, query, permissionID)
	if err != nil {
		return nil, fmt.Errorf("error consultando roles del permiso: %w", err)
	}
	defer rows.Close()

	var rolePermissions []domain_role_permission.RolePermission

	for rows.Next() {
		var rolePermission domain_role_permission.RolePermission

		err := rows.Scan(
			&rolePermission.ID,
			&rolePermission.RoleID,
			&rolePermission.PermissionID,
			&rolePermission.Created,
			&rolePermission.Updated,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando rol del permiso: %w", err)
		}

		rolePermissions = append(rolePermissions, rolePermission)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando roles del permiso: %w", err)
	}

	return rolePermissions, nil
}

// Exists verifica si existe una relación rol-permiso específica
func (r *pgRolePermissionRepository) Exists(ctx context.Context, roleID, permissionID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM role_permissions
			WHERE role_id = $1 AND permission_id = $2
		)`

	var exists bool
	err := r.db.QueryRow(ctx, query, roleID, permissionID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error verificando existencia de relación rol-permiso: %w", err)
	}

	return exists, nil
}

// SyncRolePermissions sincroniza múltiples permisos a un rol (reemplaza existentes)
func (r *pgRolePermissionRepository) SyncRolePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	// Primero eliminar permisos existentes del rol
	deleteQuery := `DELETE FROM role_permissions WHERE role_id = $1`
	err := r.db.Exec(ctx, deleteQuery, roleID)
	if err != nil {
		return fmt.Errorf("error eliminando permisos existentes del rol: %w", err)
	}

	// Insertar nuevos permisos al rol
	if len(permissionIDs) > 0 {
		insertQuery := `
			INSERT INTO role_permissions (id, role_id, permission_id, created, updated)
			VALUES ($1, $2, $3, NOW(), NOW())`

		for _, permissionID := range permissionIDs {
			err = r.db.Exec(ctx, insertQuery, uuid.New(), roleID, permissionID)
			if err != nil {
				return fmt.Errorf("error insertando permiso %s al rol: %w", permissionID, err)
			}
		}
	}

	return nil
}

// DeleteByRoleID elimina todas las relaciones de un rol
func (r *pgRolePermissionRepository) DeleteByRoleID(ctx context.Context, roleID uuid.UUID) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1`

	err := r.db.Exec(ctx, query, roleID)
	if err != nil {
		return fmt.Errorf("error eliminando permisos del rol: %w", err)
	}

	return nil
}

// DeleteByPermissionID elimina todas las relaciones de un permiso
func (r *pgRolePermissionRepository) DeleteByPermissionID(ctx context.Context, permissionID uuid.UUID) error {
	query := `DELETE FROM role_permissions WHERE permission_id = $1`

	err := r.db.Exec(ctx, query, permissionID)
	if err != nil {
		return fmt.Errorf("error eliminando roles del permiso: %w", err)
	}

	return nil
}

// GetPermissionsByRoleID obtiene información completa de permisos por rol
func (r *pgRolePermissionRepository) GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]domain_role_permission.Permission, error) {
	query := `
		SELECT p.id, p.name, p.description, p.section, p.created, p.updated
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.section, p.name`

	rows, err := r.db.Query(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("error consultando permisos completos del rol: %w", err)
	}
	defer rows.Close()

	var permissions []domain_role_permission.Permission

	for rows.Next() {
		var permission domain_role_permission.Permission
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
			return nil, fmt.Errorf("error escaneando permiso completo: %w", err)
		}

		permission.Description = description
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando permisos completos: %w", err)
	}

	return permissions, nil
}

// GetPermissionsByRoleIDs obtiene permisos por múltiples roles
func (r *pgRolePermissionRepository) GetPermissionsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]domain_role_permission.Permission, error) {
	if len(roleIDs) == 0 {
		return []domain_role_permission.Permission{}, nil
	}

	// Crear placeholders para la consulta SQL
	placeholders := ""
	args := make([]interface{}, len(roleIDs))
	for i, roleID := range roleIDs {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += fmt.Sprintf("$%d", i+1)
		args[i] = roleID
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT p.id, p.name, p.description, p.section, p.created, p.updated
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id IN (%s)
		ORDER BY p.section, p.name`, placeholders)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error consultando permisos de múltiples roles: %w", err)
	}
	defer rows.Close()

	var permissions []domain_role_permission.Permission

	for rows.Next() {
		var permission domain_role_permission.Permission
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
			return nil, fmt.Errorf("error escaneando permiso de múltiples roles: %w", err)
		}

		permission.Description = description
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando permisos de múltiples roles: %w", err)
	}

	return permissions, nil
}

// GetPermissionNamesByRoleIDs obtiene solo los nombres de permisos por múltiples roles
func (r *pgRolePermissionRepository) GetPermissionNamesByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]string, error) {
	if len(roleIDs) == 0 {
		return []string{}, nil
	}

	// Crear placeholders para la consulta SQL
	placeholders := ""
	args := make([]interface{}, len(roleIDs))
	for i, roleID := range roleIDs {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += fmt.Sprintf("$%d", i+1)
		args[i] = roleID
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT p.name
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id IN (%s)
		ORDER BY p.name`, placeholders)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error consultando nombres de permisos de múltiples roles: %w", err)
	}
	defer rows.Close()

	var permissionNames []string

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("error escaneando nombre de permiso: %w", err)
		}

		permissionNames = append(permissionNames, name)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando nombres de permisos: %w", err)
	}

	return permissionNames, nil
}
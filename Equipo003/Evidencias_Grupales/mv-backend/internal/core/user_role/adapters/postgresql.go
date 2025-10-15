package adapters

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domain_user_role "github.com/JoseLuis21/mv-backend/internal/core/user_role/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/user_role/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
)

// pgUserRoleRepository implementa UserRoleRepository para PostgreSQL
type pgUserRoleRepository struct {
	db *postgresql.PostgresqlClient
}

// NewPgUserRoleRepository crea una nueva instancia del repositorio PostgreSQL
func NewPgUserRoleRepository(db *postgresql.PostgresqlClient) ports.UserRoleRepository {
	return &pgUserRoleRepository{db: db}
}

// Create crea una nueva relación usuario-rol
func (r *pgUserRoleRepository) Create(ctx context.Context, userRole *domain_user_role.UserRole) error {
	query := `
		INSERT INTO user_roles (id, user_id, role_id, tenant_id, created, updated)
		VALUES ($1, $2, $3, $4, $5, $6)`

	err := r.db.Exec(ctx, query,
		userRole.ID,
		userRole.UserID,
		userRole.RoleID,
		userRole.TenantID,
		userRole.Created,
		userRole.Updated,
	)

	if err != nil {
		return fmt.Errorf("error creando relación usuario-rol: %w", err)
	}

	return nil
}

// Delete elimina una relación usuario-rol específica
func (r *pgUserRoleRepository) Delete(ctx context.Context, userID, roleID uuid.UUID, tenantID *uuid.UUID) error {
	query := `
		DELETE FROM user_roles
		WHERE user_id = $1 AND role_id = $2 AND ($3::uuid IS NULL OR tenant_id = $3)`

	err := r.db.Exec(ctx, query, userID, roleID, tenantID)
	if err != nil {
		return fmt.Errorf("error eliminando relación usuario-rol: %w", err)
	}

	return nil
}

// GetByUserID obtiene todos los roles de un usuario específico
func (r *pgUserRoleRepository) GetByUserID(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) ([]domain_user_role.UserRole, error) {
	query := `
		SELECT id, user_id, role_id, tenant_id, created, updated
		FROM user_roles
		WHERE user_id = $1 AND ($2::uuid IS NULL OR tenant_id = $2)
		ORDER BY created DESC`

	rows, err := r.db.Query(ctx, query, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error consultando roles del usuario: %w", err)
	}
	defer rows.Close()

	var userRoles []domain_user_role.UserRole

	for rows.Next() {
		var userRole domain_user_role.UserRole
		var tenantIDPtr *uuid.UUID

		err := rows.Scan(
			&userRole.ID,
			&userRole.UserID,
			&userRole.RoleID,
			&tenantIDPtr,
			&userRole.Created,
			&userRole.Updated,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando rol del usuario: %w", err)
		}

		userRole.TenantID = tenantIDPtr
		userRoles = append(userRoles, userRole)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando roles del usuario: %w", err)
	}

	return userRoles, nil
}

// GetByRoleID obtiene todos los usuarios de un rol específico
func (r *pgUserRoleRepository) GetByRoleID(ctx context.Context, roleID uuid.UUID, tenantID *uuid.UUID) ([]domain_user_role.UserRole, error) {
	query := `
		SELECT id, user_id, role_id, tenant_id, created, updated
		FROM user_roles
		WHERE role_id = $1 AND ($2::uuid IS NULL OR tenant_id = $2)
		ORDER BY created DESC`

	rows, err := r.db.Query(ctx, query, roleID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error consultando usuarios del rol: %w", err)
	}
	defer rows.Close()

	var userRoles []domain_user_role.UserRole

	for rows.Next() {
		var userRole domain_user_role.UserRole
		var tenantIDPtr *uuid.UUID

		err := rows.Scan(
			&userRole.ID,
			&userRole.UserID,
			&userRole.RoleID,
			&tenantIDPtr,
			&userRole.Created,
			&userRole.Updated,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando usuario del rol: %w", err)
		}

		userRole.TenantID = tenantIDPtr
		userRoles = append(userRoles, userRole)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando usuarios del rol: %w", err)
	}

	return userRoles, nil
}

// Exists verifica si existe una relación usuario-rol específica
func (r *pgUserRoleRepository) Exists(ctx context.Context, userID, roleID uuid.UUID, tenantID *uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_roles
			WHERE user_id = $1 AND role_id = $2 AND ($3::uuid IS NULL OR tenant_id = $3)
		)`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID, roleID, tenantID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error verificando existencia de relación usuario-rol: %w", err)
	}

	return exists, nil
}

// SyncRoleUsers sincroniza múltiples usuarios a un rol (reemplaza existentes)
func (r *pgUserRoleRepository) SyncRoleUsers(ctx context.Context, roleID uuid.UUID, userIDs []uuid.UUID, tenantID *uuid.UUID) error {
	// Eliminar usuarios existentes del rol
	deleteQuery := `
		DELETE FROM user_roles
		WHERE role_id = $1 AND ($2::uuid IS NULL OR tenant_id = $2)`

	err := r.db.Exec(ctx, deleteQuery, roleID, tenantID)
	if err != nil {
		return fmt.Errorf("error eliminando usuarios existentes del rol: %w", err)
	}

	// Insertar nuevos usuarios al rol
	if len(userIDs) > 0 {
		insertQuery := `
			INSERT INTO user_roles (id, user_id, role_id, tenant_id, created, updated)
			VALUES ($1, $2, $3, $4, NOW(), NOW())`

		for _, userID := range userIDs {
			err = r.db.Exec(ctx, insertQuery, uuid.New(), userID, roleID, tenantID)
			if err != nil {
				return fmt.Errorf("error insertando usuario %s al rol: %w", userID, err)
			}
		}
	}

	return nil
}

// SyncUserRoles sincroniza múltiples roles a un usuario (reemplaza existentes)
func (r *pgUserRoleRepository) SyncUserRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID, tenantID *uuid.UUID) error {
	// Eliminar roles existentes del usuario
	deleteQuery := `
		DELETE FROM user_roles
		WHERE user_id = $1 AND ($2::uuid IS NULL OR tenant_id = $2)`

	err := r.db.Exec(ctx, deleteQuery, userID, tenantID)
	if err != nil {
		return fmt.Errorf("error eliminando roles existentes del usuario: %w", err)
	}

	// Insertar nuevos roles al usuario
	if len(roleIDs) > 0 {
		insertQuery := `
			INSERT INTO user_roles (id, user_id, role_id, tenant_id, created, updated)
			VALUES ($1, $2, $3, $4, NOW(), NOW())`

		for _, roleID := range roleIDs {
			err = r.db.Exec(ctx, insertQuery, uuid.New(), userID, roleID, tenantID)
			if err != nil {
				return fmt.Errorf("error insertando rol %s al usuario: %w", roleID, err)
			}
		}
	}

	return nil
}

// DeleteByUserID elimina todas las relaciones de un usuario
func (r *pgUserRoleRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) error {
	query := `
		DELETE FROM user_roles
		WHERE user_id = $1 AND ($2::uuid IS NULL OR tenant_id = $2)`

	err := r.db.Exec(ctx, query, userID, tenantID)
	if err != nil {
		return fmt.Errorf("error eliminando roles del usuario: %w", err)
	}

	return nil
}

// DeleteByRoleID elimina todas las relaciones de un rol
func (r *pgUserRoleRepository) DeleteByRoleID(ctx context.Context, roleID uuid.UUID, tenantID *uuid.UUID) error {
	query := `
		DELETE FROM user_roles
		WHERE role_id = $1 AND ($2::uuid IS NULL OR tenant_id = $2)`

	err := r.db.Exec(ctx, query, roleID, tenantID)
	if err != nil {
		return fmt.Errorf("error eliminando usuarios del rol: %w", err)
	}

	return nil
}
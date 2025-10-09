package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/JoseLuis21/mv-backend/internal/core/expense/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/expense/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PostgreSQLCategoryRepository implementa CategoryRepository usando PostgreSQL
type PostgreSQLCategoryRepository struct {
	client *postgresql.PostgresqlClient
}

// NewPostgreSQLCategoryRepository crea una nueva instancia del repositorio
func NewPostgreSQLCategoryRepository(client *postgresql.PostgresqlClient) ports.CategoryRepository {
	return &PostgreSQLCategoryRepository{
		client: client,
	}
}

// Create crea una nueva categoría
func (r *PostgreSQLCategoryRepository) Create(ctx context.Context, category *domain.ExpenseCategory) error {
	query := `
		INSERT INTO expense_categories (
			id, name, description, icon, color, parent_id,
			daily_limit, monthly_limit, requires_receipt, is_active, created, updated
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)`

	err := r.client.Exec(ctx, query,
		category.ID,
		category.Name,
		category.Description,
		category.Icon,
		category.Color,
		category.ParentID,
		category.DailyLimit,
		category.MonthlyLimit,
		category.RequiresReceipt,
		category.IsActive,
		category.Created,
		category.Updated,
	)

	if err != nil {
		return fmt.Errorf("error creando categoría: %w", err)
	}

	return nil
}

// GetByID obtiene una categoría por su ID
func (r *PostgreSQLCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategory, error) {
	query := `
		SELECT id, name, description, icon, color, parent_id,
			   daily_limit, monthly_limit, requires_receipt, is_active,
			   created, updated, deleted_at
		FROM expense_categories
		WHERE id = $1 AND deleted_at IS NULL`

	var category domain.ExpenseCategory
	var parentID sql.NullString

	err := r.client.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Icon,
		&category.Color,
		&parentID,
		&category.DailyLimit,
		&category.MonthlyLimit,
		&category.RequiresReceipt,
		&category.IsActive,
		&category.Created,
		&category.Updated,
		&category.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sharedErrors.ErrNotFound
		}
		return nil, fmt.Errorf("error obteniendo categoría: %w", err)
	}

	if parentID.Valid {
		pid, _ := uuid.Parse(parentID.String)
		category.ParentID = &pid
	}

	return &category, nil
}

// GetAll obtiene todas las categorías
func (r *PostgreSQLCategoryRepository) GetAll(ctx context.Context, activeOnly bool) ([]domain.ExpenseCategory, error) {
	query := `
		SELECT id, name, description, icon, color, parent_id,
			   daily_limit, monthly_limit, requires_receipt, is_active,
			   created, updated, deleted_at
		FROM expense_categories
		WHERE deleted_at IS NULL`

	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY name ASC"

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo categorías: %w", err)
	}
	defer rows.Close()

	categories := []domain.ExpenseCategory{}
	for rows.Next() {
		var category domain.ExpenseCategory
		var parentID sql.NullString

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Icon,
			&category.Color,
			&parentID,
			&category.DailyLimit,
			&category.MonthlyLimit,
			&category.RequiresReceipt,
			&category.IsActive,
			&category.Created,
			&category.Updated,
			&category.DeletedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error escaneando categoría: %w", err)
		}

		if parentID.Valid {
			pid, _ := uuid.Parse(parentID.String)
			category.ParentID = &pid
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// GetByParent obtiene categorías hijas por parent_id
func (r *PostgreSQLCategoryRepository) GetByParent(ctx context.Context, parentID *uuid.UUID) ([]domain.ExpenseCategory, error) {
	query := `
		SELECT id, name, description, icon, color, parent_id,
			   daily_limit, monthly_limit, requires_receipt, is_active,
			   created, updated, deleted_at
		FROM expense_categories
		WHERE deleted_at IS NULL`

	var args []interface{}
	if parentID == nil {
		// Obtener categorías raíz (sin parent)
		query += " AND parent_id IS NULL"
	} else {
		// Obtener subcategorías
		query += " AND parent_id = $1"
		args = append(args, *parentID)
	}

	query += " ORDER BY name ASC"

	rows, err := r.client.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo categorías por parent: %w", err)
	}
	defer rows.Close()

	categories := []domain.ExpenseCategory{}
	for rows.Next() {
		var category domain.ExpenseCategory
		var parentIDNull sql.NullString

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Icon,
			&category.Color,
			&parentIDNull,
			&category.DailyLimit,
			&category.MonthlyLimit,
			&category.RequiresReceipt,
			&category.IsActive,
			&category.Created,
			&category.Updated,
			&category.DeletedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error escaneando categoría: %w", err)
		}

		if parentIDNull.Valid {
			pid, _ := uuid.Parse(parentIDNull.String)
			category.ParentID = &pid
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// Update actualiza una categoría existente
func (r *PostgreSQLCategoryRepository) Update(ctx context.Context, category *domain.ExpenseCategory) error {
	query := `
		UPDATE expense_categories SET
			name = $1,
			description = $2,
			icon = $3,
			color = $4,
			parent_id = $5,
			daily_limit = $6,
			monthly_limit = $7,
			requires_receipt = $8,
			is_active = $9,
			updated = $10
		WHERE id = $11 AND deleted_at IS NULL`

	err := r.client.Exec(ctx, query,
		category.Name,
		category.Description,
		category.Icon,
		category.Color,
		category.ParentID,
		category.DailyLimit,
		category.MonthlyLimit,
		category.RequiresReceipt,
		category.IsActive,
		category.Updated,
		category.ID,
	)

	if err != nil {
		return fmt.Errorf("error actualizando categoría: %w", err)
	}

	return nil
}

// Delete elimina lógicamente una categoría (soft delete)
func (r *PostgreSQLCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE expense_categories SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	err := r.client.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error eliminando categoría: %w", err)
	}

	return nil
}

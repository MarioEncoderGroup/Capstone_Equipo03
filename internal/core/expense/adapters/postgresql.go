package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/JoseLuis21/mv-backend/internal/core/expense/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/expense/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PostgreSQLExpenseRepository implementa ExpenseRepository usando PostgreSQL
// Conecta con la base de datos del tenant
type PostgreSQLExpenseRepository struct {
	client *postgresql.PostgresqlClient
}

// NewPostgreSQLExpenseRepository crea una nueva instancia del repositorio
func NewPostgreSQLExpenseRepository(client *postgresql.PostgresqlClient) ports.ExpenseRepository {
	return &PostgreSQLExpenseRepository{
		client: client,
	}
}

// Create crea un nuevo gasto en la base de datos del tenant
func (r *PostgreSQLExpenseRepository) Create(ctx context.Context, expense *domain.Expense) error {
	query := `
		INSERT INTO expenses (
			id, user_id, policy_id, category_id, title, description,
			amount, currency, exchange_rate, expense_date, merchant_name,
			merchant_rut, receipt_number, payment_method, status,
			is_reimbursable, violation_reason, created, updated
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19
		)`

	err := r.client.Exec(ctx, query,
		expense.ID,
		expense.UserID,
		expense.PolicyID,
		expense.CategoryID,
		expense.Title,
		expense.Description,
		expense.Amount,
		expense.Currency,
		expense.ExchangeRate,
		expense.ExpenseDate,
		expense.MerchantName,
		expense.MerchantRUT,
		expense.ReceiptNumber,
		expense.PaymentMethod,
		expense.Status,
		expense.IsReimbursable,
		expense.ViolationReason,
		expense.Created,
		expense.Updated,
	)

	if err != nil {
		return fmt.Errorf("error creando gasto: %w", err)
	}

	return nil
}

// GetByID obtiene un gasto por su ID
func (r *PostgreSQLExpenseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	query := `
		SELECT
			e.id, e.user_id, e.policy_id, e.category_id, e.title, e.description,
			e.amount, e.currency, e.exchange_rate, e.amount_clp, e.expense_date,
			e.merchant_name, e.merchant_rut, e.receipt_number, e.payment_method,
			e.status, e.is_reimbursable, e.violation_reason, e.created, e.updated, e.deleted_at,
			c.id, c.name, c.description, c.icon, c.color, c.parent_id,
			c.daily_limit, c.monthly_limit, c.requires_receipt, c.is_active
		FROM expenses e
		LEFT JOIN expense_categories c ON e.category_id = c.id
		WHERE e.id = $1 AND e.deleted_at IS NULL`

	var expense domain.Expense
	var category domain.ExpenseCategory
	var categoryParentID sql.NullString

	err := r.client.QueryRow(ctx, query, id).Scan(
		&expense.ID,
		&expense.UserID,
		&expense.PolicyID,
		&expense.CategoryID,
		&expense.Title,
		&expense.Description,
		&expense.Amount,
		&expense.Currency,
		&expense.ExchangeRate,
		&expense.AmountCLP,
		&expense.ExpenseDate,
		&expense.MerchantName,
		&expense.MerchantRUT,
		&expense.ReceiptNumber,
		&expense.PaymentMethod,
		&expense.Status,
		&expense.IsReimbursable,
		&expense.ViolationReason,
		&expense.Created,
		&expense.Updated,
		&expense.DeletedAt,
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Icon,
		&category.Color,
		&categoryParentID,
		&category.DailyLimit,
		&category.MonthlyLimit,
		&category.RequiresReceipt,
		&category.IsActive,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sharedErrors.ErrNotFound
		}
		return nil, fmt.Errorf("error obteniendo gasto: %w", err)
	}

	if categoryParentID.Valid {
		parentID, _ := uuid.Parse(categoryParentID.String)
		category.ParentID = &parentID
	}

	expense.Category = &category

	// Obtener comprobantes
	receipts, err := r.GetReceipts(ctx, expense.ID)
	if err != nil {
		return nil, err
	}
	expense.Receipts = receipts

	return &expense, nil
}

// GetAll obtiene gastos con filtros y paginación
func (r *PostgreSQLExpenseRepository) GetAll(ctx context.Context, filters domain.ExpenseFilters) ([]domain.Expense, int, error) {
	// Construir query dinámico con filtros
	whereConditions := []string{"e.deleted_at IS NULL"}
	args := []interface{}{}
	argIndex := 1

	if filters.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("e.user_id = $%d", argIndex))
		args = append(args, *filters.UserID)
		argIndex++
	}

	if filters.CategoryID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("e.category_id = $%d", argIndex))
		args = append(args, *filters.CategoryID)
		argIndex++
	}

	if filters.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("e.status = $%d", argIndex))
		args = append(args, *filters.Status)
		argIndex++
	}

	if filters.DateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("e.expense_date >= $%d", argIndex))
		args = append(args, *filters.DateFrom)
		argIndex++
	}

	if filters.DateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("e.expense_date <= $%d", argIndex))
		args = append(args, *filters.DateTo)
		argIndex++
	}

	if filters.MinAmount != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("e.amount >= $%d", argIndex))
		args = append(args, *filters.MinAmount)
		argIndex++
	}

	if filters.MaxAmount != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("e.amount <= $%d", argIndex))
		args = append(args, *filters.MaxAmount)
		argIndex++
	}

	if filters.Search != nil && *filters.Search != "" {
		searchPattern := "%" + *filters.Search + "%"
		whereConditions = append(whereConditions, fmt.Sprintf("(e.title ILIKE $%d OR e.merchant_name ILIKE $%d)", argIndex, argIndex))
		args = append(args, searchPattern)
		argIndex++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM expenses e WHERE %s", whereClause)
	var total int
	err := r.client.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error contando gastos: %w", err)
	}

	// Query principal
	query := fmt.Sprintf(`
		SELECT
			e.id, e.user_id, e.policy_id, e.category_id, e.title, e.description,
			e.amount, e.currency, e.exchange_rate, e.amount_clp, e.expense_date,
			e.merchant_name, e.merchant_rut, e.receipt_number, e.payment_method,
			e.status, e.is_reimbursable, e.violation_reason, e.created, e.updated, e.deleted_at,
			c.id, c.name, c.description, c.icon, c.color, c.parent_id,
			c.daily_limit, c.monthly_limit, c.requires_receipt, c.is_active
		FROM expenses e
		LEFT JOIN expense_categories c ON e.category_id = c.id
		WHERE %s
		ORDER BY e.created DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, filters.Limit, filters.Offset)

	rows, err := r.client.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error obteniendo gastos: %w", err)
	}
	defer rows.Close()

	expenses := []domain.Expense{}
	for rows.Next() {
		var expense domain.Expense
		var category domain.ExpenseCategory
		var categoryParentID sql.NullString

		err := rows.Scan(
			&expense.ID,
			&expense.UserID,
			&expense.PolicyID,
			&expense.CategoryID,
			&expense.Title,
			&expense.Description,
			&expense.Amount,
			&expense.Currency,
			&expense.ExchangeRate,
			&expense.AmountCLP,
			&expense.ExpenseDate,
			&expense.MerchantName,
			&expense.MerchantRUT,
			&expense.ReceiptNumber,
			&expense.PaymentMethod,
			&expense.Status,
			&expense.IsReimbursable,
			&expense.ViolationReason,
			&expense.Created,
			&expense.Updated,
			&expense.DeletedAt,
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Icon,
			&category.Color,
			&categoryParentID,
			&category.DailyLimit,
			&category.MonthlyLimit,
			&category.RequiresReceipt,
			&category.IsActive,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("error escaneando gasto: %w", err)
		}

		if categoryParentID.Valid {
			parentID, _ := uuid.Parse(categoryParentID.String)
			category.ParentID = &parentID
		}

		expense.Category = &category
		expenses = append(expenses, expense)
	}

	return expenses, total, nil
}

// Update actualiza un gasto existente
func (r *PostgreSQLExpenseRepository) Update(ctx context.Context, expense *domain.Expense) error {
	query := `
		UPDATE expenses SET
			category_id = $1,
			title = $2,
			description = $3,
			amount = $4,
			currency = $5,
			exchange_rate = $6,
			expense_date = $7,
			merchant_name = $8,
			merchant_rut = $9,
			receipt_number = $10,
			payment_method = $11,
			status = $12,
			is_reimbursable = $13,
			violation_reason = $14,
			updated = $15
		WHERE id = $16 AND deleted_at IS NULL`

	err := r.client.Exec(ctx, query,
		expense.CategoryID,
		expense.Title,
		expense.Description,
		expense.Amount,
		expense.Currency,
		expense.ExchangeRate,
		expense.ExpenseDate,
		expense.MerchantName,
		expense.MerchantRUT,
		expense.ReceiptNumber,
		expense.PaymentMethod,
		expense.Status,
		expense.IsReimbursable,
		expense.ViolationReason,
		expense.Updated,
		expense.ID,
	)

	if err != nil {
		return fmt.Errorf("error actualizando gasto: %w", err)
	}

	return nil
}

// Delete elimina lógicamente un gasto (soft delete)
func (r *PostgreSQLExpenseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE expenses SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	err := r.client.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error eliminando gasto: %w", err)
	}

	return nil
}

// AddReceipt agrega un comprobante a un gasto
func (r *PostgreSQLExpenseRepository) AddReceipt(ctx context.Context, receipt *domain.ExpenseReceipt) error {
	query := `
		INSERT INTO expense_receipts (
			id, expense_id, file_url, file_name, file_type,
			file_size, ocr_data, ocr_confidence, is_primary, created
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)`

	err := r.client.Exec(ctx, query,
		receipt.ID,
		receipt.ExpenseID,
		receipt.FileURL,
		receipt.FileName,
		receipt.FileType,
		receipt.FileSize,
		receipt.OCRData,
		receipt.OCRConfidence,
		receipt.IsPrimary,
		receipt.Created,
	)

	if err != nil {
		return fmt.Errorf("error agregando comprobante: %w", err)
	}

	return nil
}

// GetReceipts obtiene todos los comprobantes de un gasto
func (r *PostgreSQLExpenseRepository) GetReceipts(ctx context.Context, expenseID uuid.UUID) ([]domain.ExpenseReceipt, error) {
	query := `
		SELECT id, expense_id, file_url, file_name, file_type,
			   file_size, ocr_data, ocr_confidence, is_primary, created
		FROM expense_receipts
		WHERE expense_id = $1
		ORDER BY is_primary DESC, created DESC`

	rows, err := r.client.Query(ctx, query, expenseID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo comprobantes: %w", err)
	}
	defer rows.Close()

	receipts := []domain.ExpenseReceipt{}
	for rows.Next() {
		var receipt domain.ExpenseReceipt
		err := rows.Scan(
			&receipt.ID,
			&receipt.ExpenseID,
			&receipt.FileURL,
			&receipt.FileName,
			&receipt.FileType,
			&receipt.FileSize,
			&receipt.OCRData,
			&receipt.OCRConfidence,
			&receipt.IsPrimary,
			&receipt.Created,
		)

		if err != nil {
			return nil, fmt.Errorf("error escaneando comprobante: %w", err)
		}

		receipts = append(receipts, receipt)
	}

	return receipts, nil
}

// DeleteReceipt elimina un comprobante
func (r *PostgreSQLExpenseRepository) DeleteReceipt(ctx context.Context, receiptID uuid.UUID) error {
	query := `DELETE FROM expense_receipts WHERE id = $1`

	err := r.client.Exec(ctx, query, receiptID)
	if err != nil {
		return fmt.Errorf("error eliminando comprobante: %w", err)
	}

	return nil
}

// SetPrimaryReceipt establece un comprobante como principal
func (r *PostgreSQLExpenseRepository) SetPrimaryReceipt(ctx context.Context, receiptID uuid.UUID) error {
	// Obtener expense_id del receipt
	var expenseID uuid.UUID
	queryGetExpenseID := `SELECT expense_id FROM expense_receipts WHERE id = $1`
	err := r.client.QueryRow(ctx, queryGetExpenseID, receiptID).Scan(&expenseID)
	if err != nil {
		return fmt.Errorf("error obteniendo expense_id: %w", err)
	}

	// Usar transacción para asegurar atomicidad
	tx, err := r.client.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error iniciando transacción: %w", err)
	}
	defer tx.Rollback(ctx)

	// Quitar is_primary de todos los receipts del expense
	queryUnsetPrimary := `UPDATE expense_receipts SET is_primary = false WHERE expense_id = $1`
	_, err = tx.Exec(ctx, queryUnsetPrimary, expenseID)
	if err != nil {
		return fmt.Errorf("error quitando is_primary: %w", err)
	}

	// Establecer nuevo receipt como primary
	querySetPrimary := `UPDATE expense_receipts SET is_primary = true WHERE id = $1`
	_, err = tx.Exec(ctx, querySetPrimary, receiptID)
	if err != nil {
		return fmt.Errorf("error estableciendo is_primary: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committeando transacción: %w", err)
	}

	return nil
}

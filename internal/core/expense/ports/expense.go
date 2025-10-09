package ports

import (
	"context"

	"github.com/JoseLuis21/mv-backend/internal/core/expense/domain"
	"github.com/google/uuid"
)

// ExpenseRepository define el contrato para la persistencia de gastos
type ExpenseRepository interface {
	// Create crea un nuevo gasto en la base de datos del tenant
	Create(ctx context.Context, expense *domain.Expense) error

	// GetByID obtiene un gasto por su ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error)

	// GetAll obtiene gastos con filtros y paginación
	GetAll(ctx context.Context, filters domain.ExpenseFilters) ([]domain.Expense, int, error)

	// Update actualiza un gasto existente
	Update(ctx context.Context, expense *domain.Expense) error

	// Delete elimina lógicamente un gasto (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// Receipts - Gestión de comprobantes
	// AddReceipt agrega un comprobante a un gasto
	AddReceipt(ctx context.Context, receipt *domain.ExpenseReceipt) error

	// GetReceipts obtiene todos los comprobantes de un gasto
	GetReceipts(ctx context.Context, expenseID uuid.UUID) ([]domain.ExpenseReceipt, error)

	// DeleteReceipt elimina un comprobante
	DeleteReceipt(ctx context.Context, receiptID uuid.UUID) error

	// SetPrimaryReceipt establece un comprobante como principal
	SetPrimaryReceipt(ctx context.Context, receiptID uuid.UUID) error
}

// ExpenseService define el contrato para la lógica de negocio de gastos
type ExpenseService interface {
	// CRUD
	Create(ctx context.Context, dto domain.CreateExpenseDto, userID uuid.UUID) (*domain.Expense, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Expense, error)
	GetAll(ctx context.Context, filters domain.ExpenseFilters) ([]domain.Expense, int, error)
	Update(ctx context.Context, id uuid.UUID, dto domain.UpdateExpenseDto, userID uuid.UUID) (*domain.Expense, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// Receipts
	UploadReceipt(ctx context.Context, dto domain.UploadReceiptDto, fileData []byte) (*domain.ExpenseReceipt, error)
	DeleteReceipt(ctx context.Context, receiptID uuid.UUID, userID uuid.UUID) error

	// Business logic
	CanEdit(ctx context.Context, expenseID uuid.UUID, userID uuid.UUID) (bool, error)
	CanDelete(ctx context.Context, expenseID uuid.UUID, userID uuid.UUID) (bool, error)
	ChangeStatus(ctx context.Context, expenseID uuid.UUID, newStatus domain.ExpenseStatus) error
}

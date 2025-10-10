package domain

import (
	"time"

	"github.com/google/uuid"
)

// CreateExpenseDto - Request para crear gasto
type CreateExpenseDto struct {
	CategoryID     uuid.UUID     `json:"category_id" validate:"required,uuid"`
	Title          string        `json:"title" validate:"required,min=3,max=200"`
	Description    string        `json:"description" validate:"omitempty,max=500"`
	Amount         float64       `json:"amount" validate:"required,gt=0"`
	Currency       string        `json:"currency" validate:"required,oneof=CLP USD EUR"`
	ExpenseDate    time.Time     `json:"expense_date" validate:"required"`
	MerchantName   string        `json:"merchant_name" validate:"omitempty,max=200"`
	MerchantRUT    string        `json:"merchant_rut" validate:"omitempty"`
	ReceiptNumber  string        `json:"receipt_number" validate:"omitempty,max=100"`
	PaymentMethod  PaymentMethod `json:"payment_method" validate:"required,oneof=cash card transfer"`
	IsReimbursable bool          `json:"is_reimbursable"`
}

// UpdateExpenseDto - Request para actualizar gasto
type UpdateExpenseDto struct {
	CategoryID     *uuid.UUID     `json:"category_id" validate:"omitempty,uuid"`
	Title          *string        `json:"title" validate:"omitempty,min=3,max=200"`
	Description    *string        `json:"description" validate:"omitempty,max=500"`
	Amount         *float64       `json:"amount" validate:"omitempty,gt=0"`
	Currency       *string        `json:"currency" validate:"omitempty,oneof=CLP USD EUR"`
	ExpenseDate    *time.Time     `json:"expense_date"`
	MerchantName   *string        `json:"merchant_name" validate:"omitempty,max=200"`
	MerchantRUT    *string        `json:"merchant_rut" validate:"omitempty"`
	ReceiptNumber  *string        `json:"receipt_number" validate:"omitempty,max=100"`
	PaymentMethod  *PaymentMethod `json:"payment_method" validate:"omitempty,oneof=cash card transfer"`
	IsReimbursable *bool          `json:"is_reimbursable"`
}

// ExpenseFilters - Filtros para listar gastos
type ExpenseFilters struct {
	UserID     *uuid.UUID
	CategoryID *uuid.UUID
	Status     *ExpenseStatus
	DateFrom   *time.Time
	DateTo     *time.Time
	MinAmount  *float64
	MaxAmount  *float64
	Search     *string
	Limit      int
	Offset     int
}

// UploadReceiptDto - Request para subir comprobante
type UploadReceiptDto struct {
	ExpenseID uuid.UUID `json:"expense_id" validate:"required,uuid"`
	FileName  string    `json:"file_name" validate:"required"`
	FileType  string    `json:"file_type" validate:"required,oneof=image/jpeg image/png image/jpg application/pdf"`
	FileSize  int64     `json:"file_size" validate:"required,gt=0,lte=10485760"` // Max 10MB
	IsPrimary bool      `json:"is_primary"`
}

// ExpenseResponse - Response con datos completos
type ExpenseResponse struct {
	Expense
	Receipts []ExpenseReceipt  `json:"receipts,omitempty"`
	Category *ExpenseCategory  `json:"category,omitempty"`
}

// ExpensesResponse - Response para lista
type ExpensesResponse struct {
	Expenses []ExpenseResponse `json:"expenses"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

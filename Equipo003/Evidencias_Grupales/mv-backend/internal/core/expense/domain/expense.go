package domain

import (
	"time"

	"github.com/google/uuid"
)

// ExpenseStatus representa el estado del gasto
type ExpenseStatus string

const (
	ExpenseStatusDraft      ExpenseStatus = "draft"
	ExpenseStatusSubmitted  ExpenseStatus = "submitted"
	ExpenseStatusApproved   ExpenseStatus = "approved"
	ExpenseStatusRejected   ExpenseStatus = "rejected"
	ExpenseStatusReimbursed ExpenseStatus = "reimbursed"
)

// PaymentMethod representa el método de pago
type PaymentMethod string

const (
	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodCard     PaymentMethod = "card"
	PaymentMethodTransfer PaymentMethod = "transfer"
)

// Expense representa un gasto empresarial
type Expense struct {
	ID              uuid.UUID      `json:"id"`
	UserID          uuid.UUID      `json:"user_id"`
	PolicyID        *uuid.UUID     `json:"policy_id,omitempty"`
	CategoryID      uuid.UUID      `json:"category_id"`
	Title           string         `json:"title"`
	Description     *string        `json:"description,omitempty"`
	Amount          float64        `json:"amount"`
	Currency        string         `json:"currency"`
	ExchangeRate    float64        `json:"exchange_rate"`
	AmountCLP       float64        `json:"amount_clp"`
	ExpenseDate     time.Time      `json:"expense_date"`
	MerchantName    *string        `json:"merchant_name,omitempty"`
	MerchantRUT     *string        `json:"merchant_rut,omitempty"`
	ReceiptNumber   *string        `json:"receipt_number,omitempty"`
	PaymentMethod   PaymentMethod  `json:"payment_method"`
	Status          ExpenseStatus  `json:"status"`
	IsReimbursable  bool           `json:"is_reimbursable"`
	ViolationReason *string        `json:"violation_reason,omitempty"`
	Receipts        []ExpenseReceipt `json:"receipts,omitempty"`
	Category        *ExpenseCategory `json:"category,omitempty"`
	Created         time.Time      `json:"created_at"`
	Updated         time.Time      `json:"updated_at"`
	DeletedAt       *time.Time     `json:"deleted_at,omitempty"`
}

// ExpenseReceipt representa un comprobante de gasto
type ExpenseReceipt struct {
	ID            uuid.UUID       `json:"id"`
	ExpenseID     uuid.UUID       `json:"expense_id"`
	FileURL       string          `json:"file_url"`
	FileName      string          `json:"file_name"`
	FileType      string          `json:"file_type"`
	FileSize      int64           `json:"file_size"`
	OCRData       *map[string]any `json:"ocr_data,omitempty"`
	OCRConfidence *float64        `json:"ocr_confidence,omitempty"`
	IsPrimary     bool            `json:"is_primary"`
	Created       time.Time       `json:"created_at"`
}

// ExpenseCategory representa una categoría de gasto
type ExpenseCategory struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	Icon            *string    `json:"icon,omitempty"`
	Color           *string    `json:"color,omitempty"`
	ParentID        *uuid.UUID `json:"parent_id,omitempty"`
	DailyLimit      *float64   `json:"daily_limit,omitempty"`
	MonthlyLimit    *float64   `json:"monthly_limit,omitempty"`
	RequiresReceipt bool       `json:"requires_receipt"`
	IsActive        bool       `json:"is_active"`
	Created         time.Time  `json:"created_at"`
	Updated         time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// IsValidStatus verifica si el status es válido
func (s ExpenseStatus) IsValid() bool {
	switch s {
	case ExpenseStatusDraft, ExpenseStatusSubmitted, ExpenseStatusApproved, ExpenseStatusRejected, ExpenseStatusReimbursed:
		return true
	}
	return false
}

// IsValidPaymentMethod verifica si el método de pago es válido
func (p PaymentMethod) IsValid() bool {
	switch p {
	case PaymentMethodCash, PaymentMethodCard, PaymentMethodTransfer:
		return true
	}
	return false
}

// CanBeEdited verifica si el gasto puede ser editado
func (e *Expense) CanBeEdited() bool {
	return e.Status == ExpenseStatusDraft || e.Status == ExpenseStatusRejected
}

// CanBeDeleted verifica si el gasto puede ser eliminado
func (e *Expense) CanBeDeleted() bool {
	return e.Status == ExpenseStatusDraft || e.Status == ExpenseStatusRejected
}

// IsDeleted verifica si el gasto está eliminado
func (e *Expense) IsDeleted() bool {
	return e.DeletedAt != nil
}

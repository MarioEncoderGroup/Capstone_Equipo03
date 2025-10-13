package domain

import (
	"github.com/google/uuid"
)

// CreateReportDto representa los datos para crear un reporte
type CreateReportDto struct {
	Title       string     `json:"title" validate:"required,min=3,max=200"`
	Description string     `json:"description" validate:"omitempty,max=500"`
	PolicyID    *uuid.UUID `json:"policy_id" validate:"omitempty,uuid"`
	Currency    string     `json:"currency" validate:"omitempty,oneof=CLP USD EUR"`
}

// UpdateReportDto representa los datos para actualizar un reporte
type UpdateReportDto struct {
	Title       *string `json:"title" validate:"omitempty,min=3,max=200"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// AddExpensesDto representa los datos para agregar gastos a un reporte
type AddExpensesDto struct {
	ExpenseIDs []uuid.UUID `json:"expense_ids" validate:"required,min=1,dive,uuid"`
}

// ReportFilters representa los filtros para buscar reportes
type ReportFilters struct {
	UserID *uuid.UUID
	Status *ReportStatus
	Limit  int
	Offset int
}

// ReportResponse representa la respuesta con datos completos de un reporte
type ReportResponse struct {
	ExpenseReport
	Items     []ExpenseReportItem `json:"items,omitempty"`
	Approvals []Approval          `json:"approvals,omitempty"`
	Comments  []ExpenseComment    `json:"comments,omitempty"`
}

// ReportsResponse representa la respuesta para una lista de reportes
type ReportsResponse struct {
	Reports []ReportResponse `json:"reports"`
	Total   int              `json:"total"`
	Limit   int              `json:"limit"`
	Offset  int              `json:"offset"`
}

// CreateCommentDto representa los datos para crear un comentario
type CreateCommentDto struct {
	ReportID    *uuid.UUID  `json:"report_id" validate:"omitempty,uuid"`
	ExpenseID   *uuid.UUID  `json:"expense_id" validate:"omitempty,uuid"`
	Content     string      `json:"content" validate:"required,min=1,max=2000"`
	CommentType CommentType `json:"comment_type" validate:"omitempty,oneof=general question clarification approval_note rejection_note system"`
	ParentID    *uuid.UUID  `json:"parent_id" validate:"omitempty,uuid"`
	IsInternal  bool        `json:"is_internal"`
}

// ApproveReportDto representa los datos para aprobar un reporte
type ApproveReportDto struct {
	Comments       string   `json:"comments" validate:"omitempty,max=500"`
	ApprovedAmount *float64 `json:"approved_amount" validate:"omitempty,gt=0"`
}

// RejectReportDto representa los datos para rechazar un reporte
type RejectReportDto struct {
	Reason string `json:"reason" validate:"required,min=10,max=500"`
}

// SubmitReportResponse representa la respuesta al enviar un reporte
type SubmitReportResponse struct {
	Report          ExpenseReport `json:"report"`
	ApprovalsCreated []Approval    `json:"approvals_created"`
}

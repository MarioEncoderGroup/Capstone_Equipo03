package domain

import (
	"time"

	"github.com/google/uuid"
)

// ReportStatus representa el estado de un reporte de gastos
type ReportStatus string

const (
	ReportStatusDraft       ReportStatus = "draft"
	ReportStatusSubmitted   ReportStatus = "submitted"
	ReportStatusUnderReview ReportStatus = "under_review"
	ReportStatusApproved    ReportStatus = "approved"
	ReportStatusRejected    ReportStatus = "rejected"
	ReportStatusPaid        ReportStatus = "paid"
)

// ApprovalStatus representa el estado de una aprobación
type ApprovalStatus string

const (
	ApprovalStatusPending   ApprovalStatus = "pending"
	ApprovalStatusApproved  ApprovalStatus = "approved"
	ApprovalStatusRejected  ApprovalStatus = "rejected"
	ApprovalStatusEscalated ApprovalStatus = "escalated"
)

// ApprovalAction representa un tipo de acción en el historial
type ApprovalAction string

const (
	ApprovalActionCreated    ApprovalAction = "created"
	ApprovalActionApproved   ApprovalAction = "approved"
	ApprovalActionRejected   ApprovalAction = "rejected"
	ApprovalActionEscalated  ApprovalAction = "escalated"
	ApprovalActionReassigned ApprovalAction = "reassigned"
	ApprovalActionCommented  ApprovalAction = "commented"
)

// CommentType representa el tipo de comentario
type CommentType string

const (
	CommentTypeGeneral       CommentType = "general"
	CommentTypeQuestion      CommentType = "question"
	CommentTypeClarification CommentType = "clarification"
	CommentTypeApprovalNote  CommentType = "approval_note"
	CommentTypeRejectionNote CommentType = "rejection_note"
	CommentTypeSystem        CommentType = "system"
)

// ExpenseReport representa un reporte de gastos
type ExpenseReport struct {
	ID             uuid.UUID     `json:"id"`
	UserID         uuid.UUID     `json:"user_id"`
	PolicyID       *uuid.UUID    `json:"policy_id,omitempty"`
	Title          string        `json:"title"`
	Description    *string       `json:"description,omitempty"`
	Status         ReportStatus  `json:"status"`
	TotalAmount    float64       `json:"total_amount"`
	Currency       string        `json:"currency"`
	SubmissionDate *time.Time    `json:"submission_date,omitempty"`
	ApprovalDate   *time.Time    `json:"approval_date,omitempty"`
	PaymentDate    *time.Time    `json:"payment_date,omitempty"`
	RejectionReason *string      `json:"rejection_reason,omitempty"`
	Created        time.Time     `json:"created_at"`
	Updated        time.Time     `json:"updated_at"`
	DeletedAt      *time.Time    `json:"deleted_at,omitempty"`

	// Relaciones
	Items     []ExpenseReportItem `json:"items,omitempty"`
	Approvals []Approval          `json:"approvals,omitempty"`
	Comments  []ExpenseComment    `json:"comments,omitempty"`
}

// ExpenseReportItem representa un gasto dentro de un reporte
type ExpenseReportItem struct {
	ID             uuid.UUID `json:"id"`
	ReportID       uuid.UUID `json:"report_id"`
	ExpenseID      uuid.UUID `json:"expense_id"`
	SequenceNumber int       `json:"sequence_number"`
	Created        time.Time `json:"created_at"`
}

// Approval representa una aprobación en el flujo
type Approval struct {
	ID             uuid.UUID       `json:"id"`
	ReportID       uuid.UUID       `json:"report_id"`
	ApproverID     uuid.UUID       `json:"approver_id"`
	Level          int             `json:"level"`
	Status         ApprovalStatus  `json:"status"`
	Comments       *string         `json:"comments,omitempty"`
	ApprovedAmount *float64        `json:"approved_amount,omitempty"`
	DecisionDate   *time.Time      `json:"decision_date,omitempty"`
	EscalationDate *time.Time      `json:"escalation_date,omitempty"`
	EscalatedTo    *uuid.UUID      `json:"escalated_to,omitempty"`
	Created        time.Time       `json:"created_at"`
	Updated        time.Time       `json:"updated_at"`
}

// ApprovalHistory representa el historial de una aprobación
type ApprovalHistory struct {
	ID             uuid.UUID       `json:"id"`
	ApprovalID     uuid.UUID       `json:"approval_id"`
	ReportID       uuid.UUID       `json:"report_id"`
	ActorID        uuid.UUID       `json:"actor_id"`
	Action         ApprovalAction  `json:"action"`
	PreviousStatus *ApprovalStatus `json:"previous_status,omitempty"`
	NewStatus      *ApprovalStatus `json:"new_status,omitempty"`
	Comments       *string         `json:"comments,omitempty"`
	Metadata       map[string]any  `json:"metadata,omitempty"`
	Created        time.Time       `json:"created_at"`
}

// ExpenseComment representa un comentario sobre un gasto o reporte
type ExpenseComment struct {
	ID          uuid.UUID    `json:"id"`
	ReportID    *uuid.UUID   `json:"report_id,omitempty"`
	ExpenseID   *uuid.UUID   `json:"expense_id,omitempty"`
	UserID      uuid.UUID    `json:"user_id"`
	CommentType CommentType  `json:"comment_type"`
	Content     string       `json:"content"`
	ParentID    *uuid.UUID   `json:"parent_id,omitempty"`
	IsInternal  bool         `json:"is_internal"`
	Attachments map[string]any `json:"attachments,omitempty"`
	Created     time.Time    `json:"created_at"`
	Updated     time.Time    `json:"updated_at"`
	DeletedAt   *time.Time   `json:"deleted_at,omitempty"`
}

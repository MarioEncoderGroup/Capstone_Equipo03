package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/report/domain"
)

// ReportRepository define las operaciones de persistencia para reportes de gastos
type ReportRepository interface {
	// Reports
	Create(ctx context.Context, report *domain.ExpenseReport) error
	GetByID(ctx context.Context, reportID uuid.UUID) (*domain.ExpenseReport, error)
	GetByUser(ctx context.Context, filters *domain.ReportFilters) ([]domain.ExpenseReport, int, error)
	Update(ctx context.Context, report *domain.ExpenseReport) error
	Delete(ctx context.Context, reportID uuid.UUID) error
	UpdateStatus(ctx context.Context, reportID uuid.UUID, status domain.ReportStatus) error
	RecalculateTotal(ctx context.Context, reportID uuid.UUID) (float64, error)

	// Report Items (expenses)
	AddExpenseToReport(ctx context.Context, reportID, expenseID uuid.UUID) error
	RemoveExpenseFromReport(ctx context.Context, reportID, expenseID uuid.UUID) error
	GetReportExpenses(ctx context.Context, reportID uuid.UUID) ([]domain.ExpenseReportItem, error)
	IsExpenseInReport(ctx context.Context, expenseID uuid.UUID) (bool, *uuid.UUID, error)

	// Approvals
	CreateApproval(ctx context.Context, approval *domain.Approval) error
	GetApprovalsByReport(ctx context.Context, reportID uuid.UUID) ([]domain.Approval, error)
	GetApprovalByID(ctx context.Context, approvalID uuid.UUID) (*domain.Approval, error)
	UpdateApproval(ctx context.Context, approval *domain.Approval) error
	GetPendingApprovalForUser(ctx context.Context, reportID, userID uuid.UUID) (*domain.Approval, error)
	GetPendingApprovalsByApprover(ctx context.Context, approverID uuid.UUID) ([]domain.Approval, error)
	GetPendingApprovalsByReportID(ctx context.Context, reportID uuid.UUID) ([]domain.Approval, error)
	GetStaleApprovals(ctx context.Context, hoursThreshold int) ([]domain.Approval, error)

	// Approval History
	CreateApprovalHistory(ctx context.Context, history *domain.ApprovalHistory) error
	GetApprovalHistory(ctx context.Context, approvalID uuid.UUID) ([]domain.ApprovalHistory, error)

	// Comments
	CreateComment(ctx context.Context, comment *domain.ExpenseComment) error
	GetCommentsByReport(ctx context.Context, reportID uuid.UUID, includeInternal bool) ([]domain.ExpenseComment, error)
	GetCommentsByExpense(ctx context.Context, expenseID uuid.UUID) ([]domain.ExpenseComment, error)
	UpdateComment(ctx context.Context, comment *domain.ExpenseComment) error
	DeleteComment(ctx context.Context, commentID uuid.UUID) error
}

// ReportService define la lógica de negocio para reportes de gastos
type ReportService interface {
	// CRUD básico
	CreateReport(ctx context.Context, userID uuid.UUID, dto *domain.CreateReportDto) (*domain.ExpenseReport, error)
	GetReportByID(ctx context.Context, reportID uuid.UUID) (*domain.ReportResponse, error)
	GetUserReports(ctx context.Context, filters *domain.ReportFilters) (*domain.ReportsResponse, error)
	UpdateReport(ctx context.Context, reportID uuid.UUID, dto *domain.UpdateReportDto) error
	DeleteReport(ctx context.Context, reportID uuid.UUID) error

	// Gestión de gastos en reporte
	AddExpensesToReport(ctx context.Context, reportID uuid.UUID, dto *domain.AddExpensesDto) error
	RemoveExpenseFromReport(ctx context.Context, reportID, expenseID uuid.UUID) error

	// Workflow de aprobación
	SubmitReport(ctx context.Context, reportID uuid.UUID) (*domain.SubmitReportResponse, error)
	ApproveReport(ctx context.Context, approvalID, approverID uuid.UUID, dto *domain.ApproveReportDto) error
	RejectReport(ctx context.Context, approvalID, approverID uuid.UUID, dto *domain.RejectReportDto) error

	// Comentarios
	AddComment(ctx context.Context, userID uuid.UUID, dto *domain.CreateCommentDto) (*domain.ExpenseComment, error)
	GetReportComments(ctx context.Context, reportID uuid.UUID, includeInternal bool) ([]domain.ExpenseComment, error)

	// Validaciones de negocio
	CanEditReport(ctx context.Context, reportID, userID uuid.UUID) (bool, error)
	CanDeleteReport(ctx context.Context, reportID, userID uuid.UUID) (bool, error)
	CanSubmitReport(ctx context.Context, reportID uuid.UUID) (bool, error)
}

// ApprovalService define la lógica de negocio para gestión de aprobaciones
type ApprovalService interface {
	// GetPendingApprovals retorna aprobaciones pendientes del usuario
	GetPendingApprovals(ctx context.Context, approverID uuid.UUID) ([]domain.Approval, error)

	// GetPendingApprovalsByReport retorna aprobaciones pendientes de un reporte específico
	GetPendingApprovalsByReport(ctx context.Context, reportID uuid.UUID) ([]domain.Approval, error)

	// Approve aprueba una solicitud
	Approve(ctx context.Context, approvalID, approverID uuid.UUID, dto *domain.ApproveReportDto) error

	// Reject rechaza una solicitud
	Reject(ctx context.Context, approvalID, approverID uuid.UUID, dto *domain.RejectReportDto) error

	// GetHistory retorna historial de aprobaciones de un reporte
	GetHistory(ctx context.Context, reportID uuid.UUID) ([]domain.ApprovalHistory, error)

	// GetApprovalHistory retorna historial de una aprobación específica
	GetApprovalHistory(ctx context.Context, approvalID uuid.UUID) ([]domain.ApprovalHistory, error)

	// Escalate escala una aprobación a otro aprobador
	Escalate(ctx context.Context, approvalID, currentApproverID, newApproverID uuid.UUID, reason string) error

	// GetApprovalsByReport retorna todas las aprobaciones de un reporte
	GetApprovalsByReport(ctx context.Context, reportID uuid.UUID) ([]domain.Approval, error)
}

// WorkflowEngine define la lógica del motor de workflows para aprobaciones
type WorkflowEngine interface {
	// CreateApprovals crea las aprobaciones necesarias según política y monto del reporte
	CreateApprovals(ctx context.Context, report *domain.ExpenseReport) ([]domain.Approval, error)

	// ProcessApproval procesa una aprobación y determina el siguiente paso en el workflow
	ProcessApproval(ctx context.Context, approval *domain.Approval) error

	// EscalateApproval escala una aprobación al siguiente nivel automáticamente
	EscalateApproval(ctx context.Context, approvalID uuid.UUID) error

	// GetNextApprover determina el siguiente aprobador según la política y nivel
	GetNextApprover(ctx context.Context, reportID uuid.UUID, currentLevel int) (*uuid.UUID, error)
}

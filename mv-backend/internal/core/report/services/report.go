package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/report/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/report/ports"
)

var (
	ErrReportNotFound         = errors.New("report not found")
	ErrReportNotEditable      = errors.New("report cannot be edited in current status")
	ErrReportNotDeletable     = errors.New("report cannot be deleted in current status")
	ErrReportNotSubmittable   = errors.New("report cannot be submitted")
	ErrReportEmpty            = errors.New("report must have at least one expense")
	ErrExpenseAlreadyInReport = errors.New("expense is already in another report")
	ErrExpenseNotInReport     = errors.New("expense is not in this report")
	ErrApprovalNotFound       = errors.New("approval not found")
	ErrUnauthorized           = errors.New("user is not authorized for this action")
	ErrInvalidApprovalStatus  = errors.New("invalid approval status for this action")
)

type reportService struct {
	repo ports.ReportRepository
}

// NewReportService crea una nueva instancia del servicio de reportes
func NewReportService(repo ports.ReportRepository) ports.ReportService {
	return &reportService{repo: repo}
}

// CreateReport crea un nuevo reporte de gastos en estado draft
func (s *reportService) CreateReport(ctx context.Context, userID uuid.UUID, dto *domain.CreateReportDto) (*domain.ExpenseReport, error) {
	now := time.Now()
	currency := "CLP"
	if dto.Currency != "" {
		currency = dto.Currency
	}

	report := &domain.ExpenseReport{
		ID:          uuid.New(),
		UserID:      userID,
		PolicyID:    dto.PolicyID,
		Title:       dto.Title,
		Description: nil,
		Status:      domain.ReportStatusDraft,
		TotalAmount: 0,
		Currency:    currency,
		Created:     now,
		Updated:     now,
	}

	if dto.Description != "" {
		report.Description = &dto.Description
	}

	if err := s.repo.Create(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	return report, nil
}

// GetReportByID obtiene un reporte por su ID con todas sus relaciones
func (s *reportService) GetReportByID(ctx context.Context, reportID uuid.UUID) (*domain.ReportResponse, error) {
	report, err := s.repo.GetByID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return nil, ErrReportNotFound
	}

	// Obtener items, approvals y comments
	items, err := s.repo.GetReportExpenses(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report expenses: %w", err)
	}

	approvals, err := s.repo.GetApprovalsByReport(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get approvals: %w", err)
	}

	comments, err := s.repo.GetCommentsByReport(ctx, reportID, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	return &domain.ReportResponse{
		ExpenseReport: *report,
		Items:         items,
		Approvals:     approvals,
		Comments:      comments,
	}, nil
}

// GetUserReports obtiene los reportes de un usuario con filtros
func (s *reportService) GetUserReports(ctx context.Context, filters *domain.ReportFilters) (*domain.ReportsResponse, error) {
	if filters.Limit == 0 {
		filters.Limit = 20
	}

	reports, total, err := s.repo.GetByUser(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reports: %w", err)
	}

	// Construir respuestas con relaciones
	reportResponses := make([]domain.ReportResponse, 0, len(reports))
	for _, report := range reports {
		items, _ := s.repo.GetReportExpenses(ctx, report.ID)
		approvals, _ := s.repo.GetApprovalsByReport(ctx, report.ID)
		comments, _ := s.repo.GetCommentsByReport(ctx, report.ID, true)

		reportResponses = append(reportResponses, domain.ReportResponse{
			ExpenseReport: report,
			Items:         items,
			Approvals:     approvals,
			Comments:      comments,
		})
	}

	return &domain.ReportsResponse{
		Reports: reportResponses,
		Total:   total,
		Limit:   filters.Limit,
		Offset:  filters.Offset,
	}, nil
}

// UpdateReport actualiza un reporte existente
func (s *reportService) UpdateReport(ctx context.Context, reportID uuid.UUID, dto *domain.UpdateReportDto) error {
	report, err := s.repo.GetByID(ctx, reportID)
	if err != nil {
		return fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return ErrReportNotFound
	}

	// Validar que el reporte sea editable
	if report.Status != domain.ReportStatusDraft {
		return ErrReportNotEditable
	}

	// Aplicar cambios
	if dto.Title != nil {
		report.Title = *dto.Title
	}
	if dto.Description != nil {
		report.Description = dto.Description
	}

	report.Updated = time.Now()

	if err := s.repo.Update(ctx, report); err != nil {
		return fmt.Errorf("failed to update report: %w", err)
	}

	return nil
}

// DeleteReport elimina un reporte (soft delete)
func (s *reportService) DeleteReport(ctx context.Context, reportID uuid.UUID) error {
	report, err := s.repo.GetByID(ctx, reportID)
	if err != nil {
		return fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return ErrReportNotFound
	}

	// Solo se puede eliminar si está en draft
	if report.Status != domain.ReportStatusDraft {
		return ErrReportNotDeletable
	}

	if err := s.repo.Delete(ctx, reportID); err != nil {
		return fmt.Errorf("failed to delete report: %w", err)
	}

	return nil
}

// AddExpensesToReport agrega gastos a un reporte
func (s *reportService) AddExpensesToReport(ctx context.Context, reportID uuid.UUID, dto *domain.AddExpensesDto) error {
	report, err := s.repo.GetByID(ctx, reportID)
	if err != nil {
		return fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return ErrReportNotFound
	}

	// Solo se pueden agregar gastos en estado draft
	if report.Status != domain.ReportStatusDraft {
		return ErrReportNotEditable
	}

	// Verificar que los gastos no estén en otro reporte
	for _, expenseID := range dto.ExpenseIDs {
		inReport, reportID, err := s.repo.IsExpenseInReport(ctx, expenseID)
		if err != nil {
			return fmt.Errorf("failed to check expense: %w", err)
		}
		if inReport {
			return fmt.Errorf("%w (expense %s in report %s)", ErrExpenseAlreadyInReport, expenseID, *reportID)
		}
	}

	// Agregar gastos
	for _, expenseID := range dto.ExpenseIDs {
		if err := s.repo.AddExpenseToReport(ctx, reportID, expenseID); err != nil {
			return fmt.Errorf("failed to add expense to report: %w", err)
		}
	}

	// Recalcular total
	newTotal, err := s.repo.RecalculateTotal(ctx, reportID)
	if err != nil {
		return fmt.Errorf("failed to recalculate total: %w", err)
	}

	// Actualizar total en el reporte
	report.TotalAmount = newTotal
	report.Updated = time.Now()
	if err := s.repo.Update(ctx, report); err != nil {
		return fmt.Errorf("failed to update report total: %w", err)
	}

	return nil
}

// RemoveExpenseFromReport elimina un gasto de un reporte
func (s *reportService) RemoveExpenseFromReport(ctx context.Context, reportID, expenseID uuid.UUID) error {
	report, err := s.repo.GetByID(ctx, reportID)
	if err != nil {
		return fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return ErrReportNotFound
	}

	// Solo se pueden remover gastos en estado draft
	if report.Status != domain.ReportStatusDraft {
		return ErrReportNotEditable
	}

	// Verificar que el gasto esté en este reporte
	inReport, inReportID, err := s.repo.IsExpenseInReport(ctx, expenseID)
	if err != nil {
		return fmt.Errorf("failed to check expense: %w", err)
	}
	if !inReport || *inReportID != reportID {
		return ErrExpenseNotInReport
	}

	// Remover gasto
	if err := s.repo.RemoveExpenseFromReport(ctx, reportID, expenseID); err != nil {
		return fmt.Errorf("failed to remove expense from report: %w", err)
	}

	// Recalcular total
	newTotal, err := s.repo.RecalculateTotal(ctx, reportID)
	if err != nil {
		return fmt.Errorf("failed to recalculate total: %w", err)
	}

	report.TotalAmount = newTotal
	report.Updated = time.Now()
	if err := s.repo.Update(ctx, report); err != nil {
		return fmt.Errorf("failed to update report total: %w", err)
	}

	return nil
}

// SubmitReport envía un reporte para aprobación
func (s *reportService) SubmitReport(ctx context.Context, reportID uuid.UUID) (*domain.SubmitReportResponse, error) {
	report, err := s.repo.GetByID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return nil, ErrReportNotFound
	}

	// Validar que esté en draft
	if report.Status != domain.ReportStatusDraft {
		return nil, ErrReportNotSubmittable
	}

	// Validar que tenga al menos un gasto
	items, err := s.repo.GetReportExpenses(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report expenses: %w", err)
	}
	if len(items) == 0 {
		return nil, ErrReportEmpty
	}

	// Cambiar estado a submitted
	now := time.Now()
	report.Status = domain.ReportStatusSubmitted
	report.SubmissionDate = &now
	report.Updated = now

	if err := s.repo.Update(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to update report status: %w", err)
	}

	// TODO: Aquí se debería consultar la política y crear las aprobaciones necesarias
	// Por ahora, se retorna vacío - esto se implementará cuando se integre con PolicyService
	approvalsCreated := []domain.Approval{}

	return &domain.SubmitReportResponse{
		Report:           *report,
		ApprovalsCreated: approvalsCreated,
	}, nil
}

// ApproveReport aprueba un reporte
func (s *reportService) ApproveReport(ctx context.Context, approvalID, approverID uuid.UUID, dto *domain.ApproveReportDto) error {
	approval, err := s.repo.GetApprovalByID(ctx, approvalID)
	if err != nil {
		return fmt.Errorf("failed to get approval: %w", err)
	}
	if approval == nil {
		return ErrApprovalNotFound
	}

	// Verificar que el usuario sea el aprobador asignado
	if approval.ApproverID != approverID {
		return ErrUnauthorized
	}

	// Verificar que esté en estado pending
	if approval.Status != domain.ApprovalStatusPending {
		return ErrInvalidApprovalStatus
	}

	// Actualizar aprobación
	now := time.Now()
	approval.Status = domain.ApprovalStatusApproved
	approval.DecisionDate = &now
	approval.Updated = now

	if dto.Comments != "" {
		approval.Comments = &dto.Comments
	}
	if dto.ApprovedAmount != nil {
		approval.ApprovedAmount = dto.ApprovedAmount
	}

	if err := s.repo.UpdateApproval(ctx, approval); err != nil {
		return fmt.Errorf("failed to update approval: %w", err)
	}

	// Registrar en historial
	prevStatus := domain.ApprovalStatusPending
	newStatus := domain.ApprovalStatusApproved
	history := &domain.ApprovalHistory{
		ID:             uuid.New(),
		ApprovalID:     approvalID,
		ReportID:       approval.ReportID,
		ActorID:        approverID,
		Action:         domain.ApprovalActionApproved,
		PreviousStatus: &prevStatus,
		NewStatus:      &newStatus,
		Comments:       approval.Comments,
		Created:        now,
	}

	if err := s.repo.CreateApprovalHistory(ctx, history); err != nil {
		return fmt.Errorf("failed to create approval history: %w", err)
	}

	// TODO: Verificar si todas las aprobaciones están completas y actualizar estado del reporte

	return nil
}

// RejectReport rechaza un reporte
func (s *reportService) RejectReport(ctx context.Context, approvalID, approverID uuid.UUID, dto *domain.RejectReportDto) error {
	approval, err := s.repo.GetApprovalByID(ctx, approvalID)
	if err != nil {
		return fmt.Errorf("failed to get approval: %w", err)
	}
	if approval == nil {
		return ErrApprovalNotFound
	}

	// Verificar que el usuario sea el aprobador asignado
	if approval.ApproverID != approverID {
		return ErrUnauthorized
	}

	// Verificar que esté en estado pending
	if approval.Status != domain.ApprovalStatusPending {
		return ErrInvalidApprovalStatus
	}

	// Actualizar aprobación
	now := time.Now()
	approval.Status = domain.ApprovalStatusRejected
	approval.DecisionDate = &now
	approval.Comments = &dto.Reason
	approval.Updated = now

	if err := s.repo.UpdateApproval(ctx, approval); err != nil {
		return fmt.Errorf("failed to update approval: %w", err)
	}

	// Registrar en historial
	prevStatusReject := domain.ApprovalStatusPending
	newStatusReject := domain.ApprovalStatusRejected
	history := &domain.ApprovalHistory{
		ID:             uuid.New(),
		ApprovalID:     approvalID,
		ReportID:       approval.ReportID,
		ActorID:        approverID,
		Action:         domain.ApprovalActionRejected,
		PreviousStatus: &prevStatusReject,
		NewStatus:      &newStatusReject,
		Comments:       &dto.Reason,
		Created:        now,
	}

	if err := s.repo.CreateApprovalHistory(ctx, history); err != nil {
		return fmt.Errorf("failed to create approval history: %w", err)
	}

	// Actualizar estado del reporte a rejected
	if err := s.repo.UpdateStatus(ctx, approval.ReportID, domain.ReportStatusRejected); err != nil {
		return fmt.Errorf("failed to update report status: %w", err)
	}

	return nil
}

// AddComment agrega un comentario a un reporte o gasto
func (s *reportService) AddComment(ctx context.Context, userID uuid.UUID, dto *domain.CreateCommentDto) (*domain.ExpenseComment, error) {
	// Validar que al menos uno esté presente
	if dto.ReportID == nil && dto.ExpenseID == nil {
		return nil, errors.New("report_id or expense_id is required")
	}

	now := time.Now()
	comment := &domain.ExpenseComment{
		ID:          uuid.New(),
		ReportID:    dto.ReportID,
		ExpenseID:   dto.ExpenseID,
		UserID:      userID,
		CommentType: dto.CommentType,
		Content:     dto.Content,
		ParentID:    dto.ParentID,
		IsInternal:  dto.IsInternal,
		Created:     now,
		Updated:     now,
	}

	if comment.CommentType == "" {
		comment.CommentType = domain.CommentTypeGeneral
	}

	if err := s.repo.CreateComment(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

// GetReportComments obtiene los comentarios de un reporte
func (s *reportService) GetReportComments(ctx context.Context, reportID uuid.UUID, includeInternal bool) ([]domain.ExpenseComment, error) {
	comments, err := s.repo.GetCommentsByReport(ctx, reportID, includeInternal)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	return comments, nil
}

// CanEditReport verifica si un reporte puede ser editado por un usuario
func (s *reportService) CanEditReport(ctx context.Context, reportID, userID uuid.UUID) (bool, error) {
	report, err := s.repo.GetByID(ctx, reportID)
	if err != nil {
		return false, fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return false, ErrReportNotFound
	}

	// Solo el dueño puede editar y solo en estado draft
	if report.UserID != userID {
		return false, nil
	}

	return report.Status == domain.ReportStatusDraft, nil
}

// CanDeleteReport verifica si un reporte puede ser eliminado por un usuario
func (s *reportService) CanDeleteReport(ctx context.Context, reportID, userID uuid.UUID) (bool, error) {
	report, err := s.repo.GetByID(ctx, reportID)
	if err != nil {
		return false, fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return false, ErrReportNotFound
	}

	// Solo el dueño puede eliminar y solo en estado draft
	if report.UserID != userID {
		return false, nil
	}

	return report.Status == domain.ReportStatusDraft, nil
}

// CanSubmitReport verifica si un reporte puede ser enviado para aprobación
func (s *reportService) CanSubmitReport(ctx context.Context, reportID uuid.UUID) (bool, error) {
	report, err := s.repo.GetByID(ctx, reportID)
	if err != nil {
		return false, fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return false, ErrReportNotFound
	}

	// Debe estar en draft
	if report.Status != domain.ReportStatusDraft {
		return false, nil
	}

	// Debe tener al menos un gasto
	items, err := s.repo.GetReportExpenses(ctx, reportID)
	if err != nil {
		return false, fmt.Errorf("failed to get report expenses: %w", err)
	}

	return len(items) > 0, nil
}

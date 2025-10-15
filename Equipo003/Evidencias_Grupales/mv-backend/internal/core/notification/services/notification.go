package services

import (
	"context"
	"fmt"

	"github.com/JoseLuis21/mv-backend/internal/core/notification/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/notification/ports"
	expenseDomain "github.com/JoseLuis21/mv-backend/internal/core/expense/domain"
	reportDomain "github.com/JoseLuis21/mv-backend/internal/core/report/domain"
	"github.com/google/uuid"
)

type notificationService struct {
	repo ports.NotificationRepository
}

// NewNotificationService crea una nueva instancia del servicio de notificaciones
func NewNotificationService(repo ports.NotificationRepository) ports.NotificationService {
	return &notificationService{repo: repo}
}

// Create crea una nueva notificación
func (s *notificationService) Create(ctx context.Context, dto *domain.CreateNotificationDto) (*domain.Notification, error) {
	notification := domain.NewNotification(
		dto.UserID,
		dto.Type,
		dto.Title,
		dto.Message,
		dto.Data,
		dto.RelatedEntityID,
		dto.RelatedEntityType,
	)

	if err := s.repo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

// GetByID obtiene una notificación por su ID
func (s *notificationService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	notification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	return notification, nil
}

// GetByUser obtiene todas las notificaciones de un usuario
func (s *notificationService) GetByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]domain.Notification, error) {
	notifications, err := s.repo.GetByUser(ctx, userID, unreadOnly, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}
	return notifications, nil
}

// GetByUserPaginated obtiene las notificaciones de un usuario con paginación
func (s *notificationService) GetByUserPaginated(ctx context.Context, userID uuid.UUID, unreadOnly bool, limit, offset int) ([]domain.Notification, int64, error) {
	notifications, err := s.repo.GetByUser(ctx, userID, unreadOnly, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user notifications: %w", err)
	}

	count, err := s.repo.CountByUser(ctx, userID, unreadOnly)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user notifications: %w", err)
	}

	return notifications, count, nil
}

// MarkAsRead marca una notificación como leída
func (s *notificationService) MarkAsRead(ctx context.Context, notificationID uuid.UUID) error {
	if err := s.repo.MarkAsRead(ctx, notificationID); err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}

// MarkAllAsRead marca todas las notificaciones de un usuario como leídas
func (s *notificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.MarkAllAsRead(ctx, userID); err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	return nil
}

// Delete elimina una notificación (soft delete)
func (s *notificationService) Delete(ctx context.Context, notificationID uuid.UUID) error {
	if err := s.repo.Delete(ctx, notificationID); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

// NotifyExpenseApproved crea una notificación cuando un gasto es aprobado
func (s *notificationService) NotifyExpenseApproved(ctx context.Context, expense *expenseDomain.Expense) error {
	entityType := "expense"
	dto := &domain.CreateNotificationDto{
		UserID:            expense.UserID,
		Type:              domain.NotificationTypeExpenseApproved,
		Title:             "Gasto aprobado",
		Message:           fmt.Sprintf("Tu gasto '%s' por $%s ha sido aprobado.", expense.Title, formatAmount(expense.Amount, expense.Currency)),
		Data:              map[string]any{"expense_id": expense.ID, "amount": expense.Amount, "currency": expense.Currency},
		RelatedEntityID:   &expense.ID,
		RelatedEntityType: &entityType,
	}

	_, err := s.Create(ctx, dto)
	return err
}

// NotifyExpenseRejected crea una notificación cuando un gasto es rechazado
func (s *notificationService) NotifyExpenseRejected(ctx context.Context, expense *expenseDomain.Expense, reason string) error {
	entityType := "expense"
	dto := &domain.CreateNotificationDto{
		UserID:            expense.UserID,
		Type:              domain.NotificationTypeExpenseRejected,
		Title:             "Gasto rechazado",
		Message:           fmt.Sprintf("Tu gasto '%s' fue rechazado. Razón: %s", expense.Title, reason),
		Data:              map[string]any{"expense_id": expense.ID, "reason": reason},
		RelatedEntityID:   &expense.ID,
		RelatedEntityType: &entityType,
	}

	_, err := s.Create(ctx, dto)
	return err
}

// NotifyExpenseSubmitted crea una notificación cuando un gasto es enviado
func (s *notificationService) NotifyExpenseSubmitted(ctx context.Context, expense *expenseDomain.Expense) error {
	entityType := "expense"
	dto := &domain.CreateNotificationDto{
		UserID:            expense.UserID,
		Type:              domain.NotificationTypeExpenseSubmitted,
		Title:             "Gasto enviado",
		Message:           fmt.Sprintf("Tu gasto '%s' ha sido enviado para revisión.", expense.Title),
		Data:              map[string]any{"expense_id": expense.ID},
		RelatedEntityID:   &expense.ID,
		RelatedEntityType: &entityType,
	}

	_, err := s.Create(ctx, dto)
	return err
}

// NotifyApprovalNeeded crea una notificación cuando se requiere una aprobación
func (s *notificationService) NotifyApprovalNeeded(ctx context.Context, approval *reportDomain.Approval, report *reportDomain.ExpenseReport) error {
	entityType := "approval"
	dto := &domain.CreateNotificationDto{
		UserID:            approval.ApproverID,
		Type:              domain.NotificationTypeApprovalNeeded,
		Title:             "Nueva aprobación pendiente",
		Message:           fmt.Sprintf("Tienes una nueva solicitud de aprobación para el reporte '%s' (Nivel %d).", report.Title, approval.Level),
		Data:              map[string]any{"approval_id": approval.ID, "report_id": report.ID, "level": approval.Level, "amount": report.TotalAmount},
		RelatedEntityID:   &approval.ID,
		RelatedEntityType: &entityType,
	}

	_, err := s.Create(ctx, dto)
	return err
}

// NotifyApprovalApproved crea una notificación cuando una aprobación es aprobada
func (s *notificationService) NotifyApprovalApproved(ctx context.Context, approval *reportDomain.Approval, report *reportDomain.ExpenseReport) error {
	entityType := "approval"
	dto := &domain.CreateNotificationDto{
		UserID:            report.UserID,
		Type:              domain.NotificationTypeApprovalApproved,
		Title:             "Aprobación confirmada",
		Message:           fmt.Sprintf("Tu reporte '%s' ha sido aprobado en el nivel %d.", report.Title, approval.Level),
		Data:              map[string]any{"approval_id": approval.ID, "report_id": report.ID, "level": approval.Level},
		RelatedEntityID:   &approval.ID,
		RelatedEntityType: &entityType,
	}

	_, err := s.Create(ctx, dto)
	return err
}

// NotifyApprovalRejected crea una notificación cuando una aprobación es rechazada
func (s *notificationService) NotifyApprovalRejected(ctx context.Context, approval *reportDomain.Approval, report *reportDomain.ExpenseReport, reason string) error {
	entityType := "approval"
	dto := &domain.CreateNotificationDto{
		UserID:            report.UserID,
		Type:              domain.NotificationTypeApprovalRejected,
		Title:             "Aprobación rechazada",
		Message:           fmt.Sprintf("Tu reporte '%s' fue rechazado en el nivel %d. Razón: %s", report.Title, approval.Level, reason),
		Data:              map[string]any{"approval_id": approval.ID, "report_id": report.ID, "level": approval.Level, "reason": reason},
		RelatedEntityID:   &approval.ID,
		RelatedEntityType: &entityType,
	}

	_, err := s.Create(ctx, dto)
	return err
}

// NotifyApprovalEscalated crea una notificación cuando una aprobación es escalada
func (s *notificationService) NotifyApprovalEscalated(ctx context.Context, approval *reportDomain.Approval, report *reportDomain.ExpenseReport, newApproverID uuid.UUID) error {
	entityType := "approval"

	// Notificar al nuevo aprobador
	dtoNewApprover := &domain.CreateNotificationDto{
		UserID:            newApproverID,
		Type:              domain.NotificationTypeApprovalEscalated,
		Title:             "Aprobación escalada a ti",
		Message:           fmt.Sprintf("Se te ha asignado la aprobación del reporte '%s' (Nivel %d) por escalamiento.", report.Title, approval.Level),
		Data:              map[string]any{"approval_id": approval.ID, "report_id": report.ID, "level": approval.Level},
		RelatedEntityID:   &approval.ID,
		RelatedEntityType: &entityType,
	}

	if _, err := s.Create(ctx, dtoNewApprover); err != nil {
		return err
	}

	// Notificar al creador del reporte
	dtoReportOwner := &domain.CreateNotificationDto{
		UserID:            report.UserID,
		Type:              domain.NotificationTypeApprovalEscalated,
		Title:             "Aprobación escalada",
		Message:           fmt.Sprintf("La aprobación de tu reporte '%s' (Nivel %d) ha sido escalada a otro aprobador.", report.Title, approval.Level),
		Data:              map[string]any{"approval_id": approval.ID, "report_id": report.ID, "level": approval.Level},
		RelatedEntityID:   &approval.ID,
		RelatedEntityType: &entityType,
	}

	_, err := s.Create(ctx, dtoReportOwner)
	return err
}

// NotifyReportSubmitted crea una notificación cuando un reporte es enviado
func (s *notificationService) NotifyReportSubmitted(ctx context.Context, report *reportDomain.ExpenseReport) error {
	entityType := "report"
	dto := &domain.CreateNotificationDto{
		UserID:            report.UserID,
		Type:              domain.NotificationTypeReportSubmitted,
		Title:             "Reporte enviado",
		Message:           fmt.Sprintf("Tu reporte '%s' ha sido enviado para aprobación.", report.Title),
		Data:              map[string]any{"report_id": report.ID, "amount": report.TotalAmount},
		RelatedEntityID:   &report.ID,
		RelatedEntityType: &entityType,
	}

	_, err := s.Create(ctx, dto)
	return err
}

// NotifyReportApproved crea una notificación cuando un reporte es completamente aprobado
func (s *notificationService) NotifyReportApproved(ctx context.Context, report *reportDomain.ExpenseReport) error {
	entityType := "report"
	dto := &domain.CreateNotificationDto{
		UserID:            report.UserID,
		Type:              domain.NotificationTypeReportApproved,
		Title:             "Reporte aprobado",
		Message:           fmt.Sprintf("¡Tu reporte '%s' ha sido completamente aprobado! Monto total: $%s", report.Title, formatAmount(report.TotalAmount, report.Currency)),
		Data:              map[string]any{"report_id": report.ID, "amount": report.TotalAmount, "currency": report.Currency},
		RelatedEntityID:   &report.ID,
		RelatedEntityType: &entityType,
	}

	_, err := s.Create(ctx, dto)
	return err
}

// NotifyReportRejected crea una notificación cuando un reporte es rechazado
func (s *notificationService) NotifyReportRejected(ctx context.Context, report *reportDomain.ExpenseReport, reason string) error {
	entityType := "report"
	dto := &domain.CreateNotificationDto{
		UserID:            report.UserID,
		Type:              domain.NotificationTypeReportRejected,
		Title:             "Reporte rechazado",
		Message:           fmt.Sprintf("Tu reporte '%s' ha sido rechazado. Razón: %s", report.Title, reason),
		Data:              map[string]any{"report_id": report.ID, "reason": reason},
		RelatedEntityID:   &report.ID,
		RelatedEntityType: &entityType,
	}

	_, err := s.Create(ctx, dto)
	return err
}

// NotifyCommentAdded crea una notificación cuando se añade un comentario
func (s *notificationService) NotifyCommentAdded(ctx context.Context, comment *reportDomain.ExpenseComment, report *reportDomain.ExpenseReport) error {
	// No notificar al autor del comentario
	if comment.UserID == report.UserID {
		return nil
	}

	entityType := "comment"
	dto := &domain.CreateNotificationDto{
		UserID:            report.UserID,
		Type:              domain.NotificationTypeCommentAdded,
		Title:             "Nuevo comentario",
		Message:           fmt.Sprintf("Se ha añadido un comentario a tu reporte '%s'.", report.Title),
		Data:              map[string]any{"comment_id": comment.ID, "report_id": report.ID},
		RelatedEntityID:   &comment.ID,
		RelatedEntityType: &entityType,
	}

	_, err := s.Create(ctx, dto)
	return err
}

// formatAmount formatea un monto con su moneda
func formatAmount(amount float64, currency string) string {
	return fmt.Sprintf("%.2f %s", amount, currency)
}

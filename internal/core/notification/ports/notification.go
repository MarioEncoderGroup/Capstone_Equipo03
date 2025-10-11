package ports

import (
	"context"

	"github.com/JoseLuis21/mv-backend/internal/core/notification/domain"
	expenseDomain "github.com/JoseLuis21/mv-backend/internal/core/expense/domain"
	reportDomain "github.com/JoseLuis21/mv-backend/internal/core/report/domain"
	"github.com/google/uuid"
)

// NotificationRepository define la interfaz para acceso a datos de notificaciones
type NotificationRepository interface {
	// CRUD operations
	Create(ctx context.Context, notification *domain.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	GetByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool, limit, offset int) ([]domain.Notification, error)
	CountByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) (int64, error)
	Update(ctx context.Context, notification *domain.Notification) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Bulk operations
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	DeleteByUser(ctx context.Context, userID uuid.UUID) error
}

// NotificationService define la interfaz para lógica de negocio de notificaciones
type NotificationService interface {
	// CRUD operations
	Create(ctx context.Context, dto *domain.CreateNotificationDto) (*domain.Notification, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	GetByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]domain.Notification, error)
	GetByUserPaginated(ctx context.Context, userID uuid.UUID, unreadOnly bool, limit, offset int) ([]domain.Notification, int64, error)
	MarkAsRead(ctx context.Context, notificationID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, notificationID uuid.UUID) error

	// Helper methods for specific notification types
	NotifyExpenseApproved(ctx context.Context, expense *expenseDomain.Expense) error
	NotifyExpenseRejected(ctx context.Context, expense *expenseDomain.Expense, reason string) error
	NotifyExpenseSubmitted(ctx context.Context, expense *expenseDomain.Expense) error

	NotifyApprovalNeeded(ctx context.Context, approval *reportDomain.Approval, report *reportDomain.ExpenseReport) error
	NotifyApprovalApproved(ctx context.Context, approval *reportDomain.Approval, report *reportDomain.ExpenseReport) error
	NotifyApprovalRejected(ctx context.Context, approval *reportDomain.Approval, report *reportDomain.ExpenseReport, reason string) error
	NotifyApprovalEscalated(ctx context.Context, approval *reportDomain.Approval, report *reportDomain.ExpenseReport, newApproverID uuid.UUID) error

	NotifyReportSubmitted(ctx context.Context, report *reportDomain.ExpenseReport) error
	NotifyReportApproved(ctx context.Context, report *reportDomain.ExpenseReport) error
	NotifyReportRejected(ctx context.Context, report *reportDomain.ExpenseReport, reason string) error

	NotifyCommentAdded(ctx context.Context, comment *reportDomain.ExpenseComment, report *reportDomain.ExpenseReport) error
}

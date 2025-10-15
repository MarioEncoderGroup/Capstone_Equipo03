package domain

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType representa los tipos de notificaciones del sistema
type NotificationType string

const (
	NotificationTypeExpenseApproved   NotificationType = "expense_approved"
	NotificationTypeExpenseRejected   NotificationType = "expense_rejected"
	NotificationTypeExpenseSubmitted  NotificationType = "expense_submitted"
	NotificationTypeApprovalNeeded    NotificationType = "approval_needed"
	NotificationTypeApprovalApproved  NotificationType = "approval_approved"
	NotificationTypeApprovalRejected  NotificationType = "approval_rejected"
	NotificationTypeApprovalEscalated NotificationType = "approval_escalated"
	NotificationTypeReportSubmitted   NotificationType = "report_submitted"
	NotificationTypeReportApproved    NotificationType = "report_approved"
	NotificationTypeReportRejected    NotificationType = "report_rejected"
	NotificationTypeCommentAdded      NotificationType = "comment_added"
	NotificationTypeSystem            NotificationType = "system_notification"
)

// Notification representa una notificación para un usuario
type Notification struct {
	ID                 uuid.UUID        `json:"id"`
	UserID             uuid.UUID        `json:"user_id"`
	Type               NotificationType `json:"type"`
	Title              string           `json:"title"`
	Message            string           `json:"message"`
	Data               map[string]any   `json:"data,omitempty"`
	RelatedEntityID    *uuid.UUID       `json:"related_entity_id,omitempty"`
	RelatedEntityType  *string          `json:"related_entity_type,omitempty"`
	Read               bool             `json:"read"`
	ReadAt             *time.Time       `json:"read_at,omitempty"`
	Created            time.Time        `json:"created"`
	DeletedAt          *time.Time       `json:"deleted_at,omitempty"`
}

// NewNotification crea una nueva notificación
func NewNotification(
	userID uuid.UUID,
	notifType NotificationType,
	title, message string,
	data map[string]any,
	relatedEntityID *uuid.UUID,
	relatedEntityType *string,
) *Notification {
	now := time.Now()
	return &Notification{
		ID:                 uuid.New(),
		UserID:             userID,
		Type:               notifType,
		Title:              title,
		Message:            message,
		Data:               data,
		RelatedEntityID:    relatedEntityID,
		RelatedEntityType:  relatedEntityType,
		Read:               false,
		ReadAt:             nil,
		Created:            now,
		DeletedAt:          nil,
	}
}

// MarkAsRead marca la notificación como leída
func (n *Notification) MarkAsRead() {
	if !n.Read {
		now := time.Now()
		n.Read = true
		n.ReadAt = &now
	}
}

// IsRead devuelve si la notificación ha sido leída
func (n *Notification) IsRead() bool {
	return n.Read
}

// SoftDelete realiza un soft delete de la notificación
func (n *Notification) SoftDelete() {
	if n.DeletedAt == nil {
		now := time.Now()
		n.DeletedAt = &now
	}
}

// IsDeleted devuelve si la notificación ha sido eliminada
func (n *Notification) IsDeleted() bool {
	return n.DeletedAt != nil
}

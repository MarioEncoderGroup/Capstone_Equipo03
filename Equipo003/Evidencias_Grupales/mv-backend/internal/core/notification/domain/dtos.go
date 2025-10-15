package domain

import "github.com/google/uuid"

// CreateNotificationDto contiene los datos para crear una notificación
type CreateNotificationDto struct {
	UserID            uuid.UUID        `json:"user_id" validate:"required,uuid"`
	Type              NotificationType `json:"type" validate:"required"`
	Title             string           `json:"title" validate:"required,min=3,max=200"`
	Message           string           `json:"message" validate:"required,min=10,max=1000"`
	Data              map[string]any   `json:"data,omitempty"`
	RelatedEntityID   *uuid.UUID       `json:"related_entity_id,omitempty" validate:"omitempty,uuid"`
	RelatedEntityType *string          `json:"related_entity_type,omitempty" validate:"omitempty,min=3,max=50"`
}

// GetNotificationsDto contiene los filtros para obtener notificaciones
type GetNotificationsDto struct {
	UserID     uuid.UUID `json:"user_id" validate:"required,uuid"`
	UnreadOnly bool      `json:"unread_only"`
	Limit      int       `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset     int       `json:"offset" validate:"omitempty,min=0"`
}

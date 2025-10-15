package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/core/notification/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/notification/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type notificationRepository struct {
	db *postgresql.PostgresqlClient
}

// NewNotificationRepository crea una nueva instancia del repositorio de notificaciones
func NewNotificationRepository(db *postgresql.PostgresqlClient) ports.NotificationRepository {
	return &notificationRepository{db: db}
}

// Create crea una nueva notificación en la base de datos
func (r *notificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	query := `
		INSERT INTO notifications (
			id, user_id, type, title, message, data,
			related_entity_id, related_entity_type,
			read, read_at, created
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	dataJSON, err := json.Marshal(notification.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal notification data: %w", err)
	}

	_, err = r.db.Pool.Exec(ctx, query,
		notification.ID,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		dataJSON,
		notification.RelatedEntityID,
		notification.RelatedEntityType,
		notification.Read,
		notification.ReadAt,
		notification.Created,
	)

	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}

// GetByID obtiene una notificación por su ID
func (r *notificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, data,
			   related_entity_id, related_entity_type,
			   read, read_at, created, deleted_at
		FROM notifications
		WHERE id = $1 AND deleted_at IS NULL
	`

	var notification domain.Notification
	var dataJSON []byte

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Type,
		&notification.Title,
		&notification.Message,
		&dataJSON,
		&notification.RelatedEntityID,
		&notification.RelatedEntityType,
		&notification.Read,
		&notification.ReadAt,
		&notification.Created,
		&notification.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("notification not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	if len(dataJSON) > 0 {
		if err := json.Unmarshal(dataJSON, &notification.Data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal notification data: %w", err)
		}
	}

	return &notification, nil
}

// GetByUser obtiene las notificaciones de un usuario
func (r *notificationRepository) GetByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool, limit, offset int) ([]domain.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, data,
			   related_entity_id, related_entity_type,
			   read, read_at, created, deleted_at
		FROM notifications
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{userID}

	if unreadOnly {
		query += " AND read = false"
	}

	query += " ORDER BY created DESC LIMIT $2 OFFSET $3"
	args = append(args, limit, offset)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var notification domain.Notification
		var dataJSON []byte

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Type,
			&notification.Title,
			&notification.Message,
			&dataJSON,
			&notification.RelatedEntityID,
			&notification.RelatedEntityType,
			&notification.Read,
			&notification.ReadAt,
			&notification.Created,
			&notification.DeletedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}

		if len(dataJSON) > 0 {
			if err := json.Unmarshal(dataJSON, &notification.Data); err != nil {
				return nil, fmt.Errorf("failed to unmarshal notification data: %w", err)
			}
		}

		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notifications: %w", err)
	}

	return notifications, nil
}

// CountByUser cuenta las notificaciones de un usuario
func (r *notificationRepository) CountByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{userID}

	if unreadOnly {
		query += " AND read = false"
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	return count, nil
}

// Update actualiza una notificación
func (r *notificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	query := `
		UPDATE notifications
		SET type = $2,
			title = $3,
			message = $4,
			data = $5,
			related_entity_id = $6,
			related_entity_type = $7,
			read = $8,
			read_at = $9,
			deleted_at = $10
		WHERE id = $1 AND deleted_at IS NULL
	`

	dataJSON, err := json.Marshal(notification.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal notification data: %w", err)
	}

	result, err := r.db.Pool.Exec(ctx, query,
		notification.ID,
		notification.Type,
		notification.Title,
		notification.Message,
		dataJSON,
		notification.RelatedEntityID,
		notification.RelatedEntityType,
		notification.Read,
		notification.ReadAt,
		notification.DeletedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("notification not found or already deleted")
	}

	return nil
}

// Delete elimina una notificación (soft delete)
func (r *notificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE notifications
		SET deleted_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("notification not found or already deleted")
	}

	return nil
}

// MarkAsRead marca una notificación como leída
func (r *notificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE notifications
		SET read = true,
			read_at = $2
		WHERE id = $1 AND deleted_at IS NULL AND read = false
	`

	result, err := r.db.Pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("notification not found, already read, or deleted")
	}

	return nil
}

// MarkAllAsRead marca todas las notificaciones de un usuario como leídas
func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE notifications
		SET read = true,
			read_at = $2
		WHERE user_id = $1 AND deleted_at IS NULL AND read = false
	`

	_, err := r.db.Pool.Exec(ctx, query, userID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	return nil
}

// DeleteByUser elimina todas las notificaciones de un usuario (soft delete)
func (r *notificationRepository) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE notifications
		SET deleted_at = $2
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, userID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete user notifications: %w", err)
	}

	return nil
}

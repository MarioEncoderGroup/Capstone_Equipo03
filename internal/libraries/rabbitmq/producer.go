package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// NotificationMessage representa un mensaje de notificación en la cola
type NotificationMessage struct {
	NotificationID string         `json:"notification_id"`
	UserID         string         `json:"user_id"`
	Type           string         `json:"type"`
	Title          string         `json:"title"`
	Message        string         `json:"message"`
	Data           map[string]any `json:"data,omitempty"`
	Timestamp      time.Time      `json:"timestamp"`
	RetryCount     int            `json:"retry_count"`
}

// Producer gestiona la producción de mensajes a RabbitMQ
type Producer struct {
	client *Client
}

// NewProducer crea una nueva instancia del productor
func NewProducer(client *Client) *Producer {
	return &Producer{
		client: client,
	}
}

// PublishNotification publica una notificación en la cola
func (p *Producer) PublishNotification(ctx context.Context, msg *NotificationMessage) error {
	msg.Timestamp = time.Now()

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal notification message: %w", err)
	}

	// Publicar con retry automático
	return p.publishWithRetry(ctx, body, 3)
}

// publishWithRetry intenta publicar un mensaje con reintentos
func (p *Producer) publishWithRetry(ctx context.Context, body []byte, maxRetries int) error {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := p.client.channel.PublishWithContext(
			ctx,
			p.client.config.NotificationExchange, // exchange
			p.client.config.NotificationQueue,    // routing key
			false,                                 // mandatory
			false,                                 // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, // Persist message to disk
				ContentType:  "application/json",
				Body:         body,
				Timestamp:    time.Now(),
			},
		)

		if err == nil {
			return nil
		}

		lastErr = err

		if attempt < maxRetries-1 {
			time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
		}
	}

	return fmt.Errorf("failed to publish notification after %d attempts: %w", maxRetries, lastErr)
}

// PublishBulkNotifications publica múltiples notificaciones en batch
func (p *Producer) PublishBulkNotifications(ctx context.Context, messages []*NotificationMessage) error {
	for _, msg := range messages {
		if err := p.PublishNotification(ctx, msg); err != nil {
			return fmt.Errorf("failed to publish bulk notification: %w", err)
		}
	}
	return nil
}

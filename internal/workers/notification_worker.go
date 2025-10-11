package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/core/notification/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/notification/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/rabbitmq"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// MaxRetries define el número máximo de reintentos antes de enviar a DLQ
	MaxRetries = 3
	// RetryDelay define el tiempo de espera entre reintentos
	RetryDelay = 5 * time.Second
)

// NotificationWorker procesa mensajes de notificaciones desde RabbitMQ
type NotificationWorker struct {
	client              *rabbitmq.Client
	notificationService ports.NotificationService
	stopChan            chan bool
}

// NewNotificationWorker crea una nueva instancia del worker de notificaciones
func NewNotificationWorker(
	client *rabbitmq.Client,
	notificationService ports.NotificationService,
) *NotificationWorker {
	return &NotificationWorker{
		client:              client,
		notificationService: notificationService,
		stopChan:            make(chan bool),
	}
}

// Start inicia el worker para procesar mensajes
func (w *NotificationWorker) Start() error {
	log.Println("[NotificationWorker] Starting...")

	// Obtener el canal de RabbitMQ
	ch := w.client.GetChannel()

	// Comenzar a consumir mensajes
	msgs, err := ch.Consume(
		"notifications", // queue
		"",              // consumer tag
		false,           // auto-ack (disabled for manual ack)
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)

	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Println("[NotificationWorker] Started successfully. Waiting for messages...")

	// Procesar mensajes en un goroutine
	go func() {
		for {
			select {
			case <-w.stopChan:
				log.Println("[NotificationWorker] Stopping...")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("[NotificationWorker] Channel closed")
					return
				}
				w.processMessage(msg)
			}
		}
	}()

	return nil
}

// processMessage procesa un mensaje individual
func (w *NotificationWorker) processMessage(msg amqp.Delivery) {
	log.Printf("[NotificationWorker] Processing message: %s", msg.MessageId)

	ctx := context.Background()

	// Parsear el mensaje
	var notifMsg rabbitmq.NotificationMessage
	if err := json.Unmarshal(msg.Body, &notifMsg); err != nil {
		log.Printf("[NotificationWorker] Failed to unmarshal message: %v", err)
		// Rechazar el mensaje y no reencolar (mensaje malformado)
		msg.Nack(false, false)
		return
	}

	// Validar que el mensaje no exceda el límite de reintentos
	if notifMsg.RetryCount >= MaxRetries {
		log.Printf("[NotificationWorker] Message exceeded max retries (%d). Sending to DLQ.", MaxRetries)
		// Rechazar y no reencolar (irá a DLQ)
		msg.Nack(false, false)
		return
	}

	// Intentar procesar el mensaje
	if err := w.processNotification(ctx, &notifMsg); err != nil {
		log.Printf("[NotificationWorker] Failed to process notification: %v", err)

		// Incrementar contador de reintentos
		notifMsg.RetryCount++

		// Reencolar el mensaje con delay
		if notifMsg.RetryCount < MaxRetries {
			log.Printf("[NotificationWorker] Requeuing message (attempt %d/%d)", notifMsg.RetryCount, MaxRetries)

			// Rechazar y reencolar
			msg.Nack(false, true)

			// Esperar antes del siguiente reintento
			time.Sleep(RetryDelay)
		} else {
			log.Printf("[NotificationWorker] Max retries reached. Sending to DLQ.")
			// Rechazar y no reencolar (irá a DLQ)
			msg.Nack(false, false)
		}
		return
	}

	// Mensaje procesado exitosamente
	log.Printf("[NotificationWorker] Successfully processed notification: %s", notifMsg.NotificationID)
	msg.Ack(false)
}

// processNotification procesa la notificación y la guarda en la base de datos
func (w *NotificationWorker) processNotification(ctx context.Context, msg *rabbitmq.NotificationMessage) error {
	// Parsear UUIDs
	userID, err := uuid.Parse(msg.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Crear DTO para la notificación
	dto := &domain.CreateNotificationDto{
		UserID:  userID,
		Type:    domain.NotificationType(msg.Type),
		Title:   msg.Title,
		Message: msg.Message,
		Data:    msg.Data,
	}

	// Si el mensaje incluye notification_id, verificar que no exista ya
	if msg.NotificationID != "" {
		notifID, err := uuid.Parse(msg.NotificationID)
		if err == nil {
			// Verificar si ya existe
			existing, _ := w.notificationService.GetByID(ctx, notifID)
			if existing != nil {
				// Notificación ya existe, no duplicar
				log.Printf("[NotificationWorker] Notification %s already exists, skipping", msg.NotificationID)
				return nil
			}
		}
	}

	// Crear la notificación en la base de datos
	notification, err := w.notificationService.Create(ctx, dto)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	log.Printf("[NotificationWorker] Created notification %s for user %s", notification.ID, notification.UserID)

	return nil
}

// Stop detiene el worker
func (w *NotificationWorker) Stop() {
	log.Println("[NotificationWorker] Stop signal received")
	close(w.stopChan)
}

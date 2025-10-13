package rabbitmq

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Config contiene la configuración para conectar a RabbitMQ
type Config struct {
	URL                string
	ReconnectDelay     time.Duration
	ReconnectAttempts  int
	PrefetchCount      int
	NotificationQueue  string
	NotificationDLQ    string
	NotificationExchange string
}

// Client es el cliente de RabbitMQ
type Client struct {
	config     Config
	connection *amqp.Connection
	channel    *amqp.Channel
	done       chan bool
}

// NewClient crea una nueva instancia del cliente de RabbitMQ
func NewClient(config Config) (*Client, error) {
	// Set defaults
	if config.ReconnectDelay == 0 {
		config.ReconnectDelay = 5 * time.Second
	}
	if config.ReconnectAttempts == 0 {
		config.ReconnectAttempts = 10
	}
	if config.PrefetchCount == 0 {
		config.PrefetchCount = 10
	}
	if config.NotificationQueue == "" {
		config.NotificationQueue = "notifications"
	}
	if config.NotificationDLQ == "" {
		config.NotificationDLQ = "notifications.dlq"
	}
	if config.NotificationExchange == "" {
		config.NotificationExchange = "notifications.exchange"
	}

	client := &Client{
		config: config,
		done:   make(chan bool),
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
}

// connect establece la conexión con RabbitMQ
func (c *Client) connect() error {
	var err error

	log.Printf("[RabbitMQ] Connecting to %s", c.config.URL)

	c.connection, err = amqp.Dial(c.config.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	c.channel, err = c.connection.Channel()
	if err != nil {
		c.connection.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// Set QoS (prefetch count)
	if err := c.channel.Qos(c.config.PrefetchCount, 0, false); err != nil {
		c.channel.Close()
		c.connection.Close()
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Setup queues and exchanges
	if err := c.setupInfrastructure(); err != nil {
		c.channel.Close()
		c.connection.Close()
		return fmt.Errorf("failed to setup infrastructure: %w", err)
	}

	log.Println("[RabbitMQ] Connected successfully")

	// Monitor connection
	go c.monitorConnection()

	return nil
}

// setupInfrastructure declara exchanges, queues y bindings
func (c *Client) setupInfrastructure() error {
	// Declare Dead Letter Exchange
	err := c.channel.ExchangeDeclare(
		c.config.NotificationExchange+".dlx", // name
		"direct",                              // type
		true,                                  // durable
		false,                                 // auto-deleted
		false,                                 // internal
		false,                                 // no-wait
		nil,                                   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLX: %w", err)
	}

	// Declare Main Exchange
	err = c.channel.ExchangeDeclare(
		c.config.NotificationExchange, // name
		"direct",                       // type
		true,                           // durable
		false,                          // auto-deleted
		false,                          // internal
		false,                          // no-wait
		nil,                            // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare Dead Letter Queue
	_, err = c.channel.QueueDeclare(
		c.config.NotificationDLQ, // name
		true,                      // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}

	// Bind DLQ to DLX
	err = c.channel.QueueBind(
		c.config.NotificationDLQ,              // queue name
		c.config.NotificationQueue,            // routing key
		c.config.NotificationExchange+".dlx", // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind DLQ: %w", err)
	}

	// Declare Main Queue with DLQ configuration
	_, err = c.channel.QueueDeclare(
		c.config.NotificationQueue, // name
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		amqp.Table{
			"x-dead-letter-exchange":    c.config.NotificationExchange + ".dlx",
			"x-dead-letter-routing-key": c.config.NotificationQueue,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind Main Queue to Exchange
	err = c.channel.QueueBind(
		c.config.NotificationQueue,    // queue name
		c.config.NotificationQueue,    // routing key
		c.config.NotificationExchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	log.Printf("[RabbitMQ] Infrastructure setup complete (Queue: %s, DLQ: %s)",
		c.config.NotificationQueue, c.config.NotificationDLQ)

	return nil
}

// monitorConnection monitorea la conexión y reconecta si es necesario
func (c *Client) monitorConnection() {
	for {
		select {
		case <-c.done:
			return
		case err := <-c.connection.NotifyClose(make(chan *amqp.Error)):
			if err != nil {
				log.Printf("[RabbitMQ] Connection closed: %v", err)
				c.reconnect()
			}
		}
	}
}

// reconnect intenta reconectar a RabbitMQ
func (c *Client) reconnect() {
	for i := 0; i < c.config.ReconnectAttempts; i++ {
		log.Printf("[RabbitMQ] Reconnecting... (attempt %d/%d)", i+1, c.config.ReconnectAttempts)

		time.Sleep(c.config.ReconnectDelay)

		if err := c.connect(); err != nil {
			log.Printf("[RabbitMQ] Reconnect failed: %v", err)
			continue
		}

		log.Println("[RabbitMQ] Reconnected successfully")
		return
	}

	log.Fatalf("[RabbitMQ] Failed to reconnect after %d attempts", c.config.ReconnectAttempts)
}

// GetChannel devuelve el canal de RabbitMQ
func (c *Client) GetChannel() *amqp.Channel {
	return c.channel
}

// GetConnection devuelve la conexión de RabbitMQ
func (c *Client) GetConnection() *amqp.Connection {
	return c.connection
}

// Close cierra la conexión con RabbitMQ
func (c *Client) Close() error {
	close(c.done)

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return fmt.Errorf("failed to close channel: %w", err)
		}
	}

	if c.connection != nil {
		if err := c.connection.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}

	log.Println("[RabbitMQ] Connection closed")
	return nil
}

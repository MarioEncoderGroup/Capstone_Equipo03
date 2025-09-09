package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresqlClient represents a PostgreSQL database client for MisViaticos
type PostgresqlClient struct {
	Pool   *pgxpool.Pool
	config *Config
}

// Config holds the PostgreSQL connection configuration
type Config struct {
	User         string
	Password     string
	Host         string
	Port         string
	Database     string
	MaxConns     int32
	MinConns     int32
	MaxLifetime  time.Duration
	MaxIdleTime  time.Duration
	HealthCheck  bool
	SSLMode      string
}

// DefaultConfig returns a default PostgreSQL configuration for MisViaticos
func DefaultConfig() *Config {
	return &Config{
		User:        "postgres",
		Password:    "password123",
		Host:        "localhost",
		Port:        "5432",
		Database:    "misviaticos",
		MaxConns:    25,
		MinConns:    5,
		MaxLifetime: time.Hour,
		MaxIdleTime: time.Minute * 30,
		HealthCheck: true,
		SSLMode:     "disable",
	}
}

// NewPostgresqlClient creates a new PostgreSQL client with the provided configuration
func NewPostgresqlClient(config *Config) (*PostgresqlClient, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Build connection string
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.SSLMode,
	)

	// Configure pool settings
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL config: %w", err)
	}

	// Set connection pool limits
	poolConfig.MaxConns = config.MaxConns
	poolConfig.MinConns = config.MinConns
	poolConfig.MaxConnLifetime = config.MaxLifetime
	poolConfig.MaxConnIdleTime = config.MaxIdleTime

	// Create connection pool
	pool, err := pgxpool.New(context.Background(), poolConfig.ConnString())
	if err != nil {
		return nil, fmt.Errorf("failed to create PostgreSQL pool: %w", err)
	}

	client := &PostgresqlClient{
		Pool:   pool,
		config: config,
	}

	// Perform health check if enabled
	if config.HealthCheck {
		if err := client.Ping(context.Background()); err != nil {
			pool.Close()
			return nil, fmt.Errorf("PostgreSQL health check failed: %w", err)
		}
	}

	return client, nil
}

// Ping checks the database connection health
func (c *PostgresqlClient) Ping(ctx context.Context) error {
	return c.Pool.Ping(ctx)
}

// Close closes the database connection pool
func (c *PostgresqlClient) Close() {
	c.Pool.Close()
}

// GetStats returns connection pool statistics
func (c *PostgresqlClient) GetStats() *pgxpool.Stat {
	return c.Pool.Stat()
}

// Exec executes a query without returning any rows
func (c *PostgresqlClient) Exec(ctx context.Context, sql string, args ...interface{}) error {
	_, err := c.Pool.Exec(ctx, sql, args...)
	return err
}

// Query executes a query that returns rows
func (c *PostgresqlClient) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return c.Pool.Query(ctx, sql, args...)
}

// QueryRow executes a query that is expected to return at most one row
func (c *PostgresqlClient) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return c.Pool.QueryRow(ctx, sql, args...)
}
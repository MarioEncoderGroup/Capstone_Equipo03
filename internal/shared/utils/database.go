package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
)

// InitDatabaseControl initializes the control database connection for MisViaticos
func InitDatabaseControl() (*postgresql.PostgresqlClient, error) {
	config := &postgresql.Config{
		User:     getEnvOrDefault("POSTGRESQL_CONTROL_USER", "postgres"),
		Password: getEnvOrDefault("POSTGRESQL_CONTROL_PASSWORD", "password123"),
		Host:     getEnvOrDefault("POSTGRESQL_CONTROL_HOST", "localhost"),
		Port:     getEnvOrDefault("POSTGRESQL_CONTROL_PORT", "5432"),
		Database: getEnvOrDefault("POSTGRESQL_CONTROL_DATABASE", "misviaticos_control"),
		SSLMode:  getEnvOrDefault("POSTGRESQL_SSL_MODE", "disable"),
		MaxConns: 25,
		MinConns: 5,
		HealthCheck: true,
	}

	client, err := postgresql.NewPostgresqlClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize control database: %w", err)
	}

	return client, nil
}

// GetDatabaseTenantDefault gets the default tenant database connection
func GetDatabaseTenantDefault(ctx context.Context, tenantID string) (*postgresql.PostgresqlClient, error) {
	// For now, we use a simple naming convention for tenant databases
	// In production, you might query the control DB to get the actual tenant DB details
	tenantDBName := fmt.Sprintf("%s_%s", 
		getEnvOrDefault("POSTGRESQL_DATABASE_TENANT", "misviaticos_tenant"), 
		tenantID,
	)

	config := &postgresql.Config{
		User:     getEnvOrDefault("POSTGRESQL_NODE1_USER", "postgres"),
		Password: getEnvOrDefault("POSTGRESQL_NODE1_PASSWORD", "password123"),
		Host:     getEnvOrDefault("POSTGRESQL_NODE1_HOST", "localhost"),
		Port:     getEnvOrDefault("POSTGRESQL_NODE1_PORT", "5432"),
		Database: tenantDBName,
		SSLMode:  getEnvOrDefault("POSTGRESQL_SSL_MODE", "disable"),
		MaxConns: 20,
		MinConns: 3,
		HealthCheck: true,
	}

	client, err := postgresql.NewPostgresqlClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tenant database %s: %w", tenantDBName, err)
	}

	return client, nil
}

// GetDatabaseTenantByID gets a specific tenant database connection
func GetDatabaseTenantByID(ctx context.Context, tenantID string) (*postgresql.PostgresqlClient, error) {
	return GetDatabaseTenantDefault(ctx, tenantID)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
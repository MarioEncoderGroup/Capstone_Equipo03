package helpers

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/google/uuid"
)

// TestDatabaseConfig configuración para BD de testing
type TestDatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	SSLMode  string
}

// GetTestDatabaseConfig obtiene configuración de BD para testing
func GetTestDatabaseConfig() *TestDatabaseConfig {
	return &TestDatabaseConfig{
		Host:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		Port:     getEnvOrDefault("TEST_DB_PORT", "5432"),
		User:     getEnvOrDefault("TEST_DB_USER", "postgres"),
		Password: getEnvOrDefault("TEST_DB_PASSWORD", "password123"),
		SSLMode:  getEnvOrDefault("TEST_DB_SSL_MODE", "disable"),
	}
}

// CreateTestDatabase crea una base de datos temporal para testing
func CreateTestDatabase(t *testing.T) (*postgresql.PostgresqlClient, string, func()) {
	t.Helper()

	config := GetTestDatabaseConfig()

	// Generar nombre único para BD de testing
	testDBName := fmt.Sprintf("test_misviaticos_%s", uuid.New().String()[:8])

	// Conectar a PostgreSQL para crear BD
	adminClient, err := createAdminClient(config)
	if err != nil {
		t.Skipf("PostgreSQL not available for testing: %v", err)
		return nil, "", nil
	}

	// Crear BD de testing
	createQuery := fmt.Sprintf("CREATE DATABASE %s", testDBName)
	if err := adminClient.Exec(context.Background(), createQuery); err != nil {
		adminClient.Close()
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Conectar a la BD de testing
	testClient, err := createTestClient(config, testDBName)
	if err != nil {
		adminClient.Close()
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Función de limpieza
	cleanup := func() {
		testClient.Close()

		// Desconectar conexiones activas
		killConnectionsQuery := fmt.Sprintf(`
			SELECT pg_terminate_backend(pid) 
			FROM pg_stat_activity 
			WHERE datname = '%s' AND pid <> pg_backend_pid()`, testDBName)
		adminClient.Exec(context.Background(), killConnectionsQuery)

		// Eliminar BD de testing
		dropQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName)
		adminClient.Exec(context.Background(), dropQuery)
		adminClient.Close()
	}

	return testClient, testDBName, cleanup
}

// createAdminClient crea cliente administrativo para crear/eliminar BDs
func createAdminClient(config *TestDatabaseConfig) (*postgresql.PostgresqlClient, error) {
	adminConfig := &postgresql.Config{
		Host:        config.Host,
		Port:        config.Port,
		User:        config.User,
		Password:    config.Password,
		Database:    "postgres", // BD administrativa
		MaxConns:    5,
		MinConns:    1,
		MaxLifetime: 30 * time.Minute,
		MaxIdleTime: 5 * time.Minute,
		HealthCheck: false,
		SSLMode:     config.SSLMode,
	}

	return postgresql.NewPostgresqlClient(adminConfig)
}

// createTestClient crea cliente para BD de testing
func createTestClient(config *TestDatabaseConfig, dbName string) (*postgresql.PostgresqlClient, error) {
	testConfig := &postgresql.Config{
		Host:        config.Host,
		Port:        config.Port,
		User:        config.User,
		Password:    config.Password,
		Database:    dbName,
		MaxConns:    5,
		MinConns:    1,
		MaxLifetime: 30 * time.Minute,
		MaxIdleTime: 5 * time.Minute,
		HealthCheck: true,
		SSLMode:     config.SSLMode,
	}

	return postgresql.NewPostgresqlClient(testConfig)
}

// SetupTestTables ejecuta las migraciones necesarias para testing
func SetupTestTables(t *testing.T, client *postgresql.PostgresqlClient) {
	t.Helper()

	// Crear extensión UUID si no existe
	createUUIDExtension := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		
		-- Función para generar UUID v7
		CREATE OR REPLACE FUNCTION uuid_generate_v7() RETURNS UUID AS $$
		BEGIN
			RETURN uuid_generate_v4(); -- Fallback to v4 for simplicity in tests
		END;
		$$ LANGUAGE plpgsql;
	`

	if err := client.Exec(context.Background(), createUUIDExtension); err != nil {
		t.Fatalf("Failed to create UUID extension: %v", err)
	}

	// Crear tablas principales para testing
	createTablesSQL := `
		-- Tabla users
		CREATE TABLE users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
			username VARCHAR(50) NOT NULL UNIQUE,
			phone VARCHAR(20) NULL,
			full_name VARCHAR(200) NOT NULL,
			identification_number VARCHAR(50) NULL,
			email VARCHAR(150) NOT NULL UNIQUE,
			email_token VARCHAR(255) NULL,
			email_token_expires TIMESTAMP NULL,
			email_verified BOOLEAN DEFAULT FALSE NOT NULL,
			password VARCHAR(255) NOT NULL,
			password_reset_token VARCHAR(255) NULL,
			password_reset_expires TIMESTAMP NULL,
			last_password_change TIMESTAMP NULL,
			last_login TIMESTAMP NULL,
			bank_id UUID NULL,
			bank_account_number VARCHAR(50) NULL,
			bank_account_type VARCHAR(50) NULL,
			image_url TEXT NULL,
			is_active BOOLEAN DEFAULT FALSE NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP NULL
		);
		
		-- Índices para users
		CREATE INDEX idx_users_username ON users (username);
		CREATE INDEX idx_users_email ON users (email);
		CREATE INDEX idx_users_email_token ON users (email_token);
		CREATE INDEX idx_users_is_active ON users (is_active);
		
		-- Tipo enum para tenant status
		CREATE TYPE tenant_status AS ENUM ('active', 'inactive', 'suspended');
		
		-- Tabla tenants
		CREATE TABLE tenants (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
			rut VARCHAR(20) NOT NULL,
			business_name VARCHAR(150) NOT NULL,
			email VARCHAR(150) NOT NULL,
			phone VARCHAR(20) NOT NULL,
			address VARCHAR(200) NOT NULL,
			website VARCHAR(150) NOT NULL,
			logo TEXT NULL,
			region_id CHAR(2) NOT NULL,
			commune_id VARCHAR(100) NOT NULL,
			country_id UUID NOT NULL,
			status tenant_status NOT NULL DEFAULT 'active',
			node_number INT NOT NULL,
			tenant_name TEXT NOT NULL,
			created_by UUID NOT NULL,
			updated_by UUID NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP NULL
		);
		
		-- Índices para tenants
		CREATE INDEX idx_tenants_business_name ON tenants (business_name);
		CREATE INDEX idx_tenants_rut ON tenants (rut);
		
		-- Tabla tenant_users (relación many-to-many)
		CREATE TABLE tenant_users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
			tenant_id UUID NOT NULL,
			user_id UUID NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP NULL,
			UNIQUE (tenant_id, user_id)
		);
		
		-- Índices para tenant_users
		CREATE INDEX idx_tenant_users_tenant_id ON tenant_users (tenant_id);
		CREATE INDEX idx_tenant_users_user_id ON tenant_users (user_id);
		CREATE INDEX idx_tenant_users_tenant_user ON tenant_users (tenant_id, user_id);
	`

	if err := client.Exec(context.Background(), createTablesSQL); err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}
}

// TruncateTables limpia todas las tablas para testing
func TruncateTables(t *testing.T, client *postgresql.PostgresqlClient) {
	t.Helper()

	truncateSQL := `
		TRUNCATE TABLE tenant_users CASCADE;
		TRUNCATE TABLE tenants CASCADE;
		TRUNCATE TABLE users CASCADE;
	`

	if err := client.Exec(context.Background(), truncateSQL); err != nil {
		t.Fatalf("Failed to truncate tables: %v", err)
	}
}

// getEnvOrDefault obtiene variable de entorno o valor por defecto
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// CheckPostgreSQLAvailable verifica si PostgreSQL está disponible para testing
func CheckPostgreSQLAvailable() bool {
	config := GetTestDatabaseConfig()
	client, err := createAdminClient(config)
	if err != nil {
		return false
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Ping(ctx) == nil
}

// SkipIfNoPostgreSQL salta el test si PostgreSQL no está disponible
func SkipIfNoPostgreSQL(t *testing.T) {
	if !CheckPostgreSQLAvailable() {
		t.Skip("PostgreSQL not available - skipping integration test")
	}
}

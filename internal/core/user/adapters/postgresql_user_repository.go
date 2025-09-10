package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/user/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
)

// PostgreSQLUserRepository implementa UserRepository usando PostgreSQL
// Conecta con la base de datos misviaticos_control
type PostgreSQLUserRepository struct {
	client *postgresql.PostgresqlClient
}

// NewPostgreSQLUserRepository crea una nueva instancia del repositorio
func NewPostgreSQLUserRepository(client *postgresql.PostgresqlClient) ports.UserRepository {
	return &PostgreSQLUserRepository{
		client: client,
	}
}

// Create crea un nuevo usuario en la base de datos control
func (r *PostgreSQLUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			id, username, phone, full_name, identification_number, email,
			email_token, email_token_expires, email_verified, password,
			password_reset_token, password_reset_expires, last_password_change,
			last_login, bank_id, bank_account_number, bank_account_type,
			image_url, is_active, created, updated
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		)`

	err := r.client.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Phone,
		user.FullName,
		user.IdentificationNumber,
		user.Email,
		user.EmailToken,
		user.EmailTokenExpires,
		user.EmailVerified,
		user.Password,
		user.PasswordResetToken,
		user.PasswordResetExpires,
		user.LastPasswordChange,
		user.LastLogin,
		user.BankID,
		user.BankAccountNumber,
		user.BankAccountType,
		user.ImageURL,
		user.IsActive,
		user.Created,
		user.Updated,
	)

	if err != nil {
		return fmt.Errorf("error creando usuario: %w", err)
	}

	return nil
}

// GetByID obtiene un usuario por su ID
func (r *PostgreSQLUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, username, phone, full_name, identification_number, email,
			   email_token, email_token_expires, email_verified, password,
			   password_reset_token, password_reset_expires, last_password_change,
			   last_login, bank_id, bank_account_number, bank_account_type,
			   image_url, is_active, created, updated, deleted_at
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL`

	user, err := r.scanUser(r.client.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo usuario: %w", err)
	}

	return user, nil
}

// GetByEmail obtiene un usuario por su email
func (r *PostgreSQLUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, username, phone, full_name, identification_number, email,
			   email_token, email_token_expires, email_verified, password,
			   password_reset_token, password_reset_expires, last_password_change,
			   last_login, bank_id, bank_account_number, bank_account_type,
			   image_url, is_active, created, updated, deleted_at
		FROM users 
		WHERE email = $1 AND deleted_at IS NULL`

	user, err := r.scanUser(r.client.QueryRow(ctx, query, email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo usuario: %w", err)
	}

	return user, nil
}

// GetByUsername obtiene un usuario por su username
func (r *PostgreSQLUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, username, phone, full_name, identification_number, email,
			   email_token, email_token_expires, email_verified, password,
			   password_reset_token, password_reset_expires, last_password_change,
			   last_login, bank_id, bank_account_number, bank_account_type,
			   image_url, is_active, created, updated, deleted_at
		FROM users 
		WHERE username = $1 AND deleted_at IS NULL`

	user, err := r.scanUser(r.client.QueryRow(ctx, query, username))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo usuario: %w", err)
	}

	return user, nil
}

// GetByEmailToken obtiene un usuario por su token de verificación de email
func (r *PostgreSQLUserRepository) GetByEmailToken(ctx context.Context, token string) (*domain.User, error) {
	query := `
		SELECT id, username, phone, full_name, identification_number, email,
			   email_token, email_token_expires, email_verified, password,
			   password_reset_token, password_reset_expires, last_password_change,
			   last_login, bank_id, bank_account_number, bank_account_type,
			   image_url, is_active, created, updated, deleted_at
		FROM users 
		WHERE email_token = $1 AND deleted_at IS NULL`

	user, err := r.scanUser(r.client.QueryRow(ctx, query, token))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo usuario: %w", err)
	}

	return user, nil
}

// Update actualiza un usuario existente
func (r *PostgreSQLUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET
			username = $2, phone = $3, full_name = $4, identification_number = $5,
			email = $6, email_token = $7, email_token_expires = $8, email_verified = $9,
			password = $10, password_reset_token = $11, password_reset_expires = $12,
			last_password_change = $13, last_login = $14, bank_id = $15,
			bank_account_number = $16, bank_account_type = $17, image_url = $18,
			is_active = $19, updated = $20
		WHERE id = $1 AND deleted_at IS NULL`

	user.Updated = time.Now()

	err := r.client.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Phone,
		user.FullName,
		user.IdentificationNumber,
		user.Email,
		user.EmailToken,
		user.EmailTokenExpires,
		user.EmailVerified,
		user.Password,
		user.PasswordResetToken,
		user.PasswordResetExpires,
		user.LastPasswordChange,
		user.LastLogin,
		user.BankID,
		user.BankAccountNumber,
		user.BankAccountType,
		user.ImageURL,
		user.IsActive,
		user.Updated,
	)

	if err != nil {
		return fmt.Errorf("error actualizando usuario: %w", err)
	}

	return nil
}

// Delete elimina lógicamente un usuario (soft delete)
func (r *PostgreSQLUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users 
		SET deleted_at = $2, updated = $2 
		WHERE id = $1 AND deleted_at IS NULL`

	now := time.Now()
	err := r.client.Exec(ctx, query, id, now)
	if err != nil {
		return fmt.Errorf("error eliminando usuario: %w", err)
	}

	return nil
}

// ExistsByEmail verifica si existe un usuario con el email dado
func (r *PostgreSQLUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)`
	
	var exists bool
	err := r.client.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error verificando email: %w", err)
	}

	return exists, nil
}

// ExistsByUsername verifica si existe un usuario con el username dado
func (r *PostgreSQLUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND deleted_at IS NULL)`
	
	var exists bool
	err := r.client.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error verificando username: %w", err)
	}

	return exists, nil
}

// GetUserTenants obtiene todos los tenants asociados a un usuario
func (r *PostgreSQLUserRepository) GetUserTenants(ctx context.Context, userID uuid.UUID) ([]*domain.TenantUser, error) {
	query := `
		SELECT id, tenant_id, user_id, created, updated, deleted_at
		FROM tenant_users 
		WHERE user_id = $1 AND deleted_at IS NULL`

	rows, err := r.client.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo tenants del usuario: %w", err)
	}
	defer rows.Close()

	var tenantUsers []*domain.TenantUser
	for rows.Next() {
		tenantUser, err := r.scanTenantUser(rows)
		if err != nil {
			return nil, fmt.Errorf("error escaneando tenant_user: %w", err)
		}
		tenantUsers = append(tenantUsers, tenantUser)
	}

	return tenantUsers, nil
}

// AddUserToTenant asocia un usuario a un tenant
func (r *PostgreSQLUserRepository) AddUserToTenant(ctx context.Context, tenantUser *domain.TenantUser) error {
	query := `
		INSERT INTO tenant_users (id, tenant_id, user_id, created, updated)
		VALUES ($1, $2, $3, $4, $5)`

	err := r.client.Exec(ctx, query,
		tenantUser.ID,
		tenantUser.TenantID,
		tenantUser.UserID,
		tenantUser.Created,
		tenantUser.Updated,
	)

	if err != nil {
		return fmt.Errorf("error asociando usuario al tenant: %w", err)
	}

	return nil
}

// RemoveUserFromTenant desasocia un usuario de un tenant (soft delete)
func (r *PostgreSQLUserRepository) RemoveUserFromTenant(ctx context.Context, userID, tenantID uuid.UUID) error {
	query := `
		UPDATE tenant_users 
		SET deleted_at = $3, updated = $3 
		WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL`

	now := time.Now()
	err := r.client.Exec(ctx, query, userID, tenantID, now)
	if err != nil {
		return fmt.Errorf("error desasociando usuario del tenant: %w", err)
	}

	return nil
}

// scanUser escanea una fila de usuario desde la base de datos
func (r *PostgreSQLUserRepository) scanUser(row pgx.Row) (*domain.User, error) {
	user := &domain.User{}
	var deletedAt sql.NullTime

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Phone,
		&user.FullName,
		&user.IdentificationNumber,
		&user.Email,
		&user.EmailToken,
		&user.EmailTokenExpires,
		&user.EmailVerified,
		&user.Password,
		&user.PasswordResetToken,
		&user.PasswordResetExpires,
		&user.LastPasswordChange,
		&user.LastLogin,
		&user.BankID,
		&user.BankAccountNumber,
		&user.BankAccountType,
		&user.ImageURL,
		&user.IsActive,
		&user.Created,
		&user.Updated,
		&deletedAt,
	)

	if err != nil {
		return nil, err
	}

	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return user, nil
}

// scanTenantUser escanea una fila de tenant_user desde la base de datos
func (r *PostgreSQLUserRepository) scanTenantUser(rows pgx.Rows) (*domain.TenantUser, error) {
	tenantUser := &domain.TenantUser{}
	var deletedAt sql.NullTime

	err := rows.Scan(
		&tenantUser.ID,
		&tenantUser.TenantID,
		&tenantUser.UserID,
		&tenantUser.Created,
		&tenantUser.Updated,
		&deletedAt,
	)

	if err != nil {
		return nil, err
	}

	if deletedAt.Valid {
		tenantUser.DeletedAt = &deletedAt.Time
	}

	return tenantUser, nil
}
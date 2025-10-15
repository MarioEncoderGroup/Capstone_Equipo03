package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/user/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
			return nil, sharedErrors.ErrUserNotFound
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
			return nil, sharedErrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("error obteniendo usuario: %w", err)
	}

	return user, nil
}

// GetByPasswordResetToken obtiene un usuario por su token de reset de contraseña
func (r *PostgreSQLUserRepository) GetByPasswordResetToken(ctx context.Context, token string) (*domain.User, error) {
	query := `
		SELECT id, username, phone, full_name, identification_number, email,
			   email_token, email_token_expires, email_verified, password,
			   password_reset_token, password_reset_expires, last_password_change,
			   last_login, bank_id, bank_account_number, bank_account_type,
			   image_url, is_active, created, updated, deleted_at
		FROM users
		WHERE password_reset_token = $1 AND deleted_at IS NULL`

	user, err := r.scanUser(r.client.QueryRow(ctx, query, token))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sharedErrors.ErrUserNotFound
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

	// Usar Pool.Exec para obtener CommandTag y verificar rows affected
	result, err := r.client.Pool.Exec(ctx, query, id, now)
	if err != nil {
		return fmt.Errorf("error eliminando usuario: %w", err)
	}

	// Verificar que se afectó al menos una fila
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("usuario no encontrado o ya eliminado")
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

// GetTenantsByUser obtiene todos los tenant_users por userID (alias de GetUserTenants)
func (r *PostgreSQLUserRepository) GetTenantsByUser(ctx context.Context, userID uuid.UUID) ([]*domain.TenantUser, error) {
	// Delegamos a GetUserTenants ya que ambos métodos hacen lo mismo
	return r.GetUserTenants(ctx, userID)
}

// UserHasAccessToTenant verifica si un usuario tiene acceso a un tenant
func (r *PostgreSQLUserRepository) UserHasAccessToTenant(ctx context.Context, userID, tenantID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM tenant_users 
			WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)`

	var exists bool
	err := r.client.QueryRow(ctx, query, userID, tenantID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error verificando acceso del usuario al tenant: %w", err)
	}

	return exists, nil
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

// GetUsers obtiene una lista paginada de usuarios filtrados por tenant
// IMPORTANTE: Usa JOIN con tenant_users para garantizar aislamiento multi-tenant
func (r *PostgreSQLUserRepository) GetUsers(ctx context.Context, tenantID *uuid.UUID, offset, limit int, sortBy, sortDir, search string) ([]*domain.User, int64, error) {
	// Validación de seguridad: tenant_id es obligatorio
	if tenantID == nil {
		return nil, 0, fmt.Errorf("tenant_id es requerido para consultar usuarios")
	}

	// Query para contar usuarios del tenant
	countQuery := `
		SELECT COUNT(DISTINCT u.id) 
		FROM users u
		INNER JOIN tenant_users tu ON u.id = tu.user_id
		WHERE tu.tenant_id = $1 
		  AND tu.deleted_at IS NULL 
		  AND u.deleted_at IS NULL`

	var countArgs []interface{}
	countArgs = append(countArgs, *tenantID)
	argIndex := 2

	if search != "" {
		countQuery += ` AND (u.full_name ILIKE $` + fmt.Sprintf("%d", argIndex) +
			` OR u.email ILIKE $` + fmt.Sprintf("%d", argIndex) +
			` OR u.username ILIKE $` + fmt.Sprintf("%d", argIndex) + `)`
		countArgs = append(countArgs, "%"+search+"%")
	}

	var total int64
	err := r.client.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error contando usuarios: %w", err)
	}

	// Query principal con JOIN para filtrar por tenant
	query := `
		SELECT DISTINCT u.id, u.username, u.phone, u.full_name, u.identification_number, u.email,
			   u.email_token, u.email_token_expires, u.email_verified, u.password,
			   u.password_reset_token, u.password_reset_expires, u.last_password_change,
			   u.last_login, u.bank_id, u.bank_account_number, u.bank_account_type,
			   u.image_url, u.is_active, u.created, u.updated, u.deleted_at
		FROM users u
		INNER JOIN tenant_users tu ON u.id = tu.user_id
		WHERE tu.tenant_id = $1 
		  AND tu.deleted_at IS NULL 
		  AND u.deleted_at IS NULL`

	var args []interface{}
	args = append(args, *tenantID)
	argIndex = 2

	if search != "" {
		query += ` AND (u.full_name ILIKE $` + fmt.Sprintf("%d", argIndex) +
			` OR u.email ILIKE $` + fmt.Sprintf("%d", argIndex) +
			` OR u.username ILIKE $` + fmt.Sprintf("%d", argIndex) + `)`
		args = append(args, "%"+search+"%")
		argIndex++
	}

	query += ` ORDER BY u.` + sortBy + ` ` + sortDir
	query += ` LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.client.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error consultando usuarios: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user, err := r.scanUser(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("error escaneando usuario: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterando usuarios: %w", err)
	}

	return users, total, nil
}

// CheckUserExists verifica si un usuario existe por ID
func (r *PostgreSQLUserRepository) CheckUserExists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.client.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error verificando existencia de usuario: %w", err)
	}

	return exists, nil
}

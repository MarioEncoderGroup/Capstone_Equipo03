package ports

import (
	"context"

	domain_auth "github.com/JoseLuis21/mv-backend/internal/core/auth/domain"
	domain_tenant "github.com/JoseLuis21/mv-backend/internal/core/tenant/domain"
	domain_user "github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	"github.com/google/uuid"
)

// AuthService define el contrato para el servicio de autenticación
// Siguiendo el patrón del proyecto de referencia
type AuthService interface {
	// Register registra un nuevo usuario individual o con tenant
	Register(ctx context.Context, req *domain_auth.AuthRegisterDto) (*domain_auth.AuthRegisterResponse, error)
	
	// VerifyUserEmail verifica el email de un usuario usando el token
	VerifyUserEmail(ctx context.Context, token string) error
	
	// Login autentica un usuario y retorna tokens
	Login(ctx context.Context, req *domain_auth.AuthLoginDto) (*domain_auth.AuthLoginResponse, error)
	
	// ForgotPassword inicia el proceso de recuperación de contraseña
	ForgotPassword(ctx context.Context, email string) (*domain_user.User, error)
	
	// ResetPassword resetea la contraseña usando un token
	ResetPassword(ctx context.Context, token string, newPassword string) error
	
	// ResendEmailVerification reenvía el email de verificación
	ResendEmailVerification(ctx context.Context, email string) error

	// SelectTenant selecciona un tenant específico post-login y genera nuevos tokens con tenant_id
	SelectTenant(ctx context.Context, tenant *domain_tenant.Tenant, userID uuid.UUID) (*domain_tenant.SelectTenantResponseDto, error)

	// PASO 5: Métodos para manejo de refresh tokens
	// RefreshAccessToken renueva el access token usando un refresh token
	RefreshAccessToken(ctx context.Context, refreshToken string) (*domain_auth.RefreshTokenResponse, error)
	
	// RevokeRefreshToken revoca un refresh token
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
}


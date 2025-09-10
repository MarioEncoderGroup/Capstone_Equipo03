package ports

import (
	"context"
	"github.com/JoseLuis21/mv-backend/internal/core/user/domain"
)

// PasswordHasher define el contrato para el hash de contraseñas
type PasswordHasher interface {
	// Hash genera un hash seguro de la contraseña
	Hash(password string) (string, error)
	
	// Verify verifica si una contraseña coincide con su hash
	Verify(hashedPassword, password string) error
}

// TokenGenerator define el contrato para la generación de tokens
type TokenGenerator interface {
	// GenerateEmailVerificationToken genera un token para verificación de email
	GenerateEmailVerificationToken() (string, error)
	
	// GeneratePasswordResetToken genera un token para reset de contraseña
	GeneratePasswordResetToken() (string, error)
	
	// GenerateJWT genera un token JWT para autenticación
	GenerateJWT(user *domain.User) (string, error)
}

// EmailService define el contrato para el envío de emails
type EmailService interface {
	// SendEmailVerification envía un email de verificación al usuario
	SendEmailVerification(ctx context.Context, user *domain.User, token string) error
	
	// SendPasswordReset envía un email para reset de contraseña
	SendPasswordReset(ctx context.Context, user *domain.User, token string) error
	
	// SendWelcomeEmail envía email de bienvenida después del registro
	SendWelcomeEmail(ctx context.Context, user *domain.User) error
}
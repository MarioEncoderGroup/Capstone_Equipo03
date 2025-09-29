package domain

import (
	domain_user "github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	"github.com/google/uuid"
)

// AuthRegisterDto representa la solicitud de registro con campos requeridos por jefe técnico
// Incluye firstname, lastname, email, phone, password, password_confirm
type AuthRegisterDto struct {
	FirstName       string `json:"firstname" validate:"required,min=2,max=100"`
	LastName        string `json:"lastname" validate:"required,min=2,max=100"`
	Email           string `json:"email" validate:"required,email,max=150"`
	Phone           string `json:"phone" validate:"required,min=8,max=20"`
	Password        string `json:"password" validate:"required,min=8,max=255"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
}


// AuthRegisterResponse respuesta del registro
type AuthRegisterResponse struct {
	ID                        uuid.UUID `json:"id"`
	FirstName                 string    `json:"firstname"`
	LastName                  string    `json:"lastname"`
	FullName                  string    `json:"full_name"` // Para backward compatibility
	Email                     string    `json:"email"`
	Phone                     string    `json:"phone"`
	EmailToken                string    `json:"email_token"`
	RequiresEmailVerification bool       `json:"requires_email_verification"`
	Message                   string     `json:"message"`
}

// AuthLoginDto para autenticación de usuarios
type AuthLoginDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthLoginResponse respuesta del login - PASO 5: Agregado refresh token
type AuthLoginResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`  // ← AGREGADO PASO 5
	ExpiresIn    int64            `json:"expires_in"`     // ← AGREGADO PASO 5
	TokenType    string           `json:"token_type"`     // ← AGREGADO PASO 5
	User         domain_user.User `json:"user"`
}

// VerifyEmailDto para verificación de email
type VerifyEmailDto struct {
	Token string `json:"token" validate:"required,min=1"`
}

// ResendVerificationDto para reenviar verificación
type ResendVerificationDto struct {
	Email string `json:"email" validate:"required,email,max=150"`
}

// ForgotPasswordDto para solicitar reset de contraseña
type ForgotPasswordDto struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordDto para resetear contraseña
type ResetPasswordDto struct {
	Token           string `json:"token" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=NewPassword"`
}

// PASO 5: DTOs para Refresh Token System
// RefreshTokenDto para renovar tokens
type RefreshTokenDto struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse respuesta al renovar token
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

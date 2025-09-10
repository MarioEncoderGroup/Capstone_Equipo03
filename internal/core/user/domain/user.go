package domain

import (
	"time"
	"github.com/google/uuid"
)

// User representa la entidad principal de usuario en el dominio
// Mapea directamente con la tabla 'users' del control database
type User struct {
	ID                     uuid.UUID  `json:"id"`
	Username               string     `json:"username"`
	Phone                  *string    `json:"phone"`
	FullName               string     `json:"full_name"`
	IdentificationNumber   *string    `json:"identification_number"` // RUT chileno
	Email                  string     `json:"email"`
	EmailToken             *string    `json:"email_token"`
	EmailTokenExpires      *time.Time `json:"email_token_expires"`
	EmailVerified          bool       `json:"email_verified"`
	Password               string     `json:"password"`
	PasswordResetToken     *string    `json:"password_reset_token"`
	PasswordResetExpires   *time.Time `json:"password_reset_expires"`
	LastPasswordChange     *time.Time `json:"last_password_change"`
	LastLogin              *time.Time `json:"last_login"`
	BankID                 *uuid.UUID `json:"bank_id"`
	BankAccountNumber      *string    `json:"bank_account_number"`
	BankAccountType        *string    `json:"bank_account_type"`
	ImageURL               *string    `json:"image_url"`
	IsActive               bool       `json:"is_active"`
	Created                time.Time  `json:"created"`
	Updated                time.Time  `json:"updated"`
	DeletedAt              *time.Time `json:"deleted_at"`
}

// NewUser crea una nueva instancia de usuario con valores por defecto
// Aplica reglas de negocio iniciales (usuario inactivo hasta verificar email)
func NewUser(username, fullName, email, hashedPassword string) *User {
	now := time.Now()
	return &User{
		ID:            uuid.New(),
		Username:      username,
		FullName:      fullName,
		Email:         email,
		Password:      hashedPassword,
		EmailVerified: false,
		IsActive:      false, // Inactivo hasta verificar email
		Created:       now,
		Updated:       now,
	}
}

// IsEmailTokenValid valida si el token de email es válido y no ha expirado
func (u *User) IsEmailTokenValid(token string) bool {
	if u.EmailToken == nil || *u.EmailToken != token {
		return false
	}
	if u.EmailTokenExpires == nil || time.Now().After(*u.EmailTokenExpires) {
		return false
	}
	return true
}

// ActivateUser activa el usuario después de verificar el email
func (u *User) ActivateUser() {
	u.EmailVerified = true
	u.IsActive = true
	u.EmailToken = nil
	u.EmailTokenExpires = nil
	u.Updated = time.Now()
}

// SetEmailVerificationToken establece el token de verificación de email
func (u *User) SetEmailVerificationToken(token string, expiresIn time.Duration) {
	u.EmailToken = &token
	expires := time.Now().Add(expiresIn)
	u.EmailTokenExpires = &expires
	u.Updated = time.Now()
}
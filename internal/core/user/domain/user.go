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
	FirstName              string     `json:"first_name"`
	LastName               string     `json:"last_name"`
	FullName               string     `json:"full_name"` // Calculado: FirstName + " " + LastName
	Phone                  *string    `json:"phone"`
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
// Genera username automáticamente basado en el email
func NewUser(firstName, lastName, email, phone, hashedPassword string) *User {
	now := time.Now()
	
	// Generar username automáticamente del email (antes del @)
	username := generateUsernameFromEmail(email)
	
	// Construir full name a partir de first name y last name
	fullName := firstName + " " + lastName
	
	// Manejar phone como puntero (puede ser vacío)
	var phonePtr *string
	if phone != "" {
		phonePtr = &phone
	}
	
	return &User{
		ID:            uuid.New(),
		Username:      username,
		FirstName:     firstName,
		LastName:      lastName,
		FullName:      fullName,
		Phone:         phonePtr,
		Email:         email,
		Password:      hashedPassword,
		EmailVerified: false,
		IsActive:      false, // Inactivo hasta verificar email
		Created:       now,
		Updated:       now,
	}
}

// generateUsernameFromEmail genera un username único basado en el email
func generateUsernameFromEmail(email string) string {
	// Tomar la parte antes del @ y limpiar caracteres especiales
	username := email
	if atIndex := findAtIndex(email); atIndex != -1 {
		username = email[:atIndex]
	}
	
	// Limpiar caracteres especiales (mantener solo alfanuméricos)
	cleaned := ""
	for _, char := range username {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			cleaned += string(char)
		}
	}
	
	// Asegurar que tenga al menos 3 caracteres
	if len(cleaned) < 3 {
		cleaned = "user" + cleaned
	}
	
	return cleaned
}

// findAtIndex encuentra la posición del @ en el email
func findAtIndex(email string) int {
	for i, char := range email {
		if char == '@' {
			return i
		}
	}
	return -1
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

// SetPasswordResetToken establece el token de recuperación de contraseña
// Regla de negocio: Solo un token activo por usuario, expiración de 1 hora
func (u *User) SetPasswordResetToken(token string, expiresIn time.Duration) {
	u.PasswordResetToken = &token
	expires := time.Now().UTC().Add(expiresIn)
	u.PasswordResetExpires = &expires
	u.Updated = time.Now().UTC()
}

// IsPasswordResetTokenValid valida si el token de reset es válido y no ha expirado
// Regla de negocio: Token debe existir, coincidir exactamente y no estar expirado
func (u *User) IsPasswordResetTokenValid(token string) bool {
	if u.PasswordResetToken == nil || *u.PasswordResetToken != token {
		return false
	}
	if u.PasswordResetExpires == nil || time.Now().UTC().After(*u.PasswordResetExpires) {
		return false
	}
	return true
}

// ChangePassword cambia la contraseña del usuario y limpia el token de reset
// Regla de negocio: Invalidar token después del cambio, actualizar timestamp
func (u *User) ChangePassword(newHashedPassword string) {
	u.Password = newHashedPassword
	u.PasswordResetToken = nil
	u.PasswordResetExpires = nil
	now := time.Now().UTC()
	u.LastPasswordChange = &now
	u.Updated = now
}

// HasActivePasswordResetToken verifica si el usuario tiene un token de reset activo
// Regla de negocio: Para prevenir múltiples tokens simultáneos
func (u *User) HasActivePasswordResetToken() bool {
	if u.PasswordResetToken == nil || u.PasswordResetExpires == nil {
		return false
	}
	return time.Now().UTC().Before(*u.PasswordResetExpires)
}
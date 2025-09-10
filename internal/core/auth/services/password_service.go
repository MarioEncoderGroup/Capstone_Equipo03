package services

import (
	"errors"
	"github.com/JoseLuis21/mv-backend/internal/core/auth/ports"
	
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasherImpl implementa el hash seguro de contraseñas usando bcrypt
type PasswordHasherImpl struct {
	cost int // Costo de bcrypt (10-15 recomendado)
}

// NewPasswordHasher crea una nueva instancia del hasher de contraseñas
func NewPasswordHasher() ports.PasswordHasher {
	return &PasswordHasherImpl{
		cost: 12, // Costo balanceado entre seguridad y rendimiento
	}
}

// Hash genera un hash seguro de la contraseña usando bcrypt
func (ph *PasswordHasherImpl) Hash(password string) (string, error) {
	if password == "" {
		return "", errors.New("contraseña no puede estar vacía")
	}

	// Validar longitud mínima y máxima
	if len(password) < 8 {
		return "", errors.New("contraseña debe tener al menos 8 caracteres")
	}

	if len(password) > 72 { // Límite de bcrypt
		return "", errors.New("contraseña muy larga (máximo 72 caracteres)")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), ph.cost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// Verify verifica si una contraseña coincide con su hash
func (ph *PasswordHasherImpl) Verify(hashedPassword, password string) error {
	if hashedPassword == "" || password == "" {
		return errors.New("hash y contraseña son requeridos")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("contraseña incorrecta")
		}
		return err
	}

	return nil
}
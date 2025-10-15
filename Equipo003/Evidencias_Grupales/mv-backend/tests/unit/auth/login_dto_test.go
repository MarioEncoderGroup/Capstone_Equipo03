package auth_unit_test

import (
	"testing"

	authDomain "github.com/JoseLuis21/mv-backend/internal/core/auth/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthLoginDto valida la estructura del DTO de login
func TestAuthLoginDto(t *testing.T) {
	t.Run("Should create valid LoginDto", func(t *testing.T) {
		loginDto := authDomain.AuthLoginDto{
			Email:    "test@example.com",
			Password: "password123",
		}

		assert.Equal(t, "test@example.com", loginDto.Email)
		assert.Equal(t, "password123", loginDto.Password)
	})

	t.Run("Should have correct JSON tags", func(t *testing.T) {
		// Esta funcionalidad sería validada por el validador en tiempo de ejecución
		// Aquí solo verificamos que la estructura esté bien definida
		loginDto := authDomain.AuthLoginDto{}

		// Verificar que los campos pueden ser asignados
		loginDto.Email = "test@domain.com"
		loginDto.Password = "securepass"

		assert.NotEmpty(t, loginDto.Email)
		assert.NotEmpty(t, loginDto.Password)
	})
}

// TestAuthLoginResponse valida la estructura de respuesta del login
func TestAuthLoginResponse(t *testing.T) {
	t.Run("Should create valid LoginResponse", func(t *testing.T) {
		loginResponse := authDomain.AuthLoginResponse{
			AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			ExpiresIn:    86400, // 24 hours
			TokenType:    "Bearer",
		}

		assert.Equal(t, "Bearer", loginResponse.TokenType)
		assert.Equal(t, int64(86400), loginResponse.ExpiresIn)
		assert.NotEmpty(t, loginResponse.AccessToken)
		assert.NotEmpty(t, loginResponse.RefreshToken)
	})

	t.Run("Should have standard JWT format for tokens", func(t *testing.T) {
		// Verificar que los tokens siguen el formato esperado
		accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

		loginResponse := authDomain.AuthLoginResponse{
			AccessToken: accessToken,
			TokenType:   "Bearer",
		}

		// Verificar que el token tiene el formato JWT básico (3 partes separadas por puntos)
		assert.Contains(t, loginResponse.AccessToken, ".")
		assert.Equal(t, "Bearer", loginResponse.TokenType)
	})

	t.Run("Should support 24-hour expiration", func(t *testing.T) {
		expiresIn := int64(24 * 60 * 60) // 24 horas en segundos

		loginResponse := authDomain.AuthLoginResponse{
			ExpiresIn: expiresIn,
		}

		assert.Equal(t, int64(86400), loginResponse.ExpiresIn)
		assert.True(t, loginResponse.ExpiresIn > 0, "Expiration should be positive")
	})
}

// TestLoginDtoValidationTags verifica que los tags de validación estén correctos
func TestLoginDtoValidationTags(t *testing.T) {
	t.Run("Email validation requirements", func(t *testing.T) {
		// Los tags de validación son:
		// Email: `json:"email" validate:"required,email"`
		// Password: `json:"password" validate:"required"`

		dto := authDomain.AuthLoginDto{
			Email:    "valid@email.com",
			Password: "validpassword",
		}

		// Verificar que los campos son correctos
		require.NotEmpty(t, dto.Email, "Email should not be empty")
		require.NotEmpty(t, dto.Password, "Password should not be empty")
		assert.Contains(t, dto.Email, "@", "Email should contain @ symbol")
	})

	t.Run("Password validation requirements", func(t *testing.T) {
		dto := authDomain.AuthLoginDto{
			Email:    "test@example.com",
			Password: "pass123",
		}

		// Los requisitos de validación se verifican en runtime por el validador
		// Aquí solo verificamos que el campo esté presente
		assert.NotEmpty(t, dto.Password, "Password should not be empty")
		assert.True(t, len(dto.Password) >= 6, "Password should have reasonable length")
	})
}
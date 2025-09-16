package domain_test

import (
	"testing"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_SetPasswordResetToken(t *testing.T) {
	t.Run("Should set password reset token with correct expiration", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Juan", "Pérez", "juan@test.cl", "+56912345678", "hashedpass")
		token := "test-reset-token-123"
		duration := 1 * time.Hour

		// Act
		before := time.Now()
		user.SetPasswordResetToken(token, duration)
		after := time.Now()

		// Assert
		assert.NotNil(t, user.PasswordResetToken, "Password reset token should be set")
		assert.Equal(t, token, *user.PasswordResetToken, "Token should match provided value")

		require.NotNil(t, user.PasswordResetExpires, "Password reset expiration should be set")
		expectedExpiry := before.Add(duration)
		actualExpiry := *user.PasswordResetExpires

		assert.True(t, actualExpiry.After(expectedExpiry.Add(-time.Second)),
			"Expiry should be approximately 1 hour from now (lower bound)")
		assert.True(t, actualExpiry.Before(after.Add(duration).Add(time.Second)),
			"Expiry should be approximately 1 hour from now (upper bound)")

		assert.False(t, user.Updated.IsZero(), "Updated timestamp should be set")
	})

	t.Run("Should overwrite existing password reset token", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("María", "González", "maria@test.cl", "+56987654321", "hashedpass")
		oldToken := "old-token"
		newToken := "new-token"

		user.SetPasswordResetToken(oldToken, 1*time.Hour)
		oldExpiry := *user.PasswordResetExpires

		// Act
		time.Sleep(10 * time.Millisecond) // Pequeña pausa para diferencias de timestamp
		user.SetPasswordResetToken(newToken, 2*time.Hour)

		// Assert
		assert.Equal(t, newToken, *user.PasswordResetToken, "Should overwrite with new token")
		assert.True(t, user.PasswordResetExpires.After(oldExpiry), "New expiry should be later than old expiry")
	})
}

func TestUser_IsPasswordResetTokenValid(t *testing.T) {
	t.Run("Should return true for valid non-expired token", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Carlos", "Silva", "carlos@test.cl", "+56911111111", "hashedpass")
		token := "valid-token-123"
		user.SetPasswordResetToken(token, 1*time.Hour)

		// Act & Assert
		assert.True(t, user.IsPasswordResetTokenValid(token), "Valid token should return true")
	})

	t.Run("Should return false for wrong token", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Ana", "López", "ana@test.cl", "+56922222222", "hashedpass")
		correctToken := "correct-token"
		wrongToken := "wrong-token"
		user.SetPasswordResetToken(correctToken, 1*time.Hour)

		// Act & Assert
		assert.False(t, user.IsPasswordResetTokenValid(wrongToken), "Wrong token should return false")
	})

	t.Run("Should return false for expired token", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Luis", "Martín", "luis@test.cl", "+56933333333", "hashedpass")
		token := "expired-token"

		// Establecer token con expiración inmediata
		user.SetPasswordResetToken(token, 1*time.Millisecond)
		time.Sleep(10 * time.Millisecond) // Esperar a que expire

		// Act & Assert
		assert.False(t, user.IsPasswordResetTokenValid(token), "Expired token should return false")
	})

	t.Run("Should return false when no token is set", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Pedro", "Ruiz", "pedro@test.cl", "+56944444444", "hashedpass")
		token := "any-token"

		// Act & Assert
		assert.False(t, user.IsPasswordResetTokenValid(token), "Should return false when no token is set")
	})

	t.Run("Should return false when token is nil", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Sofia", "Torres", "sofia@test.cl", "+56955555555", "hashedpass")
		user.PasswordResetToken = nil
		user.PasswordResetExpires = &time.Time{}

		// Act & Assert
		assert.False(t, user.IsPasswordResetTokenValid("any-token"), "Should return false when token is nil")
	})

	t.Run("Should return false when expiry is nil", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Diego", "Vargas", "diego@test.cl", "+56966666666", "hashedpass")
		token := "test-token"
		user.PasswordResetToken = &token
		user.PasswordResetExpires = nil

		// Act & Assert
		assert.False(t, user.IsPasswordResetTokenValid(token), "Should return false when expiry is nil")
	})
}

func TestUser_ChangePassword(t *testing.T) {
	t.Run("Should change password and clear reset token", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Roberto", "Herrera", "roberto@test.cl", "+56977777777", "oldhashedpass")
		resetToken := "reset-token-123"
		user.SetPasswordResetToken(resetToken, 1*time.Hour)

		newPassword := "newhashedpassword"

		// Verificar estado inicial
		require.NotNil(t, user.PasswordResetToken, "Reset token should be set initially")
		require.NotNil(t, user.PasswordResetExpires, "Reset expiry should be set initially")

		// Act
		beforeChange := time.Now()
		user.ChangePassword(newPassword)
		afterChange := time.Now()

		// Assert
		assert.Equal(t, newPassword, user.Password, "Password should be updated")
		assert.Nil(t, user.PasswordResetToken, "Reset token should be cleared")
		assert.Nil(t, user.PasswordResetExpires, "Reset expiry should be cleared")

		require.NotNil(t, user.LastPasswordChange, "LastPasswordChange should be set")
		assert.True(t, user.LastPasswordChange.After(beforeChange.Add(-time.Second)),
			"LastPasswordChange should be recent (lower bound)")
		assert.True(t, user.LastPasswordChange.Before(afterChange.Add(time.Second)),
			"LastPasswordChange should be recent (upper bound)")

		assert.True(t, user.Updated.After(beforeChange.Add(-time.Second)),
			"Updated timestamp should be recent")
	})

	t.Run("Should work even without existing reset token", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Elena", "Castro", "elena@test.cl", "+56988888888", "oldhashedpass")
		newPassword := "newhashedpassword"

		// Act
		user.ChangePassword(newPassword)

		// Assert
		assert.Equal(t, newPassword, user.Password, "Password should be updated")
		assert.Nil(t, user.PasswordResetToken, "Reset token should remain nil")
		assert.Nil(t, user.PasswordResetExpires, "Reset expiry should remain nil")
		assert.NotNil(t, user.LastPasswordChange, "LastPasswordChange should be set")
	})
}

func TestUser_HasActivePasswordResetToken(t *testing.T) {
	t.Run("Should return true for active token", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Manuel", "Rojas", "manuel@test.cl", "+56999999999", "hashedpass")
		user.SetPasswordResetToken("active-token", 1*time.Hour)

		// Act & Assert
		assert.True(t, user.HasActivePasswordResetToken(), "Should return true for active token")
	})

	t.Run("Should return false for expired token", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Carmen", "Mendoza", "carmen@test.cl", "+56910101010", "hashedpass")
		user.SetPasswordResetToken("expired-token", 1*time.Millisecond)
		time.Sleep(10 * time.Millisecond) // Esperar a que expire

		// Act & Assert
		assert.False(t, user.HasActivePasswordResetToken(), "Should return false for expired token")
	})

	t.Run("Should return false when no token is set", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Andrés", "Soto", "andres@test.cl", "+56920202020", "hashedpass")

		// Act & Assert
		assert.False(t, user.HasActivePasswordResetToken(), "Should return false when no token is set")
	})

	t.Run("Should return false when token is nil", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Valentina", "Pérez", "valentina@test.cl", "+56930303030", "hashedpass")
		user.PasswordResetToken = nil
		futureTime := time.Now().Add(1 * time.Hour)
		user.PasswordResetExpires = &futureTime

		// Act & Assert
		assert.False(t, user.HasActivePasswordResetToken(), "Should return false when token is nil")
	})

	t.Run("Should return false when expiry is nil", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Ignacio", "Lagos", "ignacio@test.cl", "+56940404040", "hashedpass")
		token := "test-token"
		user.PasswordResetToken = &token
		user.PasswordResetExpires = nil

		// Act & Assert
		assert.False(t, user.HasActivePasswordResetToken(), "Should return false when expiry is nil")
	})
}

// Test de integración entre métodos
func TestPasswordResetFlow_Integration(t *testing.T) {
	t.Run("Complete password reset flow", func(t *testing.T) {
		// Arrange
		user := domain.NewUser("Isabella", "Morales", "isabella@test.cl", "+56950505050", "originalpassword")
		resetToken := "integration-test-token"
		newPassword := "newpasswordintegration"

		// Verificar estado inicial
		assert.False(t, user.HasActivePasswordResetToken(), "Initially should not have active token")
		assert.False(t, user.IsPasswordResetTokenValid(resetToken), "Token should not be valid initially")

		// Step 1: Set reset token
		user.SetPasswordResetToken(resetToken, 1*time.Hour)

		// Verify token is active
		assert.True(t, user.HasActivePasswordResetToken(), "Should have active token after setting")
		assert.True(t, user.IsPasswordResetTokenValid(resetToken), "Token should be valid after setting")

		// Step 2: Change password using token
		user.ChangePassword(newPassword)

		// Verify final state
		assert.Equal(t, newPassword, user.Password, "Password should be changed")
		assert.False(t, user.HasActivePasswordResetToken(), "Should not have active token after password change")
		assert.False(t, user.IsPasswordResetTokenValid(resetToken), "Token should not be valid after password change")
		assert.NotNil(t, user.LastPasswordChange, "LastPasswordChange should be set")
	})
}
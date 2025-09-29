package auth_integration_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/JoseLuis21/mv-backend/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTokenSecurity valida aspectos de seguridad de tokens
func TestTokenSecurity(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Token should be cryptographically secure", func(t *testing.T) {
		tokens := make(map[string]bool)

		// Generar múltiples usuarios y verificar unicidad de tokens
		for i := 0; i < 5; i++ {
			reqBody := helpers.CreateValidRegisterRequest()
			reqBody["email"] = fmt.Sprintf("security%d.test@misviaticos.cl", i)

			registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
			require.Equal(t, http.StatusCreated, registerResp.StatusCode)

			registerApiResp := helpers.ParseJSONResponse(t, registerResp)

			data, ok := registerApiResp.Data.(map[string]interface{})
			require.True(t, ok)

			emailToken := data["email_token"].(string)

			// Validar longitud y formato
			require.Len(t, emailToken, 64, "Token should be 64 characters")
			require.Regexp(t, "^[a-fA-F0-9]+$", emailToken, "Token should be hexadecimal")

			// Validar unicidad
			require.False(t, tokens[emailToken], "Token should be unique, got duplicate: %s", emailToken)
			tokens[emailToken] = true
		}

		t.Logf("✅ Generados %d tokens únicos y seguros", len(tokens))
	})

	t.Run("Invalid token formats should be rejected", func(t *testing.T) {
		invalidTokens := []string{
			"",                                // Empty
			"short",                           // Too short
			"toolongtoken1234567890123456789", // Too long
			"invalid-chars!@#$",               // Invalid characters
			"<script>alert('xss')</script>",   // XSS attempt
			"'; DROP TABLE users; --",         // SQL injection attempt
		}

		for _, invalidToken := range invalidTokens {
			verifyReq := map[string]interface{}{
				"token": invalidToken,
			}

			verifyResp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
			require.Equal(t, http.StatusBadRequest, verifyResp.StatusCode,
				"Invalid token should be rejected: %s", invalidToken)

			verifyApiResp := helpers.ParseJSONResponse(t, verifyResp)
			require.False(t, verifyApiResp.Success)

			t.Logf("✅ Token inválido rechazado correctamente: %s", invalidToken[:min(len(invalidToken), 20)])
		}
	})
}

// TestConcurrentVerification valida el comportamiento con verificaciones concurrentes
func TestConcurrentVerification(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Concurrent verification attempts should be handled gracefully", func(t *testing.T) {
		// Registrar usuario
		reqBody := helpers.CreateValidRegisterRequest()
		reqBody["email"] = "concurrent.test@misviaticos.cl"

		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		registerApiResp := helpers.ParseJSONResponse(t, registerResp)

		data, ok := registerApiResp.Data.(map[string]interface{})
		require.True(t, ok)
		emailToken := data["email_token"].(string)

		// Hacer múltiples verificaciones concurrentes
		successCount := 0
		errorCount := 0

		// Canal para recopilar resultados
		results := make(chan bool, 3)

		// Lanzar 3 verificaciones concurrentes
		for i := 0; i < 3; i++ {
			go func() {
				verifyReq := map[string]interface{}{
					"token": emailToken,
				}

				verifyResp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
				results <- verifyResp.StatusCode == http.StatusOK
			}()
		}

		// Recopilar resultados
		for i := 0; i < 3; i++ {
			if <-results {
				successCount++
			} else {
				errorCount++
			}
		}

		// Solo una debería ser exitosa
		assert.Equal(t, 1, successCount, "Only one verification should succeed")
		assert.Equal(t, 2, errorCount, "Two should fail due to race condition")

		t.Logf("✅ Verificación concurrente manejada: %d éxito, %d errores", successCount, errorCount)
	})
}

// TestEmailVerificationEdgeCases valida casos extremos
func TestEmailVerificationEdgeCases(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("User with already verified email", func(t *testing.T) {
		// Crear y verificar usuario
		reqBody := helpers.CreateValidRegisterRequest()
		reqBody["email"] = "alreadyverified.test@misviaticos.cl"

		// Registro
		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		registerApiResp := helpers.ParseJSONResponse(t, registerResp)

		data, ok := registerApiResp.Data.(map[string]interface{})
		require.True(t, ok)
		emailToken := data["email_token"].(string)

		// Primera verificación (exitosa)
		verifyReq := map[string]interface{}{
			"token": emailToken,
		}
		verifyResp1 := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		require.Equal(t, http.StatusOK, verifyResp1.StatusCode)

		// Segunda verificación (debería fallar elegantemente)
		verifyResp2 := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		require.Equal(t, http.StatusBadRequest, verifyResp2.StatusCode)

		verifyApiResp2 := helpers.ParseJSONResponse(t, verifyResp2)
		require.False(t, verifyApiResp2.Success)

		t.Log("✅ Usuario ya verificado manejado correctamente")
	})

	t.Run("Verify email for non-existent user token", func(t *testing.T) {
		// Token que no existe en BD
		fakeToken := "abcdefghijklmnopqrstuvwxyz123456"

		verifyReq := map[string]interface{}{
			"token": fakeToken,
		}

		verifyResp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		require.Equal(t, http.StatusBadRequest, verifyResp.StatusCode)

		verifyApiResp := helpers.ParseJSONResponse(t, verifyResp)
		require.False(t, verifyApiResp.Success)

		t.Log("✅ Token inexistente manejado correctamente")
	})
}

// TestPerformanceConstraints valida aspectos de rendimiento
func TestPerformanceConstraints(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Token generation should be fast", func(t *testing.T) {
		start := time.Now()

		reqBody := helpers.CreateValidRegisterRequest()
		reqBody["email"] = "performance.test@misviaticos.cl"

		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)

		elapsed := time.Since(start)

		require.Equal(t, http.StatusCreated, registerResp.StatusCode)
		assert.Less(t, elapsed, 2*time.Second, "Registration should complete within 2 seconds")

		t.Logf("✅ Registro completado en %v", elapsed)
	})

	t.Run("Verification should be fast", func(t *testing.T) {
		// Crear usuario primero
		reqBody := helpers.CreateValidRegisterRequest()
		reqBody["email"] = "speedtest.test@misviaticos.cl"

		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		registerApiResp := helpers.ParseJSONResponse(t, registerResp)

		data, ok := registerApiResp.Data.(map[string]interface{})
		require.True(t, ok)
		emailToken := data["email_token"].(string)

		// Verificar con medición de tiempo
		start := time.Now()

		verifyReq := map[string]interface{}{
			"token": emailToken,
		}
		verifyResp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)

		elapsed := time.Since(start)

		require.Equal(t, http.StatusOK, verifyResp.StatusCode)
		assert.Less(t, elapsed, 1*time.Second, "Verification should complete within 1 second")

		t.Logf("✅ Verificación completada en %v", elapsed)
	})
}

// TestTokenExpirationPrecision valida precisión de expiración
func TestTokenExpirationPrecision(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Token expiration should be precise", func(t *testing.T) {
		// Crear usuario
		reqBody := helpers.CreateValidRegisterRequest()
		reqBody["email"] = "precision.test@misviaticos.cl"

		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)

		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		registerApiResp := helpers.ParseJSONResponse(t, registerResp)

		data, ok := registerApiResp.Data.(map[string]interface{})
		require.True(t, ok)
		userID := data["id"].(string)

		// Verificar tiempo de expiración en BD
		ctx := context.Background()
		client := server.DBClient

		var tokenExpires time.Time
		err := client.QueryRow(ctx, "SELECT email_token_expires FROM users WHERE id = $1", userID).Scan(&tokenExpires)
		require.NoError(t, err)

		// El token debería expirar ~24 horas después del registro
		// Calcular desde el momento actual hacia el futuro
		now := time.Now().UTC()
		tokenExpiresUTC := tokenExpires.UTC()
		
		// El token debe expirar aproximadamente 24 horas DESDE EL MOMENTO DE CREACIÓN
		// Como acabamos de crear el usuario, debe expirar en aproximadamente 24 horas
		timeDiff := tokenExpiresUTC.Sub(now)
		
		// Debe estar entre 20h y 25h (amplia tolerancia para diferentes configuraciones)
		minExpectedDiff := 20*time.Hour
		maxExpectedDiff := 25*time.Hour

		assert.True(t, timeDiff >= minExpectedDiff,
			"Token should expire in at least %v, but expires in %v", minExpectedDiff, timeDiff)
		assert.True(t, timeDiff <= maxExpectedDiff,
			"Token should expire in at most %v, but expires in %v", maxExpectedDiff, timeDiff)

		t.Logf("✅ Token expira en %v (precisión validada)", tokenExpires)
	})
}

// Helper function min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

package auth_integration_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/JoseLuis21/mv-backend/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCompleteEmailVerificationFlow valida el flujo completo end-to-end
func TestCompleteEmailVerificationFlow(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Complete flow: Register -> Real Token -> Verify -> Database State", func(t *testing.T) {
		// PASO 1: REGISTRO - Crear usuario
		reqBody := helpers.CreateValidRegisterRequest()
		reqBody["email"] = "integration.test@misviaticos.cl" // Email único para test

		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		// Parsear respuesta de registro
		registerApiResp := helpers.ParseJSONResponse(t, registerResp)
		require.True(t, registerApiResp.Success)

		// Extraer datos del usuario registrado
		data, ok := registerApiResp.Data.(map[string]interface{})
		require.True(t, ok, "Response data should be a map")
		
		userID, ok := data["id"].(string)
		require.True(t, ok, "User ID should be present")
		require.NotEmpty(t, userID)

		// PASO 2: TOKEN - Validar que se generó token real
		emailToken, ok := data["email_token"].(string)
		require.True(t, ok, "Email token should be present")
		require.NotEmpty(t, emailToken, "Email token should not be empty")
		require.Len(t, emailToken, 64, "Email token should be 64 characters")
		
		requiresVerification, ok := data["requires_email_verification"].(bool)
		require.True(t, ok, "requires_email_verification should be present")
		require.True(t, requiresVerification, "Should require email verification")

		t.Logf("✅ PASO 1-2: Usuario registrado con token: %s", emailToken[:8]+"...")

		// PASO 3: VALIDAR ESTADO INICIAL EN BD
		// Verificar que el usuario está en estado no verificado
		ctx := context.Background()
		client := server.DBClient

		var emailVerified bool
		var storedToken *string
		var tokenExpires *time.Time

		err := client.QueryRow(ctx, 
			"SELECT email_verified, email_token, email_token_expires FROM users WHERE id = $1", 
			userID).Scan(&emailVerified, &storedToken, &tokenExpires)
		require.NoError(t, err)

		// Validaciones de estado inicial
		assert.False(t, emailVerified, "User should not be verified initially")
		assert.NotNil(t, storedToken, "Email token should be stored")
		assert.Equal(t, emailToken, *storedToken, "Stored token should match response token")
		assert.NotNil(t, tokenExpires, "Token expiration should be set")
		assert.True(t, tokenExpires.After(time.Now()), "Token should not be expired")
		assert.True(t, tokenExpires.Before(time.Now().Add(25*time.Hour)), "Token should expire within 24h")

		t.Log("✅ PASO 3: Estado inicial en BD validado")

		// PASO 4: VERIFICACIÓN - Usar token real para verificar
		verifyReq := map[string]interface{}{
			"token": emailToken, // Usar el token REAL generado
		}

		verifyResp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		require.Equal(t, http.StatusOK, verifyResp.StatusCode)

		verifyApiResp := helpers.ParseJSONResponse(t, verifyResp)
		require.True(t, verifyApiResp.Success)
		
		assert.Contains(t, verifyApiResp.Message, "verificado exitosamente", 
			"Verification success message should be present")

		t.Log("✅ PASO 4: Email verificado exitosamente")

		// PASO 5: VALIDAR ESTADO FINAL EN BD
		// Verificar que el usuario fue activado correctamente
		var finalEmailVerified bool
		var finalToken *string
		var finalTokenExpires *time.Time

		err = client.QueryRow(ctx, 
			"SELECT email_verified, email_token, email_token_expires FROM users WHERE id = $1", 
			userID).Scan(&finalEmailVerified, &finalToken, &finalTokenExpires)
		require.NoError(t, err)

		// Validaciones de estado final
		assert.True(t, finalEmailVerified, "User should be verified after verification")
		assert.Nil(t, finalToken, "Email token should be cleared after verification")
		assert.Nil(t, finalTokenExpires, "Token expiration should be cleared")

		t.Log("✅ PASO 5: Estado final en BD validado - Usuario activado correctamente")
	})
}

// TestTokenExpiration valida que tokens expirados no funcionen
func TestTokenExpiration(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Expired token should be rejected", func(t *testing.T) {
		// Crear usuario
		reqBody := helpers.CreateValidRegisterRequest()
		reqBody["email"] = "expired.test@misviaticos.cl"

		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		// Extraer datos
		registerApiResp := helpers.ParseJSONResponse(t, registerResp)

		data, ok := registerApiResp.Data.(map[string]interface{})
		require.True(t, ok)
		
		userID := data["id"].(string)
		emailToken := data["email_token"].(string)

		// SIMULAR EXPIRACIÓN: Modificar BD para hacer token expirado
		ctx := context.Background()
		client := server.DBClient

		// Establecer expiración en el pasado
		expiredTime := time.Now().Add(-1 * time.Hour)
		err := client.Exec(ctx,
			"UPDATE users SET email_token_expires = $1 WHERE id = $2",
			expiredTime, userID)
		require.NoError(t, err)

		t.Logf("✅ Token expirado simulado para usuario: %s", userID)

		// Intentar verificar con token expirado
		verifyReq := map[string]interface{}{
			"token": emailToken,
		}

		verifyResp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		require.Equal(t, http.StatusBadRequest, verifyResp.StatusCode)

		verifyApiResp := helpers.ParseJSONResponse(t, verifyResp)
		require.False(t, verifyApiResp.Success)
		
		assert.Contains(t, verifyApiResp.Error, "inválido", "Error should mention invalid token")

		t.Log("✅ Token expirado rechazado correctamente")
	})
}

// TestTokenReuse valida que tokens no se puedan reutilizar
func TestTokenReuse(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Used token should not work twice", func(t *testing.T) {
		// Crear y verificar usuario exitosamente
		reqBody := helpers.CreateValidRegisterRequest()
		reqBody["email"] = "reuse.test@misviaticos.cl"

		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		registerApiResp := helpers.ParseJSONResponse(t, registerResp)

		data, ok := registerApiResp.Data.(map[string]interface{})
		require.True(t, ok)
		emailToken := data["email_token"].(string)

		// Primera verificación (debería funcionar)
		verifyReq := map[string]interface{}{
			"token": emailToken,
		}

		verifyResp1 := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		require.Equal(t, http.StatusOK, verifyResp1.StatusCode)

		t.Log("✅ Primera verificación exitosa")

		// Segunda verificación con mismo token (debería fallar)
		verifyResp2 := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		require.Equal(t, http.StatusBadRequest, verifyResp2.StatusCode)

		verifyApiResp2 := helpers.ParseJSONResponse(t, verifyResp2)
		require.False(t, verifyApiResp2.Success)

		// El mensaje puede ser "email ya verificado" o "token inválido"
		errorMsg := verifyApiResp2.Error
		assert.True(t,
			contains(errorMsg, "verified") || contains(errorMsg, "inválido") || contains(errorMsg, "invalid"),
			"Should indicate token is invalid or email already verified, got: %s", errorMsg)

		t.Log("✅ Reutilización de token rechazada correctamente")
	})
}

// TestResendTokenInvalidatesPrevious valida que reenvío invalide token anterior
func TestResendTokenInvalidatesPrevious(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Resend should invalidate previous token", func(t *testing.T) {
		// Registrar usuario
		reqBody := helpers.CreateValidRegisterRequest()
		reqBody["email"] = "resend.test@misviaticos.cl"

		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		registerApiResp := helpers.ParseJSONResponse(t, registerResp)

		data, ok := registerApiResp.Data.(map[string]interface{})
		require.True(t, ok)
		originalToken := data["email_token"].(string)
		userEmail := data["email"].(string)

		t.Logf("✅ Token original: %s", originalToken[:8]+"...")

		// Reenviar verificación
		resendReq := map[string]interface{}{
			"email": userEmail,
		}

		resendResp := server.MakeRequest(t, "POST", "/api/v1/auth/resend-verification", resendReq)
		require.Equal(t, http.StatusOK, resendResp.StatusCode)

		t.Log("✅ Email de verificación reenviado")

		// Intentar usar token original (debería fallar)
		verifyReq := map[string]interface{}{
			"token": originalToken,
		}

		verifyResp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		require.Equal(t, http.StatusBadRequest, verifyResp.StatusCode)

		verifyApiResp := helpers.ParseJSONResponse(t, verifyResp)
		require.False(t, verifyApiResp.Success)

		t.Log("✅ Token original invalidado correctamente después del reenvío")
	})
}

// TestDatabaseConsistency valida consistencia de datos en BD
func TestDatabaseConsistency(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Database state should be consistent throughout flow", func(t *testing.T) {
		reqBody := helpers.CreateValidRegisterRequest()
		reqBody["email"] = "consistency.test@misviaticos.cl"

		// Registrar
		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		registerApiResp := helpers.ParseJSONResponse(t, registerResp)

		data, ok := registerApiResp.Data.(map[string]interface{})
		require.True(t, ok)
		userID := data["id"].(string)
		emailToken := data["email_token"].(string)

		// Usar la BD del servidor de testing
		ctx := context.Background()
		client := server.DBClient

		// Validar datos completos del usuario
		var fullName, email, phone string
		var emailVerified bool
		var createdAt, updatedAt time.Time
		var storedToken *string

		err := client.QueryRow(ctx, `
			SELECT full_name, email, phone, 
			       email_verified, email_token, created, updated 
			FROM users WHERE id = $1`, userID).Scan(
			&fullName, &email, &phone,
			&emailVerified, &storedToken, &createdAt, &updatedAt)
		require.NoError(t, err)

		// Validar datos del usuario
		expectedFullName := reqBody["firstname"].(string) + " " + reqBody["lastname"].(string)
		assert.Equal(t, expectedFullName, fullName)
		assert.Equal(t, reqBody["email"], email)
		assert.Equal(t, reqBody["phone"], phone)
		assert.False(t, emailVerified)
		assert.Equal(t, emailToken, *storedToken)
		assert.True(t, createdAt.Before(time.Now()))
		assert.True(t, updatedAt.Before(time.Now()))

		t.Log("✅ Datos del usuario consistentes en BD")

		// Verificar email
		verifyReq := map[string]interface{}{"token": emailToken}
		verifyResp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		require.Equal(t, http.StatusOK, verifyResp.StatusCode)

		// Validar actualización
		var finalEmailVerified bool
		var finalToken *string
		var finalUpdatedAt time.Time

		err = client.QueryRow(ctx, `
			SELECT email_verified, email_token, updated 
			FROM users WHERE id = $1`, userID).Scan(
			&finalEmailVerified, &finalToken, &finalUpdatedAt)
		require.NoError(t, err)

		assert.True(t, finalEmailVerified)
		assert.Nil(t, finalToken)
		assert.True(t, finalUpdatedAt.After(updatedAt), "updated_at should be newer after verification")

		t.Log("✅ Estado final de BD consistente después de verificación")
	})
}

// Helper function para verificar si string contiene substring
func contains(str, substr string) bool {
	if len(substr) > len(str) {
		return false
	}
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
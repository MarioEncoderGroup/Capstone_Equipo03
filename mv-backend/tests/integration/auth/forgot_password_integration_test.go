package auth_integration_test

import (
	"net/http"
	"testing"

	"github.com/JoseLuis21/mv-backend/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCompleteForgotPasswordFlow valida el flujo completo end-to-end
func TestCompleteForgotPasswordFlow(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Complete flow: Register -> Verify -> ForgotPassword -> ResetPassword", func(t *testing.T) {
		// PASO 1: REGISTRO - Crear usuario verificado
		registerReq := helpers.CreateValidRegisterRequest()
		registerReq["email"] = "forgot.test@misviaticos.cl" // Email único para test

		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", registerReq)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		// Extraer token de verificación
		registerApiResp := helpers.ParseJSONResponse(t, registerResp)
		require.True(t, registerApiResp.Success)

		data, ok := registerApiResp.Data.(map[string]interface{})
		require.True(t, ok, "Response data should be a map")

		emailToken, ok := data["email_token"].(string)
		require.True(t, ok, "Email token should be present")
		require.NotEmpty(t, emailToken)

		// PASO 2: VERIFICACIÓN - Verificar email del usuario
		verifyReq := map[string]interface{}{
			"token": emailToken,
		}

		verifyResp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		require.Equal(t, http.StatusOK, verifyResp.StatusCode)

		verifyApiResp := helpers.ParseJSONResponse(t, verifyResp)
		require.True(t, verifyApiResp.Success)

		t.Logf("✅ PASO 1-2: Usuario registrado y verificado")

		// PASO 3: FORGOT PASSWORD - Solicitar reset de contraseña
		forgotReq := map[string]interface{}{
			"email": registerReq["email"],
		}

		forgotResp := server.MakeRequest(t, "POST", "/api/v1/auth/forgot-password", forgotReq)
		require.Equal(t, http.StatusOK, forgotResp.StatusCode)

		forgotApiResp := helpers.ParseJSONResponse(t, forgotResp)
		require.True(t, forgotApiResp.Success)
		assert.Contains(t, forgotApiResp.Message, "recibirás instrucciones",
			"Should return generic message for security")

		t.Logf("✅ PASO 3: Forgot password request processed")

		// PASO 4: OBTENER TOKEN - Simular obtención del token desde email
		// En un test real, buscaríamos en la BD o interceptaríamos el email
		// Por ahora usamos el helper que puede extraer el token de la BD
		resetToken := helpers.GetLatestPasswordResetToken(t, server.DBClient, registerReq["email"].(string))
		require.NotEmpty(t, resetToken, "Password reset token should be generated")
		require.Len(t, resetToken, 64, "Reset token should be 64 characters")

		t.Logf("✅ PASO 4: Reset token obtained: %s", resetToken[:8]+"...")

		// PASO 5: RESET PASSWORD - Cambiar contraseña usando token
		newPassword := "NewSecurePassword123!"
		resetReq := map[string]interface{}{
			"token":       resetToken,
			"newPassword": newPassword,
		}

		resetResp := server.MakeRequest(t, "POST", "/api/v1/auth/reset-password", resetReq)
		require.Equal(t, http.StatusOK, resetResp.StatusCode)

		resetApiResp := helpers.ParseJSONResponse(t, resetResp)
		require.True(t, resetApiResp.Success)
		assert.Equal(t, "Contraseña actualizada exitosamente", resetApiResp.Message)

		t.Logf("✅ PASO 5: Password reset successful")

		// PASO 6: VERIFICACIÓN - Login con nueva contraseña
		loginReq := map[string]interface{}{
			"email":    registerReq["email"],
			"password": newPassword,
		}

		loginResp := server.MakeRequest(t, "POST", "/api/v1/auth/login", loginReq)
		require.Equal(t, http.StatusOK, loginResp.StatusCode)

		loginApiResp := helpers.ParseJSONResponse(t, loginResp)
		require.True(t, loginApiResp.Success)
		assert.Equal(t, "Autenticación exitosa", loginApiResp.Message)

		t.Logf("✅ PASO 6: Login with new password successful")

		// PASO 7: VERIFICACIÓN - El token de reset debe estar invalidado
		// Intentar usar el mismo token otra vez debe fallar
		resetReq2 := map[string]interface{}{
			"token":       resetToken,
			"newPassword": "AnotherPassword123!",
		}

		resetResp2 := server.MakeRequest(t, "POST", "/api/v1/auth/reset-password", resetReq2)
		require.Equal(t, http.StatusBadRequest, resetResp2.StatusCode)

		resetApiResp2 := helpers.ParseJSONResponse(t, resetResp2)
		require.False(t, resetApiResp2.Success)
		assert.Contains(t, resetApiResp2.Error, "inválido", "Used token should be invalid")

		t.Logf("✅ PASO 7: Used token correctly invalidated")
	})
}

func TestForgotPasswordSecurityScenarios(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Should handle non-existent email gracefully", func(t *testing.T) {
		forgotReq := map[string]interface{}{
			"email": "nonexistent@example.com",
		}

		forgotResp := server.MakeRequest(t, "POST", "/api/v1/auth/forgot-password", forgotReq)
		require.Equal(t, http.StatusOK, forgotResp.StatusCode) // Mismo código por seguridad

		forgotApiResp := helpers.ParseJSONResponse(t, forgotResp)
		require.True(t, forgotApiResp.Success) // Mismo mensaje por seguridad
		assert.Contains(t, forgotApiResp.Message, "recibirás instrucciones",
			"Should return same message for non-existent email")
	})

	t.Run("Should reject forgot password for unverified email", func(t *testing.T) {
		// Crear usuario sin verificar
		registerReq := helpers.CreateValidRegisterRequest()
		registerReq["email"] = "unverified@test.cl"

		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", registerReq)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		// Intentar forgot password sin verificar email
		forgotReq := map[string]interface{}{
			"email": registerReq["email"],
		}

		forgotResp := server.MakeRequest(t, "POST", "/api/v1/auth/forgot-password", forgotReq)
		require.Equal(t, http.StatusBadRequest, forgotResp.StatusCode)

		forgotApiResp := helpers.ParseJSONResponse(t, forgotResp)
		require.False(t, forgotApiResp.Success)
		assert.Contains(t, forgotApiResp.Error, "VALIDATION_ERROR", "Should reject unverified email")
	})

	t.Run("Should prevent multiple active reset tokens", func(t *testing.T) {
		// Crear y verificar usuario
		registerReq := helpers.CreateValidRegisterRequest()
		registerReq["email"] = "multiple.tokens@test.cl"

		// Registro y verificación (pasos abreviados)
		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", registerReq)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		registerApiResp := helpers.ParseJSONResponse(t, registerResp)
		data := registerApiResp.Data.(map[string]interface{})
		emailToken := data["email_token"].(string)

		verifyReq := map[string]interface{}{"token": emailToken}
		verifyResp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		require.Equal(t, http.StatusOK, verifyResp.StatusCode)

		// Primera solicitud de reset
		forgotReq := map[string]interface{}{"email": registerReq["email"]}
		forgotResp1 := server.MakeRequest(t, "POST", "/api/v1/auth/forgot-password", forgotReq)
		require.Equal(t, http.StatusOK, forgotResp1.StatusCode)

		// Intentar segunda solicitud inmediatamente
		forgotResp2 := server.MakeRequest(t, "POST", "/api/v1/auth/forgot-password", forgotReq)
		require.Equal(t, http.StatusBadRequest, forgotResp2.StatusCode)

		forgotApiResp2 := helpers.ParseJSONResponse(t, forgotResp2)
		require.False(t, forgotApiResp2.Success)
		assert.Contains(t, forgotApiResp2.Error, "VALIDATION_ERROR",
			"Should prevent multiple active tokens")
	})
}

func TestResetPasswordValidation(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Should validate reset password input", func(t *testing.T) {
		testCases := []struct {
			name           string
			token          string
			newPassword    string
			expectedStatus int
			expectedError  string
		}{
			{
				name:           "Empty token",
				token:          "",
				newPassword:    "ValidPassword123!",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "Errores de validación",
			},
			{
				name:           "Short password",
				token:          "valid-token-format-but-fake",
				newPassword:    "short",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "validación",
			},
			{
				name:           "Invalid token format",
				token:          "invalid-token",
				newPassword:    "ValidPassword123!",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "Error reseteando contraseña",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				resetReq := map[string]interface{}{
					"token":       tc.token,
					"newPassword": tc.newPassword,
				}

				resetResp := server.MakeRequest(t, "POST", "/api/v1/auth/reset-password", resetReq)
				assert.Equal(t, tc.expectedStatus, resetResp.StatusCode)

				resetApiResp := helpers.ParseJSONResponse(t, resetResp)
				assert.False(t, resetApiResp.Success)
				assert.Contains(t, resetApiResp.Message, tc.expectedError)
			})
		}
	})

	t.Run("Should reject expired reset token", func(t *testing.T) {
		// Este test requiere manipulación de tiempo o BD para simular expiración
		// Por ahora verificamos que tokens inválidos son rechazados
		resetReq := map[string]interface{}{
			"token":       "expired-or-invalid-token-12345678901234567890123456789012345678901234567890123456",
			"newPassword": "ValidPassword123!",
		}

		resetResp := server.MakeRequest(t, "POST", "/api/v1/auth/reset-password", resetReq)
		require.Equal(t, http.StatusBadRequest, resetResp.StatusCode)

		resetApiResp := helpers.ParseJSONResponse(t, resetResp)
		require.False(t, resetApiResp.Success)
		assert.Contains(t, resetApiResp.Message, "Error reseteando contraseña", "Should reject invalid/expired token")
	})
}

func TestForgotPasswordInputValidation(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Should validate forgot password input", func(t *testing.T) {
		testCases := []struct {
			name           string
			email          string
			expectedStatus int
			expectedError  string
		}{
			{
				name:           "Empty email",
				email:          "",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "Email inválido",
			},
			{
				name:           "Invalid email format",
				email:          "not-an-email",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "Email inválido",
			},
			{
				name:           "Valid email format",
				email:          "valid@example.com",
				expectedStatus: http.StatusOK,
				expectedError:  "", // No error expected
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				forgotReq := map[string]interface{}{
					"email": tc.email,
				}

				forgotResp := server.MakeRequest(t, "POST", "/api/v1/auth/forgot-password", forgotReq)
				assert.Equal(t, tc.expectedStatus, forgotResp.StatusCode)

				forgotApiResp := helpers.ParseJSONResponse(t, forgotResp)

				if tc.expectedError != "" {
					assert.False(t, forgotApiResp.Success)
					assert.Contains(t, forgotApiResp.Message, tc.expectedError)
				} else {
					assert.True(t, forgotApiResp.Success)
				}
			})
		}
	})
}

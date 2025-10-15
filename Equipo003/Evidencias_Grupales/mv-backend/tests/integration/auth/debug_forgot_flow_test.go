package auth_integration_test

import (
	"net/http"
	"testing"

	"github.com/JoseLuis21/mv-backend/tests/helpers"
	"github.com/stretchr/testify/require"
)

// TestDebugForgotFlow test simple para debuggear el flujo completo paso a paso
func TestDebugForgotFlow(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Step by step forgot password debug", func(t *testing.T) {
		// PASO 1: Registro
		registerReq := helpers.CreateValidRegisterRequest()
		registerReq["email"] = "debug.forgot@test.cl"

		t.Logf("🔵 STEP 1: Registering user with email: %s", registerReq["email"])
		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", registerReq)
		t.Logf("Register response status: %d", registerResp.StatusCode)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		registerApiResp := helpers.ParseJSONResponse(t, registerResp)
		require.True(t, registerApiResp.Success)

		data := registerApiResp.Data.(map[string]interface{})
		emailToken := data["email_token"].(string)
		t.Logf("Email verification token: %s", emailToken[:8]+"...")

		// PASO 2: Verificar email
		t.Logf("🔵 STEP 2: Verifying email")
		verifyReq := map[string]interface{}{"token": emailToken}
		verifyResp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		t.Logf("Verify response status: %d", verifyResp.StatusCode)
		require.Equal(t, http.StatusOK, verifyResp.StatusCode)

		verifyApiResp := helpers.ParseJSONResponse(t, verifyResp)
		require.True(t, verifyApiResp.Success)
		t.Logf("✅ Email verified successfully")

		// PASO 3: Forgot password
		t.Logf("🔵 STEP 3: Requesting password reset")
		forgotReq := map[string]interface{}{"email": registerReq["email"]}
		forgotResp := server.MakeRequest(t, "POST", "/api/v1/auth/forgot-password", forgotReq)
		t.Logf("Forgot password response status: %d", forgotResp.StatusCode)

		forgotApiResp := helpers.ParseJSONResponse(t, forgotResp)
		t.Logf("Forgot password response: %+v", forgotApiResp)

		require.Equal(t, http.StatusOK, forgotResp.StatusCode)
		require.True(t, forgotApiResp.Success)
		t.Logf("✅ Forgot password request successful")

		// PASO 4: Intentar obtener token de la BD
		t.Logf("🔵 STEP 4: Attempting to get reset token from database")
		resetToken := helpers.GetLatestPasswordResetToken(t, server.DBClient, registerReq["email"].(string))
		t.Logf("Reset token from DB: '%s' (length: %d)", resetToken, len(resetToken))

		if resetToken == "" {
			t.Log("❌ No reset token found in database")
			t.Log("Let's check if the user exists and what data it has...")

			// Debug: Check user in database
			t.Logf("🔍 DEBUG: Checking user in database for email: %s", registerReq["email"])
		} else {
			t.Logf("✅ Reset token found: %s", resetToken[:8]+"...")

			// PASO 5: Intentar reset de contraseña
			t.Logf("🔵 STEP 5: Attempting password reset")
			resetReq := map[string]interface{}{
				"token":       resetToken,
				"newPassword": "NewPassword123!",
			}

			resetResp := server.MakeRequest(t, "POST", "/api/v1/auth/reset-password", resetReq)
			t.Logf("Reset password response status: %d", resetResp.StatusCode)

			resetApiResp := helpers.ParseJSONResponse(t, resetResp)
			t.Logf("Reset password response: %+v", resetApiResp)
		}
	})
}
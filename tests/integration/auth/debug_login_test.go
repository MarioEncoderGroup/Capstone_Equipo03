package auth_integration_test

import (
	"net/http"
	"testing"

	"github.com/JoseLuis21/mv-backend/tests/helpers"
	"github.com/stretchr/testify/require"
)

// TestDebugLoginFlow test simple para debuggear el flujo de login paso a paso
func TestDebugLoginFlow(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Step by step login debug", func(t *testing.T) {
		// PASO 1: Registro
		registerReq := helpers.CreateValidRegisterRequest()
		registerReq["email"] = "debug.login@test.cl"

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

		// PASO 3: Login
		t.Logf("🔵 STEP 3: Attempting login")
		loginReq := map[string]interface{}{
			"email":    registerReq["email"],
			"password": registerReq["password"],
		}

		loginResp := server.MakeRequest(t, "POST", "/api/v1/auth/login", loginReq)
		t.Logf("Login response status: %d", loginResp.StatusCode)

		loginApiResp := helpers.ParseJSONResponse(t, loginResp)
		t.Logf("Login response: %+v", loginApiResp)

		require.Equal(t, http.StatusOK, loginResp.StatusCode)
		require.True(t, loginApiResp.Success)
		t.Logf("✅ Login successful")

		// PASO 4: Examinar tokens
		if loginApiResp.Data != nil {
			loginData := loginApiResp.Data.(map[string]interface{})

			if accessToken, ok := loginData["access_token"].(string); ok {
				t.Logf("🔍 Access token: %s... (length: %d)", accessToken[:20], len(accessToken))
			}

			if refreshToken, ok := loginData["refresh_token"].(string); ok {
				t.Logf("🔍 Refresh token: %s... (length: %d)", refreshToken[:20], len(refreshToken))
			}

			if expiresIn, ok := loginData["expires_in"].(float64); ok {
				t.Logf("🔍 Expires in: %f seconds (%.1f hours)", expiresIn, expiresIn/3600)
			}

			if tokenType, ok := loginData["token_type"].(string); ok {
				t.Logf("🔍 Token type: %s", tokenType)
			}

			if userData, ok := loginData["user"].(map[string]interface{}); ok {
				t.Logf("🔍 User data:")
				if email, ok := userData["email"].(string); ok {
					t.Logf("   - Email: %s", email)
				}
				if fullName, ok := userData["full_name"].(string); ok {
					t.Logf("   - Full name: %s", fullName)
				}
				if isActive, ok := userData["is_active"].(bool); ok {
					t.Logf("   - Is active: %t", isActive)
				}
				if lastLogin, ok := userData["last_login"]; ok {
					t.Logf("   - Last login: %v", lastLogin)
				}
			}

			t.Logf("✅ PASO 4: Token data examined successfully")
		}

		// PASO 5: Test de login con credenciales incorrectas
		t.Logf("🔵 STEP 5: Testing invalid credentials")
		invalidLoginReq := map[string]interface{}{
			"email":    registerReq["email"],
			"password": "wrongpassword",
		}

		invalidLoginResp := server.MakeRequest(t, "POST", "/api/v1/auth/login", invalidLoginReq)
		t.Logf("Invalid login response status: %d", invalidLoginResp.StatusCode)

		invalidLoginApiResp := helpers.ParseJSONResponse(t, invalidLoginResp)
		t.Logf("Invalid login response: %+v", invalidLoginApiResp)

		require.Equal(t, http.StatusUnauthorized, invalidLoginResp.StatusCode)
		require.False(t, invalidLoginApiResp.Success)
		t.Logf("✅ PASO 5: Invalid credentials correctly rejected")

		// PASO 6: Test segundo login para verificar actualización de last_login
		t.Logf("🔵 STEP 6: Testing second login to verify last_login update")
		secondLoginResp := server.MakeRequest(t, "POST", "/api/v1/auth/login", loginReq)
		t.Logf("Second login response status: %d", secondLoginResp.StatusCode)

		secondLoginApiResp := helpers.ParseJSONResponse(t, secondLoginResp)
		require.Equal(t, http.StatusOK, secondLoginResp.StatusCode)
		require.True(t, secondLoginApiResp.Success)

		if secondLoginApiResp.Data != nil {
			secondLoginData := secondLoginApiResp.Data.(map[string]interface{})
			if userData, ok := secondLoginData["user"].(map[string]interface{}); ok {
				if lastLogin, ok := userData["last_login"]; ok {
					t.Logf("🔍 Second login - Last login updated: %v", lastLogin)
				}
			}
		}

		t.Logf("✅ PASO 6: Second login successful with updated last_login")
	})
}

// TestDebugLoginValidation test simple para validar diferentes escenarios de validación
func TestDebugLoginValidation(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Debug login validation scenarios", func(t *testing.T) {
		// Test 1: Email vacío
		t.Logf("🔵 TEST 1: Empty email")
		emptyEmailReq := map[string]interface{}{
			"email":    "",
			"password": "password123",
		}

		emptyEmailResp := server.MakeRequest(t, "POST", "/api/v1/auth/login", emptyEmailReq)
		t.Logf("Empty email response status: %d", emptyEmailResp.StatusCode)

		emptyEmailApiResp := helpers.ParseJSONResponse(t, emptyEmailResp)
		t.Logf("Empty email response: %+v", emptyEmailApiResp)

		require.Equal(t, http.StatusBadRequest, emptyEmailResp.StatusCode)
		require.False(t, emptyEmailApiResp.Success)
		t.Logf("✅ Empty email correctly rejected")

		// Test 2: Formato de email inválido
		t.Logf("🔵 TEST 2: Invalid email format")
		invalidEmailReq := map[string]interface{}{
			"email":    "not-an-email",
			"password": "password123",
		}

		invalidEmailResp := server.MakeRequest(t, "POST", "/api/v1/auth/login", invalidEmailReq)
		t.Logf("Invalid email response status: %d", invalidEmailResp.StatusCode)

		invalidEmailApiResp := helpers.ParseJSONResponse(t, invalidEmailResp)
		t.Logf("Invalid email response: %+v", invalidEmailApiResp)

		require.Equal(t, http.StatusBadRequest, invalidEmailResp.StatusCode)
		require.False(t, invalidEmailApiResp.Success)
		t.Logf("✅ Invalid email format correctly rejected")

		// Test 3: Contraseña vacía
		t.Logf("🔵 TEST 3: Empty password")
		emptyPasswordReq := map[string]interface{}{
			"email":    "test@example.com",
			"password": "",
		}

		emptyPasswordResp := server.MakeRequest(t, "POST", "/api/v1/auth/login", emptyPasswordReq)
		t.Logf("Empty password response status: %d", emptyPasswordResp.StatusCode)

		emptyPasswordApiResp := helpers.ParseJSONResponse(t, emptyPasswordResp)
		t.Logf("Empty password response: %+v", emptyPasswordApiResp)

		require.Equal(t, http.StatusBadRequest, emptyPasswordResp.StatusCode)
		require.False(t, emptyPasswordApiResp.Success)
		t.Logf("✅ Empty password correctly rejected")

		// Test 4: Usuario inexistente
		t.Logf("🔵 TEST 4: Non-existent user")
		nonExistentReq := map[string]interface{}{
			"email":    "nonexistent@example.com",
			"password": "password123",
		}

		nonExistentResp := server.MakeRequest(t, "POST", "/api/v1/auth/login", nonExistentReq)
		t.Logf("Non-existent user response status: %d", nonExistentResp.StatusCode)

		nonExistentApiResp := helpers.ParseJSONResponse(t, nonExistentResp)
		t.Logf("Non-existent user response: %+v", nonExistentApiResp)

		require.Equal(t, http.StatusUnauthorized, nonExistentResp.StatusCode)
		require.False(t, nonExistentApiResp.Success)
		t.Logf("✅ Non-existent user correctly rejected")
	})
}
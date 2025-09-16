package auth_integration_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/JoseLuis21/mv-backend/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCompleteLoginFlow valida el flujo completo de login
func TestCompleteLoginFlow(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Complete flow: Register -> Verify -> Login", func(t *testing.T) {
		// PASO 1: REGISTRO - Crear y verificar usuario
		registerReq := helpers.CreateValidRegisterRequest()
		registerReq["email"] = "login.test@misviaticos.cl"

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

		// PASO 2: VERIFICACIÓN - Verificar email
		verifyReq := map[string]interface{}{
			"token": emailToken,
		}

		verifyResp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		require.Equal(t, http.StatusOK, verifyResp.StatusCode)

		verifyApiResp := helpers.ParseJSONResponse(t, verifyResp)
		require.True(t, verifyApiResp.Success)

		t.Logf("✅ PASO 1-2: Usuario registrado y verificado")

		// PASO 3: LOGIN - Autenticar usuario
		loginReq := map[string]interface{}{
			"email":    registerReq["email"],
			"password": registerReq["password"],
		}

		loginResp := server.MakeRequest(t, "POST", "/api/v1/auth/login", loginReq)
		require.Equal(t, http.StatusOK, loginResp.StatusCode)

		loginApiResp := helpers.ParseJSONResponse(t, loginResp)
		require.True(t, loginApiResp.Success)
		assert.Equal(t, "Autenticación exitosa", loginApiResp.Message)

		// Verificar estructura de respuesta
		loginData, ok := loginApiResp.Data.(map[string]interface{})
		require.True(t, ok, "Login response data should be a map")

		// Verificar tokens
		assert.Contains(t, loginData, "access_token", "Should contain access token")
		assert.Contains(t, loginData, "refresh_token", "Should contain refresh token")
		assert.Contains(t, loginData, "expires_in", "Should contain expiration time")
		assert.Contains(t, loginData, "token_type", "Should contain token type")

		accessToken, ok := loginData["access_token"].(string)
		require.True(t, ok, "Access token should be string")
		require.NotEmpty(t, accessToken, "Access token should not be empty")

		refreshToken, ok := loginData["refresh_token"].(string)
		require.True(t, ok, "Refresh token should be string")
		require.NotEmpty(t, refreshToken, "Refresh token should not be empty")

		tokenType, ok := loginData["token_type"].(string)
		require.True(t, ok, "Token type should be string")
		assert.Equal(t, "Bearer", tokenType, "Token type should be Bearer")

		expiresIn, ok := loginData["expires_in"].(float64)
		require.True(t, ok, "Expires in should be number")
		assert.Equal(t, float64(24*60*60), expiresIn, "Should expire in 24 hours (86400 seconds)")

		// Verificar datos del usuario
		assert.Contains(t, loginData, "user", "Should contain user data")
		userData, ok := loginData["user"].(map[string]interface{})
		require.True(t, ok, "User data should be a map")

		assert.Equal(t, registerReq["email"], userData["email"], "Email should match")
		assert.Contains(t, userData, "full_name", "Should contain full_name")
		assert.Contains(t, userData, "last_login", "Should contain last_login")
		assert.Equal(t, true, userData["is_active"], "User should be active")

		t.Logf("✅ PASO 3: Login exitoso con tokens válidos")

		// PASO 4: VERIFICAR ÚLTIMO LOGIN - Hacer login nuevamente para verificar que se actualiza
		loginResp2 := server.MakeRequest(t, "POST", "/api/v1/auth/login", loginReq)
		require.Equal(t, http.StatusOK, loginResp2.StatusCode)

		loginApiResp2 := helpers.ParseJSONResponse(t, loginResp2)
		require.True(t, loginApiResp2.Success)

		loginData2, ok := loginApiResp2.Data.(map[string]interface{})
		require.True(t, ok)

		userData2, ok := loginData2["user"].(map[string]interface{})
		require.True(t, ok)

		// Verificar que last_login se actualizó
		assert.NotNil(t, userData2["last_login"], "Last login should be updated")
		t.Logf("✅ PASO 4: Último login actualizado correctamente")
	})
}

func TestLoginSecurityScenarios(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Should reject invalid credentials", func(t *testing.T) {
		loginReq := map[string]interface{}{
			"email":    "nonexistent@test.cl",
			"password": "wrongpassword",
		}

		loginResp := server.MakeRequest(t, "POST", "/api/v1/auth/login", loginReq)
		require.Equal(t, http.StatusUnauthorized, loginResp.StatusCode)

		loginApiResp := helpers.ParseJSONResponse(t, loginResp)
		require.False(t, loginApiResp.Success)
		assert.Equal(t, "INVALID_CREDENTIALS", loginApiResp.Error)
		assert.Equal(t, "Credenciales inválidas", loginApiResp.Message)
	})

	t.Run("Should reject login for unverified email", func(t *testing.T) {
		// Crear usuario sin verificar
		registerReq := helpers.CreateValidRegisterRequest()
		registerReq["email"] = "unverified.login@test.cl"

		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", registerReq)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		// Intentar login sin verificar email
		loginReq := map[string]interface{}{
			"email":    registerReq["email"],
			"password": registerReq["password"],
		}

		loginResp := server.MakeRequest(t, "POST", "/api/v1/auth/login", loginReq)
		require.Equal(t, http.StatusForbidden, loginResp.StatusCode)

		loginApiResp := helpers.ParseJSONResponse(t, loginResp)
		require.False(t, loginApiResp.Success)
		assert.Equal(t, "EMAIL_NOT_VERIFIED", loginApiResp.Error)
		assert.Contains(t, loginApiResp.Message, "no ha sido verificado", "Should indicate email not verified")
	})

	t.Run("Should reject login with wrong password", func(t *testing.T) {
		// Crear y verificar usuario
		registerReq := helpers.CreateValidRegisterRequest()
		registerReq["email"] = "wrongpass.test@test.cl"

		// Registro
		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", registerReq)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		registerApiResp := helpers.ParseJSONResponse(t, registerResp)
		data := registerApiResp.Data.(map[string]interface{})
		emailToken := data["email_token"].(string)

		// Verificación
		verifyReq := map[string]interface{}{"token": emailToken}
		verifyResp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		require.Equal(t, http.StatusOK, verifyResp.StatusCode)

		// Intentar login con contraseña incorrecta
		loginReq := map[string]interface{}{
			"email":    registerReq["email"],
			"password": "wrongpassword123",
		}

		loginResp := server.MakeRequest(t, "POST", "/api/v1/auth/login", loginReq)
		require.Equal(t, http.StatusUnauthorized, loginResp.StatusCode)

		loginApiResp := helpers.ParseJSONResponse(t, loginResp)
		require.False(t, loginApiResp.Success)
		assert.Equal(t, "INVALID_CREDENTIALS", loginApiResp.Error)
		assert.Equal(t, "Credenciales inválidas", loginApiResp.Message)
	})
}

func TestLoginInputValidation(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Should validate login input", func(t *testing.T) {
		testCases := []struct {
			name           string
			email          string
			password       string
			expectedStatus int
			expectedError  string
		}{
			{
				name:           "Empty email",
				email:          "",
				password:       "password123",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "VALIDATION_ERROR",
			},
			{
				name:           "Invalid email format",
				email:          "not-an-email",
				password:       "password123",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "VALIDATION_ERROR",
			},
			{
				name:           "Empty password",
				email:          "valid@test.cl",
				password:       "",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "VALIDATION_ERROR",
			},
			{
				name:           "Valid format but nonexistent user",
				email:          "valid@test.cl",
				password:       "password123",
				expectedStatus: http.StatusUnauthorized,
				expectedError:  "INVALID_CREDENTIALS",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				loginReq := map[string]interface{}{
					"email":    tc.email,
					"password": tc.password,
				}

				loginResp := server.MakeRequest(t, "POST", "/api/v1/auth/login", loginReq)
				assert.Equal(t, tc.expectedStatus, loginResp.StatusCode)

				loginApiResp := helpers.ParseJSONResponse(t, loginResp)
				assert.False(t, loginApiResp.Success)
				assert.Equal(t, tc.expectedError, loginApiResp.Error)
			})
		}
	})

	t.Run("Should handle malformed JSON", func(t *testing.T) {
		// Enviar JSON malformado
		req := `{"email": "test@test.cl", "password":}`
		resp := server.MakeRequest(t, "POST", "/api/v1/auth/login", req)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		apiResp := helpers.ParseJSONResponse(t, resp)
		require.False(t, apiResp.Success)
		assert.Equal(t, "INVALID_REQUEST_FORMAT", apiResp.Error)
		assert.Contains(t, strings.ToLower(apiResp.Message), "formato", "Should indicate format error")
	})
}

func TestLoginResponseStructure(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Login response should have correct structure", func(t *testing.T) {
		// Crear y verificar usuario
		registerReq := helpers.CreateValidRegisterRequest()
		registerReq["email"] = "structure.test@test.cl"

		// Registro
		registerResp := server.MakeRequest(t, "POST", "/api/v1/auth/register", registerReq)
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		registerApiResp := helpers.ParseJSONResponse(t, registerResp)
		data := registerApiResp.Data.(map[string]interface{})
		emailToken := data["email_token"].(string)

		// Verificación
		verifyReq := map[string]interface{}{"token": emailToken}
		verifyResp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		require.Equal(t, http.StatusOK, verifyResp.StatusCode)

		// Login
		loginReq := map[string]interface{}{
			"email":    registerReq["email"],
			"password": registerReq["password"],
		}

		loginResp := server.MakeRequest(t, "POST", "/api/v1/auth/login", loginReq)
		require.Equal(t, http.StatusOK, loginResp.StatusCode)

		loginApiResp := helpers.ParseJSONResponse(t, loginResp)
		require.True(t, loginApiResp.Success)

		// Verificar estructura completa
		helpers.AssertContainsField(t, loginApiResp, "access_token")
		helpers.AssertContainsField(t, loginApiResp, "refresh_token")
		helpers.AssertContainsField(t, loginApiResp, "expires_in")
		helpers.AssertContainsField(t, loginApiResp, "token_type")
		helpers.AssertContainsField(t, loginApiResp, "user")

		// Verificar que los tokens no estén vacíos
		loginData := loginApiResp.Data.(map[string]interface{})

		accessToken := loginData["access_token"].(string)
		assert.True(t, len(accessToken) > 50, "Access token should be substantial length")

		refreshToken := loginData["refresh_token"].(string)
		assert.True(t, len(refreshToken) > 50, "Refresh token should be substantial length")

		// Verificar que los tokens sean diferentes
		assert.NotEqual(t, accessToken, refreshToken, "Access and refresh tokens should be different")

		// Verificar estructura del usuario
		userData := loginData["user"].(map[string]interface{})
		assert.Contains(t, userData, "id", "User should have ID")
		assert.Contains(t, userData, "email", "User should have email")
		assert.Contains(t, userData, "full_name", "User should have full_name")
		assert.Contains(t, userData, "is_active", "User should have is_active")
		assert.Contains(t, userData, "last_login", "User should have last_login")

		t.Logf("✅ Login response structure validated")
	})
}
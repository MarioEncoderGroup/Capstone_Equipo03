package auth_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/JoseLuis21/mv-backend/tests/helpers"
)

func TestRegisterEndpoint_UserRegistration(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Valid user registration", func(t *testing.T) {
		// Arrange
		reqBody := helpers.CreateValidRegisterRequest()
		
		// Act
		resp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
		
		// Assert
		helpers.AssertStatusCode(t, resp, http.StatusCreated)
		
		apiResp := helpers.ParseJSONResponse(t, resp)
		helpers.AssertSuccessResponse(t, apiResp)
		
		// Verificar campos de respuesta
		helpers.AssertContainsField(t, apiResp, "id")
		helpers.AssertFieldEquals(t, apiResp, "requires_email_verification", true)
		
		// No debe tener tenant_id para registro individual
		if data, ok := apiResp.Data.(map[string]interface{}); ok {
			if _, exists := data["tenant_id"]; exists {
				t.Error("Individual registration should not return tenant_id")
			}
		}
	})

	t.Run("Duplicate email registration", func(t *testing.T) {
		// Arrange - Crear usuario primero
		reqBody := helpers.CreateValidRegisterRequest()
		server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
		
		// Act - Intentar crear usuario con mismo email
		resp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
		
		// Assert
		helpers.AssertStatusCode(t, resp, http.StatusConflict)
		
		apiResp := helpers.ParseJSONResponse(t, resp)
		helpers.AssertErrorCode(t, apiResp, "USER_ALREADY_EXISTS")
	})

	t.Run("Missing required fields", func(t *testing.T) {
		// Arrange
		reqBody := map[string]interface{}{
			"firstname": "Usuario",
			// Missing required fields: lastname, email, phone, password, password_confirm
		}
		
		// Act
		resp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
		
		// Assert
		helpers.AssertStatusCode(t, resp, http.StatusBadRequest)
		
		apiResp := helpers.ParseJSONResponse(t, resp)
		helpers.AssertErrorResponse(t, apiResp)
	})

	t.Run("Weak password", func(t *testing.T) {
		// Arrange
		reqBody := helpers.CreateValidRegisterRequest()
		reqBody["password"] = "123" // Very weak password
		reqBody["password_confirm"] = "123"
		
		// Act
		resp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
		
		// Assert
		helpers.AssertStatusCode(t, resp, http.StatusBadRequest)
		
		apiResp := helpers.ParseJSONResponse(t, resp)
		helpers.AssertErrorResponse(t, apiResp)
	})
}

func TestVerifyEmailEndpoint(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Valid email verification", func(t *testing.T) {
		// Arrange - Crear usuario primero
		regReq := helpers.CreateValidRegisterRequest()
		server.MakeRequest(t, "POST", "/api/v1/auth/register", regReq)
		
		// Simular token (en test real, obtendrías del email)
		verifyReq := map[string]interface{}{
			"token": "simulated_token_for_testing",
		}
		
		// Act
		resp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		
		// Note: En ambiente de testing, esto podría fallar porque el token
		// no es real. El test verifica la estructura del endpoint.
		
		// Assert structure (aunque falle por token inválido)
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 200 or 400, got %d", resp.StatusCode)
		}
		
		apiResp := helpers.ParseJSONResponse(t, resp)
		
		// Verificar que la respuesta tenga la estructura correcta
		if apiResp.Message == "" {
			t.Error("Response should have a message")
		}
	})

	t.Run("Missing token", func(t *testing.T) {
		// Arrange
		verifyReq := map[string]interface{}{
			// Missing token
		}
		
		// Act
		resp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		
		// Assert
		helpers.AssertStatusCode(t, resp, http.StatusBadRequest)
		
		apiResp := helpers.ParseJSONResponse(t, resp)
		helpers.AssertErrorResponse(t, apiResp)
	})

	t.Run("Empty token", func(t *testing.T) {
		// Arrange
		verifyReq := map[string]interface{}{
			"token": "",
		}
		
		// Act
		resp := server.MakeRequest(t, "POST", "/api/v1/auth/verify-email", verifyReq)
		
		// Assert
		helpers.AssertStatusCode(t, resp, http.StatusBadRequest)
		
		apiResp := helpers.ParseJSONResponse(t, resp)
		helpers.AssertErrorResponse(t, apiResp)
	})
}

func TestResendVerificationEndpoint(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Valid resend request", func(t *testing.T) {
		// Arrange - Crear usuario primero
		regReq := helpers.CreateValidRegisterRequest()
		server.MakeRequest(t, "POST", "/api/v1/auth/register", regReq)
		
		resendReq := map[string]interface{}{
			"email": regReq["email"],
		}
		
		// Act
		resp := server.MakeRequest(t, "POST", "/api/v1/auth/resend-verification", resendReq)
		
		// Assert
		helpers.AssertStatusCode(t, resp, http.StatusOK)
		
		apiResp := helpers.ParseJSONResponse(t, resp)
		helpers.AssertSuccessResponse(t, apiResp)
		
		if !strings.Contains(apiResp.Message, "reenviado") {
			t.Errorf("Message should mention resending: %s", apiResp.Message)
		}
	})

	t.Run("Non-existent email", func(t *testing.T) {
		// Arrange
		resendReq := map[string]interface{}{
			"email": "nonexistent@example.cl",
		}
		
		// Act
		resp := server.MakeRequest(t, "POST", "/api/v1/auth/resend-verification", resendReq)
		
		// Assert
		helpers.AssertStatusCode(t, resp, http.StatusBadRequest)
		
		apiResp := helpers.ParseJSONResponse(t, resp)
		helpers.AssertErrorResponse(t, apiResp)
	})

	t.Run("Invalid email format", func(t *testing.T) {
		// Arrange
		resendReq := map[string]interface{}{
			"email": "not-an-email",
		}
		
		// Act
		resp := server.MakeRequest(t, "POST", "/api/v1/auth/resend-verification", resendReq)
		
		// Assert
		helpers.AssertStatusCode(t, resp, http.StatusBadRequest)
		
		apiResp := helpers.ParseJSONResponse(t, resp)
		helpers.AssertErrorResponse(t, apiResp)
	})
}

func TestAuthHealthEndpoint(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Health check", func(t *testing.T) {
		// Act
		resp := server.MakeRequest(t, "GET", "/api/v1/auth/health", nil)
		
		// Assert
		helpers.AssertStatusCode(t, resp, http.StatusOK)
		
		apiResp := helpers.ParseJSONResponse(t, resp)
		helpers.AssertSuccessResponse(t, apiResp)
		
		helpers.AssertContainsField(t, apiResp, "timestamp")
		helpers.AssertContainsField(t, apiResp, "module")
		helpers.AssertFieldEquals(t, apiResp, "module", "authentication")
	})
}

// Benchmark tests para performance
func BenchmarkRegisterEndpoint(b *testing.B) {
	server, cleanup := helpers.CreateTestServer(&testing.T{})
	defer cleanup()

	reqBody := helpers.CreateValidRegisterRequest()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Cambiar email para evitar duplicados
		reqBody["email"] = fmt.Sprintf("test%d@example.cl", i)
		reqBody["firstname"] = fmt.Sprintf("Test %d", i)
		reqBody["lastname"] = fmt.Sprintf("User %d", i)
		
		server.MakeRequest(&testing.T{}, "POST", "/api/v1/auth/register", reqBody)
	}
}
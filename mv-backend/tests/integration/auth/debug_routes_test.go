package auth_integration_test

import (
	"net/http"
	"testing"

	"github.com/JoseLuis21/mv-backend/tests/helpers"
	"github.com/stretchr/testify/require"
)

// TestDebugRoutes verifica que las rutas básicas estén funcionando
func TestDebugRoutes(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Health endpoint should work", func(t *testing.T) {
		resp := server.MakeRequest(t, "GET", "/api/v1/health", nil)
		t.Logf("Health endpoint status: %d", resp.StatusCode)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Auth health endpoint should work", func(t *testing.T) {
		resp := server.MakeRequest(t, "GET", "/api/v1/auth/health", nil)
		t.Logf("Auth health endpoint status: %d", resp.StatusCode)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Register endpoint should exist", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"firstname":        "Test",
			"lastname":         "User",
			"email":            "invalid-email",
			"phone":            "+56912345678",
			"password":         "password123",
			"password_confirm": "password123",
		}

		resp := server.MakeRequest(t, "POST", "/api/v1/auth/register", reqBody)
		t.Logf("Register endpoint status: %d", resp.StatusCode)
		// Debería retornar 400 por email inválido, no 404
		require.NotEqual(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Forgot password endpoint should exist", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email": "test@example.com",
		}

		resp := server.MakeRequest(t, "POST", "/api/v1/auth/forgot-password", reqBody)
		t.Logf("Forgot password endpoint status: %d", resp.StatusCode)
		// Debería retornar algo distinto de 404
		require.NotEqual(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Reset password endpoint should exist", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"token":       "fake-token",
			"newPassword": "newpassword123",
		}

		resp := server.MakeRequest(t, "POST", "/api/v1/auth/reset-password", reqBody)
		t.Logf("Reset password endpoint status: %d", resp.StatusCode)
		// Debería retornar algo distinto de 404
		require.NotEqual(t, http.StatusNotFound, resp.StatusCode)
	})
}
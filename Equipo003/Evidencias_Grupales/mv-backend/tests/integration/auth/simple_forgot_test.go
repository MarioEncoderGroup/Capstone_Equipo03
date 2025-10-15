package auth_integration_test

import (
	"net/http"
	"testing"

	"github.com/JoseLuis21/mv-backend/tests/helpers"
	"github.com/stretchr/testify/require"
)

// TestSimpleForgotPassword test muy básico para debuggear
func TestSimpleForgotPassword(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Basic forgot password request", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email": "test@example.com",
		}

		resp := server.MakeRequest(t, "POST", "/api/v1/auth/forgot-password", reqBody)

		// Log detailed response for debugging
		t.Logf("Response status: %d", resp.StatusCode)
		t.Logf("Response headers: %v", resp.Header)

		apiResp := helpers.ParseJSONResponse(t, resp)
		t.Logf("Response body: %+v", apiResp)

		// El endpoint debe existir (no debe ser 404)
		require.NotEqual(t, http.StatusNotFound, resp.StatusCode,
			"Endpoint should exist, got 404 - check route configuration")
	})
}
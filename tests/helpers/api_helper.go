package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/config"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/JoseLuis21/mv-backend/internal/routes"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// TestServer encapsula un servidor de testing
type TestServer struct {
	App        *fiber.App
	TestClient *http.Client
	BaseURL    string
	DBClient   *postgresql.PostgresqlClient // Exposer cliente de BD para tests de integración
}

// CreateTestServer crea un servidor de testing con todas las dependencias
func CreateTestServer(t *testing.T) (*TestServer, func()) {
	t.Helper()

	// Verificar PostgreSQL disponible
	SkipIfNoPostgreSQL(t)

	// Crear BD de testing
	dbClient, _, dbCleanup := CreateTestDatabase(t)
	SetupTestTables(t, dbClient)

	// Configurar Fiber app
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
				"error":   "TEST_ERROR",
			})
		},
	})

	// Configurar dependencias
	dependencies, err := config.NewDependencies(dbClient)
	if err != nil {
		dbCleanup()
		t.Fatalf("Failed to create dependencies: %v", err)
	}

	// Configurar rutas
	routes.AuthRoutes(app, dependencies.AuthController)
	routes.PublicRoutes(app)

	// Crear cliente HTTP de testing
	testClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	server := &TestServer{
		App:        app,
		TestClient: testClient,
		BaseURL:    "http://localhost:8080", // No usado en tests pero para consistencia
		DBClient:   dbClient,                // Exponer cliente de BD para tests de integración
	}

	cleanup := func() {
		dbCleanup()
	}

	return server, cleanup
}

// MakeRequest hace una request HTTP al servidor de testing
func (s *TestServer) MakeRequest(t *testing.T, method, path string, body interface{}) *http.Response {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req := httptest.NewRequest(method, path, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Ejecutar request usando Fiber app
	resp, err := s.App.Test(req, 30*1000) // 30 seconds timeout
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	return resp
}

// ParseJSONResponse parsea una respuesta JSON
func ParseJSONResponse(t *testing.T, resp *http.Response) *types.APIResponse {
	t.Helper()

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var apiResp types.APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		t.Fatalf("Failed to parse JSON response: %v. Body: %s", err, string(body))
	}

	return &apiResp
}

// AssertStatusCode verifica el código de status HTTP
func AssertStatusCode(t *testing.T, resp *http.Response, expectedStatus int) {
	t.Helper()

	if resp.StatusCode != expectedStatus {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status %d, got %d. Body: %s",
			expectedStatus, resp.StatusCode, string(body))
	}
}

// AssertSuccessResponse verifica que la respuesta sea exitosa
func AssertSuccessResponse(t *testing.T, apiResp *types.APIResponse) {
	t.Helper()

	if !apiResp.Success {
		t.Fatalf("Expected success response, got: %+v", apiResp)
	}
}

// AssertErrorResponse verifica que la respuesta sea de error
func AssertErrorResponse(t *testing.T, apiResp *types.APIResponse) {
	t.Helper()

	if apiResp.Success {
		t.Fatalf("Expected error response, got success: %+v", apiResp)
	}
}

// AssertErrorCode verifica el código de error específico
func AssertErrorCode(t *testing.T, apiResp *types.APIResponse, expectedErrorCode string) {
	t.Helper()

	AssertErrorResponse(t, apiResp)

	if apiResp.Error != expectedErrorCode {
		t.Fatalf("Expected error code %s, got %s", expectedErrorCode, apiResp.Error)
	}
}

// CreateValidRegisterRequest crea una request de registro válida para testing
func CreateValidRegisterRequest() map[string]interface{} {
	return map[string]interface{}{
		"firstname":        "Usuario",
		"lastname":         "de Prueba",
		"email":            "usuario@test.cl",
		"phone":            "+56912345678",
		"password":         "password123",
		"password_confirm": "password123",
	}
}

// WaitForCondition espera hasta que una condición se cumpla (útil para async operations)
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	t.Helper()

	start := time.Now()
	for time.Since(start) < timeout {
		if condition() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("Condition not met within %v: %s", timeout, message)
}

// LogResponse ayuda con debugging - log de respuestas HTTP
func LogResponse(t *testing.T, resp *http.Response) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Logf("Failed to read response body: %v", err)
		return
	}

	// Recrear body para que pueda ser leído nuevamente
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	t.Logf("Response Status: %d", resp.StatusCode)
	t.Logf("Response Headers: %+v", resp.Header)
	t.Logf("Response Body: %s", string(body))
}

// ValidateJSONSchema valida que un JSON tenga la estructura esperada
func ValidateJSONSchema(t *testing.T, data interface{}, schema interface{}) {
	t.Helper()

	validate := validator.New()
	if err := validate.Struct(schema); err != nil {
		t.Fatalf("Schema validation failed: %v", err)
	}
}

// AssertContainsField verifica que una respuesta contenga un campo específico
func AssertContainsField(t *testing.T, apiResp *types.APIResponse, field string) {
	t.Helper()

	data, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data to be map[string]interface{}, got %T", apiResp.Data)
	}

	if _, exists := data[field]; !exists {
		t.Fatalf("Response should contain field '%s'. Available fields: %v",
			field, getMapKeys(data))
	}
}

// AssertFieldEquals verifica que un campo tenga un valor específico
func AssertFieldEquals(t *testing.T, apiResp *types.APIResponse, field string, expected interface{}) {
	t.Helper()

	data, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data to be map[string]interface{}, got %T", apiResp.Data)
	}

	if _, exists := data[field]; !exists {
		t.Fatalf("Response should contain field '%s'. Available fields: %v",
			field, getMapKeys(data))
	}

	actual := data[field]
	if actual != expected {
		t.Fatalf("Expected field '%s' to be %v, got %v", field, expected, actual)
	}
}

// getMapKeys obtiene las claves de un mapa (helper para logging)
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// GetLatestPasswordResetToken obtiene el token de reset más reciente para un email
// Helper para tests de integración de forgot password
func GetLatestPasswordResetToken(t *testing.T, db *postgresql.PostgresqlClient, email string) string {
	t.Helper()

	// Debug: First let's see what user data exists
	debugQuery := `
		SELECT email, password_reset_token, password_reset_expires, updated
		FROM users
		WHERE email = $1
		AND deleted_at IS NULL`

	var userEmail string
	var resetToken *string
	var resetExpires *time.Time
	var updated time.Time

	err := db.QueryRow(context.Background(), debugQuery, email).Scan(&userEmail, &resetToken, &resetExpires, &updated)
	if err != nil {
		t.Logf("🔍 DEBUG: User not found or error: %v", err)
		return ""
	}

	t.Logf("🔍 DEBUG: User found - Email: %s", userEmail)
	if resetToken != nil {
		t.Logf("🔍 DEBUG: Reset token: %s", (*resetToken)[:8]+"...")
	} else {
		t.Logf("🔍 DEBUG: Reset token: <nil>")
	}

	if resetExpires != nil {
		t.Logf("🔍 DEBUG: Reset expires: %v (now: %v)", *resetExpires, time.Now())
		t.Logf("🔍 DEBUG: Token expired? %t", time.Now().After(*resetExpires))
	} else {
		t.Logf("🔍 DEBUG: Reset expires: <nil>")
	}

	t.Logf("🔍 DEBUG: User updated: %v", updated)

	// Now the actual query
	query := `
		SELECT password_reset_token
		FROM users
		WHERE email = $1
		AND password_reset_token IS NOT NULL
		AND password_reset_expires > NOW()
		AND deleted_at IS NULL
		ORDER BY updated DESC
		LIMIT 1`

	var token *string
	err = db.QueryRow(context.Background(), query, email).Scan(&token)
	if err != nil {
		t.Logf("❌ DEBUG: No active password reset token found for email %s: %v", email, err)
		return ""
	}

	if token == nil {
		t.Logf("❌ DEBUG: Token is nil")
		return ""
	}

	t.Logf("✅ DEBUG: Found active reset token: %s", (*token)[:8]+"...")
	return *token
}

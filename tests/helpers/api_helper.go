package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/config"
	"github.com/JoseLuis21/mv-backend/internal/routes"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// TestServer encapsula un servidor de testing
type TestServer struct {
	App        *fiber.App
	TestClient *http.Client
	BaseURL    string
}

// APIResponse representa una respuesta estándar de la API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
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
func ParseJSONResponse(t *testing.T, resp *http.Response) *APIResponse {
	t.Helper()

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var apiResp APIResponse
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
func AssertSuccessResponse(t *testing.T, apiResp *APIResponse) {
	t.Helper()

	if !apiResp.Success {
		t.Fatalf("Expected success response, got: %+v", apiResp)
	}
}

// AssertErrorResponse verifica que la respuesta sea de error
func AssertErrorResponse(t *testing.T, apiResp *APIResponse) {
	t.Helper()

	if apiResp.Success {
		t.Fatalf("Expected error response, got success: %+v", apiResp)
	}
}

// AssertErrorCode verifica el código de error específico
func AssertErrorCode(t *testing.T, apiResp *APIResponse, expectedErrorCode string) {
	t.Helper()

	AssertErrorResponse(t, apiResp)

	if apiResp.Error != expectedErrorCode {
		t.Fatalf("Expected error code %s, got %s", expectedErrorCode, apiResp.Error)
	}
}

// CreateValidRegisterRequest crea una request de registro válida para testing
func CreateValidRegisterRequest() map[string]interface{} {
	return map[string]interface{}{
		"username":              "testuser",
		"full_name":             "Usuario de Prueba",
		"email":                 "usuario@test.cl",
		"password":              "password123",
		"phone":                 "+56912345678",
		"identification_number": "12.345.678-5",
	}
}

// CreateValidTenantRegisterRequest crea una request de registro con tenant válida
func CreateValidTenantRegisterRequest() map[string]interface{} {
	baseReq := CreateValidRegisterRequest()
	baseReq["create_tenant"] = true
	baseReq["tenant_data"] = map[string]interface{}{
		"rut":           "76.123.456-0",
		"business_name": "Empresa de Prueba SpA",
		"email":         "contacto@empresa.cl",
		"phone":         "+56987654321",
		"address":       "Av. Providencia 123, Santiago",
		"website":       "https://empresa.cl",
		"region_id":     "RM",
		"commune_id":    "Santiago",
		"country_id":    "01234567-89ab-cdef-0123-456789abcdef",
	}
	return baseReq
}

// CreateInvalidRUTRequest crea una request con RUT inválido
func CreateInvalidRUTRequest() map[string]interface{} {
	req := CreateValidRegisterRequest()
	req["identification_number"] = "12.345.678-X" // RUT inválido
	return req
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
func AssertContainsField(t *testing.T, apiResp *APIResponse, field string) {
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
func AssertFieldEquals(t *testing.T, apiResp *APIResponse, field string, expected interface{}) {
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

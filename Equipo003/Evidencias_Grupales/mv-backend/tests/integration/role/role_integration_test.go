package role_integration_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/JoseLuis21/mv-backend/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCompleteRoleFlow valida el flujo completo de gestión de roles
func TestCompleteRoleFlow(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Complete flow: Create -> Get -> Update -> Delete Role", func(t *testing.T) {
		// PASO 1: CREAR ROL - Crear un nuevo rol
		createReq := map[string]interface{}{
			"name":        "Test Manager",
			"description": "Manager role for testing purposes",
		}

		createResp := server.MakeRequest(t, "POST", "/api/v1/roles", createReq)
		require.Equal(t, http.StatusCreated, createResp.StatusCode)

		createApiResp := helpers.ParseJSONResponse(t, createResp)
		require.True(t, createApiResp.Success)
		assert.Equal(t, "Rol creado exitosamente", createApiResp.Message)

		// Extraer ID del rol creado
		data, ok := createApiResp.Data.(map[string]interface{})
		require.True(t, ok, "Response data should be a map")

		roleIDStr, ok := data["id"].(string)
		require.True(t, ok, "Role ID should be present")
		require.NotEmpty(t, roleIDStr)

		t.Logf("✅ PASO 1: Rol creado con ID: %s", roleIDStr)

		// PASO 2: OBTENER ROL - Verificar que se puede obtener el rol
		getResp := server.MakeRequest(t, "GET", "/api/v1/roles/"+roleIDStr, nil)
		require.Equal(t, http.StatusOK, getResp.StatusCode)

		getApiResp := helpers.ParseJSONResponse(t, getResp)
		require.True(t, getApiResp.Success)

		// Verificar estructura de respuesta
		getRoleData, ok := getApiResp.Data.(map[string]interface{})
		require.True(t, ok, "Get response data should be a map")

		assert.Equal(t, roleIDStr, getRoleData["id"])
		assert.Equal(t, "Test Manager", getRoleData["name"])
		assert.Equal(t, "Manager role for testing purposes", getRoleData["description"])
		assert.Equal(t, false, getRoleData["is_global"])
		assert.Equal(t, false, getRoleData["is_system"])

		t.Logf("✅ PASO 2: Rol obtenido correctamente")

		// PASO 3: ACTUALIZAR ROL - Modificar el rol
		updateReq := map[string]interface{}{
			"name":        "Updated Manager",
			"description": "Updated manager role description",
		}

		updateResp := server.MakeRequest(t, "PUT", "/api/v1/roles/"+roleIDStr, updateReq)
		require.Equal(t, http.StatusOK, updateResp.StatusCode)

		updateApiResp := helpers.ParseJSONResponse(t, updateResp)
		require.True(t, updateApiResp.Success)
		assert.Equal(t, "Rol actualizado exitosamente", updateApiResp.Message)

		// Verificar que los cambios se aplicaron
		updateData, ok := updateApiResp.Data.(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, "Updated Manager", updateData["name"])
		assert.Equal(t, "Updated manager role description", updateData["description"])

		t.Logf("✅ PASO 3: Rol actualizado correctamente")

		// PASO 4: LISTAR ROLES - Verificar que aparece en la lista
		listResp := server.MakeRequest(t, "GET", "/api/v1/roles?page=1&limit=10", nil)
		require.Equal(t, http.StatusOK, listResp.StatusCode)

		listApiResp := helpers.ParseJSONResponse(t, listResp)
		require.True(t, listApiResp.Success)

		listData, ok := listApiResp.Data.(map[string]interface{})
		require.True(t, ok)

		roles, ok := listData["roles"].([]interface{})
		require.True(t, ok)
		assert.True(t, len(roles) > 0, "Should have at least one role")

		// Buscar nuestro rol en la lista
		found := false
		for _, roleInterface := range roles {
			role, ok := roleInterface.(map[string]interface{})
			require.True(t, ok)
			if role["id"] == roleIDStr {
				found = true
				assert.Equal(t, "Updated Manager", role["name"])
				break
			}
		}
		assert.True(t, found, "Updated role should be found in list")

		t.Logf("✅ PASO 4: Rol encontrado en listado")

		// PASO 5: ELIMINAR ROL - Eliminar el rol
		deleteResp := server.MakeRequest(t, "DELETE", "/api/v1/roles/"+roleIDStr, nil)
		require.Equal(t, http.StatusOK, deleteResp.StatusCode)

		deleteApiResp := helpers.ParseJSONResponse(t, deleteResp)
		require.True(t, deleteApiResp.Success)
		assert.Equal(t, "Rol eliminado exitosamente", deleteApiResp.Message)

		t.Logf("✅ PASO 5: Rol eliminado correctamente")

		// PASO 6: VERIFICAR ELIMINACIÓN - El rol no debe estar disponible
		getDeletedResp := server.MakeRequest(t, "GET", "/api/v1/roles/"+roleIDStr, nil)
		require.Equal(t, http.StatusNotFound, getDeletedResp.StatusCode)

		getDeletedApiResp := helpers.ParseJSONResponse(t, getDeletedResp)
		require.False(t, getDeletedApiResp.Success)
		assert.Equal(t, "ROLE_NOT_FOUND", getDeletedApiResp.Error)

		t.Logf("✅ PASO 6: Verificación de eliminación exitosa")
	})
}

func TestRoleSecurityScenarios(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Should prevent creation of system role names", func(t *testing.T) {
		createReq := map[string]interface{}{
			"name":        "administrator", // System role name
			"description": "Trying to create system role",
		}

		createResp := server.MakeRequest(t, "POST", "/api/v1/roles", createReq)
		require.Equal(t, http.StatusBadRequest, createResp.StatusCode)

		createApiResp := helpers.ParseJSONResponse(t, createResp)
		require.False(t, createApiResp.Success)
		assert.Equal(t, "SYSTEM_ROLE_READONLY", createApiResp.Error)
		assert.Contains(t, createApiResp.Message, "reservados del sistema")
	})

	t.Run("Should prevent duplicate role names", func(t *testing.T) {
		// Crear primer rol
		createReq := map[string]interface{}{
			"name":        "Unique Role",
			"description": "First role with this name",
		}

		createResp1 := server.MakeRequest(t, "POST", "/api/v1/roles", createReq)
		require.Equal(t, http.StatusCreated, createResp1.StatusCode)

		// Intentar crear segundo rol con mismo nombre
		createResp2 := server.MakeRequest(t, "POST", "/api/v1/roles", createReq)
		require.Equal(t, http.StatusBadRequest, createResp2.StatusCode)

		createApiResp2 := helpers.ParseJSONResponse(t, createResp2)
		require.False(t, createApiResp2.Success)
		assert.Equal(t, "ROLE_ALREADY_EXISTS", createApiResp2.Error)
		assert.Contains(t, createApiResp2.Message, "Ya existe un rol")
	})

	t.Run("Should prevent modification of system roles", func(t *testing.T) {
		// Crear un rol del sistema (esto debería fallar, pero si existiera...)
		// Por ahora simulamos obteniendo un rol del sistema existente
		// TODO: Implementar cuando tengamos roles del sistema inicializados

		updateReq := map[string]interface{}{
			"name":        "Modified Administrator",
			"description": "Trying to modify system role",
		}

		// Usar un ID ficticio de rol del sistema
		updateResp := server.MakeRequest(t, "PUT", "/api/v1/roles/system-role-id", updateReq)

		// Debería retornar 404 o 400 dependiendo de si existe
		assert.True(t, updateResp.StatusCode == http.StatusNotFound || updateResp.StatusCode == http.StatusBadRequest)
	})
}

func TestRoleInputValidation(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Should validate role creation input", func(t *testing.T) {
		testCases := []struct {
			name           string
			roleName       string
			description    string
			expectedStatus int
			expectedError  string
		}{
			{
				name:           "Empty name",
				roleName:       "",
				description:    "Valid description",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "VALIDATION_ERROR",
			},
			{
				name:           "Name too short",
				roleName:       "ab",
				description:    "Valid description",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "VALIDATION_ERROR",
			},
			{
				name:           "Name too long",
				roleName:       "this-is-a-very-long-role-name-that-exceeds-the-maximum-allowed-length",
				description:    "Valid description",
				expectedStatus: http.StatusBadRequest,
				expectedError:  "VALIDATION_ERROR",
			},
			{
				name:           "Description too long",
				roleName:       "Valid Name",
				description:    string(make([]byte, 501)), // 501 characters
				expectedStatus: http.StatusBadRequest,
				expectedError:  "VALIDATION_ERROR",
			},
			{
				name:           "Valid role",
				roleName:       "Valid Role",
				description:    "Valid description",
				expectedStatus: http.StatusCreated,
				expectedError:  "",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				createReq := map[string]interface{}{
					"name":        tc.roleName,
					"description": tc.description,
				}

				createResp := server.MakeRequest(t, "POST", "/api/v1/roles", createReq)
				assert.Equal(t, tc.expectedStatus, createResp.StatusCode)

				createApiResp := helpers.ParseJSONResponse(t, createResp)

				if tc.expectedError != "" {
					assert.False(t, createApiResp.Success)
					assert.Equal(t, tc.expectedError, createApiResp.Error)
				} else {
					assert.True(t, createApiResp.Success)
				}
			})
		}
	})

	t.Run("Should handle malformed JSON", func(t *testing.T) {
		// Enviar JSON malformado
		req := `{"name": "Test Role", "description":}`
		resp := server.MakeRequest(t, "POST", "/api/v1/roles", req)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		apiResp := helpers.ParseJSONResponse(t, resp)
		require.False(t, apiResp.Success)
		assert.Equal(t, "INVALID_REQUEST_FORMAT", apiResp.Error)
		assert.Contains(t, apiResp.Message, "formato")
	})

	t.Run("Should validate role update input", func(t *testing.T) {
		// Crear rol válido primero
		createReq := map[string]interface{}{
			"name":        "Role to Update",
			"description": "Original description",
		}

		createResp := server.MakeRequest(t, "POST", "/api/v1/roles", createReq)
		require.Equal(t, http.StatusCreated, createResp.StatusCode)

		createApiResp := helpers.ParseJSONResponse(t, createResp)
		data := createApiResp.Data.(map[string]interface{})
		roleID := data["id"].(string)

		// Probar actualizaciones inválidas
		updateReq := map[string]interface{}{
			"name":        "", // Nombre vacío
			"description": "Valid description",
		}

		updateResp := server.MakeRequest(t, "PUT", "/api/v1/roles/"+roleID, updateReq)
		require.Equal(t, http.StatusBadRequest, updateResp.StatusCode)

		updateApiResp := helpers.ParseJSONResponse(t, updateResp)
		require.False(t, updateApiResp.Success)
		assert.Equal(t, "VALIDATION_ERROR", updateApiResp.Error)
	})
}

func TestRoleListingAndFiltering(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Should list roles with pagination", func(t *testing.T) {
		// Crear varios roles para probar paginación
		for i := 1; i <= 5; i++ {
			createReq := map[string]interface{}{
				"name":        fmt.Sprintf("Test Role %d", i),
				"description": fmt.Sprintf("Description for role %d", i),
			}

			createResp := server.MakeRequest(t, "POST", "/api/v1/roles", createReq)
			require.Equal(t, http.StatusCreated, createResp.StatusCode)
		}

		// Probar paginación
		listResp := server.MakeRequest(t, "GET", "/api/v1/roles?page=1&limit=3", nil)
		require.Equal(t, http.StatusOK, listResp.StatusCode)

		listApiResp := helpers.ParseJSONResponse(t, listResp)
		require.True(t, listApiResp.Success)

		data, ok := listApiResp.Data.(map[string]interface{})
		require.True(t, ok)

		roles, ok := data["roles"].([]interface{})
		require.True(t, ok)

		assert.True(t, len(roles) <= 3, "Should respect limit")
		assert.Equal(t, float64(1), data["page"])
		assert.Equal(t, float64(3), data["limit"])
		assert.True(t, data["total"].(float64) >= 5, "Should have at least 5 roles")
	})

	t.Run("Should filter roles by name", func(t *testing.T) {
		// Crear rol con nombre específico
		createReq := map[string]interface{}{
			"name":        "Filterable Manager",
			"description": "Role for filtering test",
		}

		createResp := server.MakeRequest(t, "POST", "/api/v1/roles", createReq)
		require.Equal(t, http.StatusCreated, createResp.StatusCode)

		// Filtrar por nombre
		listResp := server.MakeRequest(t, "GET", "/api/v1/roles?name=Filterable", nil)
		require.Equal(t, http.StatusOK, listResp.StatusCode)

		listApiResp := helpers.ParseJSONResponse(t, listResp)
		require.True(t, listApiResp.Success)

		data := listApiResp.Data.(map[string]interface{})
		roles := data["roles"].([]interface{})

		// Debería encontrar el rol creado
		found := false
		for _, roleInterface := range roles {
			role := roleInterface.(map[string]interface{})
			if role["name"] == "Filterable Manager" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find the filterable role")
	})
}

func TestRoleResponseStructure(t *testing.T) {
	server, cleanup := helpers.CreateTestServer(t)
	defer cleanup()

	t.Run("Role response should have correct structure", func(t *testing.T) {
		// Crear rol
		createReq := map[string]interface{}{
			"name":        "Structure Test Role",
			"description": "Role for testing response structure",
		}

		createResp := server.MakeRequest(t, "POST", "/api/v1/roles", createReq)
		require.Equal(t, http.StatusCreated, createResp.StatusCode)

		createApiResp := helpers.ParseJSONResponse(t, createResp)
		require.True(t, createApiResp.Success)

		// Verificar estructura completa
		helpers.AssertContainsField(t, createApiResp, "id")
		helpers.AssertContainsField(t, createApiResp, "name")
		helpers.AssertContainsField(t, createApiResp, "description")
		helpers.AssertContainsField(t, createApiResp, "is_global")
		helpers.AssertContainsField(t, createApiResp, "is_system")
		helpers.AssertContainsField(t, createApiResp, "created")
		helpers.AssertContainsField(t, createApiResp, "updated")

		// Verificar tipos de datos
		roleData := createApiResp.Data.(map[string]interface{})

		assert.IsType(t, "", roleData["id"], "ID should be string")
		assert.IsType(t, "", roleData["name"], "Name should be string")
		assert.IsType(t, "", roleData["description"], "Description should be string")
		assert.IsType(t, false, roleData["is_global"], "IsGlobal should be boolean")
		assert.IsType(t, false, roleData["is_system"], "IsSystem should be boolean")
		assert.IsType(t, "", roleData["created"], "Created should be string (ISO 8601)")
		assert.IsType(t, "", roleData["updated"], "Updated should be string (ISO 8601)")

		t.Logf("✅ Role response structure validated")
	})
}
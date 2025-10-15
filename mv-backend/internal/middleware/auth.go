package middleware

import (
	"strings"

	"github.com/JoseLuis21/mv-backend/internal/core/user/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/JoseLuis21/mv-backend/internal/shared/tokens"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT tokens and user authentication for MisViaticos
func AuthMiddleware(dbControl *postgresql.PostgresqlClient) fiber.Handler {
	// Inicializar servicios necesarios
	tokenService := tokens.NewService()

	return func(c *fiber.Ctx) error {
		// 1. Obtener el header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
				Success: false,
				Message: "Autenticación requerida",
				Error:   "MISSING_AUTH_HEADER",
			})
		}

		// 2. Extraer token del formato "Bearer <token>"
		tokenString := extractBearerToken(authHeader)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
				Success: false,
				Message: "Formato de token inválido",
				Error:   "INVALID_AUTH_FORMAT",
			})
		}

		// 3. Validar JWT y obtener claims
		claims, err := tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
				Success: false,
				Message: "Token inválido o expirado",
				Error:   "INVALID_TOKEN",
			})
		}

		// 4. Verificar que el usuario existe y está activo
		// TODO: Implementar verificación de usuario en BD cuando se necesite
		// Por ahora, confiamos en el token JWT válido

		// 5. Guardar información del usuario en el contexto de Fiber
		c.Locals("userID", claims.UserID.String())

		// Crear y guardar objeto User en el contexto para RBAC middleware
		user := fiber.Map{
			"id": claims.UserID.String(),
		}

		if claims.TenantID != nil {
			c.Locals("tenantID", claims.TenantID.String())
			user["tenantId"] = claims.TenantID.String()
		}

		c.Locals("user", user)

		// Guardar roles y permisos en el contexto
		if len(claims.Roles) > 0 {
			c.Locals("roles", claims.Roles)
		}

		if len(claims.Permissions) > 0 {
			c.Locals("permissions", claims.Permissions)
		}

		// 6. Continuar con la siguiente función en la cadena
		return c.Next()
	}
}


// AuthMiddlewareWithUserService valida JWT y verifica usuario en BD
func AuthMiddlewareWithUserService(userService ports.UserService) fiber.Handler {
	tokenService := tokens.NewService()

	return func(c *fiber.Ctx) error {
		// 1. Obtener el header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
				Success: false,
				Message: "Autenticación requerida",
				Error:   "MISSING_AUTH_HEADER",
			})
		}

		// 2. Extraer token del formato "Bearer <token>"
		tokenString := extractBearerToken(authHeader)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
				Success: false,
				Message: "Formato de token inválido",
				Error:   "INVALID_AUTH_FORMAT",
			})
		}

		// 3. Validar JWT y obtener claims
		claims, err := tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
				Success: false,
				Message: "Token inválido o expirado",
				Error:   "INVALID_TOKEN",
			})
		}

		// 4. Verificar que el usuario existe y está activo en la BD
		user, err := userService.GetUserByID(c.Context(), claims.UserID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
				Success: false,
				Message: "Usuario no encontrado",
				Error:   "USER_NOT_FOUND",
			})
		}

		if !user.IsActive {
			return c.Status(fiber.StatusForbidden).JSON(types.APIResponse{
				Success: false,
				Message: "Cuenta desactivada",
				Error:   "ACCOUNT_INACTIVE",
			})
		}

		// 5. Guardar información del usuario en el contexto de Fiber
		c.Locals("userID", claims.UserID.String())
		c.Locals("user", user)

		if claims.TenantID != nil {
			c.Locals("tenantID", claims.TenantID.String())
		}

		// Guardar roles y permisos en el contexto
		if len(claims.Roles) > 0 {
			c.Locals("roles", claims.Roles)
		}

		if len(claims.Permissions) > 0 {
			c.Locals("permissions", claims.Permissions)
		}

		// 6. Continuar con la siguiente función en la cadena
		return c.Next()
	}
}

// extractBearerToken extrae el token del header Authorization con formato "Bearer <token>"
func extractBearerToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}
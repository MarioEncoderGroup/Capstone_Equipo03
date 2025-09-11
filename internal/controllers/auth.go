package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	authDomain "github.com/JoseLuis21/mv-backend/internal/core/auth/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/auth/ports"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
)

// AuthController maneja las operaciones de autenticación
type AuthController struct {
	authService ports.AuthService
	validator   *validatorapi.Validator
}

// NewAuthController crea una nueva instancia del controller de autenticación
func NewAuthController(authService ports.AuthService, validator *validatorapi.Validator) *AuthController {
	return &AuthController{
		authService: authService,
		validator:   validator,
	}
}

// APIResponse estructura estándar para respuestas de la API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ValidationErrorResponse estructura para errores de validación
type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Register maneja el registro de usuarios siguiendo el patrón de referencia
func (ac *AuthController) Register(c *fiber.Ctx) error {
	// 1. Parsear request body usando DTOs del dominio
	var req authDomain.AuthRegisterDto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Error parseando datos de entrada",
			Error:   "Formato de datos inválido",
		})
	}

	// 2. Validar estructura de datos
	if errors := ac.validator.ValidateStruct(req); len(errors) > 0 {
		var validationErrors []ValidationErrorResponse
		for _, err := range errors {
			validationErrors = append(validationErrors, ValidationErrorResponse{
				Field:   err.Field,
				Message: err.Message,
			})
		}
		
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Errores de validación",
			Data:    validationErrors,
		})
	}

	// 3. Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 4. Llamar al servicio unificado siguiendo patrón de referencia
	response, err := ac.authService.Register(ctx, &req)

	// 5. Manejar errores del servicio
	if err != nil {
		if appErr, ok := sharedErrors.IsAppError(err); ok {
			return c.Status(appErr.HTTPCode).JSON(APIResponse{
				Success: false,
				Message: appErr.Message,
				Error:   appErr.Code,
				Data: fiber.Map{
					"details": appErr.Details,
				},
			})
		}
		
		// Error genérico no categorizado
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Success: false,
			Message: "Error interno del servidor",
			Error:   "INTERNAL_SERVER_ERROR",
		})
	}

	// 6. Respuesta exitosa siguiendo patrón de referencia con nuevos campos
	responseData := fiber.Map{
		"id":                          response.ID,
		"firstname":                   response.FirstName,
		"lastname":                    response.LastName,
		"full_name":                   response.FullName, // Para backward compatibility
		"email":                       response.Email,
		"phone":                       response.Phone,
		"email_token":                 response.EmailToken,
		"requires_email_verification": response.RequiresEmailVerification,
	}
	
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Success: true,
		Message: response.Message,
		Data:    responseData,
	})
}

// VerifyUserEmail verifica el email del usuario siguiendo patrón de referencia
func (ac *AuthController) VerifyUserEmail(c *fiber.Ctx) error {
	// 1. Parsear request usando DTO del dominio
	var req authDomain.VerifyEmailDto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Error parseando datos",
			Error:   "Token requerido",
		})
	}

	// 2. Validar estructura
	if errors := ac.validator.ValidateStruct(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Token inválido",
			Error:   "Token es requerido",
		})
	}

	// 3. Crear contexto
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Verificar email
	if err := ac.authService.VerifyUserEmail(ctx, req.Token); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Error verificando email",
			Error:   err.Error(),
		})
	}

	// 5. Respuesta exitosa
	return c.JSON(APIResponse{
		Success: true,
		Message: "Email verificado exitosamente. Tu cuenta está ahora activa.",
	})
}

// ResendEmailVerification reenvía el email de verificación siguiendo patrón de referencia
func (ac *AuthController) ResendEmailVerification(c *fiber.Ctx) error {
	// 1. Parsear request usando DTO del dominio
	var req authDomain.ResendVerificationDto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Error parseando datos",
			Error:   "Email requerido",
		})
	}

	// 2. Validar estructura
	if errors := ac.validator.ValidateStruct(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Email inválido",
			Error:   "Email válido es requerido",
		})
	}

	// 3. Crear contexto
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 4. Reenviar email
	if err := ac.authService.ResendEmailVerification(ctx, req.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Error reenviando email",
			Error:   err.Error(),
		})
	}

	// 5. Respuesta exitosa
	return c.JSON(APIResponse{
		Success: true,
		Message: "Email de verificación reenviado exitosamente.",
	})
}

// Login autentica un usuario
func (ac *AuthController) Login(c *fiber.Ctx) error {
	// 1. Parsear request
	var req authDomain.AuthLoginDto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Error parseando datos",
			Error:   "Formato de datos inválido",
		})
	}

	// 2. Validar estructura
	if errors := ac.validator.ValidateStruct(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
		})
	}

	// 3. Crear contexto
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Autenticar usuario
	response, err := ac.authService.Login(ctx, &req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse{
			Success: false,
			Message: "Credenciales inválidas",
			Error:   err.Error(),
		})
	}

	// 5. Respuesta exitosa
	return c.JSON(APIResponse{
		Success: true,
		Message: "Autenticación exitosa",
		Data:    response,
	})
}

// HealthCheck endpoint de health check para el módulo de auth
func (ac *AuthController) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(APIResponse{
		Success: true,
		Message: "Auth module is healthy",
		Data: fiber.Map{
			"timestamp": time.Now().Format(time.RFC3339),
			"module":    "authentication",
			"version":   "1.0.0",
		},
	})
}

// PASO 5: RefreshToken renueva tokens de acceso usando refresh token
func (ac *AuthController) RefreshToken(c *fiber.Ctx) error {
	// 1. Parsear request usando DTO del dominio
	var req authDomain.RefreshTokenDto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Error parseando datos",
			Error:   "Refresh token requerido",
		})
	}

	// 2. Validar estructura
	if errors := ac.validator.ValidateStruct(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Refresh token inválido",
			Error:   "Refresh token es requerido",
		})
	}

	// 3. Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Renovar tokens
	response, err := ac.authService.RefreshAccessToken(ctx, req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse{
			Success: false,
			Message: "Error renovando token",
			Error:   err.Error(),
		})
	}

	// 5. Respuesta exitosa
	return c.JSON(APIResponse{
		Success: true,
		Message: "Tokens renovados exitosamente",
		Data:    response,
	})
}
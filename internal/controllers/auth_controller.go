package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/auth/ports"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
)

// AuthController maneja las operaciones de autenticación
type AuthController struct {
	registerService ports.RegisterService
	validator       *validatorapi.Validator
}

// NewAuthController crea una nueva instancia del controller de autenticación
func NewAuthController(registerService ports.RegisterService, validator *validatorapi.Validator) *AuthController {
	return &AuthController{
		registerService: registerService,
		validator:       validator,
	}
}

// RegisterRequest estructura para la request de registro con validaciones chilenas
type RegisterRequest struct {
	Username             string                    `json:"username" validate:"required,min=3,max=50,alphanum"`
	FullName             string                    `json:"full_name" validate:"required,min=2,max=200"`
	Email                string                    `json:"email" validate:"required,chilean_email"`
	Password             string                    `json:"password" validate:"required,min=8,max=255"`
	Phone                string                    `json:"phone,omitempty" validate:"omitempty,chilean_phone"`
	IdentificationNumber string                    `json:"identification_number,omitempty" validate:"omitempty,chilean_rut"`
	CreateTenant         bool                      `json:"create_tenant,omitempty"`
	TenantData           *TenantRegistrationData   `json:"tenant_data,omitempty"`
}

// TenantRegistrationData estructura para datos del tenant con validaciones chilenas
type TenantRegistrationData struct {
	RUT          string    `json:"rut" validate:"required,chilean_rut"`
	BusinessName string    `json:"business_name" validate:"required,min=2,max=150"`
	Email        string    `json:"email" validate:"required,chilean_email"`
	Phone        string    `json:"phone" validate:"required,chilean_phone"`
	Address      string    `json:"address" validate:"required,min=5,max=200"`
	Website      string    `json:"website" validate:"required,url,max=150"`
	RegionID     string    `json:"region_id" validate:"required,len=2,alpha"`
	CommuneID    string    `json:"commune_id" validate:"required,min=1,max=100"`
	CountryID    uuid.UUID `json:"country_id" validate:"required,uuid"`
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

// Register maneja el registro de usuarios
func (ac *AuthController) Register(c *fiber.Ctx) error {
	// 1. Parsear request body
	var req RegisterRequest
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

	// 3. Convertir a estructura del dominio
	registerReq := &ports.RegisterRequest{
		Username:             req.Username,
		FullName:             req.FullName,
		Email:                req.Email,
		Password:             req.Password,
		Phone:                req.Phone,
		IdentificationNumber: req.IdentificationNumber,
		CreateTenant:         req.CreateTenant,
	}

	// 4. Convertir datos del tenant si están presentes
	if req.TenantData != nil {
		registerReq.TenantData = &ports.TenantRegistrationData{
			RUT:          req.TenantData.RUT,
			BusinessName: req.TenantData.BusinessName,
			Email:        req.TenantData.Email,
			Phone:        req.TenantData.Phone,
			Address:      req.TenantData.Address,
			Website:      req.TenantData.Website,
			RegionID:     req.TenantData.RegionID,
			CommuneID:    req.TenantData.CommuneID,
			CountryID:    req.TenantData.CountryID,
		}
	}

	// 5. Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 6. Llamar al servicio apropiado
	var response *ports.RegisterResponse
	var err error

	if req.CreateTenant {
		response, err = ac.registerService.RegisterUserWithTenant(ctx, registerReq)
	} else {
		response, err = ac.registerService.RegisterUser(ctx, registerReq)
	}

	// 7. Manejar errores del servicio con sistema de errores robusto
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

	// 8. Respuesta exitosa
	responseData := fiber.Map{
		"user_id":                     response.UserID,
		"requires_email_verification": response.RequiresEmailVerification,
	}
	
	// Solo agregar tenant_id si existe
	if response.TenantID != nil {
		responseData["tenant_id"] = response.TenantID
	}
	
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Success: true,
		Message: response.Message,
		Data:    responseData,
	})
}

// VerifyEmailRequest estructura para verificación de email
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required,min=1"`
}

// VerifyUserEmail verifica el email del usuario
func (ac *AuthController) VerifyUserEmail(c *fiber.Ctx) error {
	// 1. Parsear request
	var req VerifyEmailRequest
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
	if err := ac.registerService.VerifyEmail(ctx, req.Token); err != nil {
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

// ResendEmailRequest estructura para reenvío de email
type ResendEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResendEmailVerification reenvía el email de verificación
func (ac *AuthController) ResendEmailVerification(c *fiber.Ctx) error {
	// 1. Parsear request
	var req ResendEmailRequest
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
	if err := ac.registerService.ResendEmailVerification(ctx, req.Email); err != nil {
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
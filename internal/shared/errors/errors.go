package errors

import (
	"fmt"
	"net/http"
)

// AppError representa un error de la aplicación con código HTTP y detalles
type AppError struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Details  string `json:"details,omitempty"`
	HTTPCode int    `json:"-"`
}

// Error implementa la interfaz error
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Errores predefinidos para el dominio de MisViáticos

// Errores de validación
var (
	ErrInvalidRUT = &AppError{
		Code:     "INVALID_RUT",
		Message:  "El RUT proporcionado no es válido",
		HTTPCode: http.StatusBadRequest,
	}
	
	ErrInvalidEmail = &AppError{
		Code:     "INVALID_EMAIL",
		Message:  "El email proporcionado no es válido",
		HTTPCode: http.StatusBadRequest,
	}
	
	ErrInvalidPhone = &AppError{
		Code:     "INVALID_PHONE",
		Message:  "El teléfono proporcionado no es válido para Chile",
		HTTPCode: http.StatusBadRequest,
	}
	
	ErrWeakPassword = &AppError{
		Code:     "WEAK_PASSWORD",
		Message:  "La contraseña debe tener al menos 8 caracteres",
		HTTPCode: http.StatusBadRequest,
	}
)

// Errores de negocio
var (
	ErrUserAlreadyExists = &AppError{
		Code:     "USER_ALREADY_EXISTS",
		Message:  "Ya existe un usuario con este email",
		HTTPCode: http.StatusConflict,
	}
	
	ErrUsernameAlreadyExists = &AppError{
		Code:     "USERNAME_ALREADY_EXISTS",
		Message:  "El nombre de usuario ya está en uso",
		HTTPCode: http.StatusConflict,
	}
	
	ErrTenantAlreadyExists = &AppError{
		Code:     "TENANT_ALREADY_EXISTS",
		Message:  "Ya existe una empresa registrada con este RUT",
		HTTPCode: http.StatusConflict,
	}
	
	ErrUserNotFound = &AppError{
		Code:     "USER_NOT_FOUND",
		Message:  "Usuario no encontrado",
		HTTPCode: http.StatusNotFound,
	}
	
	ErrTenantNotFound = &AppError{
		Code:     "TENANT_NOT_FOUND",
		Message:  "Empresa no encontrada",
		HTTPCode: http.StatusNotFound,
	}
	
	ErrInvalidEmailToken = &AppError{
		Code:     "INVALID_EMAIL_TOKEN",
		Message:  "Token de verificación de email inválido o expirado",
		HTTPCode: http.StatusBadRequest,
	}
	
	ErrEmailAlreadyVerified = &AppError{
		Code:     "EMAIL_ALREADY_VERIFIED",
		Message:  "El email ya está verificado",
		HTTPCode: http.StatusBadRequest,
	}
	
	ErrUserNotActive = &AppError{
		Code:     "USER_NOT_ACTIVE",
		Message:  "El usuario no está activo. Verifica tu email primero.",
		HTTPCode: http.StatusForbidden,
	}
	
	ErrTenantDataRequired = &AppError{
		Code:     "TENANT_DATA_REQUIRED",
		Message:  "Los datos de la empresa son requeridos para el registro de tenant",
		HTTPCode: http.StatusBadRequest,
	}
)

// Errores de sistema
var (
	ErrDatabaseConnection = &AppError{
		Code:     "DATABASE_CONNECTION_ERROR",
		Message:  "Error de conexión con la base de datos",
		HTTPCode: http.StatusInternalServerError,
	}
	
	ErrEmailService = &AppError{
		Code:     "EMAIL_SERVICE_ERROR",
		Message:  "Error enviando email",
		HTTPCode: http.StatusInternalServerError,
	}
	
	ErrTokenGeneration = &AppError{
		Code:     "TOKEN_GENERATION_ERROR",
		Message:  "Error generando token de seguridad",
		HTTPCode: http.StatusInternalServerError,
	}
	
	ErrPasswordHashing = &AppError{
		Code:     "PASSWORD_HASHING_ERROR",
		Message:  "Error procesando contraseña",
		HTTPCode: http.StatusInternalServerError,
	}
	
	ErrTenantDBCreation = &AppError{
		Code:     "TENANT_DB_CREATION_ERROR",
		Message:  "Error creando base de datos de la empresa",
		HTTPCode: http.StatusInternalServerError,
	}
)

// Funciones de conveniencia para crear errores con detalles

// NewValidationError crea un error de validación con detalles
func NewValidationError(message, details string) *AppError {
	return &AppError{
		Code:     "VALIDATION_ERROR",
		Message:  message,
		Details:  details,
		HTTPCode: http.StatusBadRequest,
	}
}

// NewBusinessError crea un error de lógica de negocio
func NewBusinessError(code, message, details string) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		Details:  details,
		HTTPCode: http.StatusBadRequest,
	}
}

// NewInternalError crea un error interno del sistema
func NewInternalError(message, details string) *AppError {
	return &AppError{
		Code:     "INTERNAL_ERROR",
		Message:  message,
		Details:  details,
		HTTPCode: http.StatusInternalServerError,
	}
}

// WrapError envuelve un error estándar en un AppError
func WrapError(baseError *AppError, details string) *AppError {
	return &AppError{
		Code:     baseError.Code,
		Message:  baseError.Message,
		Details:  details,
		HTTPCode: baseError.HTTPCode,
	}
}

// IsAppError verifica si un error es un AppError
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}
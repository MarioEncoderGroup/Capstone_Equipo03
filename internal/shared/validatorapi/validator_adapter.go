package validatorapi

import (
	"github.com/go-playground/validator/v10"
	"github.com/JoseLuis21/mv-backend/internal/shared/utils"
)

// Validator representa un wrapper alrededor del validador para mantener compatibilidad
type Validator struct {
	XValidator *XValidator
}

// ValidationError representa un error de validación simplificado
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidator crea una nueva instancia del validador
func NewValidator() *Validator {
	validate := validator.New()
	xValidator := &XValidator{
		Validator: validate,
	}
	
	// Registrar validaciones personalizadas existentes
	xValidator.RegisterCustomValidations()
	
	// Registrar validaciones específicas chilenas
	registerChileanValidations(validate)
	
	return &Validator{
		XValidator: xValidator,
	}
}

// registerChileanValidations registra validaciones específicas para Chile
func registerChileanValidations(validate *validator.Validate) {
	// Validación de RUT chileno
	validate.RegisterValidation("chilean_rut", func(fl validator.FieldLevel) bool {
		rut := fl.Field().String()
		if rut == "" {
			return true // Permitir campo vacío si no es required
		}
		return utils.ValidateRUT(rut)
	})
	
	// Validación de teléfono chileno
	validate.RegisterValidation("chilean_phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		if phone == "" {
			return true // Permitir campo vacío si no es required
		}
		return utils.ValidateChileanPhone(phone)
	})
	
	// Validación de email chileno (con límites específicos)
	validate.RegisterValidation("chilean_email", func(fl validator.FieldLevel) bool {
		email := fl.Field().String()
		if email == "" {
			return false // Email siempre requerido para esta validación
		}
		return utils.ValidateChileanEmail(email)
	})
}

// ValidateStruct valida una estructura y retorna errores de validación
func (v *Validator) ValidateStruct(data interface{}) []ValidationError {
	errors := v.XValidator.Validate(data)
	
	var result []ValidationError
	for _, err := range errors {
		result = append(result, ValidationError{
			Field:   err.Field,
			Message: err.Message,
		})
	}
	
	return result
}
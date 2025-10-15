package validatorapi

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// XValidator represents a wrapper around the validator
type XValidator struct {
	Validator *validator.Validate
}

// ValidatorErrors represents validation errors for MisViaticos
type ValidatorErrors struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value"`
}

// GlobalErrorHandlerResp represents the global error response structure
type GlobalErrorHandlerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// Validate validates a struct and returns formatted errors
func (v *XValidator) Validate(data interface{}) []ValidatorErrors {
	validationErrors := []ValidatorErrors{}

	errs := v.Validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ValidatorErrors

			elem.Field = strings.ToLower(err.Field())
			elem.Value = err.Param()
			elem.Message = getErrorMessage(err.Tag(), err.Field(), err.Param())

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

// ValidateStruct validates a struct and returns a simple error
func (v *XValidator) ValidateStruct(data interface{}) error {
	errs := v.Validate(data)
	if len(errs) > 0 {
		var errorMessages []string
		for _, err := range errs {
			errorMessages = append(errorMessages, err.Message)
		}
		return errors.New(strings.Join(errorMessages, "; "))
	}
	return nil
}

// getErrorMessage returns a user-friendly error message based on validation tag
func getErrorMessage(tag, field, param string) string {
	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + param + " characters long"
	case "max":
		return field + " cannot exceed " + param + " characters"
	case "len":
		return field + " must be exactly " + param + " characters long"
	case "numeric":
		return field + " must be a number"
	case "alpha":
		return field + " must contain only letters"
	case "alphanum":
		return field + " must contain only letters and numbers"
	case "uuid":
		return field + " must be a valid UUID"
	case "oneof":
		return field + " must be one of: " + param
	case "gt":
		return field + " must be greater than " + param
	case "gte":
		return field + " must be greater than or equal to " + param
	case "lt":
		return field + " must be less than " + param
	case "lte":
		return field + " must be less than or equal to " + param
	// MisViaticos specific validations
	case "currency":
		return field + " must be a valid currency code (CLP, USD, EUR)"
	case "rut":
		return field + " must be a valid Chilean RUT"
	case "amount":
		return field + " must be a valid monetary amount"
	default:
		return field + " is not valid"
	}
}

// Custom validation functions for MisViaticos

// ValidateCurrency validates Chilean currency codes
func ValidateCurrency(fl validator.FieldLevel) bool {
	validCurrencies := map[string]bool{
		"CLP": true,
		"USD": true,
		"EUR": true,
	}
	return validCurrencies[fl.Field().String()]
}

// ValidateRUT validates Chilean RUT format
func ValidateRUT(fl validator.FieldLevel) bool {
	rut := fl.Field().String()
	if len(rut) < 8 || len(rut) > 12 {
		return false
	}
	// Basic RUT format validation (X.XXX.XXX-X)
	// In a real implementation, you'd add the checksum validation
	return strings.Contains(rut, "-")
}

// ValidateAmount validates monetary amounts (positive numbers with max 2 decimals)
func ValidateAmount(fl validator.FieldLevel) bool {
	// This would validate that the amount is a positive number with max 2 decimal places
	// Implementation depends on whether you're using float64, decimal, or string
	return fl.Field().Float() > 0
}

// RegisterCustomValidations registers MisViaticos custom validation rules
func (v *XValidator) RegisterCustomValidations() {
	v.Validator.RegisterValidation("currency", ValidateCurrency)
	v.Validator.RegisterValidation("rut", ValidateRUT)
	v.Validator.RegisterValidation("amount", ValidateAmount)
}
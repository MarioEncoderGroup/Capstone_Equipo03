package types

// APIResponse estructura estándar para respuestas de la API
// Utilizada por todos los controllers para mantener consistencia en las respuestas
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ValidationErrorResponse estructura para errores de validación
// Utilizada para devolver errores específicos de validación de campos
type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
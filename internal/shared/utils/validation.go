package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// ValidateRUT valida el formato y dígito verificador de un RUT chileno
// Formato esperado: XX.XXX.XXX-Y o XXXXXXXX-Y
// Retorna true si el RUT es válido
func ValidateRUT(rut string) bool {
	if rut == "" {
		return false
	}

	// Limpiar el RUT removiendo puntos y espacios
	cleanRUT := strings.ReplaceAll(rut, ".", "")
	cleanRUT = strings.ReplaceAll(cleanRUT, " ", "")
	cleanRUT = strings.ToUpper(cleanRUT)

	// Validar formato básico con regex
	rutRegex := regexp.MustCompile(`^(\d{7,8})-([0-9K])$`)
	matches := rutRegex.FindStringSubmatch(cleanRUT)

	if len(matches) != 3 {
		return false
	}

	rutNumber := matches[1]
	checkDigit := matches[2]

	// Calcular dígito verificador
	calculatedDigit := calculateRUTCheckDigit(rutNumber)

	return calculatedDigit == checkDigit
}

// calculateRUTCheckDigit calcula el dígito verificador de un RUT
func calculateRUTCheckDigit(rutNumber string) string {
	sum := 0
	multiplier := 2

	// Recorrer el RUT de derecha a izquierda
	for i := len(rutNumber) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(rutNumber[i]))
		sum += digit * multiplier
		multiplier++
		if multiplier > 7 {
			multiplier = 2
		}
	}

	remainder := sum % 11
	checkDigit := 11 - remainder

	switch checkDigit {
	case 10:
		return "K"
	case 11:
		return "0"
	default:
		return strconv.Itoa(checkDigit)
	}
}

// FormatRUT formatea un RUT chileno al formato estándar XX.XXX.XXX-Y
func FormatRUT(rut string) string {
	if !ValidateRUT(rut) {
		return rut // Retorna sin formato si no es válido
	}

	// Limpiar el RUT
	cleanRUT := strings.ReplaceAll(rut, ".", "")
	cleanRUT = strings.ReplaceAll(cleanRUT, " ", "")
	cleanRUT = strings.ToUpper(cleanRUT)

	// Separar número y dígito verificador
	parts := strings.Split(cleanRUT, "-")
	if len(parts) != 2 {
		return rut
	}

	rutNumber := parts[0]
	checkDigit := parts[1]

	// Formatear con puntos desde la derecha
	formatted := ""
	for i, digit := range rutNumber {
		// Agregar punto cada 3 dígitos desde la derecha
		position := len(rutNumber) - i
		if i > 0 && position%3 == 0 {
			formatted += "."
		}
		formatted += string(digit)
	}

	return formatted + "-" + checkDigit
}

// ValidateChileanEmail valida email con consideraciones específicas chilenas
func ValidateChileanEmail(email string) bool {
	if email == "" {
		return false
	}

	// Regex básica para email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false
	}

	// Validaciones adicionales específicas para Chile
	// Longitud máxima
	if len(email) > 150 {
		return false
	}

	return true
}

// ValidateChileanPhone valida formato de teléfono chileno
// Acepta: +56912345678, 56912345678, 912345678, 212345678 (fijo Santiago)
func ValidateChileanPhone(phone string) bool {
	if phone == "" {
		return true // Teléfono es opcional
	}

	// Verificar que solo contenga dígitos, +, espacios, - y ()
	if !regexp.MustCompile(`^[\d\s+()-]+$`).MatchString(phone) {
		return false
	}

	// Limpiar espacios y caracteres especiales excepto +
	cleanPhone := regexp.MustCompile(`[\s()-]`).ReplaceAllString(phone, "")

	// Patrones válidos para teléfonos chilenos
	patterns := []string{
		`^\+569\d{8}$`, // +56912345678 (móvil con código país)
		`^569\d{8}$`,   // 56912345678 (móvil con código país sin +)
		`^9\d{8}$`,     // 912345678 (móvil)
		`^2\d{8}$`,     // 212345678 (fijo Santiago - 9 dígitos)
		`^[3-7]\d{7}$`, // 12345678 (fijo regiones - códigos área 3-7)
		`^4[1-9]\d{6}$`, // 41234567 (fijo regiones específicas)
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, cleanPhone); matched {
			return true
		}
	}

	return false
}

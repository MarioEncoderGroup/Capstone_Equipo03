package validation_test

import (
	"testing"
	
	"github.com/JoseLuis21/mv-backend/internal/shared/utils"
)

func TestValidateRUT(t *testing.T) {
	tests := []struct {
		name     string
		rut      string
		expected bool
	}{
		// Casos válidos
		{"RUT válido formato completo", "12.345.678-5", true},
		{"RUT válido sin puntos", "12345678-5", true},
		{"RUT válido con K", "10.600.000-K", true},
		{"RUT válido con 0", "10.900.000-0", true},
		{"RUT válido 7 dígitos", "1.234.567-4", true},
		
		// Casos inválidos
		{"RUT vacío", "", false},
		{"RUT sin guión", "12345678", false},
		{"RUT con dígito inválido", "12.345.678-X", false},
		{"RUT muy corto", "123-4", false},
		{"RUT muy largo", "123.456.789.012-3", false},
		{"RUT con letras en número", "12.A45.678-5", false},
		
		// Casos edge
		{"RUT solo espacios", "   ", false},
		{"RUT con formato incorrecto", "12-345.678-5", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ValidateRUT(tt.rut)
			if result != tt.expected {
				t.Errorf("ValidateRUT(%q) = %v, want %v", tt.rut, result, tt.expected)
			}
		})
	}
}

func TestFormatRUT(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"RUT sin formato", "12345678-5", "12.345.678-5"},
		{"RUT ya formateado", "12.345.678-5", "12.345.678-5"},
		{"RUT con espacios", "12 345 678-5", "12.345.678-5"},
		{"RUT 7 dígitos", "1234567-4", "1.234.567-4"},
		{"RUT con K", "10600000-K", "10.600.000-K"},
		{"RUT inválido", "invalid", "invalid"}, // No formatea si es inválido
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.FormatRUT(tt.input)
			if result != tt.expected {
				t.Errorf("FormatRUT(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateChileanPhone(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		// Móviles válidos
		{"Móvil con código país completo", "+56912345678", true},
		{"Móvil con código país sin +", "56912345678", true},
		{"Móvil sin código país", "912345678", true},
		
		// Fijos válidos
		{"Fijo Santiago", "212345678", true},
		{"Fijo regiones", "32345678", true},
		
		// Casos vacíos (permitidos)
		{"Teléfono vacío", "", true},
		
		// Inválidos
		{"Muy corto", "123", false},
		{"Muy largo", "12345678901234", false},
		{"Con letras", "91234567a", false},
		{"Formato incorrecto", "123-456-789", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ValidateChileanPhone(tt.phone)
			if result != tt.expected {
				t.Errorf("ValidateChileanPhone(%q) = %v, want %v", tt.phone, result, tt.expected)
			}
		})
	}
}

func TestValidateChileanEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		// Válidos
		{"Email .cl", "usuario@empresa.cl", true},
		{"Email .com", "test@example.com", true},
		{"Email con subdominios", "user@mail.empresa.cl", true},
		
		// Inválidos
		{"Email vacío", "", false},
		{"Sin @", "usuario.cl", false},
		{"Sin dominio", "usuario@", false},
		{"Muy largo", "a" + string(make([]byte, 150)) + "@test.cl", false},
		{"Formato incorrecto", "@empresa.cl", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ValidateChileanEmail(tt.email)
			if result != tt.expected {
				t.Errorf("ValidateChileanEmail(%q) = %v, want %v", tt.email, result, tt.expected)
			}
		})
	}
}

// Benchmark para validación de RUT (performance testing)
func BenchmarkValidateRUT(b *testing.B) {
	testRUT := "12.345.678-5"
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		utils.ValidateRUT(testRUT)
	}
}

func BenchmarkFormatRUT(b *testing.B) {
	testRUT := "12345678-5"
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		utils.FormatRUT(testRUT)
	}
}
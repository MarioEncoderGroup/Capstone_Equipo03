package hasher

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Service implementa el servicio genérico de hasheado de contraseñas
// Utiliza bcrypt como algoritmo de hasheado por defecto
type Service struct {
	cost int
}

// NewService crea una nueva instancia del servicio de hasher
func NewService() *Service {
	return &Service{
		cost: bcrypt.DefaultCost, // Costo por defecto de bcrypt (10)
	}
}

// NewServiceWithCost crea una nueva instancia con costo personalizado
func NewServiceWithCost(cost int) *Service {
	return &Service{
		cost: cost,
	}
}

// Hash hashea una contraseña usando bcrypt
func (s *Service) Hash(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	return string(hashedBytes), nil
}

// Verify verifica una contraseña contra su hash
func (s *Service) Verify(hashedPassword, password string) error {
	if hashedPassword == "" {
		return fmt.Errorf("hashed password cannot be empty")
	}
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password")
	}

	return nil
}

// GetCost retorna el costo actual de bcrypt
func (s *Service) GetCost() int {
	return s.cost
}

// SetCost establece un nuevo costo de bcrypt (para testing)
func (s *Service) SetCost(cost int) {
	s.cost = cost
}
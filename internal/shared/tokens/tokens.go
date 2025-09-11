package tokens

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/JoseLuis21/mv-backend/internal/shared/utils"
	"github.com/google/uuid"
)

// Service interface genérica para generación y validación de tokens
// Puede ser utilizada por cualquier dominio en el sistema
type Service interface {
	// Tokens aleatorios para verificación
	GenerateRandomToken() (string, error)
	
	// JWT tokens
	GenerateJWT(claims map[string]interface{}, expiryDuration time.Duration) (string, error)
	ValidateJWT(tokenString string) (jwt.MapClaims, error)
	
	// Tokens específicos con formatos predefinidos
	GenerateEmailVerificationToken() (string, error)
	GeneratePasswordResetToken() (string, error)
	GenerateAPIToken() (string, error)
	
	// Refresh tokens para tenant selection
	GenerateRefreshToken(userID uuid.UUID, tenantID uuid.UUID, expiryDuration time.Duration) (string, error)
}

// ServiceImpl implementa el servicio genérico de tokens
type ServiceImpl struct {
	jwtSecret []byte
}

// NewService crea una nueva instancia del servicio genérico de tokens
func NewService() Service {
	// Obtener secret del JWT desde variables de entorno
	jwtSecret := utils.GetEnvOrDefault("JWT_SECRET", "default-secret-key-change-in-production")
	
	return &ServiceImpl{
		jwtSecret: []byte(jwtSecret),
	}
}

// GenerateRandomToken genera un token aleatorio de 32 bytes (64 caracteres hex)
func (ts *ServiceImpl) GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	return hex.EncodeToString(bytes), nil
}

// GenerateEmailVerificationToken genera un token específico para verificación de email
func (ts *ServiceImpl) GenerateEmailVerificationToken() (string, error) {
	return ts.GenerateRandomToken()
}

// GeneratePasswordResetToken genera un token específico para reset de contraseña
func (ts *ServiceImpl) GeneratePasswordResetToken() (string, error) {
	return ts.GenerateRandomToken()
}

// GenerateAPIToken genera un token más largo para APIs (48 bytes = 96 caracteres hex)
func (ts *ServiceImpl) GenerateAPIToken() (string, error) {
	bytes := make([]byte, 48)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	return hex.EncodeToString(bytes), nil
}

// GenerateJWT genera un token JWT con claims personalizados
func (ts *ServiceImpl) GenerateJWT(claims map[string]interface{}, expiryDuration time.Duration) (string, error) {
	// Crear claims base
	jwtClaims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(expiryDuration).Unix(),
		"nbf": time.Now().Unix(),
		"iss": "misviaticos-backend",
	}
	
	// Agregar claims personalizados
	for key, value := range claims {
		jwtClaims[key] = value
	}
	
	// Crear token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	
	// Firmar token
	tokenString, err := token.SignedString(ts.jwtSecret)
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// ValidateJWT valida un token JWT y retorna los claims
func (ts *ServiceImpl) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inválido")
		}
		return ts.jwtSecret, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("token inválido")
}

// TokenType define los tipos de tokens disponibles
type TokenType int

const (
	TokenTypeEmailVerification TokenType = iota
	TokenTypePasswordReset
	TokenTypeAPI
	TokenTypeGeneric
)

// GenerateTypedToken genera un token de un tipo específico
func (ts *ServiceImpl) GenerateTypedToken(tokenType TokenType) (string, error) {
	switch tokenType {
	case TokenTypeEmailVerification:
		return ts.GenerateEmailVerificationToken()
	case TokenTypePasswordReset:
		return ts.GeneratePasswordResetToken()
	case TokenTypeAPI:
		return ts.GenerateAPIToken()
	case TokenTypeGeneric:
		return ts.GenerateRandomToken()
	default:
		return ts.GenerateRandomToken()
	}
}

// JWTConfig configuración para generación de JWTs
type JWTConfig struct {
	ExpiryHours   int
	Subject       string
	Audience      string
	CustomClaims  map[string]interface{}
}

// GenerateJWTWithConfig genera un JWT con configuración específica
func (ts *ServiceImpl) GenerateJWTWithConfig(config JWTConfig) (string, error) {
	if config.ExpiryHours <= 0 {
		config.ExpiryHours = 24 // Default: 24 horas
	}
	
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Duration(config.ExpiryHours) * time.Hour).Unix(),
		"nbf": time.Now().Unix(),
		"iss": "misviaticos-backend",
	}
	
	if config.Subject != "" {
		claims["sub"] = config.Subject
	}
	
	if config.Audience != "" {
		claims["aud"] = config.Audience
	}
	
	// Agregar claims personalizados
	for key, value := range config.CustomClaims {
		claims[key] = value
	}
	
	// Crear y firmar token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	return token.SignedString(ts.jwtSecret)
}

// GenerateRefreshToken genera un refresh token con información del usuario y tenant
func (ts *ServiceImpl) GenerateRefreshToken(userID uuid.UUID, tenantID uuid.UUID, expiryDuration time.Duration) (string, error) {
	claims := map[string]interface{}{
		"user_id":   userID.String(),
		"tenant_id": tenantID.String(),
		"type":      "refresh",
		"purpose":   "tenant_refresh",
	}
	
	return ts.GenerateJWT(claims, expiryDuration)
}
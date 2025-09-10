package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/JoseLuis21/mv-backend/internal/core/auth/ports"
	"github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	"github.com/JoseLuis21/mv-backend/internal/shared/utils"
)

// TokenGeneratorImpl implementa la generación de tokens
type TokenGeneratorImpl struct {
	jwtSecret []byte
}

// NewTokenGenerator crea una nueva instancia del generador de tokens
func NewTokenGenerator() ports.TokenGenerator {
	// Obtener secret del JWT desde variables de entorno
	jwtSecret := utils.GetEnvOrDefault("JWT_SECRET", "default-secret-key-change-in-production")
	
	return &TokenGeneratorImpl{
		jwtSecret: []byte(jwtSecret),
	}
}

// GenerateEmailVerificationToken genera un token aleatorio para verificación de email
func (tg *TokenGeneratorImpl) GenerateEmailVerificationToken() (string, error) {
	// Generar 32 bytes aleatorios (64 caracteres hex)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	return hex.EncodeToString(bytes), nil
}

// GeneratePasswordResetToken genera un token aleatorio para reset de contraseña
func (tg *TokenGeneratorImpl) GeneratePasswordResetToken() (string, error) {
	// Reutilizar la misma lógica que email verification
	return tg.GenerateEmailVerificationToken()
}

// Claims representa los claims del JWT
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT genera un token JWT para autenticación
func (tg *TokenGeneratorImpl) GenerateJWT(user *domain.User) (string, error) {
	if user == nil {
		return "", errors.New("usuario requerido")
	}

	// Verificar que el usuario esté activo y verificado
	if !user.IsActive || !user.EmailVerified {
		return "", errors.New("usuario debe estar activo y verificado")
	}

	// Obtener duración del token desde env (default: 24 horas)
	tokenDuration := utils.GetEnvOrDefault("JWT_EXPIRY_HOURS", "24")
	var expiryDuration time.Duration
	
	switch tokenDuration {
	case "1":
		expiryDuration = 1 * time.Hour
	case "8":
		expiryDuration = 8 * time.Hour
	case "24":
		expiryDuration = 24 * time.Hour
	case "168": // 7 días
		expiryDuration = 168 * time.Hour
	default:
		expiryDuration = 24 * time.Hour
	}

	// Crear claims
	claims := Claims{
		UserID:   user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "misviaticos-backend",
			Subject:   user.ID.String(),
		},
	}

	// Crear token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar token
	tokenString, err := token.SignedString(tg.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT valida un token JWT y retorna los claims
func (tg *TokenGeneratorImpl) ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inválido")
		}
		return tg.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}
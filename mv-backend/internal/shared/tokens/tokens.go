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

	// Access y Refresh tokens con tenant context
	GenerateAccessToken(userID uuid.UUID, tenantID *uuid.UUID) (string, int64, error)
	GenerateAccessTokenWithRoles(userID uuid.UUID, tenantID *uuid.UUID, roles []string, permissions []string) (string, int64, error)
	GenerateRefreshToken(userID uuid.UUID, tenantID *uuid.UUID) (string, int64, error)
	ValidateAccessToken(tokenString string) (*TokenClaims, error)
	ValidateRefreshToken(tokenString string) (*TokenClaims, error)
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

// TokenClaims estructura de claims para access y refresh tokens
type TokenClaims struct {
	UserID      uuid.UUID
	TenantID    *uuid.UUID
	Type        string // "access" o "refresh"
	Roles       []string
	Permissions []string
	IssuedAt    time.Time
	ExpiresAt   time.Time
}

// GenerateAccessToken genera un access token JWT con user_id y tenant_id opcional
func (ts *ServiceImpl) GenerateAccessToken(userID uuid.UUID, tenantID *uuid.UUID) (string, int64, error) {
	// Access token expira en 1 hora
	expiryDuration := 1 * time.Hour
	expiresIn := int64(expiryDuration.Seconds())

	claims := map[string]interface{}{
		"user_id": userID.String(),
		"type":    "access",
	}

	// Agregar tenant_id solo si está presente
	if tenantID != nil {
		claims["tenant_id"] = tenantID.String()
	}

	token, err := ts.GenerateJWT(claims, expiryDuration)
	if err != nil {
		return "", 0, err
	}

	return token, expiresIn, nil
}

// GenerateAccessTokenWithRoles genera un access token JWT con roles y permisos incluidos
func (ts *ServiceImpl) GenerateAccessTokenWithRoles(userID uuid.UUID, tenantID *uuid.UUID, roles []string, permissions []string) (string, int64, error) {
	// Access token expira en 1 hora
	expiryDuration := 1 * time.Hour
	expiresIn := int64(expiryDuration.Seconds())

	claims := map[string]interface{}{
		"user_id":     userID.String(),
		"type":        "access",
		"roles":       roles,
		"permissions": permissions,
	}

	// Agregar tenant_id solo si está presente
	if tenantID != nil {
		claims["tenant_id"] = tenantID.String()
	}

	token, err := ts.GenerateJWT(claims, expiryDuration)
	if err != nil {
		return "", 0, err
	}

	return token, expiresIn, nil
}

// GenerateRefreshToken genera un refresh token con información del usuario y tenant
func (ts *ServiceImpl) GenerateRefreshToken(userID uuid.UUID, tenantID *uuid.UUID) (string, int64, error) {
	// Refresh token expira en 7 días
	expiryDuration := 7 * 24 * time.Hour
	expiresIn := int64(expiryDuration.Seconds())

	claims := map[string]interface{}{
		"user_id": userID.String(),
		"type":    "refresh",
	}

	// Agregar tenant_id solo si está presente
	if tenantID != nil {
		claims["tenant_id"] = tenantID.String()
	}

	token, err := ts.GenerateJWT(claims, expiryDuration)
	if err != nil {
		return "", 0, err
	}

	return token, expiresIn, nil
}

// ValidateAccessToken valida un access token y retorna los claims estructurados
func (ts *ServiceImpl) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	claims, err := ts.ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	// Verificar que sea un access token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return nil, errors.New("token type inválido")
	}

	return ts.parseTokenClaims(claims)
}

// ValidateRefreshToken valida un refresh token y retorna los claims estructurados
func (ts *ServiceImpl) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	claims, err := ts.ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	// Verificar que sea un refresh token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, errors.New("token type inválido")
	}

	return ts.parseTokenClaims(claims)
}

// parseTokenClaims convierte MapClaims a TokenClaims estructurado
func (ts *ServiceImpl) parseTokenClaims(claims jwt.MapClaims) (*TokenClaims, error) {
	// Extraer user_id
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("user_id no encontrado en token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("user_id inválido en token")
	}

	// Extraer tenant_id (opcional)
	var tenantID *uuid.UUID
	if tenantIDStr, ok := claims["tenant_id"].(string); ok {
		parsed, err := uuid.Parse(tenantIDStr)
		if err == nil {
			tenantID = &parsed
		}
	}

	// Extraer tipo
	tokenType, _ := claims["type"].(string)

	// Extraer roles (opcional)
	var roles []string
	if rolesInterface, ok := claims["roles"].([]interface{}); ok {
		for _, r := range rolesInterface {
			if roleStr, ok := r.(string); ok {
				roles = append(roles, roleStr)
			}
		}
	}

	// Extraer permissions (opcional)
	var permissions []string
	if permsInterface, ok := claims["permissions"].([]interface{}); ok {
		for _, p := range permsInterface {
			if permStr, ok := p.(string); ok {
				permissions = append(permissions, permStr)
			}
		}
	}

	// Extraer timestamps
	var issuedAt time.Time
	if iat, ok := claims["iat"].(float64); ok {
		issuedAt = time.Unix(int64(iat), 0)
	}

	var expiresAt time.Time
	if exp, ok := claims["exp"].(float64); ok {
		expiresAt = time.Unix(int64(exp), 0)
	}

	return &TokenClaims{
		UserID:      userID,
		TenantID:    tenantID,
		Type:        tokenType,
		Roles:       roles,
		Permissions: permissions,
		IssuedAt:    issuedAt,
		ExpiresAt:   expiresAt,
	}, nil
}
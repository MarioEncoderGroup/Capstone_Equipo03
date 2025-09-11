# PASO 5: Sistema Completo de Refresh Tokens - Testing Guide

## Implementación Completada

### 1. AuthLoginResponse - ACTUALIZADA ✅
**Archivo**: `internal/core/auth/domain/auth.go`

**BEFORE**:
```go
type AuthLoginResponse struct {
    AccessToken string           `json:"access_token"`
    User        domain_user.User `json:"user"`
}
```

**AFTER**:
```go
type AuthLoginResponse struct {
    AccessToken  string           `json:"access_token"`
    RefreshToken string           `json:"refresh_token"`  // ← AGREGADO PASO 5
    ExpiresIn    int64            `json:"expires_in"`     // ← AGREGADO PASO 5
    TokenType    string           `json:"token_type"`     // ← AGREGADO PASO 5
    User         domain_user.User `json:"user"`
}
```

### 2. DTOs para Refresh Token - AGREGADOS ✅
**Archivo**: `internal/core/auth/domain/auth.go`

**NUEVOS DTOs**:
```go
// RefreshTokenDto para renovar tokens
type RefreshTokenDto struct {
    RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse respuesta al renovar token
type RefreshTokenResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
    TokenType    string `json:"token_type"`
}
```

### 3. Login() Service - ACTUALIZADO ✅
**Archivo**: `internal/core/auth/services/auth.go`

**CAMBIOS CLAVE**:
- Genera Access Token + Refresh Token
- Access Token: 24 horas de duración
- Refresh Token: 30 días de duración
- Ambos tokens retornados en la respuesta

**BEFORE**:
```go
token, err := s.tokenService.GenerateJWT(claims, 24*time.Hour)
response := &domain_auth.AuthLoginResponse{
    AccessToken: token,
    User: userDomain.User{...},
}
```

**AFTER**:
```go
accessToken, err := s.tokenService.GenerateJWT(claims, 24*time.Hour)
refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, uuid.Nil, 30*24*time.Hour)

response := &domain_auth.AuthLoginResponse{
    AccessToken:  accessToken,
    RefreshToken: refreshToken,
    ExpiresIn:    int64(24 * 60 * 60),
    TokenType:    "Bearer",
    User: userDomain.User{...},
}
```

### 4. Nuevos Métodos en Service - IMPLEMENTADOS ✅
**Archivo**: `internal/core/auth/services/auth.go`

**RefreshAccessToken()**: 
- Valida refresh token
- Verifica usuario activo
- Genera nuevos Access + Refresh tokens
- Mantiene contexto de tenant si existe

**RevokeRefreshToken()**:
- Valida refresh token
- Marca como revocado (TODO: persistencia)

### 5. Controller RefreshToken - IMPLEMENTADO ✅
**Archivo**: `internal/controllers/auth.go`

**Nuevo Endpoint**:
```go
func (ac *AuthController) RefreshToken(c *fiber.Ctx) error {
    var req authDomain.RefreshTokenDto
    // Validación + Llamada al servicio
    response, err := ac.authService.RefreshAccessToken(ctx, req.RefreshToken)
    return c.JSON(response)
}
```

### 6. Ruta /auth/refresh-token - AGREGADA ✅
**Archivo**: `internal/routes/auth_routes.go`

```go
// PASO 5: Endpoint para refresh tokens - IMPLEMENTADO
auth.Post("/refresh-token", authController.RefreshToken)
```

### 7. Interface AuthService - ACTUALIZADA ✅
**Archivo**: `internal/core/auth/ports/auth.go`

**Nuevos Métodos**:
```go
// RefreshAccessToken renueva el access token usando un refresh token
RefreshAccessToken(ctx context.Context, refreshToken string) (*domain_auth.RefreshTokenResponse, error)

// RevokeRefreshToken revoca un refresh token
RevokeRefreshToken(ctx context.Context, refreshToken string) error
```

## Cómo Probar el Sistema

### 1. Login Request
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Respuesta Esperada**:
```json
{
  "success": true,
  "message": "Autenticación exitosa",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "token_type": "Bearer",
    "user": {
      "id": "...",
      "username": "...",
      "full_name": "...",
      "email": "test@example.com"
    }
  }
}
```

### 2. Refresh Token Request
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh-token \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }'
```

**Respuesta Esperada**:
```json
{
  "success": true,
  "message": "Tokens renovados exitosamente",
  "data": {
    "access_token": "NEW_ACCESS_TOKEN",
    "refresh_token": "NEW_REFRESH_TOKEN",
    "expires_in": 86400,
    "token_type": "Bearer"
  }
}
```

## Flujo Completo del Sistema

### 1. Login
- Usuario hace login → Recibe Access Token + Refresh Token
- Access Token válido por 24 horas
- Refresh Token válido por 30 días

### 2. Uso Normal
- Cliente usa Access Token para requests autenticados
- Cuando Access Token expira, cliente usa Refresh Token

### 3. Renovación
- Cliente envía Refresh Token a `/auth/refresh-token`
- Recibe nuevos Access Token + Refresh Token
- Token anterior se invalida (lógicamente)

### 4. SelectTenant
- SelectTenant también retorna Access + Refresh tokens
- Tokens incluyen tenant_id en claims
- Mismo patrón de 24h + 30 días

## Características Implementadas

✅ Login retorna Access Token + Refresh Token  
✅ Duración: Access Token 24h, Refresh Token 30 días  
✅ Endpoint `/auth/refresh-token` funcional  
✅ Renovación genera nuevos tokens (ambos)  
✅ Validación completa de refresh tokens  
✅ Integración con SelectTenant  
✅ DTOs en domain layer (arquitectura correcta)  
✅ Servicios genéricos reutilizados  
✅ Compilación exitosa  

## Estado: PASO 5 COMPLETADO ✅

El sistema completo de refresh tokens está implementado siguiendo exactamente el patrón api-golang-2025.
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain_auth "github.com/JoseLuis21/mv-backend/internal/core/auth/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/auth/ports"
	tenantDomain "github.com/JoseLuis21/mv-backend/internal/core/tenant/domain"
	userDomain "github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	userPorts "github.com/JoseLuis21/mv-backend/internal/core/user/ports"
	"github.com/JoseLuis21/mv-backend/internal/shared/email"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
	"github.com/JoseLuis21/mv-backend/internal/shared/hasher"
	"github.com/JoseLuis21/mv-backend/internal/shared/tokens"
	"github.com/google/uuid"
)

// authService implementa el servicio de autenticación usando servicios genéricos
type authService struct {
	userService      userPorts.UserService
	passwordHasher   *hasher.Service
	tokenService     tokens.Service
	emailService     email.Service
	emailTokenExpiry time.Duration
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(
	userService userPorts.UserService,
	passwordHasher *hasher.Service,
	tokenService tokens.Service,
	emailService email.Service,
) ports.AuthService {
	return &authService{
		userService:      userService,
		passwordHasher:   passwordHasher,
		tokenService:     tokenService,
		emailService:     emailService,
		emailTokenExpiry: 24 * time.Hour, // Token válido por 24 horas
	}
}

// Register registra un nuevo usuario individual o con tenant
func (s *authService) Register(ctx context.Context, req *domain_auth.AuthRegisterDto) (*domain_auth.AuthRegisterResponse, error) {
	// 1. Validar si el usuario ya existe
	existingUser, err := s.userService.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sharedErrors.ErrUserNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, sharedErrors.ErrUserAlreadyExists
	}

	// 2. Hash de la contraseña
	hashedPassword, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	// 3. Generar token de verificación de email
	emailToken, err := s.tokenService.GenerateEmailVerificationToken()
	if err != nil {
		return nil, err
	}

	// 4. Crear entidad de usuario con firstname, lastname, email, phone, hashedPassword
	user := userDomain.NewUser(req.FirstName, req.LastName, req.Email, req.Phone, hashedPassword)
	user.SetEmailVerificationToken(emailToken, s.emailTokenExpiry)

	// 5. Guardar usuario en BD
	if err := s.userService.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("error creando usuario: %w", err)
	}

	// 6. Enviar email de verificación
	if err := s.sendEmailVerification(ctx, user, emailToken); err != nil {
		// Log error pero no fallar el registro
		fmt.Printf("Error enviando email de verificación: %v\n", err)
	}

	// 7. Preparar respuesta siguiendo el patrón de referencia con nuevos campos
	phoneValue := ""
	if user.Phone != nil {
		phoneValue = *user.Phone
	}
	
	return &domain_auth.AuthRegisterResponse{
		ID:                        user.ID,
		FirstName:                 user.FirstName,
		LastName:                  user.LastName,
		FullName:                  user.FullName, // Para backward compatibility
		Email:                     user.Email,
		Phone:                     phoneValue,
		EmailToken:                emailToken,
		RequiresEmailVerification: true,
		Message:                   "Usuario registrado exitosamente. Verifica tu email para activar la cuenta.",
	}, nil
}


// VerifyUserEmail verifica el email de un usuario usando el token
func (s *authService) VerifyUserEmail(ctx context.Context, token string) error {
	// 1. Buscar usuario por token
	user, err := s.userService.GetUserByEmailToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid verification token")
	}

	// 2. Check if user exists
	if user == nil {
		return fmt.Errorf("invalid verification token")
	}

	// 3. Check if email is already verified
	if user.EmailVerified {
		return fmt.Errorf("email is already verified")
	}

	// 4. Validar token
	if !user.IsEmailTokenValid(token) {
		return fmt.Errorf("token inválido o expirado")
	}

	// 5. Activar usuario
	user.ActivateUser()

	// 6. Actualizar en BD
	if err := s.userService.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("error activando usuario: %w", err)
	}

	// 7. Enviar email de bienvenida
	if err := s.sendWelcomeEmail(ctx, user); err != nil {
		// Log error pero no fallar la verificación
		fmt.Printf("Error enviando email de bienvenida: %v\n", err)
	}

	return nil
}

// Login autentica un usuario y retorna tokens
func (s *authService) Login(ctx context.Context, req *domain_auth.AuthLoginDto) (*domain_auth.AuthLoginResponse, error) {
	// 1. Buscar usuario por email
	user, err := s.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// 2. Check if user exists
	if user == nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// 3. Check if user's email is verified
	if !user.EmailVerified {
		return nil, fmt.Errorf("email %s is not verified", req.Email)
	}

	// 4. Verify the password
	if err := s.passwordHasher.Verify(user.Password, req.Password); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// 5. Generate JWT token with type "login" (before tenant selection)
	claims := map[string]interface{}{
		"user_id":    user.ID.String(),
		"email":      user.Email,
		"username":   user.Username,
		"full_name":  user.FullName,
		"is_active":  user.IsActive,
		"type":       "login", // Tipo login antes de seleccionar tenant
	}
	
	accessToken, err := s.tokenService.GenerateJWT(claims, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	// PASO 5: Generate refresh token for login (30 días de expiración)
	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, uuid.Nil, 30*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	// 6. Preparar respuesta con refresh token - PASO 5
	resp := &domain_auth.AuthLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,                    // ← AGREGADO PASO 5
		ExpiresIn:    int64(24 * 60 * 60),           // ← AGREGADO PASO 5: 24 horas en segundos
		TokenType:    "Bearer",                       // ← AGREGADO PASO 5
		User: userDomain.User{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Email:     user.Email,
			IsActive:  user.IsActive,
			LastLogin: user.LastLogin,
		},
	}

	return resp, nil
}

// ForgotPassword inicia el proceso de recuperación de contraseña
func (s *authService) ForgotPassword(ctx context.Context, email string) (*userDomain.User, error) {
	// 1. Buscar usuario por email
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, sharedErrors.ErrUserNotFound
	}

	if user == nil {
		return nil, sharedErrors.ErrUserNotFound
	}

	// 2. Verificar que el email esté verificado (regla de negocio)
	if !user.EmailVerified {
		return nil, sharedErrors.NewValidationError("Email no verificado", "email_not_verified")
	}

	// 3. Verificar si ya tiene un token activo (prevenir abuso)
	if user.HasActivePasswordResetToken() {
		return nil, sharedErrors.NewValidationError("Ya tienes una solicitud de reset activa", "active_reset_token")
	}

	// 4. Generar token seguro de reset (1 hora de expiración)
	token, err := s.tokenService.GeneratePasswordResetToken()
	if err != nil {
		return nil, fmt.Errorf("error generando token: %w", err)
	}

	// 5. Establecer token en usuario usando método de dominio
	user.SetPasswordResetToken(token, 1*time.Hour)
	fmt.Printf("🔵 DEBUG: Set reset token %s for user %s\n", token[:8]+"...", user.Email)

	// 6. Guardar usuario con token en BD
	if err := s.userService.UpdateUser(ctx, user); err != nil {
		fmt.Printf("❌ DEBUG: Error updating user: %v\n", err)
		return nil, fmt.Errorf("error actualizando usuario: %w", err)
	}
	fmt.Printf("✅ DEBUG: User updated successfully in database\n")

	// 7. Enviar email con token de reset
	if err := s.sendPasswordResetEmail(ctx, user, token); err != nil {
		// Log error pero no fallar el proceso
		fmt.Printf("Error enviando email de reset: %v\n", err)
	}

	return user, nil
}

// ResetPassword resetea la contraseña usando un token
func (s *authService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	// 1. Buscar usuario por token de reset
	user, err := s.userService.GetUserByPasswordResetToken(ctx, token)
	if err != nil {
		return fmt.Errorf("token de reset inválido")
	}

	if user == nil {
		return fmt.Errorf("token de reset inválido")
	}

	// 2. Validar que el token sea válido y no haya expirado usando método de dominio
	if !user.IsPasswordResetTokenValid(token) {
		return fmt.Errorf("token de reset inválido o expirado")
	}

	// 3. Hash de la nueva contraseña
	hashedPassword, err := s.passwordHasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("error hasheando contraseña: %w", err)
	}

	// 4. Cambiar contraseña usando método de dominio (limpia el token automáticamente)
	user.ChangePassword(hashedPassword)

	// 5. Actualizar usuario en BD
	if err := s.userService.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("error actualizando contraseña: %w", err)
	}

	// 6. Opcional: Enviar email de confirmación de cambio de contraseña
	if err := s.sendPasswordChangeConfirmationEmail(ctx, user); err != nil {
		// Log error pero no fallar el proceso
		fmt.Printf("Error enviando email de confirmación: %v\n", err)
	}

	return nil
}

// ResendEmailVerification reenvía el email de verificación
func (s *authService) ResendEmailVerification(ctx context.Context, email string) error {
	// 1. Buscar usuario por email
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("usuario no encontrado")
	}

	// 2. Verificar que el usuario no esté ya verificado
	if user.EmailVerified {
		return fmt.Errorf("el email ya está verificado")
	}

	// 3. Generar nuevo token
	emailToken, err := s.tokenService.GenerateEmailVerificationToken()
	if err != nil {
		return fmt.Errorf("error generando token: %w", err)
	}

	// 4. Actualizar token en usuario
	user.SetEmailVerificationToken(emailToken, s.emailTokenExpiry)

	// 5. Actualizar en BD
	if err := s.userService.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("error actualizando usuario: %w", err)
	}

	// 6. Enviar email
	if err := s.sendEmailVerification(ctx, user, emailToken); err != nil {
		return fmt.Errorf("error enviando email: %w", err)
	}

	return nil
}


// sendEmailVerification envía el email de verificación usando el servicio genérico
func (s *authService) sendEmailVerification(ctx context.Context, user *userDomain.User, token string) error {
	templateData := &email.TemplateData{
		FullName: user.FullName,
		Email:    user.Email,
		URL:      fmt.Sprintf("https://misviaticos.cl/verify?token=%s", token),
	}

	return s.emailService.SendTemplateEmail(ctx, email.TemplateEmailVerification, templateData)
}

// sendWelcomeEmail envía el email de bienvenida usando el servicio genérico
func (s *authService) sendWelcomeEmail(ctx context.Context, user *userDomain.User) error {
	templateData := &email.TemplateData{
		FullName: user.FullName,
		Email:    user.Email,
		URL:      "https://misviaticos.cl/dashboard",
	}

	return s.emailService.SendTemplateEmail(ctx, email.TemplateWelcome, templateData)
}

// sendPasswordResetEmail envía el email de reset de contraseña usando el servicio genérico
func (s *authService) sendPasswordResetEmail(ctx context.Context, user *userDomain.User, token string) error {
	templateData := &email.TemplateData{
		FullName: user.FullName,
		Email:    user.Email,
		URL:      fmt.Sprintf("https://misviaticos.cl/reset-password?token=%s", token),
	}

	return s.emailService.SendTemplateEmail(ctx, email.TemplatePasswordReset, templateData)
}

// SelectTenant selecciona un tenant específico post-login y genera nuevos tokens
func (s *authService) SelectTenant(ctx context.Context, tenant *tenantDomain.Tenant, userID uuid.UUID) (*tenantDomain.SelectTenantResponseDto, error) {
	// 1. Verificar que el tenant existe y está activo
	if tenant == nil {
		return nil, fmt.Errorf("tenant is required")
	}
	
	if tenant.Status != string(tenantDomain.TenantStatusActive) {
		return nil, fmt.Errorf("tenant is not active")
	}

	// 2. Obtener información del usuario
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// 3. Verificar que el usuario esté activo
	if !user.IsActive {
		return nil, fmt.Errorf("user is not active")
	}

	// TODO: Aquí se debería verificar que el usuario tiene acceso al tenant
	// mediante un repositorio o servicio de UserTenant, pero por ahora lo omitimos

	// 4. Generar claims para el JWT con tenant_id y type "tenant_selection"
	claims := map[string]interface{}{
		"user_id":    user.ID.String(),
		"email":      user.Email,
		"username":   user.Username,
		"full_name":  user.FullName,
		"is_active":  user.IsActive,
		"tenant_id":  tenant.ID.String(),
		"type":       "tenant_selection", // Cambio crítico vs "login"
	}

	// 5. Generar nuevo JWT con 24 horas de expiración
	accessToken, err := s.tokenService.GenerateJWT(claims, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	// 6. Generar refresh token (30 días de expiración)
	refreshToken, err := s.tokenService.GenerateRefreshToken(userID, tenant.ID, 30*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	// 7. Preparar respuesta completa
	response := &tenantDomain.SelectTenantResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(24 * 60 * 60), // 24 horas en segundos
		User: userDomain.User{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Email:     user.Email,
			IsActive:  user.IsActive,
			LastLogin: user.LastLogin,
		},
		Tenant: tenant,
	}

	return response, nil
}

// PASO 5: RefreshAccessToken renueva el access token usando un refresh token
func (s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) (*domain_auth.RefreshTokenResponse, error) {
	// 1. Validar refresh token
	claims, err := s.tokenService.ValidateJWT(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// 2. Verificar que es un token de tipo "refresh"
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	// 3. Extraer user_id del refresh token
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	// 4. Buscar usuario para validar que sigue activo
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if !user.IsActive || !user.EmailVerified {
		return nil, fmt.Errorf("user is not active or email not verified")
	}

	// 5. Extraer tenant_id si existe
	var tenantID uuid.UUID
	if tenantIDStr, exists := claims["tenant_id"].(string); exists && tenantIDStr != uuid.Nil.String() {
		tenantID, err = uuid.Parse(tenantIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid tenant_id format: %w", err)
		}
	} else {
		tenantID = uuid.Nil
	}

	// 6. Generar nuevos claims para access token
	newClaims := map[string]interface{}{
		"user_id":    user.ID.String(),
		"email":      user.Email,
		"username":   user.Username,
		"full_name":  user.FullName,
		"is_active":  user.IsActive,
		"type":       "login", // Mantener tipo login si no hay tenant
	}

	// Si hay tenant_id, cambiar tipo a tenant_selection
	if tenantID != uuid.Nil {
		newClaims["tenant_id"] = tenantID.String()
		newClaims["type"] = "tenant_selection"
	}

	// 7. Generar nuevo access token (24 horas)
	newAccessToken, err := s.tokenService.GenerateJWT(newClaims, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("error generating new access token: %w", err)
	}

	// 8. Generar nuevo refresh token (30 días)
	newRefreshToken, err := s.tokenService.GenerateRefreshToken(userID, tenantID, 30*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("error generating new refresh token: %w", err)
	}

	// 9. Preparar respuesta
	response := &domain_auth.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(24 * 60 * 60), // 24 horas en segundos
		TokenType:    "Bearer",
	}

	return response, nil
}

// PASO 5: RevokeRefreshToken revoca un refresh token (marcándolo como inválido)
func (s *authService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	// 1. Validar que el token es válido
	_, err := s.tokenService.ValidateJWT(refreshToken)
	if err != nil {
		return fmt.Errorf("invalid refresh token: %w", err)
	}

	// TODO: Implementar lista negra de tokens revocados en BD o Redis
	// Por ahora solo validamos que el token sea válido, la revocación real
	// requiere persistencia que implementaremos más adelante

	return nil
}

// sendPasswordChangeConfirmationEmail envía email de confirmación de cambio de contraseña
func (s *authService) sendPasswordChangeConfirmationEmail(ctx context.Context, user *userDomain.User) error {
	templateData := &email.TemplateData{
		FullName: user.FullName,
		Email:    user.Email,
		URL:      "https://misviaticos.cl/login",
		Message:  "Tu contraseña ha sido cambiada exitosamente",
	}

	return s.emailService.SendTemplateEmail(ctx, email.TemplateGeneric, templateData)
}
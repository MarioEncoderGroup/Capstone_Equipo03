package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/core/auth/ports"
	"github.com/JoseLuis21/mv-backend/internal/core/tenant/domain"
	tenantPorts "github.com/JoseLuis21/mv-backend/internal/core/tenant/ports"
	userDomain "github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	userPorts "github.com/JoseLuis21/mv-backend/internal/core/user/ports"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
	"github.com/JoseLuis21/mv-backend/internal/shared/utils"
	"github.com/google/uuid"
)

// RegisterServiceImpl implementa el servicio de registro
// Orquesta las operaciones de registro de usuarios y tenants
type RegisterServiceImpl struct {
	userRepo         userPorts.UserRepository
	tenantRepo       tenantPorts.TenantRepository
	passwordHasher   ports.PasswordHasher
	tokenGenerator   ports.TokenGenerator
	emailService     ports.EmailService
	emailTokenExpiry time.Duration
}

// NewRegisterService crea una nueva instancia del servicio de registro
func NewRegisterService(
	userRepo userPorts.UserRepository,
	tenantRepo tenantPorts.TenantRepository,
	passwordHasher ports.PasswordHasher,
	tokenGenerator ports.TokenGenerator,
	emailService ports.EmailService,
) ports.RegisterService {
	return &RegisterServiceImpl{
		userRepo:         userRepo,
		tenantRepo:       tenantRepo,
		passwordHasher:   passwordHasher,
		tokenGenerator:   tokenGenerator,
		emailService:     emailService,
		emailTokenExpiry: 24 * time.Hour, // Token válido por 24 horas
	}
}

// RegisterUser registra un nuevo usuario individual
func (rs *RegisterServiceImpl) RegisterUser(ctx context.Context, req *ports.RegisterRequest) (*ports.RegisterResponse, error) {
	// 1. Validar datos de entrada
	if err := rs.ValidateRegistrationData(req); err != nil {
		return nil, fmt.Errorf("datos de registro inválidos: %w", err)
	}

	// 2. Verificar que el usuario no exista
	exists, err := rs.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, sharedErrors.WrapError(sharedErrors.ErrDatabaseConnection, fmt.Sprintf("verificando email: %v", err))
	}
	if exists {
		return nil, sharedErrors.ErrUserAlreadyExists
	}

	exists, err = rs.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, sharedErrors.WrapError(sharedErrors.ErrDatabaseConnection, fmt.Sprintf("verificando username: %v", err))
	}
	if exists {
		return nil, sharedErrors.ErrUsernameAlreadyExists
	}

	// 3. Hash de la contraseña
	hashedPassword, err := rs.passwordHasher.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error generando hash de contraseña: %w", err)
	}

	// 4. Crear entidad de usuario
	user := userDomain.NewUser(req.Username, req.FullName, req.Email, hashedPassword)

	// Campos opcionales
	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	if req.IdentificationNumber != "" {
		user.IdentificationNumber = &req.IdentificationNumber
	}

	// 5. Generar token de verificación de email
	emailToken, err := rs.tokenGenerator.GenerateEmailVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("error generando token de email: %w", err)
	}
	user.SetEmailVerificationToken(emailToken, rs.emailTokenExpiry)

	// 6. Guardar usuario en BD
	if err := rs.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("error creando usuario: %w", err)
	}

	// 7. Enviar email de verificación
	if err := rs.emailService.SendEmailVerification(ctx, user, emailToken); err != nil {
		// Log error pero no fallar el registro
		// TODO: Implementar logging
		fmt.Printf("Error enviando email de verificación: %v\n", err)
	}

	return &ports.RegisterResponse{
		UserID:                    user.ID,
		Message:                   "Usuario registrado exitosamente. Verifica tu email para activar la cuenta.",
		RequiresEmailVerification: true,
	}, nil
}

// RegisterUserWithTenant registra un usuario y crea un tenant/empresa
func (rs *RegisterServiceImpl) RegisterUserWithTenant(ctx context.Context, req *ports.RegisterRequest) (*ports.RegisterResponse, error) {
	// 1. Validar datos incluyendo tenant
	if err := rs.ValidateRegistrationData(req); err != nil {
		return nil, fmt.Errorf("datos de registro inválidos: %w", err)
	}

	if req.TenantData == nil {
		return nil, sharedErrors.ErrTenantDataRequired
	}

	// 2. Validar RUT del tenant
	if !utils.ValidateRUT(req.TenantData.RUT) {
		return nil, sharedErrors.WrapError(sharedErrors.ErrInvalidRUT, "RUT del negocio inválido")
	}

	// 3. Verificar que el tenant no exista
	formattedRUT := utils.FormatRUT(req.TenantData.RUT)
	exists, err := rs.tenantRepo.ExistsByRUT(ctx, formattedRUT)
	if err != nil {
		return nil, fmt.Errorf("error verificando RUT: %w", err)
	}
	if exists {
		return nil, sharedErrors.ErrTenantAlreadyExists
	}

	// 4. Registrar usuario primero
	userResponse, err := rs.RegisterUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error registrando usuario: %w", err)
	}

	// Obtener el siguiente número de nodo
	nodeNumber, err := rs.tenantRepo.GetNextNodeNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo número de nodo: %w", err)
	}

	// Generar slug del tenant (puedes usar esta función helper)
	slug := utils.GenerateSlug(req.TenantData.BusinessName)

	// 5. Crear entidad tenant
	tenant := domain.NewTenant(
		formattedRUT,
		req.TenantData.BusinessName,
		req.TenantData.Email,
		req.TenantData.Phone,
		req.TenantData.Address,
		req.TenantData.Website,
		req.TenantData.RegionID,
		req.TenantData.CommuneID,
		req.TenantData.CountryID,
		nodeNumber,
		slug,
		userResponse.UserID,
	)

	// 6. Guardar tenant en BD
	if err := rs.tenantRepo.Create(ctx, tenant); err != nil {
		// TODO: Implementar rollback del usuario
		return nil, fmt.Errorf("error creando tenant: %w", err)
	}

	// 7. Crear relación usuario-tenant
	tenantUser := userDomain.NewTenantUser(tenant.ID, userResponse.UserID)
	if err := rs.userRepo.AddUserToTenant(ctx, tenantUser); err != nil {
		// TODO: Implementar rollback
		return nil, fmt.Errorf("error asociando usuario al tenant: %w", err)
	}

	// 8. Crear base de datos del tenant
	if err := rs.tenantRepo.CreateTenantDatabase(ctx, tenant.TenantName); err != nil {
		// TODO: Implementar rollback
		return nil, fmt.Errorf("error creando base de datos del tenant: %w", err)
	}

	userResponse.TenantID = &tenant.ID
	userResponse.Message = "Empresa y usuario registrados exitosamente. Verifica tu email para activar la cuenta."

	return userResponse, nil
}

// VerifyEmail verifica el email de un usuario usando el token
func (rs *RegisterServiceImpl) VerifyEmail(ctx context.Context, token string) error {
	// 1. Buscar usuario por token
	user, err := rs.userRepo.GetByEmailToken(ctx, token)
	if err != nil {
		return fmt.Errorf("token inválido o expirado")
	}

	// 2. Validar token
	if !user.IsEmailTokenValid(token) {
		return errors.New("token inválido o expirado")
	}

	// 3. Activar usuario
	user.ActivateUser()

	// 4. Actualizar en BD
	if err := rs.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("error activando usuario: %w", err)
	}

	// 5. Enviar email de bienvenida
	if err := rs.emailService.SendWelcomeEmail(ctx, user); err != nil {
		// Log error pero no fallar la verificación
		fmt.Printf("Error enviando email de bienvenida: %v\n", err)
	}

	return nil
}

// ResendEmailVerification reenvía el email de verificación
func (rs *RegisterServiceImpl) ResendEmailVerification(ctx context.Context, email string) error {
	// 1. Buscar usuario por email
	user, err := rs.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("usuario no encontrado")
	}

	// 2. Verificar que el usuario no esté ya verificado
	if user.EmailVerified {
		return errors.New("el email ya está verificado")
	}

	// 3. Generar nuevo token
	emailToken, err := rs.tokenGenerator.GenerateEmailVerificationToken()
	if err != nil {
		return fmt.Errorf("error generando token: %w", err)
	}

	// 4. Actualizar token en usuario
	user.SetEmailVerificationToken(emailToken, rs.emailTokenExpiry)

	// 5. Actualizar en BD
	if err := rs.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("error actualizando usuario: %w", err)
	}

	// 6. Enviar email
	if err := rs.emailService.SendEmailVerification(ctx, user, emailToken); err != nil {
		return fmt.Errorf("error enviando email: %w", err)
	}

	return nil
}

// ValidateRegistrationData valida los datos de registro
func (rs *RegisterServiceImpl) ValidateRegistrationData(req *ports.RegisterRequest) error {
	// Validaciones básicas
	if req.Username == "" || len(req.Username) < 3 || len(req.Username) > 50 {
		return errors.New("username debe tener entre 3 y 50 caracteres")
	}

	if req.FullName == "" || len(req.FullName) < 2 || len(req.FullName) > 200 {
		return errors.New("nombre completo debe tener entre 2 y 200 caracteres")
	}

	if !utils.ValidateChileanEmail(req.Email) {
		return errors.New("email inválido")
	}

	if len(req.Password) < 8 {
		return errors.New("contraseña debe tener al menos 8 caracteres")
	}

	// Validar teléfono si se proporciona
	if req.Phone != "" && !utils.ValidateChileanPhone(req.Phone) {
		return errors.New("formato de teléfono inválido")
	}

	// Validar RUT si se proporciona
	if req.IdentificationNumber != "" && !utils.ValidateRUT(req.IdentificationNumber) {
		return errors.New("RUT inválido")
	}

	// Validaciones específicas para tenant
	if req.CreateTenant && req.TenantData != nil {
		if err := rs.validateTenantData(req.TenantData); err != nil {
			return fmt.Errorf("datos del tenant inválidos: %w", err)
		}
	}

	return nil
}

// validateTenantData valida los datos específicos del tenant
func (rs *RegisterServiceImpl) validateTenantData(data *ports.TenantRegistrationData) error {
	if !utils.ValidateRUT(data.RUT) {
		return errors.New("RUT del negocio inválido")
	}

	if data.BusinessName == "" || len(data.BusinessName) < 2 || len(data.BusinessName) > 150 {
		return errors.New("nombre del negocio debe tener entre 2 y 150 caracteres")
	}

	if !utils.ValidateChileanEmail(data.Email) {
		return errors.New("email del negocio inválido")
	}

	if !utils.ValidateChileanPhone(data.Phone) {
		return errors.New("teléfono del negocio inválido")
	}

	if data.Address == "" || len(data.Address) > 200 {
		return errors.New("dirección del negocio requerida (máximo 200 caracteres)")
	}

	if data.Website == "" || len(data.Website) > 150 {
		return errors.New("sitio web del negocio requerido (máximo 150 caracteres)")
	}

	if len(data.RegionID) != 2 {
		return errors.New("ID de región debe tener 2 caracteres")
	}

	if data.CommuneID == "" || len(data.CommuneID) > 100 {
		return errors.New("ID de comuna requerido (máximo 100 caracteres)")
	}

	if data.CountryID == uuid.Nil {
		return errors.New("ID de país requerido")
	}

	return nil
}

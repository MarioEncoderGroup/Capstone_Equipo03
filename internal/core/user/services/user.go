package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/user/ports"
	"github.com/JoseLuis21/mv-backend/internal/shared/hasher"
)

// userService implementa el servicio de usuarios
type userService struct {
	userRepo    ports.UserRepository
	hasherSvc   *hasher.Service
}

// NewUserService crea una nueva instancia del servicio de usuario
func NewUserService(userRepo ports.UserRepository) ports.UserService {
	return &userService{
		userRepo:  userRepo,
		hasherSvc: hasher.NewService(),
	}
}

// CreateUser crea un nuevo usuario
func (s *userService) CreateUser(ctx context.Context, user *domain.User) error {
	return s.userRepo.Create(ctx, user)
}

// GetUserByID obtiene un usuario por su ID
func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetUserByEmail obtiene un usuario por su email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// GetUserByEmailToken obtiene un usuario por su token de verificación de email
func (s *userService) GetUserByEmailToken(ctx context.Context, token string) (*domain.User, error) {
	return s.userRepo.GetByEmailToken(ctx, token)
}

// GetUserByPasswordResetToken obtiene un usuario por su token de reset de contraseña
func (s *userService) GetUserByPasswordResetToken(ctx context.Context, token string) (*domain.User, error) {
	return s.userRepo.GetByPasswordResetToken(ctx, token)
}

// UpdateUser actualiza un usuario existente
func (s *userService) UpdateUser(ctx context.Context, user *domain.User) error {
	return s.userRepo.Update(ctx, user)
}

// ExistsByEmail verifica si existe un usuario con el email dado
func (s *userService) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return s.userRepo.ExistsByEmail(ctx, email)
}

// GetTenantsByUser obtiene todos los tenant_users por userID
func (s *userService) GetTenantsByUser(ctx context.Context, userID uuid.UUID) ([]*domain.TenantUser, error) {
	return s.userRepo.GetTenantsByUser(ctx, userID)
}

// UserHasAccessToTenant verifica si un usuario tiene acceso a un tenant
func (s *userService) UserHasAccessToTenant(ctx context.Context, userID, tenantID uuid.UUID) (bool, error) {
	return s.userRepo.UserHasAccessToTenant(ctx, userID, tenantID)
}

// AddUserToTenant asocia un usuario a un tenant
func (s *userService) AddUserToTenant(ctx context.Context, tenantUser *domain.TenantUser) error {
	return s.userRepo.AddUserToTenant(ctx, tenantUser)
}

// GetUsers obtiene una lista paginada de usuarios con validaciones de negocio
func (s *userService) GetUsers(ctx context.Context, offset, limit int, sortBy, sortDir, search string) ([]*domain.User, int64, error) {
	return s.userRepo.GetUsers(ctx, offset, limit, sortBy, sortDir, search)
}

// CheckUserExists verifica si un usuario existe por ID
func (s *userService) CheckUserExists(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.userRepo.CheckUserExists(ctx, id)
}

// DeleteUser elimina lógicamente un usuario (soft delete)
func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

// CreateUserFromDto crea un usuario desde un DTO con validaciones
func (s *userService) CreateUserFromDto(ctx context.Context, dto *domain.CreateUserDto) (*domain.User, error) {
	// Verificar que el email no existe
	exists, err := s.userRepo.ExistsByEmail(ctx, dto.Email)
	if err != nil {
		return nil, fmt.Errorf("error verificando email: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("el email %s ya está registrado", dto.Email)
	}

	// Hash de la contraseña
	hashedPassword, err := s.hasherSvc.Hash(dto.Password)
	if err != nil {
		return nil, fmt.Errorf("error hasheando contraseña: %w", err)
	}

	// Crear entidad User
	phone := ""
	if dto.Phone != nil {
		phone = *dto.Phone
	}
	
	// Determinar isActive (si no se envía, default true para usuarios creados por admin)
	isActive := true
	if dto.IsActive != nil {
		isActive = *dto.IsActive
	}
	
	user := domain.NewUser(dto.FullName, dto.Email, phone, hashedPassword, isActive)

	// Verificar y asegurar username único
	// Si el username generado ya existe, agregar sufijo incremental
	baseUsername := user.Username
	counter := 1
	for {
		// Verificar si el username está disponible
		usernameExists, err := s.userRepo.ExistsByUsername(ctx, user.Username)
		if err != nil {
			return nil, fmt.Errorf("error verificando username: %w", err)
		}

		// Si no existe, usar este username
		if !usernameExists {
			break
		}

		// Si existe, agregar sufijo numérico y reintentar
		user.Username = fmt.Sprintf("%s_%d", baseUsername, counter)
		counter++

		// Prevenir loop infinito (máximo 100 intentos)
		if counter > 100 {
			return nil, fmt.Errorf("no se pudo generar un username único después de 100 intentos")
		}
	}

	// Asignar campos opcionales
	if dto.IdentificationNumber != nil {
		user.IdentificationNumber = dto.IdentificationNumber
	}

	// Crear usuario en la base de datos
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("error creando usuario: %w", err)
	}

	return user, nil
}

// UpdateUserFromDto actualiza un usuario desde un DTO con validaciones
func (s *userService) UpdateUserFromDto(ctx context.Context, id uuid.UUID, dto *domain.UpdateUserDto) (*domain.User, error) {
	// Obtener usuario existente
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo usuario: %w", err)
	}

	// Actualizar campos si se proporcionan
	if dto.FullName != nil {
		user.FullName = *dto.FullName
	}
	if dto.Phone != nil {
		user.Phone = dto.Phone
	}
	if dto.IdentificationNumber != nil {
		user.IdentificationNumber = dto.IdentificationNumber
	}
	if dto.BankID != nil {
		bankID, err := uuid.Parse(*dto.BankID)
		if err != nil {
			return nil, fmt.Errorf("ID de banco inválido: %w", err)
		}
		user.BankID = &bankID
	}
	if dto.BankAccountNumber != nil {
		user.BankAccountNumber = dto.BankAccountNumber
	}
	if dto.BankAccountType != nil {
		user.BankAccountType = dto.BankAccountType
	}
	if dto.ImageURL != nil {
		user.ImageURL = dto.ImageURL
	}
	if dto.IsActive != nil {
		user.IsActive = *dto.IsActive
	}

	user.Updated = time.Now()

	// Actualizar en la base de datos
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("error actualizando usuario: %w", err)
	}

	return user, nil
}

// ChangeUserPassword cambia la contraseña de un usuario con validaciones
func (s *userService) ChangeUserPassword(ctx context.Context, id uuid.UUID, dto *domain.ChangePasswordDto) error {
	// Obtener usuario existente con contraseña
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error obteniendo usuario: %w", err)
	}

	// Verificar contraseña actual
	if err := s.hasherSvc.Verify(user.Password, dto.CurrentPassword); err != nil {
		return fmt.Errorf("contraseña actual incorrecta")
	}

	// Hash de la nueva contraseña
	hashedPassword, err := s.hasherSvc.Hash(dto.NewPassword)
	if err != nil {
		return fmt.Errorf("error hasheando nueva contraseña: %w", err)
	}

	// Cambiar contraseña usando método del dominio
	user.ChangePassword(hashedPassword)

	// Actualizar en la base de datos
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("error actualizando contraseña: %w", err)
	}

	return nil
}

// UpdateUserProfile actualiza el perfil de un usuario autenticado
func (s *userService) UpdateUserProfile(ctx context.Context, id uuid.UUID, dto *domain.UpdateProfileDto) (*domain.User, error) {
	// Obtener usuario existente
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo usuario: %w", err)
	}

	// Actualizar campos de perfil si se proporcionan
	if dto.FullName != nil {
		user.FullName = *dto.FullName
	}
	if dto.Phone != nil {
		user.Phone = dto.Phone
	}
	if dto.IdentificationNumber != nil {
		user.IdentificationNumber = dto.IdentificationNumber
	}
	if dto.BankID != nil {
		bankID, err := uuid.Parse(*dto.BankID)
		if err != nil {
			return nil, fmt.Errorf("ID de banco inválido: %w", err)
		}
		user.BankID = &bankID
	}
	if dto.BankAccountNumber != nil {
		user.BankAccountNumber = dto.BankAccountNumber
	}
	if dto.BankAccountType != nil {
		user.BankAccountType = dto.BankAccountType
	}
	if dto.ImageURL != nil {
		user.ImageURL = dto.ImageURL
	}

	user.Updated = time.Now()

	// Actualizar en la base de datos
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("error actualizando perfil: %w", err)
	}

	return user, nil
}

// SaveRefreshToken guarda un refresh token para un usuario
func (s *userService) SaveRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string, expiresIn int64) error {
	// TODO: Implementar guardado de refresh token en tabla user_refresh_tokens
	// Por ahora, solo retornamos nil para que compile
	// En producción, esto debe guardar el token en la base de datos
	return nil
}

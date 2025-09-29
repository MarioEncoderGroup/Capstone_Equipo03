package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/user/ports"
)

// userService implementa el servicio de usuarios
type userService struct {
	userRepo ports.UserRepository
}

// NewUserService crea una nueva instancia del servicio de usuario
func NewUserService(userRepo ports.UserRepository) ports.UserService {
	return &userService{
		userRepo: userRepo,
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
package mocks

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/user/domain"
)

// MockUserRepository implements UserRepository interface for testing
type MockUserRepository struct {
	users       map[uuid.UUID]*domain.User
	emailIndex  map[string]uuid.UUID
	tokenIndex  map[string]uuid.UUID
	tenantUsers map[uuid.UUID][]*domain.TenantUser
	
	// Control flags for testing error scenarios
	CreateError          bool
	GetByIDError         bool
	GetByEmailError      bool
	GetByEmailTokenError bool
	UpdateError          bool
	DeleteError          bool
	ExistsByEmailError   bool
}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:       make(map[uuid.UUID]*domain.User),
		emailIndex:  make(map[string]uuid.UUID),
		tokenIndex:  make(map[string]uuid.UUID),
		tenantUsers: make(map[uuid.UUID][]*domain.TenantUser),
	}
}

// Create implements UserRepository.Create
func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.CreateError {
		return errors.New("mock create error")
	}
	
	// Check if email already exists
	if _, exists := m.emailIndex[user.Email]; exists {
		return errors.New("user already exists")
	}
	
	// Store user
	m.users[user.ID] = user
	m.emailIndex[user.Email] = user.ID
	
	// Index email token if exists
	if user.EmailToken != nil {
		m.tokenIndex[*user.EmailToken] = user.ID
	}
	
	return nil
}

// GetByID implements UserRepository.GetByID
func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if m.GetByIDError {
		return nil, errors.New("mock get by id error")
	}
	
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	
	return user, nil
}

// GetByEmail implements UserRepository.GetByEmail
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.GetByEmailError {
		return nil, errors.New("mock get by email error")
	}
	
	userID, exists := m.emailIndex[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	
	return m.users[userID], nil
}

// GetByEmailToken implements UserRepository.GetByEmailToken
func (m *MockUserRepository) GetByEmailToken(ctx context.Context, token string) (*domain.User, error) {
	if m.GetByEmailTokenError {
		return nil, errors.New("mock get by email token error")
	}
	
	userID, exists := m.tokenIndex[token]
	if !exists {
		return nil, errors.New("user not found")
	}
	
	return m.users[userID], nil
}

// Update implements UserRepository.Update
func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateError {
		return errors.New("mock update error")
	}
	
	_, exists := m.users[user.ID]
	if !exists {
		return errors.New("user not found")
	}
	
	// Update user
	m.users[user.ID] = user
	
	// Update email index
	m.emailIndex[user.Email] = user.ID
	
	// Update token index
	if user.EmailToken != nil {
		m.tokenIndex[*user.EmailToken] = user.ID
	}
	
	return nil
}

// Delete implements UserRepository.Delete
func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError {
		return errors.New("mock delete error")
	}
	
	user, exists := m.users[id]
	if !exists {
		return errors.New("user not found")
	}
	
	// Remove from all indexes
	delete(m.users, id)
	delete(m.emailIndex, user.Email)
	
	if user.EmailToken != nil {
		delete(m.tokenIndex, *user.EmailToken)
	}
	
	return nil
}

// ExistsByEmail implements UserRepository.ExistsByEmail
func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.ExistsByEmailError {
		return false, errors.New("mock exists by email error")
	}
	
	_, exists := m.emailIndex[email]
	return exists, nil
}

// AddUserToTenant implements UserRepository.AddUserToTenant
func (m *MockUserRepository) AddUserToTenant(ctx context.Context, tenantUser *domain.TenantUser) error {
	m.tenantUsers[tenantUser.UserID] = append(m.tenantUsers[tenantUser.UserID], tenantUser)
	return nil
}

// RemoveUserFromTenant implements UserRepository.RemoveUserFromTenant
func (m *MockUserRepository) RemoveUserFromTenant(ctx context.Context, userID, tenantID uuid.UUID) error {
	userTenants := m.tenantUsers[userID]
	for i, tu := range userTenants {
		if tu.TenantID == tenantID {
			m.tenantUsers[userID] = append(userTenants[:i], userTenants[i+1:]...)
			break
		}
	}
	return nil
}

// GetUserTenants implements UserRepository.GetUserTenants
func (m *MockUserRepository) GetUserTenants(ctx context.Context, userID uuid.UUID) ([]*domain.TenantUser, error) {
	return m.tenantUsers[userID], nil
}

// Reset clears all data and error flags
func (m *MockUserRepository) Reset() {
	m.users = make(map[uuid.UUID]*domain.User)
	m.emailIndex = make(map[string]uuid.UUID)
	m.tokenIndex = make(map[string]uuid.UUID)
	m.tenantUsers = make(map[uuid.UUID][]*domain.TenantUser)
	
	m.CreateError = false
	m.GetByIDError = false
	m.GetByEmailError = false
	m.GetByEmailTokenError = false
	m.UpdateError = false
	m.DeleteError = false
	m.ExistsByEmailError = false
}

// SetUser manually sets a user for testing
func (m *MockUserRepository) SetUser(user *domain.User) {
	m.users[user.ID] = user
	m.emailIndex[user.Email] = user.ID
	if user.EmailToken != nil {
		m.tokenIndex[*user.EmailToken] = user.ID
	}
}

// GetUserCount returns the number of users stored
func (m *MockUserRepository) GetUserCount() int {
	return len(m.users)
}
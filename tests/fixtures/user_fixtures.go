package fixtures

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/user/domain"
)

// UserFixtures provides predefined user data for testing
type UserFixtures struct{}

// NewUserFixtures creates a new UserFixtures instance
func NewUserFixtures() *UserFixtures {
	return &UserFixtures{}
}

// ValidUser returns a valid user for testing
func (f *UserFixtures) ValidUser() *domain.User {
	return domain.NewUser(
		"testuser",
		"Usuario de Prueba",
		"test@example.cl",
		"$2a$12$hashed.password.example",
	)
}

// UserWithEmailToken returns a user with email verification token
func (f *UserFixtures) UserWithEmailToken() *domain.User {
	user := f.ValidUser()
	token := "email_token_123456"
	expiry := time.Now().Add(24 * time.Hour)
	user.EmailToken = &token
	user.EmailTokenExpires = &expiry
	return user
}

// VerifiedUser returns a user with verified email
func (f *UserFixtures) VerifiedUser() *domain.User {
	user := f.ValidUser()
	user.EmailVerified = true
	user.IsActive = true
	return user
}

// ActiveUser returns an active user ready for login
func (f *UserFixtures) ActiveUser() *domain.User {
	user := f.VerifiedUser()
	user.LastLogin = &time.Time{}
	now := time.Now()
	user.LastLogin = &now
	return user
}

// UserWithResetToken returns a user with password reset token
func (f *UserFixtures) UserWithResetToken() *domain.User {
	user := f.ValidUser()
	token := "reset_token_123456"
	expiry := time.Now().Add(1 * time.Hour)
	user.PasswordResetToken = &token
	user.PasswordResetExpires = &expiry
	return user
}

// UserWithBankInfo returns a user with bank account information
func (f *UserFixtures) UserWithBankInfo() *domain.User {
	user := f.ValidUser()
	bankID := uuid.New()
	user.BankID = &bankID
	bankAccount := "12345678901"
	user.BankAccountNumber = &bankAccount
	bankType := "CUENTA_CORRIENTE"
	user.BankAccountType = &bankType
	return user
}

// ExpiredEmailToken returns a user with expired email token
func (f *UserFixtures) ExpiredEmailToken() *domain.User {
	user := f.ValidUser()
	token := "expired_token_123456"
	expiry := time.Now().Add(-1 * time.Hour) // Expired 1 hour ago
	user.EmailToken = &token
	user.EmailTokenExpires = &expiry
	return user
}

// UserWithValidRUT returns a user with valid Chilean RUT
func (f *UserFixtures) UserWithValidRUT() *domain.User {
	user := f.ValidUser()
	rut := "12.345.678-5"
	user.IdentificationNumber = &rut
	return user
}

// UserWithInvalidRUT returns a user with invalid Chilean RUT
func (f *UserFixtures) UserWithInvalidRUT() *domain.User {
	user := f.ValidUser()
	invalidRUT := "12.345.678-X"
	user.IdentificationNumber = &invalidRUT
	return user
}

// UsersForBulkTesting returns multiple users for bulk operations testing
func (f *UserFixtures) UsersForBulkTesting(count int) []*domain.User {
	users := make([]*domain.User, count)
	for i := 0; i < count; i++ {
		users[i] = domain.NewUser(
			fmt.Sprintf("testuser%d", i),
			fmt.Sprintf("Usuario de Prueba %d", i),
			fmt.Sprintf("test%d@example.cl", i),
			"$2a$12$hashed.password.example",
		)
	}
	return users
}

// ValidRegisterRequestData returns request data for registration testing
func (f *UserFixtures) ValidRegisterRequestData() map[string]interface{} {
	return map[string]interface{}{
		"username":               "testuser",
		"full_name":              "Usuario de Prueba",
		"email":                  "test@example.cl",
		"password":               "password123",
		"phone":                  "+56912345678",
		"identification_number":  "12.345.678-5",
	}
}

// InvalidRegisterRequestData returns invalid request data for testing
func (f *UserFixtures) InvalidRegisterRequestData() map[string]interface{} {
	return map[string]interface{}{
		"username":               "tu", // Too short
		"full_name":              "",   // Empty
		"email":                  "invalid-email",
		"password":               "123", // Too weak
		"phone":                  "invalid-phone",
		"identification_number":  "12.345.678-X", // Invalid RUT
	}
}
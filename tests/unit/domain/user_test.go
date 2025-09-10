package domain_test

import (
	"testing"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/core/user/domain"
)

func TestNewUser(t *testing.T) {
	username := "testuser"
	fullName := "Test User"
	email := "test@example.cl"
	hashedPassword := "hashed_password_here"

	user := domain.NewUser(username, fullName, email, hashedPassword)

	// Verificar campos básicos
	if user.Username != username {
		t.Errorf("Expected username %s, got %s", username, user.Username)
	}
	
	if user.FullName != fullName {
		t.Errorf("Expected full name %s, got %s", fullName, user.FullName)
	}
	
	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}
	
	if user.Password != hashedPassword {
		t.Errorf("Expected password %s, got %s", hashedPassword, user.Password)
	}

	// Verificar estado inicial
	if user.EmailVerified {
		t.Error("New user should not have email verified")
	}
	
	if user.IsActive {
		t.Error("New user should not be active")
	}
	
	// Verificar UUID generado
	if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
		t.Error("User ID should be generated")
	}
	
	// Verificar timestamps
	now := time.Now()
	if user.Created.After(now) || user.Created.Before(now.Add(-time.Second)) {
		t.Error("Created timestamp should be around current time")
	}
	
	if user.Updated.After(now) || user.Updated.Before(now.Add(-time.Second)) {
		t.Error("Updated timestamp should be around current time")
	}
}

func TestUserEmailTokenValidation(t *testing.T) {
	user := domain.NewUser("test", "Test User", "test@example.cl", "password")
	
	token := "test_token_123"
	duration := 1 * time.Hour
	
	// Establecer token
	user.SetEmailVerificationToken(token, duration)
	
	// Token válido
	if !user.IsEmailTokenValid(token) {
		t.Error("Token should be valid")
	}
	
	// Token inválido
	if user.IsEmailTokenValid("wrong_token") {
		t.Error("Wrong token should be invalid")
	}
	
	// Token expirado
	user.SetEmailVerificationToken(token, -1*time.Hour) // Token expirado
	if user.IsEmailTokenValid(token) {
		t.Error("Expired token should be invalid")
	}
	
	// Sin token
	user.EmailToken = nil
	if user.IsEmailTokenValid(token) {
		t.Error("User without token should be invalid")
	}
}

func TestUserActivation(t *testing.T) {
	user := domain.NewUser("test", "Test User", "test@example.cl", "password")
	token := "test_token"
	
	// Establecer token de verificación
	user.SetEmailVerificationToken(token, 1*time.Hour)
	
	// Verificar estado inicial
	if user.EmailVerified || user.IsActive {
		t.Error("User should not be verified or active initially")
	}
	
	// Activar usuario
	beforeActivation := time.Now()
	user.ActivateUser()
	afterActivation := time.Now()
	
	// Verificar activación
	if !user.EmailVerified {
		t.Error("User should be email verified after activation")
	}
	
	if !user.IsActive {
		t.Error("User should be active after activation")
	}
	
	// Verificar que el token se eliminó
	if user.EmailToken != nil {
		t.Error("Email token should be cleared after activation")
	}
	
	if user.EmailTokenExpires != nil {
		t.Error("Email token expiration should be cleared after activation")
	}
	
	// Verificar que updated timestamp cambió
	if user.Updated.Before(beforeActivation) || user.Updated.After(afterActivation) {
		t.Error("Updated timestamp should be set during activation")
	}
}

func TestSetEmailVerificationToken(t *testing.T) {
	user := domain.NewUser("test", "Test User", "test@example.cl", "password")
	
	token := "verification_token_123"
	duration := 2 * time.Hour
	
	beforeSet := time.Now()
	user.SetEmailVerificationToken(token, duration)
	afterSet := time.Now()
	
	// Verificar token establecido
	if user.EmailToken == nil || *user.EmailToken != token {
		t.Errorf("Expected token %s, got %v", token, user.EmailToken)
	}
	
	// Verificar expiración
	if user.EmailTokenExpires == nil {
		t.Fatal("Email token expiration should be set")
	}
	
	expectedExpiration := beforeSet.Add(duration)
	actualExpiration := *user.EmailTokenExpires
	
	if actualExpiration.Before(expectedExpiration.Add(-time.Second)) || 
	   actualExpiration.After(expectedExpiration.Add(time.Second)) {
		t.Errorf("Token expiration not set correctly. Expected around %v, got %v", 
			expectedExpiration, actualExpiration)
	}
	
	// Verificar updated timestamp
	if user.Updated.Before(beforeSet) || user.Updated.After(afterSet) {
		t.Error("Updated timestamp should be set when setting email token")
	}
}
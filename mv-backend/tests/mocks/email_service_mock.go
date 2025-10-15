package mocks

import (
	"context"
	"errors"

	"github.com/JoseLuis21/mv-backend/internal/core/user/domain"
)

// MockEmailService implements EmailService interface for testing
type MockEmailService struct {
	// Tracking sent emails
	SentEmails []SentEmail

	// Control flags for testing error scenarios
	SendVerificationError  bool
	SendWelcomeError       bool
	SendPasswordResetError bool
}

// SentEmail represents an email that was sent during testing
type SentEmail struct {
	To       string
	Subject  string
	Type     string
	Template string
	Data     map[string]interface{}
}

// NewMockEmailService creates a new mock email service
func NewMockEmailService() *MockEmailService {
	return &MockEmailService{
		SentEmails: make([]SentEmail, 0),
	}
}

// SendVerificationEmail implements EmailService.SendVerificationEmail
func (m *MockEmailService) SendVerificationEmail(ctx context.Context, email, token string) error {
	if m.SendVerificationError {
		return errors.New("mock verification email error")
	}

	m.SentEmails = append(m.SentEmails, SentEmail{
		To:       email,
		Subject:  "Verificación de Email",
		Type:     "verification",
		Template: "email_verification",
		Data: map[string]interface{}{
			"token": token,
			"email": email,
		},
	})

	return nil
}

// SendWelcomeEmail implements EmailService.SendWelcomeEmail
func (m *MockEmailService) SendWelcomeEmail(ctx context.Context, user *domain.User) error {
	if m.SendWelcomeError {
		return errors.New("mock welcome email error")
	}

	m.SentEmails = append(m.SentEmails, SentEmail{
		To:       user.Email,
		Subject:  "Bienvenido a MisViáticos",
		Type:     "welcome",
		Template: "welcome",
		Data: map[string]interface{}{
			"user_id":   user.ID.String(),
			"full_name": user.FullName,
			"email":     user.Email,
		},
	})

	return nil
}

// SendPasswordResetEmail implements EmailService.SendPasswordResetEmail
func (m *MockEmailService) SendPasswordResetEmail(ctx context.Context, email, token string) error {
	if m.SendPasswordResetError {
		return errors.New("mock password reset email error")
	}

	m.SentEmails = append(m.SentEmails, SentEmail{
		To:       email,
		Subject:  "Recuperación de Contraseña",
		Type:     "password_reset",
		Template: "password_reset",
		Data: map[string]interface{}{
			"token": token,
			"email": email,
		},
	})

	return nil
}

// Reset clears all sent emails and error flags
func (m *MockEmailService) Reset() {
	m.SentEmails = make([]SentEmail, 0)
	m.SendVerificationError = false
	m.SendWelcomeError = false
	m.SendPasswordResetError = false
}

// GetSentEmailsCount returns the number of emails sent
func (m *MockEmailService) GetSentEmailsCount() int {
	return len(m.SentEmails)
}

// GetSentEmailsByType returns emails of a specific type
func (m *MockEmailService) GetSentEmailsByType(emailType string) []SentEmail {
	var emails []SentEmail
	for _, email := range m.SentEmails {
		if email.Type == emailType {
			emails = append(emails, email)
		}
	}
	return emails
}

// GetLastSentEmail returns the last email sent
func (m *MockEmailService) GetLastSentEmail() *SentEmail {
	if len(m.SentEmails) == 0 {
		return nil
	}
	return &m.SentEmails[len(m.SentEmails)-1]
}

// WasEmailSentTo checks if an email was sent to a specific recipient
func (m *MockEmailService) WasEmailSentTo(email string) bool {
	for _, sentEmail := range m.SentEmails {
		if sentEmail.To == email {
			return true
		}
	}
	return false
}

// GetEmailsSentTo returns all emails sent to a specific recipient
func (m *MockEmailService) GetEmailsSentTo(email string) []SentEmail {
	var emails []SentEmail
	for _, sentEmail := range m.SentEmails {
		if sentEmail.To == email {
			emails = append(emails, sentEmail)
		}
	}
	return emails
}

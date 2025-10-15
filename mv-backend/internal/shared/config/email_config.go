package config

import (
	"github.com/JoseLuis21/mv-backend/internal/shared/utils"
	"strconv"
)

// EmailConfig contiene la configuración para el servicio de email
type EmailConfig struct {
	// General
	Enabled     bool   `json:"enabled"`
	FromEmail   string `json:"from_email"`
	FromName    string `json:"from_name"`
	Provider    string `json:"provider"` // "smtp", "ses", "sendgrid"
	
	// URLs del frontend
	FrontendURL string `json:"frontend_url"`
	
	// SMTP Configuration
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
	SMTPUseTLS   bool   `json:"smtp_use_tls"`
	
	// AWS SES Configuration
	AWSRegion    string `json:"aws_region"`
	AWSAccessKey string `json:"aws_access_key"`
	AWSSecretKey string `json:"aws_secret_key"`
	
	// SendGrid Configuration
	SendGridAPIKey string `json:"sendgrid_api_key"`
	
	// Rate limiting
	MaxEmailsPerHour int `json:"max_emails_per_hour"`
	
	// Templates
	TemplateDir string `json:"template_dir"`
}

// NewEmailConfig crea configuración de email desde variables de entorno
func NewEmailConfig() *EmailConfig {
	// Parse SMTP port
	smtpPort := 587
	if portStr := utils.GetEnvOrDefault("EMAIL_SMTP_PORT", "587"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			smtpPort = p
		}
	}
	
	// Parse max emails per hour
	maxEmails := 100
	if maxStr := utils.GetEnvOrDefault("EMAIL_MAX_PER_HOUR", "100"); maxStr != "" {
		if m, err := strconv.Atoi(maxStr); err == nil {
			maxEmails = m
		}
	}
	
	// Parse TLS setting
	useTLS := utils.GetEnvOrDefault("EMAIL_SMTP_TLS", "true") == "true"
	
	return &EmailConfig{
		// General
		Enabled:     utils.GetEnvOrDefault("EMAIL_ENABLED", "false") == "true",
		FromEmail:   utils.GetEnvOrDefault("EMAIL_FROM", "noreply@misviaticos.cl"),
		FromName:    utils.GetEnvOrDefault("EMAIL_FROM_NAME", "MisViáticos"),
		Provider:    utils.GetEnvOrDefault("EMAIL_PROVIDER", "smtp"),
		FrontendURL: utils.GetEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
		
		// SMTP
		SMTPHost:     utils.GetEnvOrDefault("EMAIL_SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     smtpPort,
		SMTPUsername: utils.GetEnvOrDefault("EMAIL_SMTP_USERNAME", ""),
		SMTPPassword: utils.GetEnvOrDefault("EMAIL_SMTP_PASSWORD", ""),
		SMTPUseTLS:   useTLS,
		
		// AWS SES
		AWSRegion:    utils.GetEnvOrDefault("AWS_REGION", "us-east-1"),
		AWSAccessKey: utils.GetEnvOrDefault("AWS_ACCESS_KEY_ID", ""),
		AWSSecretKey: utils.GetEnvOrDefault("AWS_SECRET_ACCESS_KEY", ""),
		
		// SendGrid
		SendGridAPIKey: utils.GetEnvOrDefault("SENDGRID_API_KEY", ""),
		
		// Rate limiting
		MaxEmailsPerHour: maxEmails,
		
		// Templates
		TemplateDir: utils.GetEnvOrDefault("EMAIL_TEMPLATE_DIR", "./templates/email"),
	}
}

// IsValidProvider verifica si el proveedor de email configurado es válido
func (c *EmailConfig) IsValidProvider() bool {
	validProviders := map[string]bool{
		"smtp":     true,
		"ses":      true,
		"sendgrid": true,
	}
	return validProviders[c.Provider]
}

// HasSMTPConfig verifica si la configuración SMTP está completa
func (c *EmailConfig) HasSMTPConfig() bool {
	return c.SMTPHost != "" && c.SMTPPort > 0 && c.SMTPUsername != "" && c.SMTPPassword != ""
}

// HasAWSConfig verifica si la configuración AWS SES está completa
func (c *EmailConfig) HasAWSConfig() bool {
	return c.AWSRegion != "" && c.AWSAccessKey != "" && c.AWSSecretKey != ""
}

// HasSendGridConfig verifica si la configuración SendGrid está completa
func (c *EmailConfig) HasSendGridConfig() bool {
	return c.SendGridAPIKey != ""
}

// CanSendEmails verifica si se puede enviar emails con la configuración actual
func (c *EmailConfig) CanSendEmails() bool {
	if !c.Enabled {
		return false
	}
	
	switch c.Provider {
	case "smtp":
		return c.HasSMTPConfig()
	case "ses":
		return c.HasAWSConfig()
	case "sendgrid":
		return c.HasSendGridConfig()
	default:
		return false
	}
}
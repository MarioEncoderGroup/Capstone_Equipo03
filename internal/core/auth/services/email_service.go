package services

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/JoseLuis21/mv-backend/internal/core/auth/ports"
	"github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	"github.com/JoseLuis21/mv-backend/internal/shared/config"
	"github.com/JoseLuis21/mv-backend/internal/shared/errors"
)

// EmailServiceImpl implementa el servicio de envío de emails
// Soporta múltiples proveedores: SMTP, AWS SES, SendGrid
type EmailServiceImpl struct {
	config *config.EmailConfig
}

// NewEmailService crea una nueva instancia del servicio de email
func NewEmailService() ports.EmailService {
	emailConfig := config.NewEmailConfig()

	// Validar configuración al inicializar
	if emailConfig.Enabled && !emailConfig.CanSendEmails() {
		log.Printf("⚠️  Email habilitado pero configuración incompleta para proveedor: %s", emailConfig.Provider)
	}

	return &EmailServiceImpl{
		config: emailConfig,
	}
}

// EmailMessage representa un mensaje de email a enviar
type EmailMessage struct {
	To      []string `json:"to"`
	CC      []string `json:"cc,omitempty"`
	BCC     []string `json:"bcc,omitempty"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	IsHTML  bool     `json:"is_html"`
}

// sendEmail envía un email usando el proveedor configurado
func (es *EmailServiceImpl) sendEmail(ctx context.Context, email *EmailMessage) error {
	switch es.config.Provider {
	case "smtp":
		return es.sendViaSMTP(ctx, email)
	case "ses":
		return es.sendViaAWSSES(ctx, email)
	case "sendgrid":
		return es.sendViaSendGrid(ctx, email)
	default:
		return errors.NewInternalError("proveedor de email no soportado", es.config.Provider)
	}
}

// sendViaSMTP envía email via SMTP
func (es *EmailServiceImpl) sendViaSMTP(ctx context.Context, email *EmailMessage) error {
	if !es.config.HasSMTPConfig() {
		return errors.NewInternalError("configuración SMTP incompleta", "")
	}

	// Configurar autenticación
	auth := smtp.PlainAuth("", es.config.SMTPUsername, es.config.SMTPPassword, es.config.SMTPHost)

	// Construir mensaje
	message := es.buildMIMEMessage(email)

	// Servidor SMTP
	smtpServer := fmt.Sprintf("%s:%d", es.config.SMTPHost, es.config.SMTPPort)

	// Enviar email
	err := smtp.SendMail(smtpServer, auth, es.config.FromEmail, email.To, []byte(message))
	if err != nil {
		return errors.WrapError(errors.ErrEmailService, fmt.Sprintf("error SMTP: %v", err))
	}

	log.Printf("📧 Email enviado via SMTP a: %s", strings.Join(email.To, ", "))
	return nil
}

// sendViaAWSSES envía email via AWS SES
func (es *EmailServiceImpl) sendViaAWSSES(ctx context.Context, email *EmailMessage) error {
	// TODO: Implementar AWS SES
	// Por ahora simular el envío
	log.Printf("📧 [SES SIMULADO] Email enviado a: %s", strings.Join(email.To, ", "))
	return nil
}

// sendViaSendGrid envía email via SendGrid
func (es *EmailServiceImpl) sendViaSendGrid(ctx context.Context, email *EmailMessage) error {
	// TODO: Implementar SendGrid
	// Por ahora simular el envío
	log.Printf("📧 [SENDGRID SIMULADO] Email enviado a: %s", strings.Join(email.To, ", "))
	return nil
}

// buildMIMEMessage construye un mensaje MIME para SMTP
func (es *EmailServiceImpl) buildMIMEMessage(email *EmailMessage) string {
	var message strings.Builder

	// Headers
	message.WriteString(fmt.Sprintf("From: %s <%s>\r\n", es.config.FromName, es.config.FromEmail))
	message.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ", ")))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))

	if email.IsHTML {
		message.WriteString("MIME-Version: 1.0\r\n")
		message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	message.WriteString("\r\n")
	message.WriteString(email.Body)

	return message.String()
}

// SendEmailVerification envía un email de verificación al usuario
func (es *EmailServiceImpl) SendEmailVerification(ctx context.Context, user *domain.User, token string) error {
	if !es.config.Enabled {
		// Simular envío en desarrollo
		verificationURL := fmt.Sprintf("%s/verify-email?token=%s", es.config.FrontendURL, token)
		fmt.Printf("📧 [SIMULADO] Email de verificación enviado a: %s\n", user.Email)
		fmt.Printf("🔗 Link de verificación: %s\n", verificationURL)
		return nil
	}

	if !es.config.CanSendEmails() {
		return errors.WrapError(errors.ErrEmailService, "configuración de email incompleta")
	}

	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", es.config.FrontendURL, token)

	email := &EmailMessage{
		To:      []string{user.Email},
		Subject: "Verifica tu cuenta en MisViáticos",
		Body:    es.buildEmailVerificationHTML(user.FullName, verificationURL),
		IsHTML:  true,
	}

	return es.sendEmail(ctx, email)
}

// SendPasswordReset envía un email para reset de contraseña
func (es *EmailServiceImpl) SendPasswordReset(ctx context.Context, user *domain.User, token string) error {
	if !es.config.Enabled {
		resetURL := fmt.Sprintf("%s/reset-password?token=%s", es.config.FrontendURL, token)
		fmt.Printf("📧 [SIMULADO] Email de reset enviado a: %s\n", user.Email)
		fmt.Printf("🔗 Link de reset: %s\n", resetURL)
		return nil
	}

	if !es.config.CanSendEmails() {
		return errors.WrapError(errors.ErrEmailService, "configuración de email incompleta")
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", es.config.FrontendURL, token)

	email := &EmailMessage{
		To:      []string{user.Email},
		Subject: "Recupera tu contraseña - MisViáticos",
		Body:    es.buildPasswordResetHTML(user.FullName, resetURL),
		IsHTML:  true,
	}

	return es.sendEmail(ctx, email)
}

// SendWelcomeEmail envía email de bienvenida después del registro
func (es *EmailServiceImpl) SendWelcomeEmail(ctx context.Context, user *domain.User) error {
	if !es.config.Enabled {
		fmt.Printf("📧 [SIMULADO] Email de bienvenida enviado a: %s\n", user.Email)
		return nil
	}

	if !es.config.CanSendEmails() {
		return errors.WrapError(errors.ErrEmailService, "configuración de email incompleta")
	}

	email := &EmailMessage{
		To:      []string{user.Email},
		Subject: "¡Bienvenido a MisViáticos!",
		Body:    es.buildWelcomeEmailHTML(user.FullName),
		IsHTML:  true,
	}

	return es.sendEmail(ctx, email)
}

// buildEmailVerificationHTML construye el HTML para email de verificación
func (es *EmailServiceImpl) buildEmailVerificationHTML(fullName, verificationURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verifica tu cuenta</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #2c3e50;">¡Hola %s!</h2>
        
        <p>Gracias por registrarte en MisViáticos. Para completar tu registro, por favor verifica tu email haciendo clic en el siguiente enlace:</p>
        
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s" style="display: inline-block; padding: 12px 30px; background-color: #3498db; color: white; text-decoration: none; border-radius: 5px; font-weight: bold;">
                Verificar Email
            </a>
        </div>
        
        <p>Si no puedes hacer clic en el botón, copia y pega el siguiente enlace en tu navegador:</p>
        <p style="word-break: break-all; color: #666;">%s</p>
        
        <p style="margin-top: 30px; font-size: 12px; color: #666;">
            Este enlace expirará en 24 horas. Si no solicitaste este registro, puedes ignorar este email.
        </p>
        
        <hr style="margin: 30px 0; border: 0; border-top: 1px solid #eee;">
        <p style="font-size: 12px; color: #666; text-align: center;">
            © 2024 MisViáticos - Sistema de gestión de viáticos para Chile
        </p>
    </div>
</body>
</html>
`, fullName, verificationURL, verificationURL)
}

// buildPasswordResetHTML construye el HTML para email de reset
func (es *EmailServiceImpl) buildPasswordResetHTML(fullName, resetURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Recupera tu contraseña</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #e74c3c;">Recupera tu contraseña</h2>
        
        <p>Hola %s,</p>
        
        <p>Recibimos una solicitud para restablecer la contraseña de tu cuenta en MisViáticos. Haz clic en el siguiente enlace para crear una nueva contraseña:</p>
        
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s" style="display: inline-block; padding: 12px 30px; background-color: #e74c3c; color: white; text-decoration: none; border-radius: 5px; font-weight: bold;">
                Restablecer Contraseña
            </a>
        </div>
        
        <p>Si no puedes hacer clic en el botón, copia y pega el siguiente enlace en tu navegador:</p>
        <p style="word-break: break-all; color: #666;">%s</p>
        
        <p style="margin-top: 30px; font-size: 12px; color: #666;">
            Este enlace expirará en 1 hora. Si no solicitaste este cambio, puedes ignorar este email.
        </p>
        
        <hr style="margin: 30px 0; border: 0; border-top: 1px solid #eee;">
        <p style="font-size: 12px; color: #666; text-align: center;">
            © 2024 MisViáticos - Sistema de gestión de viáticos para Chile
        </p>
    </div>
</body>
</html>
`, fullName, resetURL, resetURL)
}

// buildWelcomeEmailHTML construye el HTML para email de bienvenida
func (es *EmailServiceImpl) buildWelcomeEmailHTML(fullName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>¡Bienvenido a MisViáticos!</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #27ae60;">¡Bienvenido a MisViáticos!</h2>
        
        <p>¡Hola %s!</p>
        
        <p>Tu cuenta ha sido verificada exitosamente. Ya puedes comenzar a usar MisViáticos para gestionar tus viáticos y gastos empresariales.</p>
        
        <div style="background-color: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0;">
            <h3 style="color: #2c3e50; margin-top: 0;">¿Qué puedes hacer ahora?</h3>
            <ul style="margin: 10px 0;">
                <li>Crear y gestionar reportes de gastos</li>
                <li>Cargar recibos y facturas</li>
                <li>Configurar políticas de viáticos</li>
                <li>Invitar a otros miembros de tu empresa</li>
            </ul>
        </div>
        
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s" style="display: inline-block; padding: 12px 30px; background-color: #27ae60; color: white; text-decoration: none; border-radius: 5px; font-weight: bold;">
                Ir a MisViáticos
            </a>
        </div>
        
        <p>Si tienes alguna pregunta o necesitas ayuda, no dudes en contactarnos.</p>
        
        <hr style="margin: 30px 0; border: 0; border-top: 1px solid #eee;">
        <p style="font-size: 12px; color: #666; text-align: center;">
            © 2024 MisViáticos - Sistema de gestión de viáticos para Chile
        </p>
    </div>
</body>
</html>
`, fullName, es.config.FrontendURL)
}

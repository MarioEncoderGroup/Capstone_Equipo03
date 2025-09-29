package email

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/JoseLuis21/mv-backend/internal/shared/config"
	"github.com/JoseLuis21/mv-backend/internal/shared/errors"
)

// Service interface genérica para envío de emails
// Puede ser utilizada por cualquier dominio en el sistema
type Service interface {
	SendEmail(ctx context.Context, message *Message) error
	SendTemplateEmail(ctx context.Context, template EmailTemplate, data interface{}) error
}

// ServiceImpl implementa el servicio genérico de envío de emails
// Soporta múltiples proveedores: SMTP, AWS SES, SendGrid
type ServiceImpl struct {
	config *config.EmailConfig
}

// NewService crea una nueva instancia del servicio genérico de email
func NewService() Service {
	emailConfig := config.NewEmailConfig()

	// Validar configuración al inicializar
	if emailConfig.Enabled && !emailConfig.CanSendEmails() {
		log.Printf("⚠️  Email habilitado pero configuración incompleta para proveedor: %s", emailConfig.Provider)
	}

	return &ServiceImpl{
		config: emailConfig,
	}
}

// Message representa un mensaje de email genérico
type Message struct {
	To      []string `json:"to"`
	CC      []string `json:"cc,omitempty"`
	BCC     []string `json:"bcc,omitempty"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	IsHTML  bool     `json:"is_html"`
}

// EmailTemplate define los tipos de templates disponibles
type EmailTemplate int

const (
	TemplateEmailVerification EmailTemplate = iota
	TemplatePasswordReset
	TemplateWelcome
	TemplateGeneric
)

// TemplateData estructura base para datos de templates
type TemplateData struct {
	FullName string
	Email    string
	URL      string
	Message  string
}

// SendEmail envía un email genérico usando el proveedor configurado
func (es *ServiceImpl) SendEmail(ctx context.Context, message *Message) error {
	switch es.config.Provider {
	case "smtp":
		return es.sendViaSMTP(ctx, message)
	case "ses":
		return es.sendViaAWSSES(ctx, message)
	case "sendgrid":
		return es.sendViaSendGrid(ctx, message)
	default:
		return errors.NewInternalError("proveedor de email no soportado", es.config.Provider)
	}
}

// SendTemplateEmail envía un email usando un template predefinido
func (es *ServiceImpl) SendTemplateEmail(ctx context.Context, template EmailTemplate, data interface{}) error {
	templateData, ok := data.(*TemplateData)
	if !ok {
		return errors.NewValidationError("datos de template inválidos", "template_data")
	}

	var message *Message
	var err error

	switch template {
	case TemplateEmailVerification:
		message, err = es.buildEmailVerificationMessage(templateData)
	case TemplatePasswordReset:
		message, err = es.buildPasswordResetMessage(templateData)
	case TemplateWelcome:
		message, err = es.buildWelcomeMessage(templateData)
	case TemplateGeneric:
		message, err = es.buildGenericMessage(templateData)
	default:
		return errors.NewValidationError("template de email no soportado", "template")
	}

	if err != nil {
		return err
	}

	// Si message es nil (email deshabilitado), no hacer nada
	if message == nil {
		return nil
	}

	return es.SendEmail(ctx, message)
}

// sendViaSMTP envía email via SMTP
func (es *ServiceImpl) sendViaSMTP(ctx context.Context, message *Message) error {
	if !es.config.HasSMTPConfig() {
		return errors.NewInternalError("configuración SMTP incompleta", "")
	}

	// Configurar autenticación
	auth := smtp.PlainAuth("", es.config.SMTPUsername, es.config.SMTPPassword, es.config.SMTPHost)

	// Construir mensaje MIME
	mimeMessage := es.buildMIMEMessage(message)

	// Servidor SMTP
	smtpServer := fmt.Sprintf("%s:%d", es.config.SMTPHost, es.config.SMTPPort)

	// Enviar email
	err := smtp.SendMail(smtpServer, auth, es.config.FromEmail, message.To, []byte(mimeMessage))
	if err != nil {
		return errors.WrapError(errors.ErrEmailService, fmt.Sprintf("error SMTP: %v", err))
	}

	log.Printf("📧 Email enviado via SMTP a: %s", strings.Join(message.To, ", "))
	return nil
}

// sendViaAWSSES envía email via AWS SES
func (es *ServiceImpl) sendViaAWSSES(ctx context.Context, message *Message) error {
	// TODO: Implementar AWS SES
	// Por ahora simular el envío
	if !es.config.Enabled {
		log.Printf("📧 [SES SIMULADO] Email enviado a: %s", strings.Join(message.To, ", "))
		return nil
	}
	
	log.Printf("📧 [SES SIMULADO] Email enviado a: %s", strings.Join(message.To, ", "))
	return nil
}

// sendViaSendGrid envía email via SendGrid
func (es *ServiceImpl) sendViaSendGrid(ctx context.Context, message *Message) error {
	// TODO: Implementar SendGrid
	// Por ahora simular el envío
	if !es.config.Enabled {
		log.Printf("📧 [SENDGRID SIMULADO] Email enviado a: %s", strings.Join(message.To, ", "))
		return nil
	}
	
	log.Printf("📧 [SENDGRID SIMULADO] Email enviado a: %s", strings.Join(message.To, ", "))
	return nil
}

// buildMIMEMessage construye un mensaje MIME para SMTP
func (es *ServiceImpl) buildMIMEMessage(message *Message) string {
	var mimeMessage strings.Builder

	// Headers
	mimeMessage.WriteString(fmt.Sprintf("From: %s <%s>\r\n", es.config.FromName, es.config.FromEmail))
	mimeMessage.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(message.To, ", ")))
	mimeMessage.WriteString(fmt.Sprintf("Subject: %s\r\n", message.Subject))

	if message.IsHTML {
		mimeMessage.WriteString("MIME-Version: 1.0\r\n")
		mimeMessage.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		mimeMessage.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	mimeMessage.WriteString("\r\n")
	mimeMessage.WriteString(message.Body)

	return mimeMessage.String()
}

// buildEmailVerificationMessage construye mensaje para verificación de email
func (es *ServiceImpl) buildEmailVerificationMessage(data *TemplateData) (*Message, error) {
	if !es.config.Enabled {
		// Simular envío en desarrollo
		fmt.Printf("📧 [SIMULADO] Email de verificación enviado a: %s\n", data.Email)
		fmt.Printf("🔗 Link de verificación: %s\n", data.URL)
		return nil, nil
	}

	if !es.config.CanSendEmails() {
		return nil, errors.WrapError(errors.ErrEmailService, "configuración de email incompleta")
	}

	return &Message{
		To:      []string{data.Email},
		Subject: "Verifica tu cuenta en MisViáticos",
		Body:    es.buildEmailVerificationHTML(data.FullName, data.URL),
		IsHTML:  true,
	}, nil
}

// buildPasswordResetMessage construye mensaje para reset de contraseña
func (es *ServiceImpl) buildPasswordResetMessage(data *TemplateData) (*Message, error) {
	if !es.config.Enabled {
		fmt.Printf("📧 [SIMULADO] Email de reset enviado a: %s\n", data.Email)
		fmt.Printf("🔗 Link de reset: %s\n", data.URL)
		return nil, nil
	}

	if !es.config.CanSendEmails() {
		return nil, errors.WrapError(errors.ErrEmailService, "configuración de email incompleta")
	}

	return &Message{
		To:      []string{data.Email},
		Subject: "Recupera tu contraseña - MisViáticos",
		Body:    es.buildPasswordResetHTML(data.FullName, data.URL),
		IsHTML:  true,
	}, nil
}

// buildWelcomeMessage construye mensaje de bienvenida
func (es *ServiceImpl) buildWelcomeMessage(data *TemplateData) (*Message, error) {
	if !es.config.Enabled {
		fmt.Printf("📧 [SIMULADO] Email de bienvenida enviado a: %s\n", data.Email)
		return nil, nil
	}

	if !es.config.CanSendEmails() {
		return nil, errors.WrapError(errors.ErrEmailService, "configuración de email incompleta")
	}

	return &Message{
		To:      []string{data.Email},
		Subject: "¡Bienvenido a MisViáticos!",
		Body:    es.buildWelcomeEmailHTML(data.FullName),
		IsHTML:  true,
	}, nil
}

// buildGenericMessage construye mensaje genérico
func (es *ServiceImpl) buildGenericMessage(data *TemplateData) (*Message, error) {
	if !es.config.Enabled {
		fmt.Printf("📧 [SIMULADO] Email genérico enviado a: %s\n", data.Email)
		return nil, nil
	}

	if !es.config.CanSendEmails() {
		return nil, errors.WrapError(errors.ErrEmailService, "configuración de email incompleta")
	}

	return &Message{
		To:      []string{data.Email},
		Subject: "Notificación - MisViáticos",
		Body:    data.Message,
		IsHTML:  false,
	}, nil
}

// buildEmailVerificationHTML construye el HTML para email de verificación
func (es *ServiceImpl) buildEmailVerificationHTML(fullName, verificationURL string) string {
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
func (es *ServiceImpl) buildPasswordResetHTML(fullName, resetURL string) string {
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
func (es *ServiceImpl) buildWelcomeEmailHTML(fullName string) string {
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
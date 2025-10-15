package email

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
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
	config    *config.EmailConfig
	templates map[EmailTemplate]*template.Template
}

// NewService crea una nueva instancia del servicio genérico de email
func NewService() Service {
	emailConfig := config.NewEmailConfig()

	// Validar configuración al inicializar
	if emailConfig.Enabled && !emailConfig.CanSendEmails() {
		log.Printf("⚠️  Email habilitado pero configuración incompleta para proveedor: %s", emailConfig.Provider)
	}

	// Cargar templates HTML
	templates := loadEmailTemplates()

	return &ServiceImpl{
		config:    emailConfig,
		templates: templates,
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

	body, err := es.renderTemplate(TemplateEmailVerification, data)
	if err != nil {
		return nil, err
	}

	return &Message{
		To:      []string{data.Email},
		Subject: "Verifica tu cuenta en MisViáticos",
		Body:    body,
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

	body, err := es.renderTemplate(TemplatePasswordReset, data)
	if err != nil {
		return nil, err
	}

	return &Message{
		To:      []string{data.Email},
		Subject: "Recupera tu contraseña - MisViáticos",
		Body:    body,
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

	body, err := es.renderTemplate(TemplateWelcome, data)
	if err != nil {
		return nil, err
	}

	return &Message{
		To:      []string{data.Email},
		Subject: "¡Bienvenido a MisViáticos!",
		Body:    body,
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



// loadEmailTemplates carga los templates HTML externos
func loadEmailTemplates() map[EmailTemplate]*template.Template {
	templates := make(map[EmailTemplate]*template.Template)
	
	// Buscar la raíz del proyecto
	projectRoot := findProjectRoot()
	
	templateFiles := map[EmailTemplate]string{
		TemplateEmailVerification: "email_templates/verification.html",
		TemplatePasswordReset:     "email_templates/password_reset.html", 
		TemplateWelcome:           "email_templates/welcome.html",
	}
	
	for templateType, relativeFilePath := range templateFiles {
		// Construir ruta absoluta
		filePath := filepath.Join(projectRoot, relativeFilePath)
		
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Printf("⚠️ Template no encontrado: %s, usando template por defecto", filePath)
			continue
		}
		
		tmpl, err := template.ParseFiles(filePath)
		if err != nil {
			log.Printf("❌ Error cargando template %s: %v", filePath, err)
			continue
		}
		
		templates[templateType] = tmpl
		log.Printf("✅ Template cargado: %s", filePath)
	}
	
	return templates
}

// findProjectRoot busca la raíz del proyecto buscando el archivo go.mod
func findProjectRoot() string {
	// Empezar desde el directorio actual
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("⚠️ Error obteniendo directorio actual: %v", err)
		return "."
	}
	
	// Buscar hacia arriba hasta encontrar go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			// Llegamos a la raíz del sistema de archivos
			break
		}
		dir = parent
	}
	
	log.Printf("⚠️ No se encontró go.mod, usando directorio actual")
	return "."
}

// renderTemplate renderiza un template con los datos proporcionados
func (es *ServiceImpl) renderTemplate(templateType EmailTemplate, data *TemplateData) (string, error) {
	tmpl, exists := es.templates[templateType]
	if !exists {
		log.Printf("⚠️ Template %v no encontrado, usando fallback mínimo", templateType)
		return es.getFallbackTemplate(templateType, data)
	}
	
	var buffer strings.Builder
	err := tmpl.Execute(&buffer, data)
	if err != nil {
		log.Printf("❌ Error renderizando template %v: %v", templateType, err)
		log.Printf("🔄 Usando fallback mínimo para template %v", templateType)
		return es.getFallbackTemplate(templateType, data)
	}
	
	return buffer.String(), nil
}

// getFallbackTemplate proporciona templates mínimos como fallback
func (es *ServiceImpl) getFallbackTemplate(templateType EmailTemplate, data *TemplateData) (string, error) {
	switch templateType {
	case TemplateEmailVerification:
		return fmt.Sprintf(`<h2>Hola %s</h2><p>Verifica tu cuenta haciendo clic aquí: <a href="%s">%s</a></p>`, data.FullName, data.URL, data.URL), nil
	case TemplatePasswordReset:
		return fmt.Sprintf(`<h2>Hola %s</h2><p>Restablece tu contraseña haciendo clic aquí: <a href="%s">%s</a></p>`, data.FullName, data.URL, data.URL), nil
	case TemplateWelcome:
		return fmt.Sprintf(`<h2>¡Bienvenido %s!</h2><p>Tu cuenta ha sido activada exitosamente.</p>`, data.FullName), nil
	default:
		return "", errors.NewValidationError("template no soportado", "template_type")
	}
}


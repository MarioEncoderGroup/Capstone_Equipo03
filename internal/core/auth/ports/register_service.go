package ports

import (
	"context"
	"github.com/google/uuid"
)

// RegisterRequest representa la solicitud de registro de usuario
type RegisterRequest struct {
	// Datos básicos del usuario
	Username    string `json:"username" validate:"required,min=3,max=50"`
	FullName    string `json:"full_name" validate:"required,min=2,max=200"`
	Email       string `json:"email" validate:"required,email,max=150"`
	Password    string `json:"password" validate:"required,min=8,max=255"`
	Phone       string `json:"phone,omitempty" validate:"omitempty,max=20"`
	
	// RUT chileno (opcional para usuarios individuales)
	IdentificationNumber string `json:"identification_number,omitempty" validate:"omitempty,max=50"`
	
	// Datos del tenant/empresa (opcional - solo para registro de empresas)
	CreateTenant bool                    `json:"create_tenant,omitempty"`
	TenantData   *TenantRegistrationData `json:"tenant_data,omitempty"`
}

// TenantRegistrationData representa los datos para crear un tenant
type TenantRegistrationData struct {
	RUT          string    `json:"rut" validate:"required,max=20"`
	BusinessName string    `json:"business_name" validate:"required,min=2,max=150"`
	Email        string    `json:"email" validate:"required,email,max=150"`
	Phone        string    `json:"phone" validate:"required,max=20"`
	Address      string    `json:"address" validate:"required,max=200"`
	Website      string    `json:"website" validate:"required,max=150"`
	RegionID     string    `json:"region_id" validate:"required,len=2"`
	CommuneID    string    `json:"commune_id" validate:"required,max=100"`
	CountryID    uuid.UUID `json:"country_id" validate:"required"`
}

// RegisterResponse representa la respuesta del registro
type RegisterResponse struct {
	UserID       uuid.UUID  `json:"user_id"`
	TenantID     *uuid.UUID `json:"tenant_id,omitempty"`
	Message      string     `json:"message"`
	RequiresEmailVerification bool `json:"requires_email_verification"`
}

// RegisterService define el contrato para el servicio de registro
type RegisterService interface {
	// RegisterUser registra un nuevo usuario individual
	RegisterUser(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	
	// RegisterUserWithTenant registra un usuario y crea un tenant/empresa
	RegisterUserWithTenant(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	
	// VerifyEmail verifica el email de un usuario usando el token
	VerifyEmail(ctx context.Context, token string) error
	
	// ResendEmailVerification reenvía el email de verificación
	ResendEmailVerification(ctx context.Context, email string) error
	
	// ValidateRegistrationData valida los datos de registro
	ValidateRegistrationData(req *RegisterRequest) error
}
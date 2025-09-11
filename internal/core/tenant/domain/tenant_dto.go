package domain_tenant

import (
	domain_user "github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	"github.com/google/uuid"
)

type CreateTenantDTO struct {
	Rut          string  `form:"rut" json:"rut" validate:"required"`
	BusinessName string  `form:"business_name" json:"business_name" validate:"required,min=3,max=100"`
	Email        string  `form:"email" json:"email" validate:"required,email"`
	Phone        string  `form:"phone" json:"phone" validate:"required,min=9,max=15"`
	Address      string  `form:"address" json:"address" validate:"required,min=3,max=100"`
	Website      *string `form:"website" json:"website,omitempty"`
	Logo         *string `form:"logo" json:"logo,omitempty"`
	RegionID     string  `form:"region_id" json:"region_id" validate:"required"`
	CommuneID    string  `form:"commune_id" json:"commune_id" validate:"required"`
	CountryID    string  `form:"country_id" json:"country_id" validate:"required"`
}

type UpdateTenantDTO struct {
	Rut          string  `form:"rut" json:"rut" validate:"required"`
	BusinessName string  `form:"business_name" json:"business_name" validate:"required,min=3,max=100"`
	Email        string  `form:"email" json:"email" validate:"required,email"`
	Phone        string  `form:"phone" json:"phone" validate:"required,min=9,max=15"`
	Address      string  `form:"address" json:"address" validate:"required,min=3,max=100"`
	Website      *string `form:"website" json:"website,omitempty"`
	Logo         *string `form:"logo" json:"logo,omitempty"`
	RegionID     string  `form:"region_id" json:"region_id" validate:"required"`
	CommuneID    string  `form:"commune_id" json:"commune_id" validate:"required"`
	CountryID    string  `form:"country_id" json:"country_id" validate:"required"`
	Status       string  `form:"status" json:"status" validate:"required,oneof=active inactive suspended"`
}

type SelectTenantDto struct {
	TenantID uuid.UUID `json:"tenant_id" validate:"required,uuid"`
}

type SelectTenantResponseDto struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	ExpiresIn    int64            `json:"expires_in"`
	User         domain_user.User `json:"user"`
	Tenant       *Tenant          `json:"tenant,omitempty"`
}

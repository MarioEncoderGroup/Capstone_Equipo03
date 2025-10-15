package domain

import (
	"time"
	"github.com/google/uuid"
)

// CreateUserDto represents the data required to create a user
type CreateUserDto struct {
	FullName             string      `json:"full_name" validate:"required,min=2,max=100"`
	Email                string      `json:"email" validate:"required,email"`
	Phone                *string     `json:"phone,omitempty" validate:"omitempty,min=8,max=15"`
	IdentificationNumber *string     `json:"identification_number,omitempty" validate:"omitempty,min=10,max=12"`
	Password             string      `json:"password" validate:"required,min=8,max=128"`
	IsActive             *bool       `json:"is_active,omitempty"` // Puntero para diferenciar nil (default true) de false (explícitamente inactivo)
	RoleIDs              []uuid.UUID `json:"role_ids,omitempty"`  // IDs de roles a asignar al usuario
}

// UpdateUserDto represents the data to update a user
type UpdateUserDto struct {
	FullName             *string `json:"full_name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone                *string `json:"phone,omitempty" validate:"omitempty,min=8,max=15"`
	IdentificationNumber *string `json:"identification_number,omitempty" validate:"omitempty,min=10,max=12"`
	BankID               *string `json:"bank_id,omitempty" validate:"omitempty,uuid"`
	BankAccountNumber    *string `json:"bank_account_number,omitempty" validate:"omitempty,min=10,max=20"`
	BankAccountType      *string `json:"bank_account_type,omitempty" validate:"omitempty,oneof=checking savings"`
	ImageURL             *string `json:"image_url,omitempty" validate:"omitempty,url"`
	IsActive             *bool   `json:"is_active,omitempty"`
}

// UserResponseDto represents the response with user information (without sensitive fields)
type UserResponseDto struct {
	ID                   uuid.UUID  `json:"id"`
	Username             string     `json:"username"`
	FullName             string     `json:"full_name"`
	Phone                *string    `json:"phone,omitempty"`
	IdentificationNumber *string    `json:"identification_number,omitempty"`
	Email                string     `json:"email"`
	EmailVerified        bool       `json:"email_verified"`
	BankID               *uuid.UUID `json:"bank_id,omitempty"`
	BankAccountNumber    *string    `json:"bank_account_number,omitempty"`
	BankAccountType      *string    `json:"bank_account_type,omitempty"`
	ImageURL             *string    `json:"image_url,omitempty"`
	IsActive             bool       `json:"is_active"`
	LastLogin            *time.Time `json:"last_login,omitempty"`
	Created              time.Time  `json:"created"`
	Updated              time.Time  `json:"updated"`
}

// RoleResponseDto representa información básica de un rol para incluir en listados
type RoleResponseDto struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
}

// UserListResponseDto representa un usuario en listados (información mínima)
type UserListResponseDto struct {
	ID            uuid.UUID          `json:"id"`
	Username      string             `json:"username"`
	FullName      string             `json:"full_name"`
	Email         string             `json:"email"`
	Phone         *string            `json:"phone,omitempty"`
	EmailVerified bool               `json:"email_verified"`
	IsActive      bool               `json:"is_active"`
	Roles         []RoleResponseDto  `json:"roles,omitempty"`
	LastLogin     *time.Time         `json:"last_login,omitempty"`
	Created       time.Time          `json:"created"`
}

// ChangePasswordDto represents the data to change password
type ChangePasswordDto struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// UpdateProfileDto represents the data to update the authenticated user's profile
type UpdateProfileDto struct {
	FullName             *string `json:"full_name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone                *string `json:"phone,omitempty" validate:"omitempty,min=8,max=15"`
	IdentificationNumber *string `json:"identification_number,omitempty" validate:"omitempty,min=10,max=12"`
	BankID               *string `json:"bank_id,omitempty" validate:"omitempty,uuid"`
	BankAccountNumber    *string `json:"bank_account_number,omitempty" validate:"omitempty,min=10,max=20"`
	BankAccountType      *string `json:"bank_account_type,omitempty" validate:"omitempty,oneof=checking savings"`
	ImageURL             *string `json:"image_url,omitempty" validate:"omitempty,url"`
}

// ToUserResponseDto converts a User to UserResponseDto (without sensitive fields)
func (u *User) ToUserResponseDto() *UserResponseDto {
	return &UserResponseDto{
		ID:                   u.ID,
		Username:             u.Username,
		FullName:             u.FullName,
		Phone:                u.Phone,
		IdentificationNumber: u.IdentificationNumber,
		Email:                u.Email,
		EmailVerified:        u.EmailVerified,
		BankID:               u.BankID,
		BankAccountNumber:    u.BankAccountNumber,
		BankAccountType:      u.BankAccountType,
		ImageURL:             u.ImageURL,
		IsActive:             u.IsActive,
		LastLogin:            u.LastLogin,
		Created:              u.Created,
		Updated:              u.Updated,
	}
}

// ToUserListResponseDto convierte un User en UserListResponseDto (información mínima)
func (u *User) ToUserListResponseDto() *UserListResponseDto {
	return &UserListResponseDto{
		ID:            u.ID,
		Username:      u.Username,
		FullName:      u.FullName,
		Email:         u.Email,
		Phone:         u.Phone,
		EmailVerified: u.EmailVerified,
		IsActive:      u.IsActive,
		Roles:         []RoleResponseDto{}, // Se cargan después en el controller
		LastLogin:     u.LastLogin,
		Created:       u.Created,
	}
}

// UsersListResponseDto representa una lista paginada de usuarios
type UsersListResponseDto struct {
	Users []UserListResponseDto `json:"users"`
	Total int64                 `json:"total"`
}
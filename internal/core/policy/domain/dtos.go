package domain

import (
	"github.com/google/uuid"
)

// CreatePolicyDto representa los datos para crear una política
type CreatePolicyDto struct {
	Name        string     `json:"name" validate:"required,min=3,max=150"`
	Description string     `json:"description" validate:"omitempty,max=500"`
	PolicyType  string     `json:"policy_type" validate:"required,oneof=travel daily project"`
	Config      any        `json:"config,omitempty"`
	CreatedBy   uuid.UUID  `json:"created_by" validate:"required,uuid"`
}

// UpdatePolicyDto representa los datos para actualizar una política
type UpdatePolicyDto struct {
	Name        *string `json:"name" validate:"omitempty,min=3,max=150"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	PolicyType  *string `json:"policy_type" validate:"omitempty,oneof=travel daily project"`
	IsActive    *bool   `json:"is_active"`
	Config      any     `json:"config,omitempty"`
}

// CreatePolicyRuleDto representa los datos para crear una regla de política
type CreatePolicyRuleDto struct {
	PolicyID   uuid.UUID `json:"policy_id" validate:"required,uuid"`
	CategoryID uuid.UUID `json:"category_id" validate:"required,uuid"`
	RuleType   string    `json:"rule_type" validate:"required,oneof=limit auto_approve require_approval reject"`
	Condition  any       `json:"condition,omitempty"`
	Action     any       `json:"action,omitempty"`
	Priority   int       `json:"priority" validate:"gte=0"`
	IsActive   bool      `json:"is_active"`
}

// UpdatePolicyRuleDto representa los datos para actualizar una regla de política
type UpdatePolicyRuleDto struct {
	CategoryID *uuid.UUID `json:"category_id" validate:"omitempty,uuid"`
	RuleType   *string    `json:"rule_type" validate:"omitempty,oneof=limit auto_approve require_approval reject"`
	Condition  any        `json:"condition,omitempty"`
	Action     any        `json:"action,omitempty"`
	Priority   *int       `json:"priority" validate:"omitempty,gte=0"`
	IsActive   *bool      `json:"is_active"`
}

// CreatePolicyApproverDto representa los datos para agregar un aprobador
type CreatePolicyApproverDto struct {
	PolicyID  uuid.UUID `json:"policy_id" validate:"required,uuid"`
	UserID    uuid.UUID `json:"user_id" validate:"required,uuid"`
	Level     int       `json:"level" validate:"required,gte=1"`
	AmountMin *float64  `json:"amount_min" validate:"omitempty,gte=0"`
	AmountMax *float64  `json:"amount_max" validate:"omitempty,gte=0"`
}

// CreatePolicySubmitterDto representa los datos para agregar un usuario que puede usar la política
type CreatePolicySubmitterDto struct {
	PolicyID   uuid.UUID  `json:"policy_id" validate:"required,uuid"`
	UserID     uuid.UUID  `json:"user_id" validate:"required,uuid"`
	RoleID     *uuid.UUID `json:"role_id" validate:"omitempty,uuid"`
	Department *string    `json:"department" validate:"omitempty,max=100"`
}

// PolicyFilters representa los filtros para buscar políticas
type PolicyFilters struct {
	PolicyType *string
	IsActive   *bool
	CreatedBy  *uuid.UUID
	Limit      int
	Offset     int
}

// PolicyResponse representa la respuesta con datos completos de una política
type PolicyResponse struct {
	Policy
	Rules      []PolicyRule      `json:"rules,omitempty"`
	Approvers  []PolicyApprover  `json:"approvers,omitempty"`
	Submitters []PolicySubmitter `json:"submitters,omitempty"`
}

// PoliciesResponse representa la respuesta para una lista de políticas
type PoliciesResponse struct {
	Policies []PolicyResponse `json:"policies"`
	Total    int              `json:"total"`
	Limit    int              `json:"limit"`
	Offset   int              `json:"offset"`
}

// ValidateExpenseDto representa los datos para validar un gasto
type ValidateExpenseDto struct {
	CategoryID  uuid.UUID `json:"category_id" validate:"required,uuid"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	ExpenseDate string    `json:"expense_date" validate:"required"` // formato: 2006-01-02
	Currency    string    `json:"currency" validate:"omitempty,oneof=CLP USD EUR"`
	Description string    `json:"description" validate:"omitempty,max=500"`
}

// ValidationResponse representa la respuesta de validación de un gasto
type ValidationResponse struct {
	IsValid          bool           `json:"is_valid"`
	Violations       []Violation    `json:"violations,omitempty"`
	RequiresApproval bool           `json:"requires_approval"`
	ApprovalLevel    int            `json:"approval_level,omitempty"`
	Approvers        []ApproverInfo `json:"approvers,omitempty"`
}

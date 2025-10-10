package domain

import (
	"time"

	"github.com/google/uuid"
)

// Policy representa una política de gastos empresariales
type Policy struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	PolicyType  string     `json:"policy_type"` // travel, daily, project
	IsActive    bool       `json:"is_active"`
	Config      any        `json:"config,omitempty"`       // JSONB flexible configuration
	CreatedBy   uuid.UUID  `json:"created_by"`
	Created     time.Time  `json:"created_at"`
	Updated     time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	// Relaciones
	Rules       []PolicyRule       `json:"rules,omitempty"`
	Approvers   []PolicyApprover   `json:"approvers,omitempty"`
	Submitters  []PolicySubmitter  `json:"submitters,omitempty"`
}

// PolicyRule representa una regla de validación de una política
type PolicyRule struct {
	ID         uuid.UUID  `json:"id"`
	PolicyID   uuid.UUID  `json:"policy_id"`
	CategoryID uuid.UUID  `json:"category_id"`
	RuleType   string     `json:"rule_type"` // limit, auto_approve, require_approval, reject
	Condition  any        `json:"condition,omitempty"` // JSONB conditions
	Action     any        `json:"action,omitempty"`    // JSONB actions
	Priority   int        `json:"priority"`
	IsActive   bool       `json:"is_active"`
	Created    time.Time  `json:"created_at"`
	Updated    time.Time  `json:"updated_at"`
}

// PolicyApprover representa un aprobador de política con niveles
type PolicyApprover struct {
	ID        uuid.UUID  `json:"id"`
	PolicyID  uuid.UUID  `json:"policy_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Level     int        `json:"level"` // Nivel de aprobación (1, 2, 3...)
	AmountMin *float64   `json:"amount_min,omitempty"`
	AmountMax *float64   `json:"amount_max,omitempty"`
	Created   time.Time  `json:"created_at"`
}

// PolicySubmitter representa quién puede usar la política
type PolicySubmitter struct {
	ID         uuid.UUID  `json:"id"`
	PolicyID   uuid.UUID  `json:"policy_id"`
	UserID     uuid.UUID  `json:"user_id"`
	RoleID     *uuid.UUID `json:"role_id,omitempty"`
	Department *string    `json:"department,omitempty"`
	Created    time.Time  `json:"created_at"`
}

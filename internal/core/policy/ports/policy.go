package ports

import (
	"context"

	"github.com/JoseLuis21/mv-backend/internal/core/policy/domain"
	"github.com/google/uuid"
)

// PolicyRepository define las operaciones de persistencia para políticas
type PolicyRepository interface {
	// Policy operations
	Create(ctx context.Context, policy *domain.Policy) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Policy, error)
	GetAll(ctx context.Context, filters domain.PolicyFilters) ([]domain.Policy, int, error)
	Update(ctx context.Context, policy *domain.Policy) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Rule operations
	CreateRule(ctx context.Context, rule *domain.PolicyRule) error
	GetRuleByID(ctx context.Context, id uuid.UUID) (*domain.PolicyRule, error)
	GetRulesByPolicy(ctx context.Context, policyID uuid.UUID) ([]domain.PolicyRule, error)
	UpdateRule(ctx context.Context, rule *domain.PolicyRule) error
	DeleteRule(ctx context.Context, id uuid.UUID) error

	// Approver operations
	CreateApprover(ctx context.Context, approver *domain.PolicyApprover) error
	GetApproversByPolicy(ctx context.Context, policyID uuid.UUID) ([]domain.PolicyApprover, error)
	DeleteApprover(ctx context.Context, id uuid.UUID) error

	// Submitter operations
	CreateSubmitter(ctx context.Context, submitter *domain.PolicySubmitter) error
	GetSubmittersByPolicy(ctx context.Context, policyID uuid.UUID) ([]domain.PolicySubmitter, error)
	DeleteSubmitter(ctx context.Context, id uuid.UUID) error
}

// PolicyService define la lógica de negocio para políticas
type PolicyService interface {
	// Policy operations
	Create(ctx context.Context, dto domain.CreatePolicyDto) (*domain.Policy, error)
	GetAll(ctx context.Context, filters domain.PolicyFilters) ([]domain.Policy, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Policy, error)
	Update(ctx context.Context, id uuid.UUID, dto domain.UpdatePolicyDto) (*domain.Policy, error)
	Delete(ctx context.Context, id uuid.UUID) error

	// Rule operations
	AddRule(ctx context.Context, dto domain.CreatePolicyRuleDto) (*domain.PolicyRule, error)
	UpdateRule(ctx context.Context, id uuid.UUID, dto domain.UpdatePolicyRuleDto) (*domain.PolicyRule, error)
	DeleteRule(ctx context.Context, id uuid.UUID) error

	// Approver operations
	AddApprover(ctx context.Context, dto domain.CreatePolicyApproverDto) (*domain.PolicyApprover, error)
	RemoveApprover(ctx context.Context, id uuid.UUID) error

	// Submitter operations
	AddSubmitter(ctx context.Context, dto domain.CreatePolicySubmitterDto) (*domain.PolicySubmitter, error)
	RemoveSubmitter(ctx context.Context, id uuid.UUID) error
}

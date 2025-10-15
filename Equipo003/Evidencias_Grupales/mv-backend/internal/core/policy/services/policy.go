package services

import (
	"context"
	"fmt"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/core/policy/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/policy/ports"
	"github.com/google/uuid"
)

type policyService struct {
	repo ports.PolicyRepository
}

// NewPolicyService creates a new policy service
func NewPolicyService(repo ports.PolicyRepository) ports.PolicyService {
	return &policyService{
		repo: repo,
	}
}

// Create creates a new policy
func (s *policyService) Create(ctx context.Context, dto domain.CreatePolicyDto) (*domain.Policy, error) {
	policy := &domain.Policy{
		ID:          uuid.New(),
		Name:        dto.Name,
		Description: &dto.Description,
		PolicyType:  dto.PolicyType,
		IsActive:    true,
		Config:      dto.Config,
		CreatedBy:   dto.CreatedBy,
		Created:     time.Now(),
		Updated:     time.Now(),
	}

	if err := s.repo.Create(ctx, policy); err != nil {
		return nil, fmt.Errorf("failed to create policy: %w", err)
	}

	return policy, nil
}

// GetAll retrieves all policies with filters
func (s *policyService) GetAll(ctx context.Context, filters domain.PolicyFilters) ([]domain.Policy, int, error) {
	policies, total, err := s.repo.GetAll(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get policies: %w", err)
	}

	return policies, total, nil
}

// GetByID retrieves a policy by ID
func (s *policyService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Policy, error) {
	policy, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy: %w", err)
	}

	return policy, nil
}

// Update updates an existing policy
func (s *policyService) Update(ctx context.Context, id uuid.UUID, dto domain.UpdatePolicyDto) (*domain.Policy, error) {
	policy, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy: %w", err)
	}

	if dto.Name != nil {
		policy.Name = *dto.Name
	}
	if dto.Description != nil {
		policy.Description = dto.Description
	}
	if dto.PolicyType != nil {
		policy.PolicyType = *dto.PolicyType
	}
	if dto.IsActive != nil {
		policy.IsActive = *dto.IsActive
	}
	if dto.Config != nil {
		policy.Config = dto.Config
	}

	policy.Updated = time.Now()

	if err := s.repo.Update(ctx, policy); err != nil {
		return nil, fmt.Errorf("failed to update policy: %w", err)
	}

	return policy, nil
}

// Delete soft deletes a policy
func (s *policyService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}

	return nil
}

// AddRule adds a new rule to a policy
func (s *policyService) AddRule(ctx context.Context, dto domain.CreatePolicyRuleDto) (*domain.PolicyRule, error) {
	// Verify policy exists
	if _, err := s.repo.GetByID(ctx, dto.PolicyID); err != nil {
		return nil, fmt.Errorf("policy not found: %w", err)
	}

	rule := &domain.PolicyRule{
		ID:         uuid.New(),
		PolicyID:   dto.PolicyID,
		CategoryID: dto.CategoryID,
		RuleType:   dto.RuleType,
		Condition:  dto.Condition,
		Action:     dto.Action,
		Priority:   dto.Priority,
		IsActive:   dto.IsActive,
		Created:    time.Now(),
		Updated:    time.Now(),
	}

	if err := s.repo.CreateRule(ctx, rule); err != nil {
		return nil, fmt.Errorf("failed to create rule: %w", err)
	}

	return rule, nil
}

// UpdateRule updates an existing rule
func (s *policyService) UpdateRule(ctx context.Context, id uuid.UUID, dto domain.UpdatePolicyRuleDto) (*domain.PolicyRule, error) {
	rule, err := s.repo.GetRuleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get rule: %w", err)
	}

	if dto.CategoryID != nil {
		rule.CategoryID = *dto.CategoryID
	}
	if dto.RuleType != nil {
		rule.RuleType = *dto.RuleType
	}
	if dto.Condition != nil {
		rule.Condition = dto.Condition
	}
	if dto.Action != nil {
		rule.Action = dto.Action
	}
	if dto.Priority != nil {
		rule.Priority = *dto.Priority
	}
	if dto.IsActive != nil {
		rule.IsActive = *dto.IsActive
	}

	rule.Updated = time.Now()

	if err := s.repo.UpdateRule(ctx, rule); err != nil {
		return nil, fmt.Errorf("failed to update rule: %w", err)
	}

	return rule, nil
}

// DeleteRule deletes a rule
func (s *policyService) DeleteRule(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteRule(ctx, id); err != nil {
		return fmt.Errorf("failed to delete rule: %w", err)
	}

	return nil
}

// AddApprover adds an approver to a policy
func (s *policyService) AddApprover(ctx context.Context, dto domain.CreatePolicyApproverDto) (*domain.PolicyApprover, error) {
	// Verify policy exists
	if _, err := s.repo.GetByID(ctx, dto.PolicyID); err != nil {
		return nil, fmt.Errorf("policy not found: %w", err)
	}

	approver := &domain.PolicyApprover{
		ID:        uuid.New(),
		PolicyID:  dto.PolicyID,
		UserID:    dto.UserID,
		Level:     dto.Level,
		AmountMin: dto.AmountMin,
		AmountMax: dto.AmountMax,
		Created:   time.Now(),
	}

	if err := s.repo.CreateApprover(ctx, approver); err != nil {
		return nil, fmt.Errorf("failed to create approver: %w", err)
	}

	return approver, nil
}

// RemoveApprover removes an approver from a policy
func (s *policyService) RemoveApprover(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteApprover(ctx, id); err != nil {
		return fmt.Errorf("failed to remove approver: %w", err)
	}

	return nil
}

// AddSubmitter adds a submitter to a policy
func (s *policyService) AddSubmitter(ctx context.Context, dto domain.CreatePolicySubmitterDto) (*domain.PolicySubmitter, error) {
	// Verify policy exists
	if _, err := s.repo.GetByID(ctx, dto.PolicyID); err != nil {
		return nil, fmt.Errorf("policy not found: %w", err)
	}

	submitter := &domain.PolicySubmitter{
		ID:         uuid.New(),
		PolicyID:   dto.PolicyID,
		UserID:     dto.UserID,
		RoleID:     dto.RoleID,
		Department: dto.Department,
		Created:    time.Now(),
	}

	if err := s.repo.CreateSubmitter(ctx, submitter); err != nil {
		return nil, fmt.Errorf("failed to create submitter: %w", err)
	}

	return submitter, nil
}

// RemoveSubmitter removes a submitter from a policy
func (s *policyService) RemoveSubmitter(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteSubmitter(ctx, id); err != nil {
		return fmt.Errorf("failed to remove submitter: %w", err)
	}

	return nil
}

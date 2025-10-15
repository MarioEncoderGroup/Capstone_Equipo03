package services

import (
	"context"
	"fmt"
	"sort"

	"github.com/JoseLuis21/mv-backend/internal/core/policy/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/policy/ports"
)

type ruleEngine struct {
	policyRepo ports.PolicyRepository
}

// NewRuleEngine creates a new rule engine
func NewRuleEngine(policyRepo ports.PolicyRepository) ports.RuleEngine {
	return &ruleEngine{
		policyRepo: policyRepo,
	}
}

// ValidateExpense valida un gasto contra una política y retorna violaciones
func (e *ruleEngine) ValidateExpense(ctx context.Context, expense *domain.ExpenseValidationInput, policy *domain.Policy) ([]domain.Violation, error) {
	violations := []domain.Violation{}

	// Get rules for this policy
	rules, err := e.policyRepo.GetRulesByPolicy(ctx, policy.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy rules: %w", err)
	}

	// Sort rules by priority (higher priority first)
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})

	// Evaluate each active rule
	for _, rule := range rules {
		if !rule.IsActive {
			continue
		}

		// Check if rule applies to this category
		if rule.CategoryID != expense.CategoryID {
			continue
		}

		// Evaluate rule based on type
		switch rule.RuleType {
		case "limit":
			if v := e.evaluateLimitRule(expense, &rule); v != nil {
				violations = append(violations, *v)
			}

		case "reject":
			if v := e.evaluateRejectRule(expense, &rule); v != nil {
				violations = append(violations, *v)
			}

		case "require_approval":
			// This is handled in CheckApprovalRequired
			continue

		case "auto_approve":
			// This is handled in CheckApprovalRequired
			continue
		}
	}

	return violations, nil
}

// CheckApprovalRequired determina si un gasto requiere aprobación y el nivel necesario
func (e *ruleEngine) CheckApprovalRequired(ctx context.Context, expense *domain.ExpenseValidationInput, policy *domain.Policy) (bool, int, error) {
	// Get rules for this policy
	rules, err := e.policyRepo.GetRulesByPolicy(ctx, policy.ID)
	if err != nil {
		return false, 0, fmt.Errorf("failed to get policy rules: %w", err)
	}

	// Sort rules by priority (higher priority first)
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})

	// Check for auto-approve rules first (highest priority)
	for _, rule := range rules {
		if !rule.IsActive || rule.CategoryID != expense.CategoryID {
			continue
		}

		if rule.RuleType == "auto_approve" {
			// Check if conditions are met
			if e.evaluateCondition(expense, rule.Condition) {
				return false, 0, nil // Auto-approved, no approval needed
			}
		}
	}

	// Check for require_approval rules
	maxLevel := 0
	requiresApproval := false

	for _, rule := range rules {
		if !rule.IsActive || rule.CategoryID != expense.CategoryID {
			continue
		}

		if rule.RuleType == "require_approval" {
			// Check if conditions are met
			if e.evaluateCondition(expense, rule.Condition) {
				requiresApproval = true
				// Extract level from action if specified
				if level := e.extractApprovalLevel(rule.Action); level > maxLevel {
					maxLevel = level
				}
			}
		}
	}

	// If no specific rule, determine by amount using approvers
	if !requiresApproval {
		approvers, err := e.policyRepo.GetApproversByPolicy(ctx, policy.ID)
		if err != nil {
			return false, 0, fmt.Errorf("failed to get approvers: %w", err)
		}

		for _, approver := range approvers {
			// Check if amount falls within approver's range
			if e.amountInRange(expense.Amount, approver.AmountMin, approver.AmountMax) {
				requiresApproval = true
				if approver.Level > maxLevel {
					maxLevel = approver.Level
				}
			}
		}
	}

	// If still no level determined, use level 1 as default
	if requiresApproval && maxLevel == 0 {
		maxLevel = 1
	}

	return requiresApproval, maxLevel, nil
}

// GetApprovers retorna los aprobadores necesarios según el monto del gasto
func (e *ruleEngine) GetApprovers(ctx context.Context, expense *domain.ExpenseValidationInput, policy *domain.Policy) ([]domain.ApproverInfo, error) {
	approvers, err := e.policyRepo.GetApproversByPolicy(ctx, policy.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get approvers: %w", err)
	}

	result := []domain.ApproverInfo{}

	// Filter approvers by amount range
	for _, approver := range approvers {
		if e.amountInRange(expense.Amount, approver.AmountMin, approver.AmountMax) {
			result = append(result, domain.ApproverInfo{
				UserID:    approver.UserID,
				Level:     approver.Level,
				AmountMin: approver.AmountMin,
				AmountMax: approver.AmountMax,
			})
		}
	}

	// Sort by level
	sort.Slice(result, func(i, j int) bool {
		return result[i].Level < result[j].Level
	})

	return result, nil
}

// Helper functions

func (e *ruleEngine) evaluateLimitRule(expense *domain.ExpenseValidationInput, rule *domain.PolicyRule) *domain.Violation {
	// Extract limit from condition (expected format: {"max_amount": 50000})
	condition, ok := rule.Condition.(map[string]interface{})
	if !ok {
		return nil
	}

	maxAmount, ok := condition["max_amount"].(float64)
	if !ok {
		return nil
	}

	if expense.Amount > maxAmount {
		return &domain.Violation{
			Field:    "amount",
			Message:  fmt.Sprintf("Amount %.2f exceeds limit of %.2f for this category", expense.Amount, maxAmount),
			Severity: "error",
		}
	}

	return nil
}

func (e *ruleEngine) evaluateRejectRule(expense *domain.ExpenseValidationInput, rule *domain.PolicyRule) *domain.Violation {
	// Evaluate if conditions are met
	if e.evaluateCondition(expense, rule.Condition) {
		message := "Expense rejected by policy rule"

		// Extract custom message from action if available
		if action, ok := rule.Action.(map[string]interface{}); ok {
			if msg, ok := action["message"].(string); ok {
				message = msg
			}
		}

		return &domain.Violation{
			Field:    "expense",
			Message:  message,
			Severity: "error",
		}
	}

	return nil
}

func (e *ruleEngine) evaluateCondition(expense *domain.ExpenseValidationInput, condition interface{}) bool {
	// Simple condition evaluation
	// Expected format: {"min_amount": 1000, "max_amount": 50000}
	if condition == nil {
		return true // No condition means always true
	}

	condMap, ok := condition.(map[string]interface{})
	if !ok {
		return true
	}

	// Check min_amount
	if minAmount, ok := condMap["min_amount"].(float64); ok {
		if expense.Amount < minAmount {
			return false
		}
	}

	// Check max_amount
	if maxAmount, ok := condMap["max_amount"].(float64); ok {
		if expense.Amount > maxAmount {
			return false
		}
	}

	return true
}

func (e *ruleEngine) extractApprovalLevel(action interface{}) int {
	if action == nil {
		return 1
	}

	actionMap, ok := action.(map[string]interface{})
	if !ok {
		return 1
	}

	if level, ok := actionMap["approval_level"].(float64); ok {
		return int(level)
	}

	return 1
}

func (e *ruleEngine) amountInRange(amount float64, min, max *float64) bool {
	if min != nil && amount < *min {
		return false
	}

	if max != nil && amount > *max {
		return false
	}

	return true
}

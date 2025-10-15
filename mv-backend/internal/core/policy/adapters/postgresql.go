package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/core/policy/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/policy/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type policyRepository struct {
	client *postgresql.PostgresqlClient
}

// NewPolicyRepository creates a new PostgreSQL policy repository
func NewPolicyRepository(client *postgresql.PostgresqlClient) ports.PolicyRepository {
	return &policyRepository{
		client: client,
	}
}

// Create creates a new policy
func (r *policyRepository) Create(ctx context.Context, policy *domain.Policy) error {
	query := `
		INSERT INTO policies (id, name, description, policy_type, is_active, config, created_by, created, updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	configJSON, err := json.Marshal(policy.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = r.client.Exec(ctx, query,
		policy.ID,
		policy.Name,
		policy.Description,
		policy.PolicyType,
		policy.IsActive,
		configJSON,
		policy.CreatedBy,
		policy.Created,
		policy.Updated,
	)

	if err != nil {
		return fmt.Errorf("failed to create policy: %w", err)
	}

	return nil
}

// GetByID retrieves a policy by ID
func (r *policyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Policy, error) {
	query := `
		SELECT id, name, description, policy_type, is_active, config, created_by, created, updated, deleted_at
		FROM policies
		WHERE id = $1 AND deleted_at IS NULL
	`

	var policy domain.Policy
	var configJSON []byte

	err := r.client.QueryRow(ctx, query, id).Scan(
		&policy.ID,
		&policy.Name,
		&policy.Description,
		&policy.PolicyType,
		&policy.IsActive,
		&configJSON,
		&policy.CreatedBy,
		&policy.Created,
		&policy.Updated,
		&policy.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("policy not found")
		}
		return nil, fmt.Errorf("failed to get policy: %w", err)
	}

	if len(configJSON) > 0 {
		if err := json.Unmarshal(configJSON, &policy.Config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	return &policy, nil
}

// GetAll retrieves all policies with filters
func (r *policyRepository) GetAll(ctx context.Context, filters domain.PolicyFilters) ([]domain.Policy, int, error) {
	baseQuery := `
		SELECT id, name, description, policy_type, is_active, config, created_by, created, updated, deleted_at
		FROM policies
		WHERE deleted_at IS NULL
	`

	countQuery := `SELECT COUNT(*) FROM policies WHERE deleted_at IS NULL`

	args := []interface{}{}
	argPos := 1

	if filters.PolicyType != nil {
		baseQuery += fmt.Sprintf(" AND policy_type = $%d", argPos)
		countQuery += fmt.Sprintf(" AND policy_type = $%d", argPos)
		args = append(args, *filters.PolicyType)
		argPos++
	}

	if filters.IsActive != nil {
		baseQuery += fmt.Sprintf(" AND is_active = $%d", argPos)
		countQuery += fmt.Sprintf(" AND is_active = $%d", argPos)
		args = append(args, *filters.IsActive)
		argPos++
	}

	if filters.CreatedBy != nil {
		baseQuery += fmt.Sprintf(" AND created_by = $%d", argPos)
		countQuery += fmt.Sprintf(" AND created_by = $%d", argPos)
		args = append(args, *filters.CreatedBy)
		argPos++
	}

	// Get total count
	var total int
	err := r.client.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count policies: %w", err)
	}

	// Add pagination
	baseQuery += " ORDER BY created DESC"
	if filters.Limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, filters.Limit)
		argPos++
	}

	if filters.Offset > 0 {
		baseQuery += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, filters.Offset)
	}

	rows, err := r.client.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query policies: %w", err)
	}
	defer rows.Close()

	policies := []domain.Policy{}
	for rows.Next() {
		var policy domain.Policy
		var configJSON []byte

		err := rows.Scan(
			&policy.ID,
			&policy.Name,
			&policy.Description,
			&policy.PolicyType,
			&policy.IsActive,
			&configJSON,
			&policy.CreatedBy,
			&policy.Created,
			&policy.Updated,
			&policy.DeletedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan policy: %w", err)
		}

		if len(configJSON) > 0 {
			if err := json.Unmarshal(configJSON, &policy.Config); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal config: %w", err)
			}
		}

		policies = append(policies, policy)
	}

	return policies, total, nil
}

// Update updates an existing policy
func (r *policyRepository) Update(ctx context.Context, policy *domain.Policy) error {
	query := `
		UPDATE policies
		SET name = $1, description = $2, policy_type = $3, is_active = $4, config = $5, updated = $6
		WHERE id = $7 AND deleted_at IS NULL
	`

	configJSON, err := json.Marshal(policy.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = r.client.Exec(ctx, query,
		policy.Name,
		policy.Description,
		policy.PolicyType,
		policy.IsActive,
		configJSON,
		time.Now(),
		policy.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update policy: %w", err)
	}

	return nil
}

// Delete soft deletes a policy
func (r *policyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE policies
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	err := r.client.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}

	return nil
}

// CreateRule creates a new policy rule
func (r *policyRepository) CreateRule(ctx context.Context, rule *domain.PolicyRule) error {
	query := `
		INSERT INTO policy_rules (id, policy_id, category_id, rule_type, condition, action, priority, is_active, created, updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	conditionJSON, err := json.Marshal(rule.Condition)
	if err != nil {
		return fmt.Errorf("failed to marshal condition: %w", err)
	}

	actionJSON, err := json.Marshal(rule.Action)
	if err != nil {
		return fmt.Errorf("failed to marshal action: %w", err)
	}

	err = r.client.Exec(ctx, query,
		rule.ID,
		rule.PolicyID,
		rule.CategoryID,
		rule.RuleType,
		conditionJSON,
		actionJSON,
		rule.Priority,
		rule.IsActive,
		rule.Created,
		rule.Updated,
	)

	if err != nil {
		return fmt.Errorf("failed to create rule: %w", err)
	}

	return nil
}

// GetRuleByID retrieves a rule by ID
func (r *policyRepository) GetRuleByID(ctx context.Context, id uuid.UUID) (*domain.PolicyRule, error) {
	query := `
		SELECT id, policy_id, category_id, rule_type, condition, action, priority, is_active, created, updated
		FROM policy_rules
		WHERE id = $1
	`

	var rule domain.PolicyRule
	var conditionJSON, actionJSON []byte

	err := r.client.QueryRow(ctx, query, id).Scan(
		&rule.ID,
		&rule.PolicyID,
		&rule.CategoryID,
		&rule.RuleType,
		&conditionJSON,
		&actionJSON,
		&rule.Priority,
		&rule.IsActive,
		&rule.Created,
		&rule.Updated,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("rule not found")
		}
		return nil, fmt.Errorf("failed to get rule: %w", err)
	}

	if len(conditionJSON) > 0 {
		if err := json.Unmarshal(conditionJSON, &rule.Condition); err != nil {
			return nil, fmt.Errorf("failed to unmarshal condition: %w", err)
		}
	}

	if len(actionJSON) > 0 {
		if err := json.Unmarshal(actionJSON, &rule.Action); err != nil {
			return nil, fmt.Errorf("failed to unmarshal action: %w", err)
		}
	}

	return &rule, nil
}

// GetRulesByPolicy retrieves all rules for a policy
func (r *policyRepository) GetRulesByPolicy(ctx context.Context, policyID uuid.UUID) ([]domain.PolicyRule, error) {
	query := `
		SELECT id, policy_id, category_id, rule_type, condition, action, priority, is_active, created, updated
		FROM policy_rules
		WHERE policy_id = $1
		ORDER BY priority ASC
	`

	rows, err := r.client.Query(ctx, query, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query rules: %w", err)
	}
	defer rows.Close()

	rules := []domain.PolicyRule{}
	for rows.Next() {
		var rule domain.PolicyRule
		var conditionJSON, actionJSON []byte

		err := rows.Scan(
			&rule.ID,
			&rule.PolicyID,
			&rule.CategoryID,
			&rule.RuleType,
			&conditionJSON,
			&actionJSON,
			&rule.Priority,
			&rule.IsActive,
			&rule.Created,
			&rule.Updated,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan rule: %w", err)
		}

		if len(conditionJSON) > 0 {
			if err := json.Unmarshal(conditionJSON, &rule.Condition); err != nil {
				return nil, fmt.Errorf("failed to unmarshal condition: %w", err)
			}
		}

		if len(actionJSON) > 0 {
			if err := json.Unmarshal(actionJSON, &rule.Action); err != nil {
				return nil, fmt.Errorf("failed to unmarshal action: %w", err)
			}
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

// UpdateRule updates an existing rule
func (r *policyRepository) UpdateRule(ctx context.Context, rule *domain.PolicyRule) error {
	query := `
		UPDATE policy_rules
		SET category_id = $1, rule_type = $2, condition = $3, action = $4, priority = $5, is_active = $6, updated = $7
		WHERE id = $8
	`

	conditionJSON, err := json.Marshal(rule.Condition)
	if err != nil {
		return fmt.Errorf("failed to marshal condition: %w", err)
	}

	actionJSON, err := json.Marshal(rule.Action)
	if err != nil {
		return fmt.Errorf("failed to marshal action: %w", err)
	}

	err = r.client.Exec(ctx, query,
		rule.CategoryID,
		rule.RuleType,
		conditionJSON,
		actionJSON,
		rule.Priority,
		rule.IsActive,
		time.Now(),
		rule.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update rule: %w", err)
	}

	return nil
}

// DeleteRule deletes a rule
func (r *policyRepository) DeleteRule(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM policy_rules WHERE id = $1`

	err := r.client.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete rule: %w", err)
	}

	return nil
}

// CreateApprover creates a new policy approver
func (r *policyRepository) CreateApprover(ctx context.Context, approver *domain.PolicyApprover) error {
	query := `
		INSERT INTO policy_approvers (id, policy_id, user_id, level, amount_min, amount_max, created)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	err := r.client.Exec(ctx, query,
		approver.ID,
		approver.PolicyID,
		approver.UserID,
		approver.Level,
		approver.AmountMin,
		approver.AmountMax,
		approver.Created,
	)

	if err != nil {
		return fmt.Errorf("failed to create approver: %w", err)
	}

	return nil
}

// GetApproversByPolicy retrieves all approvers for a policy
func (r *policyRepository) GetApproversByPolicy(ctx context.Context, policyID uuid.UUID) ([]domain.PolicyApprover, error) {
	query := `
		SELECT id, policy_id, user_id, level, amount_min, amount_max, created
		FROM policy_approvers
		WHERE policy_id = $1
		ORDER BY level ASC
	`

	rows, err := r.client.Query(ctx, query, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query approvers: %w", err)
	}
	defer rows.Close()

	approvers := []domain.PolicyApprover{}
	for rows.Next() {
		var approver domain.PolicyApprover

		err := rows.Scan(
			&approver.ID,
			&approver.PolicyID,
			&approver.UserID,
			&approver.Level,
			&approver.AmountMin,
			&approver.AmountMax,
			&approver.Created,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan approver: %w", err)
		}

		approvers = append(approvers, approver)
	}

	return approvers, nil
}

// DeleteApprover deletes an approver
func (r *policyRepository) DeleteApprover(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM policy_approvers WHERE id = $1`

	err := r.client.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete approver: %w", err)
	}

	return nil
}

// CreateSubmitter creates a new policy submitter
func (r *policyRepository) CreateSubmitter(ctx context.Context, submitter *domain.PolicySubmitter) error {
	query := `
		INSERT INTO policy_submitters (id, policy_id, user_id, role_id, department, created)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	err := r.client.Exec(ctx, query,
		submitter.ID,
		submitter.PolicyID,
		submitter.UserID,
		submitter.RoleID,
		submitter.Department,
		submitter.Created,
	)

	if err != nil {
		return fmt.Errorf("failed to create submitter: %w", err)
	}

	return nil
}

// GetSubmittersByPolicy retrieves all submitters for a policy
func (r *policyRepository) GetSubmittersByPolicy(ctx context.Context, policyID uuid.UUID) ([]domain.PolicySubmitter, error) {
	query := `
		SELECT id, policy_id, user_id, role_id, department, created
		FROM policy_submitters
		WHERE policy_id = $1
	`

	rows, err := r.client.Query(ctx, query, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query submitters: %w", err)
	}
	defer rows.Close()

	submitters := []domain.PolicySubmitter{}
	for rows.Next() {
		var submitter domain.PolicySubmitter

		err := rows.Scan(
			&submitter.ID,
			&submitter.PolicyID,
			&submitter.UserID,
			&submitter.RoleID,
			&submitter.Department,
			&submitter.Created,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan submitter: %w", err)
		}

		submitters = append(submitters, submitter)
	}

	return submitters, nil
}

// DeleteSubmitter deletes a submitter
func (r *policyRepository) DeleteSubmitter(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM policy_submitters WHERE id = $1`

	err := r.client.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete submitter: %w", err)
	}

	return nil
}

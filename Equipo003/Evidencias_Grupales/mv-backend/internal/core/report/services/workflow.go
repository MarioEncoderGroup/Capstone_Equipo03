package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/policy/domain"
	policyPorts "github.com/JoseLuis21/mv-backend/internal/core/policy/ports"
	reportDomain "github.com/JoseLuis21/mv-backend/internal/core/report/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/report/ports"
)

var (
	ErrNoPolicyAssigned     = errors.New("report has no policy assigned")
	ErrPolicyNotFound       = errors.New("policy not found")
	ErrNoApproversFound     = errors.New("no approvers found for this report")
	ErrApprovalNotCompleted = errors.New("approval is not in a completed state")
)

type workflowEngine struct {
	reportRepo    ports.ReportRepository
	policyService policyPorts.PolicyService
	ruleEngine    policyPorts.RuleEngine
}

// NewWorkflowEngine creates a new instance of the workflow engine
func NewWorkflowEngine(
	reportRepo ports.ReportRepository,
	policyService policyPorts.PolicyService,
	ruleEngine policyPorts.RuleEngine,
) ports.WorkflowEngine {
	return &workflowEngine{
		reportRepo:    reportRepo,
		policyService: policyService,
		ruleEngine:    ruleEngine,
	}
}

// CreateApprovals creates the necessary approvals for a report based on policy and amount
func (w *workflowEngine) CreateApprovals(ctx context.Context, report *reportDomain.ExpenseReport) ([]reportDomain.Approval, error) {
	// Verify report has a policy assigned
	if report.PolicyID == nil {
		return nil, ErrNoPolicyAssigned
	}

	// Get the policy
	policy, err := w.policyService.GetByID(ctx, *report.PolicyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy: %w", err)
	}
	if policy == nil {
		return nil, ErrPolicyNotFound
	}

	// Create expense validation input from report
	expenseInput := &domain.ExpenseValidationInput{
		Amount:   report.TotalAmount,
		Currency: report.Currency,
		UserID:   report.UserID,
	}

	// Get approvers from rule engine
	approvers, err := w.ruleEngine.GetApprovers(ctx, expenseInput, policy)
	if err != nil {
		return nil, fmt.Errorf("failed to get approvers: %w", err)
	}

	if len(approvers) == 0 {
		return nil, ErrNoApproversFound
	}

	// Create approval entries
	now := time.Now()
	createdApprovals := []reportDomain.Approval{}

	for _, approver := range approvers {
		approval := reportDomain.Approval{
			ID:         uuid.New(),
			ReportID:   report.ID,
			ApproverID: approver.UserID,
			Level:      approver.Level,
			Status:     reportDomain.ApprovalStatusPending,
			Created:    now,
			Updated:    now,
		}

		if err := w.reportRepo.CreateApproval(ctx, &approval); err != nil {
			return nil, fmt.Errorf("failed to create approval: %w", err)
		}

		// Create history entry for approval creation
		history := &reportDomain.ApprovalHistory{
			ID:         uuid.New(),
			ApprovalID: approval.ID,
			ReportID:   report.ID,
			ActorID:    report.UserID, // Creator is the report submitter
			Action:     reportDomain.ApprovalActionCreated,
			NewStatus:  &approval.Status,
			Created:    now,
			Metadata: map[string]any{
				"level":       approver.Level,
				"approver_id": approver.UserID.String(),
			},
		}

		if err := w.reportRepo.CreateApprovalHistory(ctx, history); err != nil {
			return nil, fmt.Errorf("failed to create approval history: %w", err)
		}

		createdApprovals = append(createdApprovals, approval)
	}

	return createdApprovals, nil
}

// ProcessApproval processes an approval and determines the next step in the workflow
func (w *workflowEngine) ProcessApproval(ctx context.Context, approval *reportDomain.Approval) error {
	// Verify approval is in a completed state (approved or rejected)
	if approval.Status != reportDomain.ApprovalStatusApproved && approval.Status != reportDomain.ApprovalStatusRejected {
		return ErrApprovalNotCompleted
	}

	// Get all approvals for the report
	allApprovals, err := w.reportRepo.GetApprovalsByReport(ctx, approval.ReportID)
	if err != nil {
		return fmt.Errorf("failed to get approvals: %w", err)
	}

	// If rejected, mark report as rejected
	if approval.Status == reportDomain.ApprovalStatusRejected {
		if err := w.reportRepo.UpdateStatus(ctx, approval.ReportID, reportDomain.ReportStatusRejected); err != nil {
			return fmt.Errorf("failed to update report status to rejected: %w", err)
		}
		return nil
	}

	// If approved, check if all approvals at this level or below are complete
	currentLevel := approval.Level
	allPreviousLevelsApproved := true

	for _, a := range allApprovals {
		// Skip escalated approvals
		if a.Status == reportDomain.ApprovalStatusEscalated {
			continue
		}

		// Check if this approval is for current level or below
		if a.Level <= currentLevel {
			if a.Status == reportDomain.ApprovalStatusPending {
				allPreviousLevelsApproved = false
				break
			}
			if a.Status == reportDomain.ApprovalStatusRejected {
				// Should not happen as rejection cascades, but check anyway
				return nil
			}
		}
	}

	// If not all previous levels approved, wait
	if !allPreviousLevelsApproved {
		return nil
	}

	// Check if there are more levels
	hasNextLevel := false
	for _, a := range allApprovals {
		if a.Level > currentLevel && a.Status != reportDomain.ApprovalStatusEscalated {
			hasNextLevel = true
			break
		}
	}

	// If no more levels, mark report as approved
	if !hasNextLevel {
		if err := w.reportRepo.UpdateStatus(ctx, approval.ReportID, reportDomain.ReportStatusApproved); err != nil {
			return fmt.Errorf("failed to update report status to approved: %w", err)
		}

		// Update approval date
		report, err := w.reportRepo.GetByID(ctx, approval.ReportID)
		if err != nil {
			return fmt.Errorf("failed to get report: %w", err)
		}
		if report != nil {
			now := time.Now()
			report.ApprovalDate = &now
			report.Updated = now
			if err := w.reportRepo.Update(ctx, report); err != nil {
				return fmt.Errorf("failed to update report approval date: %w", err)
			}
		}
	} else {
		// Update report status to under_review for next level
		if err := w.reportRepo.UpdateStatus(ctx, approval.ReportID, reportDomain.ReportStatusUnderReview); err != nil {
			return fmt.Errorf("failed to update report status to under review: %w", err)
		}
	}

	return nil
}

// EscalateApproval escalates an approval to the next level automatically
func (w *workflowEngine) EscalateApproval(ctx context.Context, approvalID uuid.UUID) error {
	// Get the approval
	approval, err := w.reportRepo.GetApprovalByID(ctx, approvalID)
	if err != nil {
		return fmt.Errorf("failed to get approval: %w", err)
	}
	if approval == nil {
		return ErrApprovalNotFound
	}

	// Only escalate pending approvals
	if approval.Status != reportDomain.ApprovalStatusPending {
		return errors.New("can only escalate pending approvals")
	}

	// Get the report to find the policy
	report, err := w.reportRepo.GetByID(ctx, approval.ReportID)
	if err != nil {
		return fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return ErrReportNotFound
	}

	// Get next approver for the next level
	nextApproverID, err := w.GetNextApprover(ctx, report.ID, approval.Level)
	if err != nil {
		return fmt.Errorf("failed to get next approver: %w", err)
	}

	// If no next approver found, keep current approval pending
	if nextApproverID == nil {
		return errors.New("no next approver found for escalation")
	}

	// Update current approval to escalated
	now := time.Now()
	approval.Status = reportDomain.ApprovalStatusEscalated
	approval.EscalationDate = &now
	approval.EscalatedTo = nextApproverID
	approval.Updated = now

	if err := w.reportRepo.UpdateApproval(ctx, approval); err != nil {
		return fmt.Errorf("failed to update approval: %w", err)
	}

	// Create history for escalation
	prevStatus := reportDomain.ApprovalStatusPending
	newStatus := reportDomain.ApprovalStatusEscalated
	history := &reportDomain.ApprovalHistory{
		ID:             uuid.New(),
		ApprovalID:     approvalID,
		ReportID:       approval.ReportID,
		ActorID:        uuid.Nil, // System action
		Action:         reportDomain.ApprovalActionEscalated,
		PreviousStatus: &prevStatus,
		NewStatus:      &newStatus,
		Created:        now,
		Metadata: map[string]any{
			"escalated_to": nextApproverID.String(),
			"reason":       "automatic escalation",
		},
	}

	if err := w.reportRepo.CreateApprovalHistory(ctx, history); err != nil {
		return fmt.Errorf("failed to create escalation history: %w", err)
	}

	// Create new approval for escalated approver
	newApproval := &reportDomain.Approval{
		ID:         uuid.New(),
		ReportID:   approval.ReportID,
		ApproverID: *nextApproverID,
		Level:      approval.Level + 1, // Increment level
		Status:     reportDomain.ApprovalStatusPending,
		Created:    now,
		Updated:    now,
	}

	if err := w.reportRepo.CreateApproval(ctx, newApproval); err != nil {
		return fmt.Errorf("failed to create escalated approval: %w", err)
	}

	// Create history for new approval
	createdStatus := reportDomain.ApprovalStatusPending
	newHistory := &reportDomain.ApprovalHistory{
		ID:         uuid.New(),
		ApprovalID: newApproval.ID,
		ReportID:   approval.ReportID,
		ActorID:    uuid.Nil, // System action
		Action:     reportDomain.ApprovalActionCreated,
		NewStatus:  &createdStatus,
		Created:    now,
		Metadata: map[string]any{
			"escalated_from": approvalID.String(),
			"level":          newApproval.Level,
		},
	}

	if err := w.reportRepo.CreateApprovalHistory(ctx, newHistory); err != nil {
		return fmt.Errorf("failed to create new approval history: %w", err)
	}

	return nil
}

// GetNextApprover determines the next approver based on policy and current level
func (w *workflowEngine) GetNextApprover(ctx context.Context, reportID uuid.UUID, currentLevel int) (*uuid.UUID, error) {
	// Get the report
	report, err := w.reportRepo.GetByID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return nil, ErrReportNotFound
	}

	// Verify report has a policy
	if report.PolicyID == nil {
		return nil, ErrNoPolicyAssigned
	}

	// Get the policy
	policy, err := w.policyService.GetByID(ctx, *report.PolicyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy: %w", err)
	}
	if policy == nil {
		return nil, ErrPolicyNotFound
	}

	// Create expense validation input
	expenseInput := &domain.ExpenseValidationInput{
		Amount:   report.TotalAmount,
		Currency: report.Currency,
		UserID:   report.UserID,
	}

	// Get all approvers from policy
	approvers, err := w.ruleEngine.GetApprovers(ctx, expenseInput, policy)
	if err != nil {
		return nil, fmt.Errorf("failed to get approvers: %w", err)
	}

	// Find approver for next level
	for _, approver := range approvers {
		if approver.Level == currentLevel+1 {
			return &approver.UserID, nil
		}
	}

	// No next level approver found
	return nil, nil
}

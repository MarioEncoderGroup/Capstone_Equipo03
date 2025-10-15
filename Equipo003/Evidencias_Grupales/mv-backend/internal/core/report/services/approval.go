package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/report/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/report/ports"
)

var (
	ErrApprovalNotPending      = errors.New("approval is not in pending status")
	ErrNotAuthorizedApprover   = errors.New("user is not the assigned approver")
	ErrAllLevelsAlreadyApproved = errors.New("all approval levels are already approved")
)

type approvalService struct {
	repo ports.ReportRepository
}

// NewApprovalService creates a new instance of the approval service
func NewApprovalService(repo ports.ReportRepository) ports.ApprovalService {
	return &approvalService{repo: repo}
}

// GetPendingApprovals retrieves all pending approvals for a specific approver
func (s *approvalService) GetPendingApprovals(ctx context.Context, approverID uuid.UUID) ([]domain.Approval, error) {
	approvals, err := s.repo.GetPendingApprovalsByApprover(ctx, approverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending approvals: %w", err)
	}

	return approvals, nil
}

// GetPendingApprovalsByReport retrieves all pending approvals for a specific report
func (s *approvalService) GetPendingApprovalsByReport(ctx context.Context, reportID uuid.UUID) ([]domain.Approval, error) {
	approvals, err := s.repo.GetPendingApprovalsByReportID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending approvals by report: %w", err)
	}

	return approvals, nil
}

// Approve approves a pending approval request
func (s *approvalService) Approve(ctx context.Context, approvalID, approverID uuid.UUID, dto *domain.ApproveReportDto) error {
	// Get approval
	approval, err := s.repo.GetApprovalByID(ctx, approvalID)
	if err != nil {
		return fmt.Errorf("failed to get approval: %w", err)
	}
	if approval == nil {
		return ErrApprovalNotFound
	}

	// Verify user is the assigned approver
	if approval.ApproverID != approverID {
		return ErrNotAuthorizedApprover
	}

	// Verify approval is pending
	if approval.Status != domain.ApprovalStatusPending {
		return ErrApprovalNotPending
	}

	// Update approval status
	now := time.Now()
	approval.Status = domain.ApprovalStatusApproved
	approval.DecisionDate = &now
	approval.Updated = now

	if dto.Comments != "" {
		approval.Comments = &dto.Comments
	}
	if dto.ApprovedAmount != nil {
		approval.ApprovedAmount = dto.ApprovedAmount
	}

	if err := s.repo.UpdateApproval(ctx, approval); err != nil {
		return fmt.Errorf("failed to update approval: %w", err)
	}

	// Create history entry
	prevStatus := domain.ApprovalStatusPending
	newStatus := domain.ApprovalStatusApproved
	history := &domain.ApprovalHistory{
		ID:             uuid.New(),
		ApprovalID:     approvalID,
		ReportID:       approval.ReportID,
		ActorID:        approverID,
		Action:         domain.ApprovalActionApproved,
		PreviousStatus: &prevStatus,
		NewStatus:      &newStatus,
		Comments:       approval.Comments,
		Created:        now,
	}

	if err := s.repo.CreateApprovalHistory(ctx, history); err != nil {
		return fmt.Errorf("failed to create approval history: %w", err)
	}

	// Check if all approvals are completed
	if err := s.checkAndUpdateReportStatus(ctx, approval.ReportID); err != nil {
		return fmt.Errorf("failed to check report status: %w", err)
	}

	return nil
}

// Reject rejects a pending approval request
func (s *approvalService) Reject(ctx context.Context, approvalID, approverID uuid.UUID, dto *domain.RejectReportDto) error {
	// Get approval
	approval, err := s.repo.GetApprovalByID(ctx, approvalID)
	if err != nil {
		return fmt.Errorf("failed to get approval: %w", err)
	}
	if approval == nil {
		return ErrApprovalNotFound
	}

	// Verify user is the assigned approver
	if approval.ApproverID != approverID {
		return ErrNotAuthorizedApprover
	}

	// Verify approval is pending
	if approval.Status != domain.ApprovalStatusPending {
		return ErrApprovalNotPending
	}

	// Update approval status
	now := time.Now()
	approval.Status = domain.ApprovalStatusRejected
	approval.DecisionDate = &now
	approval.Comments = &dto.Reason
	approval.Updated = now

	if err := s.repo.UpdateApproval(ctx, approval); err != nil {
		return fmt.Errorf("failed to update approval: %w", err)
	}

	// Create history entry
	prevStatus := domain.ApprovalStatusPending
	newStatus := domain.ApprovalStatusRejected
	history := &domain.ApprovalHistory{
		ID:             uuid.New(),
		ApprovalID:     approvalID,
		ReportID:       approval.ReportID,
		ActorID:        approverID,
		Action:         domain.ApprovalActionRejected,
		PreviousStatus: &prevStatus,
		NewStatus:      &newStatus,
		Comments:       &dto.Reason,
		Created:        now,
	}

	if err := s.repo.CreateApprovalHistory(ctx, history); err != nil {
		return fmt.Errorf("failed to create approval history: %w", err)
	}

	// Update report status to rejected
	if err := s.repo.UpdateStatus(ctx, approval.ReportID, domain.ReportStatusRejected); err != nil {
		return fmt.Errorf("failed to update report status: %w", err)
	}

	// Update report rejection reason
	report, err := s.repo.GetByID(ctx, approval.ReportID)
	if err != nil {
		return fmt.Errorf("failed to get report: %w", err)
	}
	if report != nil {
		report.RejectionReason = &dto.Reason
		report.Updated = now
		if err := s.repo.Update(ctx, report); err != nil {
			return fmt.Errorf("failed to update report: %w", err)
		}
	}

	return nil
}

// GetHistory retrieves the complete approval history for a report
func (s *approvalService) GetHistory(ctx context.Context, reportID uuid.UUID) ([]domain.ApprovalHistory, error) {
	// Get all approvals for the report
	approvals, err := s.repo.GetApprovalsByReport(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get approvals: %w", err)
	}

	// Collect all history entries
	allHistory := []domain.ApprovalHistory{}
	for _, approval := range approvals {
		history, err := s.repo.GetApprovalHistory(ctx, approval.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get approval history: %w", err)
		}
		allHistory = append(allHistory, history...)
	}

	// Sort by created date (newest first)
	// Note: This could be optimized with a database query if needed
	return allHistory, nil
}

// GetApprovalHistory retrieves the history for a specific approval
func (s *approvalService) GetApprovalHistory(ctx context.Context, approvalID uuid.UUID) ([]domain.ApprovalHistory, error) {
	history, err := s.repo.GetApprovalHistory(ctx, approvalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get approval history: %w", err)
	}

	return history, nil
}

// Escalate escalates an approval to another approver
func (s *approvalService) Escalate(ctx context.Context, approvalID, currentApproverID, newApproverID uuid.UUID, reason string) error {
	// Get approval
	approval, err := s.repo.GetApprovalByID(ctx, approvalID)
	if err != nil {
		return fmt.Errorf("failed to get approval: %w", err)
	}
	if approval == nil {
		return ErrApprovalNotFound
	}

	// Verify user is the current approver
	if approval.ApproverID != currentApproverID {
		return ErrNotAuthorizedApprover
	}

	// Verify approval is pending
	if approval.Status != domain.ApprovalStatusPending {
		return ErrApprovalNotPending
	}

	// Update approval with new approver
	now := time.Now()
	approval.Status = domain.ApprovalStatusEscalated
	approval.EscalationDate = &now
	approval.EscalatedTo = &newApproverID
	approval.Updated = now

	if err := s.repo.UpdateApproval(ctx, approval); err != nil {
		return fmt.Errorf("failed to update approval: %w", err)
	}

	// Create history entry
	prevStatus := domain.ApprovalStatusPending
	newStatus := domain.ApprovalStatusEscalated
	reasonPtr := &reason
	metadata := map[string]any{
		"previous_approver": currentApproverID.String(),
		"new_approver":      newApproverID.String(),
		"reason":            reason,
	}

	history := &domain.ApprovalHistory{
		ID:             uuid.New(),
		ApprovalID:     approvalID,
		ReportID:       approval.ReportID,
		ActorID:        currentApproverID,
		Action:         domain.ApprovalActionEscalated,
		PreviousStatus: &prevStatus,
		NewStatus:      &newStatus,
		Comments:       reasonPtr,
		Metadata:       metadata,
		Created:        now,
	}

	if err := s.repo.CreateApprovalHistory(ctx, history); err != nil {
		return fmt.Errorf("failed to create approval history: %w", err)
	}

	// Create new approval for the escalated approver
	newApproval := &domain.Approval{
		ID:         uuid.New(),
		ReportID:   approval.ReportID,
		ApproverID: newApproverID,
		Level:      approval.Level,
		Status:     domain.ApprovalStatusPending,
		Created:    now,
		Updated:    now,
	}

	if err := s.repo.CreateApproval(ctx, newApproval); err != nil {
		return fmt.Errorf("failed to create escalated approval: %w", err)
	}

	// Create history for new approval
	createdStatus := domain.ApprovalStatusPending
	createHistory := &domain.ApprovalHistory{
		ID:         uuid.New(),
		ApprovalID: newApproval.ID,
		ReportID:   approval.ReportID,
		ActorID:    currentApproverID,
		Action:     domain.ApprovalActionCreated,
		NewStatus:  &createdStatus,
		Comments:   reasonPtr,
		Metadata: map[string]any{
			"escalated_from": approvalID.String(),
			"reason":         reason,
		},
		Created: now,
	}

	if err := s.repo.CreateApprovalHistory(ctx, createHistory); err != nil {
		return fmt.Errorf("failed to create new approval history: %w", err)
	}

	return nil
}

// GetApprovalsByReport retrieves all approvals for a report
func (s *approvalService) GetApprovalsByReport(ctx context.Context, reportID uuid.UUID) ([]domain.Approval, error) {
	approvals, err := s.repo.GetApprovalsByReport(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get approvals by report: %w", err)
	}

	return approvals, nil
}

// checkAndUpdateReportStatus checks if all approvals are completed and updates report status
func (s *approvalService) checkAndUpdateReportStatus(ctx context.Context, reportID uuid.UUID) error {
	// Get all approvals for the report
	approvals, err := s.repo.GetApprovalsByReport(ctx, reportID)
	if err != nil {
		return fmt.Errorf("failed to get approvals: %w", err)
	}

	if len(approvals) == 0 {
		return nil
	}

	// Check if all non-escalated approvals are approved
	allApproved := true
	for _, approval := range approvals {
		if approval.Status == domain.ApprovalStatusPending {
			allApproved = false
			break
		}
		if approval.Status != domain.ApprovalStatusApproved && approval.Status != domain.ApprovalStatusEscalated {
			// If any approval is rejected or in another terminal state, don't change status
			return nil
		}
	}

	// If all approvals are approved, update report status
	if allApproved {
		if err := s.repo.UpdateStatus(ctx, reportID, domain.ReportStatusApproved); err != nil {
			return fmt.Errorf("failed to update report status to approved: %w", err)
		}

		// Update approval date
		report, err := s.repo.GetByID(ctx, reportID)
		if err != nil {
			return fmt.Errorf("failed to get report: %w", err)
		}
		if report != nil {
			now := time.Now()
			report.ApprovalDate = &now
			report.Updated = now
			if err := s.repo.Update(ctx, report); err != nil {
				return fmt.Errorf("failed to update report approval date: %w", err)
			}
		}
	}

	return nil
}

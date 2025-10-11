package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/JoseLuis21/mv-backend/internal/core/report/ports"
)

const (
	// StaleApprovalThresholdHours defines how many hours before an approval is considered stale
	StaleApprovalThresholdHours = 24

	// EscalationJobInterval defines how often the escalation job runs
	EscalationJobInterval = 1 * time.Hour
)

// ApprovalEscalationJob handles automatic escalation of stale approvals
type ApprovalEscalationJob struct {
	workflowEngine ports.WorkflowEngine
	reportRepo     ports.ReportRepository
	stopChan       chan struct{}
}

// NewApprovalEscalationJob creates a new instance of the approval escalation job
func NewApprovalEscalationJob(
	workflowEngine ports.WorkflowEngine,
	reportRepo ports.ReportRepository,
) *ApprovalEscalationJob {
	return &ApprovalEscalationJob{
		workflowEngine: workflowEngine,
		reportRepo:     reportRepo,
		stopChan:       make(chan struct{}),
	}
}

// Start begins the escalation job with the configured interval
func (j *ApprovalEscalationJob) Start() {
	log.Println("[ApprovalEscalationJob] Starting approval escalation job")

	// Run immediately on start
	go j.run()

	// Then run on interval
	ticker := time.NewTicker(EscalationJobInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				j.run()
			case <-j.stopChan:
				ticker.Stop()
				log.Println("[ApprovalEscalationJob] Stopping approval escalation job")
				return
			}
		}
	}()
}

// Stop stops the escalation job
func (j *ApprovalEscalationJob) Stop() {
	close(j.stopChan)
}

// run executes the escalation logic
func (j *ApprovalEscalationJob) run() {
	ctx := context.Background()

	log.Printf("[ApprovalEscalationJob] Running escalation check for approvals older than %d hours", StaleApprovalThresholdHours)

	// Get stale approvals
	staleApprovals, err := j.reportRepo.GetStaleApprovals(ctx, StaleApprovalThresholdHours)
	if err != nil {
		log.Printf("[ApprovalEscalationJob] Error fetching stale approvals: %v", err)
		return
	}

	if len(staleApprovals) == 0 {
		log.Println("[ApprovalEscalationJob] No stale approvals found")
		return
	}

	log.Printf("[ApprovalEscalationJob] Found %d stale approvals to escalate", len(staleApprovals))

	// Process each stale approval
	successCount := 0
	errorCount := 0

	for _, approval := range staleApprovals {
		if err := j.escalateApproval(ctx, approval.ID); err != nil {
			log.Printf("[ApprovalEscalationJob] Error escalating approval %s: %v", approval.ID.String(), err)
			errorCount++
		} else {
			log.Printf("[ApprovalEscalationJob] Successfully escalated approval %s (Report: %s, Level: %d)",
				approval.ID.String(), approval.ReportID.String(), approval.Level)
			successCount++
		}
	}

	log.Printf("[ApprovalEscalationJob] Escalation complete. Success: %d, Errors: %d", successCount, errorCount)
}

// escalateApproval escalates a single approval
func (j *ApprovalEscalationJob) escalateApproval(ctx context.Context, approvalID uuid.UUID) error {
	// Use workflow engine to handle escalation
	if err := j.workflowEngine.EscalateApproval(ctx, approvalID); err != nil {
		return fmt.Errorf("failed to escalate approval: %w", err)
	}

	// TODO: Send notification to new approver
	// This would integrate with a notification service when implemented
	// notificationService.SendApprovalEscalatedNotification(...)

	return nil
}

package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/JoseLuis21/mv-backend/internal/core/report/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/report/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
)

type reportRepository struct {
	client *postgresql.PostgresqlClient
}

// NewReportRepository creates a new PostgreSQL report repository
func NewReportRepository(client *postgresql.PostgresqlClient) ports.ReportRepository {
	return &reportRepository{
		client: client,
	}
}

// Create creates a new expense report
func (r *reportRepository) Create(ctx context.Context, report *domain.ExpenseReport) error {
	query := `
		INSERT INTO expense_reports (
			id, user_id, policy_id, title, description, status,
			total_amount, currency, created, updated
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	err := r.client.Exec(ctx, query,
		report.ID,
		report.UserID,
		report.PolicyID,
		report.Title,
		report.Description,
		report.Status,
		report.TotalAmount,
		report.Currency,
		report.Created,
		report.Updated,
	)

	if err != nil {
		return fmt.Errorf("failed to create report: %w", err)
	}

	return nil
}

// GetByID retrieves a report by ID
func (r *reportRepository) GetByID(ctx context.Context, reportID uuid.UUID) (*domain.ExpenseReport, error) {
	query := `
		SELECT
			id, user_id, policy_id, title, description, status,
			total_amount, currency, submission_date, approval_date,
			payment_date, rejection_reason, created, updated, deleted_at
		FROM expense_reports
		WHERE id = $1 AND deleted_at IS NULL
	`

	var report domain.ExpenseReport
	err := r.client.QueryRow(ctx, query, reportID).Scan(
		&report.ID,
		&report.UserID,
		&report.PolicyID,
		&report.Title,
		&report.Description,
		&report.Status,
		&report.TotalAmount,
		&report.Currency,
		&report.SubmissionDate,
		&report.ApprovalDate,
		&report.PaymentDate,
		&report.RejectionReason,
		&report.Created,
		&report.Updated,
		&report.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get report: %w", err)
	}

	return &report, nil
}

// GetByUser retrieves reports for a user with filters
func (r *reportRepository) GetByUser(ctx context.Context, filters *domain.ReportFilters) ([]domain.ExpenseReport, int, error) {
	baseQuery := `
		FROM expense_reports
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argPos := 1

	if filters.UserID != nil {
		baseQuery += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, *filters.UserID)
		argPos++
	}

	if filters.Status != nil {
		baseQuery += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, *filters.Status)
		argPos++
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := r.client.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reports: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id, user_id, policy_id, title, description, status,
			total_amount, currency, submission_date, approval_date,
			payment_date, rejection_reason, created, updated, deleted_at
	` + baseQuery + `
		ORDER BY created DESC
		LIMIT $` + fmt.Sprintf("%d", argPos) + ` OFFSET $` + fmt.Sprintf("%d", argPos+1)

	args = append(args, filters.Limit, filters.Offset)

	rows, err := r.client.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query reports: %w", err)
	}
	defer rows.Close()

	reports := []domain.ExpenseReport{}
	for rows.Next() {
		var report domain.ExpenseReport
		err := rows.Scan(
			&report.ID,
			&report.UserID,
			&report.PolicyID,
			&report.Title,
			&report.Description,
			&report.Status,
			&report.TotalAmount,
			&report.Currency,
			&report.SubmissionDate,
			&report.ApprovalDate,
			&report.PaymentDate,
			&report.RejectionReason,
			&report.Created,
			&report.Updated,
			&report.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan report: %w", err)
		}
		reports = append(reports, report)
	}

	return reports, total, nil
}

// Update updates an existing report
func (r *reportRepository) Update(ctx context.Context, report *domain.ExpenseReport) error {
	query := `
		UPDATE expense_reports
		SET title = $1, description = $2, total_amount = $3,
		    submission_date = $4, approval_date = $5, payment_date = $6,
		    rejection_reason = $7, updated = $8
		WHERE id = $9 AND deleted_at IS NULL
	`

	err := r.client.Exec(ctx, query,
		report.Title,
		report.Description,
		report.TotalAmount,
		report.SubmissionDate,
		report.ApprovalDate,
		report.PaymentDate,
		report.RejectionReason,
		time.Now(),
		report.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update report: %w", err)
	}

	return nil
}

// Delete soft deletes a report
func (r *reportRepository) Delete(ctx context.Context, reportID uuid.UUID) error {
	query := `
		UPDATE expense_reports
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	err := r.client.Exec(ctx, query, time.Now(), reportID)
	if err != nil {
		return fmt.Errorf("failed to delete report: %w", err)
	}

	return nil
}

// UpdateStatus updates only the status of a report
func (r *reportRepository) UpdateStatus(ctx context.Context, reportID uuid.UUID, status domain.ReportStatus) error {
	query := `
		UPDATE expense_reports
		SET status = $1, updated = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	err := r.client.Exec(ctx, query, status, time.Now(), reportID)
	if err != nil {
		return fmt.Errorf("failed to update report status: %w", err)
	}

	return nil
}

// RecalculateTotal recalculates the total amount from all expenses in the report
func (r *reportRepository) RecalculateTotal(ctx context.Context, reportID uuid.UUID) (float64, error) {
	query := `
		SELECT COALESCE(SUM(e.amount), 0)
		FROM expense_report_items eri
		JOIN expenses e ON e.id = eri.expense_id
		WHERE eri.report_id = $1
	`

	var total float64
	err := r.client.QueryRow(ctx, query, reportID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to recalculate total: %w", err)
	}

	return total, nil
}

// AddExpenseToReport adds an expense to a report
func (r *reportRepository) AddExpenseToReport(ctx context.Context, reportID, expenseID uuid.UUID) error {
	// Get next sequence number
	var nextSeq int
	seqQuery := `
		SELECT COALESCE(MAX(sequence_number), 0) + 1
		FROM expense_report_items
		WHERE report_id = $1
	`
	err := r.client.QueryRow(ctx, seqQuery, reportID).Scan(&nextSeq)
	if err != nil {
		return fmt.Errorf("failed to get next sequence: %w", err)
	}

	// Insert item
	query := `
		INSERT INTO expense_report_items (id, report_id, expense_id, sequence_number, created)
		VALUES ($1, $2, $3, $4, $5)
	`

	err = r.client.Exec(ctx, query,
		uuid.New(),
		reportID,
		expenseID,
		nextSeq,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to add expense to report: %w", err)
	}

	return nil
}

// RemoveExpenseFromReport removes an expense from a report
func (r *reportRepository) RemoveExpenseFromReport(ctx context.Context, reportID, expenseID uuid.UUID) error {
	query := `
		DELETE FROM expense_report_items
		WHERE report_id = $1 AND expense_id = $2
	`

	err := r.client.Exec(ctx, query, reportID, expenseID)
	if err != nil {
		return fmt.Errorf("failed to remove expense from report: %w", err)
	}

	return nil
}

// GetReportExpenses retrieves all expenses in a report
func (r *reportRepository) GetReportExpenses(ctx context.Context, reportID uuid.UUID) ([]domain.ExpenseReportItem, error) {
	query := `
		SELECT id, report_id, expense_id, sequence_number, created
		FROM expense_report_items
		WHERE report_id = $1
		ORDER BY sequence_number
	`

	rows, err := r.client.Query(ctx, query, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to query report expenses: %w", err)
	}
	defer rows.Close()

	items := []domain.ExpenseReportItem{}
	for rows.Next() {
		var item domain.ExpenseReportItem
		err := rows.Scan(
			&item.ID,
			&item.ReportID,
			&item.ExpenseID,
			&item.SequenceNumber,
			&item.Created,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan report item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// IsExpenseInReport checks if an expense is already in a report
func (r *reportRepository) IsExpenseInReport(ctx context.Context, expenseID uuid.UUID) (bool, *uuid.UUID, error) {
	query := `
		SELECT report_id
		FROM expense_report_items
		WHERE expense_id = $1
		LIMIT 1
	`

	var reportID uuid.UUID
	err := r.client.QueryRow(ctx, query, expenseID).Scan(&reportID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil, nil
		}
		return false, nil, fmt.Errorf("failed to check expense in report: %w", err)
	}

	return true, &reportID, nil
}

// CreateApproval creates a new approval
func (r *reportRepository) CreateApproval(ctx context.Context, approval *domain.Approval) error {
	query := `
		INSERT INTO approvals (
			id, report_id, approver_id, level, status,
			comments, approved_amount, decision_date, escalation_date,
			escalated_to, created, updated
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	err := r.client.Exec(ctx, query,
		approval.ID,
		approval.ReportID,
		approval.ApproverID,
		approval.Level,
		approval.Status,
		approval.Comments,
		approval.ApprovedAmount,
		approval.DecisionDate,
		approval.EscalationDate,
		approval.EscalatedTo,
		approval.Created,
		approval.Updated,
	)

	if err != nil {
		return fmt.Errorf("failed to create approval: %w", err)
	}

	return nil
}

// GetApprovalsByReport retrieves all approvals for a report
func (r *reportRepository) GetApprovalsByReport(ctx context.Context, reportID uuid.UUID) ([]domain.Approval, error) {
	query := `
		SELECT
			id, report_id, approver_id, level, status,
			comments, approved_amount, decision_date, escalation_date,
			escalated_to, created, updated
		FROM approvals
		WHERE report_id = $1
		ORDER BY level
	`

	rows, err := r.client.Query(ctx, query, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to query approvals: %w", err)
	}
	defer rows.Close()

	approvals := []domain.Approval{}
	for rows.Next() {
		var approval domain.Approval
		err := rows.Scan(
			&approval.ID,
			&approval.ReportID,
			&approval.ApproverID,
			&approval.Level,
			&approval.Status,
			&approval.Comments,
			&approval.ApprovedAmount,
			&approval.DecisionDate,
			&approval.EscalationDate,
			&approval.EscalatedTo,
			&approval.Created,
			&approval.Updated,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan approval: %w", err)
		}
		approvals = append(approvals, approval)
	}

	return approvals, nil
}

// GetApprovalByID retrieves an approval by ID
func (r *reportRepository) GetApprovalByID(ctx context.Context, approvalID uuid.UUID) (*domain.Approval, error) {
	query := `
		SELECT
			id, report_id, approver_id, level, status,
			comments, approved_amount, decision_date, escalation_date,
			escalated_to, created, updated
		FROM approvals
		WHERE id = $1
	`

	var approval domain.Approval
	err := r.client.QueryRow(ctx, query, approvalID).Scan(
		&approval.ID,
		&approval.ReportID,
		&approval.ApproverID,
		&approval.Level,
		&approval.Status,
		&approval.Comments,
		&approval.ApprovedAmount,
		&approval.DecisionDate,
		&approval.EscalationDate,
		&approval.EscalatedTo,
		&approval.Created,
		&approval.Updated,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get approval: %w", err)
	}

	return &approval, nil
}

// UpdateApproval updates an existing approval
func (r *reportRepository) UpdateApproval(ctx context.Context, approval *domain.Approval) error {
	query := `
		UPDATE approvals
		SET status = $1, comments = $2, approved_amount = $3,
		    decision_date = $4, escalation_date = $5, escalated_to = $6,
		    updated = $7
		WHERE id = $8
	`

	err := r.client.Exec(ctx, query,
		approval.Status,
		approval.Comments,
		approval.ApprovedAmount,
		approval.DecisionDate,
		approval.EscalationDate,
		approval.EscalatedTo,
		time.Now(),
		approval.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update approval: %w", err)
	}

	return nil
}

// GetPendingApprovalForUser retrieves a pending approval for a specific user and report
func (r *reportRepository) GetPendingApprovalForUser(ctx context.Context, reportID, userID uuid.UUID) (*domain.Approval, error) {
	query := `
		SELECT
			id, report_id, approver_id, level, status,
			comments, approved_amount, decision_date, escalation_date,
			escalated_to, created, updated
		FROM approvals
		WHERE report_id = $1 AND approver_id = $2 AND status = 'pending'
		LIMIT 1
	`

	var approval domain.Approval
	err := r.client.QueryRow(ctx, query, reportID, userID).Scan(
		&approval.ID,
		&approval.ReportID,
		&approval.ApproverID,
		&approval.Level,
		&approval.Status,
		&approval.Comments,
		&approval.ApprovedAmount,
		&approval.DecisionDate,
		&approval.EscalationDate,
		&approval.EscalatedTo,
		&approval.Created,
		&approval.Updated,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get pending approval: %w", err)
	}

	return &approval, nil
}

// GetPendingApprovalsByApprover retrieves all pending approvals for an approver
func (r *reportRepository) GetPendingApprovalsByApprover(ctx context.Context, approverID uuid.UUID) ([]domain.Approval, error) {
	query := `
		SELECT
			id, report_id, approver_id, level, status,
			comments, approved_amount, decision_date, escalation_date,
			escalated_to, created, updated
		FROM approvals
		WHERE approver_id = $1 AND status = 'pending'
		ORDER BY created ASC
	`

	rows, err := r.client.Query(ctx, query, approverID)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending approvals: %w", err)
	}
	defer rows.Close()

	approvals := []domain.Approval{}
	for rows.Next() {
		var approval domain.Approval
		err := rows.Scan(
			&approval.ID,
			&approval.ReportID,
			&approval.ApproverID,
			&approval.Level,
			&approval.Status,
			&approval.Comments,
			&approval.ApprovedAmount,
			&approval.DecisionDate,
			&approval.EscalationDate,
			&approval.EscalatedTo,
			&approval.Created,
			&approval.Updated,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan approval: %w", err)
		}
		approvals = append(approvals, approval)
	}

	return approvals, nil
}

// GetPendingApprovalsByReportID retrieves all pending approvals for a specific report
func (r *reportRepository) GetPendingApprovalsByReportID(ctx context.Context, reportID uuid.UUID) ([]domain.Approval, error) {
	query := `
		SELECT
			id, report_id, approver_id, level, status,
			comments, approved_amount, decision_date, escalation_date,
			escalated_to, created, updated
		FROM approvals
		WHERE report_id = $1 AND status = 'pending'
		ORDER BY level ASC
	`

	rows, err := r.client.Query(ctx, query, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending approvals: %w", err)
	}
	defer rows.Close()

	approvals := []domain.Approval{}
	for rows.Next() {
		var approval domain.Approval
		err := rows.Scan(
			&approval.ID,
			&approval.ReportID,
			&approval.ApproverID,
			&approval.Level,
			&approval.Status,
			&approval.Comments,
			&approval.ApprovedAmount,
			&approval.DecisionDate,
			&approval.EscalationDate,
			&approval.EscalatedTo,
			&approval.Created,
			&approval.Updated,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan approval: %w", err)
		}
		approvals = append(approvals, approval)
	}

	return approvals, nil
}

// CreateApprovalHistory creates a new approval history entry
func (r *reportRepository) CreateApprovalHistory(ctx context.Context, history *domain.ApprovalHistory) error {
	query := `
		INSERT INTO approval_history (
			id, approval_id, report_id, actor_id, action,
			previous_status, new_status, comments, metadata, created
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	var metadataJSON []byte
	var err error
	if history.Metadata != nil {
		metadataJSON, err = json.Marshal(history.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	}

	err = r.client.Exec(ctx, query,
		history.ID,
		history.ApprovalID,
		history.ReportID,
		history.ActorID,
		history.Action,
		history.PreviousStatus,
		history.NewStatus,
		history.Comments,
		metadataJSON,
		history.Created,
	)

	if err != nil {
		return fmt.Errorf("failed to create approval history: %w", err)
	}

	return nil
}

// GetApprovalHistory retrieves the history for an approval
func (r *reportRepository) GetApprovalHistory(ctx context.Context, approvalID uuid.UUID) ([]domain.ApprovalHistory, error) {
	query := `
		SELECT
			id, approval_id, report_id, actor_id, action,
			previous_status, new_status, comments, metadata, created
		FROM approval_history
		WHERE approval_id = $1
		ORDER BY created DESC
	`

	rows, err := r.client.Query(ctx, query, approvalID)
	if err != nil {
		return nil, fmt.Errorf("failed to query approval history: %w", err)
	}
	defer rows.Close()

	histories := []domain.ApprovalHistory{}
	for rows.Next() {
		var history domain.ApprovalHistory
		var metadataJSON []byte

		err := rows.Scan(
			&history.ID,
			&history.ApprovalID,
			&history.ReportID,
			&history.ActorID,
			&history.Action,
			&history.PreviousStatus,
			&history.NewStatus,
			&history.Comments,
			&metadataJSON,
			&history.Created,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan approval history: %w", err)
		}

		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &history.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		histories = append(histories, history)
	}

	return histories, nil
}

// CreateComment creates a new comment
func (r *reportRepository) CreateComment(ctx context.Context, comment *domain.ExpenseComment) error {
	query := `
		INSERT INTO expense_comments (
			id, report_id, expense_id, user_id, comment_type,
			content, parent_id, is_internal, attachments, created, updated
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	var attachmentsJSON []byte
	var err error
	if comment.Attachments != nil {
		attachmentsJSON, err = json.Marshal(comment.Attachments)
		if err != nil {
			return fmt.Errorf("failed to marshal attachments: %w", err)
		}
	}

	err = r.client.Exec(ctx, query,
		comment.ID,
		comment.ReportID,
		comment.ExpenseID,
		comment.UserID,
		comment.CommentType,
		comment.Content,
		comment.ParentID,
		comment.IsInternal,
		attachmentsJSON,
		comment.Created,
		comment.Updated,
	)

	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

// GetCommentsByReport retrieves all comments for a report
func (r *reportRepository) GetCommentsByReport(ctx context.Context, reportID uuid.UUID, includeInternal bool) ([]domain.ExpenseComment, error) {
	query := `
		SELECT
			id, report_id, expense_id, user_id, comment_type,
			content, parent_id, is_internal, attachments, created, updated, deleted_at
		FROM expense_comments
		WHERE report_id = $1 AND deleted_at IS NULL
	`

	if !includeInternal {
		query += " AND is_internal = false"
	}

	query += " ORDER BY created ASC"

	rows, err := r.client.Query(ctx, query, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	comments := []domain.ExpenseComment{}
	for rows.Next() {
		var comment domain.ExpenseComment
		var attachmentsJSON []byte

		err := rows.Scan(
			&comment.ID,
			&comment.ReportID,
			&comment.ExpenseID,
			&comment.UserID,
			&comment.CommentType,
			&comment.Content,
			&comment.ParentID,
			&comment.IsInternal,
			&attachmentsJSON,
			&comment.Created,
			&comment.Updated,
			&comment.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		if len(attachmentsJSON) > 0 {
			if err := json.Unmarshal(attachmentsJSON, &comment.Attachments); err != nil {
				return nil, fmt.Errorf("failed to unmarshal attachments: %w", err)
			}
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// GetCommentsByExpense retrieves all comments for an expense
func (r *reportRepository) GetCommentsByExpense(ctx context.Context, expenseID uuid.UUID) ([]domain.ExpenseComment, error) {
	query := `
		SELECT
			id, report_id, expense_id, user_id, comment_type,
			content, parent_id, is_internal, attachments, created, updated, deleted_at
		FROM expense_comments
		WHERE expense_id = $1 AND deleted_at IS NULL
		ORDER BY created ASC
	`

	rows, err := r.client.Query(ctx, query, expenseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	comments := []domain.ExpenseComment{}
	for rows.Next() {
		var comment domain.ExpenseComment
		var attachmentsJSON []byte

		err := rows.Scan(
			&comment.ID,
			&comment.ReportID,
			&comment.ExpenseID,
			&comment.UserID,
			&comment.CommentType,
			&comment.Content,
			&comment.ParentID,
			&comment.IsInternal,
			&attachmentsJSON,
			&comment.Created,
			&comment.Updated,
			&comment.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		if len(attachmentsJSON) > 0 {
			if err := json.Unmarshal(attachmentsJSON, &comment.Attachments); err != nil {
				return nil, fmt.Errorf("failed to unmarshal attachments: %w", err)
			}
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// UpdateComment updates an existing comment
func (r *reportRepository) UpdateComment(ctx context.Context, comment *domain.ExpenseComment) error {
	query := `
		UPDATE expense_comments
		SET content = $1, updated = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	err := r.client.Exec(ctx, query,
		comment.Content,
		time.Now(),
		comment.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

// DeleteComment soft deletes a comment
func (r *reportRepository) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	query := `
		UPDATE expense_comments
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	err := r.client.Exec(ctx, query, time.Now(), commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

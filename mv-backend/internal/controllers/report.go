package controllers

import (
	"github.com/JoseLuis21/mv-backend/internal/core/report/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/report/ports"
	"github.com/JoseLuis21/mv-backend/internal/shared/pagination"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ReportController handles HTTP operations for expense reports
type ReportController struct {
	reportService ports.ReportService
	validator     *validatorapi.Validator
}

// NewReportController creates a new instance of the report controller
func NewReportController(reportService ports.ReportService, validator *validatorapi.Validator) *ReportController {
	return &ReportController{
		reportService: reportService,
		validator:     validator,
	}
}

// CreateReport handles POST /expense-reports - Create new report
func (rc *ReportController) CreateReport(c *fiber.Ctx) error {
	// Get userID from context (from JWT)
	userIDStr := c.Locals("userID")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
			Success: false,
			Message: "User not authenticated",
			Error:   "UNAUTHORIZED",
		})
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   "INVALID_USER_ID",
		})
	}

	// Parse request body
	var dto domain.CreateReportDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if errors := rc.validator.ValidateStruct(&dto); len(errors) > 0 {
		var validationErrors []types.ValidationErrorResponse
		for _, err := range errors {
			validationErrors = append(validationErrors, types.ValidationErrorResponse{
				Field:   err.Field,
				Message: err.Message,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Validation error",
			Error:   "VALIDATION_ERROR",
			Data:    validationErrors,
		})
	}

	// Create report
	report, err := rc.reportService.CreateReport(c.Context(), userID, &dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error creating report",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusCreated).JSON(types.APIResponse{
		Success: true,
		Message: "Report created successfully",
		Data:    report,
	})
}

// GetUserReports handles GET /expense-reports - Get user's reports
func (rc *ReportController) GetUserReports(c *fiber.Ctx) error {
	// Get userID from context
	userIDStr := c.Locals("userID")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
			Success: false,
			Message: "User not authenticated",
			Error:   "UNAUTHORIZED",
		})
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   "INVALID_USER_ID",
		})
	}

	// Parse pagination parameters
	paginationReq, err := pagination.ParsePaginationFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid pagination parameters",
			Error:   err.Error(),
		})
	}

	// Build filters
	filters := &domain.ReportFilters{
		UserID: &userID,
		Limit:  paginationReq.GetLimit(),
		Offset: paginationReq.GetOffset(),
	}

	// Filter by status
	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.ReportStatus(statusStr)
		filters.Status = &status
	}

	// Get reports
	response, err := rc.reportService.GetUserReports(c.Context(), filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error fetching reports",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Reports retrieved successfully",
		Data:    response,
	})
}

// GetReportByID handles GET /expense-reports/:id - Get report details
func (rc *ReportController) GetReportByID(c *fiber.Ctx) error {
	// Parse report ID
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid report ID",
			Error:   "INVALID_REPORT_ID",
		})
	}

	// Get report
	report, err := rc.reportService.GetReportByID(c.Context(), reportID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error fetching report",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Report retrieved successfully",
		Data:    report,
	})
}

// UpdateReport handles PUT /expense-reports/:id - Update report
func (rc *ReportController) UpdateReport(c *fiber.Ctx) error {
	// Get userID from context
	userIDStr := c.Locals("userID")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
			Success: false,
			Message: "User not authenticated",
			Error:   "UNAUTHORIZED",
		})
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   "INVALID_USER_ID",
		})
	}

	// Parse report ID
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid report ID",
			Error:   "INVALID_REPORT_ID",
		})
	}

	// Check permissions
	canEdit, err := rc.reportService.CanEditReport(c.Context(), reportID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error checking permissions",
		Error:   err.Error(),
	})
	}
	if !canEdit {
		return c.Status(fiber.StatusForbidden).JSON(types.APIResponse{
			Success: false,
			Message: "You don't have permission to edit this report",
			Error:   "FORBIDDEN",
		})
	}

	// Parse request body
	var dto domain.UpdateReportDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if errors := rc.validator.ValidateStruct(&dto); len(errors) > 0 {
		var validationErrors []types.ValidationErrorResponse
		for _, err := range errors {
			validationErrors = append(validationErrors, types.ValidationErrorResponse{
				Field:   err.Field,
				Message: err.Message,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Validation error",
			Error:   "VALIDATION_ERROR",
			Data:    validationErrors,
		})
	}

	// Update report
	if err := rc.reportService.UpdateReport(c.Context(), reportID, &dto); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error updating report",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Report updated successfully",
	})
}

// DeleteReport handles DELETE /expense-reports/:id - Delete report
func (rc *ReportController) DeleteReport(c *fiber.Ctx) error {
	// Get userID from context
	userIDStr := c.Locals("userID")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
			Success: false,
			Message: "User not authenticated",
			Error:   "UNAUTHORIZED",
		})
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   "INVALID_USER_ID",
		})
	}

	// Parse report ID
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid report ID",
			Error:   "INVALID_REPORT_ID",
		})
	}

	// Check permissions
	canDelete, err := rc.reportService.CanDeleteReport(c.Context(), reportID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error checking permissions",
		Error:   err.Error(),
	})
	}
	if !canDelete {
		return c.Status(fiber.StatusForbidden).JSON(types.APIResponse{
			Success: false,
			Message: "You don't have permission to delete this report",
			Error:   "FORBIDDEN",
		})
	}

	// Delete report
	if err := rc.reportService.DeleteReport(c.Context(), reportID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error deleting report",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Report deleted successfully",
	})
}

// SubmitReport handles POST /expense-reports/:id/submit - Submit report for approval
func (rc *ReportController) SubmitReport(c *fiber.Ctx) error {
	// Parse report ID
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid report ID",
			Error:   "INVALID_REPORT_ID",
		})
	}

	// Check if report can be submitted
	canSubmit, err := rc.reportService.CanSubmitReport(c.Context(), reportID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error checking if report can be submitted",
		Error:   err.Error(),
	})
	}
	if !canSubmit {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Report cannot be submitted. It must be in draft status and have at least one expense.",
			Error:   "CANNOT_SUBMIT",
		})
	}

	// Submit report
	response, err := rc.reportService.SubmitReport(c.Context(), reportID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error submitting report",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Report submitted successfully",
		Data:    response,
	})
}

// AddExpensesToReport handles POST /expense-reports/:id/expenses - Add expenses to report
func (rc *ReportController) AddExpensesToReport(c *fiber.Ctx) error {
	// Get userID from context
	userIDStr := c.Locals("userID")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
			Success: false,
			Message: "User not authenticated",
			Error:   "UNAUTHORIZED",
		})
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   "INVALID_USER_ID",
		})
	}

	// Parse report ID
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid report ID",
			Error:   "INVALID_REPORT_ID",
		})
	}

	// Check permissions
	canEdit, err := rc.reportService.CanEditReport(c.Context(), reportID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error checking permissions",
		Error:   err.Error(),
	})
	}
	if !canEdit {
		return c.Status(fiber.StatusForbidden).JSON(types.APIResponse{
			Success: false,
			Message: "You don't have permission to edit this report",
			Error:   "FORBIDDEN",
		})
	}

	// Parse request body
	var dto domain.AddExpensesDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if errors := rc.validator.ValidateStruct(&dto); len(errors) > 0 {
		var validationErrors []types.ValidationErrorResponse
		for _, err := range errors {
			validationErrors = append(validationErrors, types.ValidationErrorResponse{
				Field:   err.Field,
				Message: err.Message,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Validation error",
			Error:   "VALIDATION_ERROR",
			Data:    validationErrors,
		})
	}

	// Add expenses
	if err := rc.reportService.AddExpensesToReport(c.Context(), reportID, &dto); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error adding expenses to report",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Expenses added to report successfully",
	})
}

// RemoveExpenseFromReport handles DELETE /expense-reports/:id/expenses/:expenseId - Remove expense from report
func (rc *ReportController) RemoveExpenseFromReport(c *fiber.Ctx) error {
	// Get userID from context
	userIDStr := c.Locals("userID")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
			Success: false,
			Message: "User not authenticated",
			Error:   "UNAUTHORIZED",
		})
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   "INVALID_USER_ID",
		})
	}

	// Parse report ID
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid report ID",
			Error:   "INVALID_REPORT_ID",
		})
	}

	// Parse expense ID
	expenseID, err := uuid.Parse(c.Params("expenseId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid expense ID",
			Error:   "INVALID_EXPENSE_ID",
		})
	}

	// Check permissions
	canEdit, err := rc.reportService.CanEditReport(c.Context(), reportID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error checking permissions",
		Error:   err.Error(),
	})
	}
	if !canEdit {
		return c.Status(fiber.StatusForbidden).JSON(types.APIResponse{
			Success: false,
			Message: "You don't have permission to edit this report",
			Error:   "FORBIDDEN",
		})
	}

	// Remove expense
	if err := rc.reportService.RemoveExpenseFromReport(c.Context(), reportID, expenseID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error removing expense from report",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Expense removed from report successfully",
	})
}

// AddComment handles POST /expense-reports/:id/comments - Add comment to report
func (rc *ReportController) AddComment(c *fiber.Ctx) error {
	// Get userID from context
	userIDStr := c.Locals("userID")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
			Success: false,
			Message: "User not authenticated",
			Error:   "UNAUTHORIZED",
		})
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   "INVALID_USER_ID",
		})
	}

	// Parse request body
	var dto domain.CreateCommentDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if errors := rc.validator.ValidateStruct(&dto); len(errors) > 0 {
		var validationErrors []types.ValidationErrorResponse
		for _, err := range errors {
			validationErrors = append(validationErrors, types.ValidationErrorResponse{
				Field:   err.Field,
				Message: err.Message,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Validation error",
			Error:   "VALIDATION_ERROR",
			Data:    validationErrors,
		})
	}

	// Add comment
	comment, err := rc.reportService.AddComment(c.Context(), userID, &dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error adding comment",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusCreated).JSON(types.APIResponse{
		Success: true,
		Message: "Comment added successfully",
		Data:    comment,
	})
}

// GetReportComments handles GET /expense-reports/:id/comments - Get report comments
func (rc *ReportController) GetReportComments(c *fiber.Ctx) error {
	// Parse report ID
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid report ID",
			Error:   "INVALID_REPORT_ID",
		})
	}

	// Parse includeInternal query param
	includeInternal := c.Query("include_internal", "false") == "true"

	// Get comments
	comments, err := rc.reportService.GetReportComments(c.Context(), reportID, includeInternal)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error fetching comments",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Comments retrieved successfully",
		Data:    comments,
	})
}

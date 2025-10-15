package controllers

import (
	"github.com/JoseLuis21/mv-backend/internal/core/report/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/report/ports"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ApprovalController handles HTTP operations for approvals
type ApprovalController struct {
	approvalService ports.ApprovalService
	validator       *validatorapi.Validator
}

// NewApprovalController creates a new instance of the approval controller
func NewApprovalController(approvalService ports.ApprovalService, validator *validatorapi.Validator) *ApprovalController {
	return &ApprovalController{
		approvalService: approvalService,
		validator:       validator,
	}
}

// GetPendingApprovals handles GET /approvals/pending - Get user's pending approvals
func (ac *ApprovalController) GetPendingApprovals(c *fiber.Ctx) error {
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

	// Get pending approvals
	approvals, err := ac.approvalService.GetPendingApprovals(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error fetching pending approvals",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Pending approvals retrieved successfully",
		Data:    approvals,
	})
}

// GetApprovalsByReport handles GET /approvals/reports/:id - Get all approvals for a report
func (ac *ApprovalController) GetApprovalsByReport(c *fiber.Ctx) error {
	// Parse report ID
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid report ID",
			Error:   "INVALID_REPORT_ID",
		})
	}

	// Get approvals
	approvals, err := ac.approvalService.GetApprovalsByReport(c.Context(), reportID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error fetching approvals",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Approvals retrieved successfully",
		Data:    approvals,
	})
}

// ApproveReport handles POST /approvals/:id/approve - Approve a report
func (ac *ApprovalController) ApproveReport(c *fiber.Ctx) error {
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

	// Parse approval ID
	approvalID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid approval ID",
			Error:   "INVALID_APPROVAL_ID",
		})
	}

	// Parse request body
	var dto domain.ApproveReportDto
	if err := c.BodyParser(&dto); err != nil {
		// Allow empty body
		dto = domain.ApproveReportDto{}
	}

	// Validate request
	if errors := ac.validator.ValidateStruct(&dto); len(errors) > 0 {
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

	// Approve report
	if err := ac.approvalService.Approve(c.Context(), approvalID, userID, &dto); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error approving report",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Report approved successfully",
	})
}

// RejectReport handles POST /approvals/:id/reject - Reject a report
func (ac *ApprovalController) RejectReport(c *fiber.Ctx) error {
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

	// Parse approval ID
	approvalID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid approval ID",
			Error:   "INVALID_APPROVAL_ID",
		})
	}

	// Parse request body
	var dto domain.RejectReportDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if errors := ac.validator.ValidateStruct(&dto); len(errors) > 0 {
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

	// Reject report
	if err := ac.approvalService.Reject(c.Context(), approvalID, userID, &dto); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error rejecting report",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Report rejected successfully",
	})
}

// GetApprovalHistory handles GET /approvals/reports/:id/history - Get approval history for a report
func (ac *ApprovalController) GetApprovalHistory(c *fiber.Ctx) error {
	// Parse report ID
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid report ID",
			Error:   "INVALID_REPORT_ID",
		})
	}

	// Get history
	history, err := ac.approvalService.GetHistory(c.Context(), reportID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error fetching approval history",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Approval history retrieved successfully",
		Data:    history,
	})
}

// EscalateApproval handles POST /approvals/:id/escalate - Escalate an approval
func (ac *ApprovalController) EscalateApproval(c *fiber.Ctx) error {
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

	// Parse approval ID
	approvalID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid approval ID",
			Error:   "INVALID_APPROVAL_ID",
		})
	}

	// Parse request body
	var requestBody struct {
		NewApproverID string `json:"new_approver_id" validate:"required,uuid"`
		Reason        string `json:"reason" validate:"required,min=10,max=500"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if errors := ac.validator.ValidateStruct(&requestBody); len(errors) > 0 {
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

	// Parse new approver ID
	newApproverID, err := uuid.Parse(requestBody.NewApproverID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Invalid new approver ID",
			Error:   "INVALID_APPROVER_ID",
		})
	}

	// Escalate approval
	if err := ac.approvalService.Escalate(c.Context(), approvalID, userID, newApproverID, requestBody.Reason); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
		Success: false,
		Message: "Error escalating approval",
		Error:   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Approval escalated successfully",
	})
}

package controllers

import (
	"time"

	"github.com/JoseLuis21/mv-backend/internal/core/policy/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/policy/ports"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PolicyController maneja las operaciones HTTP para políticas
type PolicyController struct {
	policyService ports.PolicyService
	ruleEngine    ports.RuleEngine
	validator     *validatorapi.Validator
}

// NewPolicyController crea una nueva instancia del controller de políticas
func NewPolicyController(
	policyService ports.PolicyService,
	ruleEngine ports.RuleEngine,
	validator *validatorapi.Validator,
) *PolicyController {
	return &PolicyController{
		policyService: policyService,
		ruleEngine:    ruleEngine,
		validator:     validator,
	}
}

// ValidateExpense maneja POST /policies/validate - Valida un gasto contra políticas
func (pc *PolicyController) ValidateExpense(c *fiber.Ctx) error {
	// Obtener userID del contexto (del JWT)
	userIDStr := c.Locals("userID")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
			Success: false,
			Message: "Usuario no autenticado",
			Error:   "UNAUTHORIZED",
		})
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "User ID inválido",
			Error:   "INVALID_USER_ID",
		})
	}

	// Parsear request body
	var req domain.ValidateExpenseDto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Request inválido",
			Error:   err.Error(),
		})
	}

	// Validar request
	if errs := pc.validator.ValidateStruct(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Errores de validación",
			Data:    errs,
			Error:   "VALIDATION_ERROR",
		})
	}

	// Parsear fecha
	expenseDate, err := time.Parse("2006-01-02", req.ExpenseDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Formato de fecha inválido, use YYYY-MM-DD",
			Error:   "INVALID_DATE_FORMAT",
		})
	}

	// Validar que la fecha no sea futura
	if expenseDate.After(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "La fecha del gasto no puede ser futura",
			Error:   "FUTURE_DATE_NOT_ALLOWED",
		})
	}

	// Obtener la política activa para el usuario
	// Por ahora asumimos que hay una política activa (en el futuro se podría filtrar por tipo, departamento, etc.)
	policies, _, err := pc.policyService.GetAll(c.Context(), domain.PolicyFilters{
		IsActive: boolPtr(true),
		Limit:    1,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error al obtener políticas",
			Error:   err.Error(),
		})
	}

	if len(policies) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(types.APIResponse{
			Success: false,
			Message: "No hay políticas activas configuradas",
			Error:   "NO_ACTIVE_POLICY",
		})
	}

	policy := &policies[0]

	// Crear input de validación
	currency := req.Currency
	if currency == "" {
		currency = "CLP" // Default currency
	}

	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	expenseInput := &domain.ExpenseValidationInput{
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Currency:    currency,
		ExpenseDate: expenseDate,
		UserID:      userID,
		Description: description,
	}

	// Validar contra la política
	violations, err := pc.ruleEngine.ValidateExpense(c.Context(), expenseInput, policy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error al validar gasto",
			Error:   err.Error(),
		})
	}

	// Verificar si requiere aprobación
	requiresApproval, approvalLevel, err := pc.ruleEngine.CheckApprovalRequired(c.Context(), expenseInput, policy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error al verificar aprobación",
			Error:   err.Error(),
		})
	}

	// Obtener aprobadores si requiere aprobación
	var approvers []domain.ApproverInfo
	if requiresApproval {
		approvers, err = pc.ruleEngine.GetApprovers(c.Context(), expenseInput, policy)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
				Success: false,
				Message: "Error al obtener aprobadores",
				Error:   err.Error(),
			})
		}
	}

	// Determinar si es válido (no tiene violaciones de tipo "error")
	isValid := true
	for _, v := range violations {
		if v.Severity == "error" {
			isValid = false
			break
		}
	}

	// Construir respuesta
	response := domain.ValidationResponse{
		IsValid:          isValid,
		Violations:       violations,
		RequiresApproval: requiresApproval,
		ApprovalLevel:    approvalLevel,
		Approvers:        approvers,
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Validación completada exitosamente",
		Data:    response,
	})
}

// Helper functions

func boolPtr(b bool) *bool {
	return &b
}

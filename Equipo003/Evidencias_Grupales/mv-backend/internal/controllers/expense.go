package controllers

import (
	"time"

	"github.com/JoseLuis21/mv-backend/internal/core/expense/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/expense/ports"
	"github.com/JoseLuis21/mv-backend/internal/shared/pagination"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ExpenseController maneja las operaciones HTTP para gastos
type ExpenseController struct {
	expenseService ports.ExpenseService
	validator      *validatorapi.Validator
}

// NewExpenseController crea una nueva instancia del controller de gastos
func NewExpenseController(expenseService ports.ExpenseService, validator *validatorapi.Validator) *ExpenseController {
	return &ExpenseController{
		expenseService: expenseService,
		validator:      validator,
	}
}

// GetExpenses maneja GET /expenses - Lista gastos del usuario con filtros
func (ec *ExpenseController) GetExpenses(c *fiber.Ctx) error {
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

	// Parsear parámetros de paginación
	paginationReq, err := pagination.ParsePaginationFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Parámetros de paginación inválidos",
			Error:   err.Error(),
		})
	}

	// Construir filtros
	filters := domain.ExpenseFilters{
		UserID: &userID,
		Limit:  paginationReq.GetLimit(),
		Offset: paginationReq.GetOffset(),
	}

	// Filtro por categoría
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		categoryID, err := uuid.Parse(categoryIDStr)
		if err == nil {
			filters.CategoryID = &categoryID
		}
	}

	// Filtro por estado
	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.ExpenseStatus(statusStr)
		if status.IsValid() {
			filters.Status = &status
		}
	}

	// Filtro por rango de fechas
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		dateFrom, err := time.Parse("2006-01-02", dateFromStr)
		if err == nil {
			filters.DateFrom = &dateFrom
		}
	}

	if dateToStr := c.Query("date_to"); dateToStr != "" {
		dateTo, err := time.Parse("2006-01-02", dateToStr)
		if err == nil {
			filters.DateTo = &dateTo
		}
	}

	// Filtro por búsqueda de texto
	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}

	// Obtener gastos desde el servicio
	expenses, total, err := ec.expenseService.GetAll(c.Context(), filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo gastos",
			Error:   err.Error(),
		})
	}

	// Construir respuesta paginada
	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Gastos obtenidos exitosamente",
		Data: fiber.Map{
			"expenses": expenses,
			"pagination": fiber.Map{
				"total":  total,
				"limit":  filters.Limit,
				"offset": filters.Offset,
				"page":   (filters.Offset / filters.Limit) + 1,
				"pages":  (total + filters.Limit - 1) / filters.Limit,
			},
		},
	})
}

// GetExpenseByID maneja GET /expenses/:id - Detalle de un gasto
func (ec *ExpenseController) GetExpenseByID(c *fiber.Ctx) error {
	// Obtener userID del contexto
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

	// Obtener ID del gasto
	expenseIDStr := c.Params("id")
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de gasto inválido",
			Error:   "INVALID_EXPENSE_ID",
		})
	}

	// Obtener gasto
	expense, err := ec.expenseService.GetByID(c.Context(), expenseID, userID)
	if err != nil {
		if err == sharedErrors.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(types.APIResponse{
				Success: false,
				Message: "Gasto no encontrado",
				Error:   "NOT_FOUND",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo gasto",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Gasto obtenido exitosamente",
		Data:    expense,
	})
}

// CreateExpense maneja POST /expenses - Crear nuevo gasto
func (ec *ExpenseController) CreateExpense(c *fiber.Ctx) error {
	// Obtener userID del contexto
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

	// Parsear body
	var dto domain.CreateExpenseDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
	}

	// Validar DTO
	if errors := ec.validator.ValidateStruct(dto); len(errors) > 0 {
		var validationErrors []types.ValidationErrorResponse
		for _, err := range errors {
			validationErrors = append(validationErrors, types.ValidationErrorResponse{
				Field:   err.Field,
				Message: err.Message,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Errores de validación",
			Data:    validationErrors,
		})
	}

	// Crear gasto
	expense, err := ec.expenseService.Create(c.Context(), dto, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error creando gasto",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(types.APIResponse{
		Success: true,
		Message: "Gasto creado exitosamente",
		Data:    expense,
	})
}

// UpdateExpense maneja PUT /expenses/:id - Actualizar gasto
func (ec *ExpenseController) UpdateExpense(c *fiber.Ctx) error {
	// Obtener userID del contexto
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

	// Obtener ID del gasto
	expenseIDStr := c.Params("id")
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de gasto inválido",
			Error:   "INVALID_EXPENSE_ID",
		})
	}

	// Parsear body
	var dto domain.UpdateExpenseDto
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
	}

	// Validar DTO
	if errors := ec.validator.ValidateStruct(dto); len(errors) > 0 {
		var validationErrors []types.ValidationErrorResponse
		for _, err := range errors {
			validationErrors = append(validationErrors, types.ValidationErrorResponse{
				Field:   err.Field,
				Message: err.Message,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Errores de validación",
			Data:    validationErrors,
		})
	}

	// Actualizar gasto
	expense, err := ec.expenseService.Update(c.Context(), expenseID, dto, userID)
	if err != nil {
		if err == sharedErrors.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(types.APIResponse{
				Success: false,
				Message: "Gasto no encontrado",
				Error:   "NOT_FOUND",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error actualizando gasto",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Gasto actualizado exitosamente",
		Data:    expense,
	})
}

// DeleteExpense maneja DELETE /expenses/:id - Eliminar gasto
func (ec *ExpenseController) DeleteExpense(c *fiber.Ctx) error {
	// Obtener userID del contexto
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

	// Obtener ID del gasto
	expenseIDStr := c.Params("id")
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de gasto inválido",
			Error:   "INVALID_EXPENSE_ID",
		})
	}

	// Eliminar gasto
	if err := ec.expenseService.Delete(c.Context(), expenseID, userID); err != nil {
		if err == sharedErrors.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(types.APIResponse{
				Success: false,
				Message: "Gasto no encontrado",
				Error:   "NOT_FOUND",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error eliminando gasto",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Gasto eliminado exitosamente",
		Data:    nil,
	})
}

// UploadReceipt maneja POST /expenses/:id/receipts - Subir comprobante
func (ec *ExpenseController) UploadReceipt(c *fiber.Ctx) error {
	// Obtener userID del contexto
	userIDStr := c.Locals("userID")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
			Success: false,
			Message: "Usuario no autenticado",
			Error:   "UNAUTHORIZED",
		})
	}

	// Obtener ID del gasto
	expenseIDStr := c.Params("id")
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de gasto inválido",
			Error:   "INVALID_EXPENSE_ID",
		})
	}

	// Obtener archivo
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Archivo requerido",
			Error:   "FILE_REQUIRED",
		})
	}

	// Validar tamaño (max 10MB)
	maxSize := int64(10485760) // 10MB
	if file.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Archivo demasiado grande (máx 10MB)",
			Error:   "FILE_TOO_LARGE",
		})
	}

	// Validar tipo de archivo
	allowedTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/jpg":       true,
		"application/pdf": true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Tipo de archivo no permitido (solo JPG, PNG, PDF)",
			Error:   "INVALID_FILE_TYPE",
		})
	}

	// Leer archivo
	fileData, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error leyendo archivo",
			Error:   err.Error(),
		})
	}
	defer fileData.Close()

	// Por ahora no leemos el contenido, solo creamos el DTO
	// TODO: Leer y pasar fileData al service para upload a S3

	// Crear DTO
	dto := domain.UploadReceiptDto{
		ExpenseID: expenseID,
		FileName:  file.Filename,
		FileType:  file.Header.Get("Content-Type"),
		FileSize:  file.Size,
		IsPrimary: c.FormValue("is_primary") == "true",
	}

	// Upload receipt
	receipt, err := ec.expenseService.UploadReceipt(c.Context(), dto, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error subiendo comprobante",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(types.APIResponse{
		Success: true,
		Message: "Comprobante subido exitosamente",
		Data:    receipt,
	})
}

// DeleteReceipt maneja DELETE /receipts/:id - Eliminar comprobante
func (ec *ExpenseController) DeleteReceipt(c *fiber.Ctx) error {
	// Obtener userID del contexto
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

	// Obtener ID del receipt
	receiptIDStr := c.Params("id")
	receiptID, err := uuid.Parse(receiptIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de comprobante inválido",
			Error:   "INVALID_RECEIPT_ID",
		})
	}

	// Eliminar receipt
	if err := ec.expenseService.DeleteReceipt(c.Context(), receiptID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error eliminando comprobante",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Comprobante eliminado exitosamente",
		Data:    nil,
	})
}

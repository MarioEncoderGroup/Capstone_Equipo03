package services

import (
	"context"
	"fmt"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/core/expense/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/expense/ports"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
	"github.com/google/uuid"
)

// expenseService implementa el servicio de gastos
type expenseService struct {
	expenseRepo  ports.ExpenseRepository
	categoryRepo ports.CategoryRepository
}

// NewExpenseService crea una nueva instancia del servicio de gastos
func NewExpenseService(expenseRepo ports.ExpenseRepository, categoryRepo ports.CategoryRepository) ports.ExpenseService {
	return &expenseService{
		expenseRepo:  expenseRepo,
		categoryRepo: categoryRepo,
	}
}

// Create crea un nuevo gasto con validaciones de negocio
func (s *expenseService) Create(ctx context.Context, dto domain.CreateExpenseDto, userID uuid.UUID) (*domain.Expense, error) {
	// Validar que la fecha no sea futura
	if dto.ExpenseDate.After(time.Now()) {
		return nil, sharedErrors.NewValidationError("fecha_invalida", "La fecha del gasto no puede ser futura")
	}

	// Verificar que la categoría existe y está activa
	category, err := s.categoryRepo.GetByID(ctx, dto.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("error verificando categoría: %w", err)
	}
	if !category.IsActive {
		return nil, sharedErrors.NewValidationError("categoria_inactiva", "La categoría seleccionada no está activa")
	}

	// Calcular exchange rate (por ahora siempre 1.0, luego se puede integrar con API de tasas)
	exchangeRate := 1.0
	if dto.Currency == "CLP" {
		exchangeRate = 1.0
	}
	// TODO: Integrar con API de exchange rates para USD, EUR, etc.

	// Crear el gasto
	now := time.Now()
	expense := &domain.Expense{
		ID:             uuid.New(),
		UserID:         userID,
		CategoryID:     dto.CategoryID,
		Title:          dto.Title,
		Amount:         dto.Amount,
		Currency:       dto.Currency,
		ExchangeRate:   exchangeRate,
		AmountCLP:      dto.Amount * exchangeRate,
		ExpenseDate:    dto.ExpenseDate,
		PaymentMethod:  dto.PaymentMethod,
		Status:         domain.ExpenseStatusDraft,
		IsReimbursable: dto.IsReimbursable,
		Created:        now,
		Updated:        now,
	}

	// Campos opcionales
	if dto.Description != "" {
		expense.Description = &dto.Description
	}
	if dto.MerchantName != "" {
		expense.MerchantName = &dto.MerchantName
	}
	if dto.MerchantRUT != "" {
		// TODO: Validar formato de RUT chileno
		expense.MerchantRUT = &dto.MerchantRUT
	}
	if dto.ReceiptNumber != "" {
		expense.ReceiptNumber = &dto.ReceiptNumber
	}

	// Guardar en base de datos
	if err := s.expenseRepo.Create(ctx, expense); err != nil {
		return nil, fmt.Errorf("error creando gasto: %w", err)
	}

	// Obtener el gasto completo con categoría
	return s.expenseRepo.GetByID(ctx, expense.ID)
}

// GetByID obtiene un gasto por ID con validación de ownership
func (s *expenseService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Expense, error) {
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verificar que el usuario sea el owner
	if expense.UserID != userID {
		return nil, sharedErrors.NewBusinessError("no_autorizado", "No tienes permiso para ver este gasto", "")
	}

	return expense, nil
}

// GetAll obtiene gastos con filtros
func (s *expenseService) GetAll(ctx context.Context, filters domain.ExpenseFilters) ([]domain.Expense, int, error) {
	// Establecer límite máximo
	if filters.Limit <= 0 || filters.Limit > 100 {
		filters.Limit = 20
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	return s.expenseRepo.GetAll(ctx, filters)
}

// Update actualiza un gasto con validaciones
func (s *expenseService) Update(ctx context.Context, id uuid.UUID, dto domain.UpdateExpenseDto, userID uuid.UUID) (*domain.Expense, error) {
	// Obtener gasto existente
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verificar ownership
	if expense.UserID != userID {
		return nil, sharedErrors.NewBusinessError("no_autorizado", "No tienes permiso para editar este gasto", "")
	}

	// Verificar que se puede editar (solo draft o rejected)
	if !expense.CanBeEdited() {
		return nil, sharedErrors.NewBusinessError("no_editable", "Solo se pueden editar gastos en estado 'draft' o 'rejected'", "")
	}

	// Aplicar cambios
	if dto.CategoryID != nil {
		// Verificar que la categoría existe y está activa
		category, err := s.categoryRepo.GetByID(ctx, *dto.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("error verificando categoría: %w", err)
		}
		if !category.IsActive {
			return nil, sharedErrors.NewValidationError("categoria_inactiva", "La categoría seleccionada no está activa")
		}
		expense.CategoryID = *dto.CategoryID
	}

	if dto.Title != nil {
		expense.Title = *dto.Title
	}

	if dto.Description != nil {
		expense.Description = dto.Description
	}

	if dto.Amount != nil {
		expense.Amount = *dto.Amount
		// Recalcular AmountCLP
		expense.AmountCLP = *dto.Amount * expense.ExchangeRate
	}

	if dto.Currency != nil {
		expense.Currency = *dto.Currency
		// TODO: Recalcular exchange rate
		expense.ExchangeRate = 1.0
		expense.AmountCLP = expense.Amount * expense.ExchangeRate
	}

	if dto.ExpenseDate != nil {
		if dto.ExpenseDate.After(time.Now()) {
			return nil, sharedErrors.NewValidationError("fecha_invalida", "La fecha del gasto no puede ser futura")
		}
		expense.ExpenseDate = *dto.ExpenseDate
	}

	if dto.MerchantName != nil {
		expense.MerchantName = dto.MerchantName
	}

	if dto.MerchantRUT != nil {
		// TODO: Validar formato de RUT
		expense.MerchantRUT = dto.MerchantRUT
	}

	if dto.ReceiptNumber != nil {
		expense.ReceiptNumber = dto.ReceiptNumber
	}

	if dto.PaymentMethod != nil {
		expense.PaymentMethod = *dto.PaymentMethod
	}

	if dto.IsReimbursable != nil {
		expense.IsReimbursable = *dto.IsReimbursable
	}

	expense.Updated = time.Now()

	// Actualizar en base de datos
	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return nil, fmt.Errorf("error actualizando gasto: %w", err)
	}

	// Retornar gasto actualizado
	return s.expenseRepo.GetByID(ctx, id)
}

// Delete elimina un gasto con validaciones
func (s *expenseService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Obtener gasto existente
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verificar ownership
	if expense.UserID != userID {
		return sharedErrors.NewBusinessError("no_autorizado", "No tienes permiso para eliminar este gasto", "")
	}

	// Verificar que se puede eliminar (solo draft o rejected)
	if !expense.CanBeDeleted() {
		return sharedErrors.NewBusinessError("no_eliminable", "Solo se pueden eliminar gastos en estado 'draft' o 'rejected'", "")
	}

	// Eliminar (soft delete)
	return s.expenseRepo.Delete(ctx, id)
}

// UploadReceipt sube un comprobante a un gasto
func (s *expenseService) UploadReceipt(ctx context.Context, dto domain.UploadReceiptDto, fileData []byte) (*domain.ExpenseReceipt, error) {
	// Verificar que el gasto existe
	expense, err := s.expenseRepo.GetByID(ctx, dto.ExpenseID)
	if err != nil {
		return nil, err
	}

	// Verificar que no está en estado reimbursed
	if expense.Status == domain.ExpenseStatusReimbursed {
		return nil, sharedErrors.NewBusinessError("gasto_reembolsado", "No se pueden agregar comprobantes a gastos reembolsados", "")
	}

	// TODO: Subir archivo a S3/MinIO
	// Por ahora usamos una URL placeholder
	fileURL := fmt.Sprintf("/receipts/%s/%s", dto.ExpenseID.String(), dto.FileName)

	// Crear el receipt
	receipt := &domain.ExpenseReceipt{
		ID:        uuid.New(),
		ExpenseID: dto.ExpenseID,
		FileURL:   fileURL,
		FileName:  dto.FileName,
		FileType:  dto.FileType,
		FileSize:  dto.FileSize,
		IsPrimary: dto.IsPrimary,
		Created:   time.Now(),
	}

	// Si es el primer receipt, forzar a primary
	existingReceipts, err := s.expenseRepo.GetReceipts(ctx, dto.ExpenseID)
	if err != nil {
		return nil, err
	}
	if len(existingReceipts) == 0 {
		receipt.IsPrimary = true
	}

	// Guardar receipt
	if err := s.expenseRepo.AddReceipt(ctx, receipt); err != nil {
		return nil, fmt.Errorf("error guardando comprobante: %w", err)
	}

	// Si es primary, desmarcar los demás
	if receipt.IsPrimary {
		if err := s.expenseRepo.SetPrimaryReceipt(ctx, receipt.ID); err != nil {
			return nil, err
		}
	}

	return receipt, nil
}

// DeleteReceipt elimina un comprobante
func (s *expenseService) DeleteReceipt(ctx context.Context, receiptID uuid.UUID, userID uuid.UUID) error {
	// Obtener receipts del gasto para verificar ownership
	// Primero necesitamos obtener el expense_id del receipt
	// Por simplicidad, asumimos que el receipt pertenece al usuario
	// En producción, se debería hacer un query para verificar

	// TODO: Verificar ownership del gasto antes de eliminar receipt

	// Eliminar receipt
	return s.expenseRepo.DeleteReceipt(ctx, receiptID)
}

// CanEdit verifica si un gasto puede ser editado
func (s *expenseService) CanEdit(ctx context.Context, expenseID uuid.UUID, userID uuid.UUID) (bool, error) {
	expense, err := s.expenseRepo.GetByID(ctx, expenseID)
	if err != nil {
		return false, err
	}

	// Verificar ownership
	if expense.UserID != userID {
		return false, nil
	}

	// Verificar estado
	return expense.CanBeEdited(), nil
}

// CanDelete verifica si un gasto puede ser eliminado
func (s *expenseService) CanDelete(ctx context.Context, expenseID uuid.UUID, userID uuid.UUID) (bool, error) {
	expense, err := s.expenseRepo.GetByID(ctx, expenseID)
	if err != nil {
		return false, err
	}

	// Verificar ownership
	if expense.UserID != userID {
		return false, nil
	}

	// Verificar estado
	return expense.CanBeDeleted(), nil
}

// ChangeStatus cambia el estado de un gasto
func (s *expenseService) ChangeStatus(ctx context.Context, expenseID uuid.UUID, newStatus domain.ExpenseStatus) error {
	expense, err := s.expenseRepo.GetByID(ctx, expenseID)
	if err != nil {
		return err
	}

	// Validar transiciones de estado permitidas
	if !s.isValidStatusTransition(expense.Status, newStatus) {
		return sharedErrors.NewBusinessError("transicion_invalida",
			fmt.Sprintf("No se puede cambiar de '%s' a '%s'", expense.Status, newStatus), "")
	}

	expense.Status = newStatus
	expense.Updated = time.Now()

	return s.expenseRepo.Update(ctx, expense)
}

// isValidStatusTransition valida las transiciones de estado permitidas
func (s *expenseService) isValidStatusTransition(current, new domain.ExpenseStatus) bool {
	// Matriz de transiciones permitidas
	validTransitions := map[domain.ExpenseStatus][]domain.ExpenseStatus{
		domain.ExpenseStatusDraft: {
			domain.ExpenseStatusSubmitted,
		},
		domain.ExpenseStatusSubmitted: {
			domain.ExpenseStatusApproved,
			domain.ExpenseStatusRejected,
		},
		domain.ExpenseStatusApproved: {
			domain.ExpenseStatusReimbursed,
		},
		domain.ExpenseStatusRejected: {
			domain.ExpenseStatusSubmitted, // Puede reenviar después de correcciones
		},
		domain.ExpenseStatusReimbursed: {
			// Estado final, no hay transiciones
		},
	}

	allowedStates, exists := validTransitions[current]
	if !exists {
		return false
	}

	for _, allowed := range allowedStates {
		if allowed == new {
			return true
		}
	}

	return false
}

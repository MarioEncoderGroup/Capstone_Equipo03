package controllers

import (
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/JoseLuis21/mv-backend/internal/core/ocr/ports"
	sharedErrors "github.com/JoseLuis21/mv-backend/internal/shared/errors"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
)

// OCRController maneja las operaciones de OCR
type OCRController struct {
	ocrService ports.OCRService
}

// NewOCRController crea una nueva instancia del controller de OCR
func NewOCRController(ocrService ports.OCRService) *OCRController {
	return &OCRController{
		ocrService: ocrService,
	}
}

// AnalyzeReceipt analiza un comprobante mediante OCR
// POST /api/v1/ocr/analyze
func (oc *OCRController) AnalyzeReceipt(c *fiber.Ctx) error {
	// 1. Obtener archivo del form-data
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "No se encontró el archivo en el request",
			Error:   "MISSING_FILE",
		})
	}

	// 2. Validar tipo de archivo
	contentType := file.Header.Get("Content-Type")
	allowedTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"application/pdf",
	}

	isValidType := false
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Tipo de archivo no permitido",
			Error:   "INVALID_FILE_TYPE",
			Data: fiber.Map{
				"allowed_types": allowedTypes,
			},
		})
	}

	// 3. Validar tamaño de archivo (max 10MB)
	maxSize := int64(10 * 1024 * 1024) // 10MB
	if file.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "El archivo es demasiado grande",
			Error:   "FILE_TOO_LARGE",
			Data: fiber.Map{
				"max_size_mb": 10,
				"file_size_mb": float64(file.Size) / (1024 * 1024),
			},
		})
	}

	// 4. Abrir archivo y leer contenido
	fileHandle, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error abriendo el archivo",
			Error:   "FILE_OPEN_ERROR",
		})
	}
	defer fileHandle.Close()

	imageData, err := io.ReadAll(fileHandle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error leyendo el archivo",
			Error:   "FILE_READ_ERROR",
		})
	}

	// 5. Procesar imagen con OCR Service
	result, err := oc.ocrService.ProcessReceipt(c.Context(), imageData)
	if err != nil {
		// Manejar errores del servicio
		if appErr, ok := sharedErrors.IsAppError(err); ok {
			return c.Status(appErr.HTTPCode).JSON(types.APIResponse{
				Success: false,
				Message: appErr.Message,
				Error:   appErr.Code,
				Data: fiber.Map{
					"details": appErr.Details,
				},
			})
		}

		// Error genérico
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error procesando el comprobante",
			Error:   "OCR_PROCESSING_ERROR",
			Data: fiber.Map{
				"details": err.Error(),
			},
		})
	}

	// 6. Retornar resultado exitoso
	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Comprobante analizado exitosamente",
		Data: fiber.Map{
			"amount":         result.Amount,
			"date":           result.Date,
			"merchant_rut":   result.MerchantRUT,
			"merchant_name":  result.MerchantName,
			"document_type":  result.DocumentType,
			"confidence":     result.Confidence,
			"raw_text":       result.RawText,
			"extracted_data": result.ExtractedData,
		},
	})
}

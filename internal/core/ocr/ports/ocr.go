package ports

import (
	"context"

	"github.com/JoseLuis21/mv-backend/internal/core/ocr/services"
)

// OCRService define el contrato para el servicio de OCR
type OCRService interface {
	// ProcessReceipt procesa imagen y retorna datos extraídos
	ProcessReceipt(ctx context.Context, imageData []byte) (*services.ParsedReceipt, error)

	// ProcessReceiptFromURL procesa desde URL (S3)
	ProcessReceiptFromURL(ctx context.Context, imageURL string) (*services.ParsedReceipt, error)
}

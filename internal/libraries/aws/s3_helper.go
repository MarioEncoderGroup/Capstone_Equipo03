package aws

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// S3Helper proporciona métodos de conveniencia para operaciones comunes con S3
type S3Helper struct {
	client *S3Client
}

// NewS3Helper crea una nueva instancia de S3Helper
func NewS3Helper(client *S3Client) *S3Helper {
	return &S3Helper{
		client: client,
	}
}

// UploadReceipt sube un comprobante de gasto a S3
// Organiza los archivos por usuario y fecha: receipts/{userID}/{year}/{month}/{filename}
func (h *S3Helper) UploadReceipt(ctx context.Context, userID uuid.UUID, filename string, data []byte, contentType string) (string, string, error) {
	// Generar nombre único de archivo
	ext := filepath.Ext(filename)
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Generar ruta organizacional
	now := time.Now()
	key := fmt.Sprintf("receipts/%s/%d/%02d/%s",
		userID.String(),
		now.Year(),
		now.Month(),
		uniqueFilename,
	)

	// Subir archivo
	url, err := h.client.UploadFile(ctx, key, data, contentType)
	if err != nil {
		return "", "", err
	}

	return url, key, nil
}

// DeleteReceipt elimina un comprobante usando su key
func (h *S3Helper) DeleteReceipt(ctx context.Context, key string) error {
	return h.client.DeleteFile(ctx, key)
}

// GetReceiptURL genera una URL temporal para ver un comprobante
// Default: 1 hora de validez
func (h *S3Helper) GetReceiptURL(ctx context.Context, key string) (string, error) {
	return h.client.GetPresignedURL(ctx, key, 1*time.Hour)
}

// GetReceiptURLWithDuration genera una URL temporal con duración personalizada
func (h *S3Helper) GetReceiptURLWithDuration(ctx context.Context, key string, duration time.Duration) (string, error) {
	return h.client.GetPresignedURL(ctx, key, duration)
}

// ValidateFileSize valida que el archivo no exceda el límite (10MB por defecto)
func (h *S3Helper) ValidateFileSize(data []byte, maxSizeMB int) error {
	if maxSizeMB == 0 {
		maxSizeMB = 10 // Default 10MB
	}

	maxBytes := int64(maxSizeMB * 1024 * 1024)
	if int64(len(data)) > maxBytes {
		return fmt.Errorf("el archivo excede el tamaño máximo de %dMB", maxSizeMB)
	}

	return nil
}

// ValidateContentType valida que el tipo de archivo sea permitido
func (h *S3Helper) ValidateContentType(contentType string) error {
	allowedTypes := map[string]bool{
		"image/jpeg":      true,
		"image/jpg":       true,
		"image/png":       true,
		"application/pdf": true,
	}

	if !allowedTypes[contentType] {
		return fmt.Errorf("tipo de archivo no permitido: %s. Solo se permiten JPG, PNG y PDF", contentType)
	}

	return nil
}

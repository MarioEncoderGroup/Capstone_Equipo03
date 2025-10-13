package aws

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Client maneja las operaciones con AWS S3 para MisViáticos
type S3Client struct {
	client *s3.Client
	bucket string
	region string
}

// S3Config contiene la configuración para el cliente S3
type S3Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
}

// NewS3Client crea una nueva instancia del cliente S3
func NewS3Client(cfg S3Config) (*S3Client, error) {
	// Validar configuración
	if cfg.Region == "" {
		return nil, fmt.Errorf("region es requerida")
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("access key ID es requerida")
	}
	if cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("secret access key es requerida")
	}
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket es requerido")
	}

	// Crear configuración AWS con credenciales estáticas
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error cargando configuración AWS: %w", err)
	}

	// Crear cliente S3
	client := s3.NewFromConfig(awsCfg)

	return &S3Client{
		client: client,
		bucket: cfg.Bucket,
		region: cfg.Region,
	}, nil
}

// UploadFile sube un archivo a S3 y retorna la URL pública
// key: ruta del archivo en S3 (ej: "receipts/user123/file.jpg")
// data: contenido del archivo en bytes
// contentType: tipo MIME del archivo (ej: "image/jpeg")
func (c *S3Client) UploadFile(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key no puede estar vacía")
	}
	if len(data) == 0 {
		return "", fmt.Errorf("data no puede estar vacía")
	}

	// Preparar input para upload
	input := &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        NewBytesReadSeekCloser(data),
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPrivate, // Archivo privado
	}

	// Subir archivo
	_, err := c.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("error subiendo archivo a S3: %w", err)
	}

	// Generar URL del archivo
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", c.bucket, c.region, key)

	return url, nil
}

// DeleteFile elimina un archivo de S3
func (c *S3Client) DeleteFile(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("key no puede estar vacía")
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	_, err := c.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("error eliminando archivo de S3: %w", err)
	}

	return nil
}

// GetPresignedURL genera una URL firmada para acceso temporal a un archivo
// duration: duración de validez de la URL
func (c *S3Client) GetPresignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key no puede estar vacía")
	}
	if duration <= 0 {
		return "", fmt.Errorf("duration debe ser mayor a 0")
	}

	// Crear presigner
	presignClient := s3.NewPresignClient(c.client)

	// Generar URL presignada
	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = duration
	})
	if err != nil {
		return "", fmt.Errorf("error generando URL presignada: %w", err)
	}

	return request.URL, nil
}

// FileExists verifica si un archivo existe en S3
func (c *S3Client) FileExists(ctx context.Context, key string) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key no puede estar vacía")
	}

	input := &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	_, err := c.client.HeadObject(ctx, input)
	if err != nil {
		// Si el error es "NotFound", el archivo no existe
		return false, nil
	}

	return true, nil
}

// GetFileSize obtiene el tamaño de un archivo en S3 (en bytes)
func (c *S3Client) GetFileSize(ctx context.Context, key string) (int64, error) {
	if key == "" {
		return 0, fmt.Errorf("key no puede estar vacía")
	}

	input := &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	result, err := c.client.HeadObject(ctx, input)
	if err != nil {
		return 0, fmt.Errorf("error obteniendo información del archivo: %w", err)
	}

	return aws.ToInt64(result.ContentLength), nil
}

// BytesReadSeekCloser implementa io.ReadSeekCloser para bytes
type BytesReadSeekCloser struct {
	*io.SectionReader
}

// NewBytesReadSeekCloser crea un nuevo BytesReadSeekCloser
func NewBytesReadSeekCloser(data []byte) *BytesReadSeekCloser {
	return &BytesReadSeekCloser{
		SectionReader: io.NewSectionReader(&bytesReaderAt{data: data}, 0, int64(len(data))),
	}
}

// Close implementa io.Closer
func (b *BytesReadSeekCloser) Close() error {
	return nil
}

// bytesReaderAt implementa io.ReaderAt
type bytesReaderAt struct {
	data []byte
}

// ReadAt implementa io.ReaderAt
func (b *bytesReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	if off >= int64(len(b.data)) {
		return 0, io.EOF
	}
	n = copy(p, b.data[off:])
	if n < len(p) {
		err = io.EOF
	}
	return
}
